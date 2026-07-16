package s3

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// DeleteObjectInput contains the input parameters for the DeleteObject operation.
type DeleteObjectInput struct {
	Bucket                    string
	Key                       string
	VersionId                 string
	BypassGovernanceRetention bool
}
type DeleteObjectOutput struct {
	DeleteMarker   bool
	VersionId      string
	RequestCharged string
}

// checkObjectLock verifies whether an object may be deleted under Object Lock
// rules. Returns nil when deletion is permitted, or an error describing why
// the operation is blocked.
//
// Legal hold ON always blocks deletion regardless of bypass.
// COMPLIANCE retention blocks deletion unconditionally.
// GOVERNANCE retention blocks deletion unless bypassGovernanceRetention is true.
func (o *ObjectOperations) checkObjectLock(stores *s3Stores, bucket *s3store.Bucket, objectBucket, key, versionId string, bypassGovernanceRetention bool) error {
	if !bucket.ObjectLockEnabled {
		return nil
	}

	obj, err := stores.objects.HeadWithVersion(context.Background(), objectBucket, key, versionId)
	if err != nil {
		return fmt.Errorf("failed to check Object Lock status: %w", err)
	}

	if obj.ObjectLockLegalHold != nil && obj.ObjectLockLegalHold.Status == s3store.ObjectLockLegalHoldOn {
		return ErrObjectLockedLegalHold
	}

	if obj.ObjectLockRetention != nil && obj.ObjectLockRetention.RetainUntilDate != nil {
		if obj.ObjectLockRetention.RetainUntilDate.After(time.Now()) {
			if obj.ObjectLockRetention.Mode == s3store.ObjectLockRetentionModeCompliance {
				return ErrObjectLockedRetention
			}
			if obj.ObjectLockRetention.Mode == s3store.ObjectLockRetentionModeGovernance && !bypassGovernanceRetention {
				return ErrObjectLockedRetention
			}
		}
	}

	return nil
}

// DeleteObject deletes an object from S3.
func (o *ObjectOperations) DeleteObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectInput) (*DeleteObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}
	if err := o.checkObjectLock(stores, bucket, input.Bucket, input.Key, input.VersionId, input.BypassGovernanceRetention); err != nil {
		return nil, err
	}

	marker, err := stores.objects.DeleteWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	output := &DeleteObjectOutput{}
	if marker != nil {
		output.DeleteMarker = true
		output.VersionId = marker.VersionID
		o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, 0, marker.VersionID, "", eventbus.S3ObjectRemovedDeleteMarkerCreated)
		if bucket.ReplicationConfiguration != nil {
			dmCtx, dmCancel := context.WithTimeout(context.Background(), 30*time.Second)
			go func() {
				defer dmCancel()
				o.svc.replicateDeleteMarker(dmCtx, reqCtx, stores, bucket, input.Key)
			}()
		}
	} else {
		o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
	}

	return output, nil
}

// DeleteObjectsInput contains the input parameters for the DeleteObjects operation.
type DeleteObjectsInput struct {
	Bucket                    string
	Delete                    *Delete
	BypassGovernanceRetention bool
}

// Delete contains the objects to delete.
type Delete struct {
	Objects []ObjectIdentifier `xml:"Object"`
	Quiet   bool               `xml:"Quiet"`
}

// ObjectIdentifier identifies a specific object to delete.
type ObjectIdentifier struct {
	Key       string `xml:"Key"`
	VersionId string `xml:"VersionId,omitempty"`
}

// DeleteObjectsOutput contains the output from the DeleteObjects operation.
type DeleteObjectsOutput struct {
	Deleted []DeletedObject `xml:"Deleted"`
	Error   []DeleteError   `xml:"Error"`
}

// DeletedObject contains information about a deleted object.
type DeletedObject struct {
	Key            string `xml:"Key"`
	VersionId      string `xml:"VersionId,omitempty"`
	DeleteMarker   bool   `xml:"DeleteMarker,omitempty"`
	DeleteMarkerId string `xml:"DeleteMarkerVersionId,omitempty"`
}

// DeleteError contains information about a delete error.
type DeleteError struct {
	Key     string `xml:"Key"`
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

// DeleteObjects deletes multiple objects from S3 in a single request.
func (o *ObjectOperations) DeleteObjects(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectsInput) (*DeleteObjectsOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	var deleted []DeletedObject
	var errors []DeleteError

	for _, obj := range input.Delete.Objects {
		if err := o.checkObjectLock(stores, bucket, input.Bucket, obj.Key, obj.VersionId, input.BypassGovernanceRetention); err != nil {
			errors = append(errors, DeleteError{
				Key:     obj.Key,
				Code:    "AccessDenied",
				Message: err.Error(),
			})
			continue
		}

		marker, err := stores.objects.DeleteWithVersion(ctx, input.Bucket, obj.Key, obj.VersionId)
		if err != nil {
			errors = append(errors, DeleteError{
				Key:     obj.Key,
				Code:    "InternalError",
				Message: err.Error(),
			})
		} else {
			deletedObj := DeletedObject{
				Key: obj.Key,
			}
			if marker != nil {
				deletedObj.DeleteMarker = true
				deletedObj.DeleteMarkerId = marker.VersionID
				deletedObj.VersionId = marker.VersionID
				o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, marker.VersionID, "", eventbus.S3ObjectRemovedDeleteMarkerCreated)
				if bucket.ReplicationConfiguration != nil {
					dmCtx, dmCancel := context.WithTimeout(context.Background(), 30*time.Second)
					keyVal := obj.Key
					go func() {
						defer dmCancel()
						o.svc.replicateDeleteMarker(dmCtx, reqCtx, stores, bucket, keyVal)
					}()
				}
			} else if obj.VersionId != "" {
				deletedObj.VersionId = obj.VersionId
				o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
			} else {
				o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
			}
			deleted = append(deleted, deletedObj)
		}
	}

	return &DeleteObjectsOutput{
		Deleted: deleted,
		Error:   errors,
	}, nil
}
