package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Heathcheck")
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to create")
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

	err = reg.AddUser(u.Name, u.Surname, u.Position, u.Project)
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to delete", id)
	if id == "" {
		fmt.Fprintf(w, "empty body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(w, "requested numeric format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = reg.DeleteUser(val)
	if err != nil {
		fmt.Fprintf(w, "%w", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "success")
	w.WriteHeader(http.StatusOK)
}
