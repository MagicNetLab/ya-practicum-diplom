package config

import "testing"

func TestS3Config(t *testing.T) {
	t.Run("GetEndpoint", func(t *testing.T) {
		tests := []struct {
			name     string
			config   S3Config
			expected string
		}{
			{
				name:     "empty endpoint",
				config:   S3Config{Endpoint: "", SecretKey: "secret", AccessKey: "access", Bucket: "bucket"},
				expected: "",
			},
			{
				name:     "valid endpoint",
				config:   S3Config{Endpoint: "http://localhost:9000", SecretKey: "secret", AccessKey: "access", Bucket: "bucket"},
				expected: "http://localhost:9000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetEndpoint(); got != tt.expected {
					t.Errorf("S3Config.GetEndpoint() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("GetSecretKey", func(t *testing.T) {
		tests := []struct {
			name     string
			config   S3Config
			expected string
		}{
			{
				name:     "empty secret key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "", AccessKey: "access", Bucket: "bucket"},
				expected: "",
			},
			{
				name:     "valid secret key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret123", AccessKey: "access", Bucket: "bucket"},
				expected: "secret123",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetSecretKey(); got != tt.expected {
					t.Errorf("S3Config.GetSecretKey() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("GetAccessKey", func(t *testing.T) {
		tests := []struct {
			name     string
			config   S3Config
			expected string
		}{
			{
				name:     "empty access key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "", Bucket: "bucket"},
				expected: "",
			},
			{
				name:     "valid access key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "access123", Bucket: "bucket"},
				expected: "access123",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetAccessKey(); got != tt.expected {
					t.Errorf("S3Config.GetAccessKey() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("GetBucket", func(t *testing.T) {
		tests := []struct {
			name     string
			config   S3Config
			expected string
		}{
			{
				name:     "empty bucket",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "access", Bucket: ""},
				expected: "",
			},
			{
				name:     "valid bucket",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "access", Bucket: "my-bucket"},
				expected: "my-bucket",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetBucket(); got != tt.expected {
					t.Errorf("S3Config.GetBucket() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			name     string
			config   S3Config
			expected bool
		}{
			{
				name:     "empty config",
				config:   S3Config{},
				expected: false,
			},
			{
				name:     "empty endpoint",
				config:   S3Config{Endpoint: "", SecretKey: "secret", AccessKey: "access", Bucket: "bucket"},
				expected: false,
			},
			{
				name:     "empty secret key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "", AccessKey: "access", Bucket: "bucket"},
				expected: false,
			},
			{
				name:     "empty access key",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "", Bucket: "bucket"},
				expected: false,
			},
			{
				name:     "empty bucket",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "access", Bucket: ""},
				expected: false,
			},
			{
				name:     "valid config",
				config:   S3Config{Endpoint: "endpoint", SecretKey: "secret", AccessKey: "access", Bucket: "bucket"},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.IsValid(); got != tt.expected {
					t.Errorf("S3Config.IsValid() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
