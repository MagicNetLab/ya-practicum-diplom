package readers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig_Parse(t *testing.T) {
	config := &DefaultConfig{}
	err := config.Parse()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.serverHost)
	assert.Equal(t, "8080", config.serverPort)
	assert.Equal(t, "localhost", config.dbHost)
	assert.Equal(t, "5432", config.dbPort)
	assert.Equal(t, "gophkeeper", config.dbUser)
	assert.Equal(t, "gophkeeper", config.dbPassword)
	assert.Equal(t, "gophkeeper", config.dbName)
	assert.Equal(t, "disable", config.dbSSLMode)
	assert.Equal(t, "gophkeeper", config.jwtSecret)
	assert.Equal(t, "60", config.jwtTokenLifeTime)
	assert.Equal(t, "3600", config.jwtRefreshTokenLifeTime)
	assert.Equal(t, "gophkeeper", config.encryptKey)
}

func TestDefaultConfig_GetServerHost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{serverHost: "localhost"}
		host, err := config.GetServerHost()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", host)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		host, err := config.GetServerHost()
		assert.Error(t, err)
		assert.Equal(t, "", host)
		assert.EqualError(t, err, "serverHost is not set")
	})
}

func TestDefaultConfig_GetServerPort(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{serverPort: "8080"}
		port, err := config.GetServerPort()
		assert.NoError(t, err)
		assert.Equal(t, "8080", port)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		port, err := config.GetServerPort()
		assert.Error(t, err)
		assert.Equal(t, "", port)
		assert.EqualError(t, err, "serverPort is not set")
	})
}

func TestDefaultConfig_GetS3EndPoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{s3Endpoint: "endpoint"}
		endpoint, err := config.GetS3EndPoint()
		assert.NoError(t, err)
		assert.Equal(t, "endpoint", endpoint)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		endpoint, err := config.GetS3EndPoint()
		assert.Error(t, err)
		assert.Equal(t, "", endpoint)
		assert.EqualError(t, err, "s3Endpoint is not set")
	})
}

func TestDefaultConfig_GetS3AccessKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{s3AccessKey: "key"}
		key, err := config.GetS3AccessKey()
		assert.NoError(t, err)
		assert.Equal(t, "key", key)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		key, err := config.GetS3AccessKey()
		assert.Error(t, err)
		assert.Equal(t, "", key)
		assert.EqualError(t, err, "s3AccessKey is not set")
	})
}

func TestDefaultConfig_GetS3SecretKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{s3SecretKey: "key"}
		key, err := config.GetS3SecretKey()
		assert.NoError(t, err)
		assert.Equal(t, "key", key)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		key, err := config.GetS3SecretKey()
		assert.Error(t, err)
		assert.Equal(t, "", key)
		assert.EqualError(t, err, "s3SecretKey is not set")
	})
}

func TestDefaultConfig_GetS3BucketName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{s3Bucket: "bucket"}
		bucket, err := config.GetS3BucketName()
		assert.NoError(t, err)
		assert.Equal(t, "bucket", bucket)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		bucket, err := config.GetS3BucketName()
		assert.Error(t, err)
		assert.Equal(t, "", bucket)
		assert.EqualError(t, err, "s3Bucket is not set")
	})
}

func TestDefaultConfig_GetDBHost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbHost: "localhost"}
		host, err := config.GetDBHost()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", host)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		host, err := config.GetDBHost()
		assert.Error(t, err)
		assert.Equal(t, "", host)
		assert.EqualError(t, err, "dbHost is not set")
	})
}

func TestDefaultConfig_GetDBPort(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbPort: "5432"}
		port, err := config.GetDBPort()
		assert.NoError(t, err)
		assert.Equal(t, "5432", port)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		port, err := config.GetDBPort()
		assert.Error(t, err)
		assert.Equal(t, "", port)
		assert.EqualError(t, err, "dbPort is not set")
	})
}

func TestDefaultConfig_GetDBUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbUser: "user"}
		user, err := config.GetDBUser()
		assert.NoError(t, err)
		assert.Equal(t, "user", user)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		user, err := config.GetDBUser()
		assert.Error(t, err)
		assert.Equal(t, "", user)
		assert.EqualError(t, err, "dbUser is not set")
	})
}

func TestDefaultConfig_GetDBPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbPassword: "password"}
		password, err := config.GetDBPassword()
		assert.NoError(t, err)
		assert.Equal(t, "password", password)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		password, err := config.GetDBPassword()
		assert.Error(t, err)
		assert.Equal(t, "", password)
		assert.EqualError(t, err, "dbPassword is not set")
	})
}

func TestDefaultConfig_GetDBName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbName: "dbname"}
		name, err := config.GetDBName()
		assert.NoError(t, err)
		assert.Equal(t, "dbname", name)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		name, err := config.GetDBName()
		assert.Error(t, err)
		assert.Equal(t, "", name)
		assert.EqualError(t, err, "dbName is not set")
	})
}

func TestDefaultConfig_GetDBSSLMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{dbSSLMode: "disable"}
		mode, err := config.GetDBSSLMode()
		assert.NoError(t, err)
		assert.Equal(t, "disable", mode)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		mode, err := config.GetDBSSLMode()
		assert.Error(t, err)
		assert.Equal(t, "", mode)
		assert.EqualError(t, err, "dbSSLMode is not set")
	})
}

func TestDefaultConfig_GetJWTSecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{jwtSecret: "secret"}
		secret, err := config.GetJWTSecret()
		assert.NoError(t, err)
		assert.Equal(t, "secret", secret)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		secret, err := config.GetJWTSecret()
		assert.Error(t, err)
		assert.Equal(t, "", secret)
		assert.EqualError(t, err, "jwtSecret is not set")
	})
}

func TestDefaultConfig_GetJWTTokenLifeTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{jwtTokenLifeTime: "60"}
		lifetime, err := config.GetJWTTokenLifeTime()
		assert.NoError(t, err)
		assert.Equal(t, "60", lifetime)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		lifetime, err := config.GetJWTTokenLifeTime()
		assert.Error(t, err)
		assert.Equal(t, "", lifetime)
		assert.EqualError(t, err, "jwtTokenLifeTime is not set")
	})
}

func TestDefaultConfig_GetJWTRefreshTokenLifeTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{jwtRefreshTokenLifeTime: "3600"}
		lifetime, err := config.GetJWTRefreshTokenLifeTime()
		assert.NoError(t, err)
		assert.Equal(t, "3600", lifetime)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		lifetime, err := config.GetJWTRefreshTokenLifeTime()
		assert.Error(t, err)
		assert.Equal(t, "", lifetime)
		assert.EqualError(t, err, "jwtRefreshTokenLifeTime is not set")
	})
}

func TestDefaultConfig_GetEncryptKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &DefaultConfig{encryptKey: "key"}
		key, err := config.GetEncryptKey()
		assert.NoError(t, err)
		assert.Equal(t, "key", key)
	})

	t.Run("empty", func(t *testing.T) {
		config := &DefaultConfig{}
		key, err := config.GetEncryptKey()
		assert.Error(t, err)
		assert.Equal(t, "", key)
		assert.EqualError(t, err, "encryptKey is not set")
	})
}
