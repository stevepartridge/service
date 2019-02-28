package service

import (
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Grpc
type Grpc struct {
	MaxReceiveSize int
	MaxSendSize    int

	ServerOptions      []grpc.ServerOption
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor

	Server *grpc.Server
}

// GrpcServer creates and returns the server after applying ServerOptions
func (s *Service) GrpcServer() *grpc.Server {

	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(s.Grpc.MaxReceiveSize),
		grpc.MaxSendMsgSize(s.Grpc.MaxSendSize),
		grpc.Creds(credentials.NewClientTLSFromCert(s.CertPool, "")),
	}

	options = append(options, s.Grpc.ServerOptions...)

	options = append(options, grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			s.Grpc.UnaryInterceptors...,
		)))

	options = append(options, grpc.StreamInterceptor(
		grpcMiddleware.ChainStreamServer(
			s.Grpc.StreamInterceptors...,
		)))

	s.Grpc.Server = grpc.NewServer(options...)

	return s.Grpc.Server
}

// AddUnaryInterceptors is an exposed method to append unary interceptors
func (g *Grpc) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	g.UnaryInterceptors = append(g.UnaryInterceptors, interceptors...)
}

// AddStreamInterceptors is an exposed method to append stream interceptors
func (g *Grpc) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	g.StreamInterceptors = append(g.StreamInterceptors, interceptors...)
}

// AddOptions is an exposed method to append options
func (g *Grpc) AddOptions(opts ...grpc.ServerOption) {
	g.ServerOptions = append(g.ServerOptions, opts...)
}
