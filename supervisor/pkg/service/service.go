package service

import "google.golang.org/grpc"

// RegisterableService can register a service.
type RegisterableService interface{}

// RegisterableGRPCService can register gRPC services.
type RegisterableGRPCService interface {
	// RegisterGRPC registers a gRPC service
	RegisterGRPC(*grpc.Server)
}
