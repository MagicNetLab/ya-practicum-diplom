package note

import (
	"context"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
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

func setupService() (Service, *rm.NoteRepository) {
	mockRepo := new(rm.NoteRepository)
	jwtCnf := new(mockJWTConfigurator)
	return MakeService(mockRepo, jwtCnf), mockRepo
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

// TestService_CreateNote тест создания заметки
func TestService_CreateNote(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное создание заметки", func(t *testing.T) {
		service, mockRepo := setupService()
		req := &pb.CreateNoteRequest{
			Title:   "Test Note",
			Content: "Test Content",
			Meta:    "Test Meta",
		}

		mockRepo.On("CreateNote", ctx, mock.AnythingOfType("*models.Note")).Return(nil)

		resp, err := service.CreateNote(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Note.ID)
		assert.Equal(t, req.Title, resp.Note.Title)
		assert.Equal(t, req.Content, resp.Note.Content)
		assert.Equal(t, req.Meta, resp.Note.Meta)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации данных заметки", func(t *testing.T) {
		service, _ := setupService()
		req := &pb.CreateNoteRequest{
			Title: "", // Empty title
		}

		resp, err := service.CreateNote(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Ошибка при создании заметки в репозитории", func(t *testing.T) {
		service, mockRepo := setupService()
		req := &pb.CreateNoteRequest{
			Title:   "Test Note",
			Content: "Test Content",
			Meta:    "Test Meta",
		}

		mockRepo.On("CreateNote", ctx, mock.AnythingOfType("*models.Note")).Return(assert.AnError)

		resp, err := service.CreateNote(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации токена", func(t *testing.T) {
		service, _ := setupService()
		ctx := context.Background()
		req := &pb.CreateNoteRequest{
			Title:   "Test Note",
			Content: "Test Content",
			Meta:    "Test Meta",
		}

		resp, err := service.CreateNote(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestService_GetNote тест получения заметки
func TestService_GetNote(t *testing.T) {
	service, mockRepo := setupService()
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное получение заметки", func(t *testing.T) {
		id := uuid.New().String()
		expectedNote := &models.Note{
			ID:        id,
			UID:       uid,
			Title:     "Test Note",
			Content:   "Test Content",
			Meta:      "Test Meta",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetNote", ctx, id, uid).Return(expectedNote, nil)

		resp, err := service.GetNote(ctx, &pb.GetNoteRequest{ID: id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, id, resp.Note.ID)
		assert.Equal(t, expectedNote.Title, resp.Note.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации идентификатора", func(t *testing.T) {
		resp, err := service.GetNote(ctx, &pb.GetNoteRequest{ID: "invalid-id"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Ошибка при получении заметки из репозитория", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("GetNote", ctx, id, uid).Return(nil, assert.AnError)

		resp, err := service.GetNote(ctx, &pb.GetNoteRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации токена", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.GetNoteRequest{}
		resp, err := service.GetNote(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}

// TestService_UpdateNote тест обновления заметки
func TestService_UpdateNote(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное обновление заметки", func(t *testing.T) {
		service, mockRepo := setupService()
		id := uuid.New().String()
		existingNote := &models.Note{
			ID:        id,
			UID:       uid,
			Title:     "Old Title",
			Content:   "Old Content",
			Meta:      "Old Meta",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetNote", ctx, id, uid).Return(existingNote, nil)
		mockRepo.On("UpdateNote", ctx, mock.AnythingOfType("*models.Note")).Return(nil)

		req := &pb.UpdateNoteRequest{
			ID:      id,
			Title:   "New Title",
			Content: "New Content",
			Meta:    "New Meta",
		}

		resp, err := service.UpdateNote(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Title, resp.Note.Title)
		assert.Equal(t, req.Content, resp.Note.Content)
		assert.Equal(t, req.Meta, resp.Note.Meta)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации идентификаторов", func(t *testing.T) {
		service, _ := setupService()
		req := &pb.UpdateNoteRequest{
			ID: "invalid-id",
		}

		resp, err := service.UpdateNote(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Ошибка при обновлении заметки в репозитории", func(t *testing.T) {
		service, mockRepo := setupService()
		id := uuid.New().String()
		existingNote := &models.Note{
			ID:        id,
			UID:       uid,
			Title:     "Old Title",
			Content:   "Old Content",
			Meta:      "Old Meta",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetNote", ctx, id, uid).Return(existingNote, nil)
		mockRepo.On("UpdateNote", ctx, mock.AnythingOfType("*models.Note")).Return(assert.AnError)

		req := &pb.UpdateNoteRequest{
			ID:      id,
			Title:   "New Title",
			Content: "New Content",
			Meta:    "New Meta",
		}

		resp, err := service.UpdateNote(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации токена", func(t *testing.T) {
		service, _ := setupService()
		ctx := context.Background()
		req := &pb.UpdateNoteRequest{}
		resp, err := service.UpdateNote(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}

// TestService_RemoveNote тест удаления заметки
func TestService_RemoveNote(t *testing.T) {
	service, mockRepo := setupService()
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное удаление заметки", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("RemoveNote", ctx, id, uid).Return(nil)

		resp, err := service.RemoveNote(ctx, &pb.RemoveNoteRequest{ID: id})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации идентификатора", func(t *testing.T) {
		resp, err := service.RemoveNote(ctx, &pb.RemoveNoteRequest{ID: "invalid-id"})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Ошибка при удалении заметки из репозитория", func(t *testing.T) {
		id := uuid.New().String()
		mockRepo.On("RemoveNote", ctx, id, uid).Return(assert.AnError)

		resp, err := service.RemoveNote(ctx, &pb.RemoveNoteRequest{ID: id})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации токена", func(t *testing.T) {
		service, _ := setupService()
		ctx := context.Background()
		req := &pb.RemoveNoteRequest{}
		resp, err := service.RemoveNote(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}

// TestService_ListNotes тест получения списка заметок
func TestService_ListNotes(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешное получение списка заметок", func(t *testing.T) {
		service, mockRepo := setupService()
		expectedNotes := make([]models.NoteModel, 0)
		expectedNotes = append(expectedNotes, &models.Note{
			ID:        uuid.New().String(),
			UID:       uid,
			Title:     "Note 1",
			Content:   "Content 1",
			Meta:      "Meta 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		expectedNotes = append(expectedNotes, &models.Note{
			ID:        uuid.New().String(),
			UID:       uid,
			Title:     "Note 2",
			Content:   "Content 2",
			Meta:      "Meta 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		mockRepo.On("SearchNote", ctx, mock.AnythingOfType("*models.NoteSearch")).Return(expectedNotes, nil)

		resp, err := service.ListNotes(ctx, &pb.ListNoteRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Notes, 2)
		assert.Equal(t, expectedNotes[0].GetTitle(), resp.Notes[0].Title)
		assert.Equal(t, expectedNotes[1].GetTitle(), resp.Notes[1].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при получении списка заметок из репозитория", func(t *testing.T) {
		service, mockRepo := setupService()
		mockRepo.On("SearchNote", ctx, mock.AnythingOfType("*models.NoteSearch")).Return(nil, assert.AnError)

		resp, err := service.ListNotes(ctx, &pb.ListNoteRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации идентификатора пользователя", func(t *testing.T) {
		service, _ := setupService()
		ctx := context.Background()
		resp, err := service.ListNotes(ctx, &pb.ListNoteRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}

// TestService_SearchNotes тест поиска заметок
func TestService_SearchNotes(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Успешный поиск заметок", func(t *testing.T) {
		service, mockRepo := setupService()
		expectedNotes := make([]models.NoteModel, 0)
		expectedNotes = append(expectedNotes, &models.Note{
			ID:        uuid.New().String(),
			UID:       uid,
			Title:     "Test Note",
			Content:   "Test Content",
			Meta:      "Test Meta",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		mockRepo.On("SearchNote", ctx, mock.AnythingOfType("*models.NoteSearch")).Return(expectedNotes, nil)

		req := &pb.SearchNoteRequest{
			Title:   "Test",
			Content: "Content",
			Meta:    "Meta",
			Limit:   10,
			Offset:  0,
		}
		resp, err := service.SearchNotes(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Notes, 1)
		assert.Equal(t, expectedNotes[0].GetTitle(), resp.Notes[0].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при поиске заметок в репозитории", func(t *testing.T) {
		service, mockRepo := setupService()
		mockRepo.On("SearchNote", ctx, mock.AnythingOfType("*models.NoteSearch")).Return(nil, assert.AnError)

		req := &pb.SearchNoteRequest{}
		resp, err := service.SearchNotes(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка валидации идентификатора пользователя", func(t *testing.T) {
		service, _ := setupService()
		ctx := context.Background()
		req := &pb.SearchNoteRequest{}
		resp, err := service.SearchNotes(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})
}
