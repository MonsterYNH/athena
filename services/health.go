package services

import (
	"athena/api/models"
	"athena/api/v1/health"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthService struct {
	health.UnimplementedHealthServiceServer
}

func (service *HealthService) HealthCheck(ctx context.Context, request *health.HealthRequest) (*models.Response, error) {
	return nil, status.Errorf(codes.AlreadyExists, "asdasdasd")
}
