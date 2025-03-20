package readers

import (
	"errors"
)

type DefaultConfig struct {
	serverHost       string
	serverPort       string
	fileStorageType  string
	fileStoragePath  string
	s3Endpoint       string
	s3SecretKey      string
	s3AccessKey      string
	s3Bucket         string
	dataStorageType  string
	inMemoryDumpPath string
	dbHost           string
	dbPort           string
	dbUser           string
	dbPassword       string
	dbName           string
	jwtSecret        string
	encryptKey       string
}

// Parse читает установленные параметры конфигурации
func (r *DefaultConfig) Parse() error {
	r.serverHost = "localhost"
	r.serverPort = "8080"
	r.fileStorageType = "local"
	r.fileStoragePath = "/storage"
	r.dataStorageType = "inmemory"
	r.inMemoryDumpPath = "/values/inmemory-dump"
	r.dbHost = "localhost"
	r.dbPort = "5432"
	r.dbUser = "gophkeeper"
	r.dbPassword = "gophkeeper"
	r.dbName = "gophkeeper"

	return nil
}

// GetServerHost возвращает адрес сервера
func (r *DefaultConfig) GetServerHost() (string, error) {
	if r.serverHost == "" {
		return "", errors.New("serverHost is not set")
	}
	return r.serverHost, nil
}

// GetServerPort возвращает порт сервера
func (r *DefaultConfig) GetServerPort() (string, error) {
	if r.serverPort == "" {
		return "", errors.New("serverPort is not set")
	}
	return r.serverPort, nil
}

// GetFileStorageType возвращает тип хранилища (local, s3)
func (r *DefaultConfig) GetFileStorageType() (string, error) {
	if r.fileStorageType == "" {
		return "", errors.New("fileStorageType is not set")
	}
	return r.fileStorageType, nil
}

// GetFileStoragePath возвращает путь к файлу хранилища
func (r *DefaultConfig) GetFileStoragePath() (string, error) {
	if r.fileStoragePath == "" {
		return "", errors.New("fileStoragePath is not set")
	}
	return r.fileStoragePath, nil
}

// GetS3EndPoint возвращает адрес S3
func (r *DefaultConfig) GetS3EndPoint() (string, error) {
	if r.s3Endpoint == "" {
		return "", errors.New("s3Endpoint is not set")
	}
	return r.s3Endpoint, nil
}

// GetS3AccessKey возвращает ключ S3
func (r *DefaultConfig) GetS3AccessKey() (string, error) {
	if r.s3AccessKey == "" {
		return "", errors.New("s3AccessKey is not set")
	}
	return r.s3AccessKey, nil
}

// GetS3SecretKey возвращает ключ S3
func (r *DefaultConfig) GetS3SecretKey() (string, error) {
	if r.s3SecretKey == "" {
		return "", errors.New("s3SecretKey is not set")
	}
	return r.s3SecretKey, nil
}

// GetS3BucketName возвращает имя бакета S3
func (r *DefaultConfig) GetS3BucketName() (string, error) {
	if r.s3Bucket == "" {
		return "", errors.New("s3Bucket is not set")
	}
	return r.s3Bucket, nil
}

// GetDataStorageType возвращает тип хранилища данных (inmemory, postgres)
func (r *DefaultConfig) GetDataStorageType() (string, error) {
	if r.dataStorageType == "" {
		return "", errors.New("dataStorageType is not set")
	}
	return r.dataStorageType, nil
}

// GetInMemoryDumpPath возвращает путь к папке с дампом данных
func (r *DefaultConfig) GetInMemoryDumpPath() (string, error) {
	if r.inMemoryDumpPath == "" {
		return "", errors.New("inMemoryDumpPath is not set")
	}
	return r.inMemoryDumpPath, nil
}

// GetDBHost возвращает адрес базы данных
func (r *DefaultConfig) GetDBHost() (string, error) {
	if r.dbHost == "" {
		return "", errors.New("dbHost is not set")
	}
	return r.dbHost, nil
}

// GetDBPort возвращает порт базы данных
func (r *DefaultConfig) GetDBPort() (string, error) {
	if r.dbPort == "" {
		return "", errors.New("dbPort is not set")
	}
	return r.dbPort, nil
}

// GetDBUser возвращает имя пользователя для подключения к базе данных
func (r *DefaultConfig) GetDBUser() (string, error) {
	if r.dbUser == "" {
		return "", errors.New("dbUser is not set")
	}
	return r.dbUser, nil
}

// GetDBPassword возвращает пароль для подключения к базе данных
func (r *DefaultConfig) GetDBPassword() (string, error) {
	if r.dbPassword == "" {
		return "", errors.New("dbPassword is not set")
	}
	return r.dbPassword, nil
}

// GetDBName возвращает имя базы данных
func (r *DefaultConfig) GetDBName() (string, error) {
	if r.dbName == "" {
		return "", errors.New("dbName is not set")
	}
	return r.dbName, nil
}

// GetJWTSecret возвращает секрет JWT
func (r *DefaultConfig) GetJWTSecret() (string, error) {
	if r.jwtSecret == "" {
		return "", errors.New("jwtSecret is not set")
	}
	return r.jwtSecret, nil
}

// GetEncryptKey возвращает ключ шифрования данных
func (r *DefaultConfig) GetEncryptKey() (string, error) {
	if r.encryptKey == "" {
		return "", errors.New("encryptKey is not set")
	}
	return r.encryptKey, nil
}
