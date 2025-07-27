package presentation

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse[T any](rw http.ResponseWriter, statusCode int, body T) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)

	return nil
}
