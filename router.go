package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	subrouter := router.Host("{fqdn:[a-z\\.]+}").Subrouter()
	subrouter.HandleFunc("/", Index).Name("Index")
	subrouter.HandleFunc("/nodes", NodeIndex).Name("NodeIndex")
	// subrouter.HandleFunc("/nodes", TodoCreate).Methods("POST").Name("NodeCreate")
	subrouter.HandleFunc("/nodes/{nodeId}", NodeShow).Name("NodeShow")

	return router
}
