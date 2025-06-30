package domain

type User struct {
	userID   string
	loginID  string
	password string
}

func (u *User) UserID() string   { return u.userID }
func (u *User) LoginID() string  { return u.loginID }
func (u *User) Password() string { return u.password }

type Token struct {
	accessToken  string
	refreshToken string
	userID       string
	clientID     string
	expiresAt    int64
}

func (t *Token) AccessToken() string  { return t.accessToken }
func (t *Token) RefreshToken() string { return t.refreshToken }
func (t *Token) UserID() string       { return t.userID }
func (t *Token) ClientID() string     { return t.clientID }
func (t *Token) ExpiresAt() int64     { return t.expiresAt }
