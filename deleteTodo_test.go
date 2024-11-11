package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteTodo(t *testing.T) {
	for i := 0; i < numReq; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			title := fmt.Sprintf("test todo %v", i)
			isCompleted := "false"
			todo, err := CreateTodo(title, isCompleted, t)
			if err != nil {
				t.Error(err)
				return
			}
			id := todo.ID.String()

			req := httptest.NewRequest(http.MethodDelete, url+"/"+id, nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing request body: %v", err)
					return
				}
			}(req.Body)
			rr := httptest.NewRecorder()
			DeleteTodoHandler(rr, req)

			if rr.Code != http.StatusNoContent {
				t.Errorf("expected response code %v, received %v", http.StatusNoContent, rr.Code)
				return
			}

			req = httptest.NewRequest(http.MethodDelete, url+"/"+id, nil)
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Errorf("error closing request body: %v", err)
					return
				}
			}(req.Body)
			rr = httptest.NewRecorder()
			DeleteTodoHandler(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Errorf("expected response code %v, received %v", http.StatusNoContent, rr.Code)
				return
			}
		}()
	}

	wg.Wait()
}
