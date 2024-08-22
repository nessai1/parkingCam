package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

func BuildLogger(serviceName string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	logDir := filepath.Join("/var/log/", serviceName)

	touchLogFile := func(fileName string) error {
		fullPath := filepath.Join(logDir, fileName)
		_, err := os.Stat(fullPath)

		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to open log file for check: %w", err)
		} else if err == nil {
			return nil
		}

		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("failed to create log file: %w", err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}

		return nil
	}

	err := os.MkdirAll(logDir, 0555)
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	err = touchLogFile("debug.log")
	if err != nil {
		return nil, fmt.Errorf("failed to touch debug.log: %w", err)
	}

	err = touchLogFile("error.log")
	if err != nil {
		return nil, fmt.Errorf("failed to touch error.log: %w", err)
	}

	cfg.OutputPaths = []string{
		filepath.Join(logDir, "debug.log"),
		"stderr",
	}

	cfg.ErrorOutputPaths = []string{
		filepath.Join(logDir, "error.log"),
	}

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		MessageKey: "message",
		EncodeTime: zapcore.RFC3339TimeEncoder,
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("cannot build logger: %w", err)
	}

	return logger, nil
}
