package types

type ServiceAccount struct {
	ServiceAccountName string `json:"serviceAccountName" binding:"required"`
}

type Token struct {
	Token string `json:"token"`
}
