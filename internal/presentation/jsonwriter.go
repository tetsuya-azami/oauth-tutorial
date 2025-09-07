package presentation

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse[T any](rw http.ResponseWriter, statusCode int, body T) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	rw.WriteHeader(statusCode)
	err := json.NewEncoder(rw).Encode(body)
	if err != nil {
		return err
	}

	return nil
}
