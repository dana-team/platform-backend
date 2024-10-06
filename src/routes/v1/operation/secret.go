package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const (
	secretNameKey = "secretName"
	secretTag     = "Secrets"
)

func AddCreateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-secret",
		Method:      http.MethodPost,
		Tags:        []string{secretTag},
		Path:        "/v1/namespaces/{namespaceName}/secrets",
		Summary:     "[Secret] Create a secret",
		Description: "Create a secret in a namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretNamespaceUriRequest{})),
				Example:  defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Description: "secret content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateSecretRequest{})),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateSecretResponse{})),
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

func AddGetSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-secret",
		Method:      http.MethodGet,
		Tags:        []string{secretTag},
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Get a secret",
		Description: "Get a secret by namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     secretNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.SecretName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetSecretResponse{})),
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

func AddGetSecrets(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-secrets",
		Method:      http.MethodGet,
		Tags:        []string{secretTag},
		Path:        "/v1/namespaces/{namespaceName}/secrets",
		Summary:     "[Secret] Get all secrets in a namespace",
		Description: "Get all secrets in a namespace",
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
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.NamespaceName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetSecretsResponse{})),
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

func AddUpdateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-secret",
		Method:      http.MethodPut,
		Tags:        []string{secretTag},
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Update a secret",
		Description: "Update a secret",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     secretNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.SecretName)),
				Example:  defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Description: "secret content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.UpdateSecretRequest{})),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.UpdateSecretResponse{})),
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

func AddDeleteSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-secret",
		Method:      http.MethodDelete,
		Tags:        []string{secretTag},
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Delete a secret",
		Description: "Delete a secret",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     secretNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.SecretUriRequest{}.SecretName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.DeleteSecretResponse{})),
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
