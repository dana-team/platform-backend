package operation

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"reflect"
	"strconv"
)

func AddLogin(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/v1/login",
		Summary:     "Login",
		Description: "Returns token from username and password",
		Security: []map[string][]string{
			{basicAuthKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.LoginOutput{})),
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
			strconv.Itoa(http.StatusUnauthorized): {
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
