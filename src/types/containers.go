package types

type GetContainersResponse struct {
	Containers []Container `json:"containers"`
	ListMetadata
}

type Container struct {
	ContainerName string `json:"containerName"`
}

type ContainerRequestUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	PodName       string `uri:"podName" binding:"required"`
}
