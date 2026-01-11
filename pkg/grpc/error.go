package grpc

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcError(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "internal server error"
	}
	msg := st.Message()
	var httpStatus int
	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusUnprocessableEntity
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.DeadlineExceeded:
		httpStatus = http.StatusGatewayTimeout
	case codes.Internal, codes.Unknown:
		httpStatus = http.StatusInternalServerError
		msg = "internal server error"
	case codes.Unavailable:
		httpStatus = http.StatusServiceUnavailable
		msg = "service unavailable"
	default:
		httpStatus = http.StatusInternalServerError
	}

	return httpStatus, msg
}
