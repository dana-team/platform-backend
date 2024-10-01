package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

const (
	secretNameKey = "secretName"
)

func AddCreateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-secret",
		Method:      http.MethodPost,
		Path:        "/v1/namespaces/{namespaceName}/secrets",
		Summary:     "[Secret] Create a secret",
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
	}

	api.OpenAPI().AddOperation(operation)
}

func AddGetSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-secret",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Get a secret",
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
	}

	api.OpenAPI().AddOperation(operation)
}

func AddUpdateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-secret",
		Method:      http.MethodPut,
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Update a secret",
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
	}

	api.OpenAPI().AddOperation(operation)
}

func AddDeleteSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-secret",
		Method:      http.MethodDelete,
		Path:        "/v1/namespaces/{namespaceName}/secrets/{secretName}",
		Summary:     "[Secret] Delete a secret",
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
	}

	api.OpenAPI().AddOperation(operation)
}
