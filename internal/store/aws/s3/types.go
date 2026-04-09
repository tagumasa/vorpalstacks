package s3

// Package s3 provides S3 (Simple Storage Service) data store implementations
// for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// BucketCannedACL represents a canned ACL for an S3 bucket.
type BucketCannedACL string

// BucketCannedACL constants define the supported canned ACLs for buckets.
const (
	BucketCannedACLPrivate           BucketCannedACL = "private"
	BucketCannedACLPublicRead        BucketCannedACL = "public-read"
	BucketCannedACLPublicReadWrite   BucketCannedACL = "public-read-write"
	BucketCannedACLAuthenticatedRead BucketCannedACL = "authenticated-read"
	BucketCannedACLLogDeliveryWrite  BucketCannedACL = "log-delivery-write"
)

// ObjectCannedACL represents a canned ACL for an S3 object.
type ObjectCannedACL string

// ObjectCannedACL constants define the supported canned ACLs for objects.
const (
	ObjectCannedACLPrivate                ObjectCannedACL = "private"
	ObjectCannedACLPublicRead             ObjectCannedACL = "public-read"
	ObjectCannedACLPublicReadWrite        ObjectCannedACL = "public-read-write"
	ObjectCannedACLAuthenticatedRead      ObjectCannedACL = "authenticated-read"
	ObjectCannedACLBucketOwnerRead        ObjectCannedACL = "bucket-owner-read"
	ObjectCannedACLBucketOwnerFullControl ObjectCannedACL = "bucket-owner-full-control"
	ObjectCannedACLAwsExecRead            ObjectCannedACL = "aws-exec-read"
)

// Permission represents an S3 permission.
type Permission string

// Permission constants define the supported permissions.
const (
	PermissionFullControl Permission = "FULL_CONTROL"
	PermissionRead        Permission = "READ"
	PermissionWrite       Permission = "WRITE"
	PermissionReadACP     Permission = "READ_ACP"
	PermissionWriteACP    Permission = "WRITE_ACP"
)

// GranteeType represents the type of grantee in an S3 ACL.
type GranteeType string

// GranteeType constants define the supported grantee types.
const (
	GranteeTypeCanonicalUser         GranteeType = "CanonicalUser"
	GranteeTypeGroup                 GranteeType = "Group"
	GranteeTypeAmazonCustomerByEmail GranteeType = "AmazonCustomerByEmail"
)

// Predefined S3 group URIs.
const (
	AuthenticatedUsersGroup = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
	AllUsersGroup           = "http://acs.amazonaws.com/groups/global/AllUsers"
	LogDeliveryGroup        = "http://acs.amazonaws.com/groups/s3/LogDelivery"
)

// ACLOwner represents the owner of an S3 bucket or object.
type ACLOwner struct {
	ID          string `json:"id,omitempty" xml:"ID"`
	DisplayName string `json:"display_name,omitempty" xml:"DisplayName,omitempty"`
}

// Grantee represents a grantee in an S3 ACL.
type Grantee struct {
	Type        GranteeType `json:"type" xml:"xsi:type,attr"`
	ID          string      `json:"id,omitempty" xml:"ID,omitempty"`
	DisplayName string      `json:"display_name,omitempty" xml:"DisplayName,omitempty"`
	Email       string      `json:"email,omitempty" xml:"EmailAddress,omitempty"`
	URI         string      `json:"uri,omitempty" xml:"URI,omitempty"`
}

// Grant represents a single grant in an S3 ACL.
type Grant struct {
	Grantee    *Grantee   `json:"grantee,omitempty" xml:"Grantee"`
	Permission Permission `json:"permission" xml:"Permission"`
}

// AccessControlPolicy represents the access control policy for an S3 bucket or object.
type AccessControlPolicy struct {
	Owner  *ACLOwner `json:"owner,omitempty" xml:"Owner"`
	Grants []*Grant  `json:"grants,omitempty" xml:"AccessControlList>Grant"`
}

// ObjectStorageClass represents the storage class for an S3 object.
type ObjectStorageClass string

// ObjectStorageClass constants define the supported storage classes.
const (
	StorageClassStandard           ObjectStorageClass = "STANDARD"
	StorageClassReducedRedundancy  ObjectStorageClass = "REDUCED_REDUNDANCY"
	StorageClassGlacier            ObjectStorageClass = "GLACIER"
	StorageClassStandardIA         ObjectStorageClass = "STANDARD_IA"
	StorageClassOneZoneIA          ObjectStorageClass = "ONEZONE_IA"
	StorageClassIntelligentTiering ObjectStorageClass = "INTELLIGENT_TIERING"
)

// BucketVersioningStatus represents the versioning status of an S3 bucket.
type BucketVersioningStatus string

// BucketVersioningStatus constants define the possible versioning statuses.
const (
	BucketVersioningEnabled   BucketVersioningStatus = "Enabled"
	BucketVersioningSuspended BucketVersioningStatus = "Suspended"
)

// Bucket represents an S3 bucket.
type Bucket struct {
	Name                      string                       `json:"name"`
	Region                    string                       `json:"region"`
	CreationDate              time.Time                    `json:"creation_date"`
	ACL                       *AccessControlPolicy         `json:"acl,omitempty"`
	ObjectLockEnabled         bool                         `json:"object_lock_enabled,omitempty"`
	ObjectLockConfig          *ObjectLockConfiguration     `json:"object_lock_config,omitempty"`
	VersioningStatus          BucketVersioningStatus       `json:"versioning_status,omitempty"`
	EncryptionConfig          *EncryptionConfig            `json:"encryption_config,omitempty"`
	LifecycleConfiguration    *LifecycleConfiguration      `json:"lifecycle_configuration,omitempty"`
	WebsiteConfiguration      *WebsiteConfiguration        `json:"website_configuration,omitempty"`
	CORSConfiguration         *CORSConfiguration           `json:"cors_configuration,omitempty"`
	Policy                    string                       `json:"policy,omitempty"`
	PublicAccessBlock         *PublicAccessBlockConfig     `json:"public_access_block,omitempty"`
	Tags                      []common.Tag                 `json:"tags,omitempty"`
	NotificationConfiguration *NotificationConfiguration   `json:"notification_configuration,omitempty"`
	LoggingConfiguration      *LoggingConfiguration        `json:"logging_configuration,omitempty"`
	OwnershipControls         *OwnershipControls           `json:"ownership_controls,omitempty"`
	RequestPayment            *RequestPaymentConfiguration `json:"request_payment,omitempty"`
	AccelerateConfiguration   *AccelerateConfiguration     `json:"accelerate_configuration,omitempty"`
}

// EncryptionConfig represents the encryption configuration for an S3 bucket.
type EncryptionConfig struct {
	SSEAlgorithm     string `json:"sse_algorithm"`
	KMSMasterKeyID   string `json:"kms_master_key_id,omitempty"`
	BucketKeyEnabled *bool  `json:"bucket_key_enabled,omitempty"`
}

// LifecycleConfiguration represents the lifecycle configuration for an S3 bucket.
type LifecycleConfiguration struct {
	Rules []LifecycleRule `json:"rules"`
}

// LifecycleRule represents a lifecycle rule for an S3 bucket.
type LifecycleRule struct {
	ID                             string                        `json:"id"`
	Status                         string                        `json:"status"`
	Filter                         *LifecycleRuleFilter          `json:"filter,omitempty"`
	Expiration                     *LifecycleExpiration          `json:"expiration,omitempty"`
	Transitions                    []LifecycleTransition         `json:"transitions,omitempty"`
	NoncurrentVersionExpiration    *NoncurrentVersionExpiration  `json:"noncurrent_version_expiration,omitempty"`
	NoncurrentVersionTransitions   []NoncurrentVersionTransition `json:"noncurrent_version_transitions,omitempty"`
	AbortIncompleteMultipartUpload *AbortIncompleteUpload        `json:"abort_incomplete_multipart_upload,omitempty"`
}

// LifecycleRuleFilter represents the filter criteria for a lifecycle rule.
type LifecycleRuleFilter struct {
	Prefix                string                    `json:"prefix,omitempty"`
	ObjectSizeGreaterThan *int64                    `json:"object_size_greater_than,omitempty"`
	ObjectSizeLessThan    *int64                    `json:"object_size_less_than,omitempty"`
	And                   *LifecycleRuleAndOperator `json:"and,omitempty"`
	Tag                   *common.Tag               `json:"tag,omitempty"`
}

// LifecycleRuleAndOperator represents a logical AND operator for lifecycle rule filters.
type LifecycleRuleAndOperator struct {
	Prefix                string       `json:"prefix,omitempty"`
	Tags                  []common.Tag `json:"tags,omitempty"`
	ObjectSizeGreaterThan *int64       `json:"object_size_greater_than,omitempty"`
	ObjectSizeLessThan    *int64       `json:"object_size_less_than,omitempty"`
}

// LifecycleExpiration represents the expiration configuration for a lifecycle rule.
type LifecycleExpiration struct {
	Date                      *time.Time `json:"date,omitempty"`
	Days                      *int       `json:"days,omitempty"`
	ExpiredObjectDeleteMarker *bool      `json:"expired_object_delete_marker,omitempty"`
}

// LifecycleTransition represents the transition configuration for a lifecycle rule.
type LifecycleTransition struct {
	Date         *time.Time         `json:"date,omitempty"`
	Days         *int               `json:"days,omitempty"`
	StorageClass ObjectStorageClass `json:"storage_class"`
}

// NoncurrentVersionExpiration represents the expiration configuration for noncurrent object versions.
type NoncurrentVersionExpiration struct {
	NoncurrentDays          *int `json:"noncurrent_days,omitempty"`
	NewerNoncurrentVersions *int `json:"newer_noncurrent_versions,omitempty"`
}

// NoncurrentVersionTransition represents the transition configuration for noncurrent object versions.
type NoncurrentVersionTransition struct {
	NoncurrentDays          *int               `json:"noncurrent_days,omitempty"`
	NewerNoncurrentVersions *int               `json:"newer_noncurrent_versions,omitempty"`
	StorageClass            ObjectStorageClass `json:"storage_class"`
}

// AbortIncompleteUpload represents the configuration for aborting incomplete multipart uploads.
type AbortIncompleteUpload struct {
	DaysAfterInitiation *int `json:"days_after_initiation,omitempty"`
}

// WebsiteConfiguration represents the website configuration for an S3 bucket.
type WebsiteConfiguration struct {
	IndexDocument         string                 `json:"index_document,omitempty"`
	ErrorDocument         string                 `json:"error_document,omitempty"`
	RedirectAllRequestsTo *RedirectAllRequestsTo `json:"redirect_all_requests_to,omitempty"`
	RoutingRules          []RoutingRule          `json:"routing_rules,omitempty"`
}

// RedirectAllRequestsTo represents the redirect configuration for all requests to an S3 bucket.
type RedirectAllRequestsTo struct {
	HostName string `json:"host_name"`
	Protocol string `json:"protocol,omitempty"`
}

// RoutingRule represents a routing rule for S3 bucket website configuration.
type RoutingRule struct {
	Condition *RoutingRuleCondition `json:"condition,omitempty"`
	Redirect  *RoutingRuleRedirect  `json:"redirect"`
}

// RoutingRuleCondition represents the condition for a routing rule.
type RoutingRuleCondition struct {
	HTTPErrorCodeReturnedEquals *string `json:"http_error_code_returned_equals,omitempty"`
	KeyPrefixEquals             *string `json:"key_prefix_equals,omitempty"`
}

// RoutingRuleRedirect represents the redirect configuration for a routing rule.
type RoutingRuleRedirect struct {
	HostName             *string `json:"host_name,omitempty"`
	HTTPRedirectCode     *string `json:"http_redirect_code,omitempty"`
	Protocol             *string `json:"protocol,omitempty"`
	ReplaceKeyPrefixWith *string `json:"replace_key_prefix_with,omitempty"`
	ReplaceKeyWith       *string `json:"replace_key_with,omitempty"`
}

// CORSConfiguration represents the CORS configuration for an S3 bucket.
type CORSConfiguration struct {
	CORSRules []CORSRule `json:"cors_rules"`
}

// CORSRule represents a CORS rule for an S3 bucket.
type CORSRule struct {
	AllowedHeaders []string `json:"allowed_headers"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedOrigins []string `json:"allowed_origins"`
	ExposeHeaders  []string `json:"expose_headers,omitempty"`
	MaxAgeSeconds  *int     `json:"max_age_seconds,omitempty"`
	ID             string   `json:"id,omitempty"`
}

// PublicAccessBlockConfig represents the public access block configuration for an S3 bucket.
type PublicAccessBlockConfig struct {
	BlockPublicAcls       bool `json:"block_public_acls"`
	BlockPublicPolicy     bool `json:"block_public_policy"`
	IgnorePublicAcls      bool `json:"ignore_public_acls"`
	RestrictPublicBuckets bool `json:"restrict_public_buckets"`
}

// SSEType represents the server-side encryption type.
type SSEType string

// SSEType constants define the supported server-side encryption types.
const (
	SSETypeAES256   SSEType = "AES256"
	SSETypeKMS      SSEType = "aws:kms"
	SSETypeKMSES    SSEType = "aws:kms:dsse"
	SSETypeCustomer SSEType = "CUSTOMER"
)

// Object represents an S3 object.
type Object struct {
	Key                  string               `json:"key"`
	BucketName           string               `json:"bucket_name"`
	Size                 int64                `json:"size"`
	ETag                 string               `json:"etag"`
	LastModified         time.Time            `json:"last_modified"`
	ContentType          string               `json:"content_type,omitempty"`
	ContentEncoding      string               `json:"content_encoding,omitempty"`
	ContentLanguage      string               `json:"content_language,omitempty"`
	ContentDisposition   string               `json:"content_disposition,omitempty"`
	CacheControl         string               `json:"cache_control,omitempty"`
	Expires              *time.Time           `json:"expires,omitempty"`
	Metadata             map[string]string    `json:"metadata,omitempty"`
	StorageClass         ObjectStorageClass   `json:"storage_class,omitempty"`
	ServerSideEncryption string               `json:"server_side_encryption,omitempty"`
	SSEKMSKeyID          string               `json:"sse_kms_key_id,omitempty"`
	VersionID            string               `json:"version_id,omitempty"`
	IsLatest             bool                 `json:"is_latest"`
	IsDeleteMarker       bool                 `json:"is_delete_marker"`
	Tags                 []common.Tag         `json:"tags,omitempty"`
	ACL                  *AccessControlPolicy `json:"acl,omitempty"`
	Owner                *ACLOwner            `json:"owner,omitempty"`
	ObjectLockLegalHold  *ObjectLockLegalHold `json:"object_lock_legal_hold,omitempty"`
	ObjectLockRetention  *ObjectLockRetention `json:"object_lock_retention,omitempty"`
	SSEMetadata          *SSEObjectMetadata   `json:"sse_metadata,omitempty"`
}

// SSEObjectMetadata represents the server-side encryption metadata for an S3 object.
type SSEObjectMetadata struct {
	EncryptionType    SSEType `json:"encryption_type"`
	EncryptedDataKey  []byte  `json:"encrypted_data_key,omitempty"`
	ContentNonce      []byte  `json:"content_nonce,omitempty"`
	KMSKeyID          string  `json:"kms_key_id,omitempty"`
	EncryptionContext string  `json:"encryption_context,omitempty"`
	UnencryptedMD5    string  `json:"unencrypted_md5,omitempty"`
	UnencryptedSize   int64   `json:"unencrypted_size,omitempty"`
}

// MultipartUpload represents a multipart upload to S3.
type MultipartUpload struct {
	UploadID         string             `json:"upload_id"`
	Key              string             `json:"key"`
	BucketName       string             `json:"bucket_name"`
	Initiated        time.Time          `json:"initiated"`
	StorageClass     ObjectStorageClass `json:"storage_class,omitempty"`
	Owner            string             `json:"owner,omitempty"`
	Initiator        string             `json:"initiator,omitempty"`
	Parts            []ObjectPart       `json:"parts,omitempty"`
	partIndex        map[int]int        `json:"-"` // PartNumber -> Parts slice index (not serialised)
	Metadata         map[string]string  `json:"metadata,omitempty"`
	ContentType      string             `json:"content_type,omitempty"`
	SSEMetadata      *SSEObjectMetadata `json:"sse_metadata,omitempty"`
	SSEType          SSEType            `json:"sse_type,omitempty"`
	KMSKeyID         string             `json:"kms_key_id,omitempty"`
	CustomerKeyMD5   string             `json:"customer_key_md5,omitempty"`
	PlaintextDataKey []byte             `json:"plaintext_data_key,omitempty"`
}

func (u *MultipartUpload) rebuildPartIndex() {
	u.partIndex = make(map[int]int, len(u.Parts))
	for i, p := range u.Parts {
		u.partIndex[p.PartNumber] = i
	}
}

// FindPart finds a part by its part number.
func (u *MultipartUpload) FindPart(partNumber int) (int, bool) {
	if u.partIndex == nil {
		u.rebuildPartIndex()
	}
	idx, ok := u.partIndex[partNumber]
	return idx, ok
}

// AddPart adds a part to the multipart upload.
func (u *MultipartUpload) AddPart(partNumber int, part ObjectPart) {
	if u.partIndex == nil {
		u.rebuildPartIndex()
	}
	if idx, ok := u.partIndex[partNumber]; ok {
		u.Parts[idx] = part
	} else {
		u.Parts = append(u.Parts, part)
		u.partIndex[partNumber] = len(u.Parts) - 1
	}
}

// ObjectPart represents a part in a multipart upload.
type ObjectPart struct {
	PartNumber    int       `json:"part_number"`
	ETag          string    `json:"etag"`
	Size          int64     `json:"size"`
	LastModified  time.Time `json:"last_modified"`
	EncryptedSize int64     `json:"encrypted_size,omitempty"`
	ContentNonce  []byte    `json:"content_nonce,omitempty"`
	DataKey       []byte    `json:"data_key,omitempty"`
}

// BucketListResult represents the result of listing S3 buckets.
type BucketListResult struct {
	Buckets     []*Bucket
	NextMarker  string
	IsTruncated bool
}

// ObjectListResult represents the result of listing S3 objects.
type ObjectListResult struct {
	Objects              []*Object
	CommonPrefixes       []string
	NextMarker           string
	NextVersionKeyMarker string
	NextVersionIDMarker  string
	IsTruncated          bool
}

// MultipartUploadListResult represents the result of listing multipart uploads.
type MultipartUploadListResult struct {
	Uploads            []*MultipartUpload
	NextKeyMarker      string
	NextUploadIDMarker string
	IsTruncated        bool
}

// ObjectLockConfiguration represents the object lock configuration for an S3 bucket.
type ObjectLockConfiguration struct {
	ObjectLockEnabled string          `json:"object_lock_enabled,omitempty"`
	Rule              *ObjectLockRule `json:"rule,omitempty"`
}

// ObjectLockRule represents the object lock rule for an S3 bucket.
type ObjectLockRule struct {
	DefaultRetention *DefaultRetention `json:"default_retention,omitempty"`
}

// DefaultRetention represents the default retention settings for object lock.
type DefaultRetention struct {
	Mode  ObjectLockRetentionMode `json:"mode,omitempty"`
	Days  *int                    `json:"days,omitempty"`
	Years *int                    `json:"years,omitempty"`
}

// ObjectLockRetentionMode represents the object lock retention mode.
type ObjectLockRetentionMode string

// ObjectLockRetentionMode constants define the supported retention modes.
const (
	ObjectLockRetentionModeGovernance ObjectLockRetentionMode = "GOVERNANCE"
	ObjectLockRetentionModeCompliance ObjectLockRetentionMode = "COMPLIANCE"
)

// ObjectLockLegalHold represents the legal hold status for an S3 object.
type ObjectLockLegalHold struct {
	Status ObjectLockLegalHoldStatus `json:"status"`
}

// ObjectLockLegalHoldStatus represents the legal hold status for an S3 object.
type ObjectLockLegalHoldStatus string

// ObjectLockLegalHoldStatus constants define the supported legal hold statuses.
const (
	ObjectLockLegalHoldOn  ObjectLockLegalHoldStatus = "ON"
	ObjectLockLegalHoldOff ObjectLockLegalHoldStatus = "OFF"
)

// ObjectLockRetention represents the object lock retention settings for an S3 object.
type ObjectLockRetention struct {
	Mode            ObjectLockRetentionMode `json:"mode"`
	RetainUntilDate *time.Time              `json:"retain_until_date,omitempty"`
}

// NotificationConfiguration represents the notification configuration for an S3 bucket.
type NotificationConfiguration struct {
	TopicConfigurations  []TopicNotificationConfiguration  `json:"topic_configurations,omitempty"`
	QueueConfigurations  []QueueNotificationConfiguration  `json:"queue_configurations,omitempty"`
	LambdaConfigurations []LambdaNotificationConfiguration `json:"lambda_configurations,omitempty"`
}

// TopicNotificationConfiguration represents an SNS topic notification configuration.
type TopicNotificationConfiguration struct {
	Id       string                           `json:"id,omitempty"`
	TopicArn string                           `json:"topic_arn"`
	Events   []string                         `json:"events"`
	Filter   *NotificationConfigurationFilter `json:"filter,omitempty"`
}

// QueueNotificationConfiguration represents an SQS queue notification configuration.
type QueueNotificationConfiguration struct {
	Id       string                           `json:"id,omitempty"`
	QueueArn string                           `json:"queue_arn"`
	Events   []string                         `json:"events"`
	Filter   *NotificationConfigurationFilter `json:"filter,omitempty"`
}

// LambdaNotificationConfiguration represents a Lambda function notification configuration.
type LambdaNotificationConfiguration struct {
	Id                string                           `json:"id,omitempty"`
	LambdaFunctionArn string                           `json:"lambda_function_arn"`
	Events            []string                         `json:"events"`
	Filter            *NotificationConfigurationFilter `json:"filter,omitempty"`
}

// NotificationConfigurationFilter represents a filter for notification configurations.
type NotificationConfigurationFilter struct {
	Key *S3KeyFilter `json:"key,omitempty"`
}

// S3KeyFilter represents a key-based filter for S3 notifications.
type S3KeyFilter struct {
	FilterRules []FilterRule `json:"filter_rules,omitempty"`
}

// FilterRule represents a filter rule for S3 notifications.
type FilterRule struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// LoggingConfiguration represents the logging configuration for an S3 bucket.
type LoggingConfiguration struct {
	TargetBucket string        `json:"target_bucket"`
	TargetPrefix string        `json:"target_prefix,omitempty"`
	TargetGrants []TargetGrant `json:"target_grants,omitempty"`
}

// TargetGrant represents a grant in the logging configuration.
type TargetGrant struct {
	Grantee    *Grantee   `json:"grantee,omitempty"`
	Permission Permission `json:"permission"`
}

// OwnershipControls specifies the ownership rules for a bucket.
type OwnershipControls struct {
	Rules []OwnershipControlsRule `json:"rules,omitempty"`
}

// OwnershipControlsRule specifies an ownership rule for a bucket.
type OwnershipControlsRule struct {
	ObjectOwnership string `json:"object_ownership,omitempty"`
}

// RequestPaymentConfiguration specifies who pays for the bucket's downloads.
type RequestPaymentConfiguration struct {
	Payer string `json:"payer,omitempty"`
}

// AccelerateConfiguration specifies the transfer acceleration configuration for a bucket.
type AccelerateConfiguration struct {
	Status string `json:"status,omitempty"`
}
