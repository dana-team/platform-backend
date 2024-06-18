package types

type CreateSecretRequest struct {
	Type       string     `json:"type" binding:"required" enum:"tls,opaque"`
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

type SecretUriRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	SecretName    string `uri:"secretName" binding:"required"`
}

type SecretNamespaceUriRequest struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
}

type GetSecretsResponse struct {
	Count   int      `json:"count"`
	Secrets []Secret `json:"namespaces"`
}

type GetSecretResponse struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	SecretName string     `json:"secretName"`
	Data       []KeyValue `json:"data"`
}

type Secret struct {
	Type          string `json:"type"`
	SecretName    string `json:"secretName"`
	NamespaceName string `json:"namespaceName"`
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
