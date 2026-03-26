// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

// S3Handler handles S3 HTTP requests.
type S3Handler struct {
	svc            *S3Service
	bucketOps      *BucketOperations
	objectOps      *ObjectOperations
	region         string
	storageManager *storage.RegionStorageManager
}

// NewS3Handler creates a new S3 handler.
func NewS3Handler(svc *S3Service, region string, storageMgr *storage.RegionStorageManager) *S3Handler {
	return &S3Handler{
		svc:            svc,
		bucketOps:      NewBucketOperations(svc),
		objectOps:      NewObjectOperations(svc),
		region:         region,
		storageManager: storageMgr,
	}
}

func (h *S3Handler) newRequestContext(r *http.Request) *request.RequestContext {
	ctx := request.NewRequestContext(r.Context(), h.storageManager, h.svc.accountID, h.region)
	ctx.SourceIP = extractSourceIP(r)
	ctx.UserAgent = r.UserAgent()
	ctx.Principal = h.svc.accountID
	ctx.PrincipalID = h.svc.accountID
	ctx.PrincipalType = request.PrincipalTypeUser
	return ctx
}

func extractSourceIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// ServeHTTP implements the http.Handler interface to handle S3 HTTP requests.
func (h *S3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := h.newRequestContext(r)

	requestID := fmt.Sprintf("%016X", time.Now().UnixNano())
	w.Header().Set("x-amz-request-id", requestID)
	w.Header().Set("x-amz-id-2", requestID)

	bucket, key := parseS3Path(r.URL.Path)
	query := r.URL.Query()

	if h.isPresignedURLRequest(query) {
		if err := h.verifyPresignedURL(r, bucket); err != nil {
			h.writeError(w, err, bucket, key)
			return
		}
	}

	var result interface{}
	var statusCode int
	var err error
	var header http.Header

	switch {
	case query.Has("delete") && bucket != "" && key == "":
		result, statusCode, err = h.handleDeleteObjects(reqCtx, r, bucket, r.Body)
	case bucket == "":
		result, statusCode, err = h.handleServiceRequest(reqCtx, r)
	case key == "":
		result, statusCode, err = h.handleBucketRequest(reqCtx, r, bucket)
	default:
		result, header, statusCode, err = h.handleObjectRequest(reqCtx, r, bucket, key)
	}

	if err != nil {
		h.writeError(w, err, bucket, key)
		return
	}

	for k, v := range header {
		w.Header()[k] = v
	}

	switch v := result.(type) {
	case *SelectObjectContentOutput:
		if v.Payload != nil {
			defer v.Payload.Close()
			_, _ = io.Copy(w, v.Payload)
		}
		return
	case *GetObjectOutput:
		if v.AcceptRanges != "" {
			w.Header().Set("Accept-Ranges", v.AcceptRanges)
		}
		if v.ContentRange != "" {
			w.Header().Set("Content-Range", v.ContentRange)
		}
		if v.IsPartial {
			w.WriteHeader(http.StatusPartialContent)
		} else {
			w.WriteHeader(statusCode)
		}
		if v.Body != nil {
			defer v.Body.Close()
			_, _ = io.Copy(w, v.Body)
		}
		return
	case *HeadObjectOutput:
		w.WriteHeader(statusCode)
		return
	case *ListBucketsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectsV2Output:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectVersionsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListMultipartUploadsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListPartsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *GetBucketAclOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *GetObjectAclOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *CopyObjectOutput:
		h.writeXMLResponse(w, "CopyObjectResult", v.CopyObjectResult, statusCode, "")
	case *UploadPartCopyOutput:
		h.writeXMLResponse(w, "CopyPartResult", v.CopyPartResult, statusCode, "")
	case *DeleteObjectsOutput:
		h.writeXMLResponse(w, "DeleteResult", v, statusCode, "")
	case *GetBucketNotificationOutput:
		h.writeXMLResponse(w, "NotificationConfiguration", v.NotificationConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
	case *GetBucketLoggingOutput:
		if v.LoggingConfiguration != nil {
			h.writeXMLResponse(w, "BucketLoggingStatus", v, statusCode, "")
		} else {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus xmlns="http://s3.amazonaws.com/doc/2006-03-01/"></BucketLoggingStatus>`))
		}
	case *GetObjectAttributesOutput:
		h.writeXMLResponse(w, "GetObjectAttributesOutput", v, statusCode, "")
	case *GetBucketEncryptionOutput:
		if v.ServerSideEncryptionConfiguration != nil {
			h.writeXMLResponse(w, "ServerSideEncryptionConfiguration", v.ServerSideEncryptionConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case nil:
		w.WriteHeader(statusCode)
	default:
		h.writeXMLResponse(w, "", v, statusCode, "")
	}
}

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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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

		_, err := h.bucketOps.CreateBucket(ctx, &CreateBucketInput{
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
		return nil, http.StatusOK, err

	case "GET":
		if query.Has("acl") {
			action := "s3:GetBucketAcl"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketAcl(ctx, bucket)
			return result, http.StatusOK, err
		}
		if query.Has("versioning") {
			action := "s3:GetBucketVersioning"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketVersioning(ctx, &GetBucketVersioningInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("encryption") {
			action := "s3:GetEncryptionConfiguration"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketEncryption(ctx, &GetBucketEncryptionInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("policy") {
			action := "s3:GetBucketPolicy"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketPolicy(ctx, &GetBucketPolicyInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("cors") {
			action := "s3:GetBucketCORS"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketCORS(ctx, &GetBucketCORSInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("tagging") {
			action := "s3:GetBucketTagging"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketTagging(ctx, &GetBucketTaggingInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("lifecycle") {
			action := "s3:GetLifecycleConfiguration"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLifecycleConfiguration(ctx, &GetBucketLifecycleConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("website") {
			action := "s3:GetBucketWebsite"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetObjectLockConfiguration(ctx, &GetObjectLockConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("notification") {
			action := "s3:GetBucketNotification"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketNotificationConfiguration(ctx, &GetBucketNotificationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("logging") {
			action := "s3:GetBucketLogging"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLogging(ctx, &GetBucketLoggingInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("ownershipControls") {
			action := "s3:GetBucketOwnershipControls"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketOwnershipControls(ctx, &GetBucketOwnershipControlsInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("requestPayment") {
			action := "s3:GetBucketRequestPayment"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketRequestPayment(ctx, &GetBucketRequestPaymentInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("accelerate") {
			action := "s3:GetAccelerateConfiguration"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketAccelerateConfiguration(ctx, &GetBucketAccelerateConfigurationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("location") {
			action := "s3:GetBucketLocation"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			result, err := h.bucketOps.GetBucketLocation(ctx, &GetBucketLocationInput{Bucket: bucket})
			return result, http.StatusOK, err
		}
		if query.Has("versions") {
			action := "s3:ListBucketVersions"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
		if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
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
		if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.HeadBucket(ctx, &HeadBucketInput{Bucket: bucket})
		return nil, http.StatusOK, err

	case "DELETE":
		if query.Has("encryption") {
			action := "s3:DeleteEncryptionConfiguration"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketEncryption(ctx, &DeleteBucketEncryptionInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("policy") {
			action := "s3:DeleteBucketPolicy"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketPolicy(ctx, &DeleteBucketPolicyInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("cors") {
			action := "s3:DeleteBucketCORS"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketCORS(ctx, &DeleteBucketCORSInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("tagging") {
			action := "s3:DeleteBucketTagging"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketTagging(ctx, &DeleteBucketTaggingInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("lifecycle") {
			action := "s3:DeleteLifecycleConfiguration"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketLifecycleConfiguration(ctx, &DeleteBucketLifecycleConfigurationInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("website") {
			action := "s3:DeleteBucketWebsite"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketWebsite(ctx, &DeleteBucketWebsiteInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		if query.Has("ownershipControls") {
			action := "s3:DeleteBucketOwnershipControls"
			if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
				return nil, http.StatusForbidden, err
			}
			err := h.bucketOps.DeleteBucketOwnershipControls(ctx, &DeleteBucketOwnershipControlsInput{Bucket: bucket})
			return nil, http.StatusNoContent, err
		}
		action := "s3:DeleteBucket"
		if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucket(ctx, &DeleteBucketInput{Bucket: bucket})
		return nil, http.StatusNoContent, err

	default:
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed")
	}
}

func (h *S3Handler) checkAccess(ctx *request.RequestContext, stores *s3Stores, action, bucket, key string) error {
	check := &AccessCheck{
		Principal:     ctx.Principal,
		PrincipalID:   ctx.PrincipalID,
		PrincipalType: ctx.PrincipalType,
		Action:        action,
		Resource:      buildResource(h.svc.accountID, h.region, bucket, key),
		Bucket:        bucket,
		Key:           key,
		SourceIP:      ctx.SourceIP,
	}

	if key != "" {
		return h.svc.accessController.CheckObjectAccess(ctx, stores, check)
	}
	return h.svc.accessController.CheckAccess(ctx, stores, check)
}

func (h *S3Handler) handleObjectRequest(ctx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	action := determineObjectAction(r)
	if err := h.checkAccess(ctx, stores, action, bucket, key); err != nil {
		return nil, nil, http.StatusForbidden, err
	}

	if r.Method == "POST" && r.URL.Query().Has("select") {
		return h.handleSelectObjectContent(ctx, r, bucket, key)
	}

	return h.objectOps.HandleRequest(r.Context(), ctx, r, bucket, key)
}

func (h *S3Handler) handleSelectObjectContent(ctx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	var input SelectObjectContentInput
	if err := request.NewSafeXMLDecoder(r.Body).Decode(&input); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("failed to decode request: %w", err)
	}
	input.Bucket = bucket
	input.Key = key

	store, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	objReader, obj, err := store.objects.Get(ctx, bucket, key)
	if err != nil {
		return nil, nil, http.StatusNotFound, err
	}
	defer objReader.Close()

	var dataReader io.Reader = objReader

	if obj.SSEMetadata != nil {
		encryptedData, encObj, err := store.objects.GetEncrypted(ctx, bucket, key, "")
		if err != nil {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to get encrypted object: %w", err)
		}

		if encObj.SSEMetadata.EncryptionType == "CUSTOMER" {
			if input.SSECustomerKey == "" {
				return nil, nil, http.StatusBadRequest, fmt.Errorf("customer key is required for SSE-C encrypted object")
			}
			customerKey, err := h.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, nil, http.StatusBadRequest, fmt.Errorf("invalid SSE-C customer key: %w", err)
			}
			decResult, err := h.svc.encryptionManager.DecryptWithCustomerKey(encryptedData, encObj.SSEMetadata, bucket, key, customerKey)
			if err != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", err)
			}
			dataReader = bytes.NewReader(decResult.DecryptedData)
		} else {
			decResult, err := h.svc.encryptionManager.Decrypt(encryptedData, encObj.SSEMetadata, bucket, key)
			if err != nil {
				return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to decrypt data: %w", err)
			}
			dataReader = bytes.NewReader(decResult.DecryptedData)
		}
	}

	engine, err := NewSelectEngine(&input)
	if err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	pr, pw := io.Pipe()
	header := make(http.Header)
	header.Set("Content-Type", "application/vnd.amazon.eventstream")
	header.Set("x-amz-request-id", fmt.Sprintf("%016X", time.Now().UnixNano()))

	go func() {
		defer pw.Close()
		writer := NewSelectEventStreamWriter(pw, input.RequestProgress)
		if err := engine.Execute(ctx, dataReader, writer); err != nil {
			pw.CloseWithError(err)
			return
		}
		if err := writer.WriteStats(); err != nil {
			pw.CloseWithError(err)
			return
		}
		if err := writer.WriteEnd(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	output := &SelectObjectContentOutput{
		Payload: pr,
	}

	return output, header, http.StatusOK, nil
}

func (h *S3Handler) handleDeleteObjects(ctx *request.RequestContext, r *http.Request, bucket string, body io.Reader) (interface{}, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	action := "s3:DeleteObject"
	if err := h.checkAccess(ctx, stores, action, bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}

	result, err := h.objectOps.HandleDeleteObjects(r.Context(), ctx, bucket, body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil
}

func determineObjectAction(r *http.Request) string {
	method := r.Method
	query := r.URL.Query()

	switch {
	case method == "POST" && query.Has("select"):
		return "s3:GetObject"
	case method == "GET" && query.Has("acl"):
		return "s3:GetObjectAcl"
	case method == "PUT" && query.Has("acl"):
		return "s3:PutObjectAcl"
	case method == "GET" && query.Has("tagging"):
		return "s3:GetObjectTagging"
	case method == "PUT" && query.Has("tagging"):
		return "s3:PutObjectTagging"
	case method == "DELETE" && query.Has("tagging"):
		return "s3:DeleteObjectTagging"
	case method == "GET" && query.Has("legal-hold"):
		return "s3:GetObjectLegalHold"
	case method == "PUT" && query.Has("legal-hold"):
		return "s3:PutObjectLegalHold"
	case method == "GET" && query.Has("retention"):
		return "s3:GetObjectRetention"
	case method == "PUT" && query.Has("retention"):
		return "s3:PutObjectRetention"
	case method == "GET" && query.Has("uploadId"):
		return "s3:ListMultipartUploadParts"
	case method == "PUT" && query.Has("uploadId"):
		return "s3:UploadPart"
	case method == "POST" && query.Has("uploadId"):
		return "s3:CompleteMultipartUpload"
	case method == "DELETE" && query.Has("uploadId"):
		return "s3:AbortMultipartUpload"
	case method == "POST" && query.Has("uploads"):
		return "s3:CreateMultipartUpload"
	case method == "POST" && query.Has("restore"):
		return "s3:RestoreObject"
	case method == "PUT" && r.Header.Get("x-amz-copy-source") != "" && !query.Has("uploadId"):
		return "s3:PutObject"
	case method == "GET":
		return "s3:GetObject"
	case method == "HEAD":
		return "s3:GetObject"
	case method == "PUT":
		return "s3:PutObject"
	case method == "DELETE":
		return "s3:DeleteObject"
	default:
		return "s3:PutObject"
	}
}

// #nosec G705
func (h *S3Handler) writeXMLResponse(w http.ResponseWriter, rootElement string, data interface{}, statusCode int, xmlns string) {
	var xmlData []byte
	var err error

	if rootElement != "" {
		xmlData, err = xml.Marshal(data)
		if err != nil {
			h.writeError(w, err, "", "")
			return
		}
		xmlStr := string(xmlData)
		if idx := strings.Index(xmlStr, ">"); idx != -1 {
			xmlStr = xmlStr[idx+1:]
		}
		if idx := strings.LastIndex(xmlStr, "<"); idx != -1 {
			xmlStr = xmlStr[:idx]
		}
		if xmlns != "" {
			xmlData = []byte(fmt.Sprintf(`<%s xmlns="%s">%s</%s>`, rootElement, xmlns, xmlStr, rootElement))
		} else {
			xmlData = []byte(fmt.Sprintf(`<%s>%s</%s>`, rootElement, xmlStr, rootElement))
		}
	} else {
		xmlData, err = xml.Marshal(data)
	}

	if err != nil {
		h.writeError(w, err, "", "")
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write(xmlData)
}

// #nosec G705
func (h *S3Handler) writeError(w http.ResponseWriter, err error, bucket, key string) {
	var awsErr *awserrors.AWSError

	switch {
	case errors.Is(err, storecommon.ErrNotFound):
		if bucket != "" && key == "" {
			awsErr = NewNoSuchBucketError(bucket)
		} else if key != "" {
			awsErr = NewNoSuchKeyError(key)
		} else {
			awsErr = ErrNoSuchBucket
		}
	case errors.Is(err, storecommon.ErrAlreadyExists):
		awsErr = NewBucketAlreadyExistsError(bucket)
	case errors.Is(err, storecommon.ErrConflict) || (err != nil && strings.Contains(err.Error(), "not empty")):
		awsErr = ErrBucketNotEmpty
	case errors.Is(err, storecommon.ErrInvalidInput):
		awsErr = ErrInvalidRequest
	case errors.Is(err, ErrPreconditionFailed):
		awsErr = ErrPreconditionFailed
	default:
		var castErr *awserrors.AWSError
		if errors.As(err, &castErr) {
			awsErr = castErr
		} else {
			awsErr = awserrors.NewAWSError("InternalError", err.Error(), http.StatusInternalServerError)
		}
	}

	resource := bucket
	if key != "" {
		resource = bucket + "/" + key
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(awsErr.HTTPStatus)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write([]byte(fmt.Sprintf(`<Error><Code>%s</Code><Message>%s</Message><Resource>%s</Resource><RequestId>%s</RequestId></Error>`,
		escapeXMLValue(awsErr.Code), escapeXMLValue(awsErr.Message), escapeXMLValue(resource), fmt.Sprintf("%016X", time.Now().UnixNano()))))
}

func escapeXMLValue(s string) string {
	var buf strings.Builder
	xml.Escape(&buf, []byte(s))
	return buf.String()
}

func parseS3Path(urlPath string) (bucket, key string) {
	urlPath = strings.TrimPrefix(urlPath, "/")
	if urlPath == "" {
		return "", ""
	}

	parts := strings.SplitN(urlPath, "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return bucket, key
}

func (h *S3Handler) isPresignedURLRequest(query url.Values) bool {
	return query.Has("X-Amz-Algorithm") &&
		query.Has("X-Amz-Credential") &&
		query.Has("X-Amz-Signature")
}

func (h *S3Handler) verifyPresignedURL(r *http.Request, bucket string) error {
	verifier := crypto.NewPresignedURLVerifier(h.svc.credentialsProvider)
	return verifier.VerifyPresignedURL(r, bucket, h.region)
}
