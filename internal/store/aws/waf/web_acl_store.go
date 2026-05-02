package waf

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const webACLBucketName = "waf_web_acls"

var webACLAccessor = wafResourceAccessor[WebACL]{
	getIDFn:        func(r *WebACL) string { return r.ID },
	getNameFn:      func(r *WebACL) string { return r.Name },
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
func (s *WebACLStore) Create(id, name, description, scope string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig) (*WebACL, error) {
	if existing, _ := s.FindByName(name); existing != nil {
		return nil, ErrAlreadyExists
	}
	webACL := &WebACL{
		ID:               id,
		Name:             name,
		Description:      description,
		Scope:            scope,
		Capacity:         capacity,
		Rules:            rules,
		DefaultAction:    defaultAction,
		VisibilityConfig: visibilityConfig,
		Tags:             []types.Tag{},
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}
	webACL.ARN = s.arnBuilder.BuildWebACLARN(id, scope)
	SetTimestamps(&webACLAccessor, webACL)
	if err := s.Put(id, webACL, "create_web_acl"); err != nil {
		return nil, err
	}
	return webACL, nil
}

// Update updates an existing WAF Web ACL in the store.
// Returns the updated Web ACL or an error if the Web ACL does not exist or lock token is invalid.
func (s *WebACLStore) Update(id, lockToken string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig, description string) (*WebACL, error) {
	return s.UpdateWithLockToken(id, lockToken, func(webACL *WebACL) error {
		webACL.Capacity = capacity
		webACL.Rules = rules
		webACL.DefaultAction = defaultAction
		webACL.VisibilityConfig = visibilityConfig
		if description != "" {
			webACL.Description = description
		}
		return nil
	}, "update_web_acl")
}

// List returns a list of WAF Web ACLs from the store with pagination
// support.
func (s *WebACLStore) List(marker string, maxItems int) (*WebACLListResult, error) {
	result, err := common.List[WebACL](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_web_acls", err)
	}
	return &WebACLListResult{
		WebACLs:     result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
