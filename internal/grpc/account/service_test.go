package account

import (
	"context"
	"errors"
	"testing"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	mm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService() (*Service, *rm.AccountRepository) {
	mockRepo := new(rm.AccountRepository)
	return &Service{store: mockRepo}, mockRepo
}

// TestService_GetAccount тест получения аккаунта
func TestService_GetAccount(t *testing.T) {
	service, mockRepo := setupService()
	mockAccount := new(mm.AccountModel)
	mockAccount.On("GetID").Return("test-id")
	mockAccount.On("GetUID").Return("test-uid")
	mockAccount.On("GetLogin").Return("test-login")
	mockAccount.On("GetPassword").Return("test-password")
	mockAccount.On("GetURL").Return("test-url")
	mockAccount.On("GetDescription").Return("test-description")
	ctx := context.Background()

	t.Run("Успешное получение аккаунта", func(t *testing.T) {
		mockRepo.On("GetAccount", ctx, "test-id").Return(mockAccount, nil)

		req := &pb.GetAccountRequest{Id: "test-id"}
		resp, err := service.GetAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test-id", resp.Account.Id)
		assert.Equal(t, "test-uid", resp.Account.Uid)
		assert.Equal(t, "test-login", resp.Account.Login)
		assert.Equal(t, "test-password", resp.Account.Password)
		assert.Equal(t, "test-url", resp.Account.Url)
		assert.Equal(t, "test-description", resp.Account.Description)
	})

	t.Run("Ошибка: аккаунт не найден", func(t *testing.T) {
		mockRepo.On("GetAccount", ctx, "non-existent-id").Return(nil, errors.New("account not found"))

		req := &pb.GetAccountRequest{Id: "non-existent-id"}
		resp, err := service.GetAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_CreateAccount тест создания аккаунта
func TestService_CreateAccount(t *testing.T) {
	service, mockRepo := setupService()
	mockAccount := new(mm.AccountModel)
	mockAccount.On("GetID").Return("test-id")
	mockAccount.On("GetUID").Return("test-uid")
	mockAccount.On("GetLogin").Return("test-login")
	mockAccount.On("GetPassword").Return("test-password")
	mockAccount.On("GetURL").Return("test-url")
	mockAccount.On("GetDescription").Return("test-description")
	ctx := context.Background()

	t.Run("Успешное создание аккаунта", func(t *testing.T) {
		mockRepo.On("CreateAccount", ctx, "test-uid", "test-login", "test-password", "test-url", "test-description").Return(mockAccount, nil)

		req := &pb.CreateAccountRequest{
			Uid:         "test-uid",
			Login:       "test-login",
			Password:    "test-password",
			Url:         "test-url",
			Description: "test-description",
		}
		resp, err := service.CreateAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test-id", resp.Acc.Id)
		assert.Equal(t, "test-uid", resp.Acc.Uid)
		assert.Equal(t, "test-login", resp.Acc.Login)
		assert.Equal(t, "test-password", resp.Acc.Password)
		assert.Equal(t, "test-url", resp.Acc.Url)
		assert.Equal(t, "test-description", resp.Acc.Description)
	})

	t.Run("Ошибка создания аккаунта", func(t *testing.T) {
		mockRepo.On("CreateAccount", ctx, "invalid-uid", "invalid-login", "invalid-password", "invalid-url", "invalid-description").Return(nil, errors.New("failed to create account"))

		req := &pb.CreateAccountRequest{
			Uid:         "invalid-uid",
			Login:       "invalid-login",
			Password:    "invalid-password",
			Url:         "invalid-url",
			Description: "invalid-description",
		}
		resp, err := service.CreateAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_RemoveAccount тест удаления аккаунта
func TestService_RemoveAccount(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Успешное удаление аккаунта", func(t *testing.T) {
		mockRepo.On("RemoveAccount", ctx, "test-id").Return(nil)

		req := &pb.RemoveAccountRequest{Id: "test-id"}
		resp, err := service.RemoveAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Ошибка удаления аккаунта", func(t *testing.T) {
		mockRepo.On("RemoveAccount", ctx, "non-existent-id").Return(errors.New("account not found"))

		req := &pb.RemoveAccountRequest{Id: "non-existent-id"}
		resp, err := service.RemoveAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_SearchAccounts тест успешный поиск аккаунтов
func TestService_SearchAccounts_Success(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	mockAccount1 := new(mm.AccountModel)
	mockAccount1.On("GetID").Return(uuid.New().String())
	mockAccount1.On("GetUID").Return("test-uid-1")
	mockAccount1.On("GetLogin").Return("test-login-1")
	mockAccount1.On("GetPassword").Return("test-password-1")
	mockAccount1.On("GetURL").Return("test-url-1")
	mockAccount1.On("GetDescription").Return("test-description-1")

	mockAccount2 := new(mm.AccountModel)
	mockAccount2.On("GetID").Return(uuid.New().String())
	mockAccount2.On("GetUID").Return("test-uid-2")
	mockAccount2.On("GetLogin").Return("test-login-2")
	mockAccount2.On("GetPassword").Return("test-password-2")
	mockAccount2.On("GetURL").Return("test-url-2")
	mockAccount2.On("GetDescription").Return("test-description-2")

	t.Run("Успешный поиск аккаунтов", func(t *testing.T) {
		mockRepo.On("SearchAccounts", ctx, mock.AnythingOfType("models.AccountSearch")).Return([]models.AccountModel{mockAccount1, mockAccount2}, nil)

		req := &pb.SearchAccountRequest{
			Uid:         "test-uid",
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
		assert.Equal(t, "test-uid-1", resp.Acc[0].Uid)
		assert.Equal(t, "test-uid-2", resp.Acc[1].Uid)
	})
}

// TestService_SearchAccounts тест поиска аккаунтов с пустым результатом
func TestService_SearchAccounts_Empty(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Поиск аккаунтов: пустой результат", func(t *testing.T) {
		mockRepo.On("SearchAccounts", ctx, mock.AnythingOfType("models.AccountSearch")).Return([]models.AccountModel{}, nil)

		req := &pb.SearchAccountRequest{}
		resp, err := service.SearchAccounts(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Acc, 0)
	})
}

// TestService_SearchAccounts тест поиска аккаунтов с ошибкой
func TestService_SearchAccounts_Error(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Ошибка поиска аккаунтов", func(t *testing.T) {
		mockRepo.On("SearchAccounts", ctx, mock.AnythingOfType("models.AccountSearch")).Return(nil, errors.New("search failed"))

		req := &pb.SearchAccountRequest{}
		resp, err := service.SearchAccounts(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
