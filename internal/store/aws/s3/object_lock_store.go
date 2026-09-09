package s3

import (
	"context"
)

// SetObjectLegalHold sets the legal hold status for an object version.
// The mutation skeleton keeps the versioned record and the _latest pointer
// in sync in one atomic batch.
func (s *ObjectStore) SetObjectLegalHold(ctx context.Context, bucket, key, versionId string, legalHold *ObjectLockLegalHold) error {
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.ObjectLockLegalHold = legalHold
		return nil
	})
}

// GetObjectLegalHold retrieves the legal hold status for an object version.
func (s *ObjectStore) GetObjectLegalHold(ctx context.Context, bucket, key, versionId string) (*ObjectLockLegalHold, error) {
	obj, err := s.getVersionedObjectMeta(bucket, key, versionId)
	if err != nil {
		return nil, err
	}

	if obj.ObjectLockLegalHold == nil {
		return &ObjectLockLegalHold{Status: ObjectLockLegalHoldOff}, nil
	}

	return obj.ObjectLockLegalHold, nil
}

// SetObjectRetention sets the retention policy for an object version.
func (s *ObjectStore) SetObjectRetention(ctx context.Context, bucket, key, versionId string, retention *ObjectLockRetention) error {
	return s.mutateObjectRecord(bucket, key, versionId, func(obj *Object) error {
		obj.ObjectLockRetention = retention
		return nil
	})
}

// GetObjectRetention retrieves the retention policy for an object version.
func (s *ObjectStore) GetObjectRetention(ctx context.Context, bucket, key, versionId string) (*ObjectLockRetention, error) {
	obj, err := s.getVersionedObjectMeta(bucket, key, versionId)
	if err != nil {
		return nil, err
	}

	if obj.ObjectLockRetention == nil {
		return nil, ErrRetentionNotFound
	}

	return obj.ObjectLockRetention, nil
}
