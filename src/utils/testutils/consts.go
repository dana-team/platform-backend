package testutils

const (
	Domain       = "dana-team.io"
	Hostname     = "custom-capp"
	ManagedLabel = "rcs.dana.io/managed"
)

const (
	CappRevisionNamespace = TestName + "-capp-revision-ns"
	CappRevisionName      = TestName + "-capp-revision"
)
const (
	SecretsKey             = "secrets"
	OpaqueType             = "opaque"
	SecretType             = "Opaque"
	SecretNameKey          = "secretName"
	SecretDataKey          = "test-key"
	SecretDataValue        = "fake"
	SecretDataNewValue     = "faker"
	InvalidSecretType      = "invalid"
	SecretName             = TestName + "-secret"
	SecretDataValueEncoded = "ZmFrZQ=="
	SecretNamespace        = SecretName + "-ns"
)
const (
	TestName          = "test"
	TestNamespace     = TestName + "-ns"
	NonExistentSuffix = "-non-existent"
)

const (
	LabelSelectorKey = "labelSelector"
	LabelKey         = "key"
	LabelValue       = "value"
)

const (
	NameKey          = "name"
	NameSpaceKey     = "namespaces"
	NameSpaceNameKey = "namespaceName"
	IdKey            = "id"
	OperationFailed  = "Operation failed"
	InvalidRequest   = "Invalid request"
	DetailsKey       = "details"
	ErrorKey         = "error"
	MessageKey       = "message"
	TypeKey          = "type"
)

const (
	Metadata    = "metadata"
	Labels      = "labels"
	Annotations = "annotations"
	Spec        = "spec"
	Status      = "status"
	Count       = "count"
	Data        = "data"
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
)
