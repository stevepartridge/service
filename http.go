package service

import "net/http"

type httpHandler func(http.Handler) http.Handler

func (s *Service) AddHTTPMiddleware(handler httpHandler) {
	s.httpHandlers = append(s.httpHandlers, handler)
}

func WithHTTPMiddleware(handlers ...httpHandler) func(*Service) error {
	return func(s *Service) error {
		s.httpHandlers = append(s.httpHandlers, handlers...)
		return nil
	}
}

func (s *Service) chainHandlers(h http.Handler) http.Handler {
	for i := range s.httpHandlers {
		h = s.httpHandlers[len(s.httpHandlers)-1-i](h)
	}
	return h
}
