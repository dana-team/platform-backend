package utils

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	limitQueryParameter        = "limit"
	limitDefaultValue          = "9"
	continueQueryParameter     = "continue"
	continueDefaultValue       = ""
	searchByNameQueryParameter = "name"
	metadataNameFilter         = "metadata.name=%s"
)

// GetListQueryParameters retrieves pagination and filtering parameters from the Gin context query.
// It returns the limit, continue token, and name for filtering.
func GetListQueryParameters(c *gin.Context) (limit, continueToken, search string) {
	requestLimit := c.DefaultQuery(limitQueryParameter, limitDefaultValue)
	requestContinueToken := c.DefaultQuery(continueQueryParameter, continueDefaultValue)
	searchByName := c.DefaultQuery(searchByNameQueryParameter, "")
	return requestLimit, requestContinueToken, searchByName

}

// GetPaginatedListOptions creates a ListOptions object for paginated requests,
// including field selectors for filtering based on provided parameters.
func GetPaginatedListOptions(limitStr, continueToken, search string) (metav1.ListOptions, error) {
	limit, err := ParseIntStr(limitStr)
	if err != nil {
		return metav1.ListOptions{}, fmt.Errorf("invalid limit parameter: %w", err)
	}

	if search != "" {
		return metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
			// Note: FieldSelector does not support regex, so we currently return only exact matches.
			FieldSelector: fmt.Sprintf(metadataNameFilter, search),
		}, nil

	}
	return metav1.ListOptions{
		Limit:    limit,
		Continue: continueToken,
	}, nil
}

// SetPaginationMetadata sets pagination metadata based on the resources and ListMeta.
func SetPaginationMetadata[T any](resources []T, listMeta metav1.ListMeta) types.ListMetadata {
	listMetadata := types.ListMetadata{}

	listMetadata.Count = len(resources)
	listMetadata.ContinueToken = listMeta.Continue
	listMetadata.RemainingCount = 0

	if listMeta.RemainingItemCount != nil {
		listMetadata.RemainingCount = int(*listMeta.RemainingItemCount)
	}

	return listMetadata
}
