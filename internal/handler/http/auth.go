package http

import (
	"encoding/json"
	"net/http"

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

	h.SendResponse(http.StatusOK, w, response.Success("ok", map[string]any{
		"auth": req,
	}))
}
