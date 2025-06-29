package domain

import (
	"errors"
)

type ResponseType int

const (
	responseTypeCode ResponseType = iota
	// responseTypeToken
	// responseTypeIDToken
)

// とりあえずcodeのみサポート
var codeValueMap = map[string]ResponseType{
	"code": responseTypeCode,
	// "token":    responseTypeToken,
	// "id_token": responseTypeIDToken,
}

func IsSupportedResponseType(responseType string) error {
	_, ok := codeValueMap[responseType]
	if !ok {
		return errors.New("only supports response_type 'code'")
	}
	return nil
}
