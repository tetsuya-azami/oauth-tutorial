package token

type SuccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

type ErrorType int

const (
	InvalidRequest ErrorType = iota
	InvalidClient
	InvalidGrant
	UnauthorizedClient
	UnsupportedGrantType
	InvalidScope
	ServerError
)

var (
	errorTypeToString = map[ErrorType]string{
		InvalidRequest:       "invalid_request",
		InvalidClient:        "invalid_client",
		InvalidGrant:         "invalid_grant",
		UnauthorizedClient:   "unauthorized_client",
		UnsupportedGrantType: "unsupported_grant_type",
		InvalidScope:         "invalid_scope",
	}
)

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewErrorResponse(errorType ErrorType, description string) ErrorResponse {
	return ErrorResponse{
		Error:            errorTypeToString[errorType],
		ErrorDescription: description,
	}
}
