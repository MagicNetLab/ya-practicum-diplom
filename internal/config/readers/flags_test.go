package readers

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagReader_Parse(t *testing.T) {
	// Сохраняем оригинальные аргументы
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name string
		args []string
		want FlagReader
	}{
		{
			name: "empty flags",
			args: []string{"prog"},
			want: FlagReader{},
		},
		{
			name: "all flags set",
			args: []string{
				"prog",
				"-serverHost", "localhost",
				"-serverPOrt", "8080",
				"-s3Endpoint", "http://s3.local",
				"-s3Secret", "secret",
				"-s3Access", "access",
				"-s3Bucket", "bucket",
				"-dbHost", "dbhost",
				"-dbPort", "5432",
				"-dbUser", "user",
				"-dbPassword", "pass",
				"-dbName", "dbname",
				"-dbSSLMode", "disable",
				"-jwtSecret", "secret",
				"-tokenTime", "3600",
				"-refreshTokenTime", "7200",
				"-encryptKey", "key",
			},
			want: FlagReader{
				serverPort:              "8080",
				s3Endpoint:              "http://s3.local",
				s3SecretKey:             "secret",
				s3AccessKey:             "access",
				s3Bucket:                "bucket",
				dbHost:                  "dbhost",
				dbPort:                  "5432",
				dbUser:                  "user",
				dbPassword:              "pass",
				dbName:                  "dbname",
				dbSSLMode:               "disable",
				jwtSecret:               "secret",
				jwtTokenLifeTime:        "3600",
				jwtRefreshTokenLifeTime: "7200",
				encryptKey:              "key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем флаги перед каждым тестом
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			r := &FlagReader{}
			err := r.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, *r)
		})
	}
}

func TestFlagReader_Getters(t *testing.T) {
	tests := []struct {
		name     string
		reader   FlagReader
		getters  []func(*FlagReader) (string, error)
		wantErr  bool
		wantVals []string
	}{
		{
			name:   "empty reader",
			reader: FlagReader{},
			getters: []func(*FlagReader) (string, error){
				(*FlagReader).GetServerHost,
				(*FlagReader).GetServerPort,
				(*FlagReader).GetS3EndPoint,
				(*FlagReader).GetS3AccessKey,
				(*FlagReader).GetS3SecretKey,
				(*FlagReader).GetS3BucketName,
				(*FlagReader).GetDBHost,
				(*FlagReader).GetDBPort,
				(*FlagReader).GetDBUser,
				(*FlagReader).GetDBPassword,
				(*FlagReader).GetDBName,
				(*FlagReader).GetDBSSLMode,
				(*FlagReader).GetJWTSecret,
				(*FlagReader).GetJWTTokenLifeTime,
				(*FlagReader).GetJWTRefreshTokenLifeTime,
				(*FlagReader).GetEncryptKey,
			},
			wantErr:  true,
			wantVals: make([]string, 16),
		},
		{
			name: "filled reader",
			reader: FlagReader{
				serverHost:              "localhost",
				serverPort:              "8080",
				s3Endpoint:              "http://s3.local",
				s3SecretKey:             "secret",
				s3AccessKey:             "access",
				s3Bucket:                "bucket",
				dbHost:                  "dbhost",
				dbPort:                  "5432",
				dbUser:                  "user",
				dbPassword:              "pass",
				dbName:                  "dbname",
				dbSSLMode:               "disable",
				jwtSecret:               "secret",
				jwtTokenLifeTime:        "3600",
				jwtRefreshTokenLifeTime: "7200",
				encryptKey:              "key",
			},
			getters: []func(*FlagReader) (string, error){
				(*FlagReader).GetServerHost,
				(*FlagReader).GetServerPort,
				(*FlagReader).GetS3EndPoint,
				(*FlagReader).GetS3AccessKey,
				(*FlagReader).GetS3SecretKey,
				(*FlagReader).GetS3BucketName,
				(*FlagReader).GetDBHost,
				(*FlagReader).GetDBPort,
				(*FlagReader).GetDBUser,
				(*FlagReader).GetDBPassword,
				(*FlagReader).GetDBName,
				(*FlagReader).GetDBSSLMode,
				(*FlagReader).GetJWTSecret,
				(*FlagReader).GetJWTTokenLifeTime,
				(*FlagReader).GetJWTRefreshTokenLifeTime,
				(*FlagReader).GetEncryptKey,
			},
			wantErr: false,
			wantVals: []string{
				"localhost",
				"8080",
				"http://s3.local",
				"access",
				"secret",
				"bucket",
				"dbhost",
				"5432",
				"user",
				"pass",
				"dbname",
				"disable",
				"secret",
				"3600",
				"7200",
				"key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, getter := range tt.getters {
				val, err := getter(&tt.reader)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.wantVals[i], val)
				}
			}
		})
	}
}
