package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	cappRevisionName      = testName + "-capp-revision"
	capprevisionsKey      = "capprevisions"
	cappRevisionNamespace = testNamespace + "-" + capprevisionsKey
)

// createTestCappRevision creates a test CappRevision object.
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := mocks.PrepareCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func TestGetCappRevisions(t *testing.T) {
	testNamespaceName := cappRevisionNamespace + "-get"

	type selector struct {
		keys   []string
		values []string
	}

	type requestURI struct {
		namespace     string
		labelSelector selector
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
					capprevisionsKey: []string{cappRevisionName + "-1", cappRevisionName + "-2"},
					count:            2,
				},
			},
		},
		"ShouldSucceedGettingCappRevisionsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey: []string{cappRevisionName + "-1"},
					count:            1,
				},
			},
		},
		"ShouldFailGettingCappRevisionsWithInvalidLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "found '1', expected: ',' or 'end of string'",
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldSucceedGettingNoCappRevisionsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{labelKey + "-3"},
					values: []string{labelValue + "-3"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					capprevisionsKey: nil,
					count:            0,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestCappRevision(cappRevisionName+"-1", testNamespaceName, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)
	createTestCappRevision(cappRevisionName+"-2", testNamespaceName, map[string]string{labelKey + "-2": labelValue + "-2"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestURI.labelSelector.keys {
				params.Add(labelSelectorKey, fmt.Sprintf("%s=%s", key, test.requestURI.labelSelector.values[i]))
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/", test.requestURI.namespace)
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

	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestParams: requestParams{
				namespace: testNamespaceName,
				name:      cappRevisionName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappRevisionName, Namespace: testNamespaceName},
					labels:      []types.KeyValue{{Key: labelKey, Value: labelValue}},
					annotations: nil,
					spec:        mocks.PrepareCappRevisionSpec(),
					status:      mocks.PrepareCappRevisionStatus(),
				},
			},
		},
		"ShouldHandleNotFoundCappRevision": {
			requestParams: requestParams{
				namespace: testNamespaceName,
				name:      cappRevisionName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capprevisionsKey, cappv1alpha1.GroupVersion.Group, cappRevisionName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				namespace: testNamespaceName + nonExistentSuffix,
				name:      cappRevisionName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("capprevisions.%s %q not found", cappv1alpha1.GroupVersion.Group, cappRevisionName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestCappRevision(cappRevisionName, testNamespaceName, map[string]string{labelKey: labelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/%s", test.requestParams.namespace, test.requestParams.name)
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
