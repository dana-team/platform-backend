package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/danielgtaylor/huma/v2"
)

const configMapTag = "ConfigMaps"

// AddGetConfigMap adds the GetConfigMap route to the OpenAPI scheme.
func AddGetConfigMap(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-configmap",
		Method:      http.MethodGet,
		Tags:        []string{configMapTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, configMapsKey, configMapNameKey),
		Summary:     "Get a ConfigMap in a namespace",
		Description: "Retrieves a specific ConfigMap in a specific namespace",
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
