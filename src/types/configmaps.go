package types

type ConfigMapUri struct {
	NamespaceName string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	ConfigMapName string `uri:"configMapName" json:"configMapName" binding:"required"`
}

type ConfigMap struct {
	Data []KeyValue `json:"data"`
}
