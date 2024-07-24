package testutils

import "time"

const (
	Domain       = "dana-team.io"
	Hostname     = "custom-capp"
	ManagedLabel = "rcs.dana.io/managed"
)

const (
	CappRevisionNamespace = TestName + "-capp-revision-ns"
	CappRevisionName      = TestName + "-capp-revision"
	CapprevisionsKey      = "capprevisions"
)
const (
	SecretsKey             = "secrets"
	SecretNameKey          = "secretName"
	SecretDataKey          = "test-key"
	SecretDataValue        = "fake"
	SecretDataNewValue     = "faker"
	InvalidSecretType      = "invalid"
	SecretName             = TestName + "-secret"
	SecretDataValueEncoded = "ZmFrZQ=="
	SecretNamespace        = SecretName + "-ns"
	KeyField               = "key"
	ValueField             = "value"
)
const (
	TestName          = "test"
	TestNamespace     = TestName + "-ns"
	NonExistentSuffix = "-non-existent"
)

const (
	ConfigmapsKey      = "configmaps"
	ConfigMapDataKey   = "key"
	ConfigMapDataValue = "value"
)

const (
	RoleKey              = "role"
	AdminKey             = "admin"
	ViewerKey            = "viewer"
	RoleBindingsKey      = "rolebindings"
	UsersKey             = "users"
	RoleBindingsGroupKey = "rbac.authorization.k8s.io"
	CappUserPrefix       = "capp-user-"
)

const (
	LabelSelectorKey     = "labelSelector"
	LabelKey             = "key"
	LabelValue           = "value"
	InvalidLabelSelector = ":--"
	LabelCappName        = "rcs.dana.io/cappName"
)

const (
	NameKey          = "name"
	NamespaceKey     = "namespaces"
	NamespaceNameKey = "namespaceName"
	IdKey            = "id"
	OperationFailed  = "Operation failed"
	InvalidRequest   = "Invalid request"
	DetailsKey       = "details"
	ErrorKey         = "error"
	MessageKey       = "message"
	TypeKey          = "type"
	KubeSystem       = "kube-system"
)

const (
	MetadataKey    = "metadata"
	LabelsKey      = "labels"
	AnnotationsKey = "annotations"
	SpecKey        = "spec"
	StatusKey      = "status"
	CountKey       = "count"
	DataKey        = "data"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)
const (
	CappName      = TestName + "-capp"
	CappsKey      = "capps"
	CappNamespace = TestNamespace + "-" + CappsKey
	CappImage     = "ghcr.io/dana-team/capp-gin-app:v0.2.0"
	ContainerName = "capp-container"
	StateKey      = "state"
	DisabledState = "disabled"
)

const (
	Timeout           = 300 * time.Second
	Interval          = 10 * time.Second
	DefaultEventually = 2 * time.Second
)
