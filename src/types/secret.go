package types

type CreateSecretRequest struct {
	Type          string     `json:"type" binding:"required"`
	NamespaceName string     `json:"namespaceName" binding:"required"`
	SecretName    string     `json:"secretName" binding:"required"`
	Cert          string     `json:"cert"`
	Key           string     `json:"key"`
	Data          []KeyValue `json:"data"`
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

type GetSecretsRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
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

type GetSecretRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	SecretName    string `uri:"secretName" binding:"required"`
}

type GetSecretResponse struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	SecretName string     `json:"secretName"`
	Data       []KeyValue `json:"data"`
}

type PatchSecretUriRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	SecretName    string `uri:"secretName" binding:"required"`
}

type PatchSecretJsonRequest struct {
	Data []KeyValue `json:"data"`
}

type PatchSecretRequest struct {
	NamespaceName string     `json:"namespaceName" binding:"required"`
	SecretName    string     `json:"secretName" binding:"required"`
	Data          []KeyValue `json:"data"`
}

type PatchSecretResponse struct {
	Id            string     `json:"id"`
	Type          string     `json:"type"`
	SecretName    string     `json:"secretName"`
	NamespaceName string     `json:"namespaceName"`
	Data          []KeyValue `json:"data"`
}

type DeleteSecretRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	SecretName    string `uri:"secretName" binding:"required"`
}

type DeleteSecretResponse struct {
	Message string `json:"message"`
}
