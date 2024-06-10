package middleware

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"strings"
)

// TokenAuthMiddleware validates the Authorization header and sets up Kubernetes client.
func TokenAuthMiddleware(tokenProvider auth.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger, exists := c.Get("logger")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "logger not set in context"})
			return
		}
		logger := ctxLogger.(*zap.Logger)

		token, err := validateToken(c)
		if err != nil {
			logger.Error("Failed to obtain OpenShift token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Failed to obtain OpenShift token": err.Error()})
			return
		}

		username, err := tokenProvider.ObtainUsername(token, logger)
		if err != nil {
			logger.Error("Failed to get user info", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}
		userLogger := logger.With(zap.String("user", username))

		client, err := createKubernetesClient(token, os.Getenv("KUBE_API_SERVER"))
		if err != nil {
			userLogger.Error("Failed to create Kubernetes client", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes client"})
			return
		}

		// Update the logger with the username
		c.Set("logger", userLogger)
		c.Set("kubeClient", client)

		c.Next()
	}
}

// validateToken validates the format and presence of the Authorization token.
func validateToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", fmt.Errorf("authorization token not provided")
	}

	tokenParts := strings.Split(token, " ")
	// Check if the token is in the format "Bearer <token>"
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// len(tokenParts) != 2 ensures the token is split into exactly two parts: "Bearer" and the actual token
		// tokenParts[0] != BearerTokenPrefix ensures the token starts with the prefix "Bearer"
		return "", fmt.Errorf("invalid authentication token")
	}

	// Return the actual token part, which is the second element in the split parts
	return tokenParts[1], nil
}

// createKubernetesClient creates a Kubernetes client using the provided token.
func createKubernetesClient(token, kube_api_server string) (*kubernetes.Clientset, error) {
	skipTlsVerify, err := utils.GetEnvBool("INSECURE_SKIP_VERIFY", true)
	if err != nil {
		return nil, fmt.Errorf("failed to parse INSECURE_SKIP_VERIFY: %v", err)
	}

	config := &rest.Config{
		BearerToken: token,
		Host:        kube_api_server,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: skipTlsVerify,
		},
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return client, nil
}
