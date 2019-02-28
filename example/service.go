package main

import (
	// "fmt"

	"golang.org/x/net/context"

	pb "github.com/stevepartridge/service/example/protos"
)

type exampleService struct{}

func (s *exampleService) Info(c context.Context, req *pb.ServiceInfoRequest) (*pb.ServiceInfoResponse, error) {

	return &pb.ServiceInfoResponse{
		Name:    serviceName,
		BuiltAt: builtAt,
		Version: version,
		Build:   build,
		GitHash: githash,
	}, nil

}
