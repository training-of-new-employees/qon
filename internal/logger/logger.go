// Package logger provides a logger.
package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log доступен во всём коде (паттерн Singleton).
var Log *zap.Logger = zap.NewNop()

// InitLogger инициализирует логгер.
func InitLogger(level string) error {
	// преобразование уровня логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	cfg.EncoderConfig.EncodeTime = CustomMillisTimeEncoder

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	// установить Singleton
	Log = zl

	return nil
}

// CustomMillisTimeEncoder - преобразователь времени в формат ОС Unix.
func CustomMillisTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(`2006-01-02T15:04:05.000207`))
}
