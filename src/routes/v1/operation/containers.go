package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetContainers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-containers",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/pods/{podName}/containers",
		Summary:     "[Containers] Get containers of a pod",
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
	}

	api.OpenAPI().AddOperation(operation)
}

func AddClusterGetContainers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-containers-cluster",
		Method:      http.MethodGet,
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/pods/{podName}/containers",
		Summary:     "[Containers] Get containers of a Capp with cluster",
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
	}

	api.OpenAPI().AddOperation(operation)
}
