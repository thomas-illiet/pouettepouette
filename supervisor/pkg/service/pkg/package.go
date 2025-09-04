
package

import (
	"supervisor/api"

	"google.golang.org/grpc"
)
package

import (
	"context"
	"supervisor/api"

	"google.golang.org/grpc"
)

type PackageService struct {
	api.PackageServiceServer
}

// RegisterGRPC registers the gRPC info service.
func (us *PackageService) RegisterGRPC(srv *grpc.Server) {
	api.RegisterPackageServiceServer(srv, us)
}

func (us *PackageService) List(ctx context.Context, request *api.ListPackageRequest) (*api.ListPackageResponse, error) {
	result := &api.ListPackageResponse{}
	result.Packages = append(result.Packages, &api.GetPackageResponse{
		Id:          1,
		Name:        "test",
		Description: "test",
		Version:     "1.0.0",
	})
	result.Packages = append(result.Packages, &api.GetPackageResponse{
		Id:          2,
		Name:        "test2",
		Description: "test2",
		Version:     ".0.0",
	})

	return result, nil
}

func (us *PackageService) Get(ctx context.Context, request *api.GetPackageRequest) (*api.GetPackageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *PackageService) GetContent(ctx context.Context, request *api.GetPackageContentRequest) (*api.GetPackageContentResponse, error) {
	//TODO implement me
	panic("implement me")
}
