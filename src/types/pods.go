package types

type GetPodsResponse struct {
	Pods []Pod `json:"pods"`
	ListMetadata
}

type Pod struct {
	PodName string `json:"podName"`
}

type PodRequestUri struct {
	NamespaceName string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName" json:"cappName" binding:"required"`
}
