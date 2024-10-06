package operation

import (
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

func AddHealthz(api huma.API) {
	operation := &huma.Operation{
		OperationID: "healthz",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "[Healthz] Health Check",
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Example: map[string]string{"Message": "OK"},
					},
				},
			},
		},
	}
	api.OpenAPI().AddOperation(operation)
}
