package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestGetContainerAppRevisions(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/capprevisions/", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.CappRevisionList
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.CappRevisions), 1)
	assert.GreaterOrEqual(t, response.Count, 1)
}

func TestGetContainerAppRevision(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/capprevisions/test-capprevision", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.CappRevision
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-capprevision", response.Metadata.Name)
	assert.Equal(t, "test-namespace", response.Metadata.Namespace)
}
