package accounts

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"test-id", // id
		"",        // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Remove", mock.Anything, &pb.RemoveAccountRequest{Id: "test-id"}).Return(
		&pb.RemoveAccountResponse{}, nil,
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := removeAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Аккаунт успешно удален!")
	mockClient.AssertExpectations(t)
}

func TestRemoveAction_NoID(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"", // Пустой id
		"", // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := removeAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Необходимо указать id аккаунта")
}

func TestRemoveAction_TokenError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"test-id", // id
		"",        // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)

	// Мокаем функцию чтения токена с ошибкой
	readTokenFromFile = func() (string, error) {
		return "", fmt.Errorf("token error")
	}
	defer restoreReadTokenFromFile()

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := removeAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestRemoveAction_RemoveError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"test-id", // id
		"",        // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Remove", mock.Anything, &pb.RemoveAccountRequest{Id: "test-id"}).Return(
		nil, status.Error(codes.Internal, "internal error"),
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := removeAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // Функция не возвращает ошибку в этом случае
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при удалении аккаунта.")
	mockClient.AssertExpectations(t)
}

func TestRemoveAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"test-id", // id
		"",        // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Remove", mock.Anything, &pb.RemoveAccountRequest{Id: "test-id"}).Return(
		nil, status.Error(codes.Unauthenticated, "token expired"),
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := removeAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // Функция не возвращает ошибку в этом случае
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
