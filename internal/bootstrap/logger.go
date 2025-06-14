package bootstrap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(cfg *AppConfig) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	level := cfg.LogLevel
	if level == "" {
		level = "info"
	}

	err := config.Level.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return config.Build()
}
