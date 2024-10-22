package v1

import (
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/internal/controllers"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

const (
	configMapName      = testutils.TestName + "-configmap"
	configMapNamespace = testutils.TestNamespace + testutils.ConfigmapsKey
)

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
					testutils.DataKey: []types.KeyValue{{Key: testutils.ConfigMapDataKey, Value: testutils.ConfigMapDataValue}},
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
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetConfigMap, configMapName+testutils.NonExistentSuffix),
						fmt.Sprintf("%s %q not found", testutils.ConfigmapsKey, configMapName+testutils.NonExistentSuffix)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
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
					testutils.ErrorKey: fmt.Sprintf("%v, %v",
						fmt.Sprintf(controllers.ErrCouldNotGetConfigMap, configMapName),
						fmt.Sprintf("%s %q not found", testutils.ConfigmapsKey, configMapName)),
					testutils.ReasonKey: metav1.StatusReasonNotFound,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestConfigMap(fakeClient, configMapName, testNamespaceName)

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
