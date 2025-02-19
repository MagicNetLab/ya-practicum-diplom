package readers

import (
	"errors"
	"flag"
)

const (
	serverHostKey       = "sh"
	serverPortKey       = "sp"
	fileStorageTypeKey  = "ft"
	fileStoragePathKey  = "fp"
	s3EndpointKey       = "es"
	s3SecretKeyKey      = "sk"
	s3AccessKeyKey      = "ak"
	s3BucketKey         = "bk"
	dataStorageTypeKey  = "dt"
	inMemoryDumpPathKey = "dp"
	dbHostKey           = "bh"
	dbPortKey           = "bp"
	dbUserKey           = "uu"
	dbPasswordKey       = "dp"
	dbNameKey           = "dn"
	jwtSecretKey        = "jwt"
)

type FlagReader struct {
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
}

// Parse читает установленные параметры конфигурации
func (r *FlagReader) Parse() error {
	var serverHost, serverPort, fileStorageType, fileStoragePath, s3Endpoint, s3SecretKey, s3AccessKey, s3Bucket, dataStorageType, inMemoryDumpPath, dbHost, dbPort, dbUser, dbPassword, dbName, jwtSecret string

	flag.StringVar(&serverHost, serverHostKey, "", "server host")
	flag.StringVar(&serverPort, serverPortKey, "", "server port")
	flag.StringVar(&fileStorageType, fileStorageTypeKey, "", "file storage type")
	flag.StringVar(&fileStoragePath, fileStoragePathKey, "", "file storage path")
	flag.StringVar(&s3Endpoint, s3EndpointKey, "", "s3 endpoint")
	flag.StringVar(&s3SecretKey, s3SecretKeyKey, "", "s3 secret key")
	flag.StringVar(&s3AccessKey, s3AccessKeyKey, "", "s3 access key")
	flag.StringVar(&s3Bucket, s3BucketKey, "", "s3 bucket")
	flag.StringVar(&dataStorageType, dataStorageTypeKey, "", "data storage type")
	flag.StringVar(&dbHost, dbHostKey, "", "database host")
	flag.StringVar(&dbPort, dbPortKey, "", "database port")
	flag.StringVar(&dbUser, dbUserKey, "", "database user")
	flag.StringVar(&dbPassword, dbPasswordKey, "", "database password")
	flag.StringVar(&dbName, dbNameKey, "", "database name")
	flag.StringVar(&jwtSecret, jwtSecretKey, "", "jwt secret")
	flag.StringVar(&inMemoryDumpPath, inMemoryDumpPathKey, "", "in memory dump path")
	flag.Parse()

	if serverHost != "" {
		r.serverPort = serverPort
	}

	if serverPort != "" {
		r.serverPort = serverPort
	}

	if fileStorageType != "" {
		r.fileStorageType = fileStorageType
	}

	if fileStoragePath != "" {
		r.fileStoragePath = fileStoragePath
	}

	if s3Endpoint != "" {
		r.s3Endpoint = s3Endpoint
	}

	if s3SecretKey != "" {
		r.s3SecretKey = s3SecretKey
	}

	if s3AccessKey != "" {
		r.s3AccessKey = s3AccessKey
	}

	if s3Bucket != "" {
		r.s3Bucket = s3Bucket
	}

	if dataStorageType != "" {
		r.dataStorageType = dataStorageType
	}

	if inMemoryDumpPath != "" {
		r.inMemoryDumpPath = inMemoryDumpPath
	}

	if dbHost != "" {
		r.dbHost = dbHost
	}

	if dbPort != "" {
		r.dbPort = dbPort
	}

	if dbUser != "" {
		r.dbUser = dbUser
	}

	if dbPassword != "" {
		r.dbPassword = dbPassword
	}

	if dbName != "" {
		r.dbName = dbName
	}

	return nil
}

func (r FlagReader) GetServerHost() (string, error) {
	if r.serverHost == "" {
		return "", errors.New("serverHost is not set")
	}
	return r.serverHost, nil
}

// GetServerPort возвращает порт сервера
func (r FlagReader) GetServerPort() (string, error) {
	if r.serverPort == "" {
		return "", errors.New("serverPort is not set")
	}
	return r.serverPort, nil
}

// GetFileStorageType возвращает тип хранилища (local, s3)
func (r FlagReader) GetFileStorageType() (string, error) {
	if r.fileStorageType == "" {
		return "", errors.New("fileStorageType is not set")
	}
	return r.fileStorageType, nil
}

// GetFileStoragePath возвращает путь к файлу хранилища
func (r FlagReader) GetFileStoragePath() (string, error) {
	if r.fileStoragePath == "" {
		return "", errors.New("fileStoragePath is not set")
	}
	return r.fileStoragePath, nil
}

// GetS3EndPoint возвращает адрес S3
func (r FlagReader) GetS3EndPoint() (string, error) {
	if r.s3Endpoint == "" {
		return "", errors.New("s3Endpoint is not set")
	}
	return r.s3Endpoint, nil
}

// GetS3AccessKey возвращает ключ S3
func (r FlagReader) GetS3AccessKey() (string, error) {
	if r.s3AccessKey == "" {
		return "", errors.New("s3AccessKey is not set")
	}
	return r.s3AccessKey, nil
}

// GetS3SecretKey возвращает ключ S3
func (r FlagReader) GetS3SecretKey() (string, error) {
	if r.s3SecretKey == "" {
		return "", errors.New("s3SecretKey is not set")
	}
	return r.s3SecretKey, nil
}

// GetS3BucketName возвращает имя бакета S3
func (r FlagReader) GetS3BucketName() (string, error) {
	if r.s3Bucket == "" {
		return "", errors.New("s3Bucket is not set")
	}
	return r.s3Bucket, nil
}

// GetDataStorageType возвращает тип хранилища данных (inmemory, postgres)
func (r FlagReader) GetDataStorageType() (string, error) {
	if r.dataStorageType == "" {
		return "", errors.New("dataStorageType is not set")
	}
	return r.dataStorageType, nil
}

// GetInMemoryDumpPath возвращает путь к папке с дамп
func (r FlagReader) GetInMemoryDumpPath() (string, error) {
	if r.inMemoryDumpPath == "" {
		return "", errors.New("inMemoryDumpPath is not set")
	}
	return r.inMemoryDumpPath, nil
}

// GetDBHost возвращает адрес базы данных
func (r FlagReader) GetDBHost() (string, error) {
	if r.dbHost == "" {
		return "", errors.New("dbHost is not set")
	}
	return r.dbHost, nil
}

// GetDBPort возвращает порт базы данных
func (r FlagReader) GetDBPort() (string, error) {
	if r.dbPort == "" {
		return "", errors.New("dbPort is not set")
	}
	return r.dbPort, nil
}

// GetDBUser возвращает имя пользователя для подключения к базе данных
func (r FlagReader) GetDBUser() (string, error) {
	if r.dbUser == "" {
		return "", errors.New("dbUser is not set")
	}
	return r.dbUser, nil
}

// GetDBPassword возвращает пароль для подключения к базе данных
func (r FlagReader) GetDBPassword() (string, error) {
	if r.dbPassword == "" {
		return "", errors.New("dbPassword is not set")
	}
	return r.dbPassword, nil
}

// GetDBName возвращает имя базы данных
func (r FlagReader) GetDBName() (string, error) {
	if r.dbName == "" {
		return "", errors.New("dbName is not set")
	}
	return r.dbName, nil
}

// GetJWTSecret возвращает секрет JWT
func (r FlagReader) GetJWTSecret() (string, error) {
	if r.jwtSecret == "" {
		return "", errors.New("jwtSecret is not set")
	}
	return r.jwtSecret, nil
}
