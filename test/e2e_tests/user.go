package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate Secret routes and functionality", func() {
	var namespaceName, oneUserName, secondUserName string

	BeforeEach(func() {
		namespaceName = generateName(e2eNamespace)
		createTestNamespace(k8sClient, namespaceName)

		oneUserName = generateName("a-" + testUserName)
		createTestUser(k8sClient, oneUserName, namespaceName)

		secondUserName = generateName("b-" + testUserName)
		createTestUser(k8sClient, secondUserName, namespaceName)
	})

	Context("Validate get Users route", func() {
		It("Should get users with a managed label in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.UsersKey)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.CountKey: 2,
				testutils.UsersKey: []map[string]interface{}{
					{
						testutils.NameKey: oneUserName,
						testutils.RoleKey: testutils.AdminKey,
					},
					{
						testutils.NameKey: secondUserName,
						testutils.RoleKey: testutils.AdminKey,
					},
				},
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate get User route", func() {
		It("Should get an existing User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.NameKey: oneUserName,
				testutils.RoleKey: testutils.AdminKey,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)
		})

		It("Should handle getting a not found User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			user := mocks.PrepareRoleBinding(oneUserName+testutils.NonExistentSuffix, namespaceName, "")
			Expect(doesResourceExist(k8sClient, &user)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle getting a User in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.UsersKey, oneUserName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			user := mocks.PrepareRoleBinding(oneUserName, namespaceName+testutils.NonExistentSuffix, "")
			Expect(doesResourceExist(k8sClient, &user)).To(BeFalse())
			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate create User route", func() {
		It("Should create new User in a namespace", func() {
			newUserName := generateName(testUserName)

			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.UsersKey)
			requestData := mocks.PrepareUserType(newUserName, testutils.ViewerKey)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.NameKey: newUserName,
				testutils.RoleKey: testutils.ViewerKey,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			roleBinding := mocks.PrepareRoleBinding(newUserName, namespaceName, "")
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &roleBinding)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())

			Eventually(func() bool {
				roleBinding := getRoleBinding(k8sClient, newUserName, namespaceName)
				return roleBinding.RoleRef.Name == fmt.Sprintf("%s%s", testutils.CappUserPrefix, testutils.ViewerKey)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle creation of already-existing User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.UsersKey)
			requestData := mocks.PrepareUserType(oneUserName, testutils.ViewerKey)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q already exists", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName),
				testutils.ReasonKey: testutils.ReasonAlreadyExists,
			}

			Expect(status).Should(Equal(http.StatusConflict))
			compareResponses(response, expectedResponse)
		})

		It("Should handle creation of user with a non existent role", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s", platformURL, namespaceName, testutils.UsersKey)
			requestData := mocks.PrepareUserType(oneUserName, testutils.ViewerKey+testutils.NonExistentSuffix)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPost, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'User.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate update User route", func() {
		It("Should update a User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName)
			requestData := mocks.PrepareUpdateUserDataType(testutils.ViewerKey)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.NameKey: oneUserName,
				testutils.RoleKey: testutils.ViewerKey,
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			Eventually(func() bool {
				roleBinding := getRoleBinding(k8sClient, oneUserName, namespaceName)
				return roleBinding.RoleRef.Name == fmt.Sprintf("%s%s", testutils.CappUserPrefix, testutils.ViewerKey)
			}, testutils.Timeout, testutils.Interval).Should(BeTrue())
		})

		It("Should handle update of a not found User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName+testutils.NonExistentSuffix)
			requestData := mocks.PrepareUpdateUserDataType(testutils.ViewerKey)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle update of a User in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.UsersKey, oneUserName)
			requestData := mocks.PrepareUpdateUserDataType(testutils.ViewerKey)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle update of a User with a non-existent role", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName)
			requestData := mocks.PrepareUpdateUserDataType(testutils.ViewerKey + testutils.NonExistentSuffix)
			payload, err := json.Marshal(requestData)
			Expect(err).Should(Not(HaveOccurred()))

			status, response := performHTTPRequest(httpClient, bytes.NewBuffer(payload), http.MethodPut, uri, "", "", userToken)
			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  "Key: 'UpdateUserData.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
				testutils.ReasonKey: testutils.ReasonBadRequest,
			}

			Expect(status).Should(Equal(http.StatusBadRequest))
			compareResponses(response, expectedResponse)
		})
	})

	Context("Validate delete User route", func() {
		It("Should delete a User in an existing namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.MessageKey: fmt.Sprintf("Deleted roleBinding %q in namespace %q successfully", oneUserName, namespaceName),
			}

			Expect(status).Should(Equal(http.StatusOK))
			compareResponses(response, expectedResponse)

			roleBinding := mocks.PrepareRoleBinding(oneUserName, namespaceName, "")
			Eventually(func() bool {
				return doesResourceExist(k8sClient, &roleBinding)
			}, testutils.Timeout, testutils.Interval).Should(BeFalse())
		})

		It("Should handle deletion of a not found User in a namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName, testutils.UsersKey, oneUserName+testutils.NonExistentSuffix)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName+testutils.NonExistentSuffix),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})

		It("Should handle deletion of a User in a not found namespace", func() {
			uri := fmt.Sprintf("%s/v1/namespaces/%s/%s/%s", platformURL, namespaceName+testutils.NonExistentSuffix, testutils.UsersKey, oneUserName)
			status, response := performHTTPRequest(httpClient, nil, http.MethodDelete, uri, "", "", userToken)

			expectedResponse := map[string]interface{}{
				testutils.ErrorKey:  fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, oneUserName),
				testutils.ReasonKey: testutils.ReasonNotFound,
			}

			Expect(status).Should(Equal(http.StatusNotFound))
			compareResponses(response, expectedResponse)
		})
	})
})
