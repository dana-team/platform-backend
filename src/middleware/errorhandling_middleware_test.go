package middleware

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// Mock custom error type that implements customerrors.ErrorWithStatusCode
type MockCustomError struct {
	msg        string
	statusCode int
	reason     string
}

func (e *MockCustomError) Error() string {
	return e.msg
}

func (e *MockCustomError) StatusCode() int {
	return e.statusCode
}

func (e *MockCustomError) StatusReason() string {
	return e.reason
}

func Test_ErrorHandlingMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandlingMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/custom-error", func(c *gin.Context) {
		_ = c.Error(&MockCustomError{
			msg:        "custom error message",
			statusCode: http.StatusBadRequest,
			reason:     "CustomReason",
		})
	})

	router.GET("/k8s-error", func(c *gin.Context) {
		_ = c.Error(&k8serrors.StatusError{
			ErrStatus: metav1.Status{
				Code:    http.StatusNotFound,
				Message: "k8s error message",
				Reason:  metav1.StatusReasonNotFound,
			},
		})
	})

	router.GET("/unknown-error", func(c *gin.Context) {
		_ = c.Error(errors.New("unknown error message"))
	})

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
		"ShouldReturnSuccess": {
			args: args{path: "/ping"},
			want: want{
				expectedStatus: http.StatusOK,
				expectedBody:   `{"message":"pong"}`,
			},
		},
		"ShouldReturnCustomError": {
			args: args{path: "/custom-error"},
			want: want{
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   `{"error":"custom error message","reason":"Unknown"}`,
			},
		},
		"ShouldReturnK8sError": {
			args: args{path: "/k8s-error"},
			want: want{
				expectedStatus: http.StatusNotFound,
				expectedBody:   `{"error":"k8s error message","reason":"NotFound"}`,
			},
		},
		"ShouldReturnUnknownError": {
			args: args{path: "/unknown-error"},
			want: want{
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   `{"error":"unknown error message","reason":"Unknown"}`,
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

			if w.Body.String() != tc.want.expectedBody {
				t.Errorf("Expected body %s; got %s", tc.want.expectedBody, w.Body.String())
			}
		})
	}
}
