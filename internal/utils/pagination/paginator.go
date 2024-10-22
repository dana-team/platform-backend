package pagination

import (
	"context"
	"fmt"
	"github.com/dana-team/platform-backend/internal/types"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Paginator interface with methods for pagination
type Paginator[T any] interface {
	FetchList(listOptions metav1.ListOptions) (*T, error)
}

// GenericPaginator to handle any Kubernetes list type
type GenericPaginator struct {
	Ctx    context.Context
	Logger *zap.Logger
}

// FetchList fetches the list of resources of any type
func (p *GenericPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[interface{}], error) {
	// This method needs to be adapted based on the type of resource
	return nil, fmt.Errorf("FetchList method not implemented")
}

func CreatePaginator(ctx context.Context, logger *zap.Logger) GenericPaginator {
	return GenericPaginator{
		Ctx:    ctx,
		Logger: logger,
	}
}
