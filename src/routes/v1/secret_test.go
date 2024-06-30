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
		"ShouldSuccessGettingSecret": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					"secretName": secretName,
					"type":       "Opaque",
					"id":         "",
					"data":       []interface{}{map[string]interface{}{"key": "key1", "value": "fake"}},
				},
			},
		},
		"ShouldNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					"details": fmt.Sprintf("secrets %q not found", "test-not-exists"),
					"error":   "Operation failed",
				},
			},
		},
		"ShouldNotFoundNamespace": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					"details": fmt.Sprintf("secrets %q not found", secretName),
					"error":   "Operation failed",
				},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret(secretName, secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

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
		response   types.GetSecretsResponse
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSuccessGettingSecrets": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: types.GetSecretsResponse{Count: 2, Secrets: []types.Secret{{SecretName: "test-secret1", NamespaceName: secretNamespace, Type: "Opaque"},
					{SecretName: "test-secret2", NamespaceName: secretNamespace, Type: "Opaque"}}},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret("test-secret1", secretNamespace)
	createTestSecret("test-secret2", secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/%s/secrets/",
				test.requestUri.namespace), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response types.GetSecretsResponse
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}

func TestDeleteSecrets(t *testing.T) {
	secretNamespace := "test-namespace-deleteSecret"
	secretName := "test-secret"
	type requestUri struct {
		namespace string
		secret    string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSuccessGettingSecret": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"message": fmt.Sprintf("Secret %q was deleted successfully", secretName)},
			},
		},
		"ShouldNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("secrets %q not found", "test-not-exists"),
					"error":   "Operation failed",
				},
			},
		},
		"ShouldNotFoundNamespace": {
			requestUri: requestUri{
				secret:    secretName,
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("secrets %q not found", secretName),
					"error":   "Operation failed",
				},
			},
		},
	}
	createTestNamespace(secretNamespace)
	createTestSecret(secretName, secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

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
		"ShouldSuccessCreateSecret": {
			requestUri: requestUri{
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"secretName":    "test-secret1",
					"type":          "Opaque",
					"namespaceName": secretNamespace,
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
					"details": "data is required for Opaque secrets",
					"error":   "Operation failed",
				},
			},
			requsetData: map[string]interface{}{
				"secretName": "secret-test",
				"type":       "opaque",
				"dataex":     []interface{}{map[string]interface{}{"key": "key1", "value": "fake"}},
			},
		},
	}

	createTestNamespace(secretNamespace)
	createTestSecret("test-secret", secretNamespace)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(test.requsetData)

			request, _ := http.NewRequest("POST", fmt.Sprintf("/v1/namespaces/%s/secrets/",
				test.requestUri.namespace), bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)
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
					"secretName":    "test-secret",
					"type":          "Opaque",
					"namespaceName": secretNamespace,
					"data":          []interface{}{map[string]interface{}{"key": "key1", "value": "value"}},
				},
			},
			requsetData: types.UpdateSecretRequest{Data: []types.KeyValue{{Key: "key1", Value: "value"}}},
		},
		"ShouldNotFoundSecret": {
			requestUri: requestUri{
				secret:    "test-not-exists",
				namespace: secretNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					"details": fmt.Sprintf("secrets %q not found", "test-not-exists"),
					"error":   "Operation failed",
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
					"details": "Key: 'UpdateSecretRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
					"error":   "Invalid request",
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

			request, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/namespaces/%s/secrets/%s",
				test.requestUri.namespace, test.requestUri.secret), bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}
