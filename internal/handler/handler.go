package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	mw "github.com/teacinema-go/gateway-service/internal/middleware"
	"github.com/teacinema-go/gateway-service/internal/service"
)

type Handler struct {
	logger    *slog.Logger
	validator *validator.Validate
	services  *service.Manager
}

func NewHandler(logger *slog.Logger, services *service.Manager) *Handler {
	v := validator.New()
	return &Handler{
		logger:    logger,
		validator: v,
		services:  services,
	}
}

func (h *Handler) SendResponse(w http.ResponseWriter, statusCode int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("error encoding response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
