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
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	State            string `json:"state"`
}

type Result interface {
	SuccessResponse | ErrorResponse
}
