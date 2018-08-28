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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

type Service struct {
	Port int

	CertPool       *x509.CertPool
	PublicCert     []byte
	PrivateKey     []byte
	enableInsecure bool

	Grpc            *Grpc
	gatewayMux      *runtime.ServeMux
	gatewayHandlers []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

	Mux       *http.ServeMux
	httpChain alice.Chain
	Router    *chi.Mux
}

func New(port int) (*Service, error) {

	if port <= 0 {
		return nil, ErrReplacer(ErrInvalidPort, port)
	}

	s := Service{
		Port:     port,
		CertPool: x509.NewCertPool(),

		Mux:        http.NewServeMux(),
		Router:     chi.NewMux(),
		gatewayMux: runtime.NewServeMux(),
		httpChain:  alice.New(),
	}

	s.Grpc = &Grpc{
		MaxReceiveSize: math.MaxInt32,
		MaxSendSize:    math.MaxInt32,
	}

	return &s, nil
}

// func (s *Service) Addr() string {
// 	return fmt.Sprintf("%s:%d", s.Host, s.Port)
// }

func (s *Service) AddGatewayHandler(handler ...func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) error {

	if handler == nil {
		return ErrGatewayHandlerIsNil
	}

	if s.gatewayHandlers == nil {
		s.gatewayHandlers = []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{}
	}

	s.gatewayHandlers = append(s.gatewayHandlers, handler...)

	return nil
}

func (s *Service) AddHttpMiddlware(handler func(http.Handler) http.Handler) {
	s.httpChain = s.httpChain.Append(handler)
}

func (s *Service) Serve() error {

	if s.gatewayHandlers != nil {

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				}),
			),
		}

		for i := range s.gatewayHandlers {

			err := s.gatewayHandlers[i](
				context.Background(),
				s.gatewayMux,
				fmt.Sprintf("localhost:%d", s.Port),
				opts,
			)
			if err != nil {
				return err
			}

		}

		s.Mux.Handle("/", s.gatewayMux)
	}

	s.Router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Mux.ServeHTTP(w, r)
	}))

	httpServer := s.httpChain.Then(s.Router)

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
		Handler:   handlerFunc(s.Grpc.Server, httpServer),
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
