package readers

import (
	"errors"
	"os"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

type EnvReader struct {
	serverHost       string
	serverPort       string
	fileStorageType  string
	fileStoragePath  string
	inMemoryDumpPath string
	s3Endpoint       string
	s3SecretKey      string
	s3AccessKey      string
	s3Bucket         string
	dataStorageType  string
	dbHost           string
	dbPort           string
	dbUser           string
	dbPassword       string
	dbName           string
	jwtSecret        string
	encryptKey       string
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

	fileStorageType := os.Getenv("FILE_STORAGE_TYPE")
	if fileStorageType != "" {
		r.fileStorageType = fileStorageType
	}

	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePath != "" {
		r.fileStoragePath = fileStoragePath
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

	dataStorageType := os.Getenv("DATA_STORAGE_TYPE")
	if dataStorageType != "" {
		r.dataStorageType = dataStorageType
	}

	inMemoryDumpPath := os.Getenv("IN_MEMORY_DUMP_PATH")
	if inMemoryDumpPath != "" {
		r.inMemoryDumpPath = inMemoryDumpPath
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		r.jwtSecret = jwtSecret
	}

	encryptKey := os.Getenv("ENCRIPT_KEY")
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

// GetFileStorageType возвращает тип хранилища (local, s3)
func (r *EnvReader) GetFileStorageType() (string, error) {
	if r.fileStorageType == "" {
		return "", errors.New("fileStorageType is not set")
	}
	return r.fileStorageType, nil
}

// GetFileStoragePath возвращает путь к файлу хранилища
func (r *EnvReader) GetFileStoragePath() (string, error) {
	if r.fileStoragePath == "" {
		return "", errors.New("fileStoragePath is not set")
	}
	return r.fileStoragePath, nil
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

// GetDataStorageType возвращает тип хранилища данных (inmemory, postgres)
func (r *EnvReader) GetDataStorageType() (string, error) {
	if r.dataStorageType == "" {
		return "", errors.New("dataStorageType is not set")
	}
	return r.dataStorageType, nil
}

// GetInMemoryDumpPath возвращает путь к папке с дамп
func (r *EnvReader) GetInMemoryDumpPath() (string, error) {
	if r.inMemoryDumpPath == "" {
		return "", errors.New("inMemoryDumpPath is not set")
	}
	return r.inMemoryDumpPath, nil
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

// GetJWTSecret возвращает секрет JWT
func (r *EnvReader) GetJWTSecret() (string, error) {
	if r.jwtSecret == "" {
		return "", errors.New("jwtSecret is not set")
	}
	return r.jwtSecret, nil
}

// GetEncryptKey возвращает ключ шифрования
func (r *EnvReader) GetEncryptKey() (string, error) {
	if r.encryptKey == "" {
		return "", errors.New("encryptKey is not set")
	}
	return r.encryptKey, nil
}
