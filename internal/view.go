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

// HealthCheck	 godoc
// @Summary      Checking availability
// @Description  get availability status
// @Tags         users
// @Accept       nothing
// @Produce      nothing
// @Param        nothing
// @Success      200
// @Router       /healthcheck [get]

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Heathcheck")
	w.WriteHeader(http.StatusOK)
}

// CreateUser	 godoc
// @Summary      Create new user
// @Description  set new user
// @Tags         users
// @Accept       json
// @Produce      nothing
// @Param        nothing
// @Success      200  {object}  Handlers.User
// @Failure      500  Error
// @Router       /create [post]

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to create")
	b, _ := io.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if !isText(u) {
		http.Error(w, "invalid name and/or surname", http.StatusInternalServerError)
		return
	}
	if dGrades[u.Position] == "" {
		http.Error(w, "illegal position", http.StatusInternalServerError)
		return
	}

	err = h.dbc.AddUser(u.Name, u.Surname, u.Position, u.Project)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUser	 godoc
// @Summary      Delete user
// @Description  remove user
// @Tags         users
// @Accept       nothing
// @Produce      nothing
// @Param        id   int "User ID"
// @Success      200
// @Failure      400  Error
// @Failure      500  Error
// @Router       /delete/{id}	[delete]

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to delete", id)
	if id == "" {
		http.Error(w, "empty index", http.StatusBadRequest)
		return
	}
	val, _ := strconv.Atoi(id)

	err := h.dbc.DeleteUser(val)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "delete has been successful")
}

// UpdateUser	 godoc
// @Summary      Update user
// @Description  change user
// @Tags         users
// @Accept       json
// @Produce      nothing
// @Param        id   int "User ID"
// @Success      200  {object}  Handlers.User
// @Failure      400  Error
// @Failure      500  Error
// @Router       /update/{id}	[update]

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to update", id)

	if id == "" {
		http.Error(w, "empty index", http.StatusBadRequest)
		return
	}
	val, _ := strconv.Atoi(id)
	b, _ := io.ReadAll(r.Body)
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	if !isText(u) {
		http.Error(w, "invalid name and/or surname", http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "update has been successful")
}

// GetUser	 	 godoc
// @Summary      Get user
// @Description  get user by id
// @Tags         users
// @Accept       nothing
// @Produce      json
// @Param        id				int	"User ID"
// @Success      200  {object}	Handlers.User
// @Failure      500  {object}	Error
// @Router       /get/{id}		[get]

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("Trying to get", id)
	val, _ := strconv.Atoi(id)

	u, err := h.dbc.GetUser(val)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	content, err := json.Marshal(u)
	if err != nil {
		log.Println("Marshalling error")
		http.Error(w, "marshalling error", http.StatusInternalServerError)
		return
	}
	r.Header.Set("Content-Type", "application/json")
	w.Write(content)
	/*fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s",
	u.Name, u.Surname, dGrades[u.Position], u.Project)

	*/
}

// GetUserList	 godoc
// @Summary      List users
// @Description  get users
// @Tags         users
// @Accept       nothing
// @Produce      json
// @Param        nothing
// @Success      200  {array}   Handlers.User
// @Failure      500  {object}	Error
// @Router       /getall [get]

func (h *Handlers) GetUserList(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to get all user list")
	us, err := h.dbc.GetAllUsers()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	content, err := json.Marshal(*us)
	if err != nil {
		log.Println("Marshalling error")
		http.Error(w, "marshalling error", http.StatusInternalServerError)
		return
	}
	r.Header.Set("Content-Type", "application/json")
	w.Write(content)

	/*for _, u := range *us {
		json.Marshal(u)
		fmt.Fprintf(w, "name: %s, surname: %s, position: %s, project: %s\r\n",
			u.Name, u.Surname, dGrades[u.Position], u.Project)
	}

	*/
}