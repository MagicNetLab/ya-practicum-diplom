package readers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvReader_Parse(t *testing.T) {
	// Подготовка тестовых данных
	tests := []struct {
		name     string
		envVars  map[string]string
		expected EnvReader
	}{
		{
			name: "all environment variables are set",
			envVars: map[string]string{
				"SERVER_HOST":                 "localhost",
				"SERVER_PORT":                 "8080",
				"S3_ENDPOINT":                 "http://localhost:9000",
				"S3_SECRET_KEY":               "secret",
				"S3_ACCESS_KEY":               "access",
				"S3_BUCKET":                   "bucket",
				"DB_HOST":                     "localhost",
				"DB_PORT":                     "5432",
				"DB_USER":                     "user",
				"DB_PASSWORD":                 "password",
				"DB_NAME":                     "dbname",
				"DB_SSL_MODE":                 "disable",
				"JWT_SECRET":                  "secret",
				"JWT_TOKEN_LIFE_TIME":         "3600",
				"JWT_REFRESH_TOKEN_LIFE_TIME": "7200",
				"ENCRYPT_KEY":                 "key",
			},
			expected: EnvReader{
				serverHost:              "localhost",
				serverPort:              "8080",
				s3Endpoint:              "http://localhost:9000",
				s3SecretKey:             "secret",
				s3AccessKey:             "access",
				s3Bucket:                "bucket",
				dbHost:                  "localhost",
				dbPort:                  "5432",
				dbUser:                  "user",
				dbPassword:              "password",
				dbName:                  "dbname",
				dbSSlMode:               "disable",
				jwtSecret:               "secret",
				jwtTokenLifeTime:        "3600",
				jwtRefreshTokenLifeTime: "7200",
				encryptKey:              "key",
			},
		},
		{
			name:     "no environment variables are set",
			envVars:  map[string]string{},
			expected: EnvReader{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Установка переменных окружения
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Очистка переменных окружения
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			reader := &EnvReader{}
			err := reader.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *reader)
		})
	}
}

func TestEnvReader_Getters(t *testing.T) {
	tests := []struct {
		name           string
		reader         EnvReader
		expectedValues map[string]string
		expectedErrors map[string]bool
	}{
		{
			name: "all values are set",
			reader: EnvReader{
				serverHost:              "localhost",
				serverPort:              "8080",
				s3Endpoint:              "http://localhost:9000",
				s3SecretKey:             "secret",
				s3AccessKey:             "access",
				s3Bucket:                "bucket",
				dbHost:                  "localhost",
				dbPort:                  "5432",
				dbUser:                  "user",
				dbPassword:              "password",
				dbName:                  "dbname",
				dbSSlMode:               "disable",
				jwtSecret:               "secret",
				jwtTokenLifeTime:        "3600",
				jwtRefreshTokenLifeTime: "7200",
				encryptKey:              "key",
			},
			expectedValues: map[string]string{
				"serverHost":              "localhost",
				"serverPort":              "8080",
				"s3Endpoint":              "http://localhost:9000",
				"s3SecretKey":             "secret",
				"s3AccessKey":             "access",
				"s3Bucket":                "bucket",
				"dbHost":                  "localhost",
				"dbPort":                  "5432",
				"dbUser":                  "user",
				"dbPassword":              "password",
				"dbName":                  "dbname",
				"dbSSLMode":               "disable",
				"jwtSecret":               "secret",
				"jwtTokenLifeTime":        "3600",
				"jwtRefreshTokenLifeTime": "7200",
				"encryptKey":              "key",
			},
			expectedErrors: map[string]bool{},
		},
		{
			name:   "no values are set",
			reader: EnvReader{},
			expectedValues: map[string]string{
				"serverHost":              "",
				"serverPort":              "",
				"s3Endpoint":              "",
				"s3SecretKey":             "",
				"s3AccessKey":             "",
				"s3Bucket":                "",
				"dbHost":                  "",
				"dbPort":                  "",
				"dbUser":                  "",
				"dbPassword":              "",
				"dbName":                  "",
				"dbSSLMode":               "",
				"jwtSecret":               "",
				"jwtTokenLifeTime":        "",
				"jwtRefreshTokenLifeTime": "",
				"encryptKey":              "",
			},
			expectedErrors: map[string]bool{
				"serverHost":              true,
				"serverPort":              true,
				"s3Endpoint":              true,
				"s3SecretKey":             true,
				"s3AccessKey":             true,
				"s3Bucket":                true,
				"dbHost":                  true,
				"dbPort":                  true,
				"dbUser":                  true,
				"dbPassword":              true,
				"dbName":                  true,
				"dbSSLMode":               true,
				"jwtSecret":               true,
				"jwtTokenLifeTime":        true,
				"jwtRefreshTokenLifeTime": true,
				"encryptKey":              true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Тестирование GetServerHost
			val, err := tt.reader.GetServerHost()
			assertGetterResult(t, "serverHost", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetServerPort
			val, err = tt.reader.GetServerPort()
			assertGetterResult(t, "serverPort", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetS3EndPoint
			val, err = tt.reader.GetS3EndPoint()
			assertGetterResult(t, "s3Endpoint", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetS3SecretKey
			val, err = tt.reader.GetS3SecretKey()
			assertGetterResult(t, "s3SecretKey", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetS3AccessKey
			val, err = tt.reader.GetS3AccessKey()
			assertGetterResult(t, "s3AccessKey", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetS3BucketName
			val, err = tt.reader.GetS3BucketName()
			assertGetterResult(t, "s3Bucket", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBHost
			val, err = tt.reader.GetDBHost()
			assertGetterResult(t, "dbHost", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBPort
			val, err = tt.reader.GetDBPort()
			assertGetterResult(t, "dbPort", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBUser
			val, err = tt.reader.GetDBUser()
			assertGetterResult(t, "dbUser", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBPassword
			val, err = tt.reader.GetDBPassword()
			assertGetterResult(t, "dbPassword", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBName
			val, err = tt.reader.GetDBName()
			assertGetterResult(t, "dbName", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetDBSSLMode
			val, err = tt.reader.GetDBSSLMode()
			assertGetterResult(t, "dbSSLMode", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetJWTSecret
			val, err = tt.reader.GetJWTSecret()
			assertGetterResult(t, "jwtSecret", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetJWTTokenLifeTime
			val, err = tt.reader.GetJWTTokenLifeTime()
			assertGetterResult(t, "jwtTokenLifeTime", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetJWTRefreshTokenLifeTime
			val, err = tt.reader.GetJWTRefreshTokenLifeTime()
			assertGetterResult(t, "jwtRefreshTokenLifeTime", val, err, tt.expectedValues, tt.expectedErrors)

			// Тестирование GetEncryptKey
			val, err = tt.reader.GetEncryptKey()
			assertGetterResult(t, "encryptKey", val, err, tt.expectedValues, tt.expectedErrors)
		})
	}
}

// Вспомогательная функция для проверки результатов геттеров
func assertGetterResult(t *testing.T, field string, value string, err error, expectedValues map[string]string, expectedErrors map[string]bool) {
	if expectedErrors[field] {
		assert.Error(t, err)
		assert.Equal(t, "", value)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedValues[field], value)
	}
}
