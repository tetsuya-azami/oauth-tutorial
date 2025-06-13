package domain

// User represents a resource owner
// TODO: 必要に応じて拡張

type User struct {
	ID   string
	Name string
}

type Client struct {
	ID          string
	Secret      string
	RedirectURI string
}

type AuthCode struct {
	Code      string
	UserID    string
	ClientID  string
	ExpiresAt int64
	State     string
}

type Token struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	ClientID     string
	ExpiresAt    int64
}
