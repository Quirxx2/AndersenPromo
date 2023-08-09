package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var reg *Registry

func main() {
	serverHost := ":8080"
	//httpHost := "http://localhost:8080/"
	connString := "postgres://usr:password@localhost:5432/registry"

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/createuser", createUser).Methods(http.MethodPost)
	router.HandleFunc("/deleteuser/{id}", deleteUser).Methods(http.MethodDelete)

	r, err := NewRegistry(connString)
	if err != nil {
		log.Fatalf("Failed to create Registry: %v", err)
	}
	reg = r

	srv := &http.Server{
		Handler: router,
		Addr:    serverHost,
	}

	log.Fatal(srv.ListenAndServe())
}
