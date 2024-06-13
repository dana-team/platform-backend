package types

type CreateSecretRequest struct {
	Type       string     `json:"type" binding:"required"`
	SecretName string     `json:"secretName" binding:"required"`
	Cert       string     `json:"cert"`
	Key        string     `json:"key"`
	Data       []KeyValue `json:"data"`
}

type CreateSecretResponse struct {
	Type          string `json:"type"`
	SecretName    string `json:"secretName"`
	NamespaceName string `json:"namespaceName"`
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
	Type          string `json:"type"`
	SecretName    string `json:"secretName"`
	NamespaceName string `json:"namespaceName"`
}

type GetSecretResponse struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	SecretName string     `json:"secretName"`
	Data       []KeyValue `json:"data"`
}

type PatchSecretRequest struct {
	Data []KeyValue `json:"data"`
}

type PatchSecretResponse struct {
	Id            string     `json:"id"`
	Type          string     `json:"type"`
	SecretName    string     `json:"secretName"`
	NamespaceName string     `json:"namespaceName"`
	Data          []KeyValue `json:"data"`
}

type DeleteSecretResponse struct {
	Message string `json:"message"`
}
