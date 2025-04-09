package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
)

func TestRegisterAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"newuser",  // Имя пользователя
		"password", // Пароль
		"",         // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AuthClient)
	mockClient.On("Register", mock.Anything, mock.Anything).
		Return(&pb.RegResponse{Token: "test-token"}, nil)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := registerAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Регистрация прошла успешно")
	mockClient.AssertExpectations(t)
}

func TestRegisterAction_EmptyUsername(t *testing.T) {
	// Подготовка тестовых данных с пустым именем пользователя
	testInput := strings.Join([]string{
		"", // Пустое имя пользователя
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AuthClient)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := registerAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не указанно имя пользователя")
}

func TestRegisterAction_EmptyPassword(t *testing.T) {
	// Подготовка тестовых данных с пустым паролем
	testInput := strings.Join([]string{
		"newuser", // Имя пользователя
		"",        // Пустой пароль
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента
	mockClient := new(mocks.AuthClient)

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := registerAction(ctx, cmd, mockClient)
		assert.Nil(t, err) // Функция возвращает nil при ошибке ввода
	})

	// Проверяем результаты
	assert.Contains(t, output, "Ошибка: не указан пароль")
}

func TestRegisterAction_RegisterError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"newuser",  // Имя пользователя
		"password", // Пароль
		"",         // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Создаем мок клиента с ошибкой регистрации
	mockClient := new(mocks.AuthClient)
	mockClient.On("Register", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.AlreadyExists, "user already exists"))

	// Выполняем тест и перехватываем вывод
	output := captureStdoutOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := registerAction(ctx, cmd, mockClient)
		assert.Error(t, err) // Функция возвращает ошибку при ошибке регистрации
	})

	// Проверяем результаты
	assert.Contains(t, output, "ошибка регистрации")
	mockClient.AssertExpectations(t)
}
