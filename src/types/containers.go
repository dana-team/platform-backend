package types

type GetContainersResponse struct {
	Count      int         `json:"count"`
	Containers []Container `json:"containers"`
}

type Container struct {
	ContainerName string `json:"containerName"`
}

type ContainerRequestUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	PodName       string `uri:"podName" binding:"required"`
}
