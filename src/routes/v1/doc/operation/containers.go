package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const containerTag = "Containers"

// AddGetContainers adds the GetContainers route to the OpenAPI scheme.
func AddGetContainers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-containers",
		Method:      http.MethodGet,
		Tags:        []string{containerTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, podsKey, podNameKey, containersKey),
		Summary:     "Get containers of a pod in a namespace",
		Description: "Retrieves the names of the containers of a specific pod in a specific namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ContainerRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     podNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ContainerRequestUri{}.PodName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetContainersResponse{})),
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

// AddClusterGetContainers adds the ClusterGetContainers route to the OpenAPI scheme.
func AddClusterGetContainers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-containers-cluster",
		Method:      http.MethodGet,
		Tags:        []string{containerTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, podsKey, podNameKey, containersKey),
		Summary:     "Get containers of a pod in a namespace with cluster name",
		Description: "Retrieves the names of the containers of a specific pod in a specific namespace with cluster name",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ContainerRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     podNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ContainerRequestUri{}.PodName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetContainersResponse{})),
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
