package conf

import (
	"errors"
	"testing"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf/mocks"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf/readers"
	"github.com/stretchr/testify/assert"
)

func TestGetCnf(t *testing.T) {
	// Сбрасываем глобальную конфигурацию перед каждым тестом
	Config = Configurator{}

	tests := []struct {
		name      string
		setupMock func() *mocks.Reader
		wantErr   bool
	}{
		{
			name: "успешное получение конфигурации",
			setupMock: func() *mocks.Reader {
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
				return reader
			},
			wantErr: false,
		},
		{
			name: "ошибка при парсинге конфигурации",
			setupMock: func() *mocks.Reader {
				reader := new(mocks.Reader)
				reader.On("Parse").Return(errors.New("parse error"))
				return reader
			},
			wantErr: true,
		},
		{
			name: "ошибка при получении обязательных параметров",
			setupMock: func() *mocks.Reader {
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
				return reader
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подменяем функцию getReaders для возврата нашего мока
			original := getReadersFunc
			// Override with test implementation
			getReadersFunc = func() ([]Reader, error) {
				return []Reader{tt.setupMock()}, nil
			}
			// Restore original after test
			defer func() { getReadersFunc = original }()

			// Вызываем тестируемую функцию
			config, err := GetCnf()

			// Проверяем результаты
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, config.IsValid())
			}
		})
	}
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
