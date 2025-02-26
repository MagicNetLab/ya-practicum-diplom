package postgres

import (
	"errors"
	"time"
)

// TokenModel интерфейс токена
type TokenModel struct {
	value   string    `db:"value"`
	uid     string    `db:"uid"`
	refresh bool      `db:"is_refresh"`
	expired time.Time `db:"expired"`
}

// GetValue получение значения токена
func (t *TokenModel) GetValue() string {
	return t.value
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
func (t *TokenModel) GetUID() string {
	return t.uid
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
func (t *TokenModel) SetRefresh() {
	t.refresh = true
}

// IsRefresh возвращает является ли токен refresh
func (t *TokenModel) IsRefresh() bool {
	return t.refresh
}

// SetExpired установка времени истечения действия токена
func (t *TokenModel) SetExpired(expired time.Time) error {
	if expired.IsZero() || expired.Before(time.Now()) {
		return errors.New("expired time is not correct")
	}
	t.expired = expired
	return nil
}

// GetExpired получение времени истечения действия токена
func (t *TokenModel) GetExpired() time.Time {
	return t.expired
}
