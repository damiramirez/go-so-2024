package serialization

import (
	"encoding/json"
	"net/http"
)

// Decode el body de una solicitud HTTP en generic
// Usar cuando queremos leer un JSON que viene del body y pasarlo a una struct de Go
func DecodeHTTPBody[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var data T
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
	}
	return data, nil
}

// Encode la estructura en JSON y la escribe en el body de la respuesta HTTP
// Usar cuando queres responder con un JSON
func EncodeHTTPResponse[T any](w http.ResponseWriter, data T, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
