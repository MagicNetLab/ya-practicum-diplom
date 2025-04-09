package file

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

func TestAddAction_Success(t *testing.T) {
	// Создаем временный тестовый файл
	testFileName := "test_file.txt"
	testContent := "test content"
	err := os.WriteFile(testFileName, []byte(testContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testFileName)

	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		testFileName,  // Путь к файлу
		"custom_name", // Название файла
		"test meta",   // Описание файла
		"",            // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	mockClient := new(mocks.FilesClient)
	mockClient.On("Put", mock.Anything, mock.Anything).
		Return(&pb.PutFileResponse{}, nil)

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Файл успешно добавлен")
	mockClient.AssertExpectations(t)
}

func TestAddAction_TokenError(t *testing.T) {
	// Мокаем функцию чтения токена с ошибкой
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "", fmt.Errorf("token error")
	}

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		mockClient := new(mocks.FilesClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке токена
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestAddAction_FileNotFound(t *testing.T) {
	// Подготовка тестовых данных с несуществующим файлом
	testInput := strings.Join([]string{
		"non_existent_file.txt", // Несуществующий файл
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		mockClient := new(mocks.FilesClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Файл не найден")
}

func TestAddAction_PutError(t *testing.T) {
	// Создаем временный тестовый файл
	testFileName := "test_file.txt"
	testContent := "test content"
	err := os.WriteFile(testFileName, []byte(testContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testFileName)

	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		testFileName, // Путь к файлу
		"",           // Название файла (используется имя файла по умолчанию)
		"test meta",  // Описание файла
		"",           // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой создания
	mockClient := new(mocks.FilesClient)
	mockClient.On("Put", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Internal, "internal error"))

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке создания
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при добавлении файла")
}

func TestAddAction_ExpiredToken(t *testing.T) {
	// Создаем временный тестовый файл
	testFileName := "test_file.txt"
	testContent := "test content"
	err := os.WriteFile(testFileName, []byte(testContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testFileName)

	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		testFileName, // Путь к файлу
		"",           // Название файла (используется имя файла по умолчанию)
		"test meta",  // Описание файла
		"",           // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой истекшего токена
	mockClient := new(mocks.FilesClient)
	mockClient.On("Put", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Unauthenticated, "token expired"))

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err = os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке токена
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
}
