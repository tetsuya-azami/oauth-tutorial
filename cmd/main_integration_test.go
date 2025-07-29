package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/pkg/logger"
	authorize "oauth-tutorial/internal/presentation/authorization"
	"oauth-tutorial/internal/session"
	usecase "oauth-tutorial/internal/usecase/authorization"
	"testing"
	"time"
)

const (
	mockSessionID = session.SessionID("mock-session-id")
)

type MockSessionIDGenerator struct{}

func (m *MockSessionIDGenerator) Generate() session.SessionID {
	return mockSessionID
}

func Test_認可リクエスト統合テスト(t *testing.T) {
	// given
	logger := logger.NewMyLogger()
	cr := infrastructure.NewClientRepository()
	sig := &MockSessionIDGenerator{}
	ss := infrastructure.NewSessionStorage()
	acf := usecase.NewAuthorizationCodeFlow(logger, cr, sig, ss)

	mux := http.NewServeMux()
	mux.Handle("/authorize", authorize.NewAuthorizeHandler(logger, acf))

	server := httptest.NewServer(mux)
	defer server.Close()

	// parameters
	responseType := "code"
	clientID := "iouobrnea"
	redirectURI := "https://client.example.com/callback"
	scope := "read"
	state := "test-state"

	url := fmt.Sprintf("%s/authorize?response_type=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		server.URL, responseType, clientID, redirectURI, scope, state)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// when
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// then
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if message, ok := result["message"]; !ok || message != "OK" {
		t.Errorf("Expected message 'OK', got %v", message)
	}

	sessiondata, err := ss.Get(mockSessionID)
	if err != nil {
		t.Errorf("Failed to get session parameter: %v", err)
	}
	if sessiondata == nil {
		t.Error("Expected non-nil session parameter, got nil")
	} else {
		expectedResponseType, _ := domain.GetResponseType(responseType)
		if sessiondata.AuthParam().ResponseType() != expectedResponseType {
			t.Errorf("Expected response type %d, got %d", expectedResponseType, sessiondata.AuthParam().ResponseType())
		}
		if sessiondata.AuthParam().ClientID() != clientID {
			t.Errorf("Expected client ID %s, got %s", clientID, sessiondata.AuthParam().ClientID())
		}
		if sessiondata.AuthParam().RedirectURI() != redirectURI {
			t.Errorf("Expected redirect URI %s, got %s", redirectURI, sessiondata.AuthParam().RedirectURI())
		}
		if len(sessiondata.AuthParam().Scopes()) == 0 || sessiondata.AuthParam().Scopes()[0] != scope {
			t.Errorf("Expected scope %s, got %v", scope, sessiondata.AuthParam().Scopes())
		}
		if sessiondata.AuthParam().Scopes()[0] != scope {
			t.Errorf("Expected scope %s, got %s", scope, sessiondata.AuthParam().Scopes()[0])
		}
		if sessiondata.AuthParam().State() != state {
			t.Errorf("Expected state %s, got %s", state, sessiondata.AuthParam().State())
		}
	}
}
