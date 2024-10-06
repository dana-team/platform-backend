package operation

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/danielgtaylor/huma/v2"
)

const revisionTag = "CappRevisions"

func AddGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revisions",
		Method:      http.MethodGet,
		Tags:        []string{revisionTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/capprevisions",
		Summary:     "[CappRevisions] Get CappsRevisions of a Capp",
		Description: "Get all CappRevisions of a Capp",
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
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionList{})),
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

func AddGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revision",
		Method:      http.MethodGet,
		Tags:        []string{revisionTag},
		Path:        "/v1/namespaces/{namespaceName}/capps/{cappName}/capprevisions/{cappRevisionName}",
		Summary:     "[CappRevisions] Get a CappsRevision of a Capp",
		Description: "Get a specific CappRevision of a Capp",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.CappName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevision{})),
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

func AddClusterGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevisions-cluster",
		Method:      http.MethodGet,
		Tags:        []string{revisionTag},
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/capprevisions",
		Summary:     "[CappRevisions] Get CappRevisions with cluster",
		Description: "Get all CappRevisions by namespace and cluster",
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
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionNamespaceUri{}.NamespaceName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionList{})),
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

func AddClusterGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevision-cluster",
		Method:      http.MethodGet,
		Tags:        []string{revisionTag},
		Path:        "/v1/clusters/{clusterName}/namespaces/{namespaceName}/capprevisions/{cappRevisionName}",
		Summary:     "[CappRevisions] Get a CappRevision with cluster",
		Description: "Get a specific CappRevision by namespace and cluster",
		Parameters: []*huma.Param{
			{
				Name:     clusterNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(huma.TypeString)),
				Example:  defaultExample,
			},
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappRevisionNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.CappRevisionName)),
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
						Schema: huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevision{})),
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
