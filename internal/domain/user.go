package domain

// OpenID Connectを実装する際はStandard Claimsを参照するのが良い。
// Standard Claims: https://openid.net/specs/openid-connect-core-1_0.html
// とりあえずは簡易的に実装
type User struct {
	userID   string
	loginID  string
	password string
}

func (u *User) UserID() string   { return u.userID }
func (u *User) LoginID() string  { return u.loginID }
func (u *User) Password() string { return u.password }
