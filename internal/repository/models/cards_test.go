package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCard_GetID проверка получения ID карты.
func TestCard_GetID(t *testing.T) {
	t.Run("Получение заполненого ID", func(t *testing.T) {
		card := &Card{ID: "test-id"}
		assert.Equal(t, "test-id", card.GetID())
	})

	t.Run("Получение пустого ID", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetID())
	})
}

// TestCard_GetUID проверка получения UID карты.
func TestCard_GetUID(t *testing.T) {
	t.Run("Получение заполненого UID", func(t *testing.T) {
		card := &Card{UID: "test-uid"}
		assert.Equal(t, "test-uid", card.GetUID())
	})

	t.Run("Получение пустого UID", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetUID())
	})
}

// TestCard_GetName проверка получения имени владельца карты.
func TestCard_GetName(t *testing.T) {
	t.Run("Получение заполненого имени", func(t *testing.T) {
		card := &Card{Name: "test name"}
		assert.Equal(t, "test name", card.GetName())
	})

	t.Run("Получение пустого имени", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetName())
	})
}

// TestCard_GetNumber проверка получения номера карты.
func TestCard_GetNumber(t *testing.T) {
	t.Run("Получение заполненого номера", func(t *testing.T) {
		card := &Card{Number: "1234567890123456"}
		assert.Equal(t, "1234567890123456", card.GetNumber())
	})

	t.Run("Получение пустого номера", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetNumber())
	})
}

// TestCard_GetMask проверка получения маски карты.
func TestCard_GetMask(t *testing.T) {
	t.Run("Получение заполненой маски", func(t *testing.T) {
		card := &Card{Mask: "453201******0366"}
		assert.Equal(t, "453201******0366", card.GetMask())
	})

	t.Run("Получение пустой маски", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetMask())
	})
}

// TestCard_GetMonth проверка получения месяца окончания действия карты.
func TestCard_GetMonth(t *testing.T) {
	t.Run("Получение заполненого месяца", func(t *testing.T) {
		card := &Card{Month: 12}
		assert.Equal(t, 12, card.GetMonth())
	})

	t.Run("Получение пустого месяца", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, 0, card.GetMonth())
	})
}

// TestCard_GetYear проверка получения года окончания действия карты.
func TestCard_GetYear(t *testing.T) {
	t.Run("Получение заполненого года", func(t *testing.T) {
		card := &Card{Year: 2025}
		assert.Equal(t, 2025, card.GetYear())
	})

	t.Run("Получение пустого года", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, 0, card.GetYear())
	})
}

// TestCard_GetCVC проверка получения CVC карты.
func TestCard_GetCVC(t *testing.T) {
	t.Run("Получение заполненого CVC", func(t *testing.T) {
		card := &Card{CVC: "123"}
		assert.Equal(t, "123", card.GetCVC())
	})

	t.Run("Получение пустого CVC", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetCVC())
	})
}

// TestCard_GetPin проверка получения PIN карты.
func TestCard_GetPin(t *testing.T) {
	t.Run("Получение заполненого PIN", func(t *testing.T) {
		card := &Card{PIN: "1234"}
		assert.Equal(t, "1234", card.GetPIN())
	})

	t.Run("Получение пустого PIN", func(t *testing.T) {
		card := &Card{}
		assert.Equal(t, "", card.GetPIN())
	})
}

// TestCard_Validate проверка корректности установленных данных по карте.
func TestCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    Card
		wantErr bool
		errMsg  string
	}{
		{
			name: "Успешная проверка",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "4532015112830366", // Valid Visa card number
			},
			wantErr: false,
		},
		{
			name: "Проверка карты с некорректным месяцем истечения срока действия (< 1)",
			card: Card{
				Month:  0,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid month value",
		},
		{
			name: "Проверка карты с некорректным месяцем истечения срока действия (> 12)",
			card: Card{
				Month:  13,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid month value",
		},
		{
			name: "Проверка карты с некорректно установленным годом истечения срока действия (< 1970)",
			card: Card{
				Month:  12,
				Year:   1969,
				CVC:    "123",
				PIN:    "1234",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid year value",
		},
		{
			name: "Проверка карты с некорректным CVC: меньше 3 символов",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "12",
				PIN:    "1234",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid cvc value",
		},
		{
			name: "Проверка карты с некорректным CVC: больше 4 символов",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "1234",
				PIN:    "1234",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid cvc value",
		},
		{
			name: "Проверка карты с некорректным PIN: меньше 4 символов",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "123",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid pin value",
		},
		{
			name: "Проверка карты с некорректным PIN: больше 7 символов",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234567",
				Number: "4532015112830366",
			},
			wantErr: true,
			errMsg:  "invalid pin value",
		},
		{
			name: "Проверка карты с некорректным номером: количество знаков меньше нужного",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "123456789012",
			},
			wantErr: true,
			errMsg:  "invalid number value",
		},
		{
			name: "Проверка карты с некорректным номером: количество знаков больше нужного",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "12345678901234567890",
			},
			wantErr: true,
			errMsg:  "invalid number value",
		},
		{
			name: "Проверка карты с некорректным номером: не проходит проверка алгоритма Луна",
			card: Card{
				Month:  12,
				Year:   2025,
				CVC:    "123",
				PIN:    "1234",
				Number: "4532015112830367",
			},
			wantErr: true,
			errMsg:  "invalid number value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
