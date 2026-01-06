package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
)

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	log := h.logger.With(slog.String("method", "Health"))

	pkgHTTP.SendResponse(w, log, http.StatusOK, response.Success("ok", map[string]any{
		"timestamp": time.Now().Format(time.DateTime),
	}))
}
