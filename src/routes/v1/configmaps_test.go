package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	configMapName      = testutils.TestName + "-configmap"
	configKey          = "key"
	configValue        = "value"
	configmapsKey      = "configmaps"
	configMapNamespace = testutils.TestNamespace + configmapsKey
)

// createTestConfigMap creates a test ConfigMap object.
func createTestConfigMap(name, namespace string) {
	configMap := mocks.PrepareConfigMap(name, namespace, map[string]string{configKey: configValue})
	_, err := fakeClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
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
					testutils.Data: []types.KeyValue{{Key: configKey, Value: configValue}},
				},
			},
		},
		"ShouldHandleNotFoundConfigMap": {
			requestURI: requestURI{
				namespace: testNamespaceName,
				name:      configMapName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", configmapsKey, configMapName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				namespace: testNamespaceName + testutils.NonExistentSuffix,
				name:      configMapName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s %q not found", configmapsKey, configMapName),
					testutils.ErrorKey:   testutils.OperationFailed,
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
