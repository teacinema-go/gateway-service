package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
		log.Println(
			"failed to write http response",
			"time", time.Now(),
			"status", statusCode,
			"error", err,
		)
	}
}
