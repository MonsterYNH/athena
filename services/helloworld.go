package services

import (
	"athena/api/v1/helloworld"
	"context"
	"fmt"
)

type HelloworldService struct {
	helloworld.UnimplementedHellowroldServiceServer
}

func (service *HelloworldService) SyaHello(ctx context.Context, request *helloworld.HelloworldRequest) (*helloworld.HelloworldResponse, error) {
	return &helloworld.HelloworldResponse{
		Message: fmt.Sprintf("hello %s", request.Name),
	}, nil
}
