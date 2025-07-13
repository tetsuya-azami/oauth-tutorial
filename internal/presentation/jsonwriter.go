package presentation

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Result interface {
	SuccessResponse | ErrorResponse
}

// TODO: レスポンスをロギング(ResponseWriterWithLoggingみたいなのを作るのはありかも？)
func WriteJSONResponse[T Result](rw http.ResponseWriter, statusCode int, body T) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(body)

	return nil
}
