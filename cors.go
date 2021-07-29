package service

import (
	"github.com/rs/cors"
)

var (
	defaultCorsAllowedOrigins = []string{"*"}
	defaultCorsAllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	defaultCorsAllowedHeaders = []string{
		"Access-Control-Allow-Origin",
		"Content-Type",
		"Authorization",
	}
)

// CORS sets up the Cross-Origin settings if the current configuration asks for it
func (s *Service) cors() {

	if envBool(EnvCorsEnable) {

		opts := cors.Options{
			AllowedOrigins: defaultCorsAllowedOrigins,
			AllowedMethods: defaultCorsAllowedMethods,
			AllowedHeaders: defaultCorsAllowedHeaders,
			Debug:          envBool(EnvCorsDebug),
		}

		if len(envStrings(EnvCorsAllowedOrigins)) > 0 {
			opts.AllowedOrigins = envStrings(EnvCorsAllowedOrigins)
		}

		if len(envStrings(EnvCorsAllowedMethods)) > 0 {
			opts.AllowedMethods = envStrings(EnvCorsAllowedMethods)
		}

		if len(envStrings(EnvCorsAllowedHeaders)) > 0 {
			opts.AllowedHeaders = envStrings(EnvCorsAllowedHeaders)
		}

		s.AddHTTPMiddleware(cors.New(opts).Handler)
	}

}
