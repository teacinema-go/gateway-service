package http

import (
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
)

type HealthData struct {
	Timestamp string `json:"timestamp"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("health check requested",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	h.SendResponse(http.StatusOK, w, response.Success("ok", HealthData{time.Now().Format(time.DateTime)}))
}
