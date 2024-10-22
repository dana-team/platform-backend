package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGetSecrets(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-get"

	type pagination struct {
		limit string
		page  string
	}

	type requestURI struct {
		namespace        string
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
		"ShouldSucceedGettingSecrets": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 2,
					testutils.SecretsKey: []types.Secret{
						{SecretName: testutils.SecretName + "-1", NamespaceName: testNamespaceName, Type: string(corev1.SecretTypeOpaque)},
						{SecretName: testutils.SecretName + "-2", NamespaceName: testNamespaceName, Type: string(corev1.SecretTypeOpaque)}},
				},
			},
		},
		"ShouldSucceedGettingAllSecretsWithLimitOf2": {
			requestURI: requestURI{
				namespace:        testNamespaceName,
				paginationParams: pagination{limit: "2", page: "1"},
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.CountKey: 2,
					testutils.SecretsKey: []types.Secret{
						{SecretName: testutils.SecretName + "-1", NamespaceName: testNamespaceName, Type: string(corev1.SecretTypeOpaque)},
						{SecretName: testutils.SecretName + "-2", NamespaceName: testNamespaceName, Type: string(corev1.SecretTypeOpaque)}},
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName+"-1", testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName+"-2", testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}
			if test.requestURI.paginationParams.limit != "" {
				params.Add(middleware.LimitCtxKey, test.requestURI.paginationParams.limit)
			}

			if test.requestURI.paginationParams.page != "" {
				params.Add(middleware.PageCtxKey, test.requestURI.paginationParams.page)
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets", test.requestURI.namespace)
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

func TestGetSecret(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-get-one"

	type requestURI struct {
		namespace string
		name      string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.SecretNameKey: testutils.SecretName,
					testutils.IdKey:         "",
					testutils.TypeKey:       string(corev1.SecretTypeOpaque),
					testutils.DataKey:       []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
				},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetSecret, testutils.SecretName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetSecret, testutils.SecretName),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets/%s", test.requestURI.namespace, test.requestURI.name)
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

func TestCreateSecret(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-create"

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
		"ShouldSucceedCreateSecret": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.SecretNameKey:    testutils.SecretName,
					testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					testutils.NamespaceNameKey: testNamespaceName,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(testutils.SecretName, strings.ToLower(string(corev1.SecretTypeOpaque)), "", "",
				[]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			requestData: map[string]interface{}{
				testutils.SecretNameKey: testutils.SecretName,
				testutils.TypeKey:       strings.ToLower(string(corev1.SecretTypeOpaque)),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "data is required for Opaque secrets",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotCreateSecret, testutils.SecretName+"-1"),
						fmt.Sprintf("%s %q already exists", testutils.SecretsKey, testutils.SecretName+"-1")),
					testutils.ReasonKey: metav1.StatusReasonAlreadyExists,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(testutils.SecretName+"-1", strings.ToLower(string(corev1.SecretTypeOpaque)), "", "",
				[]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName+"-1", testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets", test.requestURI.namespace)
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

func TestUpdateSecret(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-update"

	type requestURI struct {
		namespace string
		name      string
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
		"ShouldSucceedUpdatingSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.SecretNameKey:    testutils.SecretName,
					testutils.TypeKey:          string(corev1.SecretTypeOpaque),
					testutils.NamespaceNameKey: testNamespaceName,
					testutils.DataKey:          []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetSecret, testutils.SecretName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetSecret, testutils.SecretName),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.ErrorKey:  "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
					testutils.ReasonKey: metav1.StatusReasonBadRequest,
				},
			},
			requestData: map[string]interface{}{},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

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

func TestDeleteSecret(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-delete"

	type requestURI struct {
		namespace string
		name      string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"paShouldSucceedDeletingSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.MessageKey: fmt.Sprintf("Deleted secret %q in namespace %q successfully", testutils.SecretName, testNamespaceName)},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      testutils.SecretName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotDeleteSecret, testutils.SecretName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      testutils.SecretName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotDeleteSecret, testutils.SecretName),
						fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestSecret(fakeClient, testutils.SecretName, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
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
