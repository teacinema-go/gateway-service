package handlers

import (
	"context"
	"net/http"
	"time"

	authv1 "github.com/teacinema-go/contracts/gen/go/auth/v1"
	"github.com/teacinema-go/core/http/response"
	"github.com/teacinema-go/core/logger"
	"github.com/teacinema-go/gateway-service/internal/auth/dto/request"
	appContextUtil "github.com/teacinema-go/gateway-service/pkg/context"
	"github.com/teacinema-go/gateway-service/pkg/grpc"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
)

func (h *Handler) SendOtp(w http.ResponseWriter, r *http.Request) {
	log := logger.With("method", "SendOtp")

	req, ok := appContextUtil.GetValue[request.SendOtpRequest](r.Context())
	if !ok {
		log.Error("failed to get request from context")
		pkgHTTP.SendResponse(w, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.clients.Auth.SendOtp(ctx, req)
	if err != nil {
		httpStatus, msg := grpc.HandleGrpcError(err)
		log.Error("gRPC request failed", "error", err, "message", msg, "status", httpStatus)
		pkgHTTP.SendResponse(w, httpStatus, response.ErrorNoData(msg))
		return
	}

	if !res.Success {
		log.Error("gRPC request failed", "error", res.ErrorMessage)
		switch res.ErrorCode {
		case authv1.SendOtpResponse_INVALID_IDENTIFIER_TYPE:
			pkgHTTP.SendResponse(w, http.StatusBadRequest, response.ErrorNoData(res.ErrorMessage))
		case authv1.SendOtpResponse_INTERNAL_ERROR, authv1.SendOtpResponse_ERROR_CODE_UNSPECIFIED:
			pkgHTTP.SendResponse(w, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		}

		return
	}

	pkgHTTP.SendResponse(w, http.StatusOK, response.Success("ok", map[string]any{
		"expires_in_seconds": res.OtpInfo.ExpiresInSeconds,
	}))
}

func (h *Handler) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	log := logger.With("method", "VerifyOtp")

	req, ok := appContextUtil.GetValue[request.VerifyOtpRequest](r.Context())
	if !ok {
		log.Error("failed to get request from context")
		pkgHTTP.SendResponse(w, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	res, err := h.clients.Auth.VerifyOtp(ctx, req)
	if err != nil {
		httpStatus, msg := grpc.HandleGrpcError(err)
		log.Error("gRPC request failed", "error", err, "message", msg, "status", httpStatus)
		pkgHTTP.SendResponse(w, httpStatus, response.ErrorNoData(msg))
		return
	}

	if !res.Success {
		log.Error("gRPC request failed", "error", res.ErrorMessage)
		switch res.ErrorCode {
		case
			authv1.VerifyOtpResponse_INVALID_IDENTIFIER_TYPE,
			authv1.VerifyOtpResponse_INVALID_OTP,
			authv1.VerifyOtpResponse_EXPIRED_OTP,
			authv1.VerifyOtpResponse_ACCOUNT_NOT_FOUND:
			pkgHTTP.SendResponse(w, http.StatusBadRequest, response.ErrorNoData(res.ErrorMessage))
		case authv1.VerifyOtpResponse_INTERNAL_ERROR, authv1.VerifyOtpResponse_ERROR_CODE_UNSPECIFIED:
			pkgHTTP.SendResponse(w, http.StatusInternalServerError, response.ErrorNoData("internal server error"))
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.Tokens.RefreshToken,
		Path:     "/", // TODO refresh route
		HttpOnly: true,
		Secure:   false, // TODO true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	pkgHTTP.SendResponse(w, http.StatusOK, response.Success("ok", map[string]any{
		"access_token":       res.Tokens.AccessToken,
		"expires_in_seconds": res.Tokens.ExpiresInSeconds,
	}))
}
