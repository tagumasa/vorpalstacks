package request

import (
	"context"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/core/storage"
)

// AuditRecorder records audit events for CloudTrail.
type AuditRecorder interface {
	Record(event interface{}) error
}

// PrincipalType represents the type of principal making a request.
type PrincipalType string

const (
	// PrincipalTypeUser represents an IAM user.
	PrincipalTypeUser PrincipalType = "User"
	// PrincipalTypeRole represents an IAM role.
	PrincipalTypeRole PrincipalType = "Role"
	// PrincipalTypeAnonymous represents an unauthenticated request.
	PrincipalTypeAnonymous PrincipalType = "Anonymous"
)

// RequestContext holds the context information for an AWS API request.
type RequestContext struct {
	context.Context
	storageManager *storage.RegionStorageManager
	iamStore       any
	iamRoles       iam.RolePolicyProvider
	iamValidator   *iam.IAMValidator
	AccountID      string
	Region         string
	Principal      string
	PrincipalID    string
	PrincipalType  PrincipalType
	SourceIP       string
	UserAgent      string
	auditRecorder  AuditRecorder
	graphReader    any
	graphWriter    any
}

// NewRequestContext creates a new RequestContext with the given parameters.
func NewRequestContext(ctx context.Context, storageMgr *storage.RegionStorageManager, accountID, region string) *RequestContext {
	if ctx == nil {
		ctx = context.Background()
	}
	if region == "" {
		region = DefaultRegion
	}
	return &RequestContext{
		Context:        ctx,
		storageManager: storageMgr,
		AccountID:      accountID,
		Region:         region,
	}
}

// GetRegion returns the AWS region for the request.
func (c *RequestContext) GetRegion() string {
	if c.Region == "" {
		return DefaultRegion
	}
	return c.Region
}

// GetAccountID returns the AWS account ID for the request.
func (c *RequestContext) GetAccountID() string {
	return c.AccountID
}

// GetStorage returns the storage for the request's region.
func (c *RequestContext) GetStorage() (storage.BasicStorage, error) {
	return c.storageManager.GetStorage(c.GetRegion())
}

// GetGlobalStorage returns the global storage for region-independent services.
func (c *RequestContext) GetGlobalStorage() (storage.BasicStorage, error) {
	return c.storageManager.GetGlobalStorage()
}

// GetStorageForService returns the appropriate storage for the given service.
func (c *RequestContext) GetStorageForService(serviceName string) (storage.BasicStorage, error) {
	if c.IsGlobalService(serviceName) {
		return c.GetGlobalStorage()
	}
	return c.GetStorage()
}

// IsGlobalService returns true if the service is global (not region-specific).
func (c *RequestContext) IsGlobalService(serviceName string) bool {
	globalServices := map[string]bool{
		"iam":        true,
		"route53":    true,
		"cloudfront": true,
		"sts":        true,
	}
	return globalServices[serviceName]
}

// SetAuditRecorder sets the audit recorder for the request context.
func (c *RequestContext) SetAuditRecorder(recorder AuditRecorder) {
	c.auditRecorder = recorder
}

// GetAuditRecorder returns the audit recorder for the request context.
func (c *RequestContext) GetAuditRecorder() AuditRecorder {
	return c.auditRecorder
}

// HasAuditRecorder returns true if an audit recorder is set.
func (c *RequestContext) HasAuditRecorder() bool {
	return c.auditRecorder != nil
}

// SetIAMStore sets the IAM store and role provider for the request context.
// The store is retained as-is for services that need direct store access.
// The roles provider is used internally by GetIAMValidator for policy evaluation.
func (c *RequestContext) SetIAMStore(store any, roles iam.RolePolicyProvider) {
	c.iamStore = store
	c.iamRoles = roles
}

// GetIAMStore returns the IAM store. Callers should type-assert to the
// concrete IAMStoreInterface from store/aws/iam.
func (c *RequestContext) GetIAMStore() any {
	return c.iamStore
}

// GetIAMValidator returns the IAM validator for the request context.
func (c *RequestContext) GetIAMValidator() *iam.IAMValidator {
	if c.iamValidator == nil && c.iamRoles != nil {
		c.iamValidator = iam.NewIAMValidator(c.iamRoles, c.AccountID)
	}
	return c.iamValidator
}

// SetGraphDBManager sets the graph database reader and writer for
// NeptuneGraph queries. The caller (dispatcher) extracts these from
// the concrete *graphengine.DB at injection time.
func (c *RequestContext) SetGraphDBManager(reader, writer any) {
	c.graphReader = reader
	c.graphWriter = writer
}

// GraphReader returns the graph database reader for NeptuneGraph queries.
// Callers should type-assert to graphengine.GraphReader.
func (c *RequestContext) GraphReader() any {
	return c.graphReader
}

// GraphWriter returns the graph database writer for NeptuneGraph mutations.
// Callers should type-assert to graphengine.GraphWriter or graphengine.GraphStore.
func (c *RequestContext) GraphWriter() any {
	return c.graphWriter
}
