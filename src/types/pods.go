package types

type GetPodsResponse struct {
	Pods []Pod `json:"pods"`
	ListMetadata
}

type Pod struct {
	PodName string `json:"podName"`
}

type PodRequestUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName" binding:"required"`
}
