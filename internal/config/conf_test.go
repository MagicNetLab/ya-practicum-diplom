package config

import (
	"errors"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config/readers"
	"github.com/stretchr/testify/assert"
)

func TestGetCnf(t *testing.T) {
	t.Run("Успешное получение конфигурации", func(t *testing.T) {
		Config = Configurator{}
		// Подменяем функцию getReaders для возврата нашего мока
		original := getReadersFunc
		// Override with test implementation
		getReadersFunc = func() ([]Reader, error) {
			reader := new(mocks.Reader)
			reader.On("Parse").Return(nil)
			reader.On("GetServerHost").Return("localhost", nil)
			reader.On("GetServerPort").Return("8080", nil)
			reader.On("GetS3EndPoint").Return("endpoint", nil)
			reader.On("GetS3SecretKey").Return("secret", nil)
			reader.On("GetS3AccessKey").Return("access", nil)
			reader.On("GetS3BucketName").Return("bucket", nil)
			reader.On("GetDBHost").Return("dbhost", nil)
			reader.On("GetDBPort").Return("5432", nil)
			reader.On("GetDBUser").Return("user", nil)
			reader.On("GetDBPassword").Return("password", nil)
			reader.On("GetDBName").Return("dbname", nil)
			reader.On("GetDBSSLMode").Return("disable", nil)
			reader.On("GetJWTSecret").Return("jwtsecret", nil)
			reader.On("GetJWTTokenLifeTime").Return("3600", nil)
			reader.On("GetJWTRefreshTokenLifeTime").Return("7200", nil)
			reader.On("GetEncryptKey").Return("encryptkey", nil)
			return []Reader{reader}, nil
		}
		// Restore original after test
		defer func() { getReadersFunc = original }()

		// Вызываем тестируемую функцию
		config, err := MakeConfig()
		assert.NoError(t, err)
		assert.True(t, config.IsValid())
	})

	t.Run("Ошибка при получении конфигурации", func(t *testing.T) {
		Config = Configurator{}
		// Подменяем функцию getReaders для возврата нашего мока
		original := getReadersFunc
		// Override with test implementation
		getReadersFunc = func() ([]Reader, error) {
			reader := new(mocks.Reader)
			reader.On("Parse").Return(errors.New("parse error"))
			return []Reader{reader}, nil
		}
		// Restore original after test
		defer func() { getReadersFunc = original }()

		config, err := MakeConfig()
		assert.Error(t, err)
		assert.False(t, config.IsValid())

	})

	t.Run("Ошибка при получении обязательных параметров", func(t *testing.T) {
		Config = Configurator{}
		original := getReadersFunc
		getReadersFunc = func() ([]Reader, error) {
			reader := new(mocks.Reader)
			reader.On("Parse").Return(nil)
			reader.On("GetServerHost").Return("", errors.New("no host"))
			reader.On("GetServerPort").Return("", errors.New("no port"))
			reader.On("GetS3EndPoint").Return("endpoint", nil)
			reader.On("GetS3SecretKey").Return("secret", nil)
			reader.On("GetS3AccessKey").Return("access", nil)
			reader.On("GetS3BucketName").Return("bucket", nil)
			reader.On("GetDBHost").Return("dbhost", nil)
			reader.On("GetDBPort").Return("5432", nil)
			reader.On("GetDBUser").Return("user", nil)
			reader.On("GetDBPassword").Return("password", nil)
			reader.On("GetDBName").Return("dbname", nil)
			reader.On("GetDBSSLMode").Return("disable", nil)
			reader.On("GetJWTSecret").Return("jwtsecret", nil)
			reader.On("GetJWTTokenLifeTime").Return("3600", nil)
			reader.On("GetJWTRefreshTokenLifeTime").Return("7200", nil)
			reader.On("GetEncryptKey").Return("encryptkey", nil)
			return []Reader{reader}, nil
		}
		defer func() { getReadersFunc = original }()
		config, err := MakeConfig()
		assert.Error(t, err)
		assert.False(t, config.IsValid())
	})
}

func TestGetReaders(t *testing.T) {
	r, err := getReaders()

	// Проверяем, что нет ошибки
	assert.NoError(t, err)

	// Проверяем, что возвращены все три ридера
	assert.Len(t, r, 3)

	// Проверяем типы возвращаемых ридеров
	assert.IsType(t, &readers.DefaultConfig{}, r[0])
	assert.IsType(t, &readers.EnvReader{}, r[1])
	assert.IsType(t, &readers.FlagReader{}, r[2])
}
