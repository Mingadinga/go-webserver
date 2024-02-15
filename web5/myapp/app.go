package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}
	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}
	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// parse path variable "id"
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// find user from map
	foundUser, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	// make response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(foundUser)
	fmt.Fprint(w, string(data))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// parse user info from request
	newUser := new(User)
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, err)
		return
	}

	// create user and save
	lastID++
	newUser.ID = lastID
	newUser.CreatedAt = time.Now()
	userMap[newUser.ID] = newUser

	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(newUser)
	fmt.Fprint(w, string(data))
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// parse path variable "id"
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	_, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User Id:", id)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {

	updateUser := new(UpdateUser)
	err := json.NewDecoder(r.Body).Decode(updateUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[updateUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", updateUser.ID)
	}

	if updateUser.UpdatedEmail {
		user.Email = updateUser.Email
	}
	if updateUser.UpdatedFirstName {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.UpdatedLastName {
		user.LastName = updateUser.LastName
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

var userMap map[int]*User
var lastID int

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUser struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	UpdatedFirstName bool      `json:"updated_first_name"`
	LastName         string    `json:"last_name"`
	UpdatedLastName  bool      `json:"updated_last_name"`
	Email            string    `json:"email"`
	UpdatedEmail     bool      `json:"updated_email"`
	CreatedAt        time.Time `json:"created_at"`
}

// NewHandler New Handler make a new myapp handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0

	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")

	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")

	return mux
}
