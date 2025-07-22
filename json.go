package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to Marshal JSON Response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Response With 5xx error: ", msg)
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}
