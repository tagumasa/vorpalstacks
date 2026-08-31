package s3

import (
	"context"
	"io"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
)

const maxCopyObjectSize int64 = 5 * 1024 * 1024 * 1024

// PutObjectInput contains the input parameters for the PutObject operation.
type PutObjectInput struct {
	Bucket               string
	Key                  string
	Body                 io.Reader
	ContentLength        int64
	ContentType          string
	ContentEncoding      string
	ContentLanguage      string
	ContentDisposition   string
	CacheControl         string
	Metadata             map[string]string
	StorageClass         string
	IfMatch              string
	IfNoneMatch          string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
	Tagging              string
	ACLHeaders           aclHeaders
}

// PutObjectOutput contains the output from the PutObject operation.
type PutObjectOutput struct {
	ETag                 string
	VersionId            string
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// PutObject uploads an object to S3.
func (o *ObjectOperations) PutObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *PutObjectInput) (*PutObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}
	if err := validateStorageClass(input.StorageClass); err != nil {
		return nil, err
	}
	acl, err := o.svc.resolveUploadACL(ctx, stores, input.Bucket, input.ACLHeaders)
	if err != nil {
		return nil, err
	}

	coreResult, err := o.svc.putObjectStreamCore(ctx, stores.buckets, stores.objects, PutObjectStreamInput{
		Body:                 input.Body,
		ContentLength:        input.ContentLength,
		Bucket:               input.Bucket,
		Key:                  input.Key,
		ContentType:          input.ContentType,
		ContentEncoding:      input.ContentEncoding,
		ContentLanguage:      input.ContentLanguage,
		ContentDisposition:   input.ContentDisposition,
		CacheControl:         input.CacheControl,
		Metadata:             input.Metadata,
		StorageClass:         input.StorageClass,
		IfMatch:              input.IfMatch,
		IfNoneMatch:          input.IfNoneMatch,
		ServerSideEncryption: input.ServerSideEncryption,
		SSEKMSKeyId:          input.SSEKMSKeyId,
		SSECustomerAlgorithm: input.SSECustomerAlgorithm,
		SSECustomerKey:       input.SSECustomerKey,
		SSECustomerKeyMD5:    input.SSECustomerKeyMD5,
		Tagging:              input.Tagging,
		ACL:                  acl,
	})
	if err != nil {
		return nil, err
	}
	obj := coreResult.Object

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedPut)
	repCtx, repCancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer repCancel()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("s3: replication goroutine panic",
					logs.String("bucket", coreResult.Bucket.Name),
					logs.String("key", input.Key),
					logs.Any("panic", r))
			}
		}()
		o.svc.replicateObject(repCtx, reqCtx, stores, coreResult.Bucket, input.Key, obj)
	}()

	return &PutObjectOutput{
		ETag:                 formatETag(obj.ETag),
		VersionId:            obj.VersionID,
		ServerSideEncryption: coreResult.ServerSideEncryption,
		SSEKMSKeyId:          coreResult.SSEKMSKeyId,
	}, nil
}

// CopyObjectInput contains the input parameters for the CopyObject operation.
type CopyObjectInput struct {
	Bucket                      string
	Key                         string
	CopySource                  string
	CopySourceVersionId         string
	CopySourceIfMatch           string
	CopySourceIfNoneMatch       string
	CopySourceIfModifiedSince   *time.Time
	CopySourceIfUnmodifiedSince *time.Time
	MetadataDirective           string
	ContentType                 string
	Metadata                    map[string]string
	StorageClass                string
	ServerSideEncryption        string
	SSEKMSKeyId                 string
	SSECustomerAlgorithm        string
	SSECustomerKey              string
	SSECustomerKeyMD5           string
	CopySourceSSECustomerAlgo   string
	CopySourceSSECustomerKey    string
	CopySourceSSECustomerMD5    string
	ACLHeaders                  aclHeaders
}

// CopyObjectOutput contains the output from the CopyObject operation.
type CopyObjectOutput struct {
	CopyObjectResult     *CopyObjectResult `xml:"CopyObjectResult"`
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// CopyObjectResult contains the result information from a CopyObject operation.
type CopyObjectResult struct {
	ETag         string    `xml:"ETag"`
	LastModified time.Time `xml:"LastModified"`
}

func (o *ObjectOperations) CopyObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *CopyObjectInput) (*CopyObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	srcBucket, _, _, err := parseCopySource(input.CopySource)
	if err != nil {
		return nil, err
	}

	if err := o.validateBucketExists(stores, srcBucket); err != nil {
		return nil, ErrInvalidCopySource
	}

	acl, err := o.svc.resolveUploadACL(ctx, stores, input.Bucket, input.ACLHeaders)
	if err != nil {
		return nil, err
	}

	coreResult, err := o.svc.copyObjectStreamCore(ctx, stores.buckets, stores.objects, CopyObjectStreamInput{
		Bucket:                      input.Bucket,
		Key:                         input.Key,
		CopySource:                  input.CopySource,
		CopySourceVersionId:         input.CopySourceVersionId,
		CopySourceIfMatch:           input.CopySourceIfMatch,
		CopySourceIfNoneMatch:       input.CopySourceIfNoneMatch,
		CopySourceIfModifiedSince:   input.CopySourceIfModifiedSince,
		CopySourceIfUnmodifiedSince: input.CopySourceIfUnmodifiedSince,
		MetadataDirective:           input.MetadataDirective,
		ContentType:                 input.ContentType,
		Metadata:                    input.Metadata,
		StorageClass:                input.StorageClass,
		ServerSideEncryption:        input.ServerSideEncryption,
		SSEKMSKeyId:                 input.SSEKMSKeyId,
		SSECustomerAlgorithm:        input.SSECustomerAlgorithm,
		SSECustomerKey:              input.SSECustomerKey,
		SSECustomerKeyMD5:           input.SSECustomerKeyMD5,
		CopySourceSSECustomerAlgo:   input.CopySourceSSECustomerAlgo,
		CopySourceSSECustomerKey:    input.CopySourceSSECustomerKey,
		CopySourceSSECustomerMD5:    input.CopySourceSSECustomerMD5,
		ACL:                         acl,
	})
	if err != nil {
		return nil, err
	}
	obj := coreResult.Object

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedCopy)

	return &CopyObjectOutput{
		CopyObjectResult: &CopyObjectResult{
			ETag:         formatETag(obj.ETag),
			LastModified: obj.LastModified,
		},
		ServerSideEncryption: coreResult.ServerSideEncryption,
		SSEKMSKeyId:          coreResult.SSEKMSKeyId,
	}, nil
}

// RestoreObjectInput contains the parameters for restoring an archived object.
type RestoreObjectInput struct {
	Bucket    string
	Key       string
	VersionId string
	Body      io.Reader
}

// RestoreRequest specifies the parameters for a restore request, such as the number of days.
type RestoreRequest struct {
	Days int `xml:"Days"`
}

// RestoreObject creates or extends the temporary restored copy of an
// archived object and reports whether a restored copy already existed (the
// request then only extended its expiry, which the API answers with 200 OK
// instead of 202 Accepted). The object's storage class never changes; the
// restored copy's expiry is rounded up to the following midnight UTC.
func (o *ObjectOperations) RestoreObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *RestoreObjectInput) (bool, error) {
	return o.svc.restoreObjectCore(ctx, reqCtx, stores, input)
}

// parseTaggingHeader parses the x-amz-tagging header value (URL-encoded
// key=value pairs separated by &) into a slice of Tag structs.
func parseTaggingHeader(tagging string) []types.Tag {
	if tagging == "" {
		return nil
	}
	var tags []types.Tag
	for _, pair := range strings.Split(tagging, "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && kv[0] != "" {
			tags = append(tags, types.Tag{Key: kv[0], Value: kv[1]})
		}
	}
	return tags
}
