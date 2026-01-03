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

	"github.com/teacinema-go/gateway-service/internal/config"
	"github.com/teacinema-go/gateway-service/internal/handler"
	"github.com/teacinema-go/gateway-service/internal/service"
)

type App struct {
	cfg        *config.Config
	l          *slog.Logger
	httpServer *http.Server
}

func New(cfg *config.Config, l *slog.Logger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

func (a *App) Run() error {
	serviceManager, err := service.NewServiceManager(a.cfg.Service)
	if err != nil {
		return fmt.Errorf("failed to create service manager: %w", err)
	}
	defer func() {
		if err := serviceManager.Close(); err != nil {
			a.l.Error("failed to close service manager", "error", err)
		}
	}()

	h := handler.NewHandler(a.l, serviceManager)

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
		a.l.Info("starting HTTP server",
			"port", a.cfg.App.Port,
			"host", a.cfg.App.Host,
			"env", a.cfg.App.Env,
		)

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.l.Error("HTTP server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	sig := <-quit
	a.l.Info("received shutdown signal", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	a.l.Info("server stopped gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.l.Info("shutting down server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	return nil
}
