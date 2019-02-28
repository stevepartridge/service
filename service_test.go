package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/context"

	pb "github.com/stevepartridge/service/example/protos"
	"github.com/stevepartridge/service/insecure"
)

var (
	testHost1 = "service.local"
	testPort1 = 1234
)

type exampleService struct{}

func (s *exampleService) Info(c context.Context, req *pb.ServiceInfoRequest) (*pb.ServiceInfoResponse, error) {

	return &pb.ServiceInfoResponse{
		Name:    "serviceName",
		BuiltAt: "builtAt",
		Version: "version",
		Build:   "build",
		GitHash: "githash",
	}, nil

}

func Test_Unit_ServiceNew_ValidAddress(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc.Port != testPort1 {
		t.Errorf("Addr should be %d but saw %d", testPort1, svc.Port)
	}

}

func Test_Unit_ServiceAddGatewayHandler_InvalidNil(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	if svc.Port != testPort1 {
		t.Errorf("Addr should be %d but saw %d", testPort1, svc.Port)
	}

	err = svc.AddGatewayHandler()
	if err == nil {
		t.Error("Expected error but received none")
	}

	if err.Error() != ErrGatewayHandlerIsNil.Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrGatewayHandlerIsNil.Error(),
			err.Error(),
		)
	}
}

func Test_Unit_ServiceNew_InvalidPort(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if err.Error() != ErrReplacer(ErrInvalidPort, 0).Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrReplacer(ErrInvalidPort, 0).Error(),
			err.Error(),
		)
	}

}

func Test_Unit_ServiceNew_InvalidNegativePort(t *testing.T) {
	_, err := New(-123)
	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if err.Error() != ErrReplacer(ErrInvalidPort, -123).Error() {
		t.Errorf("Expected error %s but saw %s",
			ErrReplacer(ErrInvalidPort, -123).Error(),
			err.Error(),
		)
	}

}

func Test_Unit_ServiceServe_Success(t *testing.T) {

	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Error creating service %s", err.Error())
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          true,
	})
	svc.AddHttpMiddlware(c.Handler)

	err = svc.AddKeyPair([]byte(insecure.Cert), []byte(insecure.Key))
	if err != nil {
		t.Errorf("Error adding key pair %s", err.Error())
	}

	err = svc.AppendCertsFromPEM([]byte(insecure.RootCA))
	if err != nil {
		t.Errorf("Error append certs from pem %s", err.Error())
	}

	pb.RegisterExampleServer(svc.GrpcServer(), &exampleService{})

	err = svc.AddGatewayHandler(pb.RegisterExampleHandlerFromEndpoint)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		t.Errorf("Error adding gateway handler %s", err.Error())
		return
	}

	go func() {

		err := svc.Serve()
		if err != nil {
			t.Errorf("Error serving %s", err.Error())
		}

	}()

	time.Sleep(2 * time.Second)

	svc.Grpc.Server.Stop()

}
