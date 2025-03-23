package config

// SecretConfigurator интерфейс конфигуратора секрета
type SecretConfigurator interface {
	GetSecret() (string, error)
	IsValid() bool
}

// SecretConfig содержит секрет для подписи JWT-токена
type SecretConfig struct {
	secretKey string
}

// GetSecret возвращает секрет
func (s *SecretConfig) GetSecret() (string, error) {
	return s.secretKey, nil
}

// IsValid возвращает результат проверки на корректность секрета
func (s *SecretConfig) IsValid() bool {
	return s.secretKey != ""
}
