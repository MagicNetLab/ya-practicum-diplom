package card

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAction_Success(t *testing.T) {
	// Создаем тестовые карты для ответа
	testCards := []*pb.ShortCardModel{
		{
			ID:     "test-id-1",
			Name:   "Test Card 1",
			Number: "1234-5678-9012-3456",
			Meta:   "test meta 1",
		},
		{
			ID:     "test-id-2",
			Name:   "Test Card 2",
			Number: "9876-5432-1098-7654",
			Meta:   "test meta 2",
		},
	}

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	// Настраиваем мок для метода List
	mockClient.On("List", mock.Anything, &pb.ListCardRequest{}).Return(
		&pb.ListCardResponse{Cards: testCards}, nil,
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("valid-token"), 0644)
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
	assert.Contains(t, output, "Test Card 1")
	assert.Contains(t, output, "1234-5678-9012-3456")
	assert.Contains(t, output, "test meta 1")
	assert.Contains(t, output, "Test Card 2")
	assert.Contains(t, output, "9876-5432-1098-7654")
	assert.Contains(t, output, "test meta 2")
	mockClient.AssertExpectations(t)
}

func TestListAction_TokenError(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.CardClient)

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
		assert.Nil(t, err) // В функции listAction ошибка не возвращается при проблеме с токеном
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestListAction_EmptyList(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	// Настраиваем мок для метода List с пустым списком
	mockClient.On("List", mock.Anything, &pb.ListCardRequest{}).Return(
		&pb.ListCardResponse{Cards: []*pb.ShortCardModel{}}, nil,
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("valid-token"), 0644)
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
	assert.Contains(t, output, "У вас нет карт")
	mockClient.AssertExpectations(t)
}

func TestListAction_ExpiredToken(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("List", mock.Anything, &pb.ListCardRequest{}).Return(
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

func TestListAction_OtherError(t *testing.T) {
	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("List", mock.Anything, &pb.ListCardRequest{}).Return(
		nil, status.Error(codes.Internal, "internal error"),
	)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("valid-token"), 0644)
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
	assert.Contains(t, output, "Ошибка получения списка карт")
	mockClient.AssertExpectations(t)
}
