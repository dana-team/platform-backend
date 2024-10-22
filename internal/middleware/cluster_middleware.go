package middleware

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/internal/customerrors"
	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ClusterCtxKey = "cluster"
	cappName      = "cappName"
	namespaceName = "namespaceName"
	clusterName   = "clusterName"
)

const (
	errNoClusterOrNameAndNamespaceProvided = "either 'cluster' must be provided, or both 'cappName' and 'namespace' must be provided."
)

// ClusterMiddleware retrieves the cluster based on capp and adds it to the request context.
func ClusterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cluster := c.Param(clusterName)
		name := c.Param(cappName)
		namespace := c.Param(namespaceName)

		if cluster == "" && (name == "" || namespace == "") {
			err := customerrors.NewValidationError(errNoClusterOrNameAndNamespaceProvided)
			AddErrorToContext(c, err)
			c.Abort()
			return
		}

		if name != "" && namespace != "" {
			capp, err := getCapp(c, name, namespace)
			if AddErrorToContext(c, err) {
				c.Abort()
				return
			}
			cluster = capp.Status.ApplicationLinks.Site
		}

		c.Set(ClusterCtxKey, cluster)
		c.Next()
	}
}

// getCapp retrieves a Capp object from the Kubernetes cluster based on the provided name and namespace.
func getCapp(c *gin.Context, name, namespace string) (*cappv1alpha1.Capp, error) {
	capp := &cappv1alpha1.Capp{}
	kubeClient, err := GetDynClient(c)
	if err != nil {
		return nil, err
	}

	err = kubeClient.Get(c.Request.Context(), client.ObjectKey{Namespace: namespace, Name: name}, capp)
	if err != nil {
		return nil, err
	}

	return capp, nil
}
