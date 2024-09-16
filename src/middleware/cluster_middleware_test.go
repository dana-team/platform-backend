package middleware

import (
	"errors"
	"fmt"

	"github.com/dana-team/platform-backend/src/customerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	"testing"
)

// errorMiddlewareForTest is a middleware for tests that captures errors and returns a response if an error occurs, ignoring `NotFoundError`.
func errorMiddlewareForTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !hasErrors(c) {
			return
		}

		lastError := c.Errors.Last()
		var notFoundError *customerrors.NotFoundError

		// If either clusterName was received, or both cappName and namespaceName were received,
		// and Capp was not found, return valid response to avoid fakeClient setup.
		if errors.As(lastError.Err, &notFoundError) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			errorKey: lastError.Err.Error(),
			reason:   notFoundError.StatusReason(),
		})
		c.Next()
	}
}

func Test_ClusterMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(errorMiddlewareForTest()) // Use the error middleware for test scenarios

	router.Use(ClusterMiddleware()).GET("/namespaces/:namespaceName/capps/:cappName", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	var notFoundError = customerrors.NewNotFoundError("error")

	type args struct {
		path string
	}
	type want struct {
		expectedStatus int
		expectedBody   string
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccess": {
			args: args{
				path: "/namespaces/namespace/capps/capp",
			},
			want: want{
				expectedStatus: http.StatusOK,
				expectedBody:   `{"message":"pong"}`,
			},
		},
		"ShouldFailOnNoNamespace": {
			args: args{
				path: "/namespaces//capps/capp",
			},
			want: want{
				expectedStatus: http.StatusBadRequest,
				expectedBody:   fmt.Sprintf(`{"error":%q,"reason":%q}`, errNoClusterOrNameAndNamespaceProvided, notFoundError.StatusReason()),
			},
		},
		"ShouldFailOnNoName": {
			args: args{
				path: "/namespaces/namespace/capps/",
			},
			want: want{
				expectedStatus: http.StatusBadRequest,
				expectedBody:   fmt.Sprintf(`{"error":%q,"reason":%q}`, errNoClusterOrNameAndNamespaceProvided, notFoundError.StatusReason()),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.args.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.want.expectedStatus {
				t.Errorf("Expected status code %d; got %d", tc.want.expectedStatus, w.Code)
			}

			actualBody := w.Body.String()
			if actualBody != tc.want.expectedBody {
				t.Errorf("Expected body %s; got %s", tc.want.expectedBody, actualBody)
			}
		})
	}
}
