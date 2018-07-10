package service

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Grpc struct {
	MaxReceiveSize int
	MaxSendSize    int

	ServerOptions []grpc.ServerOption
	Interceptors  []grpc.UnaryServerInterceptor

	server *grpc.Server
}

func (s *Service) GrpcServer() *grpc.Server {

	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(s.Grpc.MaxReceiveSize),
		grpc.MaxSendMsgSize(s.Grpc.MaxSendSize),
		grpc.Creds(credentials.NewClientTLSFromCert(s.CertPool, s.Host)),
	}

	options = append(options, s.Grpc.ServerOptions...)

	options = append(options, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			s.Grpc.Interceptors...,
		)))

	s.Grpc.server = grpc.NewServer(options...)

	return s.Grpc.server
}

func (g *Grpc) AddInterceptors(Interceptors ...grpc.UnaryServerInterceptor) {
	g.Interceptors = append(g.Interceptors, Interceptors...)
}

func (g *Grpc) AddOptions(opts ...grpc.ServerOption) {
	g.ServerOptions = append(g.ServerOptions, opts...)
}
