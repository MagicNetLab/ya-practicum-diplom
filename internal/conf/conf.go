package conf

import (
	"errors"
	"strconv"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf/readers"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

var Config Configurator

var getReadersFunc = getReaders

// GetCnf возвращает конфигурацию приложения
func GetCnf() (Configurator, error) {
	if Config.IsValid() {
		return Config, nil
	}

	cnfReaders, err := getReadersFunc()
	if err != nil {
		logger.Error("Failed to get readers from config", logger.StrArg("error", err.Error()))
		return Config, err
	}

	for _, reader := range cnfReaders {
		err := reader.Parse()
		if err != nil {
			logger.Error("Failed to parse config", logger.StrArg("error", err.Error()))
			continue
		}

		servHost, err := reader.GetServerHost()
		if err == nil {
			Config.serverHost = servHost
		}

		servPort, err := reader.GetServerPort()
		if err == nil {
			Config.serverPort = servPort
		}

		s3Endpoint, err := reader.GetS3EndPoint()
		if err == nil {
			Config.s3Endpoint = s3Endpoint
		}

		s3SecretKey, err := reader.GetS3SecretKey()
		if err == nil {
			Config.s3SecretKey = s3SecretKey
		}

		s3AccessKey, err := reader.GetS3AccessKey()
		if err == nil {
			Config.s3AccessKey = s3AccessKey
		}

		s3Bucket, err := reader.GetS3BucketName()
		if err == nil {
			Config.s3Bucket = s3Bucket
		}

		dbHost, err := reader.GetDBHost()
		if err == nil {
			Config.dbHost = dbHost
		}

		dbPort, err := reader.GetDBPort()
		if err == nil {
			Config.dbPort = dbPort
		}

		dbUser, err := reader.GetDBUser()
		if err == nil {
			Config.dbUser = dbUser
		}

		dbPassword, err := reader.GetDBPassword()
		if err == nil {
			Config.dbPassword = dbPassword
		}

		dbName, err := reader.GetDBName()
		if err == nil {
			Config.dbName = dbName
		}

		dbSSLMode, err := reader.GetDBSSLMode()
		if err == nil {
			Config.dbSSLMode = dbSSLMode
		}

		jwtSecret, err := reader.GetJWTSecret()
		if err == nil {
			Config.jwtSecret = jwtSecret
		} else {
			Config.jwtSecret = jwt.GetRandomSecret()
		}

		jwtTokenLifeTime, err := reader.GetJWTTokenLifeTime()
		if err == nil {
			val, err := strconv.Atoi(jwtTokenLifeTime)
			if err == nil {
				Config.jwtTokenLifeTime = val
			}
		}

		jwtRefreshTokenLifeTime, err := reader.GetJWTRefreshTokenLifeTime()
		if err == nil {
			val, err := strconv.Atoi(jwtRefreshTokenLifeTime)
			if err == nil {
				Config.jwtTokenLifeTime = val
			}
		}

		encryptKey, err := reader.GetEncryptKey()
		if err == nil {
			Config.encryptKey = encryptKey
		}
	}

	if Config.IsValid() {
		return Config, nil
	}

	return Config, errors.New("invalid config")
}

func getReaders() ([]Reader, error) {
	return []Reader{
		&readers.DefaultConfig{},
		&readers.EnvReader{},
		&readers.FlagReader{},
	}, nil
}
