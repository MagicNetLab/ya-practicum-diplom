package account

import (
	"context"
	"errors"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

type mockJWTConfigurator struct {
	mock.Mock
}

func (m *mockJWTConfigurator) GetJWTSecret() string {
	return "test-secret"
}

func (m *mockJWTConfigurator) IsValid() bool { return true }

func (m *mockJWTConfigurator) GetTokenLifeTime() time.Duration {
	return time.Hour
}

func (m *mockJWTConfigurator) GetRefreshTokenLifeTime() time.Duration {
	return time.Hour
}

func setupService() (*Service, *mocks.AccountRepository, *mockJWTConfigurator) {
	mockRepo := new(mocks.AccountRepository)
	mockJWTCnf := new(mockJWTConfigurator)
	service, _ := MakeService(mockRepo, mockJWTCnf)
	return &service, mockRepo, mockJWTCnf
}

func setupAuthContext(uid string) context.Context {
	user := &models.User{
		UID:       uid,
		Login:     "test-login",
		Password:  "test-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	token, _ := jwt.GenerateToken(user, "test-secret")
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestService_GetAccount(t *testing.T) {
	service, mockRepo, _ := setupService()
	id := uuid.New().String()
	uid := uuid.New().String()
	mockAccount := new(mm.AccountModel)
	mockAccount.On("GetID").Return(id)
	mockAccount.On("GetUID").Return(uid)
	mockAccount.On("GetLogin").Return("test-login")
	mockAccount.On("GetPassword").Return("test-password")
	mockAccount.On("GetURL").Return("test-url")
	mockAccount.On("GetDescription").Return("test-description")
	ctx := setupAuthContext(uid)

	t.Run("Успешное получение аккаунта", func(t *testing.T) {
		mockRepo.On("GetAccount", ctx, id, uid).Return(mockAccount, nil)

		req := &pb.GetAccountRequest{Id: id}
		resp, err := service.GetAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, id, resp.Account.Id)
		assert.Equal(t, uid, resp.Account.Uid)
		assert.Equal(t, "test-login", resp.Account.Login)
		assert.Equal(t, "test-password", resp.Account.Password)
		assert.Equal(t, "test-url", resp.Account.Url)
		assert.Equal(t, "test-description", resp.Account.Description)
	})

	t.Run("Ошибка: аккаунт не найден", func(t *testing.T) {
		mockRepo.On("GetAccount", ctx, "non-existent-id", "test-uid").Return(nil, errors.New("account not found"))

		req := &pb.GetAccountRequest{Id: "non-existent-id"}
		resp, err := service.GetAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Ошибка: неверный формат ID", func(t *testing.T) {
		req := &pb.GetAccountRequest{Id: "invalid-id"}
		resp, err := service.GetAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Ошибка: не аутентифицирован", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.GetAccountRequest{Id: "test-id"}
		resp, err := service.GetAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestService_CreateAccount(t *testing.T) {
	service, mockRepo, _ := setupService()
	id := uuid.New().String()
	uid := uuid.New().String()
	mockAccount := new(mm.AccountModel)
	mockAccount.On("GetID").Return(id)
	mockAccount.On("GetUID").Return(uid)
	mockAccount.On("GetLogin").Return("test-login")
	mockAccount.On("GetPassword").Return("test-password")
	mockAccount.On("GetURL").Return("test-url")
	mockAccount.On("GetDescription").Return("test-description")
	ctx := setupAuthContext(uid)

	t.Run("Успешное создание аккаунта", func(t *testing.T) {
		mockRepo.On("CreateAccount", ctx, uid, "test-login", "test-password", "test-url", "test-description").Return(mockAccount, nil)

		req := &pb.CreateAccountRequest{
			Login:       "test-login",
			Password:    "test-password",
			Url:         "test-url",
			Description: "test-description",
		}
		resp, err := service.CreateAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, id, resp.Acc.Id)
		assert.Equal(t, uid, resp.Acc.Uid)
		assert.Equal(t, "test-login", resp.Acc.Login)
		assert.Equal(t, "test-password", resp.Acc.Password)
		assert.Equal(t, "test-url", resp.Acc.Url)
		assert.Equal(t, "test-description", resp.Acc.Description)
	})

	t.Run("Ошибка создания аккаунта", func(t *testing.T) {
		mockRepo.On("CreateAccount", ctx, uid, "invalid-login", "invalid-password", "invalid-url", "invalid-description").Return(nil, errors.New("failed to create account"))

		req := &pb.CreateAccountRequest{
			Login:       "invalid-login",
			Password:    "invalid-password",
			Url:         "invalid-url",
			Description: "invalid-description",
		}
		resp, err := service.CreateAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Ошибка: не аутентифицирован", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.CreateAccountRequest{
			Login:       "test-login",
			Password:    "test-password",
			Url:         "test-url",
			Description: "test-description",
		}
		resp, err := service.CreateAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestService_RemoveAccount(t *testing.T) {
	service, mockRepo, _ := setupService()
	uid := uuid.New().String()
	id := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное удаление аккаунта", func(t *testing.T) {
		mockRepo.On("RemoveAccount", ctx, id, uid).Return(nil)

		req := &pb.RemoveAccountRequest{Id: id}
		resp, err := service.RemoveAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Ошибка удаления аккаунта", func(t *testing.T) {
		nonexistentID := uuid.New().String()
		mockRepo.On("RemoveAccount", ctx, nonexistentID, uid).Return(errors.New("account not found"))

		req := &pb.RemoveAccountRequest{Id: nonexistentID}
		resp, err := service.RemoveAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Ошибка: не аутентифицирован", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.RemoveAccountRequest{Id: uuid.New().String()}
		resp, err := service.RemoveAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestService_SearchAccounts(t *testing.T) {
	service, mockRepo, _ := setupService()
	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()

	mockAccount1 := new(mm.AccountModel)
	mockAccount1.On("GetID").Return(id1)
	mockAccount1.On("GetUID").Return(uid1)
	mockAccount1.On("GetLogin").Return("test-login-1")
	mockAccount1.On("GetPassword").Return("test-password-1")
	mockAccount1.On("GetURL").Return("test-url-1")
	mockAccount1.On("GetDescription").Return("test-description-1")

	mockAccount2 := new(mm.AccountModel)
	mockAccount2.On("GetID").Return(id2)
	mockAccount2.On("GetUID").Return(uid1)
	mockAccount2.On("GetLogin").Return("test-login-2")
	mockAccount2.On("GetPassword").Return("test-password-2")
	mockAccount2.On("GetURL").Return("test-url-2")
	mockAccount2.On("GetDescription").Return("test-description-2")

	mockAccount3 := new(mm.AccountModel)
	mockAccount3.On("GetID").Return(id3)
	mockAccount3.On("GetUID").Return(uid2)
	mockAccount3.On("GetLogin").Return("test-login-3")
	mockAccount3.On("GetPassword").Return("test-password-3")
	mockAccount3.On("GetURL").Return("test-url-3")
	mockAccount3.On("GetDescription").Return("test-description-3")

	t.Run("Успешный поиск аккаунтов", func(t *testing.T) {
		ctx := setupAuthContext(uid1)
		mockRepo.On("SearchAccounts", ctx, mock.AnythingOfType("models.AccountSearch")).Return([]models.AccountModel{mockAccount1, mockAccount2}, nil)

		req := &pb.SearchAccountRequest{
			Login:       "test-login",
			Url:         "test-url",
			Description: "test-description",
			Limit:       "10",
			Offset:      "0",
		}
		resp, err := service.SearchAccounts(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Acc, 2)
		assert.Equal(t, uid1, resp.Acc[0].Uid)
		assert.Equal(t, id1, resp.Acc[0].Id)
		assert.Equal(t, uid1, resp.Acc[1].Uid)
		assert.Equal(t, id2, resp.Acc[1].Id)
	})

	t.Run("Ошибка поиска аккаунтов", func(t *testing.T) {
		ctx := setupAuthContext(uid1)
		mockRepo.On("SearchAccounts", ctx, mock.AnythingOfType("models.AccountSearch")).Return(nil, errors.New("search failed"))

		req := &pb.SearchAccountRequest{}
		resp, err := service.SearchAccounts(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Ошибка: не аутентифицирован", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.SearchAccountRequest{}
		resp, err := service.SearchAccounts(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
