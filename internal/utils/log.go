package utils

import (
	"context"
	"go.uber.org/zap"
	"goload/internal/configs"
)

func getZapLoggerLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func InitializeLogger(logConfigs configs.Log) (*zap.Logger, func(), error) {
	zapLoggerConfig := zap.NewProductionConfig()
	zapLoggerConfig.Level = getZapLoggerLevel(logConfigs.Level)

	logger, err := zapLoggerConfig.Build()
	if err != nil {
		return nil, nil, err
	}

	cleanUp := func() {
		// deliberately ignore the returned error here
		_ = logger.Sync()
	}

	return logger, cleanUp, err
}

func LoggerWithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	return logger
}
