package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	cappName      = testName + "-capp"
	cappsKey      = "capps"
	cappNamespace = testNamespace + "-" + cappsKey
)

// createTestCapp creates a test Capp object.
func createTestCapp(name, namespace string, labels, annotations map[string]string) {
	capp := mocks.PrepareCapp(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

// createTestCapp creates a test Capp object.
func createTestCappWithHostname(name, namespace string, labels, annotations map[string]string) {
	capp := mocks.PrepareCappWithHostname(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

func TestGetCapps(t *testing.T) {
	testNamespaceName := cappNamespace + "-get"

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
					count: 4,
					cappsKey: []types.CappSummary{
						{Name: cappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-1", testNamespaceName, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-2", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-2", testNamespaceName, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-3", URL: fmt.Sprintf("https://%s.%s", mocks.Hostname, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-4", URL: fmt.Sprintf("https://%s.%s", mocks.Hostname, mocks.Domain), Images: []string{mocks.CappImage}},
					},
				},
			},
		},
		"ShouldSucceedGettingCappsWithLabelSelector": {
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
					count: 1,
					cappsKey: []types.CappSummary{
						{Name: cappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-1", testNamespaceName, mocks.Domain), Images: []string{mocks.CappImage}},
					},
				},
			},
		},
		"ShouldFailGettingCappsWithInvalidLabelSelector": {
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
		"ShouldSucceedGettingNoCappsWithLabelSelector": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				labelSelector: selector{
					keys:   []string{labelKey + nonExistentSuffix},
					values: []string{labelValue + nonExistentSuffix},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count:    0,
					cappsKey: nil,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestCapp(cappName+"-1", testNamespaceName, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)
	createTestCapp(cappName+"-2", testNamespaceName, map[string]string{labelKey + "-2": labelValue + "-2"}, nil)
	createTestCappWithHostname(cappName+"-3", testNamespaceName, map[string]string{labelKey + "-3": labelValue + "-3"}, nil)
	createTestCappWithHostname(cappName+"-4", testNamespaceName, map[string]string{labelKey + "-4": labelValue + "-4"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestURI.labelSelector.keys {
				params.Add(labelSelectorKey, fmt.Sprintf("%s=%s", key, test.requestURI.labelSelector.values[i]))
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
	testNamespaceName := cappNamespace + "-get-one"

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
				name:      cappName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappName, Namespace: testNamespaceName},
					labels:      []types.KeyValue{{Key: labelKey, Value: labelValue}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      mocks.PrepareCappStatus(cappName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      cappName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + nonExistentSuffix,
				name:      cappName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestCapp(cappName, testNamespaceName, map[string]string{labelKey: labelValue}, nil)

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
	testNamespaceName := cappNamespace + "-create"

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
					metadata:    types.Metadata{Name: cappName, Namespace: testNamespaceName},
					labels:      []types.KeyValue{{Key: labelKey, Value: labelValue}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      cappv1alpha1.CappStatus{},
				},
			},
			requestData: mocks.PrepareCreateCappType(cappName, []types.KeyValue{{Key: labelKey, Value: labelValue}}, nil),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					errorKey:   invalidRequest,
				},
			},
			requestData: mocks.PrepareCreateCappType("", []types.KeyValue{{Key: labelKey, Value: labelValue}}, nil),
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q already exists", cappsKey, cappv1alpha1.GroupVersion.Group, cappName+"-1"),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareCreateCappType(cappName+"-1", []types.KeyValue{{Key: labelKey, Value: labelValue}}, nil),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestCapp(cappName+"-1", testNamespaceName, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

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
	testNamespaceName := cappNamespace + "-update"

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
				name:      cappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappName, Namespace: testNamespaceName},
					labels:      []types.KeyValue{{Key: labelKey + "-updated", Value: labelValue + "-updated"}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      mocks.PrepareCappStatus(cappName, testNamespaceName),
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: labelKey + "-updated", Value: labelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				name:      cappName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: labelKey + "-updated", Value: labelValue + "-updated"}}, nil),
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      cappName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareUpdateCappType([]types.KeyValue{{Key: labelKey + "-updated", Value: labelValue + "-updated"}}, nil),
		},
	}

	setup()
	createTestCapp(cappName, testNamespaceName, map[string]string{labelKey: labelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

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
	testNamespaceName := cappNamespace + "-delete"

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
		"ShouldSucceedDeletingCapp": {
			requestURI: requestURI{
				name:      cappName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					messageKey: fmt.Sprintf("Deleted capp %q in namespace %q successfully", cappName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestURI: requestURI{
				name:      cappName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      cappName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", cappsKey, cappv1alpha1.GroupVersion.Group, cappName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestCapp(cappName, testNamespaceName, map[string]string{labelKey: labelValue}, nil)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestURI.namespace, test.requestURI.name)
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
