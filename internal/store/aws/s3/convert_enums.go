package s3

import (
	pb "vorpalstacks/internal/pb/storage/storage_s3"
	common "vorpalstacks/internal/store/aws/common"
)

func objectStorageClassToProto(s ObjectStorageClass) pb.ObjectStorageClass {
	switch s {
	case StorageClassStandard:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD
	case StorageClassReducedRedundancy:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY
	case StorageClassGlacier:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_GLACIER
	case StorageClassStandardIA:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD_IA
	case StorageClassOneZoneIA:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_ONEZONE_IA
	case StorageClassIntelligentTiering:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING
	default:
		return pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_UNSPECIFIED
	}
}

func protoToObjectStorageClass(p pb.ObjectStorageClass) ObjectStorageClass {
	switch p {
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD:
		return StorageClassStandard
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY:
		return StorageClassReducedRedundancy
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_GLACIER:
		return StorageClassGlacier
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD_IA:
		return StorageClassStandardIA
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_ONEZONE_IA:
		return StorageClassOneZoneIA
	case pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING:
		return StorageClassIntelligentTiering
	default:
		return ""
	}
}

func bucketVersioningStatusToProto(s BucketVersioningStatus) pb.BucketVersioningStatus {
	switch s {
	case BucketVersioningEnabled:
		return pb.BucketVersioningStatus_BUCKET_VERSIONING_STATUS_ENABLED
	case BucketVersioningSuspended:
		return pb.BucketVersioningStatus_BUCKET_VERSIONING_STATUS_SUSPENDED
	default:
		return pb.BucketVersioningStatus_BUCKET_VERSIONING_STATUS_UNSPECIFIED
	}
}

func protoToBucketVersioningStatus(p pb.BucketVersioningStatus) BucketVersioningStatus {
	switch p {
	case pb.BucketVersioningStatus_BUCKET_VERSIONING_STATUS_ENABLED:
		return BucketVersioningEnabled
	case pb.BucketVersioningStatus_BUCKET_VERSIONING_STATUS_SUSPENDED:
		return BucketVersioningSuspended
	default:
		return ""
	}
}

func sseTypeToProto(s SSEType) pb.SSEType {
	switch s {
	case SSETypeAES256:
		return pb.SSEType_SSE_TYPE_AES256
	case SSETypeKMS:
		return pb.SSEType_SSE_TYPE_KMS
	case SSETypeKMSES:
		return pb.SSEType_SSE_TYPE_KMS_ES
	case SSETypeCustomer:
		return pb.SSEType_SSE_TYPE_CUSTOMER
	default:
		return pb.SSEType_SSE_TYPE_UNSPECIFIED
	}
}

func protoToSSEType(p pb.SSEType) SSEType {
	switch p {
	case pb.SSEType_SSE_TYPE_AES256:
		return SSETypeAES256
	case pb.SSEType_SSE_TYPE_KMS:
		return SSETypeKMS
	case pb.SSEType_SSE_TYPE_KMS_ES:
		return SSETypeKMSES
	case pb.SSEType_SSE_TYPE_CUSTOMER:
		return SSETypeCustomer
	default:
		return ""
	}
}

func objectLockRetentionModeToProto(m ObjectLockRetentionMode) pb.ObjectLockRetentionMode {
	switch m {
	case ObjectLockRetentionModeGovernance:
		return pb.ObjectLockRetentionMode_OBJECT_LOCK_RETENTION_MODE_GOVERNANCE
	case ObjectLockRetentionModeCompliance:
		return pb.ObjectLockRetentionMode_OBJECT_LOCK_RETENTION_MODE_COMPLIANCE
	default:
		return pb.ObjectLockRetentionMode_OBJECT_LOCK_RETENTION_MODE_UNSPECIFIED
	}
}

func protoToObjectLockRetentionMode(p pb.ObjectLockRetentionMode) ObjectLockRetentionMode {
	switch p {
	case pb.ObjectLockRetentionMode_OBJECT_LOCK_RETENTION_MODE_GOVERNANCE:
		return ObjectLockRetentionModeGovernance
	case pb.ObjectLockRetentionMode_OBJECT_LOCK_RETENTION_MODE_COMPLIANCE:
		return ObjectLockRetentionModeCompliance
	default:
		return ""
	}
}

func objectLockLegalHoldStatusToProto(s ObjectLockLegalHoldStatus) pb.ObjectLockLegalHoldStatus {
	switch s {
	case ObjectLockLegalHoldOn:
		return pb.ObjectLockLegalHoldStatus_OBJECT_LOCK_LEGAL_HOLD_STATUS_ON
	case ObjectLockLegalHoldOff:
		return pb.ObjectLockLegalHoldStatus_OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF
	default:
		return pb.ObjectLockLegalHoldStatus_OBJECT_LOCK_LEGAL_HOLD_STATUS_UNSPECIFIED
	}
}

func protoToObjectLockLegalHoldStatus(p pb.ObjectLockLegalHoldStatus) ObjectLockLegalHoldStatus {
	switch p {
	case pb.ObjectLockLegalHoldStatus_OBJECT_LOCK_LEGAL_HOLD_STATUS_ON:
		return ObjectLockLegalHoldOn
	case pb.ObjectLockLegalHoldStatus_OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF:
		return ObjectLockLegalHoldOff
	default:
		return ""
	}
}

func tagsToProto(t []common.Tag) []*pb.Tag {
	if t == nil {
		return nil
	}
	result := make([]*pb.Tag, len(t))
	for i, tag := range t {
		result[i] = &pb.Tag{Key: tag.Key, Value: tag.Value}
	}
	return result
}

func protoToTags(p []*pb.Tag) []common.Tag {
	if p == nil {
		return nil
	}
	result := make([]common.Tag, len(p))
	for i, tag := range p {
		result[i] = common.Tag{Key: tag.Key, Value: tag.Value}
	}
	return result
}

func tagToProtoPtr(t *common.Tag) *pb.Tag {
	if t == nil {
		return nil
	}
	return &pb.Tag{Key: t.Key, Value: t.Value}
}

func protoToTagPtr(p *pb.Tag) *common.Tag {
	if p == nil {
		return nil
	}
	return &common.Tag{Key: p.Key, Value: p.Value}
}

func int32PtrToInt32(p *int) int32 {
	if p == nil {
		return 0
	}
	return int32(*p)
}

func int32ToInt32Ptr(v int32) *int {
	i := int(v)
	return &i
}

func int64PtrToInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func int64ToInt64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func intPtrToInt32(p *int) int32 {
	if p == nil {
		return 0
	}
	return int32(*p)
}

func int32ToIntPtr(v int32) *int {
	if v == 0 {
		return nil
	}
	i := int(v)
	return &i
}

func boolToBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func boolToBoolPtr(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringToStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
