package s3

import (
	"net/http"
	"testing"
)

// cloudTrailEventNames and s3IAMActions are transcribed from the AWS
// CloudTrail S3 event tables and the "Required permissions for Amazon S3
// API operations" table.  Every expectation in the classification table
// below must appear in these sets, so a mistyped name fails the test even
// when the table and the classifier agree with each other.
var cloudTrailEventNames = map[string]bool{
	"ListBuckets":                      true,
	"CreateBucket":                     true,
	"DeleteBucket":                     true,
	"HeadBucket":                       true,
	"ListObjects":                      true,
	"ListObjectVersions":               true,
	"ListParts":                        true,
	"ListMultipartUploads":             true, // not logged by CloudTrail; API name kept
	"PutBucketAcl":                     true,
	"GetBucketAcl":                     true,
	"PutBucketVersioning":              true,
	"GetBucketVersioning":              true,
	"PutBucketEncryption":              true,
	"GetBucketEncryption":              true,
	"DeleteBucketEncryption":           true,
	"PutBucketPolicy":                  true,
	"GetBucketPolicy":                  true,
	"DeleteBucketPolicy":               true,
	"GetBucketPolicyStatus":            true,
	"PutBucketCors":                    true,
	"GetBucketCors":                    true,
	"DeleteBucketCors":                 true,
	"PutBucketTagging":                 true,
	"GetBucketTagging":                 true,
	"DeleteBucketTagging":              true,
	"PutBucketLifecycle":               true,
	"GetBucketLifecycle":               true,
	"DeleteBucketLifecycle":            true,
	"PutBucketWebsite":                 true,
	"GetBucketWebsite":                 true,
	"DeleteBucketWebsite":              true,
	"PutBucketObjectLockConfiguration": true,
	"GetBucketObjectLockConfiguration": true,
	"PutBucketNotification":            true,
	"GetBucketNotification":            true,
	"PutBucketLogging":                 true,
	"GetBucketLogging":                 true,
	"PutBucketOwnershipControls":       true,
	"GetBucketOwnershipControls":       true,
	"DeleteBucketOwnershipControls":    true,
	"PutBucketRequestPayment":          true,
	"GetBucketRequestPayment":          true,
	"PutAccelerateConfiguration":       true,
	"GetAccelerateConfiguration":       true,
	"GetBucketLocation":                true,
	"PutBucketPublicAccessBlock":       true,
	"GetBucketPublicAccessBlock":       true,
	"DeleteBucketPublicAccessBlock":    true,
	"PutBucketReplication":             true,
	"GetBucketReplication":             true,
	"DeleteBucketReplication":          true,
	"DeleteObjects":                    true,
	"PutObject":                        true,
	"GetObject":                        true,
	"HeadObject":                       true,
	"DeleteObject":                     true,
	"CopyObject":                       true,
	"UploadPart":                       true,
	"UploadPartCopy":                   true,
	"CreateMultipartUpload":            true,
	"CompleteMultipartUpload":          true,
	"AbortMultipartUpload":             true,
	"GetObjectAcl":                     true,
	"PutObjectAcl":                     true,
	"GetObjectTagging":                 true,
	"PutObjectTagging":                 true,
	"DeleteObjectTagging":              true,
	"GetObjectLegalHold":               true,
	"PutObjectLegalHold":               true,
	"GetObjectRetention":               true,
	"PutObjectRetention":               true,
	"GetObjectAttributes":              true,
	"RestoreObject":                    true,
	"SelectObjectContent":              true,
	"UpdateObjectEncryption":           true,
}

var s3IAMActions = map[string]bool{
	"s3:ListAllMyBuckets":                 true,
	"s3:CreateBucket":                     true,
	"s3:DeleteBucket":                     true,
	"s3:ListBucket":                       true,
	"s3:ListBucketVersions":               true,
	"s3:ListBucketMultipartUploads":       true,
	"s3:ListMultipartUploadParts":         true,
	"s3:PutBucketAcl":                     true,
	"s3:GetBucketAcl":                     true,
	"s3:PutBucketVersioning":              true,
	"s3:GetBucketVersioning":              true,
	"s3:PutEncryptionConfiguration":       true,
	"s3:GetEncryptionConfiguration":       true,
	"s3:PutBucketPolicy":                  true,
	"s3:GetBucketPolicy":                  true,
	"s3:DeleteBucketPolicy":               true,
	"s3:GetBucketPolicyStatus":            true,
	"s3:PutBucketCORS":                    true,
	"s3:GetBucketCORS":                    true,
	"s3:PutBucketTagging":                 true,
	"s3:GetBucketTagging":                 true,
	"s3:PutLifecycleConfiguration":        true,
	"s3:GetLifecycleConfiguration":        true,
	"s3:PutBucketWebsite":                 true,
	"s3:GetBucketWebsite":                 true,
	"s3:DeleteBucketWebsite":              true,
	"s3:PutBucketObjectLockConfiguration": true,
	"s3:GetBucketObjectLockConfiguration": true,
	"s3:PutBucketNotification":            true,
	"s3:GetBucketNotification":            true,
	"s3:PutBucketLogging":                 true,
	"s3:GetBucketLogging":                 true,
	"s3:PutBucketOwnershipControls":       true,
	"s3:GetBucketOwnershipControls":       true,
	"s3:PutBucketRequestPayment":          true,
	"s3:GetBucketRequestPayment":          true,
	"s3:PutAccelerateConfiguration":       true,
	"s3:GetAccelerateConfiguration":       true,
	"s3:GetBucketLocation":                true,
	"s3:PutBucketPublicAccessBlock":       true,
	"s3:GetBucketPublicAccessBlock":       true,
	"s3:PutReplicationConfiguration":      true,
	"s3:GetReplicationConfiguration":      true,
	"s3:PutObject":                        true,
	"s3:GetObject":                        true,
	"s3:GetObjectVersion":                 true,
	"s3:DeleteObject":                     true,
	"s3:DeleteObjectVersion":              true,
	"s3:GetObjectVersionAcl":              true,
	"s3:PutObjectVersionAcl":              true,
	"s3:GetObjectVersionTagging":          true,
	"s3:PutObjectVersionTagging":          true,
	"s3:DeleteObjectVersionTagging":       true,
	"s3:AbortMultipartUpload":             true,
	"s3:GetObjectAcl":                     true,
	"s3:PutObjectAcl":                     true,
	"s3:GetObjectTagging":                 true,
	"s3:PutObjectTagging":                 true,
	"s3:DeleteObjectTagging":              true,
	"s3:GetObjectLegalHold":               true,
	"s3:PutObjectLegalHold":               true,
	"s3:GetObjectRetention":               true,
	"s3:PutObjectRetention":               true,
	"s3:RestoreObject":                    true,
	"s3:UpdateObjectEncryption":           true,
}

type classifyCase struct {
	name       string
	method     string
	path       string
	rawQuery   string
	header     map[string]string
	wantEvent  string
	wantAction []s3ActionSpec // nil means no actions are emitted
}

func bucketCase(method, subresource, event, action string) classifyCase {
	name := method + " bucket"
	if subresource != "" {
		name += " ?" + subresource
	}
	query := subresource
	if subresource == "list-type=2" {
		query = subresource
	}
	return classifyCase{
		name:       name,
		method:     method,
		path:       "/bucket",
		rawQuery:   query,
		wantEvent:  event,
		wantAction: []s3ActionSpec{{Action: action, Bucket: "bucket"}},
	}
}

func objectCase(method, subresource, event, action string, header map[string]string) classifyCase {
	name := method + " object"
	if subresource != "" {
		name += " ?" + subresource
	}
	if header != nil {
		name += " +copy-source"
	}
	return classifyCase{
		name:       name,
		method:     method,
		path:       "/bucket/key",
		rawQuery:   subresource,
		header:     header,
		wantEvent:  event,
		wantAction: []s3ActionSpec{{Action: action, Bucket: "bucket", Key: "key"}},
	}
}

func TestClassifyS3Request(t *testing.T) {
	cases := []classifyCase{
		// Service level
		{name: "GET service", method: "GET", path: "/", wantEvent: "ListBuckets", wantAction: []s3ActionSpec{{Action: "s3:ListAllMyBuckets"}}},

		// Bucket PUT sub-resources
		bucketCase("PUT", "acl", "PutBucketAcl", "s3:PutBucketAcl"),
		bucketCase("PUT", "versioning", "PutBucketVersioning", "s3:PutBucketVersioning"),
		bucketCase("PUT", "encryption", "PutBucketEncryption", "s3:PutEncryptionConfiguration"),
		bucketCase("PUT", "policy", "PutBucketPolicy", "s3:PutBucketPolicy"),
		bucketCase("PUT", "cors", "PutBucketCors", "s3:PutBucketCORS"),
		bucketCase("PUT", "tagging", "PutBucketTagging", "s3:PutBucketTagging"),
		bucketCase("PUT", "lifecycle", "PutBucketLifecycle", "s3:PutLifecycleConfiguration"),
		bucketCase("PUT", "website", "PutBucketWebsite", "s3:PutBucketWebsite"),
		bucketCase("PUT", "object-lock", "PutBucketObjectLockConfiguration", "s3:PutBucketObjectLockConfiguration"),
		bucketCase("PUT", "notification", "PutBucketNotification", "s3:PutBucketNotification"),
		bucketCase("PUT", "logging", "PutBucketLogging", "s3:PutBucketLogging"),
		bucketCase("PUT", "ownershipControls", "PutBucketOwnershipControls", "s3:PutBucketOwnershipControls"),
		bucketCase("PUT", "requestPayment", "PutBucketRequestPayment", "s3:PutBucketRequestPayment"),
		bucketCase("PUT", "accelerate", "PutAccelerateConfiguration", "s3:PutAccelerateConfiguration"),
		bucketCase("PUT", "publicAccessBlock", "PutBucketPublicAccessBlock", "s3:PutBucketPublicAccessBlock"),
		bucketCase("PUT", "replication", "PutBucketReplication", "s3:PutReplicationConfiguration"),
		bucketCase("PUT", "", "CreateBucket", "s3:CreateBucket"),

		// Bucket GET sub-resources
		bucketCase("GET", "acl", "GetBucketAcl", "s3:GetBucketAcl"),
		bucketCase("GET", "versioning", "GetBucketVersioning", "s3:GetBucketVersioning"),
		bucketCase("GET", "encryption", "GetBucketEncryption", "s3:GetEncryptionConfiguration"),
		bucketCase("GET", "policy", "GetBucketPolicy", "s3:GetBucketPolicy"),
		bucketCase("GET", "policyStatus", "GetBucketPolicyStatus", "s3:GetBucketPolicyStatus"),
		bucketCase("GET", "cors", "GetBucketCors", "s3:GetBucketCORS"),
		bucketCase("GET", "tagging", "GetBucketTagging", "s3:GetBucketTagging"),
		bucketCase("GET", "lifecycle", "GetBucketLifecycle", "s3:GetLifecycleConfiguration"),
		bucketCase("GET", "website", "GetBucketWebsite", "s3:GetBucketWebsite"),
		bucketCase("GET", "replication", "GetBucketReplication", "s3:GetReplicationConfiguration"),
		bucketCase("GET", "object-lock", "GetBucketObjectLockConfiguration", "s3:GetBucketObjectLockConfiguration"),
		bucketCase("GET", "notification", "GetBucketNotification", "s3:GetBucketNotification"),
		bucketCase("GET", "logging", "GetBucketLogging", "s3:GetBucketLogging"),
		bucketCase("GET", "ownershipControls", "GetBucketOwnershipControls", "s3:GetBucketOwnershipControls"),
		bucketCase("GET", "requestPayment", "GetBucketRequestPayment", "s3:GetBucketRequestPayment"),
		bucketCase("GET", "accelerate", "GetAccelerateConfiguration", "s3:GetAccelerateConfiguration"),
		bucketCase("GET", "location", "GetBucketLocation", "s3:GetBucketLocation"),
		bucketCase("GET", "publicAccessBlock", "GetBucketPublicAccessBlock", "s3:GetBucketPublicAccessBlock"),
		bucketCase("GET", "versions", "ListObjectVersions", "s3:ListBucketVersions"),
		bucketCase("GET", "uploads", "ListMultipartUploads", "s3:ListBucketMultipartUploads"),
		bucketCase("GET", "list-type=2", "ListObjects", "s3:ListBucket"),
		bucketCase("GET", "", "ListObjects", "s3:ListBucket"),

		// Bucket HEAD / DELETE
		bucketCase("HEAD", "", "HeadBucket", "s3:ListBucket"),
		bucketCase("DELETE", "encryption", "DeleteBucketEncryption", "s3:PutEncryptionConfiguration"),
		bucketCase("DELETE", "policy", "DeleteBucketPolicy", "s3:DeleteBucketPolicy"),
		bucketCase("DELETE", "cors", "DeleteBucketCors", "s3:PutBucketCORS"),
		bucketCase("DELETE", "tagging", "DeleteBucketTagging", "s3:PutBucketTagging"),
		bucketCase("DELETE", "lifecycle", "DeleteBucketLifecycle", "s3:PutLifecycleConfiguration"),
		bucketCase("DELETE", "website", "DeleteBucketWebsite", "s3:DeleteBucketWebsite"),
		bucketCase("DELETE", "ownershipControls", "DeleteBucketOwnershipControls", "s3:PutBucketOwnershipControls"),
		bucketCase("DELETE", "publicAccessBlock", "DeleteBucketPublicAccessBlock", "s3:PutBucketPublicAccessBlock"),
		bucketCase("DELETE", "replication", "DeleteBucketReplication", "s3:PutReplicationConfiguration"),
		bucketCase("DELETE", "", "DeleteBucket", "s3:DeleteBucket"),

		// Batch delete authorises per object inside its handler
		{name: "POST bucket ?delete", method: "POST", path: "/bucket", rawQuery: "delete", wantEvent: "DeleteObjects"},

		// Object operations
		objectCase("POST", "select", "SelectObjectContent", "s3:GetObject", nil),
		objectCase("POST", "restore", "RestoreObject", "s3:RestoreObject", nil),
		objectCase("PUT", "encryption", "UpdateObjectEncryption", "s3:UpdateObjectEncryption", nil),
		objectCase("POST", "uploads", "CreateMultipartUpload", "s3:PutObject", nil),
		// Copy operations require the put-side action on the destination
		// and the read action on the source object.
		{
			name: "PUT object ?uploadId&partNumber +copy-source", method: "PUT", path: "/bucket/key",
			rawQuery: "uploadId=1&partNumber=1", header: map[string]string{"x-amz-copy-source": "/src/src-key"},
			wantEvent: "UploadPartCopy",
			wantAction: []s3ActionSpec{
				{Action: "s3:PutObject", Bucket: "bucket", Key: "key"},
				{Action: "s3:GetObject", Bucket: "src", Key: "src-key"},
			},
		},
		objectCase("PUT", "uploadId=1&partNumber=1", "UploadPart", "s3:PutObject", nil),
		objectCase("GET", "uploadId=1", "ListParts", "s3:ListMultipartUploadParts", nil),
		objectCase("POST", "uploadId=1", "CompleteMultipartUpload", "s3:PutObject", nil),
		objectCase("DELETE", "uploadId=1", "AbortMultipartUpload", "s3:AbortMultipartUpload", nil),
		objectCase("GET", "acl", "GetObjectAcl", "s3:GetObjectAcl", nil),
		objectCase("PUT", "acl", "PutObjectAcl", "s3:PutObjectAcl", nil),
		objectCase("GET", "acl&versionId=v1", "GetObjectAcl", "s3:GetObjectVersionAcl", nil),
		objectCase("PUT", "acl&versionId=v1", "PutObjectAcl", "s3:PutObjectVersionAcl", nil),
		objectCase("GET", "tagging", "GetObjectTagging", "s3:GetObjectTagging", nil),
		objectCase("PUT", "tagging", "PutObjectTagging", "s3:PutObjectTagging", nil),
		objectCase("DELETE", "tagging", "DeleteObjectTagging", "s3:DeleteObjectTagging", nil),
		objectCase("GET", "tagging&versionId=v1", "GetObjectTagging", "s3:GetObjectVersionTagging", nil),
		objectCase("PUT", "tagging&versionId=v1", "PutObjectTagging", "s3:PutObjectVersionTagging", nil),
		objectCase("DELETE", "tagging&versionId=v1", "DeleteObjectTagging", "s3:DeleteObjectVersionTagging", nil),
		objectCase("GET", "legal-hold", "GetObjectLegalHold", "s3:GetObjectLegalHold", nil),
		objectCase("PUT", "legal-hold", "PutObjectLegalHold", "s3:PutObjectLegalHold", nil),
		objectCase("GET", "retention", "GetObjectRetention", "s3:GetObjectRetention", nil),
		objectCase("PUT", "retention", "PutObjectRetention", "s3:PutObjectRetention", nil),
		objectCase("GET", "attributes", "GetObjectAttributes", "s3:GetObject", nil),
		objectCase("GET", "attributes&versionId=v1", "GetObjectAttributes", "s3:GetObjectVersion", nil),
		objectCase("GET", "", "GetObject", "s3:GetObject", nil),
		objectCase("GET", "versionId=v1", "GetObject", "s3:GetObjectVersion", nil),
		objectCase("HEAD", "", "HeadObject", "s3:GetObject", nil),
		objectCase("HEAD", "versionId=v1", "HeadObject", "s3:GetObjectVersion", nil),
		{
			name: "PUT object +copy-source", method: "PUT", path: "/bucket/key",
			header:    map[string]string{"x-amz-copy-source": "/src/src-key"},
			wantEvent: "CopyObject",
			wantAction: []s3ActionSpec{
				{Action: "s3:PutObject", Bucket: "bucket", Key: "key"},
				{Action: "s3:GetObject", Bucket: "src", Key: "src-key"},
			},
		},
		{
			name: "PUT object +copy-source with version", method: "PUT", path: "/bucket/key",
			header: map[string]string{
				"x-amz-copy-source":            "/src/src-key?versionId=vid",
				"x-amz-copy-source-version-id": "vid2",
			},
			wantEvent: "CopyObject",
			wantAction: []s3ActionSpec{
				{Action: "s3:PutObject", Bucket: "bucket", Key: "key"},
				{Action: "s3:GetObjectVersion", Bucket: "src", Key: "src-key"},
			},
		},
		objectCase("PUT", "", "PutObject", "s3:PutObject", nil),
		objectCase("DELETE", "", "DeleteObject", "s3:DeleteObject", nil),
		objectCase("DELETE", "versionId=v1", "DeleteObject", "s3:DeleteObjectVersion", nil),

		// Unrouted methods fall through without actions
		{name: "PATCH object", method: "PATCH", path: "/bucket/key", wantEvent: "UnknownS3Operation"},
		{name: "POST bucket", method: "POST", path: "/bucket", wantEvent: "UnknownS3Operation"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}
			req.URL.RawQuery = tc.rawQuery
			for k, v := range tc.header {
				req.Header.Set(k, v)
			}
			bucket, key := parseS3Path(req.URL.Path)

			gotEvent, gotActions := classifyS3Request(req, bucket, key)
			if gotEvent != tc.wantEvent {
				t.Errorf("event = %q, want %q", gotEvent, tc.wantEvent)
			}
			if len(gotActions) != len(tc.wantAction) {
				t.Fatalf("actions = %v, want %v", gotActions, tc.wantAction)
			}
			for i, want := range tc.wantAction {
				if gotActions[i] != want {
					t.Errorf("actions[%d] = %+v, want %+v", i, gotActions[i], want)
				}
			}

			// Cross-check both transcriptions: every expected name must be
			// an AWS-verified CloudTrail event / IAM action.
			if tc.wantEvent != "UnknownS3Operation" && !cloudTrailEventNames[tc.wantEvent] {
				t.Errorf("expected event %q is not in the verified CloudTrail event set", tc.wantEvent)
			}
			for _, a := range tc.wantAction {
				if !s3IAMActions[a.Action] {
					t.Errorf("expected action %q is not in the verified IAM action set", a.Action)
				}
			}
		})
	}
}
