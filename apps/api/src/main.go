// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoStore struct {
	mu      sync.RWMutex
	todos   map[int]Todo
	counter int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:   make(map[int]Todo),
		counter: 1,
	}
}

func (ts *TodoStore) Create(todo Todo) int {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	todo.ID = ts.counter
	ts.todos[todo.ID] = todo
	ts.counter++
	return todo.ID
}

func (ts *TodoStore) Read(id int) (Todo, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	todo, exists := ts.todos[id]
	return todo, exists
}

func (ts *TodoStore) Update(id int, updatedTodo Todo) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.todos[id]; !exists {
		return false
	}

	updatedTodo.ID = id
	ts.todos[id] = updatedTodo
	return true
}

func (ts *TodoStore) Delete(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.todos[id]; !exists {
		return false
	}

	delete(ts.todos, id)
	return true
}

func (ts *TodoStore) List() []Todo {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	todos := make([]Todo, 0, len(ts.todos))
	for _, todo := range ts.todos {
		todos = append(todos, todo)
	}
	return todos
}

var todoStore *TodoStore

func init() {
	todoStore = NewTodoStore()
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todoID := todoStore.Create(todo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{\"id\": %d}", todoID)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	todo, exists := todoStore.Read(id)
	if !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !todoStore.Update(id, updatedTodo) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	if !todoStore.Delete(id) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos := todoStore.List()
	json.NewEncoder(w).Encode(todos)
}

func main() {
	http.HandleFunc("/api/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createTodoHandler(w, r)
		case http.MethodGet:
			getTodoHandler(w, r)
		case http.MethodPut:
			updateTodoHandler(w, r)
		case http.MethodDelete:
			deleteTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listTodosHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server listening on port 33808")
	http.ListenAndServe(":33808", nil)
}
