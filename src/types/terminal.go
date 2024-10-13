package types

type StartTerminalUri struct {
	ClusterName   string `uri:"clusterName" json:"clusterName" binding:"required"`
	NamespaceName string `uri:"namespaceName" json:"namespaceName" binding:"required"`
	PodName       string `uri:"podName" json:"podName" binding:"required"`
	ContainerName string `uri:"containerName" json:"containerName" binding:"required"`
}

type StartTerminalBody struct {
	Shell string `json:"shell" binding:"required,oneof=bash sh powershell cmd"`
}

// StartTerminalResponse is sent by HandleStartTerminal. The ID is a random session id that binds the original REST request and the SockJS connection.
// Any client api in possession of this ID can hijack the terminal_utils session.
type StartTerminalResponse struct {
	ID string `json:"id"`
}
