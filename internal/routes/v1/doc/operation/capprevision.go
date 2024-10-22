package operation

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/danielgtaylor/huma/v2"
)

const cappRevisionTag = "CappRevisions"

// AddGetCappRevisions adds the GetCappRevisions route to the OpenAPI scheme.
func AddGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revisions",
		Method:      http.MethodGet,
		Tags:        []string{cappRevisionTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", namespacesKey, namespaceNameKey, cappsKey, cappNameKey, cappRevisionsKey),
		Summary:     "Get all CappsRevisions of a Capp in a namespace",
		Description: "Retrieves all the names of the CappRevisions of a specific Capp in a specific namespace",
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

// AddGetCappRevision adds the GetCappRevision route to the OpenAPI scheme.
func AddGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capp-revision",
		Method:      http.MethodGet,
		Tags:        []string{cappRevisionTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}", namespacesKey, namespaceNameKey, cappsKey, cappNameKey, cappRevisionsKey, cappRevisionNameKey),
		Summary:     "Get a CappsRevision of a Capp in a namespace",
		Description: "Retrieves a specific CappRevision of a specific Capp in specific namespace",
		Parameters: []*huma.Param{
			{
				Name:     namespaceNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.NamespaceName)),
				Example:  defaultExample,
			},
			{
				Name:     cappNameKey,
				In:       pathKey,
				Required: true,
				Schema:   huma.SchemaFromType(registry, reflect.TypeOf(types.CappRevisionUri{}.CappName)),
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

// AddClusterGetCappRevisions adds the ClusterGetCappRevisions route to the OpenAPI scheme.
func AddClusterGetCappRevisions(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevisions-cluster",
		Method:      http.MethodGet,
		Tags:        []string{cappRevisionTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, cappRevisionsKey),
		Summary:     "Get all CappsRevisions of a Capp in a namespace with cluster name",
		Description: "Retrieves all the names of the CappRevisions of a specific Capp in a specific namespace using a cluster name",
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

// AddClusterGetCappRevision adds the ClusterGetCappRevision route to the OpenAPI scheme.
func AddClusterGetCappRevision(api huma.API, registry huma.Registry) {
	operation := &huma.Operation{
		OperationID: "get-capprevision-cluster",
		Method:      http.MethodGet,
		Tags:        []string{cappRevisionTag},
		Path:        fmt.Sprintf("/v1/%s/{%s}/%s/{%s}/%s/{%s}", clustersKey, clusterNameKey, namespacesKey, namespaceNameKey, cappRevisionsKey, cappRevisionNameKey),
		Summary:     "Get a CappsRevision of a Capp in a namespace with cluster name",
		Description: "Retrieves a specific CappRevision of a specific Capp in specific namespace with cluster name",
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
