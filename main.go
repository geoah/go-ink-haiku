package main

import (
	"log"
	"net/http"

	"github.com/carbocation/interpose"
)

// const ContextIdentity key = 0

func main() {
	router := NewRouter()

	middle := interpose.New()
	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// context.Set(req, ContextHostname, Get)
		rw.Header().Set("X-Server-Name", "Interpose Test Server")
	}))
	middle.UseHandler(router)

	log.Fatal(http.ListenAndServe(":8080", middle))
}
