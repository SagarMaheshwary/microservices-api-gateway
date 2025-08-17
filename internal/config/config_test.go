package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigWithOptions(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	// Write a temporary .env file
	err := os.WriteFile(envFile, []byte("HTTP_PORT=5005\nAPP_ENV=staging\n"), 0644)
	require.NoError(t, err, "failed to write temp env file")

	tests := []struct {
		name           string
		envFilePath    string
		setupEnv       func()
		expectedPort   int
		expectedAppEnv string
	}{
		{
			name:           "loads from .env file",
			envFilePath:    envFile,
			setupEnv:       func() {}, // no system env
			expectedPort:   5005,
			expectedAppEnv: "staging",
		},
		{
			name:           "falls back to defaults when no .env file",
			envFilePath:    "nonexistent.env",
			setupEnv:       func() { os.Clearenv() },
			expectedPort:   4000,          // default
			expectedAppEnv: "development", // default
		},
		{
			name:        "overrides with system env vars",
			envFilePath: "",
			setupEnv: func() {
				os.Clearenv()
				os.Setenv("HTTP_PORT", "6006")
				os.Setenv("APP_ENV", "production")
			},
			expectedPort:   6006,
			expectedAppEnv: "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg := config.NewConfigWithOptions(config.LoaderOptions{EnvPath: tt.envFilePath})

			assert.Equal(t, tt.expectedPort, cfg.HTTPServer.Port, "HTTP port mismatch")
			assert.Equal(t, tt.expectedAppEnv, cfg.App.Env, "App env mismatch")

			// sanity check: durations and other defaults
			assert.NotEmpty(t, cfg.GRPCAuthenticationClient.URL, "GRPCAuthenticationClient.URL should not be empty")
			assert.Greater(t, cfg.GRPCAuthenticationClient.Timeout, 0*time.Second, "GRPCAuthenticationClient.Timeout should be > 0")
		})
	}
}
