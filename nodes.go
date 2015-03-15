package main

import (
	"sync"

	tiedotdb "github.com/HouzuoGuo/tiedot/db"
)

type NodeStore struct {
	Store
	sync.RWMutex
}

func NewNodeStore(col *tiedotdb.Col) (*NodeStore, error) {
	var nodes NodeStore = NodeStore{}
	nodes.Col = col
	return &nodes, nil
}
