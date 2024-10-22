package operation

import (
	"fmt"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"strconv"
)

const (
	logsTag = "Logs"
)

// AddGetPodLogs adds the GetPodLogs route to the OpenAPI scheme.
func AddGetPodLogs(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pod-logs",
		Method:      http.MethodGet,
		Tags:        []string{logsTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, podsKey, podNameKey, logsKey),
		Summary:     "Get logs of a pod in a namespace",
		Description: "Retrieves the logs of a specific pod in a specific namespace using websockets",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     podNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:    containerNameKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example: defaultExample,
			},
			{
				Name:    previousKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeBoolean)),
				Example: false,
			},
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
				Name:     secWebSocketProtocolHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultToken,
			},
			{
				Name:     secWebSocketKeyHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultBase64,
			},
			{
				Name:     secWebSocketVersionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  secWebSocketVersionValue,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
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

// AddClusterGetPodLogs adds the ClusterGetPodLogs route to the OpenAPI scheme.
func AddClusterGetPodLogs(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pod-logs-cluster",
		Method:      http.MethodGet,
		Tags:        []string{logsTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, podsKey, podNameKey, logsKey),
		Summary:     "Get logs of a pod in a namespace with cluster name",
		Description: "Retrieves the logs of a specific pod in a specific namespace using websockets with cluster name",
		Parameters: []*huma.Param{
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     podNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:    containerNameKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example: defaultExample,
			},
			{
				Name:    previousKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeBoolean)),
				Example: false,
			},
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
				Name:     secWebSocketProtocolHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultToken,
			},
			{
				Name:     secWebSocketKeyHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultBase64,
			},
			{
				Name:     secWebSocketVersionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  secWebSocketVersionValue,
			},
		},

		Security: []map[string][]string{
			{bearerKey: {}},
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

// AddGetCappLogs adds the GetCappLogs route to the OpenAPI scheme.
func AddGetCappLogs(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-logs",
		Method:      http.MethodGet,
		Tags:        []string{logsTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, cappsKey, cappNameKey, logsKey),
		Summary:     "Get logs of a Capp in a namespace",
		Description: "Retrieves the logs of a specific Capp in a specific namespace using websockets",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:    podNameKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example: defaultExample,
			},
			{
				Name:    containerNameKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example: defaultExample,
			},
			{
				Name:    previousKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeBoolean)),
				Example: false,
			},
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
				Name:     secWebSocketProtocolHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultToken,
			},
			{
				Name:     secWebSocketKeyHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultBase64,
			},
			{
				Name:     secWebSocketVersionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  secWebSocketVersionValue,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
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

// AddClusterGetCappLogs adds the ClusterGetCappLogs route to the OpenAPI scheme.
func AddClusterGetCappLogs(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-logs",
		Method:      http.MethodGet,
		Tags:        []string{logsTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, cappsKey, cappNameKey, logsKey),
		Summary:     "Get logs of a Capp in a namespace with cluster name",
		Description: "Retrieves the logs of a specific Capp in a specific namespace using websockets using cluster name",
		Parameters: []*huma.Param{
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:    containerNameKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example: defaultExample,
			},
			{
				Name:    previousKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeBoolean)),
				Example: false,
			},
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
				Name:     secWebSocketProtocolHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultToken,
			},
			{
				Name:     secWebSocketKeyHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultBase64,
			},
			{
				Name:     secWebSocketVersionHeaderKey,
				In:       headerKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  secWebSocketVersionValue,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
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
