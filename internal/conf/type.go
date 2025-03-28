package conf

const (
	FileStorageS3       = "s3"
	DataStoragePostgres = "postgres"
)

// Reader интерфейс для плагинов чтения конфигураций
type Reader interface {
	// Parse читает установленные параметры конфигурации
	Parse() error
	// GetServerHost возвращает адрес сервера
	GetServerHost() (string, error)
	// GetServerPort возвращает порт сервера
	GetServerPort() (string, error)
	// GetS3EndPoint возвращает адрес S3
	GetS3EndPoint() (string, error)
	// GetS3AccessKey возвращает ключ S3
	GetS3AccessKey() (string, error)
	// GetS3SecretKey возвращает ключ S3
	GetS3SecretKey() (string, error)
	// GetS3BucketName возвращает имя бакета S3
	GetS3BucketName() (string, error)
	// GetDBHost возвращает адрес базы данных
	GetDBHost() (string, error)
	// GetDBPort возвращает порт базы данных
	GetDBPort() (string, error)
	// GetDBUser возвращает имя пользователя для подключения к базе данных
	GetDBUser() (string, error)
	// GetDBPassword возвращает пароль для подключения к базе данных
	GetDBPassword() (string, error)
	// GetDBName возвращает имя базы данных
	GetDBName() (string, error)
	// GetDBSSLMode возвращает режим SSL базы данных
	GetDBSSLMode() (string, error)
	// GetJWTSecret возвращает ключ JWT
	GetJWTSecret() (string, error)
	// GetJWTTokenLifeTime возвращает время жизни токена JWT
	GetJWTTokenLifeTime() (string, error)
	// GetJWTRefreshTokenLifeTime возвращает время жизни токена обновления JWT
	GetJWTRefreshTokenLifeTime() (string, error)
	// GetEncryptKey возвращает ключ шифрования
	GetEncryptKey() (string, error)
}
