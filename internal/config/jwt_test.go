package config

import (
	"testing"
	"time"
)

// TestJWTConfig tests the JWTConfig
func TestJWTConfig(t *testing.T) {
	// Test GetJWTSecret method
	t.Run("GetJWTSecret", func(t *testing.T) {
		tests := []struct {
			name     string
			config   JWTConfig
			expected string
		}{
			{
				name:     "empty secret",
				config:   JWTConfig{secret: "", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: "",
			},
			{
				name:     "valid secret",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: "test-secret",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetJWTSecret(); got != tt.expected {
					t.Errorf("JWTConfig.GetJWTSecret() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	// Test IsValid method
	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			name     string
			config   JWTConfig
			expected bool
		}{
			{
				name:     "empty secret",
				config:   JWTConfig{secret: "", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: false,
			},
			{
				name:     "valid secret",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.IsValid(); got != tt.expected {
					t.Errorf("JWTConfig.IsValid() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	// Test GetTokenLifeTime method
	t.Run("GetTokenLifeTime", func(t *testing.T) {
		tests := []struct {
			name     string
			config   JWTConfig
			expected time.Duration
		}{
			{
				name:     "default token lifetime",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: time.Hour,
			},
			{
				name:     "custom token lifetime",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 2, refreshTokenLifeTime: 24},
				expected: 2 * time.Hour,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetTokenLifeTime(); got != tt.expected {
					t.Errorf("JWTConfig.GetTokenLifeTime() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	// Test GetRefreshTokenLifeTime method
	t.Run("GetRefreshTokenLifeTime", func(t *testing.T) {
		tests := []struct {
			name     string
			config   JWTConfig
			expected time.Duration
		}{
			{
				name:     "default refresh token lifetime",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 1, refreshTokenLifeTime: 24},
				expected: 24 * time.Hour,
			},
			{
				name:     "custom refresh token lifetime",
				config:   JWTConfig{secret: "test-secret", tokenLifeTime: 1, refreshTokenLifeTime: 48},
				expected: 48 * time.Hour,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.GetRefreshTokenLifeTime(); got != tt.expected {
					t.Errorf("JWTConfig.GetRefreshTokenLifeTime() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
