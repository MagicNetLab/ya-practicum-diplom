package config

// ServerConfigurator - конфигурация сервера
type ServerConfigurator interface {
	GetHost() string
	GetPort() string
	IsValid() bool
}

// ServerConfig - конфигурация сервера
type ServerConfig struct {
	host string `env:"SERVER_HOST" default:""`
	port string `env:"SERVER_PORT" default:""`
}

// GetHost возвращает хост сервера
func (s ServerConfig) GetHost() string {
	return s.host
}

// GetPort возвращает порт сервера
func (s ServerConfig) GetPort() string {
	return s.port
}

// IsValid проверяет корректность конфигурации сервера
func (s ServerConfig) IsValid() bool {
	return s.host != "" && s.port != ""
}
