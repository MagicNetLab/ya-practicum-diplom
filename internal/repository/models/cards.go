package models

import (
	"errors"
	"strconv"
	"time"
)

// CardModel интерфейс модели данных карты
type CardModel interface {
	GetID() string
	GetUID() string
	GetName() string
	GetNumber() string
	GetMask() string
	GetMonth() int
	GetYear() int
	GetCVC() string
	GetPIN() string
	GetMeta() string
	GetCreatedAt() time.Time
	MaskNumber()
	Validate() error
}

// Card модель данных карты
type Card struct {
	ID        string    `db:"id"`
	UID       string    `db:"uid"`
	Name      string    `db:"name"`
	Number    string    `db:"number"`
	Mask      string    `db:"mask"`
	Month     int       `db:"month"`
	Year      int       `db:"year"`
	CVC       string    `db:"cvc"`
	PIN       string    `db:"pin"`
	Meta      string    `db:"meta"`
	CreatedAt time.Time `db:"created_at"`
}

// GetID возвращает идентификатор карты
func (c *Card) GetID() string {
	return c.ID
}

// GetUID возвращает идентификатор владельца карты
func (c *Card) GetUID() string {
	return c.UID
}

// GetName возвращает имя владельца карты
func (c *Card) GetName() string {
	return c.Name
}

// GetNumber возвращает номер карты
func (c *Card) GetNumber() string {
	return c.Number
}

// GetMask возвращает маску номера карты
func (c *Card) GetMask() string {
	return c.Mask
}

// GetMonth возвращает месяц действия карты
func (c *Card) GetMonth() int {
	return c.Month
}

// GetYear возвращает год действия карты
func (c *Card) GetYear() int {
	return c.Year
}

// GetCVC возвращает код карты
func (c *Card) GetCVC() string {
	return c.CVC
}

// GetPIN возвращает пин карты
func (c *Card) GetPIN() string {
	return c.PIN
}

// GetCreatedAt возвращает дату создания карты
func (c *Card) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// MaskNumber is retrieving card number mask.
func (c *Card) MaskNumber() {
	mask := ""

	for i := 0; i < len(c.Number)-4; i++ {
		mask += "*"

		if i%4 == 3 {
			mask += " "
		}
	}

	c.Mask = mask + c.Number[len(c.Number)-4:]
}

func (c *Card) GetMeta() string {
	return c.Meta
}

// Validate проверка данных карты
func (c *Card) Validate() error {
	if c.Month < 1 || 12 < c.Month {
		return errors.New("invalid month value")
	}

	if c.Year < 1970 {
		return errors.New("invalid year value")
	}

	if len(c.CVC) != 3 {
		return errors.New("invalid cvc value")
	}

	if len(c.PIN) < 4 || len(c.PIN) > 6 {
		return errors.New("invalid pin value")
	}

	var sum int
	var alternate bool

	numberLen := len(c.Number)

	if numberLen < 13 || numberLen > 19 {
		return errors.New("invalid number value")
	}

	for i := numberLen - 1; i > -1; i-- {
		mod, _ := strconv.Atoi(string(c.Number[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate

		sum += mod
	}

	if !(sum%10 == 0) {
		return errors.New("invalid number value")
	}

	return nil
}
