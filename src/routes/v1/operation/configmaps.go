package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const configmapTag = "ConfigMaps"

func AddGetConfigMap(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-configmap",
		Method:      http.MethodGet,
		Tags:        []string{configmapTag},
		Path:        "/v1/namespaces/{namespaceName}/configmaps/{configMapName}",
		Summary:     "[ConfigMap] Get a ConfigMap",
		Description: "Get a ConfigMap by its name and namespace.",
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
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.ConfigMap{})),
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
