package e2e_tests

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	"github.com/dana-team/platform-backend/src/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Validate ConfigMap routes and functionality", func() {
	var namespaceName, configMapName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		configMapName = generateName(testConfigMapName)
		createTestConfigMap(k8sClient, configMapName, namespaceName, map[string]string{testutils.ConfigMapDataKey: testutils.ConfigMapDataValue})
	})

	Context("Validate get ConfigMap route", func() {
		It("Should get a specific ConfigMap in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ConfigmapsKey, configMapName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DataKey: []types.KeyValue{{Key: testutils.ConfigMapDataKey, Value: testutils.ConfigMapDataValue}},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a not found ConfigMap in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.ConfigmapsKey, configMapName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.ConfigmapsKey, configMapName+testutils.NonExistentSuffix),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			configMap := mocks.PrepareConfigMap(configMapName+testutils.NonExistentSuffix, namespaceName, nil)
			Expect(doesResourceExist(k8sClient, &configMap)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle a get of a ConfigMap in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.ConfigmapsKey, configMapName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.DetailsKey: fmt.Sprintf("%s %q not found", testutils.ConfigmapsKey, configMapName),
				testutils.ErrorKey:   testutils.OperationFailed,
			}

			configMap := mocks.PrepareConfigMap(configMapName, namespaceName+testutils.NonExistentSuffix, nil)
			Expect(doesResourceExist(k8sClient, &configMap)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
