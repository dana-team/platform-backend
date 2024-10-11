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
	tokenTag           = "Tokens"
	serviceAccountsKey = "serviceaccounts"
	tokensKey          = "token"
)

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
