package accounts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockStdin эмуляция ввода данных пользователем
func mockStdin(t *testing.T, input string) func() {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	os.Stdin = r

	return func() {
		os.Stdin = oldStdin
	}
}

// captureOutput перехватывает вывод в stdout
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = oldStdout
	return <-outC
}

func TestAddAction_Success(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"example.com", // URL
		"testuser",    // Login
		"testpass",    // Password
		"test meta",   // Meta information
		"",            // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Create", mock.Anything, &pb.CreateAccountRequest{
		Login:       "testuser",
		Password:    "testpass",
		Url:         "example.com",
		Description: "test meta",
	}).Return(&pb.CreateAccountResponse{}, nil)

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.NoError(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Аккаунт успешно добавлен!")
	mockClient.AssertExpectations(t)
}

func TestAddAction_TokenError(t *testing.T) {
	readTokenFromFile = func() (string, error) {
		return "", fmt.Errorf("token error")
	}

	var output string
	output = captureOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		mockClient := new(mocks.AccountsClient)
		err := addAction(ctx, cmd, mockClient)
		assert.Error(t, err)
		assert.Equal(t, "token error", err.Error())
	})

	assert.Contains(t, output, "Ошибка при получении токена")
}

func TestAddAction_CreateError(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"example.com", // URL
		"testuser",    // Login
		"testpass",    // Password
		"test meta",   // Meta information
		"",            // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой создания
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Internal, "internal error"))

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "valid-token", nil
	}

	// Выполняем тест и перехватываем вывод
	output := captureOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	assert.Contains(t, output, "Ошибка при добавлении аккаунта")
}

func TestAddAction_ExpiredToken(t *testing.T) {
	// Подготовка тестовых данных
	testInput := strings.Join([]string{
		"example.com", // URL
		"testuser",    // Login
		"testpass",    // Password
		"test meta",   // Meta information
		"",            // Новая строка в конце
	}, "\n")

	// Мокаем ввод пользователя
	cleanup := mockStdin(t, testInput)
	defer cleanup()

	// Мокаем клиент с ошибкой истекшего токена
	mockClient := new(mocks.AccountsClient)
	mockClient.On("Create", mock.Anything, mock.Anything).
		Return(nil, status.Error(codes.Unauthenticated, "token expired"))

	// Мокаем функцию чтения токена
	readTokenFromFile = func() (string, error) {
		return "expired-token", nil
	}

	// Создаем временный файл токена
	tmpTokenFile := "token.txt"
	err := os.WriteFile(tmpTokenFile, []byte("expired-token"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tmpTokenFile)

	// Выполняем тест и перехватываем вывод
	output := captureOutput(func() {
		ctx := context.Background()
		cmd := &cli.Command{}
		err := addAction(ctx, cmd, mockClient)
		assert.Error(t, err)
	})

	// Проверяем результаты
	assert.Contains(t, output, "Время жизни токена истекло")
	_, err = os.Stat(tmpTokenFile)
	assert.True(t, os.IsNotExist(err), "Token file should be removed")
}
