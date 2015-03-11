package main

import (
	"log"
	"os"

	"github.com/HouzuoGuo/tiedot/db"
)

var currentId int

var Nodes Store
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

	Nodes, err := NewStore(myDB.Use("Nodes"))
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
