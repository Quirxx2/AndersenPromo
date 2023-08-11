package promo

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

type Handlers struct {
	dbc DBConnexion
}

func NewHandlers(connString string) (*Handlers, error) {
	r, err := NewRegistry(connString)
	if err != nil {
		log.Fatalf("Failed to create Registry: %v", err)
	}
	return &Handlers{r}, nil
}

func isText(u User) bool {
	re := regexp.MustCompile(`[a-zA-Z]`)
	return re.MatchString(u.Name) && re.MatchString(u.Surname)
}

func (h *Handlers) healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Heathcheck")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to create")
	b, _ := io.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		fmt.Fprintf(w, "mandatory data is absent, error occured: %w", err)
		return
	}

	if !isText(u) {
		fmt.Fprintf(w, "invalid name and/or surname")
		return
	}
	if dGrades[u.Position] == "" {
		fmt.Fprintf(w, "illegal position", err)
		return
	}

	err = h.dbc.AddUser(u.Name, u.Surname, u.Position, u.Project)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to delete", id)
	if id == "" {
		fmt.Fprintf(w, "empty index")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, _ := strconv.Atoi(id)

	err := h.dbc.DeleteUser(val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "delete has been successful")
}

func (h *Handlers) updateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to update", id)

	if id == "" {
		fmt.Fprintf(w, "empty index")
		return
	}
	val, _ := strconv.Atoi(id)
	b, _ := io.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		fmt.Fprintf(w, "mandatory data is absent, error occured: %w", err)
		return
	}

	if !isText(u) {
		fmt.Fprintf(w, "invalid name and/or surname")
		return
	}
	m := make(map[string]string)
	if u.Name != "" {
		m["name"] = u.Name
	}
	if u.Surname != "" {
		m["surname"] = u.Surname
	}
	if dGrades[u.Position] != "" {
		m["position"] = dGrades[u.Position]
	}
	if u.Name != "" {
		m["project"] = u.Project
	}
	err = h.dbc.UpdateUser(val, m)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "update has been successful")
}

func (h *Handlers) getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to get", id)
	val, _ := strconv.Atoi(id)

	u, err := h.dbc.GetUser(val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s",
		u.Name, u.Surname, dGrades[u.Position], u.Project)
}

func (h *Handlers) getUserList(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to get all user list")
	us, err := h.dbc.GetAllUsers()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	for _, u := range *us {
		json.Marshal(u)
		fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s\r\n",
			u.Name, u.Surname, dGrades[u.Position], u.Project)
	}
}
