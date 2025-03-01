package config

// S3Configurator интерфейс конфигурации S3 хранилища
type S3Configurator interface {
	GetEndpoint() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	IsValid() bool
}

type S3Config struct {
	endpoint  string `env:"S3_ENDPOINT" envDefault:""`
	secretKey string `env:"S3_SECRET_KEY" envDefault:""`
	accessKey string `env:"S3_ACCESS_KEY" envDefault:""`
	bucket    string `env:"S3_BUCKET" envDefault:""`
}

// GetEndpoint возвращает endpoint S3 хранилища
func (s S3Config) GetEndpoint() string {
	return s.endpoint
}

// GetSecretKey возвращает secret key S3 хранилища
func (s S3Config) GetSecretKey() string {
	return s.secretKey
}

// GetAccessKey возвращает access key S3 хранилища
func (s S3Config) GetAccessKey() string {
	return s.accessKey
}

// GetBucket возвращает bucket S3 хранилища
func (s S3Config) GetBucket() string {
	return s.bucket
}

// IsValid проверяет корректность конфигурации S3 хранилища
func (s S3Config) IsValid() bool {
	return s.endpoint != "" && s.secretKey != "" && s.accessKey != "" && s.bucket != ""
}
