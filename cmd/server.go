package main

import (
	_ "github.com/Quirxx2/AndersenPromo/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	serverHost := ":8080"
	connString := "postgres://usr:password@localhost:5432/registry"

	c, err := NewHandlers(connString)
	if err != nil {
		log.Fatalf("Failed to create Handlers: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", c.healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/create", c.createUser).Methods(http.MethodPost)
	router.HandleFunc("/delete/{id}", c.deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/update/{id}", c.updateUser).Methods(http.MethodPatch)
	router.HandleFunc("/get/{id}", c.getUser).Methods(http.MethodGet)
	router.HandleFunc("/getall", c.getUserList).Methods(http.MethodGet)

	srv := &http.Server{
		Handler: router,
		Addr:    serverHost,
	}

	log.Fatal(srv.ListenAndServe())
}
