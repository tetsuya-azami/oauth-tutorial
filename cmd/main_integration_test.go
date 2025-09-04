package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/infrastructure"
	"oauth-tutorial/internal/infrastructure/dto"
	pAuthorize "oauth-tutorial/internal/presentation/authorize"
	pDecision "oauth-tutorial/internal/presentation/decision"
	"oauth-tutorial/internal/session"
	uAuthorize "oauth-tutorial/internal/usecase/authorize"
	uDecision "oauth-tutorial/internal/usecase/decision"
	"oauth-tutorial/pkg/mylogger"
	"strings"
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
	logger := mylogger.NewLogger()
	cr := infrastructure.NewClientRepository()
	sig := &MockSessionIDGenerator{}
	ss := infrastructure.NewSessionStorage()
	acf := uAuthorize.NewAuthorizationCodeFlow(logger, cr, sig, ss)

	mux := http.NewServeMux()
	mux.Handle("GET /authorize", pAuthorize.NewAuthorizeHandler(logger, acf))

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

type MockAuthzCodeGenerator struct{}

func (*MockAuthzCodeGenerator) GenerateURLSafeRandomString(n int) string {
	return "mock-authz-code"
}

func Test_認可コード発行統合テスト(t *testing.T) {
	// given
	logger := mylogger.NewLogger()
	rg := &MockAuthzCodeGenerator{}

	testRedirectURI := "http://callback.example.com"
	ss := infrastructure.NewSessionStorage()
	mockState := "mock-state"
	param, _ := domain.NewAuthorizationCodeFlowParam(logger, "code", "client_1", testRedirectURI, "read", mockState)
	ss.Save(mockSessionID, dto.NewSessionData(param, nil))

	ur := infrastructure.NewUserRepository()
	ar := infrastructure.NewAuthCodeRepository()
	pac := uDecision.NewPublishAuthorizationCodeUseCase(logger, rg, ss, ur, ar)

	mux := http.NewServeMux()
	mux.Handle("POST /decision", pDecision.NewDecisionHandler(logger, pac))

	server := httptest.NewServer(mux)
	defer server.Close()

	// リクエストの準備
	url := fmt.Sprintf("%s/decision", server.URL)
	header := "application/x-www-form-urlencoded"
	requestBody := fmt.Sprintf("approved=%s&login_id=%s&password=%s", "true", "test-user@example.com", "password")
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", header)
	req.AddCookie(&http.Cookie{Name: session.SessionIDCookieName, Value: string(mockSessionID)})

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 勝手にリダイレクトしない設定
			return http.ErrUseLastResponse
		},
	}

	// when
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// then
	// ステータスコードが302であること
	actualStatusCode := resp.StatusCode
	expectedStatusCode := http.StatusSeeOther
	if actualStatusCode != expectedStatusCode {
		t.Errorf("Expected status %d, got %d", expectedStatusCode, actualStatusCode)
	}

	// リダイレクト先が正しいこと
	expectedLocation := fmt.Sprintf("%s?code=%s&state=%s", testRedirectURI, "mock-authz-code", mockState)
	actualLocation := resp.Header.Get("Location")
	if actualLocation != expectedLocation {
		t.Errorf("Expected Location header %q, got %q", expectedLocation, actualLocation)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// 正常系ではレスポンスボディは空になること
	if string(result) != "" {
		t.Errorf("Expected non-empty response body, got %q", string(result))
	}

	// 認可リクエストのパラメーターがセッションから削除されていること
	_, err = ss.Get(mockSessionID)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	} else {
		if err != infrastructure.ErrSessionNotFound {
			t.Errorf("Expected ErrSessionNotFound, got %v", err)
		}
	}

	// 認可コードレポジトリに認可コードが保存されていること
	authzCode, err := ar.FindByCode("mock-authz-code")
	if err != nil {
		t.Errorf("Failed to find authorization code: %v", err)
	} else {
		if authzCode == nil {
			t.Error("Expected authorization code to be found, got nil")
		}
		if authzCode.Value() != "mock-authz-code" {
			t.Errorf("Expected authorization code %q, got %q", "mock-authz-code", authzCode.Value())
		}
		// ユーザーIDは現状userrepositoryのモックデータに依存
		if authzCode.UserID() != "IU7ewbuvey" {
			t.Errorf("Expected user ID %q, got %q", "IU7ewbuvey", authzCode.UserID())
		}
		// ...他の検証
	}
}
