package v1_test

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
	cappNamespace = testName + "-capp-ns"
	cappName      = testName + "-capp"
	capps         = "capps"
	message       = "message"
)

func setupCapps() {
	createTestNamespace(cappNamespace)
	createTestCapp(cappName+"-1", cappNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)
	createTestCapp(cappName+"-2", cappNamespace, map[string]string{labelKey + "-2": labelValue + "-2"}, nil)
	createTestCappWithHostname(cappName+"-3", cappNamespace, map[string]string{labelKey + "-3": labelValue + "-3"}, nil)
	createTestCappWithHostname(cappName+"-4", cappNamespace, map[string]string{labelKey + "-4": labelValue + "-4"}, nil)
}

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
	type selector struct {
		keys   []string
		values []string
	}

	type requestParams struct {
		namespace     string
		labelSelector selector
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCapps": {
			requestParams: requestParams{
				namespace: cappNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count: 4,
					capps: []types.CappSummary{
						{Name: cappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-1", cappNamespace, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-2", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-2", cappNamespace, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-3", URL: fmt.Sprintf("https://%s.%s", mocks.Hostname, mocks.Domain), Images: []string{mocks.CappImage}},
						{Name: cappName + "-4", URL: fmt.Sprintf("https://%s.%s", mocks.Hostname, mocks.Domain), Images: []string{mocks.CappImage}},
					},
				},
			},
		},
		"ShouldSucceedGettingCappsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count: 1,
					capps: []types.CappSummary{
						{Name: cappName + "-1", URL: fmt.Sprintf("https://%s-%s.%s", cappName+"-1", cappNamespace, mocks.Domain), Images: []string{mocks.CappImage}},
					},
				},
			},
		},
		"ShouldFailGettingCappsWithInvalidLabelSelector": {
			requestParams: requestParams{
				namespace: cappNamespace,
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
			requestParams: requestParams{
				namespace: cappNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-5"},
					values: []string{labelValue + "-5"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count: 0,
					capps: nil,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestParams.labelSelector.keys {
				params.Add(labelSelectorKey, key+"="+test.requestParams.labelSelector.values[i])
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/", test.requestParams.namespace)

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetCapp(t *testing.T) {
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
		"ShouldSucceedGettingCapp": {
			requestParams: requestParams{
				namespace: cappNamespace,
				name:      cappName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappName + "-1", Namespace: cappNamespace},
					labels:      []types.KeyValue{{Key: labelKey + "-1", Value: labelValue + "-1"}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      mocks.PrepareCappStatus(cappName+"-1", cappNamespace),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestParams: requestParams{
				namespace: cappNamespace,
				name:      cappName + "-1" + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-1"+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				namespace: cappNamespace + nonExistentSuffix,
				name:      cappName + "-1",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-1"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestCreateCapp(t *testing.T) {
	type bodyParams struct {
		capp types.Capp
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedCreatingCapp": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappTypeWithoutStatus(cappName+"-6", cappNamespace, map[string]string{labelKey + "-6": labelValue + "-6"}, nil),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappName + "-6", Namespace: cappNamespace},
					labels:      []types.KeyValue{{Key: labelKey + "-6", Value: labelValue + "-6"}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      cappv1alpha1.CappStatus{},
				},
			},
		},
		"ShouldFailWithBadRequestBody": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappTypeWithoutStatus("", cappNamespace, map[string]string{labelKey + "-6": labelValue + "-6"}, nil),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'CreateCapp.Metadata.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					errorKey:   invalidRequest,
				},
			},
		},
		"ShouldHandleAlreadyExists": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappTypeWithoutStatus(cappName+"-1", cappNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, nil),
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q already exists", capps, cappv1alpha1.GroupVersion.Group, cappName+"-1"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.capp)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/", test.bodyParams.capp.Metadata.Namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(body))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestUpdateCapp(t *testing.T) {
	type bodyParams struct {
		capp types.Capp
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedUpdatingCapp": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappType(cappName+"-7", cappNamespace, map[string]string{labelKey + "-7-updated": labelValue + "-7-updated"}, nil),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					metadata:    types.Metadata{Name: cappName + "-7", Namespace: cappNamespace},
					labels:      []types.KeyValue{{Key: labelKey + "-7-updated", Value: labelValue + "-7-updated"}},
					annotations: nil,
					spec:        mocks.PrepareCappSpec(),
					status:      mocks.PrepareCappStatus(cappName+"-7", cappNamespace),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappType(cappName+"-7"+nonExistentSuffix, cappNamespace, map[string]string{labelKey + "-7-updated": labelValue + "-7-updated"}, nil),
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-7"+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			bodyParams: bodyParams{
				capp: mocks.PrepareCappType(cappName+"-7", cappNamespace+nonExistentSuffix, map[string]string{labelKey + "-7-updated": labelValue + "-7-updated"}, nil),
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-7"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	createTestCapp(cappName+"-7", cappNamespace, map[string]string{labelKey + "-7": labelValue + "-7"}, nil)
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.capp)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.bodyParams.capp.Metadata.Namespace, test.bodyParams.capp.Metadata.Name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(body))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestDeleteCapp(t *testing.T) {
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
		"ShouldSucceedDeletingCapp": {
			requestParams: requestParams{
				name:      cappName + "-9",
				namespace: cappNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					message: fmt.Sprintf("Deleted capp %q in namespace %q successfully", cappName+"-9", cappNamespace),
				},
			},
		},
		"ShouldHandleNotFoundCapp": {
			requestParams: requestParams{
				name:      cappName + "-9" + nonExistentSuffix,
				namespace: cappNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-9"+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				name:      cappName + "-9",
				namespace: cappNamespace + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", capps, cappv1alpha1.GroupVersion.Group, cappName+"-9"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	createTestCapp(cappName+"-9", cappNamespace, map[string]string{labelKey + "-9": labelValue + "-9"}, nil)
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
