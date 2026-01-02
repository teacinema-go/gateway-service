package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	mw "github.com/teacinema-go/gateway-service/internal/middleware"
)

type Handler struct {
	l *slog.Logger
	v *validator.Validate
}

func NewHandler(l *slog.Logger) *Handler {
	v := validator.New()
	return &Handler{
		l: l,
		v: v,
	}
}

func (h *Handler) SendResponse(statusCode int, w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.l.Error("error encoding response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(h.l))
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
		h.l.Debug(fmt.Sprintf("[%s]: %s", method, route))
		return nil
	})

	return r
}
