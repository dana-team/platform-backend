package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const tokenTag = "Tokens"

func AddGetToken(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-token",
		Method:      http.MethodGet,
		Tags:        []string{tokenTag},
		Path:        "/v1/namespaces/{namespaceName}/serviceaccounts/{serviceAccountName}/token",
		Summary:     "[Token] Get token of a ServiceAccount",
		Description: "Get token of a ServiceAccount",
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
