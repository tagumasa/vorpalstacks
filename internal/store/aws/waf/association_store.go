package waf

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const webACLAssociationBucketName = "waf_web_acl_associations"

// WebACLAssociationStore provides storage for WAF Web ACL associations.
//
// AWS WAFv2 enforces a one-WebACL-per-resource constraint: a resource
// can be associated with at most one Web ACL at any time. Calling
// AssociateWebACL on a resource that already has an association replaces
// the existing one. The store key is therefore the resource ARN alone,
// ensuring the constraint is enforced at the storage level.
type WebACLAssociationStore struct {
	*common.BaseStore
}

// NewWebACLAssociationStore creates a new Web ACL association store.
func NewWebACLAssociationStore(store storage.BasicStorage) *WebACLAssociationStore {
	return &WebACLAssociationStore{
		BaseStore: common.NewBaseStore(store.Bucket(webACLAssociationBucketName), "waf"),
	}
}

// Associate associates a Web ACL with a resource. If the resource already
// has an association with a different Web ACL, the existing association
// is replaced.
func (s *WebACLAssociationStore) Associate(webACLArn, resourceArn string) error {
	association := &WebACLAssociation{
		WebACLArn:   webACLArn,
		ResourceArn: resourceArn,
	}
	if err := s.BaseStore.Put(resourceArn, association); err != nil {
		return NewStoreError("associate_web_acl", err)
	}
	return nil
}

// Disassociate removes the Web ACL association from the specified resource.
func (s *WebACLAssociationStore) Disassociate(resourceArn string) error {
	if err := s.BaseStore.Delete(resourceArn); err != nil {
		return NewStoreError("disassociate_web_acl", err)
	}
	return nil
}

// GetByResourceArn retrieves the Web ACL association for the specified
// resource ARN. Returns ErrNotFound if no association exists.
func (s *WebACLAssociationStore) GetByResourceArn(resourceArn string) (*WebACLAssociation, error) {
	var assoc WebACLAssociation
	if err := s.BaseStore.Get(resourceArn, &assoc); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_association_by_resource", ErrNotFound)
		}
		return nil, NewStoreError("get_association_by_resource", err)
	}
	return &assoc, nil
}

// GetByWebACLArn retrieves all associations for the specified Web ACL.
func (s *WebACLAssociationStore) GetByWebACLArn(webACLArn string) ([]*WebACLAssociation, error) {
	all, err := common.ListAll[WebACLAssociation](s.BaseStore)
	if err != nil {
		return nil, NewStoreError("list_web_acl_associations", err)
	}
	var associations []*WebACLAssociation
	for _, assoc := range all {
		if assoc.WebACLArn == webACLArn {
			associations = append(associations, assoc)
		}
	}
	return associations, nil
}

// List returns all Web ACL associations.
func (s *WebACLAssociationStore) List() ([]*WebACLAssociation, error) {
	return common.ListAll[WebACLAssociation](s.BaseStore)
}
