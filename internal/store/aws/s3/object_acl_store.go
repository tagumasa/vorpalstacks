package s3

import (
	types "vorpalstacks/internal/common/tags"
)

// SetTags sets the tags for an object. A non-empty versionId targets that
// specific object version instead of the current one.
func (s *ObjectStore) SetTags(bucket, key, versionId string, tags []types.Tag) error {
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.Tags = tags
		return nil
	})
}

// SetACL sets the access control list for the current object.
func (s *ObjectStore) SetACL(bucket, key string, acp *AccessControlPolicy) error {
	return s.mutateObjectRecord(bucket, key, "", func(obj *Object) error {
		obj.ACL = acp
		return nil
	})
}
