package config

import "testing"

// TestSecretConfig tests the SecretConfig
func TestSecretConfig(t *testing.T) {
	// Test GetSecret method
	t.Run("GetSecret", func(t *testing.T) {
		tests := []struct {
			name     string
			config   SecretConfig
			expected string
			wantErr  bool
		}{
			{
				name:     "empty secret",
				config:   SecretConfig{secretKey: ""},
				expected: "",
				wantErr:  false,
			},
			{
				name:     "valid secret",
				config:   SecretConfig{secretKey: "test-secret-key"},
				expected: "test-secret-key",
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := tt.config.GetSecret()
				if (err != nil) != tt.wantErr {
					t.Errorf("SecretConfig.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.expected {
					t.Errorf("SecretConfig.GetSecret() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	// Test IsValid method
	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			name     string
			config   SecretConfig
			expected bool
		}{
			{
				name:     "empty secret",
				config:   SecretConfig{secretKey: ""},
				expected: false,
			},
			{
				name:     "valid secret",
				config:   SecretConfig{secretKey: "test-secret-key"},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.config.IsValid(); got != tt.expected {
					t.Errorf("SecretConfig.IsValid() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
