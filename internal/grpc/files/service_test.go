package files

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	rm "github.com/MagicNetLab/ya-practicum-diplom/internal/repository/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	sm "github.com/MagicNetLab/ya-practicum-diplom/internal/services/s3/mocks"
)

type mockJWTConfigurator struct {
	mock.Mock
}

func (m *mockJWTConfigurator) GetJWTSecret() string {
	return "hfjhshfjshdkfhsjhfksjdfhskfhjsdhfhaksjdh"
}

func (m *mockJWTConfigurator) IsValid() bool { return true }

func (m *mockJWTConfigurator) GetTokenLifeTime() time.Duration {
	return time.Hour
}

func (m *mockJWTConfigurator) GetRefreshTokenLifeTime() time.Duration {
	return time.Hour
}

func setupService() (*Service, *rm.FileRepository, *sm.S3Client, *mockJWTConfigurator) {
	mockRepo := new(rm.FileRepository)
	mockStorage := new(sm.S3Client)
	mockJWTCnf := new(mockJWTConfigurator)
	return &Service{db: mockRepo, storage: mockStorage, jwt: mockJWTCnf}, mockRepo, mockStorage, mockJWTCnf
}

func setupAuthContext(uid string, secret string) context.Context {
	user := &models.User{
		UID:       uid,
		Login:     "test-login",
		Password:  "test-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	token, _ := jwt.GenerateToken(user, secret)
	md := metadata.New(map[string]string{"token": token})
	return metadata.NewIncomingContext(context.Background(), md)
}

// TestService_Put тесты добавления файла в хранилище
func TestService_Put(t *testing.T) {
	err := os.Setenv("ENCRYPT_KEY", "jhjshdahjdakhdjaskdhajshdhadaskdhad")
	assert.NoError(t, err)

	t.Run("Успешное добавление файла", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		req := &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("test content"),
			Meta:    "test-meta",
		}

		mockRepo.On("CreateFile", ctx, mock.AnythingOfType("*models.File")).Return(nil)
		mockStorage.On("PutObject", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*bytes.Reader"), mock.AnythingOfType("int64")).Return(nil)

		resp, err := service.Put(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.File.Id)
		assert.Equal(t, req.Name, resp.File.Name)
		assert.Equal(t, req.Meta, resp.File.Meta)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Ошибка аутентификации", func(t *testing.T) {
		service, _, _, _ := setupService()
		ctx := context.Background()
		req := &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("test content"),
		}

		resp, err := service.Put(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Ошибка создания файла в БД", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		req := &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("test content"),
		}

		mockRepo.On("CreateFile", ctx, mock.AnythingOfType("*models.File")).Return(assert.AnError)

		resp, err := service.Put(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка загрузки файла в хранилище", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		req := &pb.PutFileRequest{
			Name:    "test.txt",
			Content: []byte("test content"),
			Meta:    "test-meta",
		}

		mockRepo.On("CreateFile", ctx, mock.AnythingOfType("*models.File")).Return(nil)
		mockStorage.On("PutObject", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*bytes.Reader"), mock.AnythingOfType("int64")).Return(assert.AnError)
		mockRepo.On("DeleteFile", ctx, mock.AnythingOfType("string"), uid).Return(nil)

		resp, err := service.Put(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}

// TestService_List тесты получения списка файлов
func TestService_List(t *testing.T) {
	err := os.Setenv("ENCRYPT_KEY", "jhjshdahjdakhdjaskdhajshdhadaskdhad")
	assert.NoError(t, err)

	t.Run("Успешное получение списка файлов", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		expectedFiles := make([]models.FilesModel, 0)

		expectedFiles = append(expectedFiles, &models.File{
			ID:        uuid.New().String(),
			UID:       uid,
			Name:      "test1.txt",
			Path:      "/test1.txt",
			Meta:      "test-meta-1",
			Size:      100,
			CreatedAt: time.Now(),
		})
		expectedFiles = append(expectedFiles, &models.File{
			ID:        uuid.New().String(),
			UID:       uid,
			Name:      "test2.txt",
			Path:      "/test2.txt",
			Meta:      "test-meta-2",
			Size:      200,
			CreatedAt: time.Now(),
		})

		mockRepo.On("SearchFile", ctx, mock.AnythingOfType("*models.FilesSearch")).Return(expectedFiles, nil)

		resp, err := service.List(ctx, &pb.ListFilesRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Files, 2)
		assert.Equal(t, expectedFiles[0].GetID(), resp.Files[0].Id)
		assert.Equal(t, expectedFiles[0].GetName(), resp.Files[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка аутентификации", func(t *testing.T) {
		service, _, _, _ := setupService()
		ctx := context.Background()
		resp, err := service.List(ctx, &pb.ListFilesRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Ошибка поиска файлов", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		mockRepo.On("SearchFile", ctx, mock.AnythingOfType("*models.FilesSearch")).Return(nil, assert.AnError)

		resp, err := service.List(ctx, &pb.ListFilesRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}

// TestService_Download тесты скачивания файла
func TestService_Download(t *testing.T) {
	err := os.Setenv("ENCRYPT_KEY", "jhjshdahjdakhdjaskdhajshdhadaskdhad")
	assert.NoError(t, err)

	t.Run("Успешное скачивание файла", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		expectedFile := &models.File{
			ID:        fileID,
			UID:       uid,
			Name:      "test.txt",
			Path:      "/test.txt",
			Meta:      "test-meta",
			Size:      100,
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetFile", ctx, fileID, uid).Return(expectedFile, nil)
		encryptContent, err := encryptor.EncryptData("encrypted content")
		assert.NoError(t, err)
		mockStorage.On("GetObject", ctx, expectedFile.GetPath()).Return([]byte(encryptContent), nil)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: fileID})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedFile.GetName(), resp.Name)
		assert.NotEmpty(t, resp.Content)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Ошибка аутентификации", func(t *testing.T) {
		service, _, _, _ := setupService()
		ctx := context.Background()
		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: uuid.New().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Ошибка валидации ID", func(t *testing.T) {
		service, _, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: "invalid-id"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Файл не найден", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		mockRepo.On("GetFile", ctx, fileID, uid).Return(nil, assert.AnError)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: fileID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка получения файла из хранилища", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		expectedFile := &models.File{
			ID:   fileID,
			UID:  uid,
			Path: "/test.txt",
		}

		mockRepo.On("GetFile", ctx, fileID, uid).Return(expectedFile, nil)
		mockStorage.On("GetObject", ctx, expectedFile.GetPath()).Return(nil, assert.AnError)

		resp, err := service.Download(ctx, &pb.DownloadFileRequest{Id: fileID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}

// TestService_Search тесты поиска файлов
func TestService_Search(t *testing.T) {
	err := os.Setenv("ENCRYPT_KEY", "jhjshdahjdakhdjaskdhajshdhadaskdhad")
	assert.NoError(t, err)

	t.Run("Успешный поиск файлов", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		exportedFiles := make([]models.FilesModel, 0)
		expectedFiles := append(exportedFiles, &models.File{
			ID:        uuid.New().String(),
			UID:       uid,
			Name:      "test1.txt",
			Path:      "/test1.txt",
			Meta:      "test-meta-1",
			Size:      100,
			CreatedAt: time.Now(),
		})

		mockRepo.On("SearchFile", ctx, mock.AnythingOfType("*models.FilesSearch")).Return(expectedFiles, nil)

		resp, err := service.Search(ctx, &pb.SearchFilesRequest{Name: "test"})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Files, 1)
		assert.Equal(t, expectedFiles[0].GetID(), resp.Files[0].Id)
		assert.Equal(t, expectedFiles[0].GetName(), resp.Files[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка аутентификации", func(t *testing.T) {
		service, _, _, _ := setupService()
		ctx := context.Background()
		resp, err := service.Search(ctx, &pb.SearchFilesRequest{Name: "test"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Ошибка поиска файлов", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		mockRepo.On("SearchFile", ctx, mock.AnythingOfType("*models.FilesSearch")).Return(nil, assert.AnError)

		resp, err := service.Search(ctx, &pb.SearchFilesRequest{Name: "test"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})
}

// TestService_Remove тесты удаления файла
func TestService_Remove(t *testing.T) {
	err := os.Setenv("ENCRYPT_KEY", "jhjshdahjdakhdjaskdhajshdhadaskdhad")
	assert.NoError(t, err)

	t.Run("Успешное удаление файла", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		expectedFile := &models.File{
			ID:   fileID,
			UID:  uid,
			Path: "/test.txt",
		}

		mockRepo.On("GetFile", ctx, fileID, uid).Return(expectedFile, nil)
		mockStorage.On("RemoveObject", ctx, expectedFile.GetPath()).Return(nil)
		mockRepo.On("DeleteFile", ctx, fileID, uid).Return(nil)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: fileID})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Ошибка аутентификации", func(t *testing.T) {
		service, _, _, _ := setupService()
		ctx := context.Background()
		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: uuid.New().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, statusErr.Code())
	})

	t.Run("Ошибка валидации ID", func(t *testing.T) {
		service, _, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: "invalid-id"})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, statusErr.Code())
	})

	t.Run("Файл не найден", func(t *testing.T) {
		service, mockRepo, _, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		mockRepo.On("GetFile", ctx, fileID, uid).Return(nil, assert.AnError)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: fileID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, statusErr.Code())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка удаления файла из хранилища", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		expectedFile := &models.File{
			ID:   fileID,
			UID:  uid,
			Path: "/test.txt",
		}

		mockRepo.On("GetFile", ctx, fileID, uid).Return(expectedFile, nil)
		mockStorage.On("RemoveObject", ctx, expectedFile.GetPath()).Return(assert.AnError)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: fileID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("Ошибка удаления информации о файле из БД", func(t *testing.T) {
		service, mockRepo, mockStorage, _ := setupService()
		uid := uuid.New().String()
		ctx := setupAuthContext(uid, service.jwt.GetJWTSecret())
		fileID := uuid.New().String()
		expectedFile := &models.File{
			ID:   fileID,
			UID:  uid,
			Path: "/test.txt",
		}

		mockRepo.On("GetFile", ctx, fileID, uid).Return(expectedFile, nil)
		mockStorage.On("RemoveObject", ctx, expectedFile.GetPath()).Return(nil)
		mockRepo.On("DeleteFile", ctx, fileID, uid).Return(assert.AnError)

		resp, err := service.Remove(ctx, &pb.RemoveFileRequest{Id: fileID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		statusErr, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, statusErr.Code())
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
