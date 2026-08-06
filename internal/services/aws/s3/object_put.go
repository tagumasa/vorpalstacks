package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/aws/types"
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

	if input.ContentLength > maxSingleUploadSize {
		return nil, ErrEntityTooLarge
	}

	if input.IfNoneMatch != "" || input.IfMatch != "" {
		existingObj, err := stores.objects.Head(ctx, input.Bucket, input.Key)
		objectExists := err == nil && existingObj != nil

		if input.IfNoneMatch == "*" {
			if objectExists {
				return nil, ErrPreconditionFailed
			}
		} else if input.IfNoneMatch != "" {
			if objectExists && strings.Trim(existingObj.ETag, "\"") == strings.Trim(input.IfNoneMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}

		if input.IfMatch != "" {
			if !objectExists {
				return nil, ErrPreconditionFailed
			}
			if strings.Trim(existingObj.ETag, "\"") != strings.Trim(input.IfMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}
	}

	bucket, err := stores.buckets.Get(input.Bucket)
	if err != nil {
		return nil, err
	}

	var encryptionType EncryptionType
	var customerKey []byte
	if input.SSECustomerAlgorithm != "" {
		encryptionType = EncryptionTypeSSE_C
		var err error
		customerKey, err = o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
		if err != nil {
			return nil, NewInvalidArgumentError(fmt.Sprintf("invalid SSE-C customer key: %v", err))
		}
	} else {
		encryptionType = o.svc.encryptionManager.DetermineEncryptionType(
			EncryptionType(input.ServerSideEncryption),
			bucket.EncryptionConfig,
		)
	}

	var obj *s3store.Object
	storageClass := s3store.ObjectStorageClass(input.StorageClass)
	if storageClass == "" {
		storageClass = s3store.StorageClassStandard
	}

	sysMeta := &s3store.SystemMetadata{
		ContentEncoding:    input.ContentEncoding,
		ContentLanguage:    input.ContentLanguage,
		ContentDisposition: input.ContentDisposition,
		CacheControl:       input.CacheControl,
	}

	if o.svc.encryptionManager.ShouldEncrypt(encryptionType, bucket.EncryptionConfig) {
		encResult, err := o.svc.encryptionManager.EncryptStream(input.Body, encryptionType, bucket.EncryptionConfig, input.Bucket, input.Key, input.SSEKMSKeyId, customerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", err)
		}

		obj, err = stores.objects.PutEncrypted(ctx, input.Bucket, input.Key, encResult.EncryptedData, input.ContentType, input.Metadata, encResult.SSEMetadata, storageClass, sysMeta)
		if err != nil {
			return nil, err
		}

		if input.Tagging != "" {
			parsedTags := parseTaggingHeader(input.Tagging)
			if len(parsedTags) > 0 {
				obj.Tags = parsedTags
				_ = stores.objects.SetTags(input.Bucket, input.Key, parsedTags)
			}
		}

		o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedPut)
		repCtx, repCancel := context.WithTimeout(context.Background(), 30*time.Second)
		go func() {
			defer repCancel()
			defer func() {
				if r := recover(); r != nil {
					logs.Error("s3: replication goroutine panic",
						logs.String("bucket", bucket.Name),
						logs.String("key", input.Key),
						logs.Any("panic", r))
				}
			}()
			o.svc.replicateObject(repCtx, reqCtx, stores, bucket, input.Key, obj)
		}()

		return &PutObjectOutput{
			ETag:                 formatETag(obj.ETag),
			VersionId:            obj.VersionID,
			ServerSideEncryption: string(encResult.SSEMetadata.EncryptionType),
			SSEKMSKeyId:          encResult.SSEMetadata.KMSKeyID,
		}, nil
	}

	obj, err = stores.objects.PutWithVersioning(ctx, input.Bucket, input.Key, input.Body, input.ContentType, input.Metadata, false, storageClass, sysMeta)
	if err != nil {
		return nil, err
	}

	if input.Tagging != "" {
		parsedTags := parseTaggingHeader(input.Tagging)
		if len(parsedTags) > 0 {
			obj.Tags = parsedTags
			_ = stores.objects.SetTags(input.Bucket, input.Key, parsedTags)
		}
	}

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedPut)
	repCtx2, repCancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer repCancel2()
		defer func() {
			if r := recover(); r != nil {
				logs.Error("s3: replication goroutine panic",
					logs.String("bucket", bucket.Name),
					logs.String("key", input.Key),
					logs.Any("panic", r))
			}
		}()
		o.svc.replicateObject(repCtx2, reqCtx, stores, bucket, input.Key, obj)
	}()

	return &PutObjectOutput{
		ETag:      formatETag(obj.ETag),
		VersionId: obj.VersionID,
	}, nil
}

// CopyObjectInput contains the input parameters for the CopyObject operation.
type CopyObjectInput struct {
	Bucket                    string
	Key                       string
	CopySource                string
	CopySourceVersionId       string
	MetadataDirective         string
	ContentType               string
	Metadata                  map[string]string
	ServerSideEncryption      string
	SSEKMSKeyId               string
	SSECustomerAlgorithm      string
	SSECustomerKey            string
	SSECustomerKeyMD5         string
	CopySourceSSECustomerAlgo string
	CopySourceSSECustomerKey  string
	CopySourceSSECustomerMD5  string
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

// CopyObject copies an object to another location in S3.
func (o *ObjectOperations) CopyObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *CopyObjectInput) (*CopyObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	srcBucket, srcKey, srcVersionId, err := parseCopySource(input.CopySource)
	if err != nil {
		return nil, err
	}

	if input.CopySourceVersionId != "" {
		srcVersionId = input.CopySourceVersionId
	}

	if err := o.validateBucketExists(stores, srcBucket); err != nil {
		return nil, ErrInvalidCopySource
	}

	var srcObj *s3store.Object
	if srcVersionId != "" {
		srcObj, err = stores.objects.HeadWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		srcObj, err = stores.objects.GetMetadata(srcBucket, srcKey)
	}
	if err != nil {
		return nil, ErrInvalidCopySource
	}

	if srcObj.Size > maxCopyObjectSize {
		return nil, ErrEntityTooLarge
	}

	var srcReader io.Reader
	if srcObj.SSEMetadata != nil || input.CopySourceSSECustomerKey != "" {
		getInput := &GetObjectInput{
			Bucket:               srcBucket,
			Key:                  srcKey,
			VersionId:            srcVersionId,
			SSECustomerAlgorithm: input.CopySourceSSECustomerAlgo,
			SSECustomerKey:       input.CopySourceSSECustomerKey,
			SSECustomerKeyMD5:    input.CopySourceSSECustomerMD5,
		}
		getOutput, err := o.GetObject(ctx, reqCtx, stores, getInput)
		if err != nil {
			return nil, err
		}
		defer getOutput.Body.Close()
		srcReader = getOutput.Body
	} else {
		var reader io.ReadCloser
		if srcVersionId != "" {
			reader, _, err = stores.objects.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
		} else {
			reader, srcObj, err = stores.objects.Get(ctx, srcBucket, srcKey)
		}
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		srcReader = reader
	}

	bucketEncryption, err := stores.buckets.GetEncryptionConfiguration(input.Bucket)
	if err != nil {
		return nil, err
	}

	var targetEncryptionType EncryptionType
	var targetKMSKeyID string

	if input.ServerSideEncryption != "" {
		targetEncryptionType = EncryptionType(input.ServerSideEncryption)
		targetKMSKeyID = input.SSEKMSKeyId
	} else if input.SSECustomerAlgorithm != "" {
		targetEncryptionType = EncryptionTypeSSE_C
	} else {
		targetEncryptionType = o.svc.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
		if targetEncryptionType == EncryptionTypeSSE_KMS && bucketEncryption != nil {
			targetKMSKeyID = bucketEncryption.KMSMasterKeyID
		}
	}

	var obj *s3store.Object
	contentType := input.ContentType
	if contentType == "" {
		contentType = srcObj.ContentType
	}
	metadata := input.Metadata
	if input.MetadataDirective != "" && input.MetadataDirective != "COPY" && input.MetadataDirective != "REPLACE" {
		return nil, NewInvalidArgumentError(fmt.Sprintf("invalid MetadataDirective: %s (must be COPY or REPLACE)", input.MetadataDirective))
	}
	if input.MetadataDirective != "REPLACE" {
		metadata = srcObj.Metadata
	}

	if targetEncryptionType != EncryptionTypeNone {
		var customerKey []byte
		if input.SSECustomerKey != "" {
			var err error
			customerKey, err = o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, err
			}
		}

		encResult, err := o.svc.encryptionManager.EncryptStream(srcReader, targetEncryptionType, bucketEncryption, input.Bucket, input.Key, targetKMSKeyID, customerKey)
		if err != nil {
			return nil, err
		}

		targetStorageClass := srcObj.StorageClass
		if targetStorageClass == "" {
			targetStorageClass = s3store.StorageClassStandard
		}
		obj, err = stores.objects.PutEncrypted(ctx, input.Bucket, input.Key, encResult.EncryptedData, contentType, metadata, encResult.SSEMetadata, targetStorageClass, nil)
		if err != nil {
			return nil, err
		}
	} else {
		if srcVersionId != "" {
			if input.MetadataDirective == "REPLACE" {
				obj, err = stores.objects.CopyWithVersionAndMetadata(ctx, srcBucket, srcKey, srcVersionId, input.Bucket, input.Key, contentType, metadata)
			} else {
				obj, err = stores.objects.CopyWithVersion(ctx, srcBucket, srcKey, srcVersionId, input.Bucket, input.Key)
			}
		} else {
			if input.MetadataDirective == "REPLACE" {
				obj, err = stores.objects.CopyWithMetadata(ctx, srcBucket, srcKey, input.Bucket, input.Key, contentType, metadata)
			} else {
				obj, err = stores.objects.Copy(ctx, srcBucket, srcKey, input.Bucket, input.Key)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	output := &CopyObjectOutput{
		CopyObjectResult: &CopyObjectResult{
			ETag:         formatETag(obj.ETag),
			LastModified: obj.LastModified,
		},
	}

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedCopy)

	if obj.SSEMetadata != nil {
		output.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
		if obj.SSEMetadata.KMSKeyID != "" {
			output.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}
	}

	return output, nil
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

// RestoreObject restores an archived copy of an object back into S3.
func (o *ObjectOperations) RestoreObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *RestoreObjectInput) (interface{}, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	var obj *s3store.Object
	var err error
	if input.VersionId != "" {
		obj, err = stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	} else {
		obj, err = stores.objects.Head(ctx, input.Bucket, input.Key)
	}
	if err != nil {
		return nil, NewNoSuchKeyError(input.Key)
	}

	if obj.StorageClass == s3store.StorageClassStandard || obj.StorageClass == "" {
		return nil, ErrObjectAlreadyRestored
	}

	restoreDays := 1
	if input.Body != nil {
		var restoreReq RestoreRequest
		if err := request.NewSafeXMLDecoder(input.Body).Decode(&restoreReq); err != nil {
			return nil, NewInvalidArgumentError("invalid RestoreObject request body")
		}
		if err := validateRestoreDays(restoreReq.Days); err != nil {
			return nil, err
		}
		restoreDays = restoreReq.Days
	}

	if err := stores.objects.SetStorageClass(input.Bucket, input.Key, input.VersionId, s3store.StorageClassStandard); err != nil {
		return nil, err
	}

	logs.Info("s3: object restored",
		logs.String("bucket", input.Bucket),
		logs.String("key", input.Key),
		logs.String("days", fmt.Sprintf("%d", restoreDays)))

	return nil, nil
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
