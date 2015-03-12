package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func NewStore(col *db.Col) (Store, error) {
	var store Store = Store{Col: col}
	return store, nil
}

type Store struct {
	Col *db.Col
}

func (s *Store) getUidFromId(id string) (int, error) {
	var query interface{}
	json.Unmarshal([]byte(fmt.Sprintf(`[{"eq": "%s", "in": ["ID"], "limit": 1}]`, id)), &query)
	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(query, s.Col, &queryResult); err != nil {
		return 0, err
	}
	for uid := range queryResult {
		return uid, nil
	}
	return 0, errors.New("Record not found")
}

func (s *Store) One(id string) (Node, error) {
	var node Node = Node{}
	uid, err := s.getUidFromId(id)
	if err != nil {
		docContent, err := s.Col.Read(uid)
		if err != nil {
			return node, err
		}
		if err = mapstructure.Decode(docContent, &node); err != nil {
			return node, err
		}
		return node, nil
	}
	return node, err
}

func (s *Store) All() ([]Node, error) {
	var allNodes []Node = make([]Node, 0)
	s.Col.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var doc Node = Node{}
		if err := json.Unmarshal(docContent, &doc); err != nil {
			log.Fatal(err)
		}
		allNodes = append(allNodes, doc)
		return true
	})
	return allNodes, nil
}

func (s *Store) Insert(node Node) error {
	if node.ID == "" {
		return errors.New("Missing ID")
	}
	_, err := s.One(node.ID)
	if err == nil {
		return errors.New("Dublicate ID")
	}
	_, err = s.Col.Insert(structs.Map(node))
	return err
}

func (s *Store) Upsert(node Node) error {
	if node.ID == "" {
		return errors.New("Missing ID")
	}
	uid, err := s.getUidFromId(node.ID)
	if err == nil {
		err = s.Col.Update(uid, structs.Map(node))
	}
	return err
}
