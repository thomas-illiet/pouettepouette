package supervisor

import (
	"common/util"
	"context"
	"fmt"
	"supervisor/api"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SupervisorClient wraps a gRPC connection to the Supervisor service, exposing typed service clients.
type SupervisorClient struct {
	conn      *grpc.ClientConn
	closeOnce sync.Once

	// Service clients
	Package api.PackageServiceClient
	System  api.SystemServiceClient
	Utility api.UtilityServiceClient
}

// New creates a new SupervisorClient using the supervisor address.
func New(ctx context.Context) (*SupervisorClient, error) {
	address := util.GetSupervisorAddress()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for %s: %w", address, err)
	}

	return &SupervisorClient{
		conn:    conn,
		Package: api.NewPackageServiceClient(conn),
		System:  api.NewSystemServiceClient(conn),
		Utility: api.NewUtilityServiceClient(conn),
	}, nil
}

// Close terminates the underlying gRPC connection to the Supervisor service.
func (c *SupervisorClient) Close() {
	c.closeOnce.Do(func() {
		_ = c.conn.Close()
	})
}
