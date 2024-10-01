package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetToken(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-token",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/serviceaccounts/{serviceAccountName}/token",
		Summary:     "[Token] Get token of a ServiceAccount",
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
	}

	api.OpenAPI().AddOperation(operation)
}
