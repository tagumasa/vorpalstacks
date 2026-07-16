package waf

import (
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
	return common.FindFirst[WebACLAssociation](s.BaseStore, func(a *WebACLAssociation) bool { return a.ResourceArn == resourceArn })
}

// GetByWebACLArn retrieves all associations for a Web ACL.
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
