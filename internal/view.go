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
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return re.MatchString(u.Name) && re.MatchString(u.Surname)
}

// HealthCheck	 godoc
//
//	@Summary		Checking availability
//	@Description	get availability status
//	@Tags			users
//	@Success		200
//	@Router			/healthcheck [get]
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	//log.Println("Healthcheck")
	w.WriteHeader(http.StatusOK)
}

// CreateUser	 godoc
//
//	@Summary		Create new user
//	@Description	set new user
//	@Tags			users
//	@Accept			json
//	@Success		200	{object}	User
//	@Failure		500	{object}	string
//	@Router			/create [post]
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	//log.Println("Trying to create")
	b, _ := io.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	if !isText(u) {
		http.Error(w, "invalid name and/or surname", http.StatusBadRequest)
		return
	}
	if dGrades[u.Position] == "" {
		http.Error(w, "illegal position", http.StatusBadRequest)
		return
	}

	err = h.dbc.AddUser(u.Name, u.Surname, u.Position, u.Project)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUser	 godoc
//
//	@Summary		Delete user
//	@Description	remove user
//	@Tags			users
//	@Param			id	path	int	true	"User ID"
//	@Success		200
//	@Failure		400				{object}	string
//	@Failure		500				{object}	string
//	@Router			/delete/{id}	[delete]
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//log.Println("Trying to delete", id)
	if id == "" {
		http.Error(w, "empty index", http.StatusBadRequest)
		return
	}
	val, err := strconv.Atoi(id)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	err = h.dbc.DeleteUser(val)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateUser	 godoc
//
//	@Summary		Update user
//	@Description	change user
//	@Tags			users
//	@Accept			json
//	@Param			id				path		int	true	"User ID"
//	@Success		200				{object}	User
//	@Failure		400				{object}	string
//	@Failure		500				{object}	string
//	@Router			/update/{id}	[patch]
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//log.Println("Trying to update", id)

	if id == "" {
		http.Error(w, "empty index", http.StatusBadRequest)
		return
	}
	val, err := strconv.Atoi(id)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	b, _ := io.ReadAll(r.Body)
	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	if !isText(u) {
		http.Error(w, "invalid name and/or surname", http.StatusBadRequest)
		return
	}
	m := make(map[string]string)
	if u.Name != "" {
		m["name"] = u.Name
	}
	if u.Surname != "" {
		m["surname"] = u.Surname
	}
	//log.Println("dGrades[u.Position] =", dGrades[u.Position])
	if dGrades[u.Position] != "" {
		m["position"] = dGrades[u.Position]
	}
	if u.Name != "" {
		m["project"] = u.Project
	}
	err = h.dbc.UpdateUser(val, m)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUser	 	 godoc
//
//	@Summary		Get user
//	@Description	get user by id
//	@Tags			users
//	@Produce		json
//	@Param			id			path		int	true	"User ID"
//	@Success		200			{object}	User
//	@Failure		400			{object}	string
//	@Failure		500			{object}	string
//	@Router			/get/{id}																						[get]
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//log.Println("Trying to get", id)
	if id == "" {
		http.Error(w, "empty index", http.StatusBadRequest)
		return
	}
	val, err := strconv.Atoi(id)
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	u, err := h.dbc.GetUser(val)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	content, _ := json.Marshal(u)
	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

// GetUserList	 godoc
//
//	@Summary		List users
//	@Description	get users
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		User
//	@Failure		500	{object}	string
//	@Router			/getall [get]
func (h *Handlers) GetUserList(w http.ResponseWriter, r *http.Request) {
	// log.Println("Trying to get all user list")
	us, err := h.dbc.GetAllUsers()
	if err != nil {
		//log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	content, _ := json.Marshal(*us)
	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}
