package types

type Namespace struct {
	Name string `json:"name" binding:"required"`
}

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
	Count      int         `json:"count"`
}

type NamespaceUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
}
