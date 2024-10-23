package testutils

import (
	"time"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

var (
	cappAPIGroup   = cappv1alpha1.GroupVersion.Group
	ManagedLabel   = cappAPIGroup + "/managed"
	ManagedByLabel = cappAPIGroup + "/managed-by"
)

const (
	Domain      = "dana-team.io"
	Hostname    = "custom-capp"
	DefaultZone = "dana-dev.com"
	Rcs         = "rcs"
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
	TokenName              = "token"
	TokenNamespace         = TokenName + "-ns"
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
	ServiceAccountsKey     = "serviceAccounts"
	TokenKey               = "token"
	ExpirationTimestampKey = "expirationTimestamp"
	Secret                 = "Secret"
	V1                     = "v1"
	Value                  = "value"
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

var (
	ParentCappLabel      = cappAPIGroup + "/parent-capp"
	ParentCappNSLabel    = cappAPIGroup + "/parent-capp-ns"
	LastUpdatedCappLabel = cappAPIGroup + "/last-updated-by"
	HasPlacementLabel    = cappAPIGroup + "/has-placement"
	LabelCappName        = cappAPIGroup + "/cappName"
)

const (
	LabelSelectorKey     = "labelSelector"
	LabelKey             = "key"
	LabelValue           = "value"
	InvalidLabelSelector = ":--"
)

const (
	NameKey          = "name"
	NamespaceKey     = "namespaces"
	NamespaceNameKey = "namespaceName"
	IdKey            = "id"
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

var (
	PlacementRegionLabelKey      = cappAPIGroup + "/region"
	PlacementEnvironmentLabelKey = cappAPIGroup + "/environment"
)

const (
	CappName                = TestName + "-capp"
	EnvironmentName         = TestName + "-environment"
	RegionName              = TestName + "-region"
	PlacementName           = TestName + "-placement"
	SiteName                = TestName + "-site"
	PlacementEnvironmentKey = "environment"
	PlacementRegionKey      = "region"
	CappsKey                = "capps"
	RecordsKey              = "records"
	CappNamespace           = TestNamespace + "-" + CappsKey
	CappImage               = "ghcr.io/dana-team/capp-gin-app:v0.2.0"
	ContainerName           = "capp-container"
	StateKey                = "state"
	DisabledState           = "disabled"
	EnabledState            = "enabled"
	LastCreatedRevision     = "lastCreatedRevision"
	LastReadyRevision       = "lastReadyRevision"
	NoRevision              = "No revision available"
	Available               = "available"
	Unknown                 = "unknown"
	Unavailable             = "unavailable"
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
	Timeout           = 360 * time.Second
	Interval          = 10 * time.Second
	DefaultEventually = 2 * time.Second
)

const (
	ReasonNotFound      = "NotFound"
	ReasonBadRequest    = "BadRequest"
	ReasonAlreadyExists = "AlreadyExists"
)

const (
	ServiceAccountParam    = "serviceaccounts"
	ExpirationSecondsParam = "expirationSeconds"
	TokenRequestSuffix     = "token-request"
)
