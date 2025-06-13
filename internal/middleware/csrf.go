package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

// CSRF対策用のstate生成
func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// state検証用
func ValidateState(r *http.Request, expected string) bool {
	return r.URL.Query().Get("state") == expected
}
