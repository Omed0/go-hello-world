package main

import (
	"net/http"
)

type Res struct {
	Status string `json:"status"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, Res{
		Status: "ok",
	})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something Went Wring")
}
