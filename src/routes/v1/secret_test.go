package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	routev1 "github.com/dana-team/platform-backend/src/routes/v1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	router *gin.Engine
	client *fake.Clientset
)

func TestMain(m *testing.M) {
	client = fake.NewSimpleClientset()
	router = setupRouter()

	createTestSecret()

	m.Run()
}

func TestCreateSecret(t *testing.T) {
	secretRequest := types.CreateSecretRequest{
		Type: "Opaque",
		Name: "new-secret",
		Data: []types.KeyValue{{Key: "key1", Value: "value1"}},
	}
	body, _ := json.Marshal(secretRequest)
	request, _ := http.NewRequest("POST", "/v1/namespaces/default/secrets", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusOK, writer.Code)
	var response types.CreateSecretResponse
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "new-secret", response.Name)
	assert.Equal(t, "default", response.Namespace)
}

func TestGetSecrets(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/namespaces/default/secrets", nil)
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
	assert.Equal(t, "test-secret", response.Name)
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
	assert.Equal(t, "test-secret", response.Name)
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

func setupRouter() *gin.Engine {
	engine := gin.Default()

	logger, _ := zap.NewDevelopment()

	engine.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Set("kubeClient", client)
		c.Next()
	})

	v1 := engine.Group("/v1")
	{
		namespacesGroup := v1.Group("/namespaces")
		{
			namespacesGroup.GET("/:namespace/secrets", routev1.GetSecrets())
			namespacesGroup.POST("/:namespace/secrets", routev1.CreateSecret())
			namespacesGroup.GET("/:namespace/secrets/:name", routev1.GetSecret())
			namespacesGroup.PATCH("/:namespace/secrets/:name", routev1.PatchSecret())
			namespacesGroup.DELETE("/:namespace/secrets/:name", routev1.DeleteSecret())
		}
	}
	return engine
}

func createTestSecret() {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key1": []byte("ZmFrZQ=="),
		},
	}
	_, err := client.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
