package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", Index).Name("Index")
	router.HandleFunc("/nodes", NodeIndex).Name("NodeIndex")
	// router.HandleFunc("/nodes", TodoCreate).Methods("POST").Name("TodoCreate")
	router.HandleFunc("/nodes/{nodeId}", NodeShow).Name("TodoShow")

	return router
}
