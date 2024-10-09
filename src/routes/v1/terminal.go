package v1

import (
	"github.com/dana-team/platform-backend/src/controllers"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartTerminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		var startTerminalUri types.StartTerminalUri
		if err := c.BindUri(&startTerminalUri); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		var startTerminalBody types.StartTerminalBody
		if err := c.BindJSON(&startTerminalBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		rawConfig, err := middleware.GetConfig(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		clientSet, err := middleware.GetKubeClient(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		logger, err := middleware.GetLogger(c)
		if middleware.AddErrorToContext(c, err) {
			return
		}

		result, err := controllers.HandleStartTerminal(clientSet, rawConfig, startTerminalUri.ClusterName, startTerminalUri.NamespaceName, startTerminalUri.PodName,
			startTerminalUri.ContainerName, startTerminalBody.Shell, logger)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func ServeTerminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := websocket.CreateAttachHandler()

		handler.ServeHTTP(c.Writer, c.Request)

	}
}
