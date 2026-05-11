package httpapi

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{Data: data})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{Error: message})
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
