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
		Type:       "opaque",
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
	request, _ := http.NewRequest("GET", "/v1/namespaces/default/secrets", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.SecretsList
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Secrets), 1)
	assert.GreaterOrEqual(t, response.Count, 1)
	assert.NotNil(t, response.ContinueToken)
	assert.NotNil(t, response.RemainingCount)
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

func TestUpdateSecret(t *testing.T) {
	updateRequest := types.UpdateSecretRequest{
		Data: []types.KeyValue{{Key: "key2", Value: "ZmFrZQ=="}},
	}
	body, _ := json.Marshal(updateRequest)
	request, _ := http.NewRequest("PUT", "/v1/namespaces/default/secrets/test-secret", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.UpdateSecretResponse
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
