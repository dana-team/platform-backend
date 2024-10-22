package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/danielgtaylor/huma/v2"
)

const (
	serviceAccountTag  = "ServiceAccounts"
	serviceAccountsKey = "serviceaccounts"
)

// AddGetServiceAccount adds the GetServiceAccount route to the OpenAPI scheme.
func AddGetServiceAccount(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-serviceaccount",
		Method:      http.MethodGet,
		Tags:        []string{serviceAccountTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName),
		Summary:     "Get a specific ServiceAccount in a namespace",
		Description: "Retrieves a specific ServiceAccount in a namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     serviceAccountName,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.ServiceAccountName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccount{})),
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

// AddCreateServiceAccount adds the CreateServiceAccount route to the OpenAPI scheme.
func AddCreateServiceAccount(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-serviceaccount",
		Method:      http.MethodPost,
		Tags:        []string{serviceAccountTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName),
		Summary:     "Create a ServiceAccount in a namespace",
		Description: "Creates a new ServiceAccount in a specific namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     serviceAccountName,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.ServiceAccountName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccount{})),
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

// AddDeleteServiceAccount adds the DeleteServiceAccount route to the OpenAPI scheme.
func AddDeleteServiceAccount(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-serviceaccount",
		Method:      http.MethodDelete,
		Tags:        []string{serviceAccountTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName),
		Summary:     "Deletes a ServiceAccount in a namespace",
		Description: "Deletes a specific ServiceAccount in a specific namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     serviceAccountName,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountRequestUri{}.ServiceAccountName)),
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

// AddGetServiceAccounts adds the GetServiceAccounts route to the OpenAPI scheme.
func AddGetServiceAccounts(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-serviceAccounts",
		Method:      http.MethodGet,
		Tags:        []string{serviceAccountTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s", namespacesKey, namespaceNameKey, serviceAccountsKey),
		Summary:     "Get all ServiceAccounts in a namespace",
		Description: "Retrieves all ServiceAccounts in a specific namespace",
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ServiceAccountOutput{})),
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
