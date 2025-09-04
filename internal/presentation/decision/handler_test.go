package decision

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"oauth-tutorial/internal/session"
	"oauth-tutorial/internal/usecase/decision"
	"oauth-tutorial/pkg/mylogger"
	"strings"
	"testing"
)

const (
	TestBaseRedirectURI = "https://example.com/callback"
)

// モックのPublishAuthorizationCodeUseCase
type mockPublishAuthorizationCodeUseCase struct {
	executeFunc func(*decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error)
}

func (m *mockPublishAuthorizationCodeUseCase) Execute(input *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(input)
	}
	return decision.PublishAuthorizationCodeOutput{}, nil
}

func TestDecisionHandler_ServeHTTP(t *testing.T) {
	logger := mylogger.NewMockLogger()

	tests := []struct {
		name                string
		formData            url.Values
		sessionCookie       *http.Cookie
		mockUseCase         *mockPublishAuthorizationCodeUseCase
		expectedStatus      int
		expectedBody        string
		expectedRedirectURL string
	}{
		{
			name: "異常ケース - approvedパラメータが無効",
			formData: url.Values{
				"approved": {"invalid"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			mockUseCase:    &mockPublishAuthorizationCodeUseCase{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"無効なリクエストです。もう一度初めからやり直してください"}`,
		},
		{
			name: "異常ケース - セッションCookieが存在しない",
			formData: url.Values{
				"approved": {"true"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie:  nil,
			mockUseCase:    &mockPublishAuthorizationCodeUseCase{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"セッションが見つかりません。もう一度初めからやり直してください"}`,
		},
		{
			name: "正常ケース - 承認された場合",
			formData: url.Values{
				"approved": {"true"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			mockUseCase: &mockPublishAuthorizationCodeUseCase{
				executeFunc: func(input *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error) {
					return decision.NewPublishAuthorizationCodeOutput(
						TestBaseRedirectURI,
						"test-auth-code",
						"test-state",
					), nil
				},
			},
			expectedStatus:      http.StatusSeeOther,
			expectedRedirectURL: "https://example.com/callback?code=test-auth-code&state=test-state",
		},
		{
			name: "異常ケース - セッションが見つからない",
			formData: url.Values{
				"approved": {"true"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "invalid-session-id",
			},
			mockUseCase: &mockPublishAuthorizationCodeUseCase{
				executeFunc: func(input *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error) {
					return decision.PublishAuthorizationCodeOutput{}, decision.NewErrPublishAuthorizationCode(decision.ErrSessionNotFound, "", "")
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"session not found"}`,
		},
		{
			name: "異常ケース - 認可が拒否された",
			formData: url.Values{
				"approved": {"false"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			mockUseCase: &mockPublishAuthorizationCodeUseCase{
				executeFunc: func(input *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error) {
					return decision.PublishAuthorizationCodeOutput{}, decision.NewErrPublishAuthorizationCode(decision.ErrAuthorizationDenied, TestBaseRedirectURI, "test-state")
				},
			},
			expectedStatus:      http.StatusSeeOther,
			expectedBody:        "",
			expectedRedirectURL: TestBaseRedirectURI + "?error=access_denied&error_description=" + decision.ErrAuthorizationDenied.Error() + "&state=test-state",
		},
		{
			name: "異常ケース - ログイン失敗",
			formData: url.Values{
				"approved": {"true"},
				"login_id": {"unknown"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			mockUseCase: &mockPublishAuthorizationCodeUseCase{
				executeFunc: func(input *decision.PublishAuthorizationCodeInput) (decision.PublishAuthorizationCodeOutput, error) {
					return decision.PublishAuthorizationCodeOutput{}, decision.NewErrPublishAuthorizationCode(decision.ErrInvalidLoginCredentials, "", "")
				},
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedBody:        `{"message":"invalid login credentials"}`,
			expectedRedirectURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			handler := NewDecisionHandler(logger, tt.mockUseCase)

			// リクエストの準備
			req := httptest.NewRequest("POST", "/decision", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}

			recorder := httptest.NewRecorder()

			// when
			handler.ServeHTTP(recorder, req)

			// then
			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if recorder.Code == http.StatusSeeOther {
				location := recorder.Header().Get("Location")
				if location != tt.expectedRedirectURL {
					t.Errorf("expected redirect to %s, got %s", tt.expectedRedirectURL, location)
				}
			} else {
				// JSONレスポンスの場合、ボディをチェック
				actualBody := strings.TrimSpace(recorder.Body.String())
				if actualBody != tt.expectedBody {
					t.Errorf("expected body %s, got %s", tt.expectedBody, actualBody)
				}
			}
		})
	}
}

func Test_リクエストパラメーターからInputへの変換(t *testing.T) {
	logger := mylogger.NewMockLogger()
	handler := NewDecisionHandler(logger, nil)

	tests := []struct {
		name          string
		formValues    url.Values
		sessionCookie *http.Cookie
		expectError   bool
		expectedError string
	}{
		{
			name: "正常ケース",
			formValues: url.Values{
				"approved": {"true"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			expectError: false,
		},
		{
			name: "異常ケース - approvedパラメータが無効",
			formValues: url.Values{
				"approved": {"invalid"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: &http.Cookie{
				Name:  session.SessionIDCookieName,
				Value: "test-session-id",
			},
			expectError:   true,
			expectedError: "無効なリクエストです。もう一度初めからやり直してください",
		},
		{
			name: "異常ケース - セッションCookieが存在しない",
			formValues: url.Values{
				"approved": {"true"},
				"login_id": {"testuser"},
				"password": {"testpass"},
			},
			sessionCookie: nil,
			expectError:   true,
			expectedError: "セッションが見つかりません。もう一度初めからやり直してください",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := httptest.NewRequest("POST", "/decision", nil)
			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}

			// when
			param, err := handler.convertParamToInput(tt.formValues, req)

			// then
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if err != nil && err.Error() != tt.expectedError {
					t.Errorf("expected error message %s, got %s", tt.expectedError, err.Error())
				}
				if param != nil {
					t.Error("expected param to be nil when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if param == nil {
					t.Error("expected param but got nil")
				}
			}
		})
	}
}
