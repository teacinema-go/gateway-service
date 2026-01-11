package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/teacinema-go/core/logger"
)

func SendResponse(
	w http.ResponseWriter,
	statusCode int,
	resp any,
) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(resp); err != nil {
		logger.Error(
			"failed to write http response",
			"time", time.Now(),
			"status", statusCode,
			"error", err,
		)
	}
}
