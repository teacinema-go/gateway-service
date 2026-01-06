package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/teacinema-go/gateway-service/internal/clients"
	mw "github.com/teacinema-go/gateway-service/internal/middleware"
)

type Handler struct {
	logger    *slog.Logger
	validator *validator.Validate
	clients   *clients.Manager
}

func NewHandler(logger *slog.Logger, clients *clients.Manager) *Handler {
	v := validator.New()
	return &Handler{
		logger:    logger,
		validator: v,
		clients:   clients,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(h.logger))
	r.Use(middleware.Recoverer)

	r.Get("/health", h.Health)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/otp/send", h.SendOtp)
			})
		})
	})

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		h.logger.Debug(fmt.Sprintf("[%s]: %s", method, route))
		return nil
	})

	return r
}
