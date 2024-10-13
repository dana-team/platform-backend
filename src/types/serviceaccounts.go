package types

type ServiceAccountRequestUri struct {
	NamespaceName      string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	ServiceAccountName string `uri:"serviceAccountName" json:"serviceAccountName" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
