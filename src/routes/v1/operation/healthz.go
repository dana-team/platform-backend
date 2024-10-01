package operation

import (
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func AddHealthz(api huma.API) {
	operation := &huma.Operation{
		OperationID: "healthz",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "[Healthz] Health Check",
	}
	api.OpenAPI().AddOperation(operation)
}
