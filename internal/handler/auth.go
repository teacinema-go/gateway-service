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
		h.SendResponse(http.StatusBadRequest, w, response.ErrorNoData("invalid json body"))
		return
	}

	if err := h.v.Struct(req); err != nil {
		// TODO format validation messages
		h.SendResponse(http.StatusUnprocessableEntity, w, response.Error("validation failed", map[string]any{
			"error": err.Error(),
		}))
		return
	}

	if err := validator.ValidateIdentifierFormat(req); err != nil {
		h.SendResponse(http.StatusUnprocessableEntity, w, response.Error("validation failed", map[string]any{
			"error": err.Error(),
		}))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.services.Auth.SendOtp(ctx, req)
	if err != nil || !res.GetOk() {
		h.l.Error("sendOtp request failed", "error", err)
		h.SendResponse(http.StatusInternalServerError, w, response.ErrorNoData("internal server error"))
		return
	}

	h.SendResponse(http.StatusOK, w, response.SuccessNoData("ok"))
}
