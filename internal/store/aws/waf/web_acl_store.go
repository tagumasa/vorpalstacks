package waf

import (
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const webACLBucketName = "waf_web_acls"

var webACLAccessor = wafResourceAccessor[WebACL]{
	getIDFn:        func(r *WebACL) string { return r.ID },
	getNameFn:      func(r *WebACL) string { return r.Name },
	getScopeFn:     func(r *WebACL) string { return r.Scope },
	getARNFn:       func(r *WebACL) string { return r.ARN },
	setARNFn:       func(r *WebACL, arn string) { r.ARN = arn },
	getLockTokenFn: func(r *WebACL) string { return r.LockToken },
	setLockTokenFn: func(r *WebACL, lt string) { r.LockToken = lt },
	setModifiedFn:  func(r *WebACL) { r.ModifiedAt = time.Now() },
}

// WebACLStore provides storage for WAF Web ACLs.
type WebACLStore struct {
	*ResourceStore[WebACL]
}

// NewWebACLStore creates a new WebACLStore instance with the specified storage, account ID, and region.
func NewWebACLStore(store storage.BasicStorage, accountId, region string) *WebACLStore {
	return &WebACLStore{
		ResourceStore: NewResourceStore[WebACL](store, webACLBucketName, NewARNBuilder(accountId, region), webACLAccessor),
	}
}

// Create creates a new WAF Web ACL in the store.
// Returns the created Web ACL or an error if creation fails.
func (s *WebACLStore) Create(webACL *WebACL) (*WebACL, error) {
	if existing, _ := s.FindByNameAndScope(webACL.Name, webACL.Scope); existing != nil {
		return nil, ErrAlreadyExists
	}
	if webACL.Tags == nil {
		webACL.Tags = []types.Tag{}
	}
	webACL.ARN = s.arnBuilder.BuildWebACLARN(webACL.ID, webACL.Scope)
	SetTimestamps(&webACLAccessor, webACL)
	if err := s.Put(webACL.ID, webACL, "create_web_acl"); err != nil {
		return nil, err
	}
	return webACL, nil
}

// Update updates an existing WAF Web ACL in the store. UpdateWebACL is
// a full-replace operation ("completely replaces the mutable
// specifications" per the Smithy model documentation): capacity is the
// value recomputed by the service layer from the resulting rule set,
// and an empty description clears the field.
// Returns the updated Web ACL or an error if the Web ACL does not exist or lock token is invalid.
func (s *WebACLStore) Update(id, lockToken string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig, description string, extraFn ...func(*WebACL)) (*WebACL, error) {
	return s.UpdateWithLockToken(id, lockToken, func(webACL *WebACL) error {
		webACL.Capacity = capacity
		webACL.Rules = rules
		webACL.DefaultAction = defaultAction
		webACL.VisibilityConfig = visibilityConfig
		webACL.Description = description
		for _, fn := range extraFn {
			if fn != nil {
				fn(webACL)
			}
		}
		return nil
	}, "update_web_acl")
}

// List returns a list of WAF Web ACLs from the store with pagination
// support. If scope is non-empty, only Web ACLs matching the specified
// scope are returned (filter is applied during iteration, before
// pagination, to avoid empty pages).
func (s *WebACLStore) List(marker string, maxItems int, scope string) (*WebACLListResult, error) {
	var filter common.FilterFunc[WebACL]
	if scope != "" {
		filter = func(w *WebACL) bool { return w.Scope == scope }
	}
	result, err := common.List[WebACL](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, NewStoreError("list_web_acls", err)
	}
	return &WebACLListResult{
		WebACLs:     result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
