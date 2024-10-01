package types

type StartTerminalUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	PodName       string `uri:"podName" binding:"required"`
	ContainerName string `uri:"containerName" binding:"required"`
}

type StartTerminalBody struct {
	Shell string `json:"shell" enum:"bash, sh, powershell, cmd"`
}

// StartTerminalResponse is sent by HandleStartTerminal. The ID is a random session id that binds the original REST request and the SockJS connection.
// Any client api in possession of this ID can hijack the terminal_utils session.
type StartTerminalResponse struct {
	ID string `json:"id"`
}
