package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"

	"testing"
)

func Test_PaginationMiddleware(t *testing.T) {
	logger, _ := zap.NewProduction()
	router := gin.New()

	// Set up a middleware that injects the logger into the context
	router.Use(func(c *gin.Context) {
		c.Set(LoggerCtxKey, logger)
		c.Next()
	})

	router.Use(PaginationMiddleware()).GET("/ping", func(c *gin.Context) {
		_, exists := c.Get(LimitCtxKey)
		if !exists {
			t.Error("Expected limit to be set in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "limit not set in context"})
			return
		}

		_, exists = c.Get(PageCtxKey)
		if !exists {
			t.Error("Expected page to be set in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "page not set in context"})
			return
		}

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
