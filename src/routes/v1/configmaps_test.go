package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	configMapName      = testName + "-configmap"
	configKey          = "key"
	configValue        = "value"
	configmapsKey      = "configmaps"
	configMapNamespace = testNamespace + configmapsKey
)

// createTestConfigMap creates a test ConfigMap object.
func createTestConfigMap(name, namespace string) {
	configMap := mocks.PrepareConfigMap(name, namespace, map[string]string{configKey: configValue})
	_, err := client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetConfigMap(t *testing.T) {
	testNamespaceName := configMapNamespace + "-get"

	type requestURI struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingConfigMap": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      configMapName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					data: []types.KeyValue{{Key: configKey, Value: configValue}},
				},
			},
		},
		"ShouldHandleNotFoundConfigMap": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      configMapName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", configmapsKey, configMapName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + nonExistentSuffix,
				name:      configMapName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", configmapsKey, configMapName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestConfigMap(configMapName, testNamespaceName)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/configmaps/%s", test.requestURI.namespace, test.requestURI.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, err := json.Marshal(test.want.response)
			assert.NoError(t, err)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
