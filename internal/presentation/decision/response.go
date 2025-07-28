package decision

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Result interface {
	SuccessResponse | ErrorResponse
}
