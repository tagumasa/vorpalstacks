package s3

// Core functions for the object sub-resource operations (ACL, tagging, Object
// Lock legal hold / retention, versioning list, restore, attributes and the
// single/batch delete entry points). The operation methods on
// ObjectOperations are thin adapters that receive the per-region store bundle
// and delegate all validation and persistence here. The shared helpers
// (validateBucketExists, checkObjectLock, resolveTaggingTarget) also live here
// so the guard walk stops at the single validation path.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// validateBucketExists reports ErrNoSuchBucket when the named bucket does not
// exist in the store bundle's region.
func (s *S3Service) validateBucketExists(stores *s3Stores, bucket string) error {
	if !stores.buckets.Exists(bucket) {
		return ErrNoSuchBucket
	}
	return nil
}

// checkObjectLock verifies whether an object may be deleted under Object Lock
// rules. Returns nil when deletion is permitted, or an error describing why
// the operation is blocked.
//
// Legal hold ON always blocks deletion regardless of bypass.
// COMPLIANCE retention blocks deletion unconditionally.
// GOVERNANCE retention blocks deletion unless bypassGovernanceRetention is true.
func (s *S3Service) checkObjectLock(stores *s3Stores, bucket *s3store.Bucket, objectBucket, key, versionId string, bypassGovernanceRetention bool) error {
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

// resolveTaggingTarget returns the object record a tagging operation
// targets: the specific version when versionId is set, the current object
// otherwise. Both tagging reads and writes go through it so that a
// nonexistent object and a nonexistent version are reported uniformly.
func (s *S3Service) resolveTaggingTarget(ctx context.Context, stores *s3Stores, bucket, key, versionId string) (*s3store.Object, error) {
	if versionId != "" {
		return stores.objects.HeadWithVersion(ctx, bucket, key, versionId)
	}
	return stores.objects.GetMetadata(bucket, key)
}

// getObjectAclCore retrieves the Access Control List for an object and
// returns the owner and list of grants for the specified object version.
func (s *S3Service) getObjectAclCore(stores *s3Stores, bucket, key, versionId string) (*GetObjectAclOutput, error) {
	if err := s.validateBucketExists(stores, bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(key); err != nil {
		return nil, err
	}

	acp, err := stores.objects.GetACLWithVersion(bucket, key, versionId)
	if err != nil {
		return nil, mapVersionLookupError(err, versionId)
	}

	owner := &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID}

	if acp == nil {
		return &GetObjectAclOutput{
			Owner: owner,
			Grants: []*s3store.Grant{
				{
					Grantee:    &s3store.Grantee{Type: s3store.GranteeTypeCanonicalUser, ID: s.accountID, DisplayName: s.accountID},
					Permission: s3store.PermissionFullControl,
				},
			},
		}, nil
	}

	return &GetObjectAclOutput{
		Owner:  acp.Owner,
		Grants: acp.Grants,
	}, nil
}

// putObjectAclCore sets the Access Control List for an object. Accepts either
// a canned ACL string, an AccessControlPolicy, or individual grant headers.
func (s *S3Service) putObjectAclCore(ctx context.Context, stores *s3Stores, input *PutObjectAclInput) error {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	owner := &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID}

	var acp *s3store.AccessControlPolicy
	var err error

	if input.ACL != "" {
		acp, err = CannedACLToPolicy(input.ACL, owner)
		if err != nil {
			return err
		}
	} else if input.AccessControlPolicy != nil {
		acp = input.AccessControlPolicy
	} else {
		grants, err := ParseGrantHeaders(input.GrantFullControl, input.GrantRead, input.GrantReadACP, input.GrantWrite, input.GrantWriteACP)
		if err != nil {
			return NewInvalidArgumentError(err.Error())
		}
		if len(grants) > 0 {
			acp = &s3store.AccessControlPolicy{Owner: owner, Grants: grants}
		} else {
			return NewInvalidArgumentError("missing required ACL specification")
		}
	}

	publicAccessBlock, _ := stores.buckets.GetPublicAccessBlock(input.Bucket)
	if publicAccessBlock != nil && publicAccessBlock.BlockPublicAcls {
		if isPublicCannedACL(input.ACL) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
		if acpContainsPublicAccess(acp) {
			return NewInvalidArgumentError("bucket has BlockPublicAcls enabled")
		}
	}

	// With Object Ownership set to BucketOwnerEnforced, "requests to set or
	// update ACLs fail" with AccessControlListNotSupported.
	if aclsDisabled, _ := s.bucketACLsDisabled(ctx, stores, input.Bucket); aclsDisabled {
		return ErrAccessControlListNotSupported
	}

	if err := stores.objects.SetACLWithVersion(input.Bucket, input.Key, input.VersionId, acp); err != nil {
		return mapVersionLookupError(err, input.VersionId)
	}
	return nil
}

// deleteObjectOpCore deletes an object from S3 under Object Lock governance,
// publishing the removal notification and propagating delete markers to
// replication destinations.
func (s *S3Service) deleteObjectOpCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectInput) (*DeleteObjectOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}
	if err := s.checkObjectLock(stores, bucket, input.Bucket, input.Key, input.VersionId, input.BypassGovernanceRetention); err != nil {
		return nil, err
	}

	coreResult, err := s.deleteObjectCore(ctx, stores.objects, AdminDeleteObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionID: input.VersionId,
	})
	if err != nil {
		return nil, err
	}

	output := &DeleteObjectOutput{}
	if coreResult.IsDeleteMarker {
		output.DeleteMarker = true
		output.VersionId = coreResult.VersionID
		s.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, 0, coreResult.VersionID, "", eventbus.S3ObjectRemovedDeleteMarkerCreated)
		if bucket.ReplicationConfiguration != nil {
			s.goReplicationWorker(input.Bucket, input.Key, "delete-marker", func(dmCtx context.Context) {
				s.replicateDeleteMarker(dmCtx, reqCtx, stores, bucket, input.Key)
			})
		}
	} else {
		s.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
	}

	return output, nil
}

// deleteObjectsOpCore deletes multiple objects from S3 in a single request,
// applying Object Lock governance per entry and reporting per-entry results.
func (s *S3Service) deleteObjectsOpCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectsInput) (*DeleteObjectsOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	var deleted []DeletedObject
	var errors []DeleteError

	for _, obj := range input.Delete.Objects {
		if err := s.checkObjectLock(stores, bucket, input.Bucket, obj.Key, obj.VersionId, input.BypassGovernanceRetention); err != nil {
			errors = append(errors, DeleteError{
				Key:     obj.Key,
				Code:    "AccessDenied",
				Message: err.Error(),
			})
			continue
		}

		coreResult, err := s.deleteObjectCore(ctx, stores.objects, AdminDeleteObjectInput{
			Bucket:    input.Bucket,
			Key:       obj.Key,
			VersionID: obj.VersionId,
		})
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
			if coreResult.IsDeleteMarker {
				deletedObj.DeleteMarker = true
				deletedObj.DeleteMarkerId = coreResult.VersionID
				deletedObj.VersionId = coreResult.VersionID
				s.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, coreResult.VersionID, "", eventbus.S3ObjectRemovedDeleteMarkerCreated)
				if bucket.ReplicationConfiguration != nil {
					keyVal := obj.Key
					s.goReplicationWorker(input.Bucket, keyVal, "delete-marker", func(dmCtx context.Context) {
						s.replicateDeleteMarker(dmCtx, reqCtx, stores, bucket, keyVal)
					})
				}
			} else if obj.VersionId != "" {
				deletedObj.VersionId = obj.VersionId
				s.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
			} else {
				s.publishObjectNotification(ctx, reqCtx, input.Bucket, obj.Key, 0, "", "", eventbus.S3ObjectRemovedDelete)
			}
			deleted = append(deleted, deletedObj)
		}
	}

	return &DeleteObjectsOutput{
		Deleted: deleted,
		Error:   errors,
	}, nil
}

// getObjectAttributesCore retrieves the requested attributes of an object,
// honouring the ObjectAttributes part-window parameters.
func (s *S3Service) getObjectAttributesCore(ctx context.Context, stores *s3Stores, input *GetObjectAttributesInput) (*GetObjectAttributesOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	obj, err := stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, mapVersionLookupError(err, input.VersionId)
	}

	objectSize := obj.Size
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		objectSize = obj.SSEMetadata.UnencryptedSize
	}

	output := &GetObjectAttributesOutput{
		VersionId:    obj.VersionID,
		ETag:         formatETag(obj.ETag),
		ObjectSize:   objectSize,
		StorageClass: string(obj.StorageClass),
		LastModified: s3Timestamp(obj.LastModified),
	}

	for _, attr := range input.ObjectAttributes {
		switch attr {
		case "ETag":
			output.ETag = formatETag(obj.ETag)
		case "ObjectSize":
			output.ObjectSize = objectSize
		case "StorageClass":
			output.StorageClass = string(obj.StorageClass)
		case "ObjectParts":
			if obj.SSEMetadata != nil && len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
				partInfos := obj.SSEMetadata.PartEncryptionInfos
				totalParts := int32(len(partInfos))

				partNumberStart := int32(0)
				if input.PartNumberMarker != "" {
					if parsed, pErr := strconv.ParseInt(input.PartNumberMarker, 10, 32); pErr == nil && parsed > 0 {
						partNumberStart = int32(parsed)
					}
				}

				maxParts := input.MaxParts
				if maxParts <= 0 {
					maxParts = s3MaxParts
				}

				var filteredParts []GetObjectAttributesPart
				for i, pi := range partInfos {
					pn := int32(i + 1)
					if pn <= partNumberStart {
						continue
					}
					if int32(len(filteredParts)) >= maxParts {
						break
					}
					filteredParts = append(filteredParts, GetObjectAttributesPart{
						PartNumber: pn,
						Size:       pi.PlainSize,
					})
				}

				isTruncated := int32(len(partInfos)) > partNumberStart+int32(len(filteredParts))
				var nextMarker string
				if isTruncated && len(filteredParts) > 0 {
					nextMarker = strconv.FormatInt(int64(filteredParts[len(filteredParts)-1].PartNumber), 10)
				}

				output.ObjectParts = &GetObjectAttributesParts{
					IsTruncated:          isTruncated,
					MaxParts:             maxParts,
					NextPartNumberMarker: nextMarker,
					PartNumberMarker:     input.PartNumberMarker,
					Parts:                filteredParts,
					TotalPartsCount:      totalParts,
				}
			}
		case "Checksum":
			output.Checksum = &GetObjectAttributesChecksum{}
		}
	}

	return output, nil
}

// putObjectLegalHoldCore applies a legal hold to an object. The bucket must
// have Object Lock enabled; legal hold prevents object deletion or
// overwriting.
func (s *S3Service) putObjectLegalHoldCore(ctx context.Context, stores *s3Stores, input *PutObjectLegalHoldInput) error {
	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	if input.LegalHold == nil {
		return NewInvalidArgumentError("LegalHold is required")
	}

	var status s3store.ObjectLockLegalHoldStatus
	switch input.LegalHold.Status {
	case "ON":
		status = s3store.ObjectLockLegalHoldOn
	case "OFF":
		status = s3store.ObjectLockLegalHoldOff
	default:
		return NewInvalidArgumentError(fmt.Sprintf("invalid legal hold status: %s (must be ON or OFF)", input.LegalHold.Status))
	}

	legalHold := &s3store.ObjectLockLegalHold{Status: status}

	return stores.objects.SetObjectLegalHold(ctx, input.Bucket, input.Key, input.VersionId, legalHold)
}

// getObjectLegalHoldCore retrieves the current legal hold configuration for
// the specified object version.
func (s *S3Service) getObjectLegalHoldCore(ctx context.Context, stores *s3Stores, input *GetObjectLegalHoldInput) (*GetObjectLegalHoldOutput, error) {
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrObjectLockNotEnabled
	}

	legalHold, err := stores.objects.GetObjectLegalHold(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	return &GetObjectLegalHoldOutput{
		LegalHold: &LegalHoldOutput{
			Status: string(legalHold.Status),
		},
	}, nil
}

// putObjectRetentionCore applies a retention period to an object. The bucket
// must have Object Lock enabled; COMPLIANCE mode prevents deletion until
// retention expires.
func (s *S3Service) putObjectRetentionCore(ctx context.Context, stores *s3Stores, input *PutObjectRetentionInput) error {
	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return err
	}

	if !bucket.ObjectLockEnabled {
		return ErrObjectLockNotEnabled
	}

	if input.Retention == nil {
		return NewInvalidArgumentError("retention is required")
	}

	if err := validateRetentionMode(input.Retention.Mode); err != nil {
		return err
	}

	if err := validateRetainUntilDate(input.Retention.RetainUntilDate); err != nil {
		return err
	}

	retention := &s3store.ObjectLockRetention{
		Mode:            s3store.ObjectLockRetentionMode(input.Retention.Mode),
		RetainUntilDate: input.Retention.RetainUntilDate,
	}

	return stores.objects.SetObjectRetention(ctx, input.Bucket, input.Key, input.VersionId, retention)
}

// getObjectRetentionCore retrieves the retention mode and retain-until date
// for the specified object version.
func (s *S3Service) getObjectRetentionCore(ctx context.Context, stores *s3Stores, input *GetObjectRetentionInput) (*GetObjectRetentionOutput, error) {
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	if !bucket.ObjectLockEnabled {
		return nil, ErrObjectLockNotEnabled
	}

	retention, err := stores.objects.GetObjectRetention(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	return &GetObjectRetentionOutput{
		Retention: &RetentionOutput{
			Mode:            string(retention.Mode),
			RetainUntilDate: retention.RetainUntilDate,
		},
	}, nil
}

// listObjectVersionsCore lists the object versions and delete markers of a
// bucket honouring the prefix, delimiter and marker pagination parameters.
func (s *S3Service) listObjectVersionsCore(stores *s3Stores, input *ListObjectVersionsInput) (*ListObjectVersionsOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	result, err := stores.objects.ListObjectVersions(input.Bucket, input.Prefix, input.Delimiter, input.KeyMarker, input.VersionIdMarker, input.MaxKeys)
	if err != nil {
		return nil, err
	}

	var versions []*ObjectVersion
	var deleteMarkers []*DeleteMarkerEntry
	var commonPrefixes []CommonPrefix

	for _, obj := range result.Objects {
		if obj.IsDeleteMarker {
			deleteMarkers = append(deleteMarkers, &DeleteMarkerEntry{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				VersionId:    obj.VersionID,
				IsLatest:     obj.IsLatest,
			})
		} else {
			versions = append(versions, &ObjectVersion{
				Key:          obj.Key,
				LastModified: obj.LastModified,
				ETag:         formatETag(obj.ETag),
				Size:         obj.Size,
				StorageClass: string(obj.StorageClass),
				VersionId:    obj.VersionID,
				IsLatest:     obj.IsLatest,
			})
		}
	}

	for _, prefix := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, CommonPrefix{Prefix: prefix})
	}

	return &ListObjectVersionsOutput{
		Versions:            versions,
		DeleteMarkers:       deleteMarkers,
		CommonPrefixes:      commonPrefixes,
		Delimiter:           input.Delimiter,
		EncodingType:        input.EncodingType,
		IsTruncated:         result.IsTruncated,
		KeyMarker:           input.KeyMarker,
		MaxKeys:             input.MaxKeys,
		Name:                input.Bucket,
		NextKeyMarker:       result.NextVersionKeyMarker,
		NextVersionIdMarker: result.NextVersionIDMarker,
		Prefix:              input.Prefix,
		VersionIdMarker:     input.VersionIdMarker,
	}, nil
}

// restoreObjectCore creates or extends the temporary restored copy of an
// archived object and reports whether a restored copy already existed (the
// request then only extended its expiry, which the API answers with 200 OK
// instead of 202 Accepted). The object's storage class never changes; the
// restored copy's expiry is rounded up to the following midnight UTC.
func (s *S3Service) restoreObjectCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *RestoreObjectInput) (bool, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return false, err
	}

	var obj *s3store.Object
	var err error
	if input.VersionId != "" {
		obj, err = stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	} else {
		obj, err = stores.objects.Head(ctx, input.Bucket, input.Key)
	}
	if err != nil {
		return false, versionLookupError(input.Key, input.VersionId)
	}

	if !isArchiveClass(obj.StorageClass) {
		return false, ErrInvalidObjectState
	}

	restoreDays := 1
	if input.Body != nil {
		var restoreReq RestoreRequest
		if err := request.NewSafeXMLDecoder(input.Body).Decode(&restoreReq); err != nil {
			return false, NewInvalidArgumentError("invalid RestoreObject request body")
		}
		if err := validateRestoreDays(restoreReq.Days); err != nil {
			return false, err
		}
		restoreDays = restoreReq.Days
	}

	now := time.Now()
	alreadyRestored := objectRestored(obj, now)
	expiry := nextRestoreExpiry(now, restoreDays)
	if err := stores.objects.SetRestoreState(input.Bucket, input.Key, input.VersionId, &expiry); err != nil {
		return false, err
	}

	// The restore completes synchronously in this implementation, so the
	// initiation and completion notifications are published together.
	s.publishRestoreNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectRestorePost, expiry)
	s.publishRestoreNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectRestoreCompleted, expiry)

	logs.Info("s3: object restored",
		logs.String("bucket", input.Bucket),
		logs.String("key", input.Key),
		logs.String("days", fmt.Sprintf("%d", restoreDays)))

	return alreadyRestored, nil
}

// putObjectTaggingCore replaces all tags on an object with the specified tag
// set. When input.VersionId is set, only that version's tag set is replaced.
func (s *S3Service) putObjectTaggingCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectTaggingInput) error {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	obj, err := s.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return versionLookupError(input.Key, input.VersionId)
	}

	if err := validateTags(input.Tags); err != nil {
		return err
	}

	if err := stores.objects.SetTags(input.Bucket, input.Key, input.VersionId, TagsToCommon(input.Tags)); err != nil {
		return err
	}

	s.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectTaggingPut)
	return nil
}

// getObjectTaggingCore retrieves all tags associated with an object. When
// input.VersionId is set, the tags of that version are returned.
func (s *S3Service) getObjectTaggingCore(ctx context.Context, stores *s3Stores, input *GetObjectTaggingInput) (*GetObjectTaggingOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	obj, err := s.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, versionLookupError(input.Key, input.VersionId)
	}

	return &GetObjectTaggingOutput{TagSet: CommonToTags(obj.Tags)}, nil
}

// deleteObjectTaggingCore removes all tags from an object. When
// input.VersionId is set, only that version's tag set is removed.
func (s *S3Service) deleteObjectTaggingCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *DeleteObjectTaggingInput) error {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	obj, err := s.resolveTaggingTarget(ctx, stores, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return versionLookupError(input.Key, input.VersionId)
	}

	if err := stores.objects.SetTags(input.Bucket, input.Key, input.VersionId, nil); err != nil {
		return err
	}

	s.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectTaggingDelete)
	return nil
}
