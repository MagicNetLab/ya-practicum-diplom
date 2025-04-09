package card

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/mocks"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
)

func TestDetailAction(t *testing.T) {
	// Сохраняем оригинальную функцию и восстанавливаем после тестов
	defer restoreReadTokenFromFile()

	tests := []struct {
		name           string
		setupMock      func(*mocks.CardClient)
		setupToken     func()
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Успешное получение информации о карте",
			setupMock: func(m *mocks.CardClient) {
				m.On("Get", mock.Anything, &pb.GetCardRequest{ID: "test-id"}).Return(&pb.GetCardResponse{
					Card: &pb.CardModel{
						Name:   "Test Card",
						Number: "1234567890123456",
						Month:  12,
						Year:   2025,
						PIN:    "1234",
						CVC:    "123",
						Meta:   "Test Meta",
					},
				}, nil)
			},
			setupToken: func() {
				readTokenFromFile = func() (string, error) {
					return "test-token", nil
				}
			},
			input: "test-id\n",
			expectedOutput: "Введите Url сайта: Информация о карте:\n" +
				"Наименование: Test Card\n" +
				"Номер: 1234567890123456\n" +
				"Срок действия: 12/2025\n" +
				"PIN: 1234\n" +
				"CVC: 123\n" +
				"Дополнительная информация: Test Meta\n",
			expectedError: false,
		},
		{
			name:      "Ошибка при получении токена",
			setupMock: func(m *mocks.CardClient) {},
			setupToken: func() {
				readTokenFromFile = func() (string, error) {
					return "", errors.New("token error")
				}
			},
			input:          "",
			expectedOutput: "Ошибка при получении токена. Возможно, вы не авторизовались\n",
			expectedError:  false,
		},
		{
			name:      "Ошибка при вводе ID",
			setupMock: func(m *mocks.CardClient) {},
			setupToken: func() {
				readTokenFromFile = func() (string, error) {
					return "test-token", nil
				}
			},
			input:          "\n",
			expectedOutput: "Введите Url сайта: Необходимо указать id карты\n",
			expectedError:  true,
		},
		{
			name: "Ошибка аутентификации при получении карты",
			setupMock: func(m *mocks.CardClient) {
				m.On("Get", mock.Anything, &pb.GetCardRequest{ID: "test-id"}).Return(nil,
					status.Error(codes.Unauthenticated, "token expired"))
			},
			setupToken: func() {
				readTokenFromFile = func() (string, error) {
					return "test-token", nil
				}
			},
			input:          "test-id\n",
			expectedOutput: "Введите Url сайта: Время жизни токена истекло. Пожалуйста, повторите авторизацию.\n",
			expectedError:  false,
		},
		{
			name: "Общая ошибка при получении карты",
			setupMock: func(m *mocks.CardClient) {
				m.On("Get", mock.Anything, &pb.GetCardRequest{ID: "test-id"}).Return(nil,
					status.Error(codes.Internal, "internal error"))
			},
			setupToken: func() {
				readTokenFromFile = func() (string, error) {
					return "test-token", nil
				}
			},
			input:          "test-id\n",
			expectedOutput: "Введите Url сайта: Ошибка получения подробной информации о карте: rpc error: code = Internal desc = internal error\n",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок клиента
			mockClient := new(mocks.CardClient)
			tt.setupMock(mockClient)

			// Устанавливаем мок для чтения токена
			tt.setupToken()

			// Подменяем stdin для эмуляции ввода пользователя
			cleanupStdin := mockStdin(t, tt.input)
			defer cleanupStdin()

			// Перехватываем stdout
			output := captureStdoutOutput(func() {
				err := detailAction(context.Background(), &cli.Command{}, mockClient)
				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})

			// Проверяем вывод
			assert.Equal(t, tt.expectedOutput, output)

			// Проверяем, что все ожидаемые вызовы были выполнены
			mockClient.AssertExpectations(t)
		})
	}
}
