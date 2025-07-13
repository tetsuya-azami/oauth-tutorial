package domain

import (
	"reflect"
	"testing"
)

func Test_client再構築(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientName   string
		secret       string
		redirectURIs []string
	}{
		{
			name:         "正常系",
			clientID:     "client-1",
			clientName:   "Test Client",
			secret:       "secret123",
			redirectURIs: []string{"https://example.com/callback"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ReconstructClient(tt.clientID, tt.clientName, tt.secret, tt.redirectURIs)

			if client.ClientID() != tt.clientID {
				t.Errorf("ClientID() = %v, want %v", client.ClientID(), tt.clientID)
			}
			if client.ClientName() != tt.clientName {
				t.Errorf("ClientName() = %v, want %v", client.ClientName(), tt.clientName)
			}
			if client.Secret() != tt.secret {
				t.Errorf("Secret() = %v, want %v", client.Secret(), tt.secret)
			}
			if !reflect.DeepEqual(client.RedirectURI(), tt.redirectURIs) {
				t.Errorf("RedirectURI() = %v, want %v", client.RedirectURI(), tt.redirectURIs)
			}
		})
	}
}

func Test_clinetがredirectURIを持っているか検査(t *testing.T) {
	okUri := "https://example.com/callback"
	tests := []struct {
		name     string
		testURI  string
		expected bool
	}{
		{
			name:     "正常系",
			testURI:  okUri,
			expected: true,
		},
		{
			name:     "clientが持っていないredirectURIの場合",
			testURI:  "https://malicious.com/callback",
			expected: false,
		},
		{
			name:     "空のURIの場合",
			testURI:  "",
			expected: false,
		},
		{
			name:     "clientが持っているURIと前方部分一致している場合",
			testURI:  "https://example.com/callback/extra",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ReconstructClient("test-client", "Test Client", "secret", []string{okUri, "https://app.example.com/auth"})

			result := client.ContainsRedirectURI(tt.testURI)
			if result != tt.expected {
				t.Errorf("ContainsRedirectURI(%v) = %v, want %v", tt.testURI, result, tt.expected)
			}
		})
	}
}
