package accounts

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAction_Success(t *testing.T) {
	// Создаем тестовые аккаунты для ответа
	testAccounts := []*pb.Account{
		{
			Id:          "test-id-1",
			Login:       "testuser1",
			Url:         "example1.com",
			Description: "test description 1",
		},
		{
			Id:          "test-id-2",
			Login:       "testuser2",
			Url:         "example2.com",
			Description: "test description 2",
		},
	}

	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	// Настраиваем мок для метода Search
	mockClient.On("Search", mock.Anything, &pb.SearchAccountRequest{}).Return(
		&pb.SearchAccountResponse{Acc: testAccounts}, nil,
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
		err := listAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "testuser1")
	assert.Contains(t, output, "example1.com")
	assert.Contains(t, output, "test description 1")
	assert.Contains(t, output, "testuser2")
	assert.Contains(t, output, "example2.com")
	assert.Contains(t, output, "test description 2")
	mockClient.AssertExpectations(t)
}

func TestListAction_TokenError(t *testing.T) {
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
		err := listAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestListAction_EmptyList(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	// Настраиваем мок для метода Search с пустым списком
	mockClient.On("Search", mock.Anything, &pb.SearchAccountRequest{}).Return(
		&pb.SearchAccountResponse{Acc: []*pb.Account{}}, nil,
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
		err := listAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Список аккаунтов пуст")
	mockClient.AssertExpectations(t)
}

func TestListAction_ExpiredToken(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Search", mock.Anything, &pb.SearchAccountRequest{}).Return(
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
		err := listAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
