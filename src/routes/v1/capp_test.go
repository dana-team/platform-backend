package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCapps(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get"

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
		"ShouldSucceedGettingCapps": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Count: 4,
					testutils.CappsKey: []types.CappSummary{
						{Name: testutils.CappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-1", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-2", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-2", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-3", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
						{Name: testutils.CappName + "-4", URL: fmt.Sprintf("https://%s.%s", testutils.Hostname, testutils.Domain), Images: []string{testutils.CappImage}},
					},
				},
			},
		},
		"ShouldSucceedGettingCappsWithLabelSelector": {
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
					testutils.Count: 1,
					testutils.CappsKey: []types.CappSummary{
						{Name: testutils.CappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", testutils.CappName+"-1", testNamespaceName, testutils.Domain), Images: []string{testutils.CappImage}},
					},
				},
			},
		},
		"ShouldFailGettingCappsWithInvalidLabelSelector": {
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
		"ShouldSucceedGettingNoCappsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{testutils.LabelKey + testutils.NonExistentSuffix},
					values: []string{testutils.LabelValue + testutils.NonExistentSuffix},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Count:    0,
					testutils.CappsKey: nil,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	mocks.CreateTestCapp(testutils.CappName+"-1", testNamespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, nil, dynClient)
	mocks.CreateTestCapp(testutils.CappName+"-2", testNamespaceName, map[string]string{testutils.LabelKey + "-2": testutils.LabelValue + "-2"}, nil, dynClient)
	mocks.CreateTestCappWithHostname(testutils.CappName+"-3", testNamespaceName, map[string]string{testutils.LabelKey + "-3": testutils.LabelValue + "-3"}, nil, dynClient)
	mocks.CreateTestCappWithHostname(testutils.CappName+"-4", testNamespaceName, map[string]string{testutils.LabelKey + "-4": testutils.LabelValue + "-4"}, nil, dynClient)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestURI.labelSelector.keys {
				params.Add(testutils.LabelSelectorKey, fmt.Sprintf("%s=%s", key, test.requestURI.labelSelector.values[i]))
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps", test.requestURI.namespace)

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

func TestGetCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-get-one"

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
		"ShouldSucceedGettingCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Metadata:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.Labels:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
					testutils.Annotations: nil,
					testutils.Spec:        mocks.PrepareCappSpec(),
					testutils.Status:      mocks.PrepareCappStatus(testutils.CappName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      testutils.CappName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      testutils.CappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	mocks.CreateTestCapp(testutils.CappName, testNamespaceName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil, dynClient)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
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

func TestCreateCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-create"

	type requestURI struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedCreatingCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Metadata:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.Labels:      []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}},
					testutils.Annotations: nil,
					testutils.Spec:        mocks.PrepareCappSpec(),
					testutils.Status:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName, []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.DetailsKey: "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					testutils.ErrorKey:   testutils.InvalidRequest,
				},
			},
			requestData: mocks.PrepareCreateCappType("", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q already exists", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+"-1"),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareCreateCappType(testutils.CappName+"-1", []types.KeyValue{{Key: testutils.LabelKey, Value: testutils.LabelValue}}, nil),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	mocks.CreateTestCapp(testutils.CappName+"-1", testNamespaceName, map[string]string{testutils.LabelKey + "-1": testutils.LabelValue + "-1"}, nil, dynClient)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

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

func TestUpdateCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-update"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedUpdatingCapp": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Metadata:    types.Metadata{Name: testutils.CappName, Namespace: testNamespaceName},
					testutils.Labels:      []types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}},
					testutils.Annotations: nil,
					testutils.Spec:        mocks.PrepareCappSpec(),
					testutils.Status:      mocks.PrepareCappStatus(testutils.CappName, testNamespaceName),
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: testutils.LabelKey + "-updated", Value: testutils.LabelValue + "-updated"}}, nil),
		},
	}

	setup()
	mocks.CreateTestCapp(testutils.CappName, testNamespaceName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil, dynClient)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

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

func TestDeleteCapp(t *testing.T) {
	testNamespaceName := testutils.CappNamespace + "-delete"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestURI
		want          want
	}{
		"ShouldSucceedDeletingCapp": {
			requestParams: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MessageKey: fmt.Sprintf("Deleted capp %q in namespace %q successfully", testutils.CappName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestParams: requestURI{
				name:      testutils.CappName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestURI{
				name:      testutils.CappName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.CappsKey, cappv1alpha1.GroupVersion.Group, testutils.CappName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	mocks.CreateTestCapp(testutils.CappName, testNamespaceName, map[string]string{testutils.LabelKey: testutils.LabelValue}, nil, dynClient)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
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
