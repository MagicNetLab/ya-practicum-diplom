package note

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSearchAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test search\n" // Поисковый запрос

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем тестовые заметки для ответа
	testNotes := []*pb.NoteModel{
		{
			ID:      "test-id-1",
			Title:   "Test Title 1",
			Meta:    "Test Meta 1",
			Content: "Test Content 1",
		},
		{
			ID:      "test-id-2",
			Title:   "Test Title 2",
			Meta:    "Test Meta 2",
			Content: "Test Content 2",
		},
	}

	// Создаем мок клиента
	mockClient := new(mocks.NoteClient)
	// Настраиваем мок для метода Search
	mockClient.On("Search", mock.Anything, &pb.SearchNoteRequest{Search: "test search"}).Return(
		&pb.SearchNoteResponse{Notes: testNotes}, nil,
	)

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
	assert.Contains(t, output, "test-id-1")
	assert.Contains(t, output, "Test Title 1")
	assert.Contains(t, output, "Test Meta 1")
	assert.Contains(t, output, "test-id-2")
	assert.Contains(t, output, "Test Title 2")
	assert.Contains(t, output, "Test Meta 2")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_TokenError(t *testing.T) {
	// Мокаем функцию чтения токена с ошибкой
	readTokenFromFile = func() (string, error) {
		return "", fmt.Errorf("token error")
	}
	defer restoreReadTokenFromFile()

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		mockClient := new(mocks.NoteClient)
		err := searchAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestSearchAction_EmptyResult(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test search\n" // Поисковый запрос

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с пустым результатом поиска
	mockClient := new(mocks.NoteClient)
	mockClient.On("Search", mock.Anything, &pb.SearchNoteRequest{Search: "test search"}).Return(
		&pb.SearchNoteResponse{Notes: []*pb.NoteModel{}}, nil,
	)

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
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ничего не найдено")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_SearchError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test search\n" // Поисковый запрос

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с ошибкой поиска
	mockClient := new(mocks.NoteClient)
	mockClient.On("Search", mock.Anything, &pb.SearchNoteRequest{Search: "test search"}).Return(
		nil, status.Error(codes.Internal, "internal error"),
	)

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
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении заметок")
	mockClient.AssertExpectations(t)
}

func TestSearchAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test search\n" // Поисковый запрос

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с ошибкой истекшего токена
	mockClient := new(mocks.NoteClient)
	mockClient.On("Search", mock.Anything, &pb.SearchNoteRequest{Search: "test search"}).Return(
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
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
