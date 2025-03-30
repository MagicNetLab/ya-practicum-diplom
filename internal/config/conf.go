package config

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

const (
	defaultTokenLifeTime        = 1
	defaultRefreshTokenLifeTime = 24
)

// InitConfiguration инициализирует конфигурацию приложения
func InitConfiguration() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}

// GetAppConfig возвращает параметры приложения
func GetAppConfig() AppConfigurator {
	return &AppConfig{
		serverConf: GetServerConfig(),
		dbConf:     GetDBConfig(),
		s3Conf:     GetS3Config(),
		jwtConf:    GetJWTConfig(),
		secretConf: GetSecretConfig(),
	}
}

// GetServerConfig возвращает секреты для приложения
func GetServerConfig() ServerConfigurator {
	return &ServerConfig{
		host: os.Getenv("SERVER_HOST"),
		port: os.Getenv("SERVER_PORT"),
	}
}

// GetDBConfig возвращает параметры для подключения к БД
func GetDBConfig() DataBaseConfigurator {
	return DatabaseConfig{
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASSWORD"),
		dbname:   os.Getenv("DB_NAME"),
		sslMode:  os.Getenv("DB_SSL_MODE"),
	}
}

// GetS3Config возвращает параметры для работы с S3
func GetS3Config() S3Configurator {
	return S3Config{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		Bucket:    os.Getenv("S3_BUCKET"),
	}
}

// GetJWTConfig возвращает параметры для генерации JWT токенов
func GetJWTConfig() JWTConfigurator {
	tokenLifeTime := os.Getenv("JWT_TOKEN_LIFE_TIME")
	tlf, err := strconv.Atoi(tokenLifeTime)
	if err != nil {
		logger.Error("error parsing JWT_TOKEN_LIFE_TIME", logger.StrArg("error", err.Error()))
		tlf = defaultTokenLifeTime
	}

	refreshTokenLifeTime := os.Getenv("JWT_REFRESH_TOKEN_LIFE_TIME")
	rtlf, err := strconv.Atoi(refreshTokenLifeTime)
	if err != nil {
		logger.Error("error parsing JWT_REFRESH_TOKEN_LIFE_TIME", logger.StrArg("error", err.Error()))
		rtlf = defaultRefreshTokenLifeTime
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = jwt.GetRandomSecret()
	}
	return &JWTConfig{
		secret:               secretKey,
		tokenLifeTime:        tlf,
		refreshTokenLifeTime: rtlf,
	}
}

// GetSecretConfig возвращает секрет для приложения
func GetSecretConfig() SecretConfigurator {
	return &SecretConfig{
		secretKey: os.Getenv("ENCRYPT_KEY"),
	}
}
