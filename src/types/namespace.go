package types

type Namespace struct {
	Name string `json:"name"  binding:"required"`
}

type OutputNamespaces struct {
	Namespaces []Namespace `json:"namespaces"`
	Count      int         `json:"count"`
}
