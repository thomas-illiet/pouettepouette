package utility

import (
	"context"
	"supervisor/api"

	"google.golang.org/grpc"
)

type UtilityService struct {
	api.UtilityServiceServer
}

// RegisterGRPC registers the gRPC info service.
func (us *UtilityService) RegisterGRPC(srv *grpc.Server) {
	api.RegisterUtilityServiceServer(srv, us)
}

func (us *UtilityService) Ping(ctx context.Context, request *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Message: "pong"}, nil
}
