package domain

type GrantType int

const (
	GrantTypeNotSupported GrantType = iota
	GrantTypeAuthorizationCode
	GrantTypeRefreshToken
)

var grantTypeValueMap = map[string]GrantType{
	"authorization_code": GrantTypeAuthorizationCode,
	"refresh_token":      GrantTypeRefreshToken,
}

type UnsupportedGrantTypeError struct {
	GrantType string
}

func (e *UnsupportedGrantTypeError) Error() string {
	return "unsupported grant_type: " + e.GrantType
}

func ResolveGrantType(grantType string) (GrantType, error) {
	g, ok := grantTypeValueMap[grantType]
	if !ok {
		return GrantTypeNotSupported, &UnsupportedGrantTypeError{GrantType: grantType}
	}
	return g, nil
}
