package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/danielgtaylor/huma/v2"
)

const cappTag = "Capps"

// AddGetCapps adds the GetCapps route to the OpenAPI scheme.
func AddGetCapps(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capps",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s", namespacesKey, namespaceNameKey, cappsKey),
		Summary:     "Get all Capps in a namespace",
		Description: "Retrieves a summary of all the Capps in a specific namespace",
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

// AddCreateCapp adds the CreateCapp route to the OpenAPI scheme.
func AddCreateCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "create-capp",
		Method:      http.MethodPost,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s", namespacesKey, namespaceNameKey, cappsKey),
		Summary:     "Create a Capp in a namespace",
		Description: "Creates a new Capp in a specific namespace",
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
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateCapp{})),
					Examples: map[string]*huma.Example{
						"Full scheme": {
							Value: huma.SchemaFromType(registry, reflect.TypeOf(types.CreateCapp{})),
						},
						"Capp with site": {
							Value: mocks.PrepareCreateCappType("test-name", "test-site", nil, nil),
						},
						"Capp without site": {
							Value: mocks.PrepareCreateCappType("test-name", "", nil, nil),
						},
					},
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

// AddGetCapp adds the GetCapp route to the OpenAPI scheme.
func AddGetCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Get a Capp in a namespace",
		Description: "Retrieves a specific Capp in a specific namespace",
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

// AddUpdateCapp adds the UpdateCapp route to the OpenAPI scheme.
func AddUpdateCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "update-capp",
		Method:      http.MethodPut,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Update a Capp in a namespace",
		Description: "Updates a specific Capp in a specific namespace",
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
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Examples: map[string]*huma.Example{},
					Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.UpdateCapp{})),
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

// AddEditCappState adds the EditCappState route to the OpenAPI scheme.
func AddEditCappState(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "edit-capp-state",
		Method:      http.MethodPut,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/state", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Edit state of a Capp in a namespace",
		Description: "Changes the state field a specific Capp in a specific namespace to disabled or enabled",
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
			Content: map[string]*huma.MediaType{
				applicationJSONKey: {
					Examples: map[string]*huma.Example{
						"Change to disabled": {
							Value: types.CappState{
								State: "disabled",
							},
						},
						"Change to enabled": {
							Value: types.CappState{
								State: "enabled",
							},
						},
					},
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

// AddGetCappState adds the GetCappState route to the OpenAPI scheme.
func AddGetCappState(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-state",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/state", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Get state of a Capp in a namespace",
		Description: "Retrieves the state of a specific Capp in a specific namespace",
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

// AddDeleteCapp adds the DeleteCapp route to the OpenAPI scheme.
func AddDeleteCapp(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "delete-capp",
		Method:      http.MethodDelete,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Delete a Capp in a namespace",
		Description: "Deletes a specific Capp in a specific namespace",
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

// AddGetCappDNS adds the GetCappDNS route to the OpenAPI scheme.
func AddGetCappDNS(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-dns",
		Method:      http.MethodGet,
		Tags:        []string{cappTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/dns", namespacesKey, namespaceNameKey, cappsKey, cappNameKey),
		Summary:     "Get DNS records of a Capp in a namespace",
		Description: "Retrieves the DNS records of a specific Capp in a specific namespace",
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
