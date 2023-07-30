package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	startTime := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("accept gRPC request")

	return
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rc *ResponseRecorder) WriteHeader(statusCode int) {
	rc.StatusCode = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}
func (rc *ResponseRecorder) Write(body []byte) (int, error) {
	rc.Body = body
	return rc.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rc := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rc, r)
		duration := time.Since(startTime)

		logger := log.Info()
		if rc.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rc.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Int("status_code", rc.StatusCode).
			Str("status_text", http.StatusText(rc.StatusCode)).
			Dur("duration", duration).
			Msg("accept HTTP request")

	})
}
