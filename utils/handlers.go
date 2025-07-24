package utils

import (
	"fmt"
	"net/http"
)

type Res struct {
	Status string `json:"status"`
}

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, Res{Status: "ok"})
}

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}

func HandleRequestError(w http.ResponseWriter, err error, msg string, code int) bool {
	if err != nil {
		RespondWithError(w, code, fmt.Sprintf("%s: %v", msg, err))
		return true
	}
	return false
}
