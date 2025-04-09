package readers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestConfig() *ConfigReader {
	return &ConfigReader{
		Server: ServerConfig{
			Host: "localhost",
			Port: "3200",
		},
		S3: S3Config{
			Endpoint:  "localhost:9000",
			SecretKey: "minioadmin",
			AccessKey: "minioadmin",
			Bucket:    "gophkeeper",
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "gophkeeper",
			Password: "gophkeeper",
			Name:     "gophkeeper",
			SSLMode:  "disable",
		},
		JWT: JWTConfig{
			Secret:               "test-secret",
			TokenLifeTime:        "1",
			RefreshTokenLifeTime: "24",
		},
		EncryptKey: "test-encrypt-key",
	}
}

func TestConfigReader_Parse(t *testing.T) {
	t.Run("Успешное чтение конфигурации", func(t *testing.T) {
		// Создаем временный файл конфигурации
		content := []byte(`
server:
  host: localhost
  port: 3200
s3:
  endpoint: localhost:9000
  secretKey: minioadmin
  accessKey: minioadmin
  bucket: gophkeeper
`)
		err := os.WriteFile("config.yaml", content, 0644)
		assert.NoError(t, err)
		defer os.Remove("config.yaml")

		config := &ConfigReader{}
		err = config.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", config.Server.Host)
	})

	t.Run("Ошибка: файл конфигурации не найден", func(t *testing.T) {
		config := &ConfigReader{}
		err := config.Parse()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config.yaml not found")
	})
}

func TestConfigReader_GetServerHost(t *testing.T) {
	t.Run("Успешное получение хоста сервера", func(t *testing.T) {
		config := setupTestConfig()
		host, err := config.GetServerHost()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", host)
	})

	t.Run("Ошибка: пустой хост сервера", func(t *testing.T) {
		config := setupTestConfig()
		config.Server.Host = ""
		host, err := config.GetServerHost()
		assert.Error(t, err)
		assert.Empty(t, host)
	})
}

func TestConfigReader_GetServerPort(t *testing.T) {
	t.Run("Успешное получение порта сервера", func(t *testing.T) {
		config := setupTestConfig()
		port, err := config.GetServerPort()
		assert.NoError(t, err)
		assert.Equal(t, "3200", port)
	})

	t.Run("Ошибка: пустой порт сервера", func(t *testing.T) {
		config := setupTestConfig()
		config.Server.Port = ""
		port, err := config.GetServerPort()
		assert.Error(t, err)
		assert.Empty(t, port)
	})
}

func TestConfigReader_GetS3Config(t *testing.T) {
	t.Run("Успешное получение конфигурации S3", func(t *testing.T) {
		config := setupTestConfig()

		endpoint, err := config.GetS3EndPoint()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:9000", endpoint)

		accessKey, err := config.GetS3AccessKey()
		assert.NoError(t, err)
		assert.Equal(t, "minioadmin", accessKey)

		secretKey, err := config.GetS3SecretKey()
		assert.NoError(t, err)
		assert.Equal(t, "minioadmin", secretKey)

		bucket, err := config.GetS3BucketName()
		assert.NoError(t, err)
		assert.Equal(t, "gophkeeper", bucket)
	})

	t.Run("Ошибки: пустые значения S3", func(t *testing.T) {
		config := setupTestConfig()
		config.S3 = S3Config{}

		_, err := config.GetS3EndPoint()
		assert.Error(t, err)

		_, err = config.GetS3AccessKey()
		assert.Error(t, err)

		_, err = config.GetS3SecretKey()
		assert.Error(t, err)

		_, err = config.GetS3BucketName()
		assert.Error(t, err)
	})
}

func TestConfigReader_GetDBConfig(t *testing.T) {
	t.Run("Успешное получение конфигурации БД", func(t *testing.T) {
		config := setupTestConfig()

		host, err := config.GetDBHost()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", host)

		port, err := config.GetDBPort()
		assert.NoError(t, err)
		assert.Equal(t, "5432", port)

		user, err := config.GetDBUser()
		assert.NoError(t, err)
		assert.Equal(t, "gophkeeper", user)

		password, err := config.GetDBPassword()
		assert.NoError(t, err)
		assert.Equal(t, "gophkeeper", password)

		name, err := config.GetDBName()
		assert.NoError(t, err)
		assert.Equal(t, "gophkeeper", name)

		sslMode, err := config.GetDBSSLMode()
		assert.NoError(t, err)
		assert.Equal(t, "disable", sslMode)
	})

	t.Run("Ошибки: пустые значения БД", func(t *testing.T) {
		config := setupTestConfig()
		config.DB = DBConfig{}

		_, err := config.GetDBHost()
		assert.Error(t, err)

		_, err = config.GetDBPort()
		assert.Error(t, err)

		_, err = config.GetDBUser()
		assert.Error(t, err)

		_, err = config.GetDBPassword()
		assert.Error(t, err)

		_, err = config.GetDBName()
		assert.Error(t, err)

		_, err = config.GetDBSSLMode()
		assert.Error(t, err)
	})
}

func TestConfigReader_GetJWTConfig(t *testing.T) {
	t.Run("Успешное получение конфигурации JWT", func(t *testing.T) {
		config := setupTestConfig()

		secret, err := config.GetJWTSecret()
		assert.NoError(t, err)
		assert.Equal(t, "test-secret", secret)

		tokenLifeTime, err := config.GetJWTTokenLifeTime()
		assert.NoError(t, err)
		assert.Equal(t, "1", tokenLifeTime)

		refreshTokenLifeTime, err := config.GetJWTRefreshTokenLifeTime()
		assert.NoError(t, err)
		assert.Equal(t, "24", refreshTokenLifeTime)
	})

	t.Run("Ошибки: пустые значения JWT", func(t *testing.T) {
		config := setupTestConfig()
		config.JWT = JWTConfig{}

		_, err := config.GetJWTSecret()
		assert.Error(t, err)

		_, err = config.GetJWTTokenLifeTime()
		assert.Error(t, err)

		_, err = config.GetJWTRefreshTokenLifeTime()
		assert.Error(t, err)
	})
}

func TestConfigReader_GetEncryptKey(t *testing.T) {
	t.Run("Успешное получение ключа шифрования", func(t *testing.T) {
		config := setupTestConfig()
		key, err := config.GetEncryptKey()
		assert.NoError(t, err)
		assert.Equal(t, "test-encrypt-key", key)
	})

	t.Run("Ошибка: пустой ключ шифрования", func(t *testing.T) {
		config := setupTestConfig()
		config.EncryptKey = ""
		key, err := config.GetEncryptKey()
		assert.Error(t, err)
		assert.Empty(t, key)
	})
}
