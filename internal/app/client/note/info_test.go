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

func TestDetailAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-id\n" // ID заметки

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем тестовую заметку для ответа
	testNote := &pb.NoteModel{
		ID:      "test-id",
		Title:   "Test Title",
		Meta:    "Test Meta",
		Content: "Test Content",
	}

	// Создаем мок клиента
	mockClient := new(mocks.NoteClient)
	// Настраиваем мок для метода Get
	mockClient.On("Get", mock.Anything, &pb.GetNoteRequest{ID: "test-id"}).Return(
		&pb.GetNoteResponse{Note: testNote}, nil,
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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Заметка: test-id")
	assert.Contains(t, output, "Заголовок: Test Title")
	assert.Contains(t, output, "Мета информация: Test Meta")
	assert.Contains(t, output, "Содержание: \nTest Content")
	mockClient.AssertExpectations(t)
}

func TestDetailAction_TokenError(t *testing.T) {
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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestDetailAction_ScanError(t *testing.T) {
	// Подготовка тестовых данных с некорректным вводом
	testInput := "" // Пустой ввод вызовет ошибку сканирования

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}
	defer restoreReadTokenFromFile()

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		mockClient := new(mocks.NoteClient)
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка чтения ID заметки")
}

func TestDetailAction_GetError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-id\n" // ID заметки

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с ошибкой получения заметки
	mockClient := new(mocks.NoteClient)
	mockClient.On("Get", mock.Anything, &pb.GetNoteRequest{ID: "test-id"}).Return(
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
		err := detailAction(ctx, cmd, mockClient)
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении заметки")
	mockClient.AssertExpectations(t)
}

func TestDetailAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := "test-id\n" // ID заметки

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с ошибкой истекшего токена
	mockClient := new(mocks.NoteClient)
	mockClient.On("Get", mock.Anything, &pb.GetNoteRequest{ID: "test-id"}).Return(
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
		assert.NoError(t, err) // В этой функции всегда возвращается nil
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	// Проверяем, что файл токена был удален
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
	mockClient.AssertExpectations(t)
}
