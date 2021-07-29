package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func (s *Service) AddGatewayHandlers(handlers ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) error {

	if len(handlers) == 0 {
		return ErrMissingGatewayHandlers
	}

	s.gatewayHandlers = append(s.gatewayHandlers, handlers...)

	return nil

}

func WithGatewayHandlers(handlers ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) func(*Service) error {
	return func(s *Service) error {
		return s.AddGatewayHandlers(handlers...)
	}
}

func WithServerMuxOptions(opts ...runtime.ServeMuxOption) func(*Service) error {
	return func(s *Service) error {
		s.gatewayMux = runtime.NewServeMux(opts...)
		return nil
	}
}
