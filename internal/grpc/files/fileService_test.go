package files

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/s3/mocks"
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

func setupService() (*Service, *rm.FileRepository, *mocks.S3Client) {
	mockRepo := new(rm.FileRepository)
	mockStorage := new(mocks.S3Client)
	jwtCnf := new(mockJWTConfigurator)
	return &Service{db: mockRepo, storage: mockStorage, jwt: jwtCnf}, mockRepo, mockStorage
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

// TestService_Put проверка успешной загрузки файла.
func TestService_Put(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("успешная загрузка файла", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("CreateFile", ctx, mock.Anything).Return(nil)
		mockStorage.On("PutObject", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		resp, err := service.Put(ctx, &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("content"),
			Meta:    "meta",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp.File)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Ошибка при записи данных в БД", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("CreateFile", ctx, mock.Anything).Return(assert.AnError)

		resp, err := service.Put(ctx, &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("content"),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("ОШибка валидации UID", func(t *testing.T) {
		service, _, _ := setupService()
		resp, err := service.Put(context.Background(), &pb.PutFileRequest{})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Некорректные данные в запросе", func(t *testing.T) {
		service, _, _ := setupService()
		tests := []struct {
			name string
			req  *pb.PutFileRequest
		}{
			{"EmptyName", &pb.PutFileRequest{Content: []byte("content")}},
			{"EmptyContent", &pb.PutFileRequest{Name: "test.txt"}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := service.Put(ctx, tt.req)
				assert.Error(t, err)
				assert.Nil(t, resp)
				statusErr, _ := status.FromError(err)
				assert.Equal(t, codes.InvalidArgument, statusErr.Code())
			})
		}
	})

	t.Run("Ошибка при загрузке файла в хранилище", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("CreateFile", ctx, mock.Anything).Return(nil)
		mockStorage.On("PutObject", ctx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
		mockRepo.On("DeleteFile", ctx, mock.Anything, mock.Anything).Return(nil)

		resp, err := service.Put(ctx, &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("content"),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}

// TestService_List проверка успешного получения списка файлов.
func TestService_List(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("SearchFile", ctx, mock.Anything).Return([]models.FilesModel{}, nil)

		resp, err := service.List(ctx, &pb.ListFilesRequest{Limit: 10})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("SearchFile", ctx, mock.Anything).Return(nil, assert.AnError)

		resp, err := service.List(ctx, &pb.ListFilesRequest{})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, statusErr.Code())
	})
}

// TestService_Download проверка успешного скачивания файла.
func TestService_Download(t *testing.T) {
	uid := uuid.New().String()
	id := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("StorageGetFailure", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("GetFile", ctx, id, mock.Anything).Return(models.File{}, nil)
		mockStorage.On("GetObject", ctx, mock.Anything).Return(nil, assert.AnError)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: id})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, statusErr.Code())
	})

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("GetFile", ctx, id, mock.Anything).Return(models.File{}, nil)
		mockStorage.On("GetObject", ctx, mock.Anything).Return([]byte("content"), nil)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: id})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		inorrectID := "incorrect-id"
		mockRepo.On("GetFile", ctx, inorrectID, uid).Return(nil, assert.AnError)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: inorrectID})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("NotFound", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("GetFile", ctx, id, mock.Anything).Return(nil, assert.AnError)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: id})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, statusErr.Code())
	})
}

// TestService_Search проверка успешного поиска файлов.
func TestService_Search(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)

	t.Run("RepositoryError", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("SearchFile", ctx, mock.Anything).Return(nil, assert.AnError)

		resp, err := service.Search(ctx, &pb.SearchFilesRequest{})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, statusErr.Code())
	})

	t.Run("WithFilters", func(t *testing.T) {
		service, mockRepo, _ := setupService()

		extends := make([]models.FilesModel, 0)
		extends = append(extends, &models.File{ID: "id1", UID: uid, Name: "test1.txt", Path: "path1", Meta: "meta1", CreatedAt: time.Now()})
		mockRepo.On("SearchFile", ctx, mock.Anything).Return(extends, nil)

		resp, err := service.Search(ctx, &pb.SearchFilesRequest{
			Name:  "test",
			Meta:  "meta",
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyResults", func(t *testing.T) {
		service, mockRepo, _ := setupService()
		mockRepo.On("SearchFile", ctx, mock.Anything).Return([]models.FilesModel{}, nil)

		resp, err := service.Search(ctx, &pb.SearchFilesRequest{})
		assert.NoError(t, err)
		assert.Empty(t, resp.Files)
	})
}

// TestService_Remove проверка успешного удаления файла.
func TestService_Remove(t *testing.T) {
	uid := uuid.New().String()
	ctx := setupAuthContext(uid)
	validID := uuid.NewString()

	t.Run("Success", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("GetFile", ctx, validID, mock.Anything).Return(&models.File{}, nil)
		mockStorage.On("RemoveObject", ctx, mock.Anything).Return(nil)
		mockRepo.On("DeleteFile", ctx, validID, mock.Anything).Return(nil)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: validID})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("RepositoryDeleteFailure", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("GetFile", ctx, validID, mock.Anything).Return(&models.File{}, nil)
		mockStorage.On("RemoveObject", ctx, mock.Anything).Return(nil)
		mockRepo.On("DeleteFile", ctx, validID, mock.Anything).Return(assert.AnError)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: validID})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, statusErr.Code())
	})

	t.Run("StorageRemoveFailure", func(t *testing.T) {
		service, mockRepo, mockStorage := setupService()
		mockRepo.On("GetFile", ctx, validID, mock.Anything).Return(&models.File{}, nil)
		mockStorage.On("RemoveObject", ctx, mock.Anything).Return(assert.AnError)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: validID})
		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, statusErr.Code())
	})
}
