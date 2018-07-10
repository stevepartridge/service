package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/justinas/alice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Service struct {
	Host string
	Port int

	CertPool       *x509.CertPool
	PublicCert     []byte
	PrivateKey     []byte
	enableInsecure bool

	Grpc           *Grpc
	Mux            *http.ServeMux
	gatewayMux     *runtime.ServeMux
	gatewayHandler func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

	httpChain alice.Chain
}

func New(host string, port int) (*Service, error) {

	if host == "" {
		return nil, ErrInvalidHost
	}

	if port <= 0 {
		return nil, ErrReplacer(ErrInvalidPort, port)
	}

	s := Service{
		Host:       host,
		Port:       port,
		CertPool:   x509.NewCertPool(),
		httpChain:  alice.New(),
		Mux:        http.NewServeMux(),
		gatewayMux: runtime.NewServeMux(),
	}

	s.Grpc = &Grpc{
		MaxReceiveSize: math.MaxInt32,
		MaxSendSize:    math.MaxInt32,
	}

	return &s, nil
}

func (s *Service) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *Service) EnableGatewayHandler(handler func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) error {

	if handler == nil {
		return ErrGatewayHandlerIsNil
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:         s.Host,
		RootCAs:            s.CertPool,
		InsecureSkipVerify: s.enableInsecure,
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err := handler(context.Background(), s.gatewayMux, s.Addr(), opts)
	if err != nil {
		return err
	}

	s.Mux.Handle("/", s.gatewayMux)

	return nil
}

func (s *Service) AddHttpHandler(handler func(http.Handler) http.Handler) {
	s.httpChain = s.httpChain.Append(handler)
}

func (s *Service) Serve() error {

	httpServer := s.httpChain.Then(s.Mux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		panic(err)
	}

	cert, err := s.GetCertificate()
	if err != nil {
		return err
	}

	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"h2"},
		InsecureSkipVerify: s.enableInsecure,
	}

	tlsConfig.BuildNameToCertificate()

	srv := &http.Server{
		Addr:      strconv.Itoa(s.Port),
		Handler:   handlerFunc(s.Grpc.server, httpServer),
		TLSConfig: &tlsConfig,
	}

	return srv.Serve(tls.NewListener(conn, srv.TLSConfig))

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
