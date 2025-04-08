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

func TestSearchAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-card"

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("Search", mock.Anything, &pb.SearchCardRequest{Name: "test-card"}).Return(
		&pb.SearchCardResponse{Cards: []*pb.ShortCardModel{{ID: "1", Name: "test-card", Number: "1234-5678-9012-3456", Meta: "test description"}}}, nil,
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
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "test-card")
	assert.Contains(t, output, "1234-5678-9012-3456")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_EmptyName(t *testing.T) {
	// Подготовка тестовых данных с пустым вводом
	testInput := ""

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнено имя карты")
}

func TestSearchAction_TokenError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-card"

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

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
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestSearchAction_SearchError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-card"

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("Search", mock.Anything, &pb.SearchCardRequest{Name: "test-card"}).Return(
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
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка получения списка карт")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_EmptyResult(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-card"

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("Search", mock.Anything, &pb.SearchCardRequest{Name: "test-card"}).Return(
		&pb.SearchCardResponse{Cards: []*pb.ShortCardModel{}}, nil,
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
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ничего не найдено")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-card"

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.CardClient)
	mockClient.On("Search", mock.Anything, &pb.SearchCardRequest{Name: "test-card"}).Return(
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
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
