package config

import "time"

// JWTConfigurator интерфейс конфигурации JWT
type JWTConfigurator interface {
	GetJWTSecret() string
	IsValid() bool
	GetTokenLifeTime() time.Duration
	GetRefreshTokenLifeTime() time.Duration
}

// JWTConfig конфигурация JWT
type JWTConfig struct {
	secret               string `env:"JWT_SECRET" default:""`
	tokenLifeTime        int    `env:"JWT_TOKEN_LIFE_TIME" default:"1"`
	refreshTokenLifeTime int    `env:"JWT_REFRESH_TOKEN_LIFE_TIME" default:"24"`
}

// GetJWTSecret возвращает секрет JWT
func (j *JWTConfig) GetJWTSecret() string {
	return j.secret
}

// IsValid возвращает утверждение корректности конфигурации JWT
func (j *JWTConfig) IsValid() bool {
	return j.secret != ""
}

// GetTokenLifeTime возвращает время жизни токена авторизации
func (j *JWTConfig) GetTokenLifeTime() time.Duration {
	return time.Duration(j.tokenLifeTime) * time.Hour
}

// GetRefreshTokenLifeTime возвращает время жизни токена обновления
func (j *JWTConfig) GetRefreshTokenLifeTime() time.Duration {
	return time.Duration(j.refreshTokenLifeTime) * time.Hour
}
