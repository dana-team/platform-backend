package testutils

import "time"

const (
	Domain       = "dana-team.io"
	Hostname     = "custom-capp"
	DefaultZone  = "dana-dev.com"
	ManagedLabel = "rcs.dana.io/managed"
	cappAPIGroup = "rcs.dana.io"
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
	ServiceAccountsKey = "serviceaccounts"
	TokenKey           = "token"
	Secret             = "Secret"
	V1 = "v1"
	Value = "value"
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
	ParentCappLabel      = cappAPIGroup + "/parent-capp"
	ParentCappNSLabel    = cappAPIGroup + "/parent-capp-ns"
)

const (
	NameKey          = "name"
	NamespaceKey     = "namespaces"
	NamespaceNameKey = "namespaceName"
	IdKey            = "id"
	InvalidRequest   = "Invalid request"
	ReasonKey        = "reason"
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
	CappName            = TestName + "-capp"
	CappsKey            = "capps"
	RecordsKey          = "records"
	CappNamespace       = TestNamespace + "-" + CappsKey
	CappImage           = "ghcr.io/dana-team/capp-gin-app:v0.2.0"
	ContainerName       = "capp-container"
	StateKey            = "state"
	DisabledState       = "disabled"
	EnabledState        = "enabled"
	LastCreatedRevision = "lastCreatedRevision"
	LastReadyRevision   = "lastReadyRevision"
	NoRevision          = "No revision available"
	Available           = "available"
	Unknown             = "unknown"
	UnAvailable         = "unavailable"
)

const (
	CappResourceKey = "rcs.dana.io/parent-capp"
)

const (
	ServiceAccountName       = TestName + "-serviceAccount"
	ServiceAccountAnnotation = "kubernetes.io/service-account.name"
)

const (
	ContainersKey     = "containers"
	ContainerNameKey  = "containerName"
	TestContainerName = "test-container"
	Image             = "nginx"
)

const (
	PodsKey    = "pods"
	PodNameKey = "podName"
	PodName    = TestName + "-pod"
)

const (
	Timeout           = 300 * time.Second
	Interval          = 10 * time.Second
	DefaultEventually = 2 * time.Second
)

const (
	ReasonNotFound      = "NotFound"
	ReasonBadRequest    = "BadRequest"
	ReasonAlreadyExists = "AlreadyExists"
)
