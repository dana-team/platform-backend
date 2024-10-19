package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const (
	tokenTag  = "Tokens"
	tokensKey = "token"
)

// TODO: figure out if this function is neccessary. It is sort of confusing since this is relating to dockercfgTokens, not auth tokens.
// AddGetToken adds the GetToken route to the OpenAPI scheme.
func AddGetToken(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-token",
		Method:      http.MethodGet,
		Tags:        []string{tokenTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName, tokensKey),
		Summary:     "Get token of a ServiceAccount in a namespace",
		Description: "Retrieves the token of a specific ServiceAccount in a namespace",
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.TokenResponse{})),
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

// AddCreateToken adds the CreateToken route to the OpenAPI scheme.
func AddCreateToken(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-token",
		Method:      http.MethodPost,
		Tags:        []string{tokenTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName, tokensKey),
		Summary:     "Creates an auth token for a service account",
		Description: "Creates an auth token for a service account in a namespace",
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
			{
				Name:     expirationSecondsKey,
				In:       queryKey,
				Required: false,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf("")),
				Example:  defaultExampleSeconds,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.TokenRequestResponse{})),
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

// AddRevokeToken adds the RevokeToken route to the OpenAPI scheme.
func AddRevokeToken(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "revoke-token",
		Method:      http.MethodDelete,
		Tags:        []string{tokenTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, serviceAccountsKey, serviceAccountName, tokensKey),
		Summary:     "Revokes the tokens for a ServiceAccount",
		Description: "Revokes all tokens for a specific ServiceAccount in a namespace",
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
