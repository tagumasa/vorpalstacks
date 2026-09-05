package s3

import "net/http"

// s3ActionSpec names one IAM policy action together with the resource
// (bucket or object) that the action is evaluated against.
type s3ActionSpec struct {
	Action string
	Bucket string
	Key    string
}

// CloudTrail records bucket-level API calls under event names that
// sometimes differ from the API operation name (for example the
// PutBucketLifecycleConfiguration API is logged as the PutBucketLifecycle
// event), and the IAM action required by an operation is frequently not
// derived from the operation name at all (DeleteBucketEncryption requires
// s3:PutEncryptionConfiguration, every multipart-upload operation requires
// s3:PutObject, HeadBucket requires s3:ListBucket).  Keeping those mappings
// in private copies next to each consumer has historically let newly added
// operations drift out of sync, so classifyS3Request is the single place
// that derives both the CloudTrail event name and the required IAM actions
// for a request.  The audit recorder and the pre-dispatch access checks
// must consume it instead of re-deriving names.
func classifyS3Request(r *http.Request, bucket, key string) (string, []s3ActionSpec) {
	if bucket == "" && key == "" {
		return "ListBuckets", []s3ActionSpec{{Action: "s3:ListAllMyBuckets"}}
	}
	if key == "" {
		return classifyBucketRequest(r, bucket)
	}
	return classifyObjectRequest(r, bucket, key)
}

func classifyBucketRequest(r *http.Request, bucket string) (string, []s3ActionSpec) {
	query := r.URL.Query()
	bucketAction := func(name string) []s3ActionSpec {
		return []s3ActionSpec{{Action: name, Bucket: bucket}}
	}

	switch r.Method {
	case "PUT":
		switch {
		case query.Has("acl"):
			return "PutBucketAcl", bucketAction("s3:PutBucketAcl")
		case query.Has("versioning"):
			return "PutBucketVersioning", bucketAction("s3:PutBucketVersioning")
		case query.Has("encryption"):
			return "PutBucketEncryption", bucketAction("s3:PutEncryptionConfiguration")
		case query.Has("policy"):
			return "PutBucketPolicy", bucketAction("s3:PutBucketPolicy")
		case query.Has("cors"):
			return "PutBucketCors", bucketAction("s3:PutBucketCORS")
		case query.Has("tagging"):
			return "PutBucketTagging", bucketAction("s3:PutBucketTagging")
		case query.Has("lifecycle"):
			return "PutBucketLifecycle", bucketAction("s3:PutLifecycleConfiguration")
		case query.Has("website"):
			return "PutBucketWebsite", bucketAction("s3:PutBucketWebsite")
		case query.Has("object-lock"):
			return "PutBucketObjectLockConfiguration", bucketAction("s3:PutBucketObjectLockConfiguration")
		case query.Has("notification"):
			return "PutBucketNotification", bucketAction("s3:PutBucketNotification")
		case query.Has("logging"):
			return "PutBucketLogging", bucketAction("s3:PutBucketLogging")
		case query.Has("ownershipControls"):
			return "PutBucketOwnershipControls", bucketAction("s3:PutBucketOwnershipControls")
		case query.Has("requestPayment"):
			return "PutBucketRequestPayment", bucketAction("s3:PutBucketRequestPayment")
		case query.Has("accelerate"):
			return "PutAccelerateConfiguration", bucketAction("s3:PutAccelerateConfiguration")
		case query.Has("publicAccessBlock"):
			return "PutBucketPublicAccessBlock", bucketAction("s3:PutBucketPublicAccessBlock")
		case query.Has("replication"):
			return "PutBucketReplication", bucketAction("s3:PutReplicationConfiguration")
		case query.Has("inventory"):
			return "PutBucketInventoryConfiguration", bucketAction("s3:PutInventoryConfiguration")
		case query.Has("metrics"):
			return "PutBucketMetricsConfiguration", bucketAction("s3:PutMetricsConfiguration")
		default:
			return "CreateBucket", bucketAction("s3:CreateBucket")
		}
	case "GET":
		switch {
		case query.Has("acl"):
			return "GetBucketAcl", bucketAction("s3:GetBucketAcl")
		case query.Has("versioning"):
			return "GetBucketVersioning", bucketAction("s3:GetBucketVersioning")
		case query.Has("encryption"):
			return "GetBucketEncryption", bucketAction("s3:GetEncryptionConfiguration")
		case query.Has("policy"):
			return "GetBucketPolicy", bucketAction("s3:GetBucketPolicy")
		case query.Has("policyStatus"):
			return "GetBucketPolicyStatus", bucketAction("s3:GetBucketPolicyStatus")
		case query.Has("cors"):
			return "GetBucketCors", bucketAction("s3:GetBucketCORS")
		case query.Has("tagging"):
			return "GetBucketTagging", bucketAction("s3:GetBucketTagging")
		case query.Has("lifecycle"):
			return "GetBucketLifecycle", bucketAction("s3:GetLifecycleConfiguration")
		case query.Has("website"):
			return "GetBucketWebsite", bucketAction("s3:GetBucketWebsite")
		case query.Has("replication"):
			return "GetBucketReplication", bucketAction("s3:GetReplicationConfiguration")
		case query.Has("inventory"):
			// The Delete and List variants of both configuration families
			// share the Put/Get IAM actions; the List operations are not
			// CloudTrail-logged and fall back to their API names.
			if query.Has("id") {
				return "GetBucketInventoryConfiguration", bucketAction("s3:GetInventoryConfiguration")
			}
			return "ListBucketInventoryConfigurations", bucketAction("s3:GetInventoryConfiguration")
		case query.Has("metrics"):
			if query.Has("id") {
				return "GetBucketMetricsConfiguration", bucketAction("s3:GetMetricsConfiguration")
			}
			return "ListBucketMetricsConfigurations", bucketAction("s3:GetMetricsConfiguration")
		case query.Has("object-lock"):
			return "GetBucketObjectLockConfiguration", bucketAction("s3:GetBucketObjectLockConfiguration")
		case query.Has("notification"):
			return "GetBucketNotification", bucketAction("s3:GetBucketNotification")
		case query.Has("logging"):
			return "GetBucketLogging", bucketAction("s3:GetBucketLogging")
		case query.Has("ownershipControls"):
			return "GetBucketOwnershipControls", bucketAction("s3:GetBucketOwnershipControls")
		case query.Has("requestPayment"):
			return "GetBucketRequestPayment", bucketAction("s3:GetBucketRequestPayment")
		case query.Has("accelerate"):
			return "GetAccelerateConfiguration", bucketAction("s3:GetAccelerateConfiguration")
		case query.Has("location"):
			return "GetBucketLocation", bucketAction("s3:GetBucketLocation")
		case query.Has("publicAccessBlock"):
			return "GetBucketPublicAccessBlock", bucketAction("s3:GetBucketPublicAccessBlock")
		case query.Has("versions"):
			return "ListObjectVersions", bucketAction("s3:ListBucketVersions")
		case query.Has("uploads"):
			// CloudTrail does not log ListMultipartUploads in either its
			// bucket-level or data-event lists, so the API name is used.
			return "ListMultipartUploads", bucketAction("s3:ListBucketMultipartUploads")
		default:
			// Both ListObjects and ListObjectsV2 are logged by CloudTrail
			// as the ListObjects data event on general purpose buckets.
			return "ListObjects", bucketAction("s3:ListBucket")
		}
	case "HEAD":
		return "HeadBucket", bucketAction("s3:ListBucket")
	case "DELETE":
		switch {
		case query.Has("encryption"):
			return "DeleteBucketEncryption", bucketAction("s3:PutEncryptionConfiguration")
		case query.Has("policy"):
			return "DeleteBucketPolicy", bucketAction("s3:DeleteBucketPolicy")
		case query.Has("cors"):
			return "DeleteBucketCors", bucketAction("s3:PutBucketCORS")
		case query.Has("tagging"):
			return "DeleteBucketTagging", bucketAction("s3:PutBucketTagging")
		case query.Has("lifecycle"):
			return "DeleteBucketLifecycle", bucketAction("s3:PutLifecycleConfiguration")
		case query.Has("website"):
			return "DeleteBucketWebsite", bucketAction("s3:DeleteBucketWebsite")
		case query.Has("ownershipControls"):
			return "DeleteBucketOwnershipControls", bucketAction("s3:PutBucketOwnershipControls")
		case query.Has("publicAccessBlock"):
			return "DeleteBucketPublicAccessBlock", bucketAction("s3:PutBucketPublicAccessBlock")
		case query.Has("replication"):
			return "DeleteBucketReplication", bucketAction("s3:PutReplicationConfiguration")
		case query.Has("inventory"):
			return "DeleteBucketInventoryConfiguration", bucketAction("s3:PutInventoryConfiguration")
		case query.Has("metrics"):
			return "DeleteBucketMetricsConfiguration", bucketAction("s3:PutMetricsConfiguration")
		default:
			return "DeleteBucket", bucketAction("s3:DeleteBucket")
		}
	case "POST":
		if query.Has("delete") {
			// DeleteObjects authorises each object individually inside its
			// own handler, so no request-level action is emitted here.
			return "DeleteObjects", nil
		}
	}
	return "UnknownS3Operation", nil
}

func classifyObjectRequest(r *http.Request, bucket, key string) (string, []s3ActionSpec) {
	query := r.URL.Query()
	objectAction := func(name string) []s3ActionSpec {
		return []s3ActionSpec{{Action: name, Bucket: bucket, Key: key}}
	}
	copySource := r.Header.Get("x-amz-copy-source")

	// Versioned requests require the *Version IAM action instead of the
	// base one (the required-permissions table pairs GetObject with
	// s3:GetObject or s3:GetObjectVersion depending on whether versionId
	// is specified; the same either/or applies to DeleteObject, the object
	// ACL operations, and the object tagging operations).
	getAction := "s3:GetObject"
	if query.Get("versionId") != "" {
		getAction = "s3:GetObjectVersion"
	}
	deleteAction := "s3:DeleteObject"
	if query.Get("versionId") != "" {
		deleteAction = "s3:DeleteObjectVersion"
	}
	getAclAction := "s3:GetObjectAcl"
	if query.Get("versionId") != "" {
		getAclAction = "s3:GetObjectVersionAcl"
	}
	putAclAction := "s3:PutObjectAcl"
	if query.Get("versionId") != "" {
		putAclAction = "s3:PutObjectVersionAcl"
	}
	getTaggingAction := "s3:GetObjectTagging"
	if query.Get("versionId") != "" {
		getTaggingAction = "s3:GetObjectVersionTagging"
	}
	putTaggingAction := "s3:PutObjectTagging"
	if query.Get("versionId") != "" {
		putTaggingAction = "s3:PutObjectVersionTagging"
	}
	deleteTaggingAction := "s3:DeleteObjectTagging"
	if query.Get("versionId") != "" {
		deleteTaggingAction = "s3:DeleteObjectVersionTagging"
	}

	switch {
	case r.Method == "POST" && query.Has("select"):
		return "SelectObjectContent", objectAction("s3:GetObject")
	case r.Method == "POST" && query.Has("restore"):
		return "RestoreObject", objectAction("s3:RestoreObject")
	case r.Method == "PUT" && query.Has("encryption"):
		return "UpdateObjectEncryption", objectAction("s3:UpdateObjectEncryption")
	case r.Method == "POST" && query.Has("uploads"):
		return "CreateMultipartUpload", objectAction("s3:PutObject")
	case r.Method == "PUT" && query.Has("uploadId") && query.Has("partNumber") && copySource != "":
		return "UploadPartCopy", append(objectAction("s3:PutObject"), copySourceReadAction(r)...)
	case r.Method == "PUT" && query.Has("uploadId") && query.Has("partNumber"):
		return "UploadPart", objectAction("s3:PutObject")
	case r.Method == "GET" && query.Has("uploadId"):
		return "ListParts", objectAction("s3:ListMultipartUploadParts")
	case r.Method == "POST" && query.Has("uploadId"):
		return "CompleteMultipartUpload", objectAction("s3:PutObject")
	case r.Method == "DELETE" && query.Has("uploadId"):
		return "AbortMultipartUpload", objectAction("s3:AbortMultipartUpload")
	case r.Method == "GET" && query.Has("acl"):
		return "GetObjectAcl", objectAction(getAclAction)
	case r.Method == "PUT" && query.Has("acl"):
		return "PutObjectAcl", objectAction(putAclAction)
	case r.Method == "GET" && query.Has("tagging"):
		return "GetObjectTagging", objectAction(getTaggingAction)
	case r.Method == "PUT" && query.Has("tagging"):
		return "PutObjectTagging", objectAction(putTaggingAction)
	case r.Method == "DELETE" && query.Has("tagging"):
		return "DeleteObjectTagging", objectAction(deleteTaggingAction)
	case r.Method == "GET" && query.Has("legal-hold"):
		return "GetObjectLegalHold", objectAction("s3:GetObjectLegalHold")
	case r.Method == "PUT" && query.Has("legal-hold"):
		return "PutObjectLegalHold", objectAction("s3:PutObjectLegalHold")
	case r.Method == "GET" && query.Has("retention"):
		return "GetObjectRetention", objectAction("s3:GetObjectRetention")
	case r.Method == "PUT" && query.Has("retention"):
		return "PutObjectRetention", objectAction("s3:PutObjectRetention")
	case r.Method == "GET" && query.Has("attributes"):
		return "GetObjectAttributes", objectAction(getAction)
	case r.Method == "GET":
		return "GetObject", objectAction(getAction)
	case r.Method == "HEAD":
		return "HeadObject", objectAction(getAction)
	case r.Method == "PUT" && copySource != "" && !query.Has("uploadId"):
		return "CopyObject", append(objectAction("s3:PutObject"), copySourceReadAction(r)...)
	case r.Method == "PUT":
		return "PutObject", objectAction("s3:PutObject")
	case r.Method == "DELETE":
		return "DeleteObject", objectAction(deleteAction)
	}
	return "UnknownS3Operation", nil
}

// copySourceReadAction returns the read action required on the copy source
// object: s3:GetObjectVersion when a specific source version is requested
// and s3:GetObject otherwise.  A nil result means the copy source could not
// be parsed; the dispatcher rejects such requests before any data moves.
func copySourceReadAction(r *http.Request) []s3ActionSpec {
	srcBucket, srcKey, srcVersionId, err := parseCopySource(r.Header.Get("x-amz-copy-source"))
	if err != nil {
		return nil
	}
	if v := r.Header.Get("x-amz-copy-source-version-id"); v != "" {
		srcVersionId = v
	}
	action := "s3:GetObject"
	if srcVersionId != "" {
		action = "s3:GetObjectVersion"
	}
	return []s3ActionSpec{{Action: action, Bucket: srcBucket, Key: srcKey}}
}
