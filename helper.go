package main

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type FT map[string]interface{}

func formatTodo(t *Todo) *FT {
	return &FT{
		"id":           t.ID,
		"title":        t.Title,
		"is_completed": t.IsCompleted,
	}
}

func getID(r *http.Request) (uuid.UUID, error) {
	// Remove the "/api/v1/todos/" prefix to get the UUID part
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
	if idStr == "" {
		return uuid.Nil, errors.New("missing todo id in url")
	}

	return uuid.Parse(idStr)
}
