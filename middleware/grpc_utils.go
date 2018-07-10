package middleware

import (
	"golang.org/x/net/context" // have to use this context because of grpc lib
	"google.golang.org/grpc/metadata"
)

func getUserAgentFromContext(ctx context.Context) string {

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if key, ok := md["grpcgateway-user-agent"]; ok && len(key) > 0 {
			return key[0]
		}

		if key, ok := md["user-agent"]; ok && len(key) > 0 {
			return key[0]
		}
	}

	return ""
}

func getForwardedForFromContext(ctx context.Context) string {

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if key, ok := md["x-forwarded-for"]; ok && len(key) > 0 {
			return key[0]
		}
	}

	return ""
}

func getForwardedHostFromContext(ctx context.Context) string {

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		if key, ok := md["x-forwarded-host"]; ok && len(key) > 0 {
			return key[0]
		}
	}
	return ""
}
