package v1

import (
	"errors"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Login handles user authentication and issues a token on successful login.
// It expects a JSON payload containing the user's credentials:
//
//	{
//	   "username": "your_username", "password": "your_password"
//	}
//
// It responds with a JSON object containing the token on successful login:
//
//	{
//	  "token": "your_generated_token"
//	}
//
// If there's an error during authentication, it responds with an appropriate error message.
func Login(tokenProvider auth.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger, _ := c.Get("logger")
		logger := ctxLogger.(*zap.Logger)

		var credentials types.Credentials
		if err := c.ShouldBindJSON(&credentials); err != nil {
			logger.Error("Invalid request payload", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		logger = logger.With(zap.String("user", credentials.Username))
		token, err := tokenProvider.ObtainToken(credentials.Username, credentials.Password, logger, c)
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
