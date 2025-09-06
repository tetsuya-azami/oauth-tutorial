package infrastructure

import (
	"oauth-tutorial/internal/domain"
	"oauth-tutorial/pkg/mycrypto"
	"sort"
	"sync"
	"testing"
	"time"
)

type MockRandomGenerator struct{}

func (r *MockRandomGenerator) GenerateURLSafeRandomString(n int) string {
	return "mocked-random-string"
}

func TestNewAuthCodeRepository(t *testing.T) {
	repo := NewAuthCodeRepository()

	if repo == nil {
		t.Error("NewAuthCodeRepository() should not return nil")
		return
	}

	if repo.authCodeStore == nil {
		t.Error("authCodeStore should be initialized")
	}

	if len(repo.authCodeStore) != 0 {
		t.Error("authCodeStore should be empty initially")
	}
}

func Test_認可コードの保存(t *testing.T) {
	repo := NewAuthCodeRepository()

	// テスト用の認可コードを作成
	authCode := domain.NewAuthorizationCode(&MockRandomGenerator{}, "test-user-id", "test-client-id", []string{"read"}, "https://example.com/callback", time.Now())

	originalLength := len(repo.authCodeStore)

	// Save操作をテスト
	repo.Save(authCode)

	// 保存されたことを確認
	if len(repo.authCodeStore) != originalLength+1 {
		t.Errorf("authCodeStore length = %d, want 1", len(repo.authCodeStore))
	}

	// 正しいキーで保存されていることを確認
	if _, exists := repo.authCodeStore[authCode.Value()]; !exists {
		t.Error("authCode should be saved with correct key")
	}
}

func Test_認可コードの検索(t *testing.T) {
	repo := NewAuthCodeRepository()

	// テスト用の認可コードを作成・保存
	expectedAuthCode := domain.NewAuthorizationCode(&MockRandomGenerator{}, "test-user", "test-client", []string{"read"}, "https://example.com/callback", time.Now())
	repo.Save(expectedAuthCode)

	tests := []struct {
		name      string
		code      string
		wantErr   bool
		expecterr error
	}{
		{
			name:      "存在するコードを検索",
			code:      expectedAuthCode.Value(),
			wantErr:   false,
			expecterr: nil,
		},
		{
			name:      "存在しないコードを検索",
			code:      "non-existing-code",
			wantErr:   true,
			expecterr: ErrAuthorizationCodeNotFound,
		},
		{
			name:      "空のコードを検索",
			code:      "",
			wantErr:   true,
			expecterr: ErrAuthorizationCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.FindByCode(tt.code)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FindByCode() error = nil, wantErr %v", tt.wantErr)
				}
				if err != tt.expecterr {
					t.Errorf("FindByCode() error = %v, want %v", err, tt.expecterr)
				}
			} else {
				if err != nil {
					t.Errorf("FindByCode() error = %v, wantErr %v", err, tt.wantErr)
				}
				if actual == nil {
					t.Error("FindByCode() should return a valid AuthorizationCode")
				}
				if actual.Value() != expectedAuthCode.Value() {
					t.Errorf("FindByCode() result.Value() = %v, want %v", actual.Value(), expectedAuthCode.Value())
				}
				if actual.UserID() != expectedAuthCode.UserID() {
					t.Errorf("FindByCode() result.UserID() = %v, want %v", actual.UserID(), expectedAuthCode.UserID())
				}
				if actual.ClientID() != expectedAuthCode.ClientID() {
					t.Errorf("FindByCode() result.ClientID() = %v, want %v", actual.ClientID(), expectedAuthCode.ClientID())
				}
				if len(actual.Scopes()) != len(expectedAuthCode.Scopes()) {
					t.Errorf("FindByCode() result.Scopes() length = %d, want %d", len(actual.Scopes()), len(expectedAuthCode.Scopes()))
				}
				// Scopesの内容をソートして比較
				actualScopes := actual.Scopes()
				expectedScopes := expectedAuthCode.Scopes()
				if len(actualScopes) != len(expectedScopes) {
					t.Errorf("FindByCode() result.Scopes() length = %d, want %d", len(actualScopes), len(expectedScopes))
				} else {
					// ソートしてから比較
					as := append([]string{}, actualScopes...)
					es := append([]string{}, expectedScopes...)
					sort.Strings(as)
					sort.Strings(es)
					for i := range as {
						if as[i] != es[i] {
							t.Errorf("FindByCode() result.Scopes[%d] = %v, want %v", i, as[i], es[i])
						}
					}
				}

				if actual.RedirectURI() != expectedAuthCode.RedirectURI() {
					t.Errorf("FindByCode() result.RedirectURI() = %v, want %v", actual.RedirectURI(), expectedAuthCode.RedirectURI())
				}
				if actual.ExpiresAt() != expectedAuthCode.ExpiresAt() {
					t.Errorf("FindByCode() result.ExpiresAt() = %d, want %d", actual.ExpiresAt(), expectedAuthCode.ExpiresAt())
				}
			}
		})
	}
}

func TestAuthCodeRepository_Delete(t *testing.T) {
	repo := NewAuthCodeRepository()

	// テスト用の認可コードを作成・保存
	authCode := domain.NewAuthorizationCode(&MockRandomGenerator{}, "test-user", "test-client", []string{"read"}, "https://example.com/callback", time.Now())
	repo.Save(authCode)

	// 削除前の確認
	value, _ := repo.FindByCode(authCode.Value())
	if value == nil {
		t.Error("FindByCode() should return a valid AuthorizationCode before deletion")
	}

	// Delete操作をテスト
	repo.Delete(authCode.Value())

	value, _ = repo.FindByCode(authCode.Value())
	if value != nil {
		t.Error("FindByCode() should return nil after deletion")
	}
}

func TestAuthCodeRepository_ConcurrentAccess(t *testing.T) {
	repo := NewAuthCodeRepository()

	// 複数のゴルーチンで同時にアクセス
	var wg sync.WaitGroup
	numGoroutines := 100

	// 同時に複数の認可コードを保存
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			authCode := domain.NewAuthorizationCode(
				&mycrypto.RandomGenerator{},
				"test-user",
				"test-client",
				[]string{"read"},
				"https://example.com/callback",
				time.Now(),
			)
			repo.Save(authCode)
		}(i)
	}

	wg.Wait()

	// 全てのコードが保存されていることを確認
	if len(repo.authCodeStore) != numGoroutines {
		t.Errorf("authCodeStore length = %d, want %d", len(repo.authCodeStore), numGoroutines)
	}
}
