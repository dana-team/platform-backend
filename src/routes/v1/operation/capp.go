package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/danielgtaylor/huma/v2"
)

const (
	cappTag = "Capps"
)

func AddGetCapps(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capps",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps",
		Summary:     "[Capps] Get all Capps",
		Description: "Get all Capps in a namespace",
		Parameters: []*huma.Param{
			{
				Name:    paginationPageKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Page)),
				Example: 1,
			},
			{
				Name:    paginationLimitKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.PaginationParams{}.Limit)),
				Example: 1,
			},
			{
				Name:    labelSelectorKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.GetCappQuery{}.LabelSelector)),
				Example: "app=example",
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappNamespaceUri{}.NamespaceName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappList{})),
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

func AddCreateCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-capp",
		Method:      http.MethodPost,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps",
		Summary:     "[Capps] Create a Capp",
		Description: "Create a new capp",
		Parameters: []*huma.Param{
			{
				Name:    namespaceNameKey,
				In:      pathKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.CappNamespaceUri{}.NamespaceName)),
				Example: defaultExample,
			},
			{
				Name:    utils.PlacementEnvironmentKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.CreateCappQuery{}.Environment)),
				Example: defaultExample,
			},
			{
				Name:    utils.PlacementRegionKey,
				In:      queryKey,
				Schema:  huma.SchemaFromType(registry, reflect.TypeOf(types.CreateCappQuery{}.Region)),
				Example: defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Description: "capp content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateCapp{})),
				},
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Capp{})),
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

func AddGetCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}",
		Summary:     "[Capps] Get a Capp",
		Description: "Get a specific capp in a namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Capp{})),
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

func AddUpdateCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-capp",
		Method:      http.MethodPut,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}",
		Summary:     "[Capps] Update a Capp",
		Description: "Update a capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Description: "capp content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.UpdateCapp{})),
				},
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.Capp{})),
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

func AddEditCappState(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "edit-capp-state",
		Method:      http.MethodPut,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/state",
		Summary:     "[Capps] Edit a Capp state",
		Description: "Enables or disables a capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
				Example:  defaultExample,
			},
		},
		RequestBody: &huma.RequestBody{
			Description: "capp content",
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappState{})),
				},
			},
		},
		Security: []map[string][]string{
			{bearerKey: {}},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Content: map[string]*huma.MediaType{
					applicationJSONKey: {
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappStateResponse{})),
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

func AddGetCappState(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-state",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/state",
		Summary:     "[Capps] Get a Capp state",
		Description: "Get the state of a capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetCappStateResponse{})),
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

func AddDeleteCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-capp",
		Method:      http.MethodDelete,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}",
		Summary:     "[Capps] Delete a Capp",
		Description: "Delete a capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.MessageResponse{})),
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

func AddGetCappDNS(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-dns",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/dns",
		Summary:     "[Capps] Get a Capp DNS",
		Description: "Get the DNS record of a capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.GetDNSResponse{})),
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
