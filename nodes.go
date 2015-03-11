package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

var currentId int

var Nodes NodesRepo = NodesRepo{}
var myDB db.DB

func init() {
	myDBDir := "/tmp/tiedot_test_embeddedExample"
	os.RemoveAll(myDBDir)
	defer os.RemoveAll(myDBDir)

	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := myDB.Create("Nodes"); err != nil {
		log.Fatal(err)
	}
	Nodes.Col = myDB.Use("Nodes")
	if err := Nodes.Col.Index([]string{"ID"}); err != nil {
		panic(err)
	}
	Nodes.Insert(Node{
		ID:       "a1",
		Hostname: "http://1",
	})
	Nodes.Insert(Node{
		ID:       "a2",
		Hostname: "http://xx",
	})
	Nodes.Upsert(Node{
		ID:       "a2",
		Hostname: "http://2",
	})
}

type NodesRepo struct {
	DB  db.DB
	Col *db.Col
}

func (s *NodesRepo) getUidFromId(id string) (int, error) {
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

func (s *NodesRepo) One(id string) (Node, error) {
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

func (s *NodesRepo) All() ([]Node, error) {
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

func (s *NodesRepo) Insert(node Node) error {
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

func (s *NodesRepo) Upsert(node Node) error {
	if node.ID == "" {
		return errors.New("Missing ID")
	}
	uid, err := s.getUidFromId(node.ID)
	if err == nil {
		err = s.Col.Update(uid, structs.Map(node))
	}
	return err
}
