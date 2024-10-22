package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

const (
	nsName = testutils.TestNamespace + "-" + testutils.NamespaceKey
)

func TestGetNamespaces(t *testing.T) {
	testNamespaceName := nsName + "-get"

	type pagination struct {
		limit string
		page  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	type args struct {
		paginationParams pagination
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSucceedGettingNamespaces": {
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey:     2,
					testutils.NamespaceKey: []types.Namespace{{Name: testNamespaceName + "-1"}, {Name: testNamespaceName + "-2"}},
				},
			},
		},
		"ShouldSucceedGettingAllNamespacesWithLimitOf2": {
			args: args{
				paginationParams: pagination{limit: "2", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey:     2,
					testutils.NamespaceKey: []types.Namespace{{Name: testNamespaceName + "-1"}, {Name: testNamespaceName + "-2"}},
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName+"-1")
	mocks.CreateTestNamespace(fakeClient, testNamespaceName+"-2")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			if test.args.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.args.paginationParams.limit)
			}

			if test.args.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.args.paginationParams.page)
			}

			baseURI := "/v1/namespaces"
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
					testutils.ReasonKey: metav1.StatusReasonNotFound,
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotFetchNamespace, testNamespaceName+"-1"+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.NamespaceKey, testNamespaceName+"-1"+testutils.NonExistentSuffix)),
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName+"-1")

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
					testutils.ErrorKey:  "Key: 'Namespace.Name' Error:Field validation for 'Name' failed on the 'required' tag",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: map[string]interface{}{},
		},
		"ShouldHandleAlreadyExists": {
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotCreateNamespace, testNamespaceName+"-1"),
						fmt.Sprintf("%s %q already exists", testutils.NamespaceKey, testNamespaceName+"-1")),
					testutils.ReasonKey: metav1.StatusReasonAlreadyExists,
				},
			},
			requestData: mocks.PrepareNamespaceType(testNamespaceName + "-1"),
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName+"-1")

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
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotDeleteNamespace, testNamespaceName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.NamespaceKey, testNamespaceName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)

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
