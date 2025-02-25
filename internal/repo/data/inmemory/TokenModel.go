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
func (t TokenModel) GetValue() (string, error) {
	if t.Value == "" {
		return "", errors.New("value is not set")
	}
	return t.Value, nil
}

// SetValue установка значения токена
func (t TokenModel) SetValue(value string) error {
	if value == "" {
		return errors.New("value is not be empty")
	}
	t.Value = value
	return nil
}

// GetUID получение UID пользователя для которого выдан токен
func (t TokenModel) GetUID() (string, error) {
	if t.UID == "" {
		return "", errors.New("uid is not set")
	}
	return t.UID, nil
}

// SetUID установка UID пользователя для которого выдан токен
func (t TokenModel) SetUID(uid string) error {
	if uid == "" {
		return errors.New("uid is not be empty")
	}
	t.UID = uid
	return nil
}

// SetRefresh маркировка токена как refresh
func (t TokenModel) SetRefresh() error {
	t.Refresh = true
	return nil
}

// IsRefresh возвращает является ли токен refresh
func (t TokenModel) IsRefresh() (bool, error) {
	return t.Refresh, nil
}

// SetExpired установка времени истечения действия токена
func (t TokenModel) SetExpired(expired time.Time) error {
	if expired.IsZero() || expired.Before(time.Now()) {
		return errors.New("expired time is not correct")
	}
	t.Expired = expired.Format(time.RFC3339)
	return nil
}

// GetExpired получение времени истечения действия токена
func (t TokenModel) GetExpired() (time.Time, error) {
	if t.Expired == "" {
		return time.Time{}, errors.New("expired time is not set")
	}

	val, err := time.Parse(time.RFC3339, t.Expired)
	if err != nil {
		return time.Time{}, err
	}
	return val, nil
}
