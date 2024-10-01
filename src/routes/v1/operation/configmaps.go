package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
)

func AddGetConfigMap(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-configmap",
		Method:      http.MethodGet,
		Path:        "/v1/namespaces/{namespaceName}/configmaps/{configMapName}",
		Summary:     "[ConfigMap] Get a ConfigMap",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ConfigMapUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     configMapNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.ConfigMapUri{}.ConfigMapName)),
				Example:  defaultExample,
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
	}

	api.OpenAPI().AddOperation(operation)
}
