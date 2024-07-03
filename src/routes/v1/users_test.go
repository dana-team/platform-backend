package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	userNamespace := "test-namespace-user"
	type requestUri struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSuccessGettingUser": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"name": "test-user",
					"role": "admin",
				},
			},
		},
		"ShouldNotFoundUser": {
			requestUri: requestUri{
				username:  "test-not-exists",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q not found", "test-not-exists"),
					"error":   "Operation failed",
				},
			},
		},
		"ShouldNotFoundNamespace": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q not found", "test-user"),
					"error":   "Operation failed",
				},
			},
		},
	}

	setup()
	createTestNamespace(userNamespace)
	createRoleBinding(userNamespace, "test-user")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/%s/users/%s",
				test.requestUri.namespace, test.requestUri.username), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}

func TestGetUsers(t *testing.T) {
	usesNamespace := "test-namespace-users"
	type requestUri struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   types.UsersOutput
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldSuccessGettingUsers": {
			requestUri: requestUri{
				namespace: usesNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.UsersOutput{Count: 2, Users: []types.User{{Name: "test-user1", Role: "admin"}, {Name: "test-user2", Role: "admin"}}},
			},
		},
	}

	setup()
	createTestNamespace(usesNamespace)
	createRoleBinding(usesNamespace, "test-user1")
	createRoleBinding(usesNamespace, "test-user2")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", fmt.Sprintf("/v1/namespaces/%s/users/",
				test.requestUri.namespace), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response types.UsersOutput
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}

func TestUpdateUser(t *testing.T) {
	userNamespace := "test-update-user"
	type requestUri struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestUri  requestUri
		want        want
		requsetData map[string]string
	}{
		"ShouldSuccessUpdateUser": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"name": "test-user",
					"role": "viewer",
				},
			},
			requsetData: map[string]string{"role": "viewer", "name": "test-user"},
		},
		"ShouldNotFoundUser": {
			requestUri: requestUri{
				username:  "test-not-exists",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q not found", "test-not-exists"),
					"error":   "Operation failed",
				},
			},
			requsetData: map[string]string{"role": "viewer"},
		},
		"NotExistsRole": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]string{
					"details": "Key: 'PatchUserData.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					"error":   "Invalid request",
				},
			},
			requsetData: map[string]string{"role": "baladi"},
		},
	}

	setup()
	createTestNamespace(userNamespace)
	createRoleBinding(userNamespace, "test-user")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(test.requsetData)

			request, _ := http.NewRequest("PUT", fmt.Sprintf("/v1/namespaces/%s/users/%s",
				test.requestUri.namespace, test.requestUri.username), bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}

func TestCreateUser(t *testing.T) {
	userNamespace := "create-user-namespace"
	type requestUri struct {
		namespace string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestUri  requestUri
		want        want
		requsetData map[string]string
	}{
		"ShouldSuccessCreateUser": {
			requestUri: requestUri{
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"name": "test-user",
					"role": "viewer",
				},
			},
			requsetData: map[string]string{"role": "viewer", "name": "test-user"},
		},
		"AlreadyExists": {
			requestUri: requestUri{
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusConflict,
				response:   map[string]string{"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q already exists", "exists-user"), "error": "Operation failed"},
			},
			requsetData: map[string]string{"role": "viewer", "name": "exists-user"},
		},
		"NotExistsRole": {
			requestUri: requestUri{
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response: map[string]string{
					"details": "Key: 'User.Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
					"error":   "Invalid request",
				},
			},
			requsetData: map[string]string{"role": "baladi", "name": "test-user1"},
		},
	}

	setup()
	createTestNamespace(userNamespace)
	createRoleBinding(userNamespace, "exists-user")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			payload, _ := json.Marshal(test.requsetData)

			request, _ := http.NewRequest("POST", fmt.Sprintf("/v1/namespaces/%s/users/",
				test.requestUri.namespace), bytes.NewBuffer(payload))
			request.Header.Set("Content-Type", "application/json")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}

func TestDeleteUser(t *testing.T) {
	userNamespace := "delete-user-namespace"
	type requestUri struct {
		namespace string
		username  string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestUri requestUri
		want       want
	}{
		"ShouldDeleteGettingUser": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response: map[string]string{
					"message": fmt.Sprintf("deleted roleBinding %q successfully", "test-user"),
				},
			},
		},
		"ShouldNotFoundUser": {
			requestUri: requestUri{
				username:  "test-not-exists",
				namespace: userNamespace,
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q not found", "test-not-exists"),
					"error":   "Operation failed",
				},
			},
		},
		"ShouldNotFoundNamespace": {
			requestUri: requestUri{
				username:  "test-user",
				namespace: "test-not-exists",
			},
			want: want{
				statusCode: http.StatusNotFound,
				response: map[string]string{
					"details": fmt.Sprintf("rolebindings.rbac.authorization.k8s.io %q not found", "test-user"),
					"error":   "Operation failed",
				},
			},
		},
	}

	setup()
	createTestNamespace(userNamespace)
	createRoleBinding(userNamespace, "test-user")

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/namespaces/%s/users/%s",
				test.requestUri.namespace, test.requestUri.username), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			var response map[string]string
			if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
				panic(err)
			}
			assert.Equal(t, test.want.response, response)

		})
	}
}
