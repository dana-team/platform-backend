package operation

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"strconv"
)

const (
	terminalTag = "Terminal"
)

// AddStartTerminal adds the StartTerminal route to the OpenAPI scheme.
func AddStartTerminal(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "start-pod-terminal",
		Method:      http.MethodPost,
		Tags:        []string{terminalTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, podsKey, podNameKey, containersKey, containerNameKey, terminalKey),
		Summary:     "Start terminal for a pod in a namespace",
		Description: "Returns a sessionID for the terminal of a specific pod in a specific namespace",
		Parameters: []*huma.Param{
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalUri{}.ClusterName)),
				Example:  defaultExample,
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     podNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalUri{}.PodName)),
				Example:  defaultExample,
			},
			{
				Name:     containerNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalUri{}.ContainerName)),
				Example:  defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalBody{})),
				},
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.StartTerminalResponse{})),
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ErrorResponse{})),
					},
				},
			},
			strconv.Itoa(http.StatusInternalServerError): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ErrorResponse{})),
					},
				},
			},
		},
	}
	api.OpenAPI().AddOperation(operation)
}

// AddServeTerminal adds the StartTerminal route to the OpenAPI scheme.
func AddServeTerminal(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "serve-terminal",
		Method:      http.MethodGet,
		Tags:        []string{terminalTag},
		Path:        fmt.Sprintf("/ws/%s", terminalKey),
		Summary:     "Serves terminal for a pod in a namespace",
		Description: "Serves an interactive terminal for a pod using websockets",
		Parameters: []*huma.Param{
			{
				Name:     connectionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  upgradeHeaderKey,
			},
			{
				Name:     upgradeHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  webSocketValue,
			},
			{
				Name:     secWebSocketVersionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  secWebSocketVersionValue,
			},
		},
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Examples: map[string]*huma.Example{
						"Bind to SessionID": {
							Value: map[string]any{
								opKey:        "bind",
								dataKey:      "bash",
								sessionIDKey: "<id>",
								rowsKey:      10,
								colsKey:      10,
							},
						},
						"Run whoami inside the container": {
							Value: map[string]any{
								opKey:        "stdin",
								dataKey:      "whoami\r",
								sessionIDKey: "<id>",
								rowsKey:      10,
								colsKey:      10,
							},
						},
					},
				},
			},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusSwitchingProtocols): {
				Description: "Switching protocols",
			},
			strconv.Itoa(http.StatusBadRequest): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ErrorResponse{})),
					},
				},
			},
			strconv.Itoa(http.StatusInternalServerError): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ErrorResponse{})),
					},
				},
			},
		},
	}
	api.OpenAPI().AddOperation(operation)
}
