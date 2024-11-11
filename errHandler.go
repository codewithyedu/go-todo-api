package main

import (
	"log"
	"net/http"
)

func RespondWithErr(w http.ResponseWriter, r *http.Request, code int, msg string) {
	log.Println("error:", msg)
	RespondWithJSON(w, r, code, FT{"error": msg})
}
