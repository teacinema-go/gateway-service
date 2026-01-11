package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/teacinema-go/core/logger"
	"github.com/teacinema-go/gateway-service/internal/clients"
	mw "github.com/teacinema-go/gateway-service/internal/middleware"
)

type Handler struct {
	validator *validator.Validate
	clients   *clients.Manager
}

func NewHandler(clients *clients.Manager) *Handler {
	v := validator.New()
	return &Handler{
		validator: v,
		clients:   clients,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(logger.With()))
	r.Use(middleware.Recoverer)

	r.Get("/health", h.Health)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/otp/send", h.SendOtp)
				r.Post("/otp/verify", h.VerifyOtp)
			})
		})
	})

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Debug(fmt.Sprintf("[%s]: %s", method, route))
		return nil
	})

	return r
}
