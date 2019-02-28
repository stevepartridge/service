package service

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/stevepartridge/service/middleware"
)

func Test_Unit_NewGrpcServer_Success(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Service should not result in error %s", err.Error())
	}

	if svc == nil {
		t.Error("service should not be nil")
	}

	server := svc.GrpcServer()
	if server == nil {
		t.Error("Grpc Server should not be nil")
	}

}

func Test_Unit_NewGrpcServerAddUnaryInterceptor_Success(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Service should not result in error %s", err.Error())
	}

	if svc == nil {
		t.Error("service should not be nil")
	}

	svc.Grpc.AddUnaryInterceptors(middleware.RequestInterceptor())

	if len(svc.Grpc.UnaryInterceptors) != 1 {
		t.Errorf("Expected 1 but saw %d", len(svc.Grpc.UnaryInterceptors))
	}

}

func Test_Unit_NewGrpcServerAddOption_Success(t *testing.T) {
	svc, err := New(testPort1)
	if err != nil {
		t.Errorf("Service should not result in error %s", err.Error())
	}

	if svc == nil {
		t.Error("service should not be nil")
	}

	svc.Grpc.AddOptions(grpc.MaxRecvMsgSize(1024))

	if len(svc.Grpc.ServerOptions) != 1 {
		t.Errorf("Expected 1 but saw %d", len(svc.Grpc.ServerOptions))
	}

}

// func Test_Unit_ServiceNewEnableHandler_Success(t *testing.T) {
// 	svc, err := New(testPort1)
// 	if err != nil {
// 		t.Errorf("Service should not result in error %s", err.Error())
// 	}

// 	if svc == nil {
// 		t.Error("service should not be nil")
// 	}

// 	err = svc.EnableGatewayHandler(dummyGatewayHandler)
// 	if err != nil {
// 		t.Errorf("Error should be nil not %s", err.Error())
// 	}
// }

// func dummyGatewayHandler(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error {
// 	return nil
// }
