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

func TestDetailAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"test-id", // id
		"",        // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем тестовый аккаунт для ответа
	testAccount := &pb.Account{
		Id:          "test-id",
		Uid:         "test-uid",
		Login:       "testuser",
		Password:    "testpass",
		Url:         "example.com",
		Description: "test description",
	}

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	// Настраиваем мок для метода Get
	mockClient.On("Get", mock.Anything, &pb.GetAccountRequest{Id: "test-id"}).Return(
		&pb.GetAccountResponse{Account: testAccount}, nil,
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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Данные аккаунта:")
	assert.Contains(t, output, "Url: example.com")
	assert.Contains(t, output, "Логин: testuser")
	assert.Contains(t, output, "Пароль: testpass")
	assert.Contains(t, output, "Meta: test description")
	mockClient.AssertExpectations(t)
}

func TestDetailAction_NoID(t *testing.T) {
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
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := detailAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Необходимо указать id аккаунта")
}

func TestDetailAction_TokenError(t *testing.T) {
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
		err := detailAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestDetailAction_GetError(t *testing.T) {
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

	// Настраиваем мок для метода Get с ошибкой
	mockClient.On("Get", mock.Anything, &pb.GetAccountRequest{Id: "test-id"}).Return(
		nil, status.Error(codes.Internal, "internal error"),
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // Функция не возвращает ошибку в этом случае
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении информации об аккаунте")
	mockClient.AssertExpectations(t)
}

func TestDetailAction_ExpiredToken(t *testing.T) {
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
	mockClient.On("Get", mock.Anything, &pb.GetAccountRequest{Id: "test-id"}).Return(
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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // Функция не возвращает ошибку в этом случае
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
