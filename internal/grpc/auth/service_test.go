package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repo/mocks"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService() (*Service, *rm.Repository) {
	mockRepo := new(rm.Repository)
	cnf, _ := conf.GetCnf()
	return &Service{store: mockRepo, cnf: cnf}, mockRepo
}

// TestService_Auth тест авторизации пользователя
func TestService_Auth(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Успешная авторизация", func(t *testing.T) {
		mockRepo.On("GetUserByLoginAndPassword", ctx, "testuser", "password").Return(&mm.UserEntity{UID: "test_uid"}, nil)
		mockRepo.On("CreateToken", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

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
		mockRepo.On("HasLogin", ctx, "newuser").Return(false, nil)
		mockRepo.On("CreateUser", ctx, "newuser", "password").Return(&mm.UserEntity{}, nil)

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
		token, _ := jwt.GenerateToken(&mm.UserEntity{UID: "test_uid"}, service.cnf.JWTSecret())
		mockRepo.On("HasToken", ctx, token, mock.Anything, true).Return(true, nil)
		mockRepo.On("GetUserByUID", ctx, mock.Anything).Return(&mm.UserEntity{UID: "test_uid"}, nil)
		mockRepo.On("CreateToken", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

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
