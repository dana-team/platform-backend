package types

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type Capp struct {
	Metadata    Metadata                `json:"metadata" binding:"required"`
	Annotations []KeyValue              `json:"annotations"`
	Labels      []KeyValue              `json:"labels"`
	Spec        cappv1alpha1.CappSpec   `json:"spec" binding:"required"`
	Status      cappv1alpha1.CappStatus `json:"status" binding:"required"`
}

type CreateCapp struct {
	Metadata    CreateMetadata        `json:"metadata" binding:"required"`
	Annotations []KeyValue            `json:"annotations"`
	Labels      []KeyValue            `json:"labels"`
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

type GetCappQuery struct {
	LabelSelector string `form:"labelSelector" json:"labelSelector"`
}

type CreateCappQuery struct {
	Environment string `form:"environment" json:"environment"`
	Region      string `form:"region" json:"region"`
}

type CappList struct {
	Capps []CappSummary `json:"capps"`
	ListMetadata
}

type CappSummary struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Images []string `json:"images"`
}

type CappStateResponse struct {
	Name  string `json:"name"`
	State string `json:"state" binding:"oneof=enabled disabled"`
}

type CappNamespaceUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
}

type CappUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName" binding:"required"`
}

type CappState struct {
	State string `json:"state" binding:"required,oneof=enabled disabled"`
}
type GetCappStateResponse struct {
	LastCreatedRevision string `json:"lastCreatedRevision"`
	LastReadyRevision   string `json:"lastReadyRevision"`
	State               string `json:"state" binding:"required,oneof=enabled disabled"`
}

type GetDNSResponse struct {
	Records []DNS `json:"records"`
}

type DNS struct {
	Status corev1.ConditionStatus `json:"status"`
	Name   string                 `json:"name"`
}
