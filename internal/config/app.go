package config

// AppConfigurator интерфейс общей конфигурации приложения
type AppConfigurator interface {
	GetServerConf() ServerConfigurator
	GetDBConf() DataBaseConfigurator
	GetS3Conf() S3Configurator
	GetJWTConf() JWTConfigurator
	IsValid() bool
}

// AppConfig общая конфигурация приложения
type AppConfig struct {
	serverConf ServerConfigurator
	dbConf     DataBaseConfigurator
	s3Conf     S3Configurator
	jwtConf    JWTConfigurator
}

// GetServerConf возвращает конфигурацию сервера
func (a AppConfig) GetServerConf() ServerConfigurator {
	return a.serverConf
}

// GetDBConf возвращает конфигурацию подключения к БД
func (a AppConfig) GetDBConf() DataBaseConfigurator {
	return a.dbConf
}

// GetS3Conf возвращает конфигурацию S3
func (a AppConfig) GetS3Conf() S3Configurator {
	return a.s3Conf
}

// GetJWTConf возвращает конфигурацию JWT
func (a AppConfig) GetJWTConf() JWTConfigurator {
	return a.jwtConf
}

// IsValid возвращает утверждение корректности конфигурации приложения
func (a AppConfig) IsValid() bool {
	return a.serverConf.IsValid() && a.dbConf.IsValid() && a.s3Conf.IsValid() && a.jwtConf.IsValid()
}
