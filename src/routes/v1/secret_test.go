package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

// createTestSecret creates a test Secret object.
func createTestSecret(name, namespace string) {
	secret := mocks.PrepareSecret(name, namespace, testutils.SecretDataKey, testutils.SecretDataValueEncoded)
	_, err := fakeClient.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetSecrets(t *testing.T) {
	testNamespaceName := testutils.SecretNamespace + "-get"

	type requestURI struct {
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
		"ShouldSucceedGettingSecrets": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					testutils.Count: 2,
					testutils.SecretsKey: []types.Secret{
						{SecretName: testutils.SecretName + "-1", NamespaceName: testNamespaceName, Type: testutils.SecretType},
						{SecretName: testutils.SecretName + "-2", NamespaceName: testNamespaceName, Type: testutils.SecretType}},
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(testutils.SecretName+"-1", testNamespaceName)
	createTestSecret(testutils.SecretName+"-2", testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets", test.requestURI.namespace)
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
					testutils.TypeKey:       testutils.SecretType,
					testutils.Data:          []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(testutils.SecretName, testNamespaceName)

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
		response   map[string]string
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
				response: map[string]string{
					testutils.SecretNameKey:    testutils.SecretName,
					testutils.TypeKey:          testutils.SecretType,
					testutils.NameSpaceNameKey: testNamespaceName,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(testutils.SecretName, strings.ToLower(testutils.OpaqueType), "", "",
				[]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			requestData: map[string]interface{}{
				testutils.SecretNameKey: testutils.SecretName,
				testutils.TypeKey:       strings.ToLower(testutils.OpaqueType),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]string{
					testutils.DetailsKey: "data is required for Opaque secrets",
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]string{
					testutils.DetailsKey: fmt.Sprintf("%s %q already exists", testutils.SecretsKey, testutils.SecretName+"-1"),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(testutils.SecretName+"-1", strings.ToLower(testutils.OpaqueType), "", "",
				[]types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}}),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(testutils.SecretName+"-1", testNamespaceName)

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
					testutils.TypeKey:          testutils.SecretType,
					testutils.NameSpaceNameKey: testNamespaceName,
					testutils.Data:             []types.KeyValue{{Key: testutils.SecretDataKey, Value: testutils.SecretDataValue}},
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName),
					testutils.ErrorKey:   testutils.OperationFailed,
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
					testutils.DetailsKey: "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
					testutils.ErrorKey:   testutils.InvalidRequest,
				},
			},
			requestData: map[string]interface{}{},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(testutils.SecretName, testNamespaceName)

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
		"ShouldSucceedDeletingSecret": {
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
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
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.SecretsKey, testutils.SecretName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(testutils.SecretName, testNamespaceName)

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
