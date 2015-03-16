package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"

	tiedotdb "github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func NewStore(col *tiedotdb.Col) (*Store, error) {
	var store Store = Store{Col: col}
	return &store, nil
}

type Store struct {
	Col *tiedotdb.Col
	sync.RWMutex
}

func (s *Store) getUidFromId(id string) (int, error) {
	var query interface{}
	json.Unmarshal([]byte(fmt.Sprintf(`[{"eq": "%s", "in": ["ID"], "limit": 1}]`, id)), &query)
	queryResult := make(map[int]struct{})
	if err := tiedotdb.EvalQuery(query, s.Col, &queryResult); err != nil {
		return 0, err
	}
	for uid := range queryResult {
		return uid, nil
	}
	return 0, errors.New("Record not found")
}

func (s *Store) One(doc ModelInterface, id string) error {
	s.Lock()
	defer s.Unlock()

	uid, err := s.getUidFromId(id)
	if err == nil {
		docContent, err := s.Col.Read(uid)
		if err != nil {
			return err
		}
		if err = mapstructure.Decode(docContent, doc); err != nil {
			return err
		}
		return nil
	}
	return err
}

func (s *Store) All(sliceptr interface{}) error {
	s.Lock()
	defer s.Unlock()

	// http://play.golang.org/p/Hn7kLc_Qm6
	// http://play.golang.org/p/IgTRevpxEP !!!

	// TODO Check if is struct etc

	sv := reflect.ValueOf(sliceptr).Elem()
	et := sv.Type().Elem()

	s.Col.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		tmp := reflect.New(et).Interface()
		if err := json.Unmarshal(docContent, tmp); err != nil {
			log.Fatal(err)
		}
		tmpv := reflect.ValueOf(tmp).Elem()
		sv.Set(reflect.Append(sv, tmpv))
		return true
	})

	return nil
}

func (s *Store) Insert(doc ModelInterface) error {
	s.Lock()
	defer s.Unlock()

	if doc.String() == "" {
		return errors.New("Missing ID")
	}
	_, err := s.getUidFromId(doc.String())
	if err == nil {
		return errors.New("Dublicate ID")
	}
	_, err = s.Col.Insert(structs.Map(doc))
	return err
}

func (s *Store) Upsert(doc ModelInterface) error {
	s.Lock()
	defer s.Unlock()

	if doc.String() == "" {
		return errors.New("Missing ID")
	}
	uid, err := s.getUidFromId(doc.String())
	if err == nil {
		err = s.Col.Update(uid, structs.Map(doc))
	}
	return err
}
