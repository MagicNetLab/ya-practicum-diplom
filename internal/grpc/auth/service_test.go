package auth

import (
	"context"
	"errors"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models/mocks"
)

func setupService() (*Service, *rm.AuthRepository) {
	mockRepo := new(rm.AuthRepository)
	cnf := config.GetJWTConfig()
	return &Service{store: mockRepo, jwt: cnf}, mockRepo
}

// TestService_Auth тест авторизации пользователя
func TestService_Auth(t *testing.T) {
	service, mockRepo := setupService()
	mockUser := new(mm.UserModel)
	mockUser.On("GetUID").Return(uuid.New().String())
	ctx := context.Background()

	t.Run("Успешная авторизация", func(t *testing.T) {
		mockRepo.On("GetUserByLoginAndPassword", ctx, "testuser", "password").Return(mockUser, nil)
		mockRepo.On("CreateToken", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := &pb.AuthRequest{Login: "testuser", Secret: "password"}
		resp, err := service.Auth(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("Ошибка авторизации: неверный логин или пароль", func(t *testing.T) {
		mockRepo.On("GetUserByLoginAndPassword", ctx, "wronguser", "wrongpassword").Return(nil, errors.New("user not found"))

		req := &pb.AuthRequest{Login: "wronguser", Secret: "wrongpassword"}
		resp, err := service.Auth(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Register тест регистрации пользователя
func TestService_Register(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Успешная регистрация", func(t *testing.T) {
		userMock := new(mm.UserModel)
		userMock.On("GetUID").Return(uuid.New().String())
		mockRepo.On("HasLogin", ctx, "newuser").Return(false, nil)
		mockRepo.On("CreateUser", ctx, "newuser", "password").Return(userMock, nil)

		req := &pb.RegRequest{Login: "newuser", Secret: "password"}
		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("Ошибка регистрации: логин уже существует", func(t *testing.T) {
		mockRepo.On("HasLogin", ctx, "existinguser").Return(true, nil)

		req := &pb.RegRequest{Login: "existinguser", Secret: "password"}
		resp, err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Refresh тест обновления токена
func TestService_Refresh(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Успешное обновление токена", func(t *testing.T) {
		mockUser := new(mm.UserModel)
		mockUser.On("GetUID").Return(uuid.New().String())
		token, _ := jwt.GenerateToken(mockUser, service.jwt.GetJWTSecret())
		mockRepo.On("HasToken", ctx, token).Return(true, nil)
		mockRepo.On("GetUserByUID", ctx, mock.Anything).Return(mockUser, nil)
		mockRepo.On("CreateToken", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := &pb.RefreshRequest{Token: token}
		resp, err := service.Refresh(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("Ошибка обновления токена: неверный токен", func(t *testing.T) {
		req := &pb.RefreshRequest{Token: "invalidtoken"}
		resp, err := service.Refresh(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
