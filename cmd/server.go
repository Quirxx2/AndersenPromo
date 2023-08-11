package main

import (
	"github.com/Quirxx2/AndersenPromo/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var reg *Registry

func main() {
	serverHost := ":8080"
	connString := "postgres://usr:password@localhost:5432/registry"

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/create", createUser).Methods(http.MethodPost)
	router.HandleFunc("/delete/{id}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/update/{id}", updateUser).Methods(http.MethodPatch)
	router.HandleFunc("/get/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/getall", getUserList).Methods(http.MethodGet)

	r, err := NewRegistry(connString)
	if err != nil {
		log.Fatalf("Failed to create Registry: %v", err)
	}
	reg = r
	log.Println(r)

	srv := &http.Server{
		Handler: router,
		Addr:    serverHost,
	}

	log.Fatal(srv.ListenAndServe())
}
