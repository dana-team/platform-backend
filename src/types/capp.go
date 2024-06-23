package types

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

type Capp struct {
	Metadata    Metadata                `json:"metadata" binding:"required"`
	Annotations []KeyValue              `json:"annotations" binding:"required"`
	Labels      []KeyValue              `json:"labels" binding:"required"`
	Spec        cappv1alpha1.CappSpec   `json:"spec" binding:"required"`
	Status      cappv1alpha1.CappStatus `json:"status" binding:"required"`
}

type CreateCapp struct {
	Metadata    CreateMetadata        `json:"metadata" binding:"required"`
	Annotations []KeyValue            `json:"annotations" binding:"required"`
	Labels      []KeyValue            `json:"labels" binding:"required"`
	Spec        cappv1alpha1.CappSpec `json:"spec" binding:"required"`
}

type CreateMetadata struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCapp struct {
	Annotations []KeyValue            `json:"annotations"`
	Labels      []KeyValue            `json:"labels"`
	Spec        cappv1alpha1.CappSpec `json:"spec"`
}

type CappQuery struct {
	LabelSelector string `form:"labelSelector"`
}

type CappList struct {
	Capps []Capp `json:"capps"`
	ListMetadata
}

type CappNamespaceUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
}

type CappUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName" binding:"required"`
}

type CappError struct {
	Message string `json:"message"`
}
