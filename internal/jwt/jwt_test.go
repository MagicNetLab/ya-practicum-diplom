package jwt

import (
	"context"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// mockUserModel mock для UserModel
type mockUserModel struct {
	uid string
}

// GetUID возвращает UID пользователя
func (m mockUserModel) GetUID() string {
	return m.uid
}

// TestGenerateToken тестирование генерации токена
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		user      models.UserModel
		jwtSecret string
		wantErr   bool
	}{
		{
			name:      "Valid token generation",
			user:      &models.User{UID: uuid.New().String()},
			jwtSecret: "test-secret",
			wantErr:   false,
		},
		{
			name:      "Empty secret",
			user:      &models.User{UID: uuid.New().String()},
			jwtSecret: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.user, tt.jwtSecret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

// TestParseToken тестирование парсинга токена
func TestParseToken(t *testing.T) {
	validUID := uuid.New().String()
	validSecret := "test-secret"
	validToken, _ := GenerateToken(&models.User{UID: validUID}, validSecret)

	tests := []struct {
		name       string
		token      string
		jwtSecret  string
		wantUID    string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			jwtSecret: validSecret,
			wantUID:   validUID,
			wantErr:   false,
		},
		{
			name:       "Invalid token format",
			token:      "invalid-token",
			jwtSecret:  validSecret,
			wantErr:    true,
			wantErrMsg: "failed check token: invalid token",
		},
		{
			name:       "Wrong secret",
			token:      validToken,
			jwtSecret:  "wrong-secret",
			wantErr:    true,
			wantErrMsg: "failed check token: invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token, tt.jwtSecret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUID, claims.UID)
			}
		})
	}
}

// TestVerifyToken тестирование проверки токена
func TestVerifyToken(t *testing.T) {
	validUID := uuid.New().String()
	validSecret := "test-secret"
	validToken, _ := GenerateToken(&models.User{UID: validUID}, validSecret)

	tests := []struct {
		name      string
		token     string
		jwtSecret string
		want      bool
	}{
		{
			name:      "Valid token",
			token:     validToken,
			jwtSecret: validSecret,
			want:      true,
		},
		{
			name:      "Invalid token",
			token:     "invalid-token",
			jwtSecret: validSecret,
			want:      false,
		},
		{
			name:      "Wrong secret",
			token:     validToken,
			jwtSecret: "wrong-secret",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyToken(tt.token, tt.jwtSecret)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestGetRandomSecret тестирование генерации случайного секрета
func TestGetRandomSecret(t *testing.T) {
	secret1 := GetRandomSecret()
	secret2 := GetRandomSecret()

	assert.NotEmpty(t, secret1)
	assert.NotEmpty(t, secret2)
	assert.NotEqual(t, secret1, secret2)
	assert.Len(t, secret1, 64) // 32 bytes in hex = 64 characters
}

// TestExtractToken тестирование извлечения токена из контекста
func TestExtractToken(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		wantToken  string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Valid token in metadata",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{"token": "test-token"}),
			),
			wantToken: "test-token",
			wantErr:   false,
		},
		{
			name:       "No metadata",
			ctx:        context.Background(),
			wantErr:    true,
			wantErrMsg: "failed to parse metadata",
		},
		{
			name: "Empty token",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{}),
			),
			wantErr:    true,
			wantErrMsg: "no token provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractToken(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

// TestGetUIDFromContext тестирование извлечения UID из контекста
func TestGetUIDFromContext(t *testing.T) {
	validUID := uuid.New().String()
	validSecret := "test-secret"
	validToken, _ := GenerateToken(&models.User{UID: validUID}, validSecret)

	tests := []struct {
		name       string
		ctx        context.Context
		jwtSecret  string
		wantUID    string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Valid token and UID",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{"token": validToken}),
			),
			jwtSecret: validSecret,
			wantUID:   validUID,
			wantErr:   false,
		},
		{
			name:       "No metadata",
			ctx:        context.Background(),
			jwtSecret:  validSecret,
			wantErr:    true,
			wantErrMsg: "failed to parse metadata",
		},
		{
			name: "Invalid token",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{"token": "invalid-token"}),
			),
			jwtSecret:  validSecret,
			wantErr:    true,
			wantErrMsg: "failed check token: invalid token",
		},
		{
			name: "Empty UID in token",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.New(map[string]string{"token": validToken}),
			),
			jwtSecret:  "wrong-secret",
			wantErr:    true,
			wantErrMsg: "failed check token: invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := GetUIDFromContext(tt.ctx, tt.jwtSecret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUID, uid)
			}
		})
	}
}
