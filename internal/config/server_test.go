package config

import "testing"

func TestServerConfig(t *testing.T) {
	t.Run("GetHost", func(t *testing.T) {
		tests := []struct {
			name     string
			config   ServerConfig
			expected string
		}{
			{
				name:     "empty host",
				config:   ServerConfig{host: "", port: "8080"},
				expected: "",
			},
			{
				name:     "localhost",
				config:   ServerConfig{host: "localhost", port: "8080"},
				expected: "localhost",
			},
			{
				name:     "ip address",
				config:   ServerConfig{host: "127.0.0.1", port: "8080"},
				expected: "127.0.0.1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetHost(); got != tt.expected {
					t.Errorf("ServerConfig.GetHost() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("GetPort", func(t *testing.T) {
		tests := []struct {
			name     string
			config   ServerConfig
			expected string
		}{
			{
				name:     "empty port",
				config:   ServerConfig{host: "localhost", port: ""},
				expected: "",
			},
			{
				name:     "valid port",
				config:   ServerConfig{host: "localhost", port: "8080"},
				expected: "8080",
			},
			{
				name:     "system port",
				config:   ServerConfig{host: "localhost", port: "80"},
				expected: "80",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetPort(); got != tt.expected {
					t.Errorf("ServerConfig.GetPort() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			name     string
			config   ServerConfig
			expected bool
		}{
			{
				name:     "empty config",
				config:   ServerConfig{host: "", port: ""},
				expected: false,
			},
			{
				name:     "empty host",
				config:   ServerConfig{host: "", port: "8080"},
				expected: false,
			},
			{
				name:     "empty port",
				config:   ServerConfig{host: "localhost", port: ""},
				expected: false,
			},
			{
				name:     "valid config",
				config:   ServerConfig{host: "localhost", port: "8080"},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.IsValid(); got != tt.expected {
					t.Errorf("ServerConfig.IsValid() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
