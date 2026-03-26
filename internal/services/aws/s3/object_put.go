package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

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
}

// PutObjectOutput contains the output from the PutObject operation.
type PutObjectOutput struct {
	ETag                 string
	VersionId            string
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// PutObject uploads an object to S3.
func (o *ObjectOperations) PutObject(ctx context.Context, reqCtx *request.RequestContext, input *PutObjectInput) (*PutObjectOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	if input.ContentLength > maxSingleUploadSize {
		return nil, fmt.Errorf("object size exceeds maximum single upload size (5GB), use multipart upload")
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if input.IfNoneMatch != "" || input.IfMatch != "" {
		existingObj, err := store.objects.Head(ctx, input.Bucket, input.Key)
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

	bucket, err := store.buckets.Get(input.Bucket)
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
			return nil, fmt.Errorf("invalid SSE-C customer key: %w", err)
		}
	} else {
		encryptionType = o.svc.encryptionManager.DetermineEncryptionType(
			EncryptionType(input.ServerSideEncryption),
			bucket.EncryptionConfig,
		)
	}

	var obj *s3store.Object
	if o.svc.encryptionManager.ShouldEncrypt(encryptionType, bucket.EncryptionConfig) {
		data, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		encResult, err := o.svc.encryptionManager.EncryptWithCustomerKey(data, encryptionType, bucket.EncryptionConfig, input.Bucket, input.Key, input.SSEKMSKeyId, customerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", err)
		}

		sseMetadata := &s3store.SSEObjectMetadata{
			EncryptionType:   s3store.SSEType(encResult.EncryptionType),
			EncryptedDataKey: encResult.EncryptedDataKey,
			ContentNonce:     encResult.ContentNonce,
			KMSKeyID:         encResult.KMSKeyID,
			UnencryptedMD5:   encResult.UnencryptedMD5,
			UnencryptedSize:  encResult.UnencryptedSize,
		}

		obj, err = store.objects.PutEncrypted(ctx, input.Bucket, input.Key, encResult.EncryptedData, input.ContentType, input.Metadata, sseMetadata)
		if err != nil {
			return nil, err
		}

		return &PutObjectOutput{
			ETag:                 formatETag(obj.ETag),
			VersionId:            obj.VersionID,
			ServerSideEncryption: string(encResult.EncryptionType),
			SSEKMSKeyId:          encResult.KMSKeyID,
		}, nil
	}

	obj, err = store.objects.Put(ctx, input.Bucket, input.Key, input.Body, input.ContentType, input.Metadata)
	if err != nil {
		return nil, err
	}

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
func (o *ObjectOperations) CopyObject(ctx context.Context, reqCtx *request.RequestContext, input *CopyObjectInput) (*CopyObjectOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
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

	if err := o.validateBucketExists(reqCtx, srcBucket); err != nil {
		return nil, ErrInvalidCopySource
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var srcObj *s3store.Object
	if srcVersionId != "" {
		srcObj, err = store.objects.HeadWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		srcObj, err = store.objects.GetMetadata(srcBucket, srcKey)
	}
	if err != nil {
		return nil, ErrInvalidCopySource
	}

	var data []byte
	if srcObj.SSEMetadata != nil || input.CopySourceSSECustomerKey != "" {
		getInput := &GetObjectInput{
			Bucket:               srcBucket,
			Key:                  srcKey,
			VersionId:            srcVersionId,
			SSECustomerAlgorithm: input.CopySourceSSECustomerAlgo,
			SSECustomerKey:       input.CopySourceSSECustomerKey,
			SSECustomerKeyMD5:    input.CopySourceSSECustomerMD5,
		}
		getOutput, err := o.GetObject(ctx, reqCtx, getInput)
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(getOutput.Body)
		getOutput.Body.Close()
		if err != nil {
			return nil, err
		}
	} else {
		var reader io.ReadCloser
		if srcVersionId != "" {
			reader, _, err = store.objects.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
		} else {
			reader, srcObj, err = store.objects.Get(ctx, srcBucket, srcKey)
		}
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, err
		}
	}

	bucketEncryption, _ := store.buckets.GetEncryptionConfiguration(input.Bucket)

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

		encResult, err := o.svc.encryptionManager.EncryptWithCustomerKey(data, targetEncryptionType, bucketEncryption, input.Bucket, input.Key, targetKMSKeyID, customerKey)
		if err != nil {
			return nil, err
		}

		sseMetadata := &s3store.SSEObjectMetadata{
			EncryptionType:   s3store.SSEType(targetEncryptionType),
			EncryptedDataKey: encResult.EncryptedDataKey,
			ContentNonce:     encResult.ContentNonce,
			KMSKeyID:         encResult.KMSKeyID,
			UnencryptedMD5:   encResult.UnencryptedMD5,
			UnencryptedSize:  encResult.UnencryptedSize,
		}

		obj, err = store.objects.PutEncrypted(ctx, input.Bucket, input.Key, encResult.EncryptedData, contentType, metadata, sseMetadata)
		if err != nil {
			return nil, err
		}
	} else {
		if srcVersionId != "" {
			if input.MetadataDirective == "REPLACE" {
				obj, err = store.objects.CopyWithVersionAndMetadata(ctx, srcBucket, srcKey, srcVersionId, input.Bucket, input.Key, contentType, metadata)
			} else {
				obj, err = store.objects.CopyWithVersion(ctx, srcBucket, srcKey, srcVersionId, input.Bucket, input.Key)
			}
		} else {
			if input.MetadataDirective == "REPLACE" {
				obj, err = store.objects.CopyWithMetadata(ctx, srcBucket, srcKey, input.Bucket, input.Key, contentType, metadata)
			} else {
				obj, err = store.objects.Copy(ctx, srcBucket, srcKey, input.Bucket, input.Key)
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
func (o *ObjectOperations) RestoreObject(ctx context.Context, reqCtx *request.RequestContext, input *RestoreObjectInput) (interface{}, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var obj *s3store.Object
	if input.VersionId != "" {
		obj, err = store.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	} else {
		obj, err = store.objects.Head(ctx, input.Bucket, input.Key)
	}
	if err != nil {
		return nil, NewNoSuchKeyError(input.Key)
	}

	if obj.StorageClass == s3store.StorageClassStandard || obj.StorageClass == "" {
		return nil, ErrObjectAlreadyRestored
	}

	if input.Body != nil {
		var restoreReq RestoreRequest
		if err := request.NewSafeXMLDecoder(input.Body).Decode(&restoreReq); err == nil {
			if restoreReq.Days < 1 || restoreReq.Days > 3653 {
				return nil, NewInvalidArgumentError("Days must be between 1 and 3653")
			}
		}
	}

	return nil, nil
}
