package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const webACLBucketName = "waf_web_acls"

// WebACLStore provides storage for WAF Web ACLs.
type WebACLStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
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
	var found *WebACL
	err := s.ForEach(func(key string, value []byte) error {
		var webACL WebACL
		if err := json.Unmarshal(value, &webACL); err != nil {
			return err
		}
		if webACL.ARN == arn {
			found = &webACL
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_web_acl_by_arn", err)
	}
	if found == nil {
		return nil, NewStoreError("get_web_acl_by_arn", ErrNotFound)
	}
	return found, nil
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

// List returns a list of WAF Web ACLs from the store with pagination support.
func (s *WebACLStore) List(marker string, maxItems int) (*WebACLListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var webACLs []*WebACL
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var webACL WebACL
		if err := json.Unmarshal(value, &webACL); err != nil {
			return err
		}

		if !started {
			if webACL.ID == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			webACLs = append(webACLs, &webACL)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_web_acls", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &WebACLListResult{
		WebACLs:     webACLs,
		IsTruncated: hasMore,
		NextMarker:  nextMarker,
	}, nil
}
