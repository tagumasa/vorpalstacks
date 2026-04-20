package s3

import (
	"encoding/xml"
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// dispatchPutBucket handles PUT requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
func (h *S3Handler) dispatchPutBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("acl") {
		action := "s3:PutBucketAcl"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		input := &PutBucketAclInput{Bucket: bucket}
		input.ACL = r.Header.Get("x-amz-acl")
		input.GrantFullControl = r.Header.Get("x-amz-grant-full-control")
		input.GrantRead = r.Header.Get("x-amz-grant-read")
		input.GrantReadACP = r.Header.Get("x-amz-grant-read-acp")
		input.GrantWrite = r.Header.Get("x-amz-grant-write")
		input.GrantWriteACP = r.Header.Get("x-amz-grant-write-acp")
		if input.ACL == "" && input.GrantFullControl == "" && input.GrantRead == "" && input.GrantWrite == "" {
			var acp s3store.AccessControlPolicy
			if err := request.NewSafeXMLDecoder(r.Body).Decode(&acp); err == nil {
				input.AccessControlPolicy = &acp
			}
		}
		err := h.bucketOps.PutBucketAcl(ctx, input)
		return nil, http.StatusOK, err
	}
	if query.Has("versioning") {
		action := "s3:PutBucketVersioning"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config struct {
			Status string `xml:"Status"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketVersioning(ctx, &PutBucketVersioningInput{
			Bucket: bucket,
			Status: config.Status,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("encryption") {
		action := "s3:PutEncryptionConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config ServerSideEncryptionConfiguration
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketEncryption(ctx, &PutBucketEncryptionInput{
			Bucket:                            bucket,
			ServerSideEncryptionConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("policy") {
		action := "s3:PutBucketPolicy"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		err = h.bucketOps.PutBucketPolicy(ctx, &PutBucketPolicyInput{
			Bucket: bucket,
			Policy: string(body),
		})
		return nil, http.StatusNoContent, err
	}
	if query.Has("cors") {
		action := "s3:PutBucketCORS"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config CORSConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketCORS(ctx, &PutBucketCORSInput{
			Bucket:            bucket,
			CORSConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("tagging") {
		action := "s3:PutBucketTagging"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var tagSet struct {
			Tags []Tag `xml:"TagSet>Tag"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&tagSet); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketTagging(ctx, &PutBucketTaggingInput{
			Bucket: bucket,
			Tags:   tagSet.Tags,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("lifecycle") {
		action := "s3:PutLifecycleConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config LifecycleConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketLifecycleConfiguration(ctx, &PutBucketLifecycleConfigurationInput{
			Bucket:                 bucket,
			LifecycleConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("website") {
		action := "s3:PutBucketWebsite"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config WebsiteConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketWebsite(ctx, &PutBucketWebsiteInput{
			Bucket:               bucket,
			WebsiteConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("object-lock") {
		action := "s3:PutBucketObjectLockConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config ObjectLockConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutObjectLockConfiguration(ctx, &PutObjectLockConfigurationInput{
			Bucket:                  bucket,
			ObjectLockConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("notification") {
		action := "s3:PutBucketNotification"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config NotificationConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketNotificationConfiguration(ctx, &PutBucketNotificationInput{
			Bucket:                    bucket,
			NotificationConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("logging") {
		action := "s3:PutBucketLogging"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var wrapper struct {
			LoggingEnabled *LoggingConfigurationInput `xml:"LoggingEnabled"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&wrapper); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketLogging(ctx, &PutBucketLoggingInput{
			Bucket:               bucket,
			LoggingConfiguration: wrapper.LoggingEnabled,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("ownershipControls") {
		action := "s3:PutBucketOwnershipControls"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config OwnershipControlsInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketOwnershipControls(ctx, &PutBucketOwnershipControlsInput{
			Bucket:            bucket,
			OwnershipControls: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("requestPayment") {
		action := "s3:PutBucketRequestPayment"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config RequestPaymentConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketRequestPayment(ctx, &PutBucketRequestPaymentInput{
			Bucket:                      bucket,
			RequestPaymentConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("accelerate") {
		action := "s3:PutAccelerateConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config AccelerateConfigurationInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketAccelerateConfiguration(ctx, &PutBucketAccelerateConfigurationInput{
			Bucket:                  bucket,
			AccelerateConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("publicAccessBlock") {
		action := "s3:PutBucketPublicAccessBlock"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		var config PublicAccessBlockConfiguration
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutPublicAccessBlock(ctx, &PutPublicAccessBlockInput{
			Bucket:                         bucket,
			PublicAccessBlockConfiguration: &config,
		})
		return nil, http.StatusOK, err
	}

	var createConfig struct {
		XMLName            xml.Name `xml:"CreateBucketConfiguration"`
		LocationConstraint string   `xml:"LocationConstraint"`
		ObjectLockEnabled  bool     `xml:"ObjectLockEnabledForBucket"`
		Tags               []Tag    `xml:"Tags>Tag"`
	}
	bodyBytes, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		return nil, http.StatusInternalServerError, bodyErr
	}
	_ = xml.Unmarshal(bodyBytes, &createConfig)

	if r.Header.Get("x-amz-bucket-object-lock-enabled") == "true" {
		createConfig.ObjectLockEnabled = true
	}

	result, err := h.bucketOps.CreateBucket(ctx, &CreateBucketInput{
		Bucket:                     bucket,
		LocationConstraint:         createConfig.LocationConstraint,
		ObjectLockEnabledForBucket: createConfig.ObjectLockEnabled,
	})
	if err != nil {
		return nil, http.StatusOK, err
	}
	if len(createConfig.Tags) > 0 {
		err = h.bucketOps.PutBucketTagging(ctx, &PutBucketTaggingInput{
			Bucket: bucket,
			Tags:   createConfig.Tags,
		})
	}
	return result, http.StatusOK, err
}
