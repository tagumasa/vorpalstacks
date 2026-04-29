// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// AccessKeyStatus represents the status of an IAM access key.
type AccessKeyStatus string

const (
	// AccessKeyStatusActive indicates the access key is active.
	AccessKeyStatusActive AccessKeyStatus = "Active"
	// AccessKeyStatusInactive indicates the access key is inactive.
	AccessKeyStatusInactive AccessKeyStatus = "Inactive"
)

// User represents an IAM user.
type User struct {
	ID                  string               `json:"id"`
	Path                string               `json:"path"`
	UserName            string               `json:"user_name"`
	Arn                 string               `json:"arn"`
	AccountId           string               `json:"account_id"`
	CreateDate          time.Time            `json:"create_date"`
	PasswordLastUsed    *time.Time           `json:"password_last_used,omitempty"`
	PermissionsBoundary *PermissionsBoundary `json:"permissions_boundary,omitempty"`
	Tags                []types.Tag          `json:"tags,omitempty"`
}

// AccessKey represents an IAM access key.
type AccessKey struct {
	AccessKeyId     string          `json:"access_key_id"`
	UserName        string          `json:"user_name"`
	Status          AccessKeyStatus `json:"status"`
	SecretAccessKey string          `json:"secret_access_key,omitempty"`
	CreateDate      time.Time       `json:"create_date"`
	LastUsedDate    *time.Time      `json:"last_used_date,omitempty"`
	LastUsedRegion  string          `json:"last_used_region,omitempty"`
	LastUsedService string          `json:"last_used_service,omitempty"`
}

// LoginProfile represents a login profile for an IAM user.
type LoginProfile struct {
	UserName              string    `json:"user_name"`
	PasswordHash          string    `json:"password_hash"`
	PasswordResetRequired bool      `json:"password_reset_required"`
	CreateDate            time.Time `json:"create_date"`
}

// PermissionsBoundary represents a permissions boundary for an IAM entity.
type PermissionsBoundary struct {
	PermissionsBoundaryType string `json:"permissions_boundary_type"`
	PermissionsBoundaryArn  string `json:"permissions_boundary_arn"`
}

// UserListResult represents the result of listing IAM users.
type UserListResult struct {
	Users       []*User
	IsTruncated bool
	Marker      string
}

// AccessKeyListResult represents the result of listing IAM access keys.
type AccessKeyListResult struct {
	AccessKeys  []*AccessKey
	IsTruncated bool
	Marker      string
}

// Group represents an IAM group.
type Group struct {
	ID                  string               `json:"id"`
	Path                string               `json:"path"`
	GroupName           string               `json:"group_name"`
	Arn                 string               `json:"arn"`
	AccountId           string               `json:"account_id"`
	CreateDate          time.Time            `json:"create_date"`
	PermissionsBoundary *PermissionsBoundary `json:"permissions_boundary,omitempty"`
	Tags                []types.Tag          `json:"tags,omitempty"`
}

// GroupListResult represents the result of listing IAM groups.
type GroupListResult struct {
	Groups      []*Group
	IsTruncated bool
	Marker      string
}

// UserGroupMembership represents membership of a user in a group.
type UserGroupMembership struct {
	UserName  string    `json:"user_name"`
	GroupName string    `json:"group_name"`
	JoinDate  time.Time `json:"join_date"`
}

// Role represents an IAM role.
type Role struct {
	ID                       string               `json:"id"`
	Path                     string               `json:"path"`
	RoleName                 string               `json:"role_name"`
	Arn                      string               `json:"arn"`
	AccountId                string               `json:"account_id"`
	CreateDate               time.Time            `json:"create_date"`
	AssumeRolePolicyDocument string               `json:"assume_role_policy_document,omitempty"`
	Description              string               `json:"description,omitempty"`
	MaxSessionDuration       int                  `json:"max_session_duration,omitempty"`
	PermissionsBoundary      *PermissionsBoundary `json:"permissions_boundary,omitempty"`
	Tags                     []types.Tag          `json:"tags,omitempty"`
	RoleLastUsed             *RoleLastUsed        `json:"role_last_used,omitempty"`
}

// RoleLastUsed represents information about when the role was last used.
type RoleLastUsed struct {
	LastUsedDate *time.Time `json:"last_used_date,omitempty"`
	Region       string     `json:"region,omitempty"`
}

// RoleListResult represents the result of listing IAM roles.
type RoleListResult struct {
	Roles       []*Role
	IsTruncated bool
	Marker      string
}

// InstanceProfile represents an IAM instance profile.
type InstanceProfile struct {
	ID                  string      `json:"id"`
	Path                string      `json:"path"`
	InstanceProfileName string      `json:"instance_profile_name"`
	Arn                 string      `json:"arn"`
	AccountId           string      `json:"account_id"`
	CreateDate          time.Time   `json:"create_date"`
	Roles               []string    `json:"roles,omitempty"`
	Tags                []types.Tag `json:"tags,omitempty"`
}

// InstanceProfileListResult represents the result of listing IAM instance profiles.
type InstanceProfileListResult struct {
	InstanceProfiles []*InstanceProfile
	IsTruncated      bool
	Marker           string
}

// RolePolicy represents an inline policy attached to an IAM role.
type RolePolicy struct {
	RoleName       string `json:"role_name"`
	PolicyName     string `json:"policy_name"`
	PolicyDocument string `json:"policy_document"`
}

// AttachedPolicy represents a managed policy attached to an IAM entity.
type AttachedPolicy struct {
	PolicyName string `json:"policy_name"`
	PolicyArn  string `json:"policy_arn"`
}

// Policy represents an IAM managed policy.
type Policy struct {
	ID               string      `json:"id"`
	Path             string      `json:"path"`
	PolicyName       string      `json:"policy_name"`
	Arn              string      `json:"arn"`
	AccountId        string      `json:"account_id"`
	CreateDate       time.Time   `json:"create_date"`
	UpdateDate       time.Time   `json:"update_date"`
	DefaultVersionId string      `json:"default_version_id"`
	AttachmentCount  int         `json:"attachment_count"`
	IsAttachable     bool        `json:"is_attachable"`
	Description      string      `json:"description,omitempty"`
	Tags             []types.Tag `json:"tags,omitempty"`
}

// PolicyVersion represents a version of an IAM managed policy.
type PolicyVersion struct {
	VersionId        string    `json:"version_id"`
	PolicyArn        string    `json:"policy_arn"`
	IsDefaultVersion bool      `json:"is_default_version"`
	CreateDate       time.Time `json:"create_date"`
	Document         string    `json:"document"`
}

// PolicyListResult represents the result of listing IAM policies.
type PolicyListResult struct {
	Policies    []*Policy
	IsTruncated bool
	Marker      string
}

// PolicyVersionListResult represents the result of listing IAM policy versions.
type PolicyVersionListResult struct {
	Versions    []*PolicyVersion
	IsTruncated bool
	Marker      string
}

// InlinePolicy represents an inline policy attached to an IAM principal.
type InlinePolicy struct {
	PrincipalType  string `json:"principal_type"`
	PrincipalName  string `json:"principal_name"`
	PolicyName     string `json:"policy_name"`
	PolicyDocument string `json:"policy_document"`
}

// AttachedPolicyRef represents a reference to an attached managed policy.
type AttachedPolicyRef struct {
	PrincipalType string `json:"principal_type"`
	PrincipalName string `json:"principal_name"`
	PolicyArn     string `json:"policy_arn"`
}

// VirtualMFADevice represents a virtual MFA device.
type VirtualMFADevice struct {
	SerialNumber     string             `json:"serial_number"`
	FriendlyName     string             `json:"friendly_name"`
	AccountId        string             `json:"account_id"`
	Base32StringSeed string             `json:"base32_string_seed,omitempty"`
	QRCodePNG        []byte             `json:"qr_code_png,omitempty"`
	EnableDate       *time.Time         `json:"enable_date,omitempty"`
	UserAssignment   *MFAUserAssignment `json:"user_assignment,omitempty"`
	Tags             []types.Tag        `json:"tags,omitempty"`
}

// MFAUserAssignment represents assignment of a user to a virtual MFA device.
type MFAUserAssignment struct {
	UserName   string    `json:"user_name"`
	EnableDate time.Time `json:"enable_date"`
}

// MFADeviceListResult represents the result of listing virtual MFA devices.
type MFADeviceListResult struct {
	MFADevices  []*VirtualMFADevice
	IsTruncated bool
	Marker      string
}

// AccountPasswordPolicy represents the account password policy.
type AccountPasswordPolicy struct {
	MinimumPasswordLength      int  `json:"minimum_password_length"`
	RequireSymbols             bool `json:"require_symbols"`
	RequireNumbers             bool `json:"require_numbers"`
	RequireUppercaseCharacters bool `json:"require_uppercase_characters"`
	RequireLowercaseCharacters bool `json:"require_lowercase_characters"`
	AllowUsersToChangePassword bool `json:"allow_users_to_change_password"`
	MaxPasswordAge             int  `json:"max_password_age"`
	PasswordReusePrevention    int  `json:"password_reuse_prevention"`
	HardExpiry                 bool `json:"hard_expiry"`
	ExpirePasswords            bool `json:"expire_passwords"`
}

// AccountAlias represents an IAM account alias.
type AccountAlias struct {
	AccountAlias string    `json:"account_alias"`
	CreateDate   time.Time `json:"create_date"`
}

// ServerCertificate represents an IAM server certificate.
type ServerCertificate struct {
	ID                    string      `json:"id"`
	Path                  string      `json:"path"`
	ServerCertificateName string      `json:"server_certificate_name"`
	Arn                   string      `json:"arn"`
	AccountId             string      `json:"account_id"`
	CreateDate            time.Time   `json:"create_date"`
	Expiration            *time.Time  `json:"expiration,omitempty"`
	CertificateBody       string      `json:"certificate_body,omitempty"`
	CertificateChain      string      `json:"certificate_chain,omitempty"`
	Tags                  []types.Tag `json:"tags,omitempty"`
}

// ServerCertificateListResult represents the result of listing server certificates.
type ServerCertificateListResult struct {
	ServerCertificateMetadataList []*ServerCertificate
	IsTruncated                   bool
	Marker                        string
}

// SigningCertificate represents an IAM signing certificate.
type SigningCertificate struct {
	CertificateId   string    `json:"certificate_id"`
	UserName        string    `json:"user_name"`
	CertificateBody string    `json:"certificate_body"`
	Status          string    `json:"status"`
	UploadDate      time.Time `json:"upload_date"`
}

// SSHPublicKey represents an IAM SSH public key.
type SSHPublicKey struct {
	SSHPublicKeyId   string    `json:"ssh_public_key_id"`
	UserName         string    `json:"user_name"`
	SSHPublicKeyBody string    `json:"ssh_public_key_body"`
	Fingerprint      string    `json:"fingerprint"`
	Status           string    `json:"status"`
	UploadDate       time.Time `json:"upload_date"`
}

// ServiceSpecificCredential represents an IAM service-specific credential.
type ServiceSpecificCredential struct {
	ServiceSpecificCredentialId   string    `json:"service_specific_credential_id"`
	ServiceSpecificCredentialName string    `json:"service_specific_credential_name"`
	ServiceName                   string    `json:"service_name"`
	UserName                      string    `json:"user_name"`
	ServicePassword               string    `json:"service_password,omitempty"`
	ServiceSpecificCredentialArn  string    `json:"service_specific_credential_arn,omitempty"`
	CreateDate                    time.Time `json:"create_date"`
	Status                        string    `json:"status"`
}

// SAMLProvider represents an IAM SAML provider.
type SAMLProvider struct {
	Arn                  string      `json:"arn"`
	AccountId            string      `json:"account_id"`
	SAMLMetadataDocument string      `json:"saml_metadata_document,omitempty"`
	ValidUntil           *time.Time  `json:"valid_until,omitempty"`
	CreateDate           time.Time   `json:"create_date"`
	Tags                 []types.Tag `json:"tags,omitempty"`
}

// SAMLProviderListResult represents the result of listing SAML providers.
type SAMLProviderListResult struct {
	SAMLProviders []*SAMLProvider
}

// OpenIDConnectProvider represents an IAM OIDC provider.
type OpenIDConnectProvider struct {
	Arn              string      `json:"arn"`
	AccountId        string      `json:"account_id"`
	URL              string      `json:"url"`
	ThumbprintList   []string    `json:"thumbprint_list"`
	ClientIDList     []string    `json:"client_id_list"`
	CreateDate       time.Time   `json:"create_date"`
	LastModifiedDate *time.Time  `json:"last_modified_date,omitempty"`
	Tags             []types.Tag `json:"tags,omitempty"`
}

// OpenIDConnectProviderListResult represents the result of listing OIDC providers.
type OpenIDConnectProviderListResult struct {
	OpenIDConnectProviderList []*OpenIDConnectProvider
}
