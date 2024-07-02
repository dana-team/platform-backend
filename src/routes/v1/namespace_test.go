package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	namespacesKey = "namespaces"
	nameKey       = "name"
	nsName        = testNamespace + "-" + namespacesKey
)

func TestGetNamespaces(t *testing.T) {
	testNamespaceName := nsName + "-get"

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		want want
	}{
		"ShouldSucceedGettingNamespaces": {
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count:         2,
					namespacesKey: []types.Namespace{{Name: testNamespaceName + "-1"}, {Name: testNamespaceName + "-2"}},
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")
	createTestNamespace(testNamespaceName + "-2")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := "/v1/namespaces/"
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

func TestGetNamespace(t *testing.T) {
	testNamespaceName := nsName + "-get-one"

	type requestParams struct {
		name string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingNamespace": {
			requestParams: requestParams{
				name: testNamespaceName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					nameKey: testNamespaceName + "-1",
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				name: testNamespaceName + "-1" + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					errorKey:   operationFailed,
					detailsKey: fmt.Sprintf("%s %q not found", namespacesKey, testNamespaceName+"-1"+nonExistentSuffix),
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestParams.name)
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

func TestCreateNamespace(t *testing.T) {
	testNamespaceName := nsName + "-create"

	type bodyParams struct {
		ns types.Namespace
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedCreatingNamespace": {
			bodyParams: bodyParams{
				ns: mocks.PrepareNamespaceType(testNamespaceName),
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					nameKey: testNamespaceName,
				},
			},
		},
		"ShouldFailWithBadRequestBody": {
			bodyParams: bodyParams{
				ns: types.Namespace{},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'Namespace.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					errorKey:   invalidRequest,
				},
			},
		},
		"ShouldHandleAlreadyExists": {
			bodyParams: bodyParams{
				ns: mocks.PrepareNamespaceType(testNamespaceName + "-1"),
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q already exists", namespacesKey, testNamespaceName+"-1"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.ns)
			assert.NoError(t, err)

			baseURI := "/v1/namespaces/"
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(body))
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

func TestDeleteNamespace(t *testing.T) {
	testNamespaceName := nsName + "-delete"

	type requestParams struct {
		name string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedDeletingNamespace": {
			requestParams: requestParams{
				name: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					messageKey: fmt.Sprintf("Deleted namespace successfully %q", testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				name: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", namespacesKey, testNamespaceName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
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
