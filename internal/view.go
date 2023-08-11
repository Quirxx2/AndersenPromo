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

func isText(u User) bool {
	re := regexp.MustCompile(`[a-zA-Z]`)
	return re.MatchString(u.Name) && re.MatchString(u.Surname)
}

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

	err = reg.AddUser(u.Name, u.Surname, u.Position, u.Project)
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to delete", id)
	if id == "" {
		fmt.Fprintf(w, "empty index")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, _ := strconv.Atoi(id)

	err := reg.DeleteUser(val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "delete has been successful")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
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
	err = reg.UpdateUser(val, m)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "update has been successful")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to get", id)
	val, _ := strconv.Atoi(id)

	u, err := reg.GetUser(val)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s",
		u.Name, u.Surname, dGrades[u.Position], u.Project)
}

func getUserList(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to get all user list")
	us, err := reg.GetAllUsers()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err)
		return
	}
	for _, u := range *us {
		fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s\r\n",
			u.Name, u.Surname, dGrades[u.Position], u.Project)
	}
}
