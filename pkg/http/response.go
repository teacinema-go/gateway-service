package http

import (
	"encoding/json"
	"net/http"
)

type Logger interface {
	Error(msg string, args ...any)
}

func SendResponse(
	w http.ResponseWriter,
	log Logger,
	statusCode int,
	resp any,
) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(resp); err != nil {
		log.Error(
			"failed to write http response",
			"status", statusCode,
			"error", err,
		)
	}
}
