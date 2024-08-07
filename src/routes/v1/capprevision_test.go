package v1

import (
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	cappRevisionName      = testutils.TestName + "-capp-revision"
	capprevisionsKey      = "capprevisions"
	cappRevisionNamespace = testutils.TestNamespace + "-" + capprevisionsKey
)

func TestGetCappRevisions(t *testing.T) {
	testNamespaceName := cappRevisionNamespace + "-get"

	type selector struct {
		keys   []string
		values []string
	}

	type pagination struct {
		limit string
		page  string
	}

	type requestURI struct {
		namespace        string
		labelSelector    selector
		paginationParams pagination
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCappRevisions": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldSucceedGettingAllCappRevisionsWithLimitOf2": {
			requestURI: requestURI{
				namespace:        testNamespaceName,
				paginationParams: pagination{limit: "2", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					testutils.CountKey: 2,
				},
			},
		},
		"ShouldSucceedGettingCappRevisionsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + "-1"},
					values: []string{testutils.LabelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   []string{cappRevisionName + "-1"},
					testutils.CountKey: 1,
				},
			},
		},
		"ShouldFailGettingCappRevisionsWithInvalidLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + "-1"},
					values: []string{testutils.LabelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.DetailsKey: "found '1', expected: ',' or 'end of string'",
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldSucceedGettingNoCappRevisionsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + "-3"},
					values: []string{testutils.LabelValue + "-3"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey:   nil,
					testutils.CountKey: 0,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCappRevision(dynClient, cappRevisionName+"-1", testNamespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, nil)
	mocks.CreateTestCappRevision(dynClient, cappRevisionName+"-2", testNamespaceName, map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestURI.labelSelector.keys {
				params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", key, test.requestURI.labelSelector.values[i]))
			}

			if test.requestURI.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.requestURI.paginationParams.limit)
			}

			if test.requestURI.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.requestURI.paginationParams.page)
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCappRevision(t *testing.T) {
	testNamespaceName := cappRevisionNamespace + "-get-one"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      cappRevisionName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MetadataKey:    types.Metadata{Name: cappRevisionName, Namespace: testNamespaceName},
					testutils.LabelsKey:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
					testutils.AnnotationsKey: nil,
					testutils.SpecKey:        mocks.PrepareCappRevisionSpec(),
					testutils.StatusKey:      mocks.PrepareCappRevisionStatus(),
				},
			},
		},
		"ShouldHandleNotFoundCappRevision": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      cappRevisionName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", capprevisionsKey, cappv1alpha1.GroupVersion.Group, cappRevisionName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      cappRevisionName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("capprevisions.%s %q not found", cappv1alpha1.GroupVersion.Group, cappRevisionName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestCappRevision(dynClient, cappRevisionName, testNamespaceName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
