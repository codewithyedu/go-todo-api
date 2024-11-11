package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

var (
	todos    = make(map[uuid.UUID]*Todo)
	validate = validator.New(validator.WithRequiredStructEnabled())
	mu       = &sync.RWMutex{}
)

type Todo struct {
	Title       string    `json:"title" validate:"required"`
	IsCompleted string    `json:"is_completed" validate:"required,oneof=true false"`
	ID          uuid.UUID `json:"id"`
}

func HandlerOne(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTodosHandler(w, r)
	case http.MethodPost:
		CreateTodoHandler(w, r)
	default:
		RespondWithErr(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func HandlerTwo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTodoHandler(w, r)
	case http.MethodPut:
		UpdateTodoHandler(w, r)
	case http.MethodDelete:
		DeleteTodoHandler(w, r)
	default:
		RespondWithErr(w, r, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithErr(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	RespondWithJSON(w, r, http.StatusOK, FT{"message": "OK"})
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	t := make([]*Todo, 0, len(todos))
	for _, todo := range todos {
		t = append(t, todo)
	}

	RespondWithJSON(w, r, http.StatusOK, FT{"todos": t})
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	t := &Todo{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(t); err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err := validate.Struct(t)
	if err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	t.ID = uuid.New()

	mu.Lock()
	defer mu.Unlock()

	todos[t.ID] = t
	RespondWithJSON(w, r, http.StatusCreated, formatTodo(t))
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	t, ok := todos[id]
	if !ok {
		RespondWithErr(w, r, http.StatusNotFound, fmt.Sprintf("todo with id '%v' not found", id))
		return
	}

	RespondWithJSON(w, r, http.StatusOK, formatTodo(t))
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	respT := &Todo{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(respT); err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err := validate.Struct(respT)
	if err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getID(r)
	if err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	t, ok := todos[id]
	if !ok {
		RespondWithErr(w, r, http.StatusNotFound, fmt.Sprintf("todo with id '%v' not found", id))
		return
	}

	t.Title = respT.Title
	t.IsCompleted = respT.IsCompleted

	RespondWithJSON(w, r, http.StatusOK, formatTodo(t))
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		RespondWithErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	t, ok := todos[id]
	if !ok {
		RespondWithErr(w, r, http.StatusNotFound, fmt.Sprintf("todo with id '%v' not found", id))
		return
	}

	delete(todos, t.ID)
	log.Printf("::: %v ::: %v %v", http.StatusNoContent, r.Method, r.URL)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	v1Router := http.NewServeMux()
	v1Router.HandleFunc("/api/v1/health", HealthHandler)
	v1Router.HandleFunc("/api/v1/todos", HandlerOne)
	v1Router.HandleFunc("/api/v1/todos/", HandlerTwo)

	log.Println("server is running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", v1Router))
}
