package types

type GetPodsResponse struct {
	Count int   `json:"count"`
	Pods  []Pod `json:"pods"`
}

type Pod struct {
	PodName string `json:"podName"`
}

type PodRequestUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName" binding:"required"`
}
