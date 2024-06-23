package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigMaps(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/configmaps/test-configmap", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.ConfigMap
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Data), 1)
}
