package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/danielgtaylor/huma/v2"
)

const namespaceTag = "Namespaces"

// AddGetNamespaces adds the GetNamespaces route to the OpenAPI scheme.
func AddGetNamespaces(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-namespaces",
		Method:      http.MethodGet,
		Tags:        []string{namespaceTag},
		Path:        fmt.Sprintf("/v1/%s", namespacesKey),
		Summary:     "Get all namespaces",
		Description: "Retrieves all namespaces",
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
		},

		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.NamespaceList{})),
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

// AddGetNamespace adds the GetNamespace route to the OpenAPI scheme.
func AddGetNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-namespace",
		Method:      http.MethodGet,
		Tags:        []string{namespaceTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}", namespacesKey, namespaceNameKey),
		Summary:     "Get a namespace",
		Description: "Retrieves a specific namespace",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Namespace{})),
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

// AddCreateNamespace adds the CreateNamespace route to the OpenAPI scheme.
func AddCreateNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-namespace",
		Method:      http.MethodPost,
		Tags:        []string{namespaceTag},
		Path:        fmt.Sprintf("/v1/%s", namespacesKey),
		Summary:     "Create a namespace",
		Description: "Creates a new namespace",
		RequestBody: &huma.RequestBody{
			Required: true,
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Namespace{})),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Namespace{})),
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

// AddDeleteNamespace adds the DeleteNamespace route to the OpenAPI scheme.
func AddDeleteNamespace(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-namespace",
		Method:      http.MethodDelete,
		Tags:        []string{namespaceTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}", namespacesKey, namespaceNameKey),
		Summary:     "Delete a namespace",
		Description: "Deletes a specific namespace",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.MessageResponse{})),
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
