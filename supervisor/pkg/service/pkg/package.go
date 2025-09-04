package pkg

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
		Status:      api.PackageStatus_PACKAGE_STATUS_INSTALLED,
		Version:     "1.0.0",
	})
	result.Packages = append(result.Packages, &api.GetPackageResponse{
		Id:          2,
		Name:        "test2",
		Description: "test2",
		Status:      api.PackageStatus_PACKAGE_STATUS_NOT_INSTALLED,
		Version:     "2.0.0",
	})

	return result, nil
}

func (us *PackageService) Get(ctx context.Context, request *api.GetPackageRequest) (*api.GetPackageResponse, error) {
	return &api.GetPackageResponse{
		Id:          1,
		Name:        "test",
		Description: "test",
		Status:      api.PackageStatus_PACKAGE_STATUS_UPDATABLE,
		Version:     "1.0.0",
	}, nil
}

func (us *PackageService) Install(ctx context.Context, request *api.InstallPackageRequest) (*api.InstallPackageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *PackageService) Remove(ctx context.Context, request *api.RemovePackageRequest) (*api.RemovePackageResponse, error) {
	//TODO implement me
	panic("implement me")
}
