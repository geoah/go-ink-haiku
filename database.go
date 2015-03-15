package main

import (
	"fmt"
	"log"
	"os"

	tiedotdb "github.com/HouzuoGuo/tiedot/db"
)

var db DB = NewDB()

type DB struct {
	db    *tiedotdb.DB
	Nodes *NodeStore
	// Identities tiedotdb.Col
	// Instances tiedotdb.Col
	// Schemas tiedotdb.Col
}

func NewDB() DB {
	var dbPath = fmt.Sprintf("/tmp/%s", "ink-v1")
	os.RemoveAll(dbPath)
	defer os.RemoveAll(dbPath)

	tdb, err := tiedotdb.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := tdb.Create("Nodes"); err != nil {
		log.Fatal(err)
	}

	nodesCol := tdb.Use("Nodes")
	if err := nodesCol.Index([]string{"ID"}); err != nil {
		panic(err)
	}
	nodeStore, err := NewNodeStore(nodesCol)
	if err != nil {
		panic(err)
	}

	var idb DB = DB{db: tdb, Nodes: nodeStore}

	nodeStore.Insert(Node{
		ID:       "a1",
		Hostname: "http://1",
	})
	nodeStore.Insert(Node{
		ID:       "a2",
		Hostname: "http://xx",
	})
	nodeStore.Upsert(Node{
		ID:       "a2",
		Hostname: "http://2",
	})

	return idb
}
