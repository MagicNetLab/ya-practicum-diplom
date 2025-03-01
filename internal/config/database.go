package config

import "fmt"

// DataBaseConfigurator интерфейс конфигурации подключения к БД
type DataBaseConfigurator interface {
	GetDSN() string
	IsValid() bool
}

// DatabaseConfig конфигурации подключения к БД
type DatabaseConfig struct {
	host     string `env:"DB_HOST" envDefault:""`
	port     string `env:"DB_PORT" envDefault:""`
	user     string `env:"DB_USER" envDefault:""`
	password string `env:"DB_PASSWORD" envDefault:""`
	dbname   string `env:"DB_NAME" envDefault:""`
	sslMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

// GetDSN возвращает DSN для подключения к БД
func (d DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", d.user, d.password, d.host, d.port, d.dbname, d.sslMode)
}

// IsValid проверяет корректность конфигурации подключения к БД
func (d DatabaseConfig) IsValid() bool {
	return d.host != "" && d.port != "" && d.user != "" && d.password != "" && d.dbname != ""
}
