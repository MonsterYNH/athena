package services

import (
	"athena/api/models"
	"athena/api/v1/helloworld"
	"context"
	"fmt"
)

type HelloworldService struct {
	helloworld.UnimplementedHellowroldServiceServer
}

func (service *HelloworldService) SyaHello(ctx context.Context, request *helloworld.HelloworldRequest) (*models.Response, error) {
	return &models.Response{
		Message: fmt.Sprintf("hello %s", request.Name),
	}, nil
}
