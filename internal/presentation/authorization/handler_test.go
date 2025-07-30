package authorize

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/internal/session"
	usecase "oauth-tutorial/internal/usecase/authorize"
	"oauth-tutorial/pkg/mylogger"
	"testing"
)

type MockAuthorizationFlow struct {
	err error
}

func NewMockAuthorizationFlow(err error) *MockAuthorizationFlow {
	return &MockAuthorizationFlow{
		err: err,
	}
}

func (m *MockAuthorizationFlow) Execute(param *domain.AuthorizationCodeFlowParam) (session.SessionID, error) {
	return "test-session-id", m.err
}

func TestAuthorizeHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockErr        error
		wantStatusCode int
		wantHeader     map[string]string
		wantResponse   any
	}{
		{
			name: "正常ケース - 認可成功",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "test-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantHeader: map[string]string{
				"Set-Cookie":   session.SessionIDCookieName + "=test-session-id; Path=/; HttpOnly; Secure",
				"Content-Type": "application/json",
			},
			wantResponse: SuccessResponse{
				Message: "OK",
			},
		},
		{
			name: "異常ケース - 不正なリクエストパラメーター",
			queryParams: map[string]string{
				"response_type": "invalid",
				"client_id":     "test-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        nil,
			wantStatusCode: http.StatusBadRequest,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "unsupported response_type: invalid",
			},
		},
		{
			name: "異常ケース - クライアントが見つからない",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "nonexistent-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        usecase.ErrClientNotFound,
			wantStatusCode: http.StatusBadRequest,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "client not found",
			},
		},
		{
			name: "異常ケース - 無効なリダイレクトURI",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "test-client",
				"redirect_uri":  "https://malicious.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        usecase.ErrInvalidRedirectURI,
			wantStatusCode: http.StatusBadRequest,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "invalid redirect URI",
			},
		},
		{
			name: "異常ケース - 予期しないエラー",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "test-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        errors.New("database connection failed"),
			wantStatusCode: http.StatusInternalServerError,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "database connection failed",
			},
		},
		{
			name: "異常ケース - サーバーエラー",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "test-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        usecase.ErrServer,
			wantStatusCode: http.StatusInternalServerError,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "server error occurred",
			},
		},
		{
			name: "異常ケース - 予期しないエラー",
			queryParams: map[string]string{
				"response_type": "code",
				"client_id":     "test-client",
				"redirect_uri":  "https://example.com/callback",
				"scope":         "read write",
				"state":         "test-state",
			},
			mockErr:        usecase.ErrUnExpected,
			wantStatusCode: http.StatusInternalServerError,
			wantHeader:     map[string]string{"Content-Type": "application/json"},
			wantResponse: ErrorResponse{
				Message: "unexpected error occurred",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			logger := mylogger.NewMockLogger()
			flow := NewMockAuthorizationFlow(tt.mockErr)
			handler := NewAuthorizeHandler(logger, flow)

			reqURL := buildRequestURL(tt.queryParams)

			req := httptest.NewRequest(http.MethodGet, reqURL, nil)

			rr := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rr, req)

			// then
			if rr.Code != tt.wantStatusCode {
				t.Errorf("Status code = %d, want %d", rr.Code, tt.wantStatusCode)
			}

			// headersが正しいこと
			for key, value := range tt.wantHeader {
				if rr.Header().Get(key) != value {
					t.Errorf("Header %s = %s, want %s", key, rr.Header().Get(key), value)
				}
			}

			var actualResponse any
			switch tt.wantResponse.(type) {
			case ErrorResponse:
				actualResponse = &ErrorResponse{}
			case SuccessResponse:
				actualResponse = &SuccessResponse{}
			}

			if err := json.NewDecoder(rr.Body).Decode(actualResponse); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// response bodyが正しいこと
			expectedJSONResponse, _ := json.Marshal(tt.wantResponse)
			actualJSONResponse, _ := json.Marshal(actualResponse)

			if string(actualJSONResponse) != string(expectedJSONResponse) {
				t.Errorf("Response body = %s, want %s", actualJSONResponse, expectedJSONResponse)
			}
		})
	}
}

func buildRequestURL(queryParams map[string]string) string {
	reqURL := "http://example.com/authorize"
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		reqURL += "?" + params.Encode()
	}
	return reqURL
}
