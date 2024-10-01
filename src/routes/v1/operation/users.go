package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetUsers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/users",
		Summary:     "[Users] Get all Users",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddCreateUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/v1/namespaces/{namespaceName}/users",
		Summary:     "[Users] Create a User",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		RequestBody: &huma.RequestBody{
			Description: "user content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.User{})),
				},
			},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddGetUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Get a User",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     userNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.UserName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddUpdateUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Update a User",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     userNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.UserName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		RequestBody: &huma.RequestBody{
			Description: "user content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.UpdateUserData{})),
				},
			},
		},
	}

	api.OpenAPI().AddOperation(operation)
}

func AddDeleteUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Delete a User",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     userNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UserIdentifier{}.UserName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}
