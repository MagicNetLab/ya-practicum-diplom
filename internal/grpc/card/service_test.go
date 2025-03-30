package card

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// mockJWTConfigurator mock для JWTConfigurator
type mockJWTConfigurator struct {
	mock.Mock
}

// GetJWTSecret mock для метода получения секрета JWT
func (m *mockJWTConfigurator) GetJWTSecret() string {
	return "test-secret"
}

// GetJWTSecret mock для метода проверки валидности JWT
func (m *mockJWTConfigurator) IsValid() bool { return true }

// GetTokenLifeTime mock для метода получения времени жизни токена
func (m *mockJWTConfigurator) GetTokenLifeTime() time.Duration {
	return time.Hour
}

// GetRefreshTokenLifeTime mock для метода получения времени жизни refresh токена
func (m *mockJWTConfigurator) GetRefreshTokenLifeTime() time.Duration {
	return time.Hour
}

// setupService - настройка сервиса для тестов
func setupService() (*Service, *rm.CardRepository, *mockJWTConfigurator) {
	mockRepo := new(rm.CardRepository)
	mockJWTCnf := new(mockJWTConfigurator)
	return &Service{store: mockRepo, jwt: mockJWTCnf}, mockRepo, mockJWTCnf
}

// setupAuthContext - настройка контекста с валидным токеном авторизации для тестов
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

// TestService_Get тесты получения карты по идентификатору
func TestService_Get(t *testing.T) {
	service, mockRepo, _ := setupService()
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Проверка успешного получения карты по идентификатору", func(t *testing.T) {
		id := uuid.New().String()
		expectedCard := &models.Card{
			ID:     id,
			UID:    uid,
			Name:   "Test Card",
			Number: "4111111111111111",
			Mask:   "test-mask",
			Month:  12,
			Year:   2025,
			CVC:    "test-cvc",
			PIN:    "test-pin",
			Meta:   "test-meta",
		}

		mockRepo.On("GetCardByID", ctx, id, uid).Return(expectedCard, nil)

		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, id, resp.Card.ID)
		assert.Equal(t, expectedCard.Name, resp.Card.Name)
		assert.Equal(t, expectedCard.Number, resp.Card.Number)
		assert.Equal(t, expectedCard.Meta, resp.Card.Meta)
		assert.Equal(t, int32(expectedCard.Month), resp.Card.Month)
		assert.Equal(t, int32(expectedCard.Year), resp.Card.Year)
		assert.Equal(t, expectedCard.CVC, resp.Card.CVC)
		assert.Equal(t, expectedCard.PIN, resp.Card.PIN)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации идентификатора", func(t *testing.T) {
		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: "invalid-id"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Проверка ошибки при получении карты из репозитория", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("GetCardByID", ctx, id, uid).Return(nil, assert.AnError)

		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации токена", func(t *testing.T) {
		ctx = context.Background()
		req := &pb.GetCardRequest{ID: uuid.New().String()}
		resp, err := service.Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Create тест создания карты
func TestService_Create(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Проверка успешного создания карты", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		req := &pb.CreateCardRequest{
			Name:   "Test Card",
			Number: "4111111111111111",
			Month:  12,
			Meta:   "test-meta",
			Year:   2025,
			CVC:    "123",
			PIN:    "1234",
		}

		mockRepo.On("CreateCard", ctx, mock.AnythingOfType("*models.Card")).Return(nil)

		resp, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Card.ID)
		assert.Equal(t, req.Name, resp.Card.Name)
		assert.Equal(t, "**** **** **** 1111", resp.Card.Number)
		assert.Equal(t, req.Meta, resp.Card.Meta)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации данных карты", func(t *testing.T) {
		service, _, _ := setupService()
		req := &pb.CreateCardRequest{
			Name:   "Test Card",
			Number: "invalid",
			Month:  13,   // Invalid month
			Year:   2020, // Past year
		}

		resp, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Проверка ошибки при создании карты в репозитории", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		req := &pb.CreateCardRequest{
			Name:   "Test Card",
			Number: "4532015112830366",
			Month:  12,
			Year:   2025,
			CVC:    "123",
			PIN:    "1234",
		}

		mockRepo.On("CreateCard", ctx, mock.AnythingOfType("*models.Card")).Return(assert.AnError)

		resp, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации токена", func(t *testing.T) {
		service, _, _ := setupService()
		req := &pb.CreateCardRequest{
			Name:   "Test Card",
			Number: "4532015112830366",
			Month:  12,
			Year:   2025,
			CVC:    "123",
			PIN:    "1234",
		}
		ctx := context.Background()
		resp, err := service.Create(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_Delete тест удаления карты
func TestService_Delete(t *testing.T) {
	service, mockRepo, _ := setupService()
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Проверка успешного удаления карты", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("DeleteCard", ctx, id, uid).Return(nil)

		resp, err := service.Delete(ctx, &pb.DeleteCardRequest{ID: id})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации идентификатора", func(t *testing.T) {
		resp, err := service.Delete(ctx, &pb.DeleteCardRequest{ID: "invalid-id"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Проверка ошибки отсутствия карты в репозитории", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("DeleteCard", ctx, id, uid).Return(assert.AnError)

		resp, err := service.Delete(ctx, &pb.DeleteCardRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации токена", func(t *testing.T) {
		ctx = metadata.NewIncomingContext(context.Background(), metadata.Pairs("token", "invalid"))
		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: uuid.New().String()})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}

// TestService_Search tests the Search method
func TestService_Search(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Проверка успешного поиска карты с корректными параметрами поиска", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		searchReq := &pb.SearchCardRequest{
			Name:   "test",
			Limit:  10,
			Offset: 0,
		}

		expectedCards := []models.CardModel{
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card 1",
				Number: "4111111111111111",
				Meta:   "test-meta",
				Month:  12,
				Year:   2025,
			},
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card 2",
				Number: "4111222233334444",
				Meta:   "test-meta",
				Month:  11,
				Year:   2025,
			},
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(expectedCards, nil)

		resp, err := service.Search(ctx, searchReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, len(expectedCards))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка успешного поиска карты с пустыми параметрами поиска", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid)
		searchReq := &pb.SearchCardRequest{}

		expectedCards := []models.CardModel{
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card",
				Number: "4111111111111111",
				Month:  12,
				Year:   2025,
			},
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(expectedCards, nil)

		resp, err := service.Search(ctx, searchReq)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка успешного поиска карт с пагинацией", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid)
		searchReq := &pb.SearchCardRequest{
			Limit:  2,
			Offset: 1,
		}

		expectedCards := []models.CardModel{
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card",
				Number: "4111111111111111",
				Month:  12,
				Year:   2025,
			},
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(expectedCards, nil)

		resp, err := service.Search(ctx, searchReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка успешного поиска карты по имени с пустым результатом", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid)
		searchReq := &pb.SearchCardRequest{
			Name: "NonExistentCard",
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return([]models.CardModel{}, nil)

		resp, err := service.Search(ctx, searchReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Cards)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка поиска с  ошибкой от репозитория", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid)
		searchReq := &pb.SearchCardRequest{}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(nil, assert.AnError)

		resp, err := service.Search(ctx, searchReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации токена", func(t *testing.T) {
		service, _, _ := setupService()
		ctx := context.Background()
		searchReq := &pb.SearchCardRequest{}

		resp, err := service.Search(ctx, searchReq)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_List тестирование метода List сервиса CardService
func TestService_List(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Проверка успешного получения списка карт", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		expectedCards := []models.CardModel{
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card 1",
				Number: "4111111111111111",
				Month:  12,
				Year:   2025,
			},
			&models.Card{
				ID:     uuid.New().String(),
				Name:   "Test Card 2",
				Number: "4111222233334444",
				Month:  11,
				Year:   2025,
			},
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(expectedCards, nil)

		resp, err := service.List(ctx, &pb.ListCardRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, len(expectedCards))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки при получении списка карт из репозитория", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(nil, assert.AnError)

		resp, err := service.List(ctx, &pb.ListCardRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации UID", func(t *testing.T) {
		service, _, _ := setupService()
		ctx = context.Background()
		resp, err := service.List(ctx, &pb.ListCardRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}
