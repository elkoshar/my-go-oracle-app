package logger_test

import (
	"testing"

	config "oracle.com/oracle/my-go-oracle-app/configs"
	"oracle.com/oracle/my-go-oracle-app/pkg/logger"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "Test Success",
			config: &config.Config{
				LogLevel: "INFO",
			},
		},
		{
			name: "Test Success",
			config: &config.Config{
				LogLevel: "debug",
			},
		},
		{
			name: "Test Success",
			config: &config.Config{
				LogLevel: "warN",
			},
		},
		{
			name: "Test Success",
			config: &config.Config{
				LogLevel: "error",
			},
		},
		{
			name: "Test error",
			config: &config.Config{
				LogLevel: "INFOX",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.InitLogger(tt.config)
		})
	}
}
