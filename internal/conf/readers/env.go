package readers

import (
	"errors"
	"os"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

type EnvReader struct {
	serverHost              string
	serverPort              string
	s3Endpoint              string
	s3SecretKey             string
	s3AccessKey             string
	s3Bucket                string
	dbHost                  string
	dbPort                  string
	dbUser                  string
	dbPassword              string
	dbName                  string
	dbSSlMode               string
	jwtSecret               string
	jwtTokenLifeTime        string
	jwtRefreshTokenLifeTime string
	encryptKey              string
}

// Parse читает установленные параметры конфигурации
func (r *EnvReader) Parse() error {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("error loading .env file", logger.StrArg("error", err.Error()))
	}

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost != "" {
		r.serverHost = serverHost
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort != "" {
		r.serverPort = serverPort
	}

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint != "" {
		r.s3Endpoint = s3Endpoint
	}

	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	if s3SecretKey != "" {
		r.s3SecretKey = s3SecretKey
	}

	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	if s3AccessKey != "" {
		r.s3AccessKey = s3AccessKey
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket != "" {
		r.s3Bucket = s3Bucket
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		r.dbHost = dbHost
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort != "" {
		r.dbPort = dbPort
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser != "" {
		r.dbUser = dbUser
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword != "" {
		r.dbPassword = dbPassword
	}

	dbName := os.Getenv("DB_NAME")
	if dbName != "" {
		r.dbName = dbName
	}

	dbSSLMode := os.Getenv("DB_SSL_MODE")
	if dbSSLMode != "" {
		r.dbSSlMode = dbSSLMode
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		r.jwtSecret = jwtSecret
	}

	jwtTokenLifeTime := os.Getenv("JWT_TOKEN_LIFE_TIME")
	if jwtTokenLifeTime != "" {
		r.jwtTokenLifeTime = jwtTokenLifeTime
	}

	jwtRefreshTokenLifeTime := os.Getenv("JWT_REFRESH_TOKEN_LIFE_TIME")
	if jwtRefreshTokenLifeTime != "" {
		r.jwtRefreshTokenLifeTime = jwtRefreshTokenLifeTime
	}

	encryptKey := os.Getenv("ENCRYPT_KEY")
	if encryptKey != "" {
		r.encryptKey = encryptKey
	}

	return nil
}

// GetServerHost возвращает адрес сервера
func (r *EnvReader) GetServerHost() (string, error) {
	if r.serverHost == "" {
		return "", errors.New("serverHost is not set")
	}
	return r.serverHost, nil
}

// GetServerPort возвращает порт сервера
func (r *EnvReader) GetServerPort() (string, error) {
	if r.serverPort == "" {
		return "", errors.New("serverPort is not set")
	}
	return r.serverPort, nil
}

// GetS3EndPoint возвращает адрес S3
func (r *EnvReader) GetS3EndPoint() (string, error) {
	if r.s3Endpoint == "" {
		return "", errors.New("s3Endpoint is not set")
	}
	return r.s3Endpoint, nil
}

// GetS3AccessKey возвращает ключ S3
func (r *EnvReader) GetS3AccessKey() (string, error) {
	if r.s3AccessKey == "" {
		return "", errors.New("s3AccessKey is not set")
	}
	return r.s3AccessKey, nil
}

// GetS3SecretKey возвращает ключ S3
func (r *EnvReader) GetS3SecretKey() (string, error) {
	if r.s3SecretKey == "" {
		return "", errors.New("s3SecretKey is not set")
	}
	return r.s3SecretKey, nil
}

// GetS3BucketName возвращает имя бакета S3
func (r *EnvReader) GetS3BucketName() (string, error) {
	if r.s3Bucket == "" {
		return "", errors.New("s3Bucket is not set")
	}
	return r.s3Bucket, nil
}

// GetDBHost возвращает адрес базы данных
func (r *EnvReader) GetDBHost() (string, error) {
	if r.dbHost == "" {
		return "", errors.New("dbHost is not set")
	}
	return r.dbHost, nil
}

// GetDBPort возвращает порт базы данных
func (r *EnvReader) GetDBPort() (string, error) {
	if r.dbPort == "" {
		return "", errors.New("dbPort is not set")
	}
	return r.dbPort, nil
}

// GetDBUser возвращает имя пользователя для подключения к базе данных
func (r *EnvReader) GetDBUser() (string, error) {
	if r.dbUser == "" {
		return "", errors.New("dbUser is not set")
	}
	return r.dbUser, nil
}

// GetDBPassword возвращает пароль для подключения к базе данных
func (r *EnvReader) GetDBPassword() (string, error) {
	if r.dbPassword == "" {
		return "", errors.New("dbPassword is not set")
	}
	return r.dbPassword, nil
}

// GetDBName возвращает имя базы данных
func (r *EnvReader) GetDBName() (string, error) {
	if r.dbName == "" {
		return "", errors.New("dbName is not set")
	}
	return r.dbName, nil
}

func (r *EnvReader) GetDBSSLMode() (string, error) {
	if r.dbSSlMode == "" {
		return "", errors.New("dbSSLMode is not set")
	}
	return r.dbSSlMode, nil
}

// GetJWTSecret возвращает секрет JWT
func (r *EnvReader) GetJWTSecret() (string, error) {
	if r.jwtSecret == "" {
		return "", errors.New("jwtSecret is not set")
	}
	return r.jwtSecret, nil
}

func (r *EnvReader) GetJWTTokenLifeTime() (string, error) {
	if r.jwtTokenLifeTime == "" {
		return "", errors.New("jwtTokenLifeTime is not set")
	}
	return r.jwtTokenLifeTime, nil
}

func (r *EnvReader) GetJWTRefreshTokenLifeTime() (string, error) {
	if r.jwtRefreshTokenLifeTime == "" {
		return "", errors.New("jwtRefreshTokenLifeTime is not set")
	}
	return r.jwtRefreshTokenLifeTime, nil
}

// GetEncryptKey возвращает ключ шифрования
func (r *EnvReader) GetEncryptKey() (string, error) {
	if r.encryptKey == "" {
		return "", errors.New("encryptKey is not set")
	}
	return r.encryptKey, nil
}
