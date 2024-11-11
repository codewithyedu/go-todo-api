package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestUpdateTodo(t *testing.T) {
	for i := 0; i < numReq; i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			title := fmt.Sprintf("test todo %v", i)
			isCompleted := "false"
			todo, err := CreateTodo(title, isCompleted, t)
			if err != nil {
				t.Error(err)
				return
			}
			id := todo.ID.String()

			// Test valid id and body
			title = fmt.Sprintf("testing update todo %v", i)
			isCompleted = "true"
			updateBody := fmt.Sprintf(`{ "title": "%v", "is_completed": "%v" }`, title, isCompleted)
			nr := strings.NewReader(updateBody)

			req := httptest.NewRequest(http.MethodPut, url+"/"+id, nr)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing req body: %v", err)
					return
				}
			}(req.Body)
			rr := httptest.NewRecorder()
			UpdateTodoHandler(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected response code %v, received %v", http.StatusOK, rr.Code)
				return
			}

			todo = &TestTodo{}
			dec := json.NewDecoder(rr.Body)
			if err := dec.Decode(todo); err != nil {
				t.Errorf("error decoding marshal: %v", err)
				return
			}

			if todo.Title != title {
				t.Errorf("expected title %v, received %v", title, todo.Title)
				return
			}

			if todo.IsCompleted != isCompleted {
				t.Errorf("expected completion status %v, received %v", isCompleted, todo.IsCompleted)
				return
			}

			// Test malformed body
			req = httptest.NewRequest(http.MethodPut, url+"/"+id, nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing req body: %v", err)
					return
				}
			}(req.Body)
			rr = httptest.NewRecorder()

			UpdateTodoHandler(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected response code %v, received %v", http.StatusBadRequest, rr.Code)
				return
			}

			// Test invalid ID
			req = httptest.NewRequest(http.MethodPut, url+"/invalid-id", nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing req body: %v", err)
					return
				}
			}(req.Body)
			rr = httptest.NewRecorder()

			UpdateTodoHandler(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected response code %v, received %v", http.StatusBadRequest, rr.Code)
				return
			}

		}(i, wg)
	}

	wg.Wait()
}
