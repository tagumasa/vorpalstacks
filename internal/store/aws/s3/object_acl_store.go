package s3

import (
	"vorpalstacks/internal/utils/aws/types"
)

// updateObjectMetadata is a shared helper for SetTags and SetACL. It acquires
// the key lock, reads the current object metadata, applies the mutation
// callback, then writes back to both the versioned record and the _latest
// pointer when versioning is enabled.
func (s *ObjectStore) updateObjectMetadata(bucket, key string, mutate func(obj *Object)) error {
	lockKey := bucket + "#" + key
	s.keyLocker.Lock(lockKey)
	defer func() {
		s.keyLocker.Unlock(lockKey)
		s.keyLocker.Delete(lockKey)
	}()

	obj, err := s.GetMetadata(bucket, key)
	if err != nil {
		return err
	}

	previousProto := ObjectToProto(obj)

	mutate(obj)

	if s.isVersioningEnabled(bucket) {
		versionedKey := s.versionedStorageKey(bucket, key, obj.VersionID)
		if err := s.BaseStore.PutProto(versionedKey, ObjectToProto(obj)); err != nil {
			return err
		}
		latestKey := s.latestKeyStorageKey(bucket, key)
		if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
			_ = s.BaseStore.PutProto(versionedKey, previousProto)
			return err
		}
		return nil
	}
	storageKey := s.versionedStorageKey(bucket, key, "null")
	return s.BaseStore.PutProto(storageKey, ObjectToProto(obj))
}

// SetTags sets the tags for an object.
func (s *ObjectStore) SetTags(bucket, key string, tags []types.Tag) error {
	return s.updateObjectMetadata(bucket, key, func(obj *Object) {
		obj.Tags = tags
	})
}

// SetACL sets the access control list for an object.
func (s *ObjectStore) SetACL(bucket, key string, acp *AccessControlPolicy) error {
	return s.updateObjectMetadata(bucket, key, func(obj *Object) {
		obj.ACL = acp
	})
}

// GetACL retrieves the access control list for an object.
func (s *ObjectStore) GetACL(bucket, key string) (*AccessControlPolicy, error) {
	obj, err := s.GetMetadata(bucket, key)
	if err != nil {
		return nil, err
	}
	return obj.ACL, nil
}
