package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetTodos(t *testing.T) {
	for i := 0; i < numReq; i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodGet, url, nil)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Errorf("error closing request body: %v", err)
					return
				}
			}(req.Body)
			rr := httptest.NewRecorder()
			GetTodosHandler(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected response code %v status code, received %v", http.StatusOK, rr.Code)
				return
			}
		}(i, wg)
	}

	wg.Wait()
}
