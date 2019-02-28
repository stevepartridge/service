package middleware

import (
	"time"

	"golang.org/x/net/context" // have to use this context because of grpc lib
	"google.golang.org/grpc"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"
)

// RequestInterceptor is an example request logger middleware
func RequestInterceptor() func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		log.Info().
			Str("request_method", info.FullMethod).
			Str("user_agent", getUserAgentFromContext(ctx)).
			Str("ip", getForwardedForFromContext(ctx)).
			Str("referer", getForwardedHostFromContext(ctx)).
			Msg("")

		resp, err := handler(ctx, req)

		return resp, err
	}
}

// TelemetryInterceptor is an example telementry logger middleware
func TelemetryInterceptor() func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		started := time.Now().UTC()

		resp, err := handler(ctx, req)

		finished := time.Now().UTC()

		log.Info().
			Str("request_method", info.FullMethod).
			Str("start", started.Format(time.RFC3339Nano)).
			Str("finish", finished.Format(time.RFC3339Nano)).
			Str("duration", finished.Sub(started).String()).
			Str("duration", finished.Sub(started).String()).
			Msg("")

		return resp, err
	}
}
