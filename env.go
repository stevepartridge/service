package service

import (
	"os"
	"strconv"
	"strings"
)

const (
	// EnvPort is the primary port to serve everything on.
	EnvPort = "PORT"

	// EnvGRPCPort serves GRPC on a separate port
	EnvGRPCPort = "GRPC_PORT"

	// EnvDebug enables debug output
	EnvDebug = "SERVICE_DEBUG"

	// EnvInsecureVerifySkip disables TLS verification
	EnvInsecureVerifySkip = "INSECURE_VERIFY_SKIP"

	// EnvTLSDisabled will disable TLS (will only support HTTP/1 requests)
	EnvTLSDisabled = "TLS_DISABLED"

	// EnvCorsEnable env var (default: false)
	EnvCorsEnable = "CORS_ENABLE"

	// EnvCorsDebug env var (default: false)
	EnvCorsDebug = "CORS_DEBUG"

	// EnvCorsAllowedOrigins env var
	EnvCorsAllowedOrigins = "CORS_ALLOWED_ORIGINS"

	// EnvCorsAllowedMethods env var
	EnvCorsAllowedMethods = "CORS_ALLOWED_METHODS"

	// EnvCorsAllowedHeaders env var
	EnvCorsAllowedHeaders = "CORS_ALLOWED_HEADERS"
)

func envBool(env string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(env)))
	switch v {
	case "false", "0", "":
		return false
	}

	return true
}

func envInt(env string) int {
	val, _ := strconv.Atoi(os.Getenv(env))
	return val
}

func envStrings(env string) []string {
	var (
		v      = os.Getenv(env)
		result []string
	)

	switch {
	case strings.Contains(env, ","):
		result = strings.Split(v, ",")
	case strings.Contains(env, "|"):
		result = strings.Split(v, "|")
	case strings.Contains(env, ":"):
		result = strings.Split(v, ":")
	}

	return result
}
