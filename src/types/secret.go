package types

type CreateSecretRequest struct {
	Type string     `json:"type" binding:"required"`
	Name string     `json:"name" binding:"required"`
	Cert string     `json:"cert"`
	Key  string     `json:"key"`
	Data []KeyValue `json:"data"`
}

type CreateSecretResponse struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type KeyValue struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type GetSecretsResponse struct {
	Count   int      `json:"count"`
	Secrets []Secret `json:"namespaces"`
}

type Secret struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type GetSecretResponse struct {
	Id   string     `json:"id"`
	Type string     `json:"type"`
	Name string     `json:"name"`
	Data []KeyValue `json:"data"`
}

type PatchSecretRequest struct {
	Data []KeyValue `json:"data"`
}

type PatchSecretResponse struct {
	Id        string     `json:"id"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Data      []KeyValue `json:"data"`
}

type DeleteSecretResponse struct {
	Message string `json:"message"`
}
