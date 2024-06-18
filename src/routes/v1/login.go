package v1

import (
	"errors"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Login handles user authentication and issues a token on successful login.
// If there's an error during authentication, it responds with an appropriate error message.
func Login(tokenProvider auth.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger, _ := c.Get("logger")
		logger := ctxLogger.(*zap.Logger)

		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			logger.Error("Authorization header not found")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authorization header not found"})
			return
		}

		logger = logger.With(zap.String("user", username))
		token, err := tokenProvider.ObtainToken(username, password, logger, c)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				logger.Warn("Invalid credentials provided", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				logger.Error("Failed to obtain OpenShift token", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
