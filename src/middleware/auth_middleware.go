package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	httpAuthorizationHeader     = "Authorization"
	httpBearerTokenPrefix       = "Bearer"
	validBearerTokenPartsLength = 2
	httpBearerTokenPrefixIndex  = 0
	httpBearerTokenIndex        = 1
)

const (
	envKubeAPIServer      = "KUBE_API_SERVER"
	envInsecureSkipVerify = "INSECURE_SKIP_VERIFY"
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

		config, err := createKubernetesConfig(token, os.Getenv(envKubeAPIServer))
		if err != nil {
			userLogger.Error("Failed to create Kubernetes client config", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes client config"})
			return
		}

		kubeClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			userLogger.Error("Failed to create Kubernetes client", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes client"})
			return
		}

		schema := scheme.Scheme
		if err := cappv1.AddToScheme(scheme.Scheme); err != nil {
			userLogger.Error("Failed to create Kubernetes dynamic client schema", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes dynamic client schema"})
			return
		}

		dynClient, err := client.New(config, client.Options{Scheme: schema})
		if err != nil {
			userLogger.Error("Failed to create Kubernetes dynamic client", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kubernetes dynamic client"})
			return
		}

		log.SetLogger(logr.New(logr.Discard().GetSink()))

		// Update the logger with the username
		c.Set("logger", userLogger)
		c.Set("kubeClient", kubeClient)
		c.Set("dynClient", dynClient)

		c.Next()
	}
}

// validateToken validates the format and presence of the Authorization token.
func validateToken(c *gin.Context) (string, error) {
	token := c.GetHeader(httpAuthorizationHeader)
	if token == "" {
		return "", fmt.Errorf("authorization token not provided")
	}

	tokenParts := strings.Split(token, " ")

	// Check if the token is in the format "Bearer <token>"
	if len(tokenParts) != validBearerTokenPartsLength || tokenParts[httpBearerTokenPrefixIndex] != httpBearerTokenPrefix {
		return "", fmt.Errorf("invalid authentication token")
	}

	return tokenParts[httpBearerTokenIndex], nil
}

// createKubernetesConfig creates a new Kubernetes client config using the provided token.
func createKubernetesConfig(token, kubeApiServer string) (*rest.Config, error) {
	skipTlsVerify, err := utils.GetEnvBool(envInsecureSkipVerify, true)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		BearerToken: token,
		Host:        kubeApiServer,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: skipTlsVerify,
		},
	}
	return config, nil
}
