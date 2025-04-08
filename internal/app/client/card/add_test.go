package card

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

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
)

func TestAddAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"Test Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
		"123",        // CVC
		"1234",       // PIN
		"test meta",  // Дополнительная информация
		"",           // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	mockClient := new(mocks.CardClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
		Return(&pb.CreateCardResponse{}, nil)

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Карта успешно добавлена")
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке токена
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestAddAction_EmptyName(t *testing.T) {
	// Подготовка тестовых данных с пустым названием
	testInput := strings.Join([]string{
		"", // Пустое название карты
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнено название карты")
}

func TestAddAction_EmptyNumber(t *testing.T) {
	// Подготовка тестовых данных с пустым номером
	testInput := "Test_Card\n"

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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнен номер карты")
}

func TestAddAction_InvalidMonth(t *testing.T) {
	// Подготовка тестовых данных с некорректным месяцем
	testInput := strings.Join([]string{
		"Test_Card",   // Название карты
		"74583658467", // Номер карты
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнен месяц")
}

func TestAddAction_InvalidYear(t *testing.T) {
	// Подготовка тестовых данных с некорректным годом
	testInput := strings.Join([]string{
		"Test_Card",  // Название карты
		"1234567890", // Номер карты
		"01",         // Месяц
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнен год")
}

func TestAddAction_EmptyCVC(t *testing.T) {
	// Подготовка тестовых данных с пустым CVC
	testInput := strings.Join([]string{
		"Test_Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнен CVC")
}

func TestAddAction_EmptyPIN(t *testing.T) {
	// Подготовка тестовых данных с пустым PIN
	testInput := strings.Join([]string{
		"Test_Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
		"342",        // CVC
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
		mockClient := new(mocks.CardClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не заполнен PIN")
}

func TestAddAction_CreateError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"Test_Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
		"123",        // CVC
		"1234",       // PIN
		"",           // Дополнительная информация
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой создания
	mockClient := new(mocks.CardClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
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
	assert.Contains(t, output, "Ошибка при добавлении карты")
}

func TestAddAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"Test Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
		"123",        // CVC
		"1234",       // PIN
		"test meta",  // Дополнительная информация
		"",           // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой истекшего токена
	mockClient := new(mocks.CardClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Unauthenticated, "token expired"))

	// Мокаем функцию чтения токена
	originalFunc := readTokenFromFile
	defer func() { readTokenFromFile = originalFunc }()
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
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

func TestAddAction_InvalidArgument(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"Test Card",  // Название карты
		"1234567890", // Номер карты
		"12",         // Месяц
		"2025",       // Год
		"123",        // CVC
		"1234",       // PIN
		"test meta",  // Дополнительная информация
		"",           // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой неверных аргументов
	mockClient := new(mocks.CardClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.InvalidArgument, "invalid argument"))

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
		assert.Nil(t, err) // Функция возвращает nil при ошибке аргументов
	})

	// Проверяем результаты
	assert.Contains(t, output, "Карта не добавлена. Проверьте правильность заполнения данных.")
}
