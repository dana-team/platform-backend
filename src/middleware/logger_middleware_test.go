package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestLoggerMiddleware(t *testing.T) {
	logger, _ := zap.NewProduction()

	router := gin.New()
	router.Use(LoggerMiddleware(logger))
	router.GET("/ping", func(c *gin.Context) {
		ctxLogger, exists := c.Get("logger")
		if !exists {
			t.Error("Expected logger to be set in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "logger not set in context"})
			return
		}

		ctxLogger.(*zap.Logger).Info("Ping endpoint hit")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	type args struct {
		path string
	}
	type want struct {
		expectedStatus int
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccess": {
			args: args{
				path: "/ping",
			},
			want: want{
				expectedStatus: http.StatusOK,
			},
		},
	}

	// Run test cases
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.args.path, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.want.expectedStatus {
				t.Errorf("Expected status code %d; got %d", tc.want.expectedStatus, w.Code)
			}
		})
	}
}
