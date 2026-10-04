package handlers

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, r *http.Request, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respondWithJSON(w, r, statusCode, map[string]string{"message": message})
}