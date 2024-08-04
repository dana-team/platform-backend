package mocks

import (
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// GinContext creates a new *gin.Context for testing purposes
func GinContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{}
	return c
}

func SetPaginationValues(c *gin.Context, limit, page int) {
	c.Set(middleware.LimitCtxKey, limit)
	c.Set(middleware.PageCtxKey, page)
}
