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

func TestRemoveAction(t *testing.T) {
	// Сохраняем оригинальную функцию и восстанавливаем ее после тестов
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()

	// Создаем контекст и команду для тестов
	ctx := context.Background()
	cmd := &cli.Command{}

	t.Run("Success_RemoveFile", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал успешный ответ
		mockClient.On("Remove", mock.Anything, mock.Anything).Return(
			&pb.RemoveFileResponse{}, nil,
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "file_id")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := removeAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит сообщение об успешном удалении
		assert.Contains(t, output, "Файл успешно удален")

		// Проверяем, что мок был вызван с правильными параметрами
		mockClient.AssertCalled(t, "Remove", mock.Anything, &pb.RemoveFileRequest{Id: "file_id"})
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
			err := removeAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при получении токена")

		// Проверяем, что мок не вызывался
		mockClient.AssertNotCalled(t, "Remove")
	})

	t.Run("Error_ReadFileID", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Эмулируем ввод пользователя с ошибкой
		restore := mockStdin(t, "")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := removeAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка: не заполнено имя карты")

		// Проверяем, что мок не вызывался
		mockClient.AssertNotCalled(t, "Remove")
	})

	t.Run("Error_RemoveFailed", func(t *testing.T) {
		// Подменяем функцию чтения токена
		readTokenFromFile = func() (string, error) {
			return "valid_token", nil
		}

		// Создаем мок клиента
		mockClient := new(mocks.FilesClient)

		// Настраиваем мок, чтобы он возвращал ошибку
		mockClient.On("Remove", mock.Anything, mock.Anything).Return(
			nil, errors.New("remove error"),
		)

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "file_id")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := removeAction(ctx, cmd, mockClient)
			assert.NoError(t, err)
		})

		// Проверяем, что вывод содержит ожидаемое сообщение об ошибке
		assert.Contains(t, output, "Ошибка при удалении файла")

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
		mockClient.On("Remove", mock.Anything, mock.Anything).Return(
			nil, status.Error(codes.Unauthenticated, "token expired"),
		)

		// Создаем временный файл token.txt для проверки его удаления
		tokenFile := "token.txt"
		f, err := os.Create(tokenFile)
		assert.NoError(t, err)
		f.Close()

		// Эмулируем ввод пользователя
		restore := mockStdin(t, "file_id")
		defer restore()

		// Перехватываем вывод в stdout
		output := captureStdoutOutput(func() {
			err := removeAction(ctx, cmd, mockClient)
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
}
