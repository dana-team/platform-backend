package types

type GetContainersResponse struct {
	Containers []Container `json:"containers"`
	ListMetadata
}

type Container struct {
	ContainerName string `json:"containerName"`
}

type ContainerRequestUri struct {
	NamespaceName string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	PodName       string `uri:"podName" json:"podName" binding:"required"`
}
