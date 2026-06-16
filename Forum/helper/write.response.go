package helper

import (
	"encoding/json"
	"exemple_api/dto"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, dto.ApiError{Status: status, Error: message})
}

// WriteErrorResponse écrit une réponse d'erreur
func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	WriteError(w, status, message)
}

// WriteSuccessResponse écrit une réponse de succès
func WriteSuccessResponse(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, data)
}
