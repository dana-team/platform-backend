package types

type ConfigMapUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	ConfigMapName string `uri:"configMapName" binding:"required"`
}

type ConfigMap struct {
	Data []KeyValue `json:"data"`
}
