package types

type ServiceAccountRequestUri struct {
	NamespaceName      string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	ServiceAccountName string `uri:"serviceAccountName" json:"serviceAccountName" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type CreateServiceAccountRequest ServiceAccountRequestUri

type ServiceAccountOutput struct {
	ServiceAccounts []string `json:"serviceAccounts"`
	ListMetadata
}

type ServiceAccount struct {
	Name string `json:"name" binding:"required"`
}
