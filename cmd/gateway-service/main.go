package main

import (
	"log"

	"github.com/teacinema-go/core/logger"
	"github.com/teacinema-go/gateway-service/internal/app"
	"github.com/teacinema-go/gateway-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	logger.Init(cfg.App.Env)
	logger.Info("config loaded successfully", "env", cfg.App.Env)

	application := app.New(cfg)

	if err = application.Run(); err != nil {
		logger.Error("application stopped with error", "error", err)
	}
}
