package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const podTag = "Pods"

func AddGetPods(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pods",
		Method:      http.MethodGet,
		Tags:        []string{podTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/pods",
		Summary:     "[Pods] Get pods of a Capp",
		Description: "Get pods of a Capp",
		Parameters: []*huma.Param{
			{
				Name:    paginationPageKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Page)),
				Example: 1,
			},
			{
				Name:    paginationLimitKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Limit)),
				Example: 1,
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetPodsResponse{})),
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

// TODO: Research websocket docs
func AddGetPodLogs(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pod-logs",
		Method:      http.MethodGet,
		Tags:        []string{podTag},
		Path:        "/v1/namespaces/{namespaceName}/pod/{podName}/logs",
		Summary:     "[Pods] Get logs of a pod",
		Description: "Get logs of a pod",
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
				Name:     "Connection",
				In:       "header",
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  "upgrade",
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},

		Middlewares: huma.Middlewares{
			func(ctx huma.Context, next func(huma.Context)) {
			},
		},
		Responses: map[string]*huma.Response{
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

func AddClusterGetPods(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pods-cluster",
		Method:      http.MethodGet,
		Tags:        []string{podTag},
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/capps/{cappName}/pods",
		Summary:     "[Pods] Get pods of a Capp with cluster",
		Description: "Get pods of a Capp by cluster",
		Parameters: []*huma.Param{
			{
				Name:    paginationPageKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Page)),
				Example: 1,
			},
			{
				Name:    paginationLimitKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Limit)),
				Example: 1,
			},
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.PodRequestUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetPodsResponse{})),
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
