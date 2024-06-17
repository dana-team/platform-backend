package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	secretRequest := types.CreateSecretRequest{
		Type:       "Opaque",
		SecretName: "new-secret",
		Data:       []types.KeyValue{{Key: "key1", Value: "ZmFrZQ=="}},
	}
	body, _ := json.Marshal(secretRequest)
	request, _ := http.NewRequest("POST", "/v1/namespaces/default/secrets/", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.CreateSecretResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new-secret", response.SecretName)
	assert.Equal(t, "default", response.NamespaceName)
}

func TestGetSecrets(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/default/secrets/", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.GetSecretsResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestGetSpecificSecret(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/default/secrets/test-secret", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.GetSecretResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-secret", response.SecretName)
}

func TestPatchSecret(t *testing.T) {
	patchRequest := types.PatchSecretRequest{
		Data: []types.KeyValue{{Key: "key2", Value: "ZmFrZQ=="}},
	}
	body, _ := json.Marshal(patchRequest)
	request, _ := http.NewRequest("PATCH", "/v1/namespaces/default/secrets/test-secret", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.PatchSecretResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-secret", response.SecretName)
}

func TestDeleteSecret(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/v1/namespaces/default/secrets/test-secret", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.DeleteSecretResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Secret \"test-secret\" was deleted successfully", response.Message)
}
