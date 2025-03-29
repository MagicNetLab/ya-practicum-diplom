package interceptors

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/google/uuid"
	"os"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Handle(ctx context.Context, req any) (any, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestAuthInterceptor(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "secret")
	assert.NoError(t, err)
	cnf := config.GetJWTConfig()

	t.Run("успешная авторизация с валидным токеном", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)
		user := models.User{UID: uuid.New().String()}
		token, err := jwt.GenerateToken(&user, cnf.GetJWTSecret())
		assert.NoError(t, err)

		ctx := context.Background()
		md := metadata.New(map[string]string{"token": token})
		ctx = metadata.NewIncomingContext(ctx, md)

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Method"}
		result, err := AuthInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("пропуск авторизации для метода Auth", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)

		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: "/service/Auth"}
		result, err := AuthInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("пропуск авторизации для метода Register", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)

		ctx := context.Background()
		info := &grpc.UnaryServerInfo{FullMethod: "/service/Register"}
		result, err := AuthInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("отсутствие метаданных", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)

		ctx := context.Background()

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Method"}
		result, err := AuthInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.Error(t, err)
		assert.NotEqual(t, "success", result)
	})

	t.Run("невалидный токен", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)
		token := "invalid_token"

		ctx := context.Background()
		md := metadata.New(map[string]string{"token": token})
		ctx = metadata.NewIncomingContext(ctx, md)

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Method"}
		result, err := AuthInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.Error(t, err)
		assert.NotEqual(t, "success", result)
	})

	err = os.Unsetenv("JWT_SECRET")
	assert.NoError(t, err)
}

func TestGuestInterceptor(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "secret")
	assert.NoError(t, err)
	cnf := config.GetJWTConfig()

	t.Run("успешный доступ для гостя", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)

		ctx := context.Background()

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Auth"}
		result, err := GuestInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	t.Run("запрет доступа для авторизованного пользователя", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)
		user := models.User{UID: uuid.New().String()}
		token, err := jwt.GenerateToken(&user, cnf.GetJWTSecret())
		assert.NoError(t, err)

		ctx := context.Background()
		md := metadata.New(map[string]string{"token": token})
		ctx = metadata.NewIncomingContext(ctx, md)

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Auth"}
		result, err := GuestInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.Error(t, err)
		assert.NotEqual(t, "success", result)
	})

	t.Run("пропуск для не-Auth методов", func(t *testing.T) {
		mockHandler := &mockHandler{}
		mockHandler.On("Handle", mock.Anything, mock.Anything).Return("success", nil)

		ctx := context.Background()

		info := &grpc.UnaryServerInfo{FullMethod: "/service/Method"}
		result, err := GuestInterceptor(ctx, "request", info, mockHandler.Handle)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})

	err = os.Unsetenv("JWT_SECRET")
	assert.NoError(t, err)
}
