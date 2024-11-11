package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetTodo(t *testing.T) {

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

			req := httptest.NewRequest(http.MethodGet, url+"/"+id, nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing request body: %v", err)
					return
				}
			}(req.Body)
			rr := httptest.NewRecorder()
			GetTodoHandler(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected resposnse code %v, received %v", http.StatusOK, rr.Code)
				return
			}

			respTodo := &TestTodo{}
			dec := json.NewDecoder(rr.Body)
			if err = dec.Decode(respTodo); err != nil {
				t.Errorf("error decoding json: %v", err)
				return
			}

			if todo.ID != respTodo.ID {
				t.Errorf("expected ID %v, received %v", todo.ID, respTodo.ID)
				return
			}

			if todo.Title != respTodo.Title {
				t.Errorf("expected title %v, received %v", todo.Title, respTodo.Title)
				return
			}

			if todo.IsCompleted != respTodo.IsCompleted {
				t.Errorf("expected completion status %v, received %v", todo.IsCompleted, respTodo.IsCompleted)
				return
			}

			req = httptest.NewRequest(http.MethodGet, url+"/invalid-id", nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing request body: %v", err)
					return
				}
			}(req.Body)
			rr = httptest.NewRecorder()
			GetTodoHandler(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected response code %v, received %v", http.StatusBadRequest, rr.Code)
				return
			}

		}(i, wg)
	}

	wg.Wait()
}
