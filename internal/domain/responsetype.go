package domain

import (
	"fmt"
)

type UnsupportedResponseTypeError struct {
	ResponseType string
}

func (e *UnsupportedResponseTypeError) Error() string {
	return fmt.Sprintf("unsupported response_type: %s", e.ResponseType)
}

type ResponseType int

const (
	notSupported ResponseType = iota
	responseTypeCode
	// responseTypeToken
	// responseTypeIDToken
)

// とりあえずcodeのみサポート
var codeValueMap = map[string]ResponseType{
	"code": responseTypeCode,
	// "token":    responseTypeToken,
	// "id_token": responseTypeIDToken,
}

func GetResponseType(responseType string) (ResponseType, error) {
	r, ok := codeValueMap[responseType]
	if !ok {
		return notSupported, &UnsupportedResponseTypeError{ResponseType: responseType}
	}

	return r, nil
}
