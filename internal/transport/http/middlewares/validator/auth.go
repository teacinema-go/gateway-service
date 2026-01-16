package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/teacinema-go/core/http/response"
	"github.com/teacinema-go/gateway-service/internal/auth/dto/request"
	appContextUtil "github.com/teacinema-go/gateway-service/pkg/context"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
)

func SendOtp(v *validator.Validate) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, valErr := DecodeAndValidate[request.SendOtpRequest](r, v)
			if valErr != nil {
				if valErr.Fields != nil {
					pkgHTTP.SendResponse(w, valErr.HTTPStatus, response.Error(valErr.Message, map[string]any{
						"fields": valErr.Fields,
					}))
				} else {
					pkgHTTP.SendResponse(w, valErr.HTTPStatus, response.ErrorNoData(valErr.Message))
				}
				return
			}

			if err := req.IdentifierType.Validate(req.Identifier); err != nil {
				pkgHTTP.SendResponse(w, http.StatusUnprocessableEntity, response.Error("validation failed", map[string]any{
					"fields": map[string]string{
						"identifier": err.Error(),
					},
				}))
				return
			}

			ctx := appContextUtil.WithValue[request.SendOtpRequest](r.Context(), req)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func VerifyOtp(v *validator.Validate) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, valErr := DecodeAndValidate[request.VerifyOtpRequest](r, v)
			if valErr != nil {
				if valErr.Fields != nil {
					pkgHTTP.SendResponse(w, valErr.HTTPStatus, response.Error(valErr.Message, map[string]any{
						"fields": valErr.Fields,
					}))
				} else {
					pkgHTTP.SendResponse(w, valErr.HTTPStatus, response.ErrorNoData(valErr.Message))
				}
				return
			}

			if err := req.IdentifierType.Validate(req.Identifier); err != nil {
				pkgHTTP.SendResponse(w, http.StatusUnprocessableEntity, response.Error("validation failed", map[string]any{
					"fields": map[string]string{
						"identifier": err.Error(),
					},
				}))
				return
			}

			ctx := appContextUtil.WithValue[request.VerifyOtpRequest](r.Context(), req)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
