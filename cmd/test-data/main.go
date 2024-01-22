package main

import (
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
)

func main() {
	logger.InitLogger("debug")
	cfg := InitFlags()
	logger.Log.Info("Config", zap.Any("config", cfg))
	err := upTestEnv(cfg)
	if err != nil {
		logger.Log.Error("Error", zap.Error(err))
	}
}
