package types

type ServiceAccountRequestUri struct {
	NamespaceName      string `uri:"namespaceName" binding:"required"`
	ServiceAccountName string `uri:"serviceAccountName" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
