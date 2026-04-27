package s3

import "net/http"

// determineS3EventName maps an S3 HTTP request to a CloudTrail event name
// based on HTTP method, query parameters, and path.
func determineS3EventName(r *http.Request, bucket, key string) string {
	if bucket == "" && key == "" {
		if r.Method == "GET" {
			return "ListBuckets"
		}
		return "ListBuckets"
	}

	if bucket != "" && key == "" {
		return determineBucketEventName(r)
	}

	return determineObjectEventName(r)
}

func determineBucketEventName(r *http.Request) string {
	query := r.URL.Query()
	method := r.Method

	switch method {
	case "PUT":
		switch {
		case query.Has("acl"):
			return "PutBucketAcl"
		case query.Has("versioning"):
			return "PutBucketVersioning"
		case query.Has("encryption"):
			return "PutEncryptionConfiguration"
		case query.Has("policy"):
			return "PutBucketPolicy"
		case query.Has("cors"):
			return "PutBucketCORS"
		case query.Has("tagging"):
			return "PutBucketTagging"
		case query.Has("lifecycle"):
			return "PutLifecycleConfiguration"
		case query.Has("website"):
			return "PutBucketWebsite"
		case query.Has("logging"):
			return "PutBucketLogging"
		case query.Has("requestPayment"):
			return "PutBucketRequestPayment"
		case query.Has("accelerate"):
			return "PutAccelerateConfiguration"
		case query.Has("publicAccessBlock"):
			return "PutBucketPublicAccessBlock"
		case query.Has("ownershipControls"):
			return "PutBucketOwnershipControls"
		case query.Has("object-lock"):
			return "PutBucketObjectLockConfiguration"
		case query.Has("notification"):
			return "PutBucketNotification"
		default:
			return "CreateBucket"
		}
	case "GET":
		switch {
		case query.Has("acl"):
			return "GetBucketAcl"
		case query.Has("versioning"):
			return "GetBucketVersioning"
		case query.Has("encryption"):
			return "GetEncryptionConfiguration"
		case query.Has("policy"):
			return "GetBucketPolicy"
		case query.Has("cors"):
			return "GetBucketCORS"
		case query.Has("tagging"):
			return "GetBucketTagging"
		case query.Has("lifecycle"):
			return "GetLifecycleConfiguration"
		case query.Has("website"):
			return "GetBucketWebsite"
		case query.Has("location"):
			return "GetBucketLocation"
		case query.Has("logging"):
			return "GetBucketLogging"
		case query.Has("requestPayment"):
			return "GetBucketRequestPayment"
		case query.Has("accelerate"):
			return "GetAccelerateConfiguration"
		case query.Has("publicAccessBlock"):
			return "GetBucketPublicAccessBlock"
		case query.Has("ownershipControls"):
			return "GetBucketOwnershipControls"
		case query.Has("object-lock"):
			return "GetObjectLockConfiguration"
		case query.Has("notification"):
			return "GetBucketNotification"
		case query.Has("versions"):
			return "ListObjectVersions"
		case query.Has("uploads"):
			return "ListMultipartUploads"
		case query.Has("list-type"):
			return "ListObjectsV2"
		default:
			return "ListObjects"
		}
	case "HEAD":
		return "HeadBucket"
	case "DELETE":
		switch {
		case query.Has("policy"):
			return "DeleteBucketPolicy"
		case query.Has("cors"):
			return "DeleteBucketCORS"
		case query.Has("tagging"):
			return "DeleteBucketTagging"
		case query.Has("lifecycle"):
			return "DeleteLifecycleConfiguration"
		case query.Has("website"):
			return "DeleteBucketWebsite"
		case query.Has("encryption"):
			return "DeleteEncryptionConfiguration"
		case query.Has("publicAccessBlock"):
			return "DeleteBucketPublicAccessBlock"
		case query.Has("ownershipControls"):
			return "DeleteBucketOwnershipControls"
		default:
			return "DeleteBucket"
		}
	case "POST":
		if query.Has("delete") {
			return "DeleteObjects"
		}
		return "DeleteObjects"
	}
	return "UnknownS3Operation"
}

func determineObjectEventName(r *http.Request) string {
	query := r.URL.Query()
	method := r.Method

	switch {
	case method == "POST" && query.Has("select"):
		return "SelectObjectContent"
	case method == "POST" && query.Has("restore"):
		return "RestoreObject"
	case method == "GET" && query.Has("acl"):
		return "GetObjectAcl"
	case method == "PUT" && query.Has("acl"):
		return "PutObjectAcl"
	case method == "GET" && query.Has("tagging"):
		return "GetObjectTagging"
	case method == "PUT" && query.Has("tagging"):
		return "PutObjectTagging"
	case method == "DELETE" && query.Has("tagging"):
		return "DeleteObjectTagging"
	case method == "GET" && query.Has("legal-hold"):
		return "GetObjectLegalHold"
	case method == "PUT" && query.Has("legal-hold"):
		return "PutObjectLegalHold"
	case method == "GET" && query.Has("retention"):
		return "GetObjectRetention"
	case method == "PUT" && query.Has("retention"):
		return "PutObjectRetention"
	case method == "GET" && query.Has("uploadId"):
		return "ListMultipartUploadParts"
	case method == "PUT" && query.Has("uploadId"):
		return "UploadPart"
	case method == "POST" && query.Has("uploadId"):
		return "CompleteMultipartUpload"
	case method == "DELETE" && query.Has("uploadId"):
		return "AbortMultipartUpload"
	case method == "POST" && query.Has("uploads"):
		return "CreateMultipartUpload"
	case method == "GET":
		return "GetObject"
	case method == "HEAD":
		return "HeadObject"
	case method == "PUT":
		return "PutObject"
	case method == "DELETE":
		return "DeleteObject"
	}
	return "UnknownS3Operation"
}
