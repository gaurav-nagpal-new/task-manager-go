package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger() {
	// by this we can use zap.L() in our code for logging
	zap.ReplaceGlobals(zap.Must(createLogger()))
}

func createLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel), // Change to zap.DebugLevel
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	return config.Build()
}
