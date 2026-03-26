package waf

// Package waf provides WAF (Web Application Firewall) data store implementations
// for vorpalstacks.

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const webACLAssociationBucketName = "waf_web_acl_associations"

// WebACLAssociationStore provides storage for WAF Web ACL associations.
type WebACLAssociationStore struct {
	*common.BaseStore
}

// NewWebACLAssociationStore creates a new Web ACL association store.
func NewWebACLAssociationStore(store storage.BasicStorage) *WebACLAssociationStore {
	return &WebACLAssociationStore{
		BaseStore: common.NewBaseStore(store.Bucket(webACLAssociationBucketName), "waf"),
	}
}

// Associate associates a Web ACL with a resource.
func (s *WebACLAssociationStore) Associate(webACLArn, resourceArn string) error {
	key := webACLArn + ":" + resourceArn
	association := &WebACLAssociation{
		WebACLArn:   webACLArn,
		ResourceArn: resourceArn,
	}
	if err := s.BaseStore.Put(key, association); err != nil {
		return NewStoreError("associate_web_acl", err)
	}
	return nil
}

// Disassociate removes a Web ACL association from a resource.
func (s *WebACLAssociationStore) Disassociate(webACLArn, resourceArn string) error {
	key := webACLArn + ":" + resourceArn
	if err := s.BaseStore.Delete(key); err != nil {
		return NewStoreError("disassociate_web_acl", err)
	}
	return nil
}

// GetByResourceArn retrieves a Web ACL association by resource ARN.
func (s *WebACLAssociationStore) GetByResourceArn(resourceArn string) (*WebACLAssociation, error) {
	var found *WebACLAssociation
	err := s.ForEach(func(key string, value []byte) error {
		var assoc WebACLAssociation
		if err := json.Unmarshal(value, &assoc); err != nil {
			return err
		}
		if assoc.ResourceArn == resourceArn {
			found = &assoc
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_web_acl_association", err)
	}
	if found == nil {
		return nil, NewStoreError("get_web_acl_association", ErrNotFound)
	}
	return found, nil
}

// GetByWebACLArn retrieves all associations for a Web ACL.
func (s *WebACLAssociationStore) GetByWebACLArn(webACLArn string) ([]*WebACLAssociation, error) {
	var associations []*WebACLAssociation
	err := s.ForEach(func(key string, value []byte) error {
		var assoc WebACLAssociation
		if err := json.Unmarshal(value, &assoc); err != nil {
			return err
		}
		if assoc.WebACLArn == webACLArn {
			associations = append(associations, &assoc)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_web_acl_associations", err)
	}
	return associations, nil
}

// List returns all Web ACL associations.
func (s *WebACLAssociationStore) List() ([]*WebACLAssociation, error) {
	var associations []*WebACLAssociation
	err := s.ForEach(func(key string, value []byte) error {
		var assoc WebACLAssociation
		if err := json.Unmarshal(value, &assoc); err != nil {
			return err
		}
		associations = append(associations, &assoc)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_web_acl_associations", err)
	}
	return associations, nil
}
