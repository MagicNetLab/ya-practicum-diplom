package postgres

import (
	"errors"
	"time"
)

// TokenModel интерфейс токена
type TokenModel struct {
	value   string `sql:"value"`
	uid     string `sql:"uid"`
	refresh bool   `sql:"is_refresh"`
	expired string `sql:"expired"`
}

// GetValue получение значения токена
func (t TokenModel) GetValue() (string, error) {
	if t.value == "" {
		return "", errors.New("value is not set")
	}
	return t.value, nil
}

// SetValue установка значения токена
func (t *TokenModel) SetValue(value string) error {
	if value == "" {
		return errors.New("value is not be empty")
	}
	t.value = value
	return nil
}

// GetUID получение UID пользователя для которого выдан токен
func (t TokenModel) GetUID() (string, error) {
	if t.uid == "" {
		return "", errors.New("uid is not set")
	}
	return t.uid, nil
}

// SetUID установка UID пользователя для которого выдан токен
func (t *TokenModel) SetUID(uid string) error {
	if uid == "" {
		return errors.New("uid is not be empty")
	}
	t.uid = uid
	return nil
}

// SetRefresh маркировка токена как refresh
func (t *TokenModel) SetRefresh() error {
	t.refresh = true
	return nil
}

// IsRefresh возвращает является ли токен refresh
func (t TokenModel) IsRefresh() (bool, error) {
	return t.refresh, nil
}

// SetExpired установка времени истечения действия токена
func (t *TokenModel) SetExpired(expired time.Time) error {
	if expired.IsZero() || expired.Before(time.Now()) {
		return errors.New("expired time is not correct")
	}
	t.expired = expired.Format(time.RFC3339)
	return nil
}

// GetExpired получение времени истечения действия токена
func (t TokenModel) GetExpired() (time.Time, error) {
	if t.expired == "" {
		return time.Time{}, errors.New("expired time is not set")
	}

	val, err := time.Parse(time.RFC3339, t.expired)
	if err != nil {
		return time.Time{}, err
	}
	return val, nil
}
