package http

import (
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse("ok", response.Data{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})

	h.l.Debug("health check requested",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	h.SendResponse(http.StatusOK, w, resp)
}
