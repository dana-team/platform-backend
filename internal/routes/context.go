package routes

import (
	"context"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	multicluster "github.com/oam-dev/cluster-gateway/pkg/apis/cluster/transport"
)

// GetContext returns the context from the gin.Context, containing the managed cluster if available.
func GetContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	cluster, exists := middleware.GetCluster(c)
	if !exists {
		return ctx
	}

	return multicluster.WithMultiClusterContext(ctx, cluster)
}
