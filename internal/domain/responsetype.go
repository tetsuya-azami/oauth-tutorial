package domain

import (
	"errors"
)

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
		return notSupported, errors.New("unsupported response_type: " + responseType)
	}

	return r, nil
}
