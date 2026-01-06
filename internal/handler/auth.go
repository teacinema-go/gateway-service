package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"github.com/teacinema-go/gateway-service/pkg/grpc"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
)

func (h *Handler) SendOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOtpRequest
	log := h.logger.With(slog.String("method", "SendOtp"))
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgHTTP.SendResponse(w, log, http.StatusBadRequest, response.ErrorNoData("invalid json body"))
		return
	}

	if err := req.Validate(h.validator); err != nil {
		pkgHTTP.SendResponse(w, log, http.StatusUnprocessableEntity, response.Error("validation failed", map[string]any{
			"error": err.Error(),
		}))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.clients.Auth.SendOtp(ctx, req)
	if err != nil {
		httpStatus, msg := grpc.HandleGrpcError(err)
		log.Error("gRPC request failed", "error", err, "message", msg, "status", httpStatus)
		pkgHTTP.SendResponse(w, log, httpStatus, response.ErrorNoData(msg))
		return
	}

	if !res.GetOk() {
		log.Warn("gRPC request rejected by auth service")
		pkgHTTP.SendResponse(w, log, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		return
	}

	pkgHTTP.SendResponse(w, log, http.StatusOK, response.SuccessNoData("ok"))
}
