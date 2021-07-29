package service

import (
	"fmt"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (s *Service) GRPC() *grpc.Server {

	if s.grpc != nil {
		return s.grpc
	}

	opts := []grpc.ServerOption{}

	if s.disableTLSCerts {
		opts = append(opts, grpc.Creds(
			credentials.NewServerTLSFromCert(s.insecureCert()),
		))
	}

	opts = append(opts,
		grpc.MaxRecvMsgSize(s.maxReceiveSize),
		grpc.MaxSendMsgSize(s.maxSendSize),
	)

	opts = append(opts, s.grpcServerOptions...)

	opts = append(opts, grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			s.unaryInterceptors...,
		),
	))

	opts = append(opts, grpc.StreamInterceptor(
		grpcMiddleware.ChainStreamServer(
			s.streamInterceptors...,
		),
	))

	s.grpc = grpc.NewServer(opts...)

	return s.grpc
}

func WithGRPCPort(port int) func(*Service) error {
	return func(s *Service) error {
		if port <= 0 {
			return fmt.Errorf(ErrInvalidPort.Error(), port)
		}
		s.grpcPort = port
		return nil
	}
}

func WithMaxRecieveSize(size int) func(*Service) error {
	return func(s *Service) error {
		if size <= 0 {
			return fmt.Errorf(ErrInvalidRecieveSize.Error(), size)
		}
		return nil
	}
}

func WithMaxSendSize(size int) func(*Service) error {
	return func(s *Service) error {
		if size <= 0 {
			return fmt.Errorf(ErrInvalidSendSize.Error(), size)
		}
		return nil
	}
}

func WithGRPCServerOptions(opts ...grpc.ServerOption) func(*Service) error {
	return func(s *Service) error {
		if len(opts) == 0 {
			return ErrMissingGrpcServerOptions
		}

		s.grpcServerOptions = append(s.grpcServerOptions, opts...)

		return nil
	}
}

func WithGRPCServer(server *grpc.Server) func(*Service) error {
	return func(s *Service) error {
		if server == nil {
			return ErrWithGRPCServerIsNil
		}
		s.grpc = server
		return nil
	}
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) func(*Service) error {
	return func(s *Service) error {
		if len(interceptors) == 0 {
			return ErrMissingGrpcUnaryInterceptors
		}

		s.unaryInterceptors = interceptors

		return nil
	}
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) func(*Service) error {
	return func(s *Service) error {
		if len(interceptors) == 0 {
			return ErrMissingGrpcStreamInterceptors
		}

		s.streamInterceptors = interceptors

		return nil
	}
}

func (s *Service) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) error {
	if len(interceptors) == 0 {
		return ErrMissingGrpcUnaryInterceptors
	}

	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)

	return nil
}

func (s *Service) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) error {
	if len(interceptors) == 0 {
		return ErrMissingGrpcStreamInterceptors
	}

	s.streamInterceptors = append(s.streamInterceptors, interceptors...)

	return nil
}
