package readers

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type ConfigReader struct {
	Server     ServerConfig `yaml:"server"`
	S3         S3Config     `yaml:"s3"`
	DB         DBConfig     `yaml:"db"`
	JWT        JWTConfig    `yaml:"jwt"`
	EncryptKey string       `yaml:"encryptKey" env:"ENCRYPT_KEY" env-default:"askjHJ83JIUbjkgygJfh64ghVhbmbvt664ghbjkgygJfh64ghVhbmbvt6"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"3200"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint" env:"S3_ENDPOINT" env-default:"localhost:9000"`
	SecretKey string `yaml:"secretKey" env:"S3_SECRET_KEY" env-default:"minioadmin"`
	AccessKey string `yaml:"accessKey" env:"S3_ACCESS_KEY" env-default:"minioadmin"`
	Bucket    string `yaml:"bucket" env:"S3_BUCKET" env-default:"gophkeeper"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"gophkeeper"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"gophkeeper"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"gophkeeper"`
	SSLMode  string `yaml:"sslMode" env:"DB_SSL_MODE" env-default:"disable"`
}

type JWTConfig struct {
	Secret               string `yaml:"secret" env:"JWT_SECRET" env-default:"11sb6GFNBhas86BjjndHJBROcxd4adrmc"`
	TokenLifeTime        string `yaml:"tokenLifeTime" env:"JWT_TOKEN_LIFE_TIME" env-default:"1"`
	RefreshTokenLifeTime string `yaml:"refreshTokenLifeTime" env:"JWT_REFRESH_TOKEN_LIFE_TIME" env-default:"24"`
}

func (y *ConfigReader) Parse() error {
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config.yaml not found")
	}

	if err := cleanenv.ReadConfig(configPath, y); err != nil {
		return err
	}

	return nil
}

// GetServerHost возвращает хост сервера
func (y *ConfigReader) GetServerHost() (string, error) {
	if y.Server.Host == "" {
		return "", fmt.Errorf("Server Host is empty")
	}
	return y.Server.Host, nil
}

// GetServerPort возвращает порт сервера
func (y *ConfigReader) GetServerPort() (string, error) {
	if y.Server.Port == "" {
		return "", fmt.Errorf("Server Port is empty")
	}
	return y.Server.Port, nil
}

// GetS3EndPoint возвращает S3 Endpoint
func (y *ConfigReader) GetS3EndPoint() (string, error) {
	if y.S3.Endpoint == "" {
		return "", fmt.Errorf("S3 Endpoint is empty")
	}
	return y.S3.Endpoint, nil
}

// GetS3AccessKey возвращает ключ доступа к S3
func (y *ConfigReader) GetS3AccessKey() (string, error) {
	if y.S3.AccessKey == "" {
		return "", fmt.Errorf("S3 access key is empty")
	}
	return y.S3.AccessKey, nil
}

// GetS3SecretKey возвращает секретный ключ доступа к S3
func (y *ConfigReader) GetS3SecretKey() (string, error) {
	if y.S3.SecretKey == "" {
		return "", fmt.Errorf("S3 secret key is empty")
	}
	return y.S3.SecretKey, nil
}

// GetS3BucketName возвращает название бакета на S3
func (y *ConfigReader) GetS3BucketName() (string, error) {
	if y.S3.Bucket == "" {
		return "", fmt.Errorf("S3 Bucket Name is empty")
	}
	return y.S3.Bucket, nil
}

// GetDBHost возвращает хост базы данных
func (y *ConfigReader) GetDBHost() (string, error) {
	if y.DB.Host == "" {
		return "", fmt.Errorf("DB Host is empty")
	}
	return y.DB.Host, nil
}

// GetDBPort возвращает порт базы данных
func (y *ConfigReader) GetDBPort() (string, error) {
	if y.DB.Port == "" {
		return "", fmt.Errorf("DB Port is empty")
	}
	return y.DB.Port, nil
}

// GetDBUser возвращает имя пользователя базы данных
func (y *ConfigReader) GetDBUser() (string, error) {
	if y.DB.User == "" {
		return "", fmt.Errorf("DB User is empty")
	}
	return y.DB.User, nil
}

// GetDBPassword возвращает пароль пользователя базы данных
func (y *ConfigReader) GetDBPassword() (string, error) {
	if y.DB.Password == "" {
		return "", fmt.Errorf("DB Password is empty")
	}
	return y.DB.Password, nil
}

// GetDBName возвращает название базы данных
func (y *ConfigReader) GetDBName() (string, error) {
	if y.DB.Name == "" {
		return "", fmt.Errorf("DB Name is empty")
	}
	return y.DB.Name, nil
}

// GetDBSSLMode возвращает режим подключения к базе данных
func (y *ConfigReader) GetDBSSLMode() (string, error) {
	if y.DB.SSLMode == "" {
		return "", fmt.Errorf("DB ssl mode is empty")
	}
	return y.DB.SSLMode, nil
}

// GetJWTSecret возвращает секретный ключ для генерации токенов
func (y *ConfigReader) GetJWTSecret() (string, error) {
	if y.JWT.Secret == "" {
		return "", fmt.Errorf("JWT secret is empty")
	}
	return y.JWT.Secret, nil
}

// GetJWTTokenLifeTime возвращает время жизни токена
func (y *ConfigReader) GetJWTTokenLifeTime() (string, error) {
	if y.JWT.TokenLifeTime == "" {
		return "", fmt.Errorf("JWT token life time is empty")
	}
	return y.JWT.TokenLifeTime, nil
}

// GetJWTRefreshTokenLifeTime возвращает время жизни токена обновления
func (y *ConfigReader) GetJWTRefreshTokenLifeTime() (string, error) {
	if y.JWT.RefreshTokenLifeTime == "" {
		return "", fmt.Errorf("JWT refresh token life time is empty")
	}
	return y.JWT.RefreshTokenLifeTime, nil
}

func (y *ConfigReader) GetEncryptKey() (string, error) {
	if y.EncryptKey == "" {
		return "", fmt.Errorf("encrypt key is empty")
	}
	return y.EncryptKey, nil
}
