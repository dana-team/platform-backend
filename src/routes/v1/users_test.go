package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	roleBindingsKey      = "rolebindings"
	usersKey             = "users"
	userNamespace        = testNamespace + "-" + usersKey
	roleBindingsGroupKey = "rbac.authorization.k8s.io"
	userName             = testName + "-user"
	roleKey              = "role"
	adminKey             = "admin"
	viewerKey            = "viewer"
)

// createTestRoleBinding creates a test RoleBinding object.
func createTestRoleBinding(name, namespace, role string) {
	roleBinding := mocks.PrepareRoleBinding(name, namespace, role)

	_, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetUsers(t *testing.T) {
	testNamespaceName := userNamespace + "-get"
	type requestURI struct {
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
		"ShouldSucceedGettingUsers": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					count: 2,
					usersKey: []types.User{
						{Name: userName + "-1", Role: adminKey},
						{Name: userName + "-2", Role: adminKey},
					},
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestRoleBinding(userName+"-1", testNamespaceName, adminKey)
	createTestRoleBinding(userName+"-2", testNamespaceName, adminKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/users", test.requestURI.namespace)
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

func TestGetUser(t *testing.T) {
	testNamespaceName := userNamespace + "-get-one"

	type requestURI struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedGettingUser": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					nameKey: userName,
					roleKey: adminKey,
				},
			},
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", roleBindingsKey, roleBindingsGroupKey, userName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", roleBindingsKey, roleBindingsGroupKey, userName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestRoleBinding(userName, testNamespaceName, adminKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/users/%s", test.requestURI.namespace, test.requestURI.username)
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

func TestCreateUser(t *testing.T) {
	testNamespaceName := userNamespace + "-create"

	type requestURI struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedCreateUser": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					nameKey: userName,
					roleKey: viewerKey,
				},
			},
			requestData: mocks.PrepareUserType(userName, viewerKey),
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q already exists", roleBindingsKey, roleBindingsGroupKey, userName+"-1"),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareUserType(userName+"-1", viewerKey),
		},
		"ShouldHandleNonExistentRole": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'User.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					errorKey:   invalidRequest,
				},
			},
			requestData: mocks.PrepareUserType(userName, viewerKey+nonExistentSuffix),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestRoleBinding(userName+"-1", testNamespaceName, viewerKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/users", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

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

func TestUpdateUser(t *testing.T) {
	testNamespaceName := userNamespace + "-update"

	type requestURI struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI  requestURI
		want        want
		requestData interface{}
	}{
		"ShouldSucceedUpdateUser": {
			requestURI: requestURI{
				username:  "test-user",
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					nameKey: userName,
					roleKey: viewerKey,
				},
			},
			requestData: mocks.PrepareUserType(userName, viewerKey),
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", roleBindingsKey, roleBindingsGroupKey, userName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
			requestData: mocks.PrepareUpdateUserDataType(viewerKey),
		},
		"ShouldHandleNotExistentRole": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					detailsKey: "Key: 'UpdateUserData.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					errorKey:   invalidRequest,
				},
			},
			requestData: mocks.PrepareUpdateUserDataType(viewerKey + nonExistentSuffix),
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestRoleBinding(userName, testNamespaceName, adminKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/users/%s", test.requestURI.namespace, test.requestURI.username)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

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

func TestDeleteUser(t *testing.T) {
	testNamespaceName := userNamespace + "-delete"

	type requestURI struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]interface{}
	}

	cases := map[string]struct {
		requestURI requestURI
		want       want
	}{
		"ShouldSucceedDeletingUser": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]interface{}{
					messageKey: fmt.Sprintf("Deleted roleBinding %q in namespace %q successfully", userName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + nonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", roleBindingsKey, roleBindingsGroupKey, userName+nonExistentSuffix),
					errorKey:   operationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName + nonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					detailsKey: fmt.Sprintf("%s.%s %q not found", roleBindingsKey, roleBindingsGroupKey, userName),
					errorKey:   operationFailed,
				},
			},
		},
	}

	setup()
	createTestNamespace(testNamespaceName)
	createTestRoleBinding(userName, testNamespaceName, adminKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/users/%s", test.requestURI.namespace, test.requestURI.username)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
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
