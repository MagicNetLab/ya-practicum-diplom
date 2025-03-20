package conf

// Configurator конфигуратор сервера приложения
type Configurator struct {
	serverHost       string
	serverPort       string
	fileStorageType  string
	fileStoragePath  string
	s3Endpoint       string
	s3SecretKey      string
	s3AccessKey      string
	s3Bucket         string
	dataStorageType  string
	inMemoryDumpPath string
	dbHost           string
	dbPort           string
	dbUser           string
	dbPassword       string
	dbName           string
	jwtSecret        string
	encryptKey       string
}

// ServerHost возвращает адрес сервера
func (c Configurator) ServerHost() string {
	return c.serverHost
}

// ServerPort возвращает порт сервера
func (c Configurator) ServerPort() string {
	return c.serverPort
}

// FileStorageType возвращает тип используемого хранилища
func (c Configurator) FileStorageType() string {
	return c.fileStorageType
}

// FileStoragePath возвращает путь к файлу хранилища
func (c Configurator) FileStoragePath() string {
	return c.fileStoragePath
}

// S3Endpoint возвращает адрес сервера s3
func (c Configurator) S3Endpoint() string {
	return c.s3Endpoint
}

// S3SecretKey возвращает secret key s3
func (c Configurator) S3SecretKey() string {
	return c.s3SecretKey
}

// S3AccessKey возвращает access key s3
func (c Configurator) S3AccessKey() string {
	return c.s3AccessKey
}

// S3Bucket возвращает bucket s3
func (c Configurator) S3Bucket() string {
	return c.s3Bucket
}

// DataStorageType возвращает тип используемого хранилища данных
func (c Configurator) DataStorageType() string {
	return c.dataStorageType
}

// InMemoryDumpPath возвращает путь к файлу в памяти
func (c Configurator) InMemoryDumpPath() string {
	return c.inMemoryDumpPath
}

// DBHost возвращает адрес базы данных
func (c Configurator) DBHost() string {
	return c.dbHost
}

// DBPort возвращает порт базы данных
func (c Configurator) DBPort() string {
	return c.dbPort
}

// DBUser возвращает имя пользователя для подключения к базе данных
func (c Configurator) DBUser() string {
	return c.dbUser
}

// DBPassword возвращает пароль для подключения к базе данных
func (c Configurator) DBPassword() string {
	return c.dbPassword
}

// DBName возвращает имя базы данных
func (c Configurator) DBName() string {
	return c.dbName
}

// JWTSecret возвращает secret key jwt
func (c Configurator) JWTSecret() string {
	return c.jwtSecret
}

// IsValid проверяет корректность конфигурации
func (c Configurator) IsValid() bool {
	isValidHost := c.serverHost != "" && c.serverPort != ""

	isValidS3FileStorage := c.fileStorageType == FileStorageS3 && c.s3Endpoint != "" && c.s3SecretKey != "" && c.s3AccessKey != "" && c.s3Bucket != ""

	isValidDataStoragePostgres := c.dataStorageType == DataStoragePostgres && c.dbHost != "" && c.dbPort != "" && c.dbUser != "" && c.dbPassword != "" && c.dbName != ""

	isValidJWTSecret := c.jwtSecret != ""

	isValidEncrytKey := c.encryptKey != ""

	return isValidHost && isValidS3FileStorage && isValidDataStoragePostgres && isValidJWTSecret && isValidEncrytKey
}
