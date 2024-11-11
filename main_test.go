package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	wg     = &sync.WaitGroup{}
	numReq = 100
	url    string
)

type TestTodo struct {
	Title       string `json:"title"`
	IsCompleted string `json:"is_completed"`
	ID          uuid.UUID
}

func TestMain(m *testing.M) {
	url = "http://localhost:8080/api/v1/todos"
	code := m.Run()
	os.Exit(code)
}

func CreateTodo(title, isCompleted string, t *testing.T) (*TestTodo, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	body := strings.NewReader(
		fmt.Sprintf(`{ "title": "%s", "is_completed": "%s" }`, title, isCompleted),
	)
	req := httptest.NewRequest(http.MethodPost, url, body)
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.Errorf("error closing request body: %v", err)
			return
		}
	}(req.Body)
	rr := httptest.NewRecorder()
	CreateTodoHandler(rr, req)

	todo := &TestTodo{}
	dec := json.NewDecoder(rr.Body)
	if err := dec.Decode(todo); err != nil {
		return nil, fmt.Errorf("error decoding marshal: %v", err)
	}

	return todo, nil
}
