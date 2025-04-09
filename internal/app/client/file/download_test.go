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
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

func TestDownloadAction(t *testing.T) {
	// Подменяем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid_token", nil
	}
	// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
	defer restoreReadTokenFromFile()

	// Создаем контекст и команду для тестов
	ctx := context.Background()
	cmd := &cli.Command{}

	t.Run("Error_ReadTokenFromFile", func(t *testing.T) {
		// Подменяем функцию чтения токена, чтобы она возвращала ошибку
		readTokenFromFile = func() (string, error) {
			return "", errors.New("token error")
		}
		defer restoreReadTokenFromFile()

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при получении токена")

		// Проверяем, что мок не вызывался
		mockClient.AssertNotCalled(t, "Download")
	})

	t.Run("Error_ReadFileID", func(t *testing.T) {
		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}
		// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
		defer restoreReadTokenFromFile()

		// Эмулируем ввод пользователя с ошибкой
		restore := mockStdin(t, "")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка чтения ID файла")

		// Проверяем, что мок не вызывался
		mockClient.AssertNotCalled(t, "Download")
	})

	t.Run("Error_AuthenticationFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}
		// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
		defer restoreReadTokenFromFile()

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал ошибку аутентификации
		mockClient.On("Download", mock.Anything, mock.Anything).Return(
			nil, status.Error(codes.Unauthenticated, "token expired"),
		)

		// Создаем временный файл token.txt для проверки его удаления
		tokenFile := "token.txt"
		f, err := os.Create(tokenFile)
		assert.NoError(t, err)
		f.Close()

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "file_id\n/path/to/")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
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

	t.Run("Error_DownloadFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}
		// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
		defer restoreReadTokenFromFile()

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал ошибку
		mockClient.On("Download", mock.Anything, mock.Anything).Return(
			nil, errors.New("download error"),
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "file_id\n/path/to/")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при скачивании файла")

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})

	t.Run("Error_CreateFileFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}
		// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
		defer restoreReadTokenFromFile()

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Создаем ответ с недопустимым путем для файла
		response := &pb.DownloadFileResponse{
			Name:    "test.txt",
			Content: []byte("test content"),
			Size:    12,
		}

		// Настраиваем мок, чтобы он возвращал успешный ответ
		mockClient.On("Download", mock.Anything, mock.Anything).Return(
			response, nil,
		)

		// Эмулируем ввод пользователя с недопустимым путем
		restore := mockStdin(t, "file_id\n/invalid/path/")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка создания файла")

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})

	t.Run("Success_DownloadFile", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}
		// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
		defer restoreReadTokenFromFile()

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Создаем временную директорию для тестов
		tempDir, err := os.MkdirTemp("", "test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Создаем ответ
		response := &pb.DownloadFileResponse{
			Name:    "test.txt",
			Content: []byte("test content"),
			Size:    12,
		}

		// Настраиваем мок, чтобы он возвращал успешный ответ
		mockClient.On("Download", mock.Anything, mock.Anything).Return(
			response, nil,
		)

		// Эмулируем ввод пользователя с путем к временной директории
		restore := mockStdin(t, "file_id\n"+tempDir+"/")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := downloadAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит сообщение об успешном скачивании
		assert.Contains(t, output, "Файл успешно скачан")

		// Проверяем, что файл был создан
		filePath := tempDir + "/test.txt"
		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		// Проверяем содержимое файла
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		// Проверяем, что мок был вызван
		mockClient.AssertExpectations(t)
	})
}
