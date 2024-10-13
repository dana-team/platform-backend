package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const podTag = "Pods"

// AddGetPods adds the GetPods route to the OpenAPI scheme.
func AddGetPods(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pods",
		Method:      http.MethodGet,
		Tags:        []string{podTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, cappsKey, cappNameKey, podsKey),
		Summary:     "Get pods of a Capp in a namespace",
		Description: "Retrieves pods of a specific in a Capp in a namespace",
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

// AddClusterGetPods adds the ClusterGetPods route to the OpenAPI scheme.
func AddClusterGetPods(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-pods-cluster",
		Method:      http.MethodGet,
		Tags:        []string{podTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, cappsKey, cappNameKey, podsKey),
		Summary:     "Get pods of a Capp in a namespace with cluster name",
		Description: "Retrieves pods of a specific in a Capp in a namespace with cluster name",
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
