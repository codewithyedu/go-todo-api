package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestCreateTodo(t *testing.T) {

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

			if todo.Title != title {
				t.Errorf("expected  title %s, recieved %s", title, todo.Title)
				return
			}

			if todo.IsCompleted != isCompleted {
				t.Errorf("expected completion status %s, recieved %s", isCompleted, todo.IsCompleted)
				return
			}
		}(i, wg)
	}

	wg.Wait()
}
