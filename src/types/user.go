package types

type User struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"  enum:"admin,viewer,contributor"`
}

type UserIdentifier struct {
	UserName      string `json:"userName" uri:"userName" binding:"required"`
	NamespaceName string `json:"namespaceName" binding:"required" uri:"namespaceName"`
}

type UserInput struct {
	Namespace string `json:"namespace" binding:"required"`
	User
}

type PatchUserData struct {
	Role string `json:"role" binding:"required" validate:"one0f=admin viewer contributor"`
}

type UsersOutput struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}
