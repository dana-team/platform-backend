package types

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

type CappRevision struct {
	Metadata    Metadata                        `json:"metadata" binding:"required"`
	Annotations []KeyValue                      `json:"annotations" binding:"required"`
	Labels      []KeyValue                      `json:"labels" binding:"required"`
	Spec        cappv1alpha1.CappRevisionSpec   `json:"spec" binding:"required"`
	Status      cappv1alpha1.CappRevisionStatus `json:"status" binding:"required"`
}

type CappRevisionList struct {
	CappRevisions []string `json:"capprevisions"`
	ListMetadata
}

type CappRevisionNamespaceUri struct {
	ClusterName   string `uri:"clusterName"`
	NamespaceName string `uri:"namespaceName" binding:"required"`
	CappName      string `uri:"cappName"`
}

type CappRevisionUri struct {
	NamespaceName    string `uri:"namespaceName" binding:"required"`
	CappRevisionName string `uri:"cappRevisionName" binding:"required"`
}

type CappRevisionQuery struct {
	LabelSelector string `form:"labelSelector"`
}
