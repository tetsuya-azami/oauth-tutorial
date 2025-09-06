package domain

type GrantType int

const (
	grantTypeNotSupported GrantType = iota
	grantTypeAuthorizationCode
	grantTypeRefreshToken
)

var grantTypeValueMap = map[string]GrantType{
	"authorization_code": grantTypeAuthorizationCode,
	"refresh_token":      grantTypeRefreshToken,
}

type UnsupportedGrantTypeError struct {
	GrantType string
}

func (e *UnsupportedGrantTypeError) Error() string {
	return "unsupported grant_type: " + e.GrantType
}

func GetGrantType(grantType string) (GrantType, error) {
	g, ok := grantTypeValueMap[grantType]
	if !ok {
		return grantTypeNotSupported, &UnsupportedGrantTypeError{GrantType: grantType}
	}
	return g, nil
}
