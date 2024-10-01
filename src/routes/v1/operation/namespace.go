package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetNamespaces(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-namespaces",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces",
		Summary:     "[Namespace] Get all namespaces",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceUri{})),
				Example:  defaultExample,
			},
		},

		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddGetNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-namespace",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}",
		Summary:     "[Namespace] Get a namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceUri{})),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddCreateNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-namespace",
		Method:      http.MethodPost,
		Path:        "/v1/namespaces",
		Summary:     "[Namespace] Create a namespace",
		RequestBody: &huma.RequestBody{
			Description: "namespace name",
			Required:    true,
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Namespace{})),
				},
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}
	api.OpenAPI().AddOperation(operation)
}

func AddDeleteNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-namespace",
		Method:      http.MethodDelete,
		Path:        "/v1/namespaces/{namespaceName}",
		Summary:     "[Namespace] Delete a namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceUri{})),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}
