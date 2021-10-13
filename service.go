package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Defaults
var (
	defaultPort = 8000
)

type Service struct {
	HTTP   *http.ServeMux
	Router *chi.Mux

	CertPool   *x509.CertPool
	PublicCert []byte
	PrivateKey []byte

	port int
	grpc *grpc.Server

	server *http.Server

	gatewayMux *runtime.ServeMux
	grpcPort   int // to serve grpc on a separate port

	maxReceiveSize int
	maxSendSize    int

	insecureSkipVerify bool
	disableTLSCerts    bool

	httpHandlers []httpHandler

	gatewayHandlers []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

	grpcServerOptions  []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	debug bool
}

type Option func(*Service) error

func New(opts ...Option) (*Service, error) {

	s := Service{
		port:     envInt(EnvPort),
		HTTP:     http.NewServeMux(),
		Router:   chi.NewMux(),
		CertPool: x509.NewCertPool(),

		gatewayMux: runtime.NewServeMux(),

		gatewayHandlers: []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{},
		httpHandlers:    []httpHandler{},

		grpcServerOptions:  []grpc.ServerOption{},
		unaryInterceptors:  []grpc.UnaryServerInterceptor{},
		streamInterceptors: []grpc.StreamServerInterceptor{},

		maxReceiveSize: math.MaxInt32,
		maxSendSize:    math.MaxInt32,

		grpcPort:           envInt(EnvGRPCPort),
		insecureSkipVerify: envBool(EnvInsecureVerifySkip),
		disableTLSCerts:    envBool(EnvTLSDisabled),
	}

	if envBool(EnvDebug) {
		opts = append(opts, WithDebug())
	}

	err := s.WithOptions(opts...)
	if err != nil {
		return nil, err
	}

	return &s, nil

}

func (s *Service) Serve(ctx context.Context) error {

	if s.debug {
		s.PrintDebug()
		log := grpclog.NewLoggerV2(os.Stdout, os.Stdout, ioutil.Discard)
		grpclog.SetLoggerV2(log)

	}

	if s.port < 1 || s.port > 65535 {
		fmt.Printf("Invalid Port %d falling back to default port %d\n", s.port, defaultPort)
		s.port = defaultPort
	}

	if s.grpcPort > 65535 {
		fmt.Printf("Invalid GRPC Port %d falling back to primary port %d\n", s.grpcPort, s.port)
		s.grpcPort = 0
	}

	if s.grpc == nil {
		return ErrServeGRPCNotYetDefined
	}

	if s.disableTLSCerts && s.grpcPort == 0 {
		return ErrDisableTLSCertsMissingGRPCPort
	}

	cert, err := s.GetCertificate()
	if err != nil {
		return err
	}

	tlsConfig := tls.Config{
		NextProtos:         []string{"h2"},
		InsecureSkipVerify: s.insecureSkipVerify,
	}

	grpcPort := s.port
	if s.grpcPort > 0 {
		grpcPort = s.grpcPort
	}

	// grpcAddr := fmt.Sprintf(":%d", grpcPort)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}

	// Serving HTTP/1 and gRPC on separate ports
	if s.grpcPort > 0 && s.grpcPort != s.port {

		go func() {

			if s.disableTLSCerts {
				fmt.Printf("Serving gRPC (No TLS) on Port: %d\n", s.grpcPort)
				err := s.grpc.Serve(conn)
				fmt.Println("error serving grpc: ", err.Error())
				return
			}

			fmt.Printf("Serving gRPC on Port: %d\n", s.grpcPort)
			err := s.grpc.Serve(tls.NewListener(conn, &tlsConfig))
			fmt.Println("error serving grpc: ", err.Error())

		}()

	}

	if s.gatewayHandlers != nil {

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true, // talk to grpc within the service at localhost/127.0.0.1
				}),
			),
		}

		if s.disableTLSCerts {
			opts = []grpc.DialOption{
				// grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(s.CertPool, "")),
				grpc.WithInsecure(),
				grpc.WithBlock(),
			}
		}

		clientConn, err := grpc.DialContext(
			context.Background(),
			fmt.Sprintf("dns:///127.0.0.1:%d", grpcPort),
			opts...,
		)
		if err != nil {
			return err
		}

		for i := range s.gatewayHandlers {

			err := s.gatewayHandlers[i](
				ctx,
				s.gatewayMux,
				clientConn,
			)
			if err != nil {
				return err
			}

		}

		s.HTTP.Handle("/", s.gatewayMux)
	}

	s.Router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.HTTP.ServeHTTP(w, r)
	}))

	s.cors()

	s.server = &http.Server{
		Addr:    strconv.Itoa(s.port),
		Handler: handlerFunc(s.grpc, s.chainHandlers(s.Router)),
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	// tlsConfig.BuildNameToCertificate()

	if s.grpcPort > 0 && s.grpcPort != s.port {

		c, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err != nil {
			return err
		}

		if s.disableTLSCerts {
			fmt.Printf("Serving HTTP on Port: %d\n", s.port)
			return s.server.Serve(c)
		}

		fmt.Printf("Serving HTTPS on Port: %d\n", s.port)
		return s.server.Serve(tls.NewListener(c, &tlsConfig))
	}

	fmt.Printf("Serving HTTPS and gRPC on Port: %d\n", s.port)
	return s.server.Serve(tls.NewListener(conn, &tlsConfig))
}

// Shutdown attempts to gracefully shutdown server given a context timeout
func (s *Service) Shutdown(ctx context.Context) error {
	defer s.grpc.GracefulStop()
	if s.server == nil {
		return errors.New("service server is nil")
	}
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	s.grpc.Stop()
	return nil
}

func (s *Service) WithOptions(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}
	return nil
}

func WithPort(port int) func(*Service) error {
	return func(s *Service) error {
		if port <= 0 {
			return fmt.Errorf(ErrInvalidPort.Error(), port)
		}
		s.port = port
		return nil
	}
}

func WithDebug() func(*Service) error {
	return func(s *Service) error {
		s.debug = true
		return nil
	}
}

func (s *Service) PrintDebug() {
	fmt.Printf(`
=== Service Info =====================

Port             : %d
Gateway Handlers : %d
HTTP Handlers    : %d

--------------------------------------
--- GRPC -----------------------------

Port : %d
Service Options     : %d
Unary Interceptors  : %d
Stream Interceptors : %d
Max Receive Size    : %d
Max Send Size       : %d

--------------------------------------
--- TLS ------------------------------

Insecure Verify Skip : %t
Disable TLS          : %t

======================================
`,
		s.port,
		len(s.gatewayHandlers),
		len(s.httpHandlers),
		s.grpcPort,
		len(s.grpcServerOptions),
		len(s.unaryInterceptors),
		len(s.streamInterceptors),
		s.maxReceiveSize,
		s.maxSendSize,
		s.insecureSkipVerify,
		s.disableTLSCerts,
	)
}

func handlerFunc(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
