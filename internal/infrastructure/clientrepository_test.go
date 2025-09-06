package infrastructure

import (
	"oauth-tutorial/internal/domain"
	"testing"
)

func Test_クライアントIDでのクライアント取得(t *testing.T) {
	tests := []struct {
		name         string
		clientID     domain.ClientID
		expectError  error
		expectedName string
	}{
		{
			name:         "existing client",
			clientID:     domain.ClientID("iouobrnea"),
			expectError:  nil,
			expectedName: "client-1",
		},
		{
			name:        "non-existing client",
			clientID:    domain.ClientID("nonexistent"),
			expectError: ErrClientNotFound,
		},
		{
			name:        "empty client ID",
			clientID:    domain.ClientID(""),
			expectError: ErrClientNotFound,
		},
	}

	repo := NewClientRepository()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := repo.SelectByClientID(tt.clientID)

			// エラーが期待される場合
			if tt.expectError != nil {
				if err != tt.expectError {
					t.Errorf("expected error %v, got %v", tt.expectError, err)
				}
				if client != nil {
					t.Errorf("expected nil client when error occurs, got %v", client)
				}
				return
			}

			// 正常ケース
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Errorf("expected client but got nil")
				return
			}

			if client.ClientName() != tt.expectedName {
				t.Errorf("expected client name %s, got %s", tt.expectedName, client.ClientName())
			}
		})
	}
}
