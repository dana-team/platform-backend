package v1

import (
	"errors"
	"github.com/dana-team/platform-backend/internal/auth"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/dana-team/platform-backend/internal/middleware"
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	errAuthorizationHeaderNotFound = "Authorization header not found"
)

// Login handles user authentication and issues a token on successful login.
// If there's an error during authentication, it responds with an appropriate error message.
func Login(tokenProvider auth.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger, _ := c.Get("logger")
		logger := ctxLogger.(*zap.Logger)

		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			logger.Error(errAuthorizationHeaderNotFound)
			middleware.AddErrorToContext(c, customerrors.NewValidationError(errAuthorizationHeaderNotFound))
			return
		}

		logger = logger.With(zap.String("user", username))
		token, err := tokenProvider.ObtainToken(username, password, logger, c)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				logger.Warn("Invalid credentials provided", zap.Error(err))
				middleware.AddErrorToContext(c, customerrors.NewUnauthorizedError(err.Error()))
			} else {
				logger.Error("Failed to obtain OpenShift token", zap.Error(err))
				middleware.AddErrorToContext(c, customerrors.NewInternalServerError(err.Error()))
			}
			return
		}

		result := types.LoginOutput{
			Token: token,
		}

		c.JSON(http.StatusOK, result)
	}
}
