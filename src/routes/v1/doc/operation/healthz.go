package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"net/http"
	"reflect"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

func AddHealthz(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "healthz",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Health Check",
		Description: "Returns OK if API is up and running",
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.MessageResponse{})),
					},
				},
			},
		},
	}
	api.OpenAPI().AddOperation(operation)
}
