package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_FILE         = "DBFILE.db"
	LISTEN_PORT     = ":8080"
	ErrAlreadyExist = "item already exist"
)

type Server struct {
	db *gorm.DB
}

type Todo struct {
    ID       uint64 `json:"id"`
    TodoItem string `json:"todoItem"`
    Complete bool   `json:"complete"`
}

func (db *Server) GetTodoItemsList(w http.ResponseWriter, r *http.Request) {

    var todoList []Todo
	
	if result := db.db.Find(&todoList) ;result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(todoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (db *Server) GetTodoItem(w http.ResponseWriter, r *http.Request) {

	var todo Todo
	vars := mux.Vars(r)
	id := vars["id"]
	
	if result := db.db.First(&todo, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (db *Server) AddTodoItem(w http.ResponseWriter, r *http.Request) {

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	if err := validateTaskFields(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	if result := db.db.First(&todo, todo.ID); result.Error == nil {
		http.Error(w, ErrAlreadyExist, http.StatusConflict)
		return
	}

	if result := db.db.Create(&todo); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return 
	}

	data, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

}

func (db *Server) UpdateTodoItem(w http.ResponseWriter, r *http.Request) {

	var newTodo Todo

	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	var oldTodo Todo
	if result := db.db.First(&oldTodo, newTodo.ID); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusConflict)
		return
	}

	if result := db.db.Save(&newTodo); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return 
	}

	data, err := json.Marshal(newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(data)

}

func (db *Server) DeleteTodoItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var todo Todo
	if result := db.db.First(&todo, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	if result := db.db.Delete(&todo); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusAccepted)

}

func validateTaskFields(t *Todo) error {

    if t.ID == 0 || len(t.TodoItem) == 0 {
        return errors.New(fmt.Sprintln(http.StatusBadRequest))
    }
    return nil

}

func main () {

	s := Server{}
	var err error
	
	if s.db, err = gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{}); err != nil {
		fmt.Println("Error")
		return 
	}
	s.db.AutoMigrate(&Todo{})

	router := mux.NewRouter()
	router.HandleFunc("/todolist", s.GetTodoItemsList).Methods("GET")
	router.HandleFunc("/todolist", s.AddTodoItem).Methods("POST")
	router.HandleFunc("/todolist", s.UpdateTodoItem).Methods("PATCH")
	router.HandleFunc("/todolist/{id}", s.GetTodoItem).Methods("GET")
	router.HandleFunc("/todolist/{id}", s.DeleteTodoItem).Methods("DELETE")

	http.ListenAndServe(LISTEN_PORT, router)

}






