package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Grade int

const (
	trainee Grade = iota + 1
	junior
	middle
	senior
)

var grades = map[Grade]string{
	trainee: "trainee",
	junior:  "junior",
	middle:  "middle",
	senior:  "senior",
}

func main() {
	serverHost := ":8080"
	httpHost := "http://localhost:8080/"
	connString := "postgres://user:password@localhost:5432/registry"

	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.HandleFunc("/createuser", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		var u User
		err := json.Unmarshal(b, &u)
		if err != nil {
			fmt.Fprintf(w, "mandatory data is absent, error occured: %w", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		re := regexp.MustCompile(`[a-zA-Z]`)
		isText := re.MatchString(u.Name) && re.MatchString(u.Surname)
		if !isText {
			fmt.Fprintf(w, "invalid name and/or surname")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if grades[u.Position] == "" {
			fmt.Fprintf(w, "illegal position", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	c, err := NewConnexion(connString)
	if err != nil {
		log.Fatalf("Failed to create Connexion: %v", err)
	}

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
