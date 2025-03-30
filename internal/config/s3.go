package config

// S3Configurator интерфейс конфигурации S3 хранилища
type S3Configurator interface {
	GetEndpoint() string
	GetSecretKey() string
	GetAccessKey() string
	GetBucket() string
	IsValid() bool
}

// S3Config конфигурация S3 хранилища
type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT" envDefault:""`
	SecretKey string `env:"S3_SECRET_KEY" envDefault:""`
	AccessKey string `env:"S3_ACCESS_KEY" envDefault:""`
	Bucket    string `env:"S3_BUCKET" envDefault:""`
}

// GetEndpoint возвращает Endpoint S3 хранилища
func (s S3Config) GetEndpoint() string {
	return s.Endpoint
}

// GetSecretKey возвращает secret key S3 хранилища
func (s S3Config) GetSecretKey() string {
	return s.SecretKey
}

// GetAccessKey возвращает access key S3 хранилища
func (s S3Config) GetAccessKey() string {
	return s.AccessKey
}

// GetBucket возвращает Bucket S3 хранилища
func (s S3Config) GetBucket() string {
	return s.Bucket
}

// IsValid проверяет корректность конфигурации S3 хранилища
func (s S3Config) IsValid() bool {
	return s.Endpoint != "" && s.SecretKey != "" && s.AccessKey != "" && s.Bucket != ""
}
