package authorize

type SuccessResponse struct {
	Message string `json:"message"`
}

var (
	ErrInvalidRequest          = "invalid_request"
	ErrUnauthorized            = "unauthorized_client"
	ErrAccessDenied            = "access_denied"
	ErrUnsupportedResponseType = "unsupported_response_type"
	ErrInvalidScope            = "invalid_scope"
	ErrServerError             = "server_error"
	ErrTemporarilyUnavailable  = "temporarily_unavailable"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Result interface {
	SuccessResponse | ErrorResponse
}
