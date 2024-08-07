package pagination

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	envDefaultPaginationLimit = "DEFAULT_PAGINATION_LIMIT"
	defaultPaginationLimit    = 100

	firstPage = 1
)

// buildListOptions builds the ListOptions based on page and limit
func buildListOptions(limit int64, continueToken string) v1.ListOptions {
	return v1.ListOptions{
		Limit:         limit,
		Continue:      continueToken,
		LabelSelector: utils.ManagedLabelSelector,
	}
}

// FetchPage fetches the specified page with given limit
func FetchPage[T any](limit, page int, paginator Paginator[types.List[T]]) ([]T, error) {
	var results []T
	var continueToken string

	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}

	// Fetch items until the specified page is reached
	for currentPage := firstPage; currentPage <= page; currentPage++ {
		listOptions := buildListOptions(int64(limit), continueToken)
		list, err := paginator.FetchList(listOptions)
		if err != nil {
			return nil, err
		}

		if currentPage == page {
			return list.Items, nil
		}

		continueToken = list.Continue
		if continueToken == "" {
			break
		}
	}

	return results, nil
}

// extractLimitFromCtx retrieves the pagination limit from the Gin context or defaults to an environment variable
func extractLimitFromCtx(c *gin.Context) (int, error) {
	limit, exists := c.Get(middleware.LimitCtxKey)
	if !exists || limit == 0 {
		return utils.GetEnvNumber(envDefaultPaginationLimit, defaultPaginationLimit)
	}

	return limit.(int), nil
}

// extractPageFromCtx retrieves the page number from the context; returns 1 if not set or an error if conversion fails
func extractPageFromCtx(c *gin.Context) (int, error) {
	page, exists := c.Get(middleware.PageCtxKey)
	if !exists || page == 0 {
		return firstPage, nil
	}

	return page.(int), nil
}

// ExtractPaginationParamsFromCtx retrieves the page and limit from the context
func ExtractPaginationParamsFromCtx(c *gin.Context) (limit, page int, err error) {
	ctxLimit, err := extractLimitFromCtx(c)
	if err != nil {
		return 0, 0, err
	}

	ctxPage, err := extractPageFromCtx(c)
	if err != nil {
		return 0, 0, err
	}

	return ctxLimit, ctxPage, nil
}
