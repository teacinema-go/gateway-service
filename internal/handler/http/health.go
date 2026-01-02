package http

import (
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
)

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	h.SendResponse(http.StatusOK, w, response.Success("ok", map[string]any{
		"timestamp": time.Now().Format(time.DateTime),
	}))
}
