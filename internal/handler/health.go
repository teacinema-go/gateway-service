package handler

import (
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
)

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	pkgHTTP.SendResponse(w, http.StatusOK, response.Success("ok", map[string]any{
		"timestamp": time.Now().Format(time.DateTime),
	}))
}
