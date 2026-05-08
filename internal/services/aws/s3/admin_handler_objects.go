package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3store "vorpalstacks/internal/store/aws/s3"
)

func (h *AdminHandler) getObjectStore(headers http.Header) s3store.ObjectStoreInterface {
	region := svccommon.GetRegionFromHeader(headers)
	return h.s3Store.Objects(region)
}

func (h *AdminHandler) getStores(headers http.Header) (s3store.BucketStoreInterface, s3store.ObjectStoreInterface) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.s3Store.Buckets(region), h.s3Store.Objects(region)
}

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

// ListObjectsV2 retrieves objects in a bucket for the admin console.
func (h *AdminHandler) ListObjectsV2(ctx context.Context, req *connect.Request[pb.ListObjectsV2Request]) (*connect.Response[pb.ListObjectsV2Output], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	maxKeys := int(req.Msg.Maxkeys)
	if maxKeys <= 0 {
		maxKeys = 1000
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}

	marker := req.Msg.Continuationtoken
	if marker == "" {
		marker = req.Msg.Startafter
	}

	result, err := objectStore.List(req.Msg.Bucket, req.Msg.Prefix, req.Msg.Delimiter, marker, maxKeys)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var contents []*pb.Object
	for _, obj := range result.Objects {
		if obj.IsDeleteMarker {
			continue
		}
		contents = append(contents, &pb.Object{
			Key:          obj.Key,
			Lastmodified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
			Etag:         formatETag(obj.ETag),
			Size:         obj.Size,
			Storageclass: pb.ObjectStorageClass_OBJECT_STORAGE_CLASS_STANDARD,
		})
	}

	var commonPrefixes []*pb.CommonPrefix
	for _, p := range result.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, &pb.CommonPrefix{Prefix: p})
	}

	output := &pb.ListObjectsV2Output{
		Name:              req.Msg.Bucket,
		Prefix:            req.Msg.Prefix,
		Delimiter:         req.Msg.Delimiter,
		Maxkeys:           int32(maxKeys),
		Keycount:          int32(len(contents) + len(commonPrefixes)),
		Istruncated:       result.IsTruncated,
		Contents:          contents,
		Commonprefixes:    commonPrefixes,
		Continuationtoken: req.Msg.Continuationtoken,
		Startafter:        req.Msg.Startafter,
	}
	if result.IsTruncated && result.NextMarker != "" {
		output.Nextcontinuationtoken = result.NextMarker
	}

	return connect.NewResponse(output), nil
}

// HeadObject retrieves metadata for an object without returning the body.
func (h *AdminHandler) HeadObject(ctx context.Context, req *connect.Request[pb.HeadObjectRequest]) (*connect.Response[pb.HeadObjectOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}
	if req.Msg.Key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	obj, err := objectStore.HeadWithVersion(ctx, req.Msg.Bucket, req.Msg.Key, req.Msg.Versionid)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	contentLength := obj.Size
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		contentLength = obj.SSEMetadata.UnencryptedSize
	}

	output := &pb.HeadObjectOutput{
		Contentlength:      contentLength,
		Contenttype:        obj.ContentType,
		Contentencoding:    obj.ContentEncoding,
		Contentlanguage:    obj.ContentLanguage,
		Contentdisposition: obj.ContentDisposition,
		Cachecontrol:       obj.CacheControl,
		Etag:               formatETag(obj.ETag),
		Lastmodified:       obj.LastModified.Format(timeutils.ISO8601UTCFormat),
		Storageclass:       pb.StorageClass_STORAGE_CLASS_STANDARD,
		Versionid:          obj.VersionID,
		Acceptranges:       "bytes",
	}
	if obj.Metadata != nil {
		output.Metadata = obj.Metadata
	}
	if obj.SSEMetadata != nil {
		output.Serversideencryption = storeSSETypeToProto(obj.SSEMetadata.EncryptionType)
		output.Ssekmskeyid = obj.SSEMetadata.KMSKeyID
	}

	return connect.NewResponse(output), nil
}

// GetObject retrieves an object body for the admin console.
// For encrypted objects, returns metadata only without the body.
func (h *AdminHandler) GetObject(ctx context.Context, req *connect.Request[pb.GetObjectRequest]) (*connect.Response[pb.GetObjectOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}
	if req.Msg.Key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	obj, err := objectStore.HeadWithVersion(ctx, req.Msg.Bucket, req.Msg.Key, req.Msg.Versionid)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	contentLength := obj.Size
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		contentLength = obj.SSEMetadata.UnencryptedSize
	}

	output := &pb.GetObjectOutput{
		Contentlength:      contentLength,
		Contenttype:        obj.ContentType,
		Contentencoding:    obj.ContentEncoding,
		Contentlanguage:    obj.ContentLanguage,
		Contentdisposition: obj.ContentDisposition,
		Cachecontrol:       obj.CacheControl,
		Etag:               formatETag(obj.ETag),
		Lastmodified:       obj.LastModified.Format(timeutils.ISO8601UTCFormat),
		Versionid:          obj.VersionID,
		Acceptranges:       "bytes",
	}
	if obj.Metadata != nil {
		output.Metadata = obj.Metadata
	}
	if obj.SSEMetadata != nil {
		output.Serversideencryption = storeSSETypeToProto(obj.SSEMetadata.EncryptionType)
		output.Ssekmskeyid = obj.SSEMetadata.KMSKeyID
	}

	var reader io.ReadCloser
	if req.Msg.Versionid != "" {
		reader, _, err = objectStore.GetWithVersion(ctx, req.Msg.Bucket, req.Msg.Key, req.Msg.Versionid)
	} else {
		reader, _, err = objectStore.Get(ctx, req.Msg.Bucket, req.Msg.Key)
	}
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	output.Body = data

	return connect.NewResponse(output), nil
}
