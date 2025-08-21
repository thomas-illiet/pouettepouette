package system

import (
	"supervisor/api"
	"supervisor/pkg/config"

	"google.golang.org/grpc"
)

type SystemService struct {
	Cfg *config.Config
	api.SystemServiceServer
}

// RegisterGRPC registers the gRPC info service.
func (is *SystemService) RegisterGRPC(srv *grpc.Server) {
	api.RegisterSystemServiceServer(srv, is)
}
