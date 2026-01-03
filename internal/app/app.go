package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teacinema-go/gateway-service/internal/client/grpc"
	"github.com/teacinema-go/gateway-service/internal/config"
	"github.com/teacinema-go/gateway-service/internal/handler"
)

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	authClient, err := grpc.NewAuthClient(a.cfg.Service.AuthServiceURL, a.logger)
	if err != nil {
		a.logger.Error("could not create auth client", "error", err)
	}
	defer func(authClient *grpc.AuthClient) {
		err = authClient.Close()
		if err != nil {
			a.logger.Warn("Failed to close Auth Client")
		}
	}(authClient)

	h := handler.NewHandler(a.logger, authClient)

	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:      h.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info("starting HTTP server",
			"port", a.cfg.App.Port,
			"host", a.cfg.App.Host,
			"env", a.cfg.App.Env,
		)

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	sig := <-quit
	a.logger.Info("received shutdown signal", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	a.logger.Info("server stopped gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	return nil
}
