package log

import log "github.com/sirupsen/logrus"

const (
	UserIDField         = "userId"
	ContainerIDField    = "containerId"
	ClusterIDField      = "clusterId"
	ServiceContextField = "serviceContext"
)

func UserID(userID string) log.Fields {
	return String(UserIDField, userID)
}

func ContainerID(containerID string) log.Fields {
	return String(ContainerIDField, containerID)
}

func ClusterID(clusterID string) log.Fields {
	return String(ClusterIDField, clusterID)
}

type serviceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func ServiceContext(service, version string) log.Fields {
	return log.Fields{
		ServiceContextField: serviceContext{service, version},
	}
}

func String(key, value string) log.Fields {
	return log.Fields{key: value}
}
