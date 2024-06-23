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

func TestCreateServiceAccount(t *testing.T) {
	serviceAccount := types.ServiceAccount{
		ServiceAccountName: "new-sa",
	}
	body, _ := json.Marshal(serviceAccount)
	request, _ := http.NewRequest("POST", "/v1/namespaces/default/serviceaccounts/", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.ServiceAccount
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new-sa", response.ServiceAccountName)
}

func TestGetToken(t *testing.T) {
	serviceAccount := types.ServiceAccount{
		ServiceAccountName: "test-sa",
	}
	body, _ := json.Marshal(serviceAccount)
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/token/", bytes.NewBuffer(body))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Token
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
}
