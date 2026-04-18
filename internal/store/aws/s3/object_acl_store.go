package s3

import (
	"vorpalstacks/internal/store/aws/common"
)

// SetTags sets the tags for an object.
func (s *ObjectStore) SetTags(bucket, key string, tags []common.Tag) error {
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
	obj.Tags = tags

	if s.isVersioningEnabled(bucket) {
		storageKey := s.versionedStorageKey(bucket, key, obj.VersionID)
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return err
		}
		latestKey := s.latestKeyStorageKey(bucket, key)
		return s.BaseStore.PutProto(latestKey, ObjectToProto(obj))
	}
	storageKey := s.storageKey(bucket, key)
	return s.BaseStore.PutProto(storageKey, ObjectToProto(obj))
}

// SetACL sets the access control list for an object.
func (s *ObjectStore) SetACL(bucket, key string, acp *AccessControlPolicy) error {
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
	obj.ACL = acp

	if s.isVersioningEnabled(bucket) {
		storageKey := s.versionedStorageKey(bucket, key, obj.VersionID)
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return err
		}
		latestKey := s.latestKeyStorageKey(bucket, key)
		return s.BaseStore.PutProto(latestKey, ObjectToProto(obj))
	}
	storageKey := s.storageKey(bucket, key)
	return s.BaseStore.PutProto(storageKey, ObjectToProto(obj))
}

// GetACL retrieves the access control list for an object.
func (s *ObjectStore) GetACL(bucket, key string) (*AccessControlPolicy, error) {
	obj, err := s.GetMetadata(bucket, key)
	if err != nil {
		return nil, err
	}
	return obj.ACL, nil
}
