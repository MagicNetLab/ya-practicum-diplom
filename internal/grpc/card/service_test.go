package card

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

func setupService() (*Service, *rm.CardRepository) {
	mockRepo := new(rm.CardRepository)
	return &Service{store: mockRepo}, mockRepo
}

// TestService_Get тесты получения карты по идентификатору
func TestService_Get(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Проверка успешного получения карты по идентификатору", func(t *testing.T) {
		id := uuid.New().String()
		expectedCard := &models.Card{
			ID:     id,
			Name:   "Test Card",
			Number: "4111111111111111",
			Month:  12,
			Year:   2025,
		}

		mockRepo.On("GetCardByID", ctx, id).Return(expectedCard, nil)

		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, id, resp.Card.ID)
		assert.Equal(t, expectedCard.Name, resp.Card.Name)
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
		mockRepo.On("GetCardByID", ctx, id).Return(nil, assert.AnError)

		resp, err := service.Get(ctx, &pb.GetCardRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}

// TestService_Create тест создания карты
func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Проверка успешного создания карты", func(t *testing.T) {
		service, mockRepo := setupService()
		req := &pb.CreateCardRequest{
			UID:    uuid.New().String(),
			Name:   "Test Card",
			Number: "4111111111111111",
			Month:  12,
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
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации данных карты", func(t *testing.T) {
		service, _ := setupService()
		req := &pb.CreateCardRequest{
			UID:    "", // Invalid UID
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
		service, mockRepo := setupService()
		req := &pb.CreateCardRequest{
			UID:    uuid.New().String(),
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
}

// TestService_Delete тест удаления карты
func TestService_Delete(t *testing.T) {
	service, mockRepo := setupService()
	ctx := context.Background()

	t.Run("Проверка успешного удаления карты", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("DeleteCard", ctx, id).Return(nil)

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
		mockRepo.On("DeleteCard", ctx, id).Return(assert.AnError)

		resp, err := service.Delete(ctx, &pb.DeleteCardRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}

// TestService_Search tests the Search method
func TestService_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("Success with all search parameters", func(t *testing.T) {
		service, mockRepo := setupService()
		searchReq := &pb.SearchCardRequest{
			UID:    uuid.New().String(),
			Number: "4111",
			Name:   "Test",
			Year:   2025,
			Limit:  10,
			Offset: 0,
		}

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

		resp, err := service.Search(ctx, searchReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, len(expectedCards))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with only UID", func(t *testing.T) {
		service, mockRepo := setupService()
		searchReq := &pb.SearchCardRequest{
			UID: uuid.New().String(),
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

	t.Run("Success with only Number", func(t *testing.T) {
		service, mockRepo := setupService()
		searchReq := &pb.SearchCardRequest{
			Number: "4111",
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

	t.Run("Success with pagination", func(t *testing.T) {
		service, mockRepo := setupService()
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

	t.Run("Success with empty result", func(t *testing.T) {
		service, mockRepo := setupService()
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

	t.Run("Repository Error", func(t *testing.T) {
		service, mockRepo := setupService()
		searchReq := &pb.SearchCardRequest{
			UID: uuid.New().String(),
		}

		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(nil, assert.AnError)

		resp, err := service.Search(ctx, searchReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}

// TestService_List тестирование метода List сервиса CardService
func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("Проверка успешного получения списка карт", func(t *testing.T) {
		service, mockRepo := setupService()
		uid := uuid.New().String()
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

		resp, err := service.List(ctx, &pb.ListCardRequest{UID: uid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Cards, len(expectedCards))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Проверка ошибки валидации UID", func(t *testing.T) {
		service, _ := setupService()
		resp, err := service.List(ctx, &pb.ListCardRequest{UID: "invalid-uid"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Проверка ошибки при получении списка карт из репозитория", func(t *testing.T) {
		service, mockRepo := setupService()
		uid := uuid.New().String()
		mockRepo.On("SearchCards", ctx, mock.AnythingOfType("models.CardSearch")).Return(nil, assert.AnError)

		resp, err := service.List(ctx, &pb.ListCardRequest{UID: uid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}
