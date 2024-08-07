package types

type KeyValue struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Metadata struct {
	Name              string `json:"name" binding:"required"`
	Namespace         string `json:"namespace"`
	CreationTimestamp string `json:"creationTimestamp"`
}

type ListMetadata struct {
	Count int `json:"count"`
}
