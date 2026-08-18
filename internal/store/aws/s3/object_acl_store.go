package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
	"vorpalstacks/internal/utils/aws/types"
)

// updateObjectMetadata is a shared helper for SetTags and SetACL. It acquires
// the key lock, reads the object metadata, applies the mutation callback, then
// writes back. With an empty versionId it targets the current object and, when
// versioning is enabled, writes both the versioned record and the _latest
// pointer; with an explicit versionId it targets that version only and keeps
// the _latest pointer in sync just when the targeted version is the latest.
func (s *ObjectStore) updateObjectMetadata(bucket, key, versionId string, mutate func(obj *Object)) error {
	lockKey := bucket + keySep + key
	s.keyLocker.Lock(lockKey)
	defer s.keyLocker.Unlock(lockKey)

	if versionId != "" {
		storageKey := s.versionedStorageKey(bucket, key, versionId)
		var pbObj pb.Object
		if err := s.BaseStore.GetProto(storageKey, &pbObj); err != nil {
			return ErrObjectNotFound
		}
		obj := ProtoToObject(&pbObj)
		mutate(obj)
		if err := s.BaseStore.PutProto(storageKey, ObjectToProto(obj)); err != nil {
			return err
		}
		if obj.IsLatest {
			latestKey := s.latestKeyStorageKey(bucket, key)
			if latestKey != storageKey {
				if err := s.BaseStore.PutProto(latestKey, ObjectToProto(obj)); err != nil {
					return err
				}
			}
		}
		return nil
	}

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

// SetTags sets the tags for an object. A non-empty versionId targets that
// specific object version instead of the current one.
func (s *ObjectStore) SetTags(bucket, key, versionId string, tags []types.Tag) error {
	return s.updateObjectMetadata(bucket, key, versionId, func(obj *Object) {
		obj.Tags = tags
	})
}

// SetACL sets the access control list for the current object.
func (s *ObjectStore) SetACL(bucket, key string, acp *AccessControlPolicy) error {
	return s.updateObjectMetadata(bucket, key, "", func(obj *Object) {
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
