package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils/testutils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	userNamespace = testutils.TestNamespace + "-" + testutils.UsersKey
	userName      = testutils.TestName + "-user"
)

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
					testutils.CountKey: 2,
					testutils.UsersKey: []types.User{
						{Name: userName + "-1", Role: testutils.AdminKey},
						{Name: userName + "-2", Role: testutils.AdminKey},
					},
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestRoleBinding(fakeClient, userName+"-1", testNamespaceName, testutils.AdminKey)
	mocks.CreateTestRoleBinding(fakeClient, userName+"-2", testNamespaceName, testutils.AdminKey)

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
					testutils.NameKey: userName,
					testutils.RoleKey: testutils.AdminKey,
				},
			},
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestRoleBinding(fakeClient, userName, testNamespaceName, testutils.AdminKey)

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
					testutils.NameKey: userName,
					testutils.RoleKey: testutils.ViewerKey,
				},
			},
			requestData: mocks.PrepareUserType(userName, testutils.ViewerKey),
		},
		"ShouldHandleAlreadyExists": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusConflict,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q already exists", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName+"-1"),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareUserType(userName+"-1", testutils.ViewerKey),
		},
		"ShouldHandleNonExistentRole": {
			requestURI: requestURI{
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.DetailsKey: "Key: 'User.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					testutils.ErrorKey:   testutils.InvalidRequest,
				},
			},
			requestData: mocks.PrepareUserType(userName, testutils.ViewerKey+testutils.NonExistentSuffix),
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestRoleBinding(fakeClient, userName+"-1", testNamespaceName, testutils.ViewerKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/users", test.requestURI.namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

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
					testutils.NameKey: userName,
					testutils.RoleKey: testutils.ViewerKey,
				},
			},
			requestData: mocks.PrepareUserType(userName, testutils.ViewerKey),
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
			requestData: mocks.PrepareUpdateUserDataType(testutils.ViewerKey),
		},
		"ShouldHandleNotExistentRole": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]interface{}{
					testutils.DetailsKey: "Key: 'UpdateUserData.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					testutils.ErrorKey:   testutils.InvalidRequest,
				},
			},
			requestData: mocks.PrepareUpdateUserDataType(testutils.ViewerKey + testutils.NonExistentSuffix),
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestRoleBinding(fakeClient, userName, testNamespaceName, testutils.AdminKey)

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(test.requestData)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/users/%s", test.requestURI.namespace, test.requestURI.username)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(payload))
			assert.NoError(t, err)
			request.Header.Set(testutils.ContentType, testutils.ApplicationJson)

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
					testutils.MessageKey: fmt.Sprintf("Deleted roleBinding %q in namespace %q successfully", userName, testNamespaceName),
				},
			},
		},
		"ShouldHandleNotFoundUser": {
			requestURI: requestURI{
				username:  userName + testutils.NonExistentSuffix,
				namespace: testNamespaceName,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName+testutils.NonExistentSuffix),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
		"ShouldHandleNotFoundNamespace": {
			requestURI: requestURI{
				username:  userName,
				namespace: testNamespaceName + testutils.NonExistentSuffix,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]interface{}{
					testutils.DetailsKey: fmt.Sprintf("%s.%s %q not found", testutils.RoleBindingsKey, testutils.RoleBindingsGroupKey, userName),
					testutils.ErrorKey:   testutils.OperationFailed,
				},
			},
		},
	}

	setup()
	mocks.CreateTestNamespace(fakeClient, testNamespaceName)
	mocks.CreateTestRoleBinding(fakeClient, userName, testNamespaceName, testutils.AdminKey)

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
