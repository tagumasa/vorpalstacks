package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

func (h *S3Handler) handleServiceRequest(ctx *request.RequestContext, r *http.Request) (interface{}, int, error) {
	if r.Method == "GET" {
		result, err := h.bucketOps.ListBuckets(ctx, &ListBucketsInput{})
		return result, http.StatusOK, err
	}
	return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed")
}

func (h *S3Handler) handleBucketRequest(ctx *request.RequestContext, r *http.Request, bucket string) (interface{}, int, error) {
	query := r.URL.Query()

	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	switch r.Method {
	case "PUT":
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
			var config LoggingConfigurationInput
			if err := request.NewSafeXMLDecoder(r.Body).Decode(&config); err != nil {
				return nil, http.StatusBadRequest, err
			}
			err := h.bucketOps.PutBucketLogging(ctx, &PutBucketLoggingInput{
				Bucket:               bucket,
				LoggingConfiguration: &config,
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

	case "GET":
		if query.Has("acl") {
			action := "s3:GetBucketAcl"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketAcl(ctx, bucket)
			return result, http.StatusOK, err
		}
		if query.Has("versioning") {
			action := "s3:GetBucketVersioning"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketVersioning(ctx, &GetBucketVersioningInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("encryption") {
			action := "s3:GetEncryptionConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketEncryption(ctx, &GetBucketEncryptionInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("policy") {
			action := "s3:GetBucketPolicy"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketPolicy(ctx, &GetBucketPolicyInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("cors") {
			action := "s3:GetBucketCORS"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketCORS(ctx, &GetBucketCORSInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("tagging") {
			action := "s3:GetBucketTagging"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketTagging(ctx, &GetBucketTaggingInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("lifecycle") {
			action := "s3:GetLifecycleConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLifecycleConfiguration(ctx, &GetBucketLifecycleConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("website") {
			action := "s3:GetBucketWebsite"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketWebsite(ctx, &GetBucketWebsiteInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("replication") {
			return nil, http.StatusNotFound, ErrNoSuchReplication
		}
		if query.Has("object-lock") {
			action := "s3:GetBucketObjectLockConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetObjectLockConfiguration(ctx, &GetObjectLockConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("notification") {
			action := "s3:GetBucketNotification"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketNotificationConfiguration(ctx, &GetBucketNotificationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("logging") {
			action := "s3:GetBucketLogging"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLogging(ctx, &GetBucketLoggingInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("ownershipControls") {
			action := "s3:GetBucketOwnershipControls"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketOwnershipControls(ctx, &GetBucketOwnershipControlsInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("requestPayment") {
			action := "s3:GetBucketRequestPayment"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketRequestPayment(ctx, &GetBucketRequestPaymentInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("accelerate") {
			action := "s3:GetAccelerateConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketAccelerateConfiguration(ctx, &GetBucketAccelerateConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("location") {
			action := "s3:GetBucketLocation"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLocation(ctx, &GetBucketLocationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("versions") {
			action := "s3:ListBucketVersions"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			input := &ListObjectVersionsInput{
				Bucket:          bucket,
				Delimiter:       query.Get("delimiter"),
				Prefix:          query.Get("prefix"),
				KeyMarker:       query.Get("key-marker"),
				VersionIdMarker: query.Get("version-id-marker"),
			}
			if maxKeys := query.Get("max-keys"); maxKeys != "" {
				mk, err := strconv.Atoi(maxKeys)
				if err != nil {
					return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
				}
				input.MaxKeys = mk
			}
			result, err := h.objectOps.ListObjectVersions(r.Context(), ctx, input)
			return result, http.StatusOK, err
		}
		if query.Has("uploads") {
			action := "s3:ListBucketMultipartUploads"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			input := &ListMultipartUploadsInput{
				Bucket:         bucket,
				Delimiter:      query.Get("delimiter"),
				Prefix:         query.Get("prefix"),
				KeyMarker:      query.Get("key-marker"),
				UploadIdMarker: query.Get("upload-id-marker"),
			}
			if maxUploads := query.Get("max-uploads"); maxUploads != "" {
				mu, err := strconv.Atoi(maxUploads)
				if err != nil {
					return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-uploads not an integer")
				}
				input.MaxUploads = mu
			}
			result, err := h.objectOps.ListMultipartUploads(r.Context(), ctx, input)
			return result, http.StatusOK, err
		}
		if query.Has("list-type") && query.Get("list-type") == "2" {
			action := "s3:ListBucket"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			input := &ListObjectsV2Input{
				Bucket:            bucket,
				Delimiter:         query.Get("delimiter"),
				Prefix:            query.Get("prefix"),
				ContinuationToken: query.Get("continuation-token"),
				StartAfter:        query.Get("start-after"),
			}
			if maxKeys := query.Get("max-keys"); maxKeys != "" {
				mk, err := strconv.Atoi(maxKeys)
				if err != nil {
					return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
				}
				input.MaxKeys = mk
			}
			result, err := h.objectOps.ListObjectsV2(r.Context(), ctx, input)
			return result, http.StatusOK, err
		}
		action := "s3:ListBucket"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		input := &ListObjectsInput{
			Bucket:    bucket,
			Delimiter: query.Get("delimiter"),
			Prefix:    query.Get("prefix"),
			Marker:    query.Get("marker"),
		}
		if maxKeys := query.Get("max-keys"); maxKeys != "" {
			mk, err := strconv.Atoi(maxKeys)
			if err != nil {
				return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
			}
			input.MaxKeys = mk
		}
		result, err := h.objectOps.ListObjects(r.Context(), ctx, input)
		return result, http.StatusOK, err

	case "HEAD":
		action := "s3:ListBucket"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.HeadBucket(ctx, &HeadBucketInput{Bucket: bucket})
		return result, http.StatusOK, err

	case "DELETE":
		if query.Has("encryption") {
			action := "s3:DeleteEncryptionConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketEncryption(ctx, &DeleteBucketEncryptionInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("policy") {
			action := "s3:DeleteBucketPolicy"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketPolicy(ctx, &DeleteBucketPolicyInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("cors") {
			action := "s3:DeleteBucketCORS"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketCORS(ctx, &DeleteBucketCORSInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("tagging") {
			action := "s3:DeleteBucketTagging"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketTagging(ctx, &DeleteBucketTaggingInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("lifecycle") {
			action := "s3:DeleteLifecycleConfiguration"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketLifecycleConfiguration(ctx, &DeleteBucketLifecycleConfigurationInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("website") {
			action := "s3:DeleteBucketWebsite"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketWebsite(ctx, &DeleteBucketWebsiteInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("ownershipControls") {
			action := "s3:DeleteBucketOwnershipControls"
			if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketOwnershipControls(ctx, &DeleteBucketOwnershipControlsInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		action := "s3:DeleteBucket"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucket(ctx, &DeleteBucketInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	default:
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed")
	}
}
