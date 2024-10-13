package operation

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/danielgtaylor/huma/v2"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	secretNameKey = "secretName"
	secretsName   = "secrets"
	secretTag     = "Secrets"
)

// AddCreateSecret adds the CreateSecret route to the OpenAPI scheme.
func AddCreateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-secret",
		Method:      http.MethodPost,
		Tags:        []string{secretTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s", namespacesKey, namespaceNameKey, secretsName),
		Summary:     "Create a secret in a namespace",
		Description: "Creates a new secret in a specific namespace",
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
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Examples: map[string]*huma.Example{
						"Full scheme": {
							Value: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateSecretRequest{})),
						},
						"Opaque secret": {
							Value: mocks.PrepareCreateSecretRequestType(
								testutils.SecretName,
								strings.ToLower(string(corev1.SecretTypeOpaque)),
								"",
								"",
								[]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
							),
						},
					},
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

// AddGetSecret adds the GetSecret route to the OpenAPI scheme.
func AddGetSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-secret",
		Method:      http.MethodGet,
		Tags:        []string{secretTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, secretsName, secretNameKey),
		Summary:     "Get a secret in a namespace",
		Description: "Retrieves a specific secret in a specific namespace",
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

// AddGetSecrets adds the GetSecrets route to the OpenAPI scheme.
func AddGetSecrets(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-secrets",
		Method:      http.MethodGet,
		Tags:        []string{secretTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s", namespacesKey, namespaceNameKey, secretsName),
		Summary:     "Get all secrets in a namespace",
		Description: "Retrieves all the secrets in a specific namespace",
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

// AddUpdateSecret adds the UpdateSecret route to the OpenAPI scheme.
func AddUpdateSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-secret",
		Method:      http.MethodPut,
		Tags:        []string{secretTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, secretsName, secretNameKey),
		Summary:     "Update a secret in a namespace",
		Description: "Updates a specific secret in a specific namespace",
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

// AddDeleteSecret adds the DeleteSecret route to the OpenAPI scheme.
func AddDeleteSecret(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-secret",
		Method:      http.MethodDelete,
		Tags:        []string{secretTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, secretsName, secretNameKey),
		Summary:     "Delete a secret in a namespace",
		Description: "Deletes a specific secret in a specific namespace",
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
