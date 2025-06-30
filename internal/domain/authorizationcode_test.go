package domain

import (
	"testing"
	"time"
)

func TestNewAuthorizationCode(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		userID      string
		clientID    string
		scopes      []string
		redirectURI string
		now         time.Time
		expiresAt   int64
	}{
		{
			name:        "basic",
			value:       "test-value",
			userID:      "user-1",
			clientID:    "client-1",
			scopes:      []string{"openid profile"},
			redirectURI: "https://example.com/cb",
			now:         time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC).UTC(),
			expiresAt:   time.Date(2000, 1, 2, 3, 14, 5, 0, time.UTC).Local().Unix(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := NewAuthorizationCode(tt.value, tt.userID, tt.clientID, tt.scopes, tt.redirectURI, tt.now)

			if ac.Value() != tt.value {
				t.Errorf("Value() = %v, want %v", ac.Value(), tt.value)
			}
			if ac.UserID() != tt.userID {
				t.Errorf("UserID() = %v, want %v", ac.UserID(), tt.userID)
			}
			if ac.ClientID() != tt.clientID {
				t.Errorf("ClientID() = %v, want %v", ac.ClientID(), tt.clientID)
			}
			for i, scope := range ac.scopes {
				if scope != tt.scopes[i] {
					t.Errorf("Scopes() = %v, want %v", ac.Scopes(), tt.scopes)
				}
			}
			if ac.RedirectURI() != tt.redirectURI {
				t.Errorf("RedirectURI() = %v, want %v", ac.RedirectURI(), tt.redirectURI)
			}
			if ac.ExpiresAt() != tt.expiresAt {
				t.Errorf("ExpiresAt() = %v, want %v", ac.ExpiresAt(), tt.expiresAt)
			}
		})
	}
}
