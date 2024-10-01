package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revisions",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/capprevisions",
		Summary:     "[CappRevisions] Get all CappsRevisions of a Capp",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.ClusterName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revision",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/capprevisions/{cappRevisionName}",
		Summary:     "[CappRevisions] Get a CappsRevision of a Capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.ClusterName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddClusterGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevisions-cluster",
		Method:      http.MethodGet,
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/capprevisions",
		Summary:     "[CappRevisions] Get CappRevisions with cluster",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddClusterGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevision-cluster",
		Method:      http.MethodGet,
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/capprevisions/{cappRevisionName}",
		Summary:     "[CappRevisions] Get a CappRevision with cluster",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappRevisionNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.CappRevisionName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}
