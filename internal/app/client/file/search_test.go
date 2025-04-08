package file

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

func TestSearchAction(t *testing.T) {
	// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()

	// Создаем контекст и команду для тестов
	ctx := context.Background()
	cmd := &cli.Command{}

	t.Run("Success_SearchFiles", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Создаем тестовые данные
		files := []*pb.FileModel{
			{
				Id:   "file1",
				Name: "test1.txt",
				Size: 100,
				Meta: "Test file 1",
			},
			{
				Id:   "file2",
				Name: "test2.txt",
				Size: 200,
				Meta: "Test file 2",
			},
		}

		// Настраиваем мок, чтобы он возвращал успешный ответ
		mockClient.On("Search", mock.Anything, mock.Anything).Return(
			&pb.SearchFilesResponse{Files: files}, nil,
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "test.txt\n")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемые данные
		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "Name")
		assert.Contains(t, output, "Size")
		assert.Contains(t, output, "Desc")
		assert.Contains(t, output, "file1")
		assert.Contains(t, output, "test1.txt")
		assert.Contains(t, output, "Test file 1")
		assert.Contains(t, output, "file2")
		assert.Contains(t, output, "test2.txt")
		assert.Contains(t, output, "Test file 2")

		// Проверяем, что мок был вызван с правильными параметрами
		mockClient.AssertCalled(t, "Search", mock.Anything, &pb.SearchFilesRequest{Name: "test.txt"})
	})

	t.Run("Error_ReadTokenFromFile", func(t *testing.T) {
		// Подменяем функцию чтения токена, чтобы она возвращала ошибку
		readTokenFromFile = func() (string, error) {
			return "", errors.New("token error")
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при получении токена")

		// Проверяем, что мок не вызывался
		mockClient.AssertNotCalled(t, "Search")
	})

	t.Run("Error_SearchFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал ошибку
		mockClient.On("Search", mock.Anything, mock.Anything).Return(
			nil, errors.New("search error"),
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "test.txt\n")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при поиске файла")

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})

	t.Run("Error_AuthenticationFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "expired_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал ошибку аутентификации
		mockClient.On("Search", mock.Anything, mock.Anything).Return(
			nil, status.Error(codes.Unauthenticated, "token expired"),
		)

		// Создаем временный файл token.txt для проверки его удаления
		tokenFile := "token.txt"
		f, err := os.Create(tokenFile)
		assert.NoError(t, err)
		f.Close()

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "test.txt\n")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Время жизни токена истекло")

		// Проверяем, что файл token.txt был удален
		_, err = os.Stat(tokenFile)
		assert.True(t, os.IsNotExist(err))

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})

	t.Run("Success_EmptyResult", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал пустой результат
		mockClient.On("Search", mock.Anything, mock.Anything).Return(
			&pb.SearchFilesResponse{Files: []*pb.FileModel{}}, nil,
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "test.txt\n")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение
		assert.Contains(t, output, "Ничего не найдено")

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})

	// Проверяем, что метаданные правильно передаются в запрос
	t.Run("Check_Metadata", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "test_token", nil
		}

		// Создаем мок клиента с проверкой метаданных
		mockClient := new(mocks.FilesClient)
		mockClient.On("Search", mock.MatchedBy(func(ctx context.Context) bool {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				return false
			}
			tokens := md.Get("token")
			return len(tokens) == 1 && tokens[0] == "test_token"
		}), mock.Anything).Return(
			&pb.SearchFilesResponse{Files: []*pb.FileModel{}}, nil,
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "test.txt\n")
		defer restore()

		// Выполняем функцию
		_ = captureStdoutOutput(func() {
			err := searchAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что мок был вызван с правильными метаданными
		mockClient.AssertExpectations(t)
	})
}
