package models

import (
	"errors"
	"strconv"
	"time"
)

// CardModel интерфейс модели данных карты
type CardModel interface {
	GetID() string
	SetID(id string) error
	GetUID() string
	SetUID(uid string) error
	GetName() string
	SetName(name string) error
	GetNumber() string
	SetNumber(number string) error
	GetMask() string
	SetMask(mask string) error
	GetMonth() int
	SetMonth(month int) error
	GetYear() int
	SetYear(year int) error
	GetCVC() string
	SetCVC(cvc string) error
	GetPIN() string
	SetPIN(pin string) error
	GetMeta() string
	SetMeta(meta string) error
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time) error
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

// SetID устанавливает идентификатор карты
func (c *Card) SetID(id string) error {
	c.ID = id
	return nil
}

// GetUID возвращает идентификатор владельца карты
func (c *Card) GetUID() string {
	return c.UID
}

// SetUID устанавливает идентификатор
func (c *Card) SetUID(uid string) error {
	c.UID = uid
	return nil
}

// GetName возвращает имя владельца карты
func (c *Card) GetName() string {
	return c.Name
}

// SetName устанавливает имя владельца карты
func (c *Card) SetName(name string) error {
	c.Name = name
	return nil
}

// GetNumber возвращает номер карты
func (c *Card) GetNumber() string {
	return c.Number
}

// SetNumber устанавливает номер карты
func (c *Card) SetNumber(number string) error {
	c.Number = number
	return nil
}

// GetMask возвращает маску номера карты
func (c *Card) GetMask() string {
	return c.Mask
}

// SetMask устанавливает маску номера карты
func (c *Card) SetMask(mask string) error {
	c.Mask = mask
	return nil
}

// GetMonth возвращает месяц действия карты
func (c *Card) GetMonth() int {
	return c.Month
}

// SetMonth устанавливает месяц
func (c *Card) SetMonth(month int) error {
	c.Month = month
	return nil
}

// GetYear возвращает год действия карты
func (c *Card) GetYear() int {
	return c.Year
}

// SetYear устанавливает год
func (c *Card) SetYear(year int) error {
	c.Year = year
	return nil
}

// GetCVC возвращает код карты
func (c *Card) GetCVC() string {
	return c.CVC
}

// SetCVC устанавливает код карты
func (c *Card) SetCVC(cvc string) error {
	c.CVC = cvc
	return nil
}

// GetPIN возвращает пин карты
func (c *Card) GetPIN() string {
	return c.PIN
}

// SetPIN устанавливает пин карты
func (c *Card) SetPIN(pin string) error {
	c.PIN = pin
	return nil
}

// GetCreatedAt возвращает дату создания карты
func (c *Card) GetCreatedAt() time.Time {
	return c.CreatedAt
}

// SetCreatedAt устанавливает дату создания карты
func (c *Card) SetCreatedAt(createdAt time.Time) error {
	c.CreatedAt = createdAt
	return nil
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

// GetMeta возвращает мета данные карты
func (c *Card) GetMeta() string {
	return c.Meta
}

// SetMeta устанавливает мета данные карты
func (c *Card) SetMeta(meta string) error {
	c.Meta = meta
	return nil
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
