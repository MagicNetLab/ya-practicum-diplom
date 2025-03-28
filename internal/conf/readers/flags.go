package readers

import (
	"errors"
	"flag"
)

const (
	serverHostKey      = "serverHost"
	serverPortKey      = "serverPOrt"
	s3EndpointKey      = "s3Endpoint"
	s3SecretKeyKey     = "s3Secret"
	s3AccessKeyKey     = "s3Access"
	s3BucketKey        = "s3Bucket"
	dbHostKey          = "dbHost"
	dbPortKey          = "dbPort"
	dbUserKey          = "dbUser"
	dbPasswordKey      = "dbPassword"
	dbNameKey          = "dbName"
	dbSSLModeKey       = "dbSSLMode"
	jwtSecretKey       = "jwtSecret"
	jwtTokenKey        = "tokenTime"
	jwtRefreshTokenKey = "refreshTokenTime"
	encryptSecretKey   = "encryptKey"
)

type FlagReader struct {
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
	dbSSLMode               string
	jwtSecret               string
	jwtTokenLifeTime        string
	jwtRefreshTokenLifeTime string
	encryptKey              string
}

// Parse читает установленные параметры конфигурации
func (r *FlagReader) Parse() error {
	var serverHost, serverPort string
	var s3Endpoint, s3SecretKey, s3AccessKey, s3Bucket string
	var dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode string
	var jwtSecret, jwtTokenLifeTime, jwtRefreshTokenLifeTime string
	var encryptKey string

	flag.StringVar(&serverHost, serverHostKey, "", "server host")
	flag.StringVar(&serverPort, serverPortKey, "", "server port")

	flag.StringVar(&s3Endpoint, s3EndpointKey, "", "s3 endpoint")
	flag.StringVar(&s3SecretKey, s3SecretKeyKey, "", "s3 secret key")
	flag.StringVar(&s3AccessKey, s3AccessKeyKey, "", "s3 access key")
	flag.StringVar(&s3Bucket, s3BucketKey, "", "s3 bucket")

	flag.StringVar(&dbHost, dbHostKey, "", "database host")
	flag.StringVar(&dbPort, dbPortKey, "", "database port")
	flag.StringVar(&dbUser, dbUserKey, "", "database user")
	flag.StringVar(&dbPassword, dbPasswordKey, "", "database password")
	flag.StringVar(&dbName, dbNameKey, "", "database name")
	flag.StringVar(&dbSSLMode, dbSSLModeKey, "", "database ssl mode")

	flag.StringVar(&jwtSecret, jwtSecretKey, "", "jwt secret")
	flag.StringVar(&jwtTokenLifeTime, jwtTokenKey, "", "jwt token life time")
	flag.StringVar(&jwtRefreshTokenLifeTime, jwtRefreshTokenKey, "", "jwt refresh token life time")

	flag.StringVar(&encryptKey, encryptSecretKey, "", "encrypt key")

	flag.Parse()

	if serverHost != "" {
		r.serverPort = serverPort
	}

	if serverPort != "" {
		r.serverPort = serverPort
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

	if dbSSLMode != "" {
		r.dbSSLMode = dbSSLMode
	}

	if jwtSecret != "" {
		r.jwtSecret = jwtSecret
	}

	if jwtTokenLifeTime != "" {
		r.jwtTokenLifeTime = jwtTokenLifeTime
	}

	if jwtRefreshTokenLifeTime != "" {
		r.jwtRefreshTokenLifeTime = jwtRefreshTokenLifeTime
	}

	if encryptKey != "" {
		r.encryptKey = encryptKey
	}

	return nil
}

func (r *FlagReader) GetServerHost() (string, error) {
	if r.serverHost == "" {
		return "", errors.New("serverHost is not set")
	}
	return r.serverHost, nil
}

// GetServerPort возвращает порт сервера
func (r *FlagReader) GetServerPort() (string, error) {
	if r.serverPort == "" {
		return "", errors.New("serverPort is not set")
	}
	return r.serverPort, nil
}

// GetS3EndPoint возвращает адрес S3
func (r *FlagReader) GetS3EndPoint() (string, error) {
	if r.s3Endpoint == "" {
		return "", errors.New("s3Endpoint is not set")
	}
	return r.s3Endpoint, nil
}

// GetS3AccessKey возвращает ключ S3
func (r *FlagReader) GetS3AccessKey() (string, error) {
	if r.s3AccessKey == "" {
		return "", errors.New("s3AccessKey is not set")
	}
	return r.s3AccessKey, nil
}

// GetS3SecretKey возвращает ключ S3
func (r *FlagReader) GetS3SecretKey() (string, error) {
	if r.s3SecretKey == "" {
		return "", errors.New("s3SecretKey is not set")
	}
	return r.s3SecretKey, nil
}

// GetS3BucketName возвращает имя бакета S3
func (r *FlagReader) GetS3BucketName() (string, error) {
	if r.s3Bucket == "" {
		return "", errors.New("s3Bucket is not set")
	}
	return r.s3Bucket, nil
}

// GetDBHost возвращает адрес базы данных
func (r *FlagReader) GetDBHost() (string, error) {
	if r.dbHost == "" {
		return "", errors.New("dbHost is not set")
	}
	return r.dbHost, nil
}

// GetDBPort возвращает порт базы данных
func (r *FlagReader) GetDBPort() (string, error) {
	if r.dbPort == "" {
		return "", errors.New("dbPort is not set")
	}
	return r.dbPort, nil
}

// GetDBUser возвращает имя пользователя для подключения к базе данных
func (r *FlagReader) GetDBUser() (string, error) {
	if r.dbUser == "" {
		return "", errors.New("dbUser is not set")
	}
	return r.dbUser, nil
}

// GetDBPassword возвращает пароль для подключения к базе данных
func (r *FlagReader) GetDBPassword() (string, error) {
	if r.dbPassword == "" {
		return "", errors.New("dbPassword is not set")
	}
	return r.dbPassword, nil
}

// GetDBName возвращает имя базы данных
func (r *FlagReader) GetDBName() (string, error) {
	if r.dbName == "" {
		return "", errors.New("dbName is not set")
	}
	return r.dbName, nil
}

// GetDBSSLMode возвращает режим SSL для подключения к базе данных
func (r *FlagReader) GetDBSSLMode() (string, error) {
	if r.dbSSLMode == "" {
		return "", errors.New("dbSSLMode is not set")
	}
	return r.dbSSLMode, nil
}

// GetJWTSecret возвращает секрет JWT
func (r *FlagReader) GetJWTSecret() (string, error) {
	if r.jwtSecret == "" {
		return "", errors.New("jwtSecret is not set")
	}
	return r.jwtSecret, nil
}

// GetJWTTokenLifeTime возвращает время жизни JWT
func (r *FlagReader) GetJWTTokenLifeTime() (string, error) {
	if r.jwtTokenLifeTime == "" {
		return "", errors.New("jwtTokenLifeTime is not set")
	}
	return r.jwtTokenLifeTime, nil
}

// GetJWTRefreshTokenLifeTime возвращает время жизни JWT
func (r *FlagReader) GetJWTRefreshTokenLifeTime() (string, error) {
	if r.jwtRefreshTokenLifeTime == "" {
		return "", errors.New("jwtRefreshTokenLifeTime is not set")
	}
	return r.jwtRefreshTokenLifeTime, nil
}

// GetEncryptKey возвращает ключ шифрования
func (r *FlagReader) GetEncryptKey() (string, error) {
	if r.encryptKey == "" {
		return "", errors.New("encryptKey is not set")
	}
	return r.encryptKey, nil
}
