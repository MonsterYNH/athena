package services

import (
	"athena/api/v1/health"
	"context"
)

type HealthService struct {
	health.UnimplementedHealthServiceServer
}

func (service *HealthService) HealthCheck(ctx context.Context, request *health.HealthResponse) (*health.HealthResponse, error) {
	return &health.HealthResponse{
		Message: "OK",
	}, nil
}
