package main

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/routes/v1"
	dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env file")
	}

	logger := initializeLogger()
	defer syncLogger(logger)

	tokenProvider := auth.DefaultTokenProvider{}
	engine := initializeRouter(logger, tokenProvider)
	if err := engine.Run(); err != nil {
		panic(err.Error())
	}
}

// initializeLogger initializes a new Zap logger instance.
func initializeLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	return logger
}

// syncLogger syncs the logger to ensure all pending logs are written before shutdown.
func syncLogger(logger *zap.Logger) {
	if err := logger.Sync(); err != nil {
		log.Fatalf("Error syncing logger: %v", err)
	}
}

// initializeRouter initializes the Gin router with routes for API v1.
func initializeRouter(logger *zap.Logger, tokenProvider auth.TokenProvider) *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.LoggerMiddleware(logger))
	v1.SetupRoutes(engine, tokenProvider, newScheme())

	return engine
}

// newScheme adds the relevant APIs to the scheme for the K8S client.
func newScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cappv1alpha1.AddToScheme(scheme))
	utilruntime.Must(dnsrecordv1alpha1.AddToScheme(scheme))

	return scheme
}
