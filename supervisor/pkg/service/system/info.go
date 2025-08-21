package system

import (
	"context"
	"supervisor/api"
)

func (is *SystemService) WorkspaceInfo(ctx context.Context, request *api.WorkspaceInfoRequest) (*api.WorkspaceInfoResponse, error) {
	resp := &api.WorkspaceInfoResponse{
		WorkspaceId:      is.Cfg.WorkspaceID,
		CheckoutLocation: is.Cfg.WorkspaceLocation + "/devel",
		UserHome:         is.Cfg.WorkspaceLocation,
		OwnerId:          is.Cfg.OwnerId,
		ClusterHost:      is.Cfg.WorkspaceClusterHost,
		IdeAlias:         is.Cfg.Editor.Name,
		IdePort:          3000,
	}

	return resp, nil
}
