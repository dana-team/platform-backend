package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type List[T any] struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of objects.
	Items []T `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type PaginationParams struct {
	Limit int `form:"limit,omitempty" json:"limit" binding:"min=0"`
	Page  int `form:"page,default=1,omitempty" json:"page" binding:"min=1"`
}
