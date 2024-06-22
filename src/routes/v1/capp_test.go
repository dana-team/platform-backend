package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateCapp(t *testing.T) {
	cappRequest := types.Capp{
		Metadata: types.Metadata{
			Name:      "new-capp",
			Namespace: "test-namespace",
		},
		Annotations: []types.KeyValue{
			{Key: "annotation1", Value: "value1"},
		},
		Labels: []types.KeyValue{
			{Key: "label1", Value: "value1"},
		},
		Spec:   cappv1alpha1.CappSpec{},
		Status: cappv1alpha1.CappStatus{},
	}
	body, _ := json.Marshal(cappRequest)
	request, _ := http.NewRequest("POST", "/v1/namespaces/test-namespace/capps/", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Capp
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new-capp", response.Metadata.Name)
	assert.Equal(t, "test-namespace", response.Metadata.Namespace)
}

func TestGetCapps(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/capps/", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.CappList
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Capps), 1)
	assert.GreaterOrEqual(t, response.Count, 1)
}

func TestGetCapp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/test-namespace/capps/test-capp", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Capp
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-capp", response.Metadata.Name)
	assert.Equal(t, "test-namespace", response.Metadata.Namespace)
}

func TestUpdateCapp(t *testing.T) {
	cappRequest := types.Capp{
		Metadata: types.Metadata{
			Name:      "test-capp",
			Namespace: "test-namespace",
		},
		Annotations: []types.KeyValue{
			{Key: "annotation1", Value: "updated-value"},
		},
		Labels: []types.KeyValue{
			{Key: "label1", Value: "updated-value"},
		},
		Spec:   cappv1alpha1.CappSpec{},
		Status: cappv1alpha1.CappStatus{},
	}
	body, _ := json.Marshal(cappRequest)
	request, _ := http.NewRequest("PUT", "/v1/namespaces/test-namespace/capps/test-capp", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.Capp
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-capp", response.Metadata.Name)
	assert.Equal(t, "test-namespace", response.Metadata.Namespace)
}

func TestDeleteCapp(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/v1/namespaces/test-namespace/capps/test-capp", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response map[string]string
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Deleted capp \"test-capp\" in namespace \"test-namespace\" successfully", response["message"])
}
