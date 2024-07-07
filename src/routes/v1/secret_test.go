package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	secretsKey             = "secrets"
	secretNamespace        = testNamespace + "-" + secretsKey
	secretName             = testName + "-secret"
	secretDataKey          = "test-key"
	secretDataValue        = "fake"
	secretDataValueEncoded = "ZmFrZQ=="
	secretNameKey          = "secretName"
	dataKey                = "data"
	idKey                  = "id"
	typeKey                = "type"
	namespaceKey           = "namespaceName"
	opaqueKey              = "Opaque"
)

// createTestSecret creates a test Secret object.
func createTestSecret(name, namespace string) {
	secret := mocks.PrepareSecret(name, namespace, secretDataKey, secretDataValueEncoded)
	_, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetSecrets(t *testing.T) {
	testNamespaceName := secretNamespace + "-get"

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
					count: 2,
					secretsKey: []types.Secret{
						{SecretName: secretName + "-1", NamespaceName: testNamespaceName, Type: opaqueKey},
						{SecretName: secretName + "-2", NamespaceName: testNamespaceName, Type: opaqueKey}},
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(secretName+"-1", testNamespaceName)
	createTestSecret(secretName+"-2", testNamespaceName)

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

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}

func TestGetSecret(t *testing.T) {
	testNamespaceName := secretNamespace + "-get-one"

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
				name:      secretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					secretNameKey: secretName,
					idKey:         "",
					typeKey:       opaqueKey,
					dataKey:       []types.KeyValue{{Key: secretDataKey, Value: secretDataValue}},
				},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      secretName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      secretName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(secretName, testNamespaceName)

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
	testNamespaceName := secretNamespace + "-create"

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
					secretNameKey: secretName,
					typeKey:       opaqueKey,
					namespaceKey:  testNamespaceName,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(secretName, strings.ToLower(opaqueKey), "", "",
				[]types.KeyValue{{Key: secretDataKey, Value: secretDataValue}}),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			requestData: map[string]interface{}{
				secretNameKey: secretName,
				typeKey:       strings.ToLower(opaqueKey),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]string{
					detailsKey: "data is required for Opaque secrets",
					errorKey:   operationFailed,
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
					detailsKey: fmt.Sprintf("%s %q already exists", secretsKey, secretName+"-1"),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareCreateSecretRequestType(secretName+"-1", strings.ToLower(opaqueKey), "", "",
				[]types.KeyValue{{Key: secretDataKey, Value: secretDataValue}}),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(secretName+"-1", testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets", test.requestURI.namespace)
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

func TestUpdateSecret(t *testing.T) {
	testNamespaceName := secretNamespace + "-update"

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
				name:      secretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					secretNameKey: secretName,
					typeKey:       opaqueKey,
					namespaceKey:  testNamespaceName,
					dataKey:       []types.KeyValue{{Key: secretDataKey, Value: secretDataValue}},
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: secretDataKey, Value: secretDataValue}}),
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      secretName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: secretDataKey, Value: secretDataValue}}),
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      secretName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareSecretRequestType([]types.KeyValue{{Key: secretDataKey, Value: secretDataValue}}),
		},
		"ShouldFailWithBadRequestBody": {
			requestURI: requestURI{
				name:      secretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
					errorKey:   invalidRequest,
				},
			},
			requestData: map[string]interface{}{},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(secretName, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/secrets/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			request.Header.Set(contentType, applicationJson)

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

func TestDeleteSecret(t *testing.T) {
	testNamespaceName := secretNamespace + "-delete"

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
				name:      secretName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					messageKey: fmt.Sprintf("Deleted secret %q in namespace %q successfully", secretName, testNamespaceName)},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestURI: requestURI{
				name:      secretName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				name:      secretName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", secretsKey, secretName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestSecret(secretName, testNamespaceName)

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

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)

		})
	}
}
