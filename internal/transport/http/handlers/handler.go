package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/teacinema-go/core/constants"
	"github.com/teacinema-go/core/logger"
	"github.com/teacinema-go/gateway-service/internal/infra/grpc/clients"
	mw "github.com/teacinema-go/gateway-service/internal/transport/http/middlewares"
	validatorMW "github.com/teacinema-go/gateway-service/internal/transport/http/middlewares/validator"
)

type Handler struct {
	clients *clients.Manager
	env     constants.Env
}

func NewHandler(clients *clients.Manager, env constants.Env) *Handler {
	return &Handler{
		clients: clients,
		env:     env,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(logger.With()))
	r.Use(middleware.Recoverer)

	v := validator.New()

	r.Get("/health", h.Health)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.With(validatorMW.SendOtp(v)).Post("/otp/send", h.SendOtp)
				r.With(validatorMW.VerifyOtp(v)).Post("/otp/verify", h.VerifyOtp)
			})
		})
	})

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Debug(fmt.Sprintf("[%s]: %s", method, route))
		return nil
	})

	return r
}
