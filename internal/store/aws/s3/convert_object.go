package s3

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "vorpalstacks/internal/pb/storage/storage_s3"
)

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func protoToTime(p *timestamppb.Timestamp) *time.Time {
	if p == nil {
		return nil
	}
	t := p.AsTime()
	return &t
}

// ObjectToProto converts an internal Object to a protobuf Object.
// Returns nil if the input object is nil.
func ObjectToProto(o *Object) *pb.Object {
	if o == nil {
		return nil
	}
	return &pb.Object{
		Key:                  o.Key,
		BucketName:           o.BucketName,
		Size:                 o.Size,
		Etag:                 o.ETag,
		LastModified:         timestamppb.New(o.LastModified),
		ContentType:          o.ContentType,
		ContentEncoding:      o.ContentEncoding,
		ContentLanguage:      o.ContentLanguage,
		ContentDisposition:   o.ContentDisposition,
		CacheControl:         o.CacheControl,
		Expires:              timeToProto(o.Expires),
		Metadata:             o.Metadata,
		StorageClass:         objectStorageClassToProto(o.StorageClass),
		ServerSideEncryption: o.ServerSideEncryption,
		SseKmsKeyId:          o.SSEKMSKeyID,
		VersionId:            o.VersionID,
		IsLatest:             o.IsLatest,
		IsDeleteMarker:       o.IsDeleteMarker,
		Tags:                 tagsToProto(o.Tags),
		Acl:                  accessControlPolicyToProto(o.ACL),
		Owner:                aclOwnerToProto(o.Owner),
		ObjectLockLegalHold:  objectLockLegalHoldToProto(o.ObjectLockLegalHold),
		ObjectLockRetention:  objectLockRetentionToProto(o.ObjectLockRetention),
		SseMetadata:          sseObjectMetadataToProto(o.SSEMetadata),
	}
}

// ProtoToObject converts a protobuf Object to an internal Object.
// Returns nil if the input protobuf object is nil.
func ProtoToObject(p *pb.Object) *Object {
	if p == nil {
		return nil
	}
	return &Object{
		Key:                  p.Key,
		BucketName:           p.BucketName,
		Size:                 p.Size,
		ETag:                 p.Etag,
		LastModified:         p.LastModified.AsTime(),
		ContentType:          p.ContentType,
		ContentEncoding:      p.ContentEncoding,
		ContentLanguage:      p.ContentLanguage,
		ContentDisposition:   p.ContentDisposition,
		CacheControl:         p.CacheControl,
		Expires:              protoToTime(p.Expires),
		Metadata:             p.Metadata,
		StorageClass:         protoToObjectStorageClass(p.StorageClass),
		ServerSideEncryption: p.ServerSideEncryption,
		SSEKMSKeyID:          p.SseKmsKeyId,
		VersionID:            p.VersionId,
		IsLatest:             p.IsLatest,
		IsDeleteMarker:       p.IsDeleteMarker,
		Tags:                 protoToTags(p.Tags),
		ACL:                  protoToAccessControlPolicy(p.Acl),
		Owner:                protoToACLOwner(p.Owner),
		ObjectLockLegalHold:  protoToObjectLockLegalHold(p.ObjectLockLegalHold),
		ObjectLockRetention:  protoToObjectLockRetention(p.ObjectLockRetention),
		SSEMetadata:          protoToSSEObjectMetadata(p.SseMetadata),
	}
}

// MultipartUploadToProto converts an internal MultipartUpload to a protobuf MultipartUpload.
// Returns nil if the input upload is nil.
func MultipartUploadToProto(u *MultipartUpload) *pb.MultipartUpload {
	if u == nil {
		return nil
	}
	return &pb.MultipartUpload{
		UploadId:         u.UploadID,
		Key:              u.Key,
		BucketName:       u.BucketName,
		Initiated:        timestamppb.New(u.Initiated),
		StorageClass:     objectStorageClassToProto(u.StorageClass),
		Owner:            u.Owner,
		Initiator:        u.Initiator,
		Parts:            objectPartsToProto(u.Parts),
		Metadata:         u.Metadata,
		ContentType:      u.ContentType,
		SseMetadata:      sseObjectMetadataToProto(u.SSEMetadata),
		SseType:          sseTypeToProto(u.SSEType),
		KmsKeyId:         u.KMSKeyID,
		CustomerKeyMd5:   u.CustomerKeyMD5,
		PlaintextDataKey: u.PlaintextDataKey,
	}
}

// ProtoToMultipartUpload converts a protobuf MultipartUpload to an internal MultipartUpload.
// Returns nil if the input protobuf upload is nil.
func ProtoToMultipartUpload(p *pb.MultipartUpload) *MultipartUpload {
	if p == nil {
		return nil
	}
	return &MultipartUpload{
		UploadID:         p.UploadId,
		Key:              p.Key,
		BucketName:       p.BucketName,
		Initiated:        p.Initiated.AsTime(),
		StorageClass:     protoToObjectStorageClass(p.StorageClass),
		Owner:            p.Owner,
		Initiator:        p.Initiator,
		Parts:            protoToObjectParts(p.Parts),
		Metadata:         p.Metadata,
		ContentType:      p.ContentType,
		SSEMetadata:      protoToSSEObjectMetadata(p.SseMetadata),
		SSEType:          protoToSSEType(p.SseType),
		KMSKeyID:         p.KmsKeyId,
		CustomerKeyMD5:   p.CustomerKeyMd5,
		PlaintextDataKey: p.PlaintextDataKey,
	}
}

func objectPartToProto(p ObjectPart) *pb.ObjectPart {
	return &pb.ObjectPart{
		PartNumber:    int32(p.PartNumber),
		Etag:          p.ETag,
		Size:          p.Size,
		LastModified:  timestamppb.New(p.LastModified),
		EncryptedSize: p.EncryptedSize,
		ContentNonce:  p.ContentNonce,
		DataKey:       p.DataKey,
	}
}

func protoToObjectPart(p *pb.ObjectPart) ObjectPart {
	return ObjectPart{
		PartNumber:    int(p.PartNumber),
		ETag:          p.Etag,
		Size:          p.Size,
		LastModified:  p.LastModified.AsTime(),
		EncryptedSize: p.EncryptedSize,
		ContentNonce:  p.ContentNonce,
		DataKey:       p.DataKey,
	}
}

func objectPartsToProto(parts []ObjectPart) []*pb.ObjectPart {
	if parts == nil {
		return nil
	}
	result := make([]*pb.ObjectPart, len(parts))
	for i, p := range parts {
		result[i] = objectPartToProto(p)
	}
	return result
}

func protoToObjectParts(p []*pb.ObjectPart) []ObjectPart {
	if p == nil {
		return nil
	}
	result := make([]ObjectPart, len(p))
	for i, part := range p {
		result[i] = protoToObjectPart(part)
	}
	return result
}

func encryptionConfigToProto(e *EncryptionConfig) *pb.EncryptionConfig {
	if e == nil {
		return nil
	}
	result := &pb.EncryptionConfig{
		SseAlgorithm:   e.SSEAlgorithm,
		KmsMasterKeyId: e.KMSMasterKeyID,
	}
	if e.BucketKeyEnabled != nil {
		result.BucketKeyEnabled = &wrapperspb.BoolValue{Value: *e.BucketKeyEnabled}
	}
	return result
}

func protoToEncryptionConfig(p *pb.EncryptionConfig) *EncryptionConfig {
	if p == nil {
		return nil
	}
	result := &EncryptionConfig{
		SSEAlgorithm:   p.SseAlgorithm,
		KMSMasterKeyID: p.KmsMasterKeyId,
	}
	if p.BucketKeyEnabled != nil {
		result.BucketKeyEnabled = &p.BucketKeyEnabled.Value
	}
	return result
}

func sseObjectMetadataToProto(m *SSEObjectMetadata) *pb.SSEObjectMetadata {
	if m == nil {
		return nil
	}
	var partInfos []*pb.PartEncryptionInfo
	for _, pi := range m.PartEncryptionInfos {
		partInfos = append(partInfos, &pb.PartEncryptionInfo{
			EncryptedSize: pi.EncryptedSize,
			PlainSize:     pi.PlainSize,
			ContentNonce:  pi.ContentNonce,
			DataKey:       pi.DataKey,
		})
	}
	return &pb.SSEObjectMetadata{
		EncryptionType:      sseTypeToProto(m.EncryptionType),
		EncryptedDataKey:    m.EncryptedDataKey,
		ContentNonce:        m.ContentNonce,
		KmsKeyId:            m.KMSKeyID,
		EncryptionContext:   m.EncryptionContext,
		UnencryptedMd5:      m.UnencryptedMD5,
		UnencryptedSize:     m.UnencryptedSize,
		PartEncryptionInfos: partInfos,
	}
}

func protoToSSEObjectMetadata(p *pb.SSEObjectMetadata) *SSEObjectMetadata {
	if p == nil {
		return nil
	}
	var partInfos []PartEncryptionInfo
	for _, pi := range p.PartEncryptionInfos {
		partInfos = append(partInfos, PartEncryptionInfo{
			EncryptedSize: pi.EncryptedSize,
			PlainSize:     pi.PlainSize,
			ContentNonce:  pi.ContentNonce,
			DataKey:       pi.DataKey,
		})
	}
	return &SSEObjectMetadata{
		EncryptionType:      protoToSSEType(p.EncryptionType),
		EncryptedDataKey:    p.EncryptedDataKey,
		ContentNonce:        p.ContentNonce,
		KMSKeyID:            p.KmsKeyId,
		EncryptionContext:   p.EncryptionContext,
		UnencryptedMD5:      p.UnencryptedMd5,
		UnencryptedSize:     p.UnencryptedSize,
		PartEncryptionInfos: partInfos,
	}
}

func objectLockLegalHoldToProto(h *ObjectLockLegalHold) *pb.ObjectLockLegalHold {
	if h == nil {
		return nil
	}
	return &pb.ObjectLockLegalHold{
		Status: objectLockLegalHoldStatusToProto(h.Status),
	}
}

func protoToObjectLockLegalHold(p *pb.ObjectLockLegalHold) *ObjectLockLegalHold {
	if p == nil {
		return nil
	}
	return &ObjectLockLegalHold{
		Status: protoToObjectLockLegalHoldStatus(p.Status),
	}
}

func objectLockRetentionToProto(r *ObjectLockRetention) *pb.ObjectLockRetention {
	if r == nil {
		return nil
	}
	return &pb.ObjectLockRetention{
		Mode:            objectLockRetentionModeToProto(r.Mode),
		RetainUntilDate: timeToProto(r.RetainUntilDate),
	}
}

func protoToObjectLockRetention(p *pb.ObjectLockRetention) *ObjectLockRetention {
	if p == nil {
		return nil
	}
	return &ObjectLockRetention{
		Mode:            protoToObjectLockRetentionMode(p.Mode),
		RetainUntilDate: protoToTime(p.RetainUntilDate),
	}
}

func objectLockConfigurationToProto(c *ObjectLockConfiguration) *pb.ObjectLockConfiguration {
	if c == nil {
		return nil
	}
	return &pb.ObjectLockConfiguration{
		ObjectLockEnabled: c.ObjectLockEnabled,
		Rule:              objectLockRuleToProto(c.Rule),
	}
}

func protoToObjectLockConfiguration(p *pb.ObjectLockConfiguration) *ObjectLockConfiguration {
	if p == nil {
		return nil
	}
	return &ObjectLockConfiguration{
		ObjectLockEnabled: p.ObjectLockEnabled,
		Rule:              protoToObjectLockRule(p.Rule),
	}
}

func objectLockRuleToProto(r *ObjectLockRule) *pb.ObjectLockRule {
	if r == nil {
		return nil
	}
	return &pb.ObjectLockRule{
		DefaultRetention: defaultRetentionToProto(r.DefaultRetention),
	}
}

func protoToObjectLockRule(p *pb.ObjectLockRule) *ObjectLockRule {
	if p == nil {
		return nil
	}
	return &ObjectLockRule{
		DefaultRetention: protoToDefaultRetention(p.DefaultRetention),
	}
}

func defaultRetentionToProto(r *DefaultRetention) *pb.DefaultRetention {
	if r == nil {
		return nil
	}
	return &pb.DefaultRetention{
		Mode:  objectLockRetentionModeToProto(r.Mode),
		Days:  int32PtrToInt32(r.Days),
		Years: int32PtrToInt32(r.Years),
	}
}

func protoToDefaultRetention(p *pb.DefaultRetention) *DefaultRetention {
	if p == nil {
		return nil
	}
	return &DefaultRetention{
		Mode:  protoToObjectLockRetentionMode(p.Mode),
		Days:  int32ToInt32Ptr(p.Days),
		Years: int32ToInt32Ptr(p.Years),
	}
}
