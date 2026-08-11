package s3

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/s3"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// ---------------------------------------------------------------------------
// Store accessors — extract region from headers and return regional stores.
// This is the only admin_handler file that imports the s3store package.
// ---------------------------------------------------------------------------

// getBucketStoreFromHeaders returns the bucket store for the region in the
// request headers.
func (h *AdminHandler) getBucketStoreFromHeaders(headers http.Header) s3store.BucketStoreInterface {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.s3Store.Buckets(region)
}

// getObjectStoreFromHeaders returns the object store for the region in the
// request headers.
func (h *AdminHandler) getObjectStoreFromHeaders(headers http.Header) s3store.ObjectStoreInterface {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.s3Store.Objects(region)
}

// getStoresFromHeaders returns both bucket and object stores for the region.
func (h *AdminHandler) getStoresFromHeaders(headers http.Header) (s3store.BucketStoreInterface, s3store.ObjectStoreInterface) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.s3Store.Buckets(region), h.service.s3Store.Objects(region)
}

// ---------------------------------------------------------------------------
// SSE helpers
// ---------------------------------------------------------------------------

// storeSSETypeToProto converts a store SSE type to the proto enum value.
func storeSSETypeToProto(sseType s3store.SSEType) pb.ServerSideEncryption {
	switch sseType {
	case s3store.SSETypeKMS:
		return pb.ServerSideEncryption_SERVER_SIDE_ENCRYPTION_AWS_KMS
	case s3store.SSETypeDSSEKMS:
		return pb.ServerSideEncryption_SERVER_SIDE_ENCRYPTION_AWS_KMS_DSSE
	case s3store.SSETypeAES256:
		return pb.ServerSideEncryption_SERVER_SIDE_ENCRYPTION_AES256
	default:
		return pb.ServerSideEncryption_SERVER_SIDE_ENCRYPTION_AES256
	}
}

// sseContentLength returns the effective content length, preferring the
// unencrypted size when SSE metadata is present.
func sseContentLength(obj *s3store.Object) int64 {
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		return obj.SSEMetadata.UnencryptedSize
	}
	return obj.Size
}

// fillObjectMetadata populates proto output fields common to HeadObject and
// GetObject from a store Object.
func fillObjectMetadata(out *pbMetaFields, obj *s3store.Object) {
	out.contentLength = sseContentLength(obj)
	out.contentType = obj.ContentType
	out.contentEncoding = obj.ContentEncoding
	out.contentLanguage = obj.ContentLanguage
	out.contentDisposition = obj.ContentDisposition
	out.cacheControl = obj.CacheControl
	out.etag = formatETag(obj.ETag)
	out.lastModified = obj.LastModified.Format(timeutils.ISO8601UTCFormat)
	out.versionID = obj.VersionID
	out.metadata = obj.Metadata
	if obj.SSEMetadata != nil {
		out.sseType = storeSSETypeToProto(obj.SSEMetadata.EncryptionType)
		out.kmsKeyID = obj.SSEMetadata.KMSKeyID
	}
}

// pbMetaFields is a helper struct used by fillObjectMetadata to populate
// proto output messages without duplicating the field-mapping logic.
type pbMetaFields struct {
	contentLength      int64
	contentType        string
	contentEncoding    string
	contentLanguage    string
	contentDisposition string
	cacheControl       string
	etag               string
	lastModified       string
	versionID          string
	metadata           map[string]string
	sseType            pb.ServerSideEncryption
	kmsKeyID           string
}

// ---------------------------------------------------------------------------
// Proto → DTO input conversion functions
// ---------------------------------------------------------------------------

func pbToCreateBucketInput(msg *pb.CreateBucketRequest, region string) AdminCreateBucketInput {
	return AdminCreateBucketInput{
		Bucket: msg.Bucket,
		Region: region,
	}
}

func pbToDeleteBucketInput(msg *pb.DeleteBucketRequest) AdminDeleteBucketInput {
	return AdminDeleteBucketInput{
		Bucket: msg.Bucket,
	}
}

func pbToListObjectsInput(msg *pb.ListObjectsV2Request) AdminListObjectsInput {
	marker := msg.Continuationtoken
	if marker == "" {
		marker = msg.Startafter
	}
	return AdminListObjectsInput{
		Bucket:    msg.Bucket,
		Prefix:    msg.Prefix,
		Delimiter: msg.Delimiter,
		Marker:    marker,
		MaxKeys:   int(msg.GetMaxkeys()),
	}
}

func pbToHeadObjectInput(msg *pb.HeadObjectRequest) AdminHeadObjectInput {
	return AdminHeadObjectInput{
		Bucket:    msg.Bucket,
		Key:       msg.Key,
		VersionID: msg.Versionid,
	}
}

func pbToGetObjectInput(msg *pb.GetObjectRequest) AdminGetObjectInput {
	return AdminGetObjectInput{
		Bucket:    msg.Bucket,
		Key:       msg.Key,
		VersionID: msg.Versionid,
	}
}

func pbToPutObjectInput(msg *pb.PutObjectRequest) AdminPutObjectInput {
	return AdminPutObjectInput{
		Bucket:      msg.Bucket,
		Key:         msg.Key,
		Body:        msg.Body,
		ContentType: msg.Contenttype,
		Metadata:    msg.Metadata,
	}
}

func pbToDeleteObjectInput(msg *pb.DeleteObjectRequest) AdminDeleteObjectInput {
	return AdminDeleteObjectInput{
		Bucket:    msg.Bucket,
		Key:       msg.Key,
		VersionID: msg.Versionid,
	}
}

func pbToDeleteObjectsInput(msg *pb.DeleteObjectsRequest) AdminDeleteObjectsInput {
	objects := make([]AdminObjectIdentifier, 0, len(msg.Delete.Objects))
	for _, o := range msg.Delete.Objects {
		objects = append(objects, AdminObjectIdentifier{
			Key:       o.Key,
			VersionID: o.Versionid,
		})
	}
	return AdminDeleteObjectsInput{
		Bucket:  msg.Bucket,
		Objects: objects,
	}
}

func pbToCopyObjectInput(msg *pb.CopyObjectRequest) AdminCopyObjectInput {
	return AdminCopyObjectInput{
		Bucket:      msg.Bucket,
		Key:         msg.Key,
		CopySource:  msg.Copysource,
		ContentType: msg.Contenttype,
	}
}

// ---------------------------------------------------------------------------
// DTO → Proto result conversion functions
// ---------------------------------------------------------------------------

func listBucketsResultToPb(result *AdminListBucketsResult) *pb.ListBucketsOutput {
	var buckets []*pb.Bucket
	for _, b := range result.Buckets {
		buckets = append(buckets, &pb.Bucket{
			Name:         b.Name,
			Creationdate: b.CreationDate.Format(timeutils.ISO8601UTCFormat),
			Bucketregion: b.Region,
		})
	}
	return &pb.ListBucketsOutput{Buckets: buckets}
}

func createBucketResultToPb(result *AdminCreateBucketResult) *pb.CreateBucketOutput {
	return &pb.CreateBucketOutput{Location: result.Location}
}

func listObjectsResultToPb(result *AdminListObjectsResult, in AdminListObjectsInput, maxKeys int) *pb.ListObjectsV2Output {
	var contents []*pb.Object
	for _, obj := range result.Objects {
		if obj.IsDeleteMarker {
			continue
		}
		contents = append(contents, &pb.Object{
			Key:          obj.Key,
			Lastmodified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
			Etag:         formatETag(obj.ETag),
			Size:         proto.Int64(obj.Size),
			Storageclass: pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD,
		})
	}

	var commonPrefixes []*pb.CommonPrefix
	for _, p := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, &pb.CommonPrefix{Prefix: p})
	}

	output := &pb.ListObjectsV2Output{
		Name:              in.Bucket,
		Prefix:            in.Prefix,
		Delimiter:         in.Delimiter,
		Maxkeys:           proto.Int32(int32(maxKeys)),
		Keycount:          proto.Int32(int32(len(contents) + len(commonPrefixes))),
		Istruncated:       proto.Bool(result.IsTruncated),
		Contents:          contents,
		Commonprefixes:    commonPrefixes,
		Continuationtoken: in.Marker,
		Startafter:        in.Marker,
	}
	if result.IsTruncated && result.NextMarker != "" {
		output.Nextcontinuationtoken = result.NextMarker
	}
	return output
}

func headObjectResultToPb(result *AdminHeadObjectResult) *pb.HeadObjectOutput {
	var mf pbMetaFields
	fillObjectMetadata(&mf, result.Object)

	out := &pb.HeadObjectOutput{
		Contentlength:      proto.Int64(mf.contentLength),
		Contenttype:        mf.contentType,
		Contentencoding:    mf.contentEncoding,
		Contentlanguage:    mf.contentLanguage,
		Contentdisposition: mf.contentDisposition,
		Cachecontrol:       mf.cacheControl,
		Etag:               mf.etag,
		Lastmodified:       mf.lastModified,
		Storageclass:       pb.StorageClass_STORAGE_CLASS_STANDARD,
		Versionid:          mf.versionID,
		Acceptranges:       "bytes",
	}
	if mf.metadata != nil {
		out.Metadata = mf.metadata
	}
	out.Serversideencryption = mf.sseType
	out.Ssekmskeyid = mf.kmsKeyID
	return out
}

func getObjectResultToPb(result *AdminGetObjectResult) *pb.GetObjectOutput {
	var mf pbMetaFields
	fillObjectMetadata(&mf, result.Object)

	out := &pb.GetObjectOutput{
		Contentlength:      proto.Int64(mf.contentLength),
		Contenttype:        mf.contentType,
		Contentencoding:    mf.contentEncoding,
		Contentlanguage:    mf.contentLanguage,
		Contentdisposition: mf.contentDisposition,
		Cachecontrol:       mf.cacheControl,
		Etag:               mf.etag,
		Lastmodified:       mf.lastModified,
		Versionid:          mf.versionID,
		Acceptranges:       "bytes",
		Body:               result.Body,
	}
	if mf.metadata != nil {
		out.Metadata = mf.metadata
	}
	out.Serversideencryption = mf.sseType
	out.Ssekmskeyid = mf.kmsKeyID
	return out
}

func putObjectResultToPb(result *AdminPutObjectResult) *pb.PutObjectOutput {
	out := &pb.PutObjectOutput{
		Etag:      result.ETag,
		Versionid: result.VersionID,
		Size:      proto.Int64(result.Size),
	}
	out.Ssekmskeyid = result.KMSKeyID
	return out
}

func deleteObjectResultToPb(result *AdminDeleteObjectResult) *pb.DeleteObjectOutput {
	return &pb.DeleteObjectOutput{
		Versionid:    result.VersionID,
		Deletemarker: proto.Bool(result.IsDeleteMarker),
	}
}

func deleteObjectsResultToPb(result *AdminDeleteObjectsResult) *pb.DeleteObjectsOutput {
	var deleted []*pb.DeletedObject
	for _, d := range result.Deleted {
		deleted = append(deleted, &pb.DeletedObject{
			Key:                   d.Key,
			Versionid:             d.VersionID,
			Deletemarker:          proto.Bool(d.DeleteMarker),
			Deletemarkerversionid: d.DeleteMarkerVersionID,
		})
	}

	var errors []*pb.Error
	for _, e := range result.Errors {
		errors = append(errors, &pb.Error{
			Key:     e.Key,
			Code:    e.Code,
			Message: e.Message,
		})
	}

	return &pb.DeleteObjectsOutput{
		Deleted: deleted,
		Errors:  errors,
	}
}

func copyObjectResultToPb(result *AdminCopyObjectResult) *pb.CopyObjectOutput {
	return &pb.CopyObjectOutput{
		Copyobjectresult: &pb.CopyObjectResult{
			Etag:         result.ETag,
			Lastmodified: result.LastModified,
		},
	}
}

// requireBucket validates that a bucket name is non-empty.
func requireBucket(bucket string) error {
	if bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	return nil
}

// requireKey validates that an object key is non-empty.
func requireKey(key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}
