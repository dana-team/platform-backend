package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	secretNameKey = "secretName"
	dataKey       = "data"
	idKey         = "id"
	typeKey       = "type"
	namespaceKey  = "namespaceName"
	secretsKey    = "secrets"
)

func TestGetSecret(t *testing.T) {
	secretNamespace := "test-namespace-secret"
	secretName := "test-secret"
	type requestUri struct {
		namespace string
		secret    string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSucceedGettingSecret": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					secretNameKey: secretName,
					typeKey:       "Opaque",
					idKey:         "",
					dataKey:       []interface{}{map[string]interface{}{"key": "key1", "value": "fake"}},
				},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("secrets %q not found", "test-not-exists"),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("secrets %q not found", secretName),
					errorKey:   operationFailed,
				},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret(secretName, secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), nil)
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

func TestGetSecrets(t *testing.T) {
	secretNamespace := "test-namespace-secrets"
	type requestUri struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSucceedGettingSecrets": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count: 2,
					secretsKey: []types.Secret{
						{SecretName: "test-secret1", NamespaceName: secretNamespace, Type: "Opaque"},
						{SecretName: "test-secret2", NamespaceName: secretNamespace, Type: "Opaque"}},
				},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret("test-secret1", secretNamespace)
	createTestSecret("test-secret2", secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/namespaces/%s/secrets/",
				test.requestUri.namespace), nil)
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
	secretNamespace := "test-namespace-deleteSecret"
	secretName := "test-secret"
	type requestUri struct {
		namespace string
		secret    string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}
	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSucceedGettingSecret": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					message: fmt.Sprintf("Secret %q was deleted successfully", secretName)},
			},
		},
		"ShouldHandleNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("secrets %q not found", "test-not-exists"),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("secrets %q not found", secretName),
					errorKey:   operationFailed,
				},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret(secretName, secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), nil)
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

func TestCreateSecret(t *testing.T) {
	secretNamespace := "create-secret-namespace"
	type requestUri struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}
	cases := map[string]struct {
		requestUri  requestUri
		want        want
		requsetData interface{}
	}{
		"ShouldSucceedCreateSecret": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					secretNameKey: "test-secret1",
					typeKey:       "Opaque",
					namespaceKey:  secretNamespace,
				},
			},
			requsetData: types.CreateSecretRequest{Type: "opaque", Data: []types.KeyValue{{Key: "key1", Value: "value"}}, SecretName: "test-secret1"},
		},
		"AlreadyExists": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusConflict,
				response:   map[string]string{"details": fmt.Sprintf("secrets %q already exists", "test-secret"), "error": "Operation failed"},
			},
			requsetData: types.CreateSecretRequest{Type: "opaque", Data: []types.KeyValue{{Key: "key1", Value: "value"}}, SecretName: "test-secret"},
		},
		"DidNotProvideData": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]string{
					detailsKey: "data is required for Opaque secrets",
					errorKey:   operationFailed,
				},
			},
			requsetData: map[string]interface{}{
				secretNameKey: "secret-test",
				typeKey:       "opaque",
				"dataex":      []interface{}{map[string]interface{}{"key": "key1", "value": "fake"}},
			},
		},
	}

	createTestNamespace(secretNamespace)
	createTestSecret("test-secret", secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(test.requsetData)

			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/namespaces/%s/secrets/",
				test.requestUri.namespace), bytes.NewBuffer(payload))
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

func TestUpdateSecret(t *testing.T) {
	secretNamespace := "test-update-secret"
	type requestUri struct {
		namespace string
		secret    string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}
	cases := map[string]struct {
		requestUri  requestUri
		want        want
		requsetData interface{}
	}{
		"ShouldSucceedUpdateSecret": {
			requestUri: requestUri{
				secret:    "test-secret",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					secretNameKey: "test-secret",
					typeKey:       "Opaque",
					namespaceKey:  secretNamespace,
					dataKey:       []interface{}{map[string]interface{}{"key": "key1", "value": "value"}},
				},
			},
			requsetData: types.UpdateSecretRequest{Data: []types.KeyValue{{Key: "key1", Value: "value"}}},
		},
		"ShouldHandleNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("secrets %q not found", "test-not-exists"),
					errorKey:   operationFailed,
				},
			},
			requsetData: types.UpdateSecretRequest{Data: []types.KeyValue{{Key: "key1", Value: "value"}}},
		},
		"InvalidData": {
			requestUri: requestUri{
				secret:    "test-secret",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
					errorKey:   invalidRequest,
				},
			},
			requsetData: map[string]interface{}{
				"datasd": []interface{}{map[string]interface{}{"key": "key1", "value": "value"}},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret("test-secret", secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(test.requsetData)

			request, err := http.NewRequest("PUT", fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), bytes.NewBuffer(payload))
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
