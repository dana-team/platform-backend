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
	configMapNamespace = testName + "-configmap-ns"
	configMapName      = testName + "-configmap"
	configKey          = "key"
	configValue        = "value"
	configmaps         = "configmaps"
)

func setupConfigMaps() {
	createTestNamespace(configMapNamespace)
	createTestConfigMap(configMapName+"-1", configMapNamespace)
}

// createTestConfigMap creates a test ConfigMap object.
func createTestConfigMap(name, namespace string) {
	configMap := mocks.PrepareConfigMap(name, namespace, map[string]string{configKey + "-1": configValue + "-1"})
	_, err := client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetConfigMap(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingConfigMap": {
			requestParams: requestParams{
				namespace: configMapNamespace,
				name:      configMapName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					data: []types.KeyValue{{Key: configKey + "-1", Value: configValue + "-1"}},
				},
			},
		},
		"ShouldHandleNotFoundConfigMap": {
			requestParams: requestParams{
				namespace: configMapNamespace,
				name:      configMapName + "-1" + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", configmaps, configMapName+"-1"+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestParams: requestParams{
				namespace: configMapNamespace + nonExistentSuffix,
				name:      configMapName + "-1",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s %q not found", configmaps, configMapName+"-1"),
					errorKey:   operationFailed,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/configmaps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]interface{}
			err = json.Unmarshal(writer.Body.Bytes(), &response)
			assert.NoError(t, err)

			wantResponseJSON, _ := json.Marshal(test.want.response)
			var wantResponseNormalized map[string]interface{}
			err = json.Unmarshal(wantResponseJSON, &wantResponseNormalized)
			assert.NoError(t, err)
			assert.Equal(t, wantResponseNormalized, response)
		})
	}
}
