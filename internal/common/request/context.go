package request

import (
	"context"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/pkg/graphengine"
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
	iamStore       iamstore.IAMStoreInterface
	iamValidator   *iam.IAMValidator
	AccountID      string
	Region         string
	Principal      string
	PrincipalID    string
	PrincipalType  PrincipalType
	SourceIP       string
	UserAgent      string
	auditRecorder  AuditRecorder
	graphDBManager *graphengine.DB
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

// SetIAMStore sets the IAM store for the request context.
func (c *RequestContext) SetIAMStore(store iamstore.IAMStoreInterface) {
	c.iamStore = store
}

// GetIAMStore returns the IAM store.
func (c *RequestContext) GetIAMStore() iamstore.IAMStoreInterface {
	return c.iamStore
}

// GetIAMValidator returns the IAM validator for the request context.
func (c *RequestContext) GetIAMValidator() *iam.IAMValidator {
	if c.iamValidator == nil && c.iamStore != nil {
		c.iamValidator = iam.NewIAMValidator(c.iamStore.Roles(), c.AccountID)
	}
	return c.iamValidator
}

// SetGraphDBManager sets the graph database manager for NeptuneGraph queries.
func (c *RequestContext) SetGraphDBManager(db *graphengine.DB) {
	c.graphDBManager = db
}

// GraphReader returns the graph database reader for NeptuneGraph queries.
func (c *RequestContext) GraphReader() graphengine.GraphReader {
	if c.graphDBManager == nil {
		return nil
	}
	return c.graphDBManager
}

// GraphWriter returns the graph database writer for NeptuneGraph mutations.
func (c *RequestContext) GraphWriter() graphengine.GraphWriter {
	if c.graphDBManager == nil {
		return nil
	}
	return c.graphDBManager
}
