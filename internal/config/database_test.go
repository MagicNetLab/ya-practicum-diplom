package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDatabaseConfig_GetDSN tests the GetDSN method of the DatabaseConfig
func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "Valid DSN",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "5432",
				user:     "postgres",
				password: "secret",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: "postgres://postgres:secret@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "Empty Fields",
			config: DatabaseConfig{
				host:     "",
				port:     "",
				user:     "",
				password: "",
				dbname:   "",
				sslMode:  "disable",
			},
			expected: "postgres://:@:/?sslmode=disable",
		},
		{
			name: "Custom SSL Mode",
			config: DatabaseConfig{
				host:     "db.example.com",
				port:     "5432",
				user:     "admin",
				password: "pass123",
				dbname:   "proddb",
				sslMode:  "require",
			},
			expected: "postgres://admin:pass123@db.example.com:5432/proddb?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetDSN()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDatabaseConfig_IsValid tests the IsValid method of the DatabaseConfig
func TestDatabaseConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected bool
	}{
		{
			name: "Valid Config",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "5432",
				user:     "postgres",
				password: "secret",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: true,
		},
		{
			name: "Empty Host",
			config: DatabaseConfig{
				host:     "",
				port:     "5432",
				user:     "postgres",
				password: "secret",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: false,
		},
		{
			name: "Empty Port",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "",
				user:     "postgres",
				password: "secret",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: false,
		},
		{
			name: "Empty User",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "5432",
				user:     "",
				password: "secret",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: false,
		},
		{
			name: "Empty Password",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "5432",
				user:     "postgres",
				password: "",
				dbname:   "testdb",
				sslMode:  "disable",
			},
			expected: false,
		},
		{
			name: "Empty DBName",
			config: DatabaseConfig{
				host:     "localhost",
				port:     "5432",
				user:     "postgres",
				password: "secret",
				dbname:   "",
				sslMode:  "disable",
			},
			expected: false,
		},
		{
			name: "All Empty Fields",
			config: DatabaseConfig{
				host:     "",
				port:     "",
				user:     "",
				password: "",
				dbname:   "",
				sslMode:  "disable",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}
