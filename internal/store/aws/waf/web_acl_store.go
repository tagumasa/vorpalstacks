package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const webACLBucketName = "waf_web_acls"

// WebACLStore provides storage for WAF Web ACLs.
type WebACLStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewWebACLStore creates a new WebACLStore instance with the specified storage, account ID, and region.
func NewWebACLStore(store storage.BasicStorage, accountId, region string) *WebACLStore {
	return &WebACLStore{
		BaseStore:  common.NewBaseStore(store.Bucket(webACLBucketName), "waf"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Get retrieves a WAF Web ACL by its ID from the store.
// Returns the Web ACL or an error if not found.
func (s *WebACLStore) Get(id string) (*WebACL, error) {
	var webACL WebACL
	if err := s.BaseStore.Get(id, &webACL); err != nil {
		return nil, NewStoreError("get_web_acl", err)
	}
	return &webACL, nil
}

// GetByARN retrieves a WAF Web ACL by its ARN from the store.
// Returns the Web ACL or an error if not found.
func (s *WebACLStore) GetByARN(arn string) (*WebACL, error) {
	return common.FindFirst[WebACL](s.BaseStore, func(w *WebACL) bool { return w.ARN == arn })
}

// Create creates a new WAF Web ACL in the store.
// Returns the created Web ACL or an error if creation fails.
func (s *WebACLStore) Create(id, name, description, scope string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig) (*WebACL, error) {
	webACL := &WebACL{
		ID:               id,
		Name:             name,
		Description:      description,
		Scope:            scope,
		Capacity:         capacity,
		Rules:            rules,
		DefaultAction:    defaultAction,
		VisibilityConfig: visibilityConfig,
		ARN:              s.arnBuilder.BuildWebACLARN(id, scope),
		LockToken:        GenerateLockToken(),
		CreatedAt:        time.Now(),
		ModifiedAt:       time.Now(),
	}

	if err := s.BaseStore.Put(id, webACL); err != nil {
		return nil, NewStoreError("create_web_acl", err)
	}
	return webACL, nil
}

// Update updates an existing WAF Web ACL in the store.
// Returns the updated Web ACL or an error if the Web ACL does not exist or lock token is invalid.
func (s *WebACLStore) Update(id, lockToken string, capacity int64, rules []*Rule, defaultAction *Action, visibilityConfig *VisibilityConfig, description string) (*WebACL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	webACL, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if webACL.LockToken != lockToken {
		return nil, NewStoreError("update_web_acl", ErrLockTokenMismatch)
	}

	webACL.Capacity = capacity
	webACL.Rules = rules
	webACL.DefaultAction = defaultAction
	webACL.VisibilityConfig = visibilityConfig
	if description != "" {
		webACL.Description = description
	}
	webACL.ModifiedAt = time.Now()
	webACL.LockToken = GenerateLockToken()

	if err := s.BaseStore.Put(id, webACL); err != nil {
		return nil, NewStoreError("update_web_acl", err)
	}
	return webACL, nil
}

// Delete deletes a WAF Web ACL by its ID from the store.
// Returns an error if the Web ACL does not exist or lock token is invalid.
func (s *WebACLStore) Delete(id, lockToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	webACL, err := s.Get(id)
	if err != nil {
		return err
	}

	if lockToken != "" && webACL.LockToken != lockToken {
		return NewStoreError("delete_web_acl", ErrLockTokenMismatch)
	}

	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_web_acl", err)
	}
	return nil
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
