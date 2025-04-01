package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cm "github.com/MagicNetLab/ya-practicum-diplom/internal/config/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"
)

// setupService - инициализирует сервис и моки
func setupService() (*Service, *rm.AuthRepository, *cm.AppConfig) {
	mockRepo := new(rm.AuthRepository)
	cnf := new(cm.AppConfig)
	return &Service{store: mockRepo, cnf: cnf}, mockRepo, cnf
}

// TestService_Auth тест авторизации пользователя
func TestService_Auth(t *testing.T) {
	service, mockRepo, _ := setupService()
	mockUser := new(mm.UserModel)
	mockUser.On("GetUID").Return(uuid.New().String())
	pass, err := encryptor.EncryptPassword("password")
	assert.NoError(t, err)
	mockUser.On("GetPassword").Return(pass)
	ctx := context.Background()

	t.Run("Успешная авторизация", func(t *testing.T) {
		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(mockUser, nil)
		mockRepo.On("CreateToken", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := &pb.AuthRequest{Login: "testuser", Secret: "password"}
		resp, err := service.Auth(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("Ошибка авторизации: неверный логин или пароль", func(t *testing.T) {
		mockRepo.On("GetUserByLogin", ctx, "wronguser").Return(nil, errors.New("user not found"))

		req := &pb.AuthRequest{Login: "wronguser", Secret: "wrongpassword"}
		resp, err := service.Auth(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Register тест регистрации пользователя
func TestService_Register(t *testing.T) {
	service, mockRepo, _ := setupService()
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
	service, mockRepo, cnf := setupService()
	ctx := context.Background()

	t.Run("Успешное обновление токена", func(t *testing.T) {
		cnf.On("GetJWTSecret").Return("secret")
		mockUser := new(mm.UserModel)
		mockUser.On("GetUID").Return(uuid.New().String())
		token, _ := jwt.GenerateToken(mockUser, service.cnf.JWTSecret())
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
