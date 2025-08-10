package domain

import (
	"reflect"
	"strings"
	"testing"
)

type testLogger struct{}

func (l *testLogger) Info(msg string, args ...any) {
	// do nothing
}
func (l *testLogger) Error(msg string, args ...any) {
	// do nothing
}
func (l *testLogger) Warn(msg string, args ...any) {
	// do nothing
}

func Test_認可コードフローで使用するparameter構築(t *testing.T) {
	logger := &testLogger{}

	tests := []struct {
		name         string
		responseType string
		clientID     string
		redirectURI  string
		scope        string
		state        string
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "正常系",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read",
			state:        "state123",
			wantErr:      false,
		},
		{
			name:         "サポートされていないresponse_type",
			responseType: "hoge",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read write",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "unsupported response_type: hoge",
		},
		{
			name:         "空のresponse_type",
			responseType: "",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read write",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "unsupported response_type: ",
		},
		{
			name:         "空のclient_id",
			responseType: "code",
			clientID:     "",
			redirectURI:  "https://example.com/callback",
			scope:        "read write",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "client_id is required",
		},
		{
			name:         "空のredirect_uri",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "",
			scope:        "read write",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "redirect_uri is required",
		},
		{
			name:         "空のscope",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "scope is required",
		},
		{
			name:         "正常系 scopeが複数",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read write",
			state:        "state123",
			wantErr:      false,
		},
		{
			name:         "空のscopeが含まれる場合",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read ",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "invalid scopes. Supported scopes are: read, write",
		},
		{
			name:         "サポートされていないscope",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "invalid-scope",
			state:        "state123",
			wantErr:      true,
			expectedErr:  "invalid scopes. Supported scopes are: read, write",
		},
		{
			name:         "空のstate",
			responseType: "code",
			clientID:     "client-1",
			redirectURI:  "https://example.com/callback",
			scope:        "read write",
			state:        "",
			wantErr:      true,
			expectedErr:  "state is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewAuthorizationCodeFlowParam(logger, tt.responseType, tt.clientID, tt.redirectURI, tt.scope, tt.state)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewAuthorizationCodeFlowParam() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("NewAuthorizationCodeFlowParam() error = %v, want %v", err.Error(), tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewAuthorizationCodeFlowParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			expectedResponseType, _ := GetResponseType(tt.responseType)
			if actual.ResponseType() != expectedResponseType {
				t.Errorf("ResponseType() = %v, want %v", actual.ResponseType(), expectedResponseType)
			}
			if actual.ClientID() != tt.clientID {
				t.Errorf("ClientID() = %v, want %v", actual.ClientID(), tt.clientID)
			}
			if actual.RedirectURI() != tt.redirectURI {
				t.Errorf("RedirectURI() = %v, want %v", actual.RedirectURI(), tt.redirectURI)
			}
			expectedScopes := strings.Split(tt.scope, " ")
			if !reflect.DeepEqual(actual.Scopes(), expectedScopes) {
				t.Errorf("Scopes() = %v, want %v", actual.Scopes(), expectedScopes)
			}
			if actual.State() != tt.state {
				t.Errorf("State() = %v, want %v", actual.State(), tt.state)
			}
		})
	}
}
