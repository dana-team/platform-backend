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

func TestListNamespaces(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.NamespaceList
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Namespaces), 1)
	assert.GreaterOrEqual(t, response.Count, 1)
	assert.NotNil(t, response.ContinueToken)
	assert.NotNil(t, response.RemainingCount)
}

func TestGetNamespace(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Namespace
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-namespace", response.Name)
}

func TestCreateNamespace(t *testing.T) {
	namespaceRequest := types.Namespace{
		Name: "new-namespace",
	}
	body, _ := json.Marshal(namespaceRequest)
	request, _ := http.NewRequest("POST", "/v1/namespaces/", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Namespace
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new-namespace", response.Name)
}

func TestDeleteNamespace(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/v1/namespaces/test-namespace", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response map[string]string
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deleted namespace successfully test-namespace", response["message"])
}
