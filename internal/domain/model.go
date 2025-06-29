package domain

type User struct {
	userID   string
	loginID  string
	password string
}

func (u *User) UserID() string   { return u.userID }
func (u *User) LoginID() string  { return u.loginID }
func (u *User) Password() string { return u.password }

type Client struct {
	clientID    string
	secret      string
	redirectURI string
}

func (c *Client) ClientID() string    { return c.clientID }
func (c *Client) Secret() string      { return c.secret }
func (c *Client) RedirectURI() string { return c.redirectURI }

type AuthCode struct {
	code      string
	userID    string
	clientID  string
	expiresAt int64
	state     string
}

func (a *AuthCode) Code() string     { return a.code }
func (a *AuthCode) UserID() string   { return a.userID }
func (a *AuthCode) ClientID() string { return a.clientID }
func (a *AuthCode) ExpiresAt() int64 { return a.expiresAt }
func (a *AuthCode) State() string    { return a.state }

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
