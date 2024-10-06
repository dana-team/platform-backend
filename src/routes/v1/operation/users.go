package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const userTag = "Users"

func AddGetUsers(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Tags:        []string{userTag},
		Path:        "/v1/namespaces/{namespaceName}/users",
		Summary:     "[Users] Get all Users",
		Description: "Get all users in a namespace",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.UsersOutput{})),
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

func AddCreateUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Tags:        []string{userTag},
		Path:        "/v1/namespaces/{namespaceName}/users",
		Summary:     "[Users] Create a User",
		Description: "Create a new user in a namespace",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.User{})),
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

func AddGetUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Tags:        []string{userTag},
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Get a User",
		Description: "Get a user in a namespace",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.User{})),
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

func AddUpdateUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodGet,
		Tags:        []string{userTag},
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Update a User",
		Description: "Update a user",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.User{})),
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

func AddDeleteUser(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Tags:        []string{userTag},
		Path:        "/v1/namespaces/{namespaceName}/users/{userName}",
		Summary:     "[Users] Delete a User",
		Description: "Delete a user",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.DeleteUserResponse{})),
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
