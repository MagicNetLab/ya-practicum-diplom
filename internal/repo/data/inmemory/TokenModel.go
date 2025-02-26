package inmemory

import (
	"errors"
	"time"
)

// TokenModel интерфейс токена
type TokenModel struct {
	Value   string `json:"value"`
	UID     string `json:"uid"`
	Refresh bool   `json:"is_refresh"`
	Expired string `json:"expired"`
}

// GetValue получение значения токена
func (t *TokenModel) GetValue() string {
	return t.Value
}

// SetValue установка значения токена
func (t *TokenModel) SetValue(value string) error {
	if value == "" {
		return errors.New("value is not be empty")
	}
	t.Value = value
	return nil
}

// GetUID получение UID пользователя для которого выдан токен
func (t *TokenModel) GetUID() string {
	return t.UID
}

// SetUID установка UID пользователя для которого выдан токен
func (t *TokenModel) SetUID(uid string) error {
	if uid == "" {
		return errors.New("uid is not be empty")
	}
	t.UID = uid
	return nil
}

// SetRefresh маркировка токена как refresh
func (t *TokenModel) SetRefresh() {
	t.Refresh = true
}

// IsRefresh возвращает является ли токен refresh
func (t *TokenModel) IsRefresh() bool {
	return t.Refresh
}

// SetExpired установка времени истечения действия токена
func (t *TokenModel) SetExpired(expired time.Time) error {
	if expired.IsZero() || expired.Before(time.Now()) {
		return errors.New("expired time is not correct")
	}
	t.Expired = expired.Format(time.RFC3339)
	return nil
}

// GetExpired получение времени истечения действия токена
func (t *TokenModel) GetExpired() time.Time {
	val, err := time.Parse(time.RFC3339, t.Expired)
	if err != nil {
		t.Expired = time.Now().Add(time.Hour * 3).Format(time.RFC3339)
		return time.Now().Add(time.Hour * 3)
	}
	return val
}
