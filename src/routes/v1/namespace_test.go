package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	nsName = testutils.TestNamespace + "-" + testutils.NameSpaceKey
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
					testutils.Count:        2,
					testutils.NameSpaceKey: []types.Namespace{{Name: testNamespaceName + "-1"}, {Name: testNamespaceName + "-2"}},
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")
	createTestNamespace(testNamespaceName + "-2")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := "/v1/namespaces"
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

	type requestURI struct {
		name string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingNamespace": {
			requestURI: requestURI{
				name: testNamespaceName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.NameKey: testNamespaceName + "-1",
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name: testNamespaceName + "-1" + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey:   testutils.OperationFailed,
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.NameSpaceKey, testNamespaceName+"-1"+testutils.NonExistentSuffix),
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestURI.name)
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

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		want        want
		requestData interface{}
	}{
		"ShouldSucceedCreatingNamespace": {
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.NameKey: testNamespaceName,
				},
			},
			requestData: mocks.PrepareNamespaceType(testNamespaceName),
		},
		"ShouldFailWithBadRequestBody": {
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.DetailsKey: "Key: 'Namespace.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					testutils.ErrorKey:   testutils.InvalidRequest,
				},
			},
			requestData: map[string]interface{}{},
		},
		"ShouldHandleAlreadyExists": {
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s %q already exists", testutils.NameSpaceKey, testNamespaceName+"-1"),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareNamespaceType(testNamespaceName + "-1"),
		},
	}

	setup()
	createTestNamespace(testNamespaceName + "-1")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := "/v1/namespaces"
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

func TestDeleteNamespace(t *testing.T) {
	testNamespaceName := nsName + "-delete"

	type requestURI struct {
		name string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedDeletingNamespace": {
			requestURI: requestURI{
				name: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MessageKey: fmt.Sprintf("Deleted namespace successfully %q", testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.NameSpaceKey, testNamespaceName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestURI.name)
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
