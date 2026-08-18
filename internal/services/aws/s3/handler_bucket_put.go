package s3

import (
	"encoding/xml"
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// aclHeaders carries the x-amz-acl and x-amz-grant-* request header values
// shared by PutBucketAcl and CreateBucket.
type aclHeaders struct {
	ACL              string
	GrantFullControl string
	GrantRead        string
	GrantReadACP     string
	GrantWrite       string
	GrantWriteACP    string
}

func parseACLHeaders(r *http.Request) aclHeaders {
	return aclHeaders{
		ACL:              r.Header.Get("x-amz-acl"),
		GrantFullControl: r.Header.Get("x-amz-grant-full-control"),
		GrantRead:        r.Header.Get("x-amz-grant-read"),
		GrantReadACP:     r.Header.Get("x-amz-grant-read-acp"),
		GrantWrite:       r.Header.Get("x-amz-grant-write"),
		GrantWriteACP:    r.Header.Get("x-amz-grant-write-acp"),
	}
}

// dispatchPutBucket handles PUT requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
// IAM action checks run in handleBucketRequest before dispatch.
func (h *S3Handler) dispatchPutBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("acl") {
		headers := parseACLHeaders(r)
		input := &PutBucketAclInput{
			Bucket:           bucket,
			ACL:              headers.ACL,
			GrantFullControl: headers.GrantFullControl,
			GrantRead:        headers.GrantRead,
			GrantReadACP:     headers.GrantReadACP,
			GrantWrite:       headers.GrantWrite,
			GrantWriteACP:    headers.GrantWriteACP,
		}
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
		var config struct {
			Status    string `xml:"Status"`
			MfaDelete string `xml:"MfaDelete"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketVersioning(ctx, &PutBucketVersioningInput{
			Bucket:    bucket,
			Status:    config.Status,
			MFADelete: config.MfaDelete,
		})
		return nil, http.StatusOK, err
	}
	if query.Has("encryption") {
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
	if query.Has("replication") {
		var config ReplicationConfigurationXML
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
			return nil, http.StatusBadRequest, err
		}
		err := h.bucketOps.PutBucketReplication(ctx, &PutBucketReplicationInput{
			Bucket:                   bucket,
			ReplicationConfiguration: &config,
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
	if len(bodyBytes) > 0 {
		if err := xml.Unmarshal(bodyBytes, &createConfig); err != nil {
			return nil, 0, ErrMalformedXML
		}
	}

	if r.Header.Get("x-amz-bucket-object-lock-enabled") == "true" {
		createConfig.ObjectLockEnabled = true
	}

	// CreateBucket accepts the same ACL headers as PutBucketAcl plus the
	// x-amz-object-ownership header that seeds the ownership controls.
	aclHdrs := parseACLHeaders(r)
	result, err := h.bucketOps.CreateBucket(ctx, &CreateBucketInput{
		Bucket:                     bucket,
		ACL:                        aclHdrs.ACL,
		GrantFullControl:           aclHdrs.GrantFullControl,
		GrantRead:                  aclHdrs.GrantRead,
		GrantReadACP:               aclHdrs.GrantReadACP,
		GrantWrite:                 aclHdrs.GrantWrite,
		GrantWriteACP:              aclHdrs.GrantWriteACP,
		ObjectOwnership:            r.Header.Get("x-amz-object-ownership"),
		LocationConstraint:         createConfig.LocationConstraint,
		ObjectLockEnabledForBucket: createConfig.ObjectLockEnabled,
	})
	if err != nil {
		return nil, http.StatusConflict, err
	}
	if len(createConfig.Tags) > 0 {
		err = h.bucketOps.PutBucketTagging(ctx, &PutBucketTaggingInput{
			Bucket: bucket,
			Tags:   createConfig.Tags,
		})
	}
	return result, http.StatusOK, err
}
