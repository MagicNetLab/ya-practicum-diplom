package interceptors

import (
	"context"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models/mocks"
)

// TestAuthInterceptor тестирование AuthInterceptor
func TestAuthInterceptor(t *testing.T) {
	cnf := config.GetJWTConfig()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	t.Run("Успешная аутентификация", func(t *testing.T) {
		mockUser := new(mm.UserModel)
		mockUser.On("GetUID").Return(uuid.New().String())
		token, _ := jwt.GenerateToken(mockUser, cnf.GetJWTSecret())

		md := metadata.New(map[string]string{"token": token})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := AuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("Ошибка аутентификации: отсутствует токен", func(t *testing.T) {
		ctx := context.Background()

		_, err := AuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("Ошибка аутентификации: неверный токен", func(t *testing.T) {
		md := metadata.New(map[string]string{"token": "invalid_token"})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := AuthInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}

// TestGuestInterceptor тестирование GuestInterceptor
func TestGuestInterceptor(t *testing.T) {
	cnf := config.GetJWTConfig()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	t.Run("Успешный доступ для гостя", func(t *testing.T) {
		ctx := context.Background()

		resp, err := GuestInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp)
	})

	t.Run("Ошибка доступа для авторизованного пользователя", func(t *testing.T) {
		mockUser := new(mm.UserModel)
		mockUser.On("GetUID").Return(uuid.New().String())
		token, _ := jwt.GenerateToken(mockUser, cnf.GetJWTSecret())

		md := metadata.New(map[string]string{"token": token})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := GuestInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
		assert.Error(t, err)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))
	})
}
