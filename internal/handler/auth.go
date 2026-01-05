package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/teacinema-go/core/http/response"
	"github.com/teacinema-go/gateway-service/internal/dto"
	"github.com/teacinema-go/gateway-service/internal/validator"
)

func (h *Handler) SendOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOtpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.SendResponse(w, http.StatusBadRequest, response.ErrorNoData("invalid json body"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		// TODO format validation messages
		h.SendResponse(w, http.StatusUnprocessableEntity, response.Error("validation failed", map[string]any{
			"error": err.Error(),
		}))
		return
	}

	if err := validator.ValidateIdentifierFormat(req); err != nil {
		h.SendResponse(w, http.StatusUnprocessableEntity, response.Error("validation failed", map[string]any{
			"error": err.Error(),
		}))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.clients.Auth.SendOtp(ctx, req)
	if err != nil || !res.GetOk() {
		h.logger.Error("sendOtp request failed", "error", err)
		h.SendResponse(w, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		return
	}

	h.SendResponse(w, http.StatusOK, response.SuccessNoData("ok"))
}
