package middlewares

import (
	"net/http"
	"strings"

	"github.com/teacinema-go/core/http/response"
	"github.com/teacinema-go/gateway-service/internal/auth/dto/request"
	appContextUtil "github.com/teacinema-go/gateway-service/pkg/context"
	pkgHTTP "github.com/teacinema-go/gateway-service/pkg/http"
	"github.com/teacinema-go/passport"
)

func BearerAuth(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				pkgHTTP.SendResponse(w, http.StatusUnauthorized, response.ErrorNoData("unauthorized"))
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				pkgHTTP.SendResponse(w, http.StatusUnauthorized, response.ErrorNoData("unauthorized"))
				return
			}

			token, err := passport.ParseToken(parts[1])
			if err != nil {
				pkgHTTP.SendResponse(w, http.StatusUnauthorized, response.ErrorNoData("unauthorized"))
				return
			}

			verified := token.VerifyToken(secretKey)
			if !verified {
				pkgHTTP.SendResponse(w, http.StatusUnauthorized, response.ErrorNoData("unauthorized"))
				return
			}

			ctx := appContextUtil.WithValue[request.UserID](r.Context(), request.UserID(token.UserID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
