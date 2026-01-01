package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/teacinema-go/core/http/response"
	mw "github.com/teacinema-go/gateway-service/internal/middleware"
)

type Handler struct {
	l *slog.Logger
}

func NewHandler(l *slog.Logger) *Handler {
	return &Handler{
		l: l,
	}
}

func (h *Handler) SendResponse(statusCode int, w http.ResponseWriter, resp *response.Response) {
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

		})
	})

	return r
}
