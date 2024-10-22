package types

import "time"

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
	Name  string `json:"name" binding:"required"`
	Token string `json:"token,omitempty"`
}

type TokenRequestResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires,omitempty"`
}

type CreateTokenQuery struct {
	ExpirationSeconds string `form:"expirationSeconds" json:"expirationSeconds"`
}
