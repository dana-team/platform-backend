package types

type KeyValue struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Metadata struct {
	Name              string `json:"name" binding:"required"`
	Namespace         string `json:"namespace" binding:"required"`
	CreationTimestamp string `json:"creationTimestamp"`
}
