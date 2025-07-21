package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/logger"
	authorize "oauth-tutorial/internal/presentation/authorization"
	usecase "oauth-tutorial/internal/usecase/authorization"
	"testing"
	"time"
)

func Test_認可リクエスト統合テスト(t *testing.T) {
	// given
	logger := logger.NewMyLogger()
	cr := infrastructure.NewClientRepository()
	aps := infrastructure.NewAuthParamSession()
	acf := usecase.NewAuthorizationCodeFlow(logger, cr, aps)

	mux := http.NewServeMux()
	mux.Handle("/authorize", authorize.NewAuthorizeHandler(logger, acf))

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(
		"GET",
		server.URL+"/authorize?response_type=code&client_id=iouobrnea&redirect_uri=https://client.example.com/callback&scope=read&state=test-state",
		nil,
	)
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

	t.Log("Integration test passed: authorization flow completed successfully")
}
