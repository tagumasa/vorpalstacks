package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// KeyState represents the state of a KMS key.
type KeyState string

// KeyState constants define the possible states of a KMS key. Per Smithy
// com.amazonaws.kms#KeyState the enum has 8 values: Creating, Enabled,
// Disabled, PendingDeletion, PendingImport, PendingReplicaDeletion,
// Unavailable, Updating. Creating/Updating are not currently produced by
// the platform (no asynchronous key-material generation), and
// PendingReplicaDeletion is reserved for multi-region replica teardown;
// they are declared here so that values returned by AWS-imported metadata
// are not silently re-serialised as the zero value.
const (
	KeyStateEnabled                KeyState = "Enabled"
	KeyStateDisabled               KeyState = "Disabled"
	KeyStatePendingDeletion        KeyState = "PendingDeletion"
	KeyStatePendingImport          KeyState = "PendingImport"
	KeyStateUnavailable            KeyState = "Unavailable"
	KeyStateCreating               KeyState = "Creating"
	KeyStatePendingReplicaDeletion KeyState = "PendingReplicaDeletion"
	KeyStateUpdating               KeyState = "Updating"
)

// KeyUsage represents the intended use of a KMS key.
type KeyUsage string

// KeyUsage constants define the possible usages for a KMS key.
const (
	KeyUsageEncryptDecrypt    KeyUsage = "ENCRYPT_DECRYPT"
	KeyUsageSignVerify        KeyUsage = "SIGN_VERIFY"
	KeyUsageGenerateVerifyMAC KeyUsage = "GENERATE_VERIFY_MAC"
)

// KeySpec represents the cryptographic configuration of a KMS key.
type KeySpec string

// KeySpec constants define the supported key specifications.
const (
	KeySpecSymmetricDefault KeySpec = "SYMMETRIC_DEFAULT"
	KeySpecHMAC224          KeySpec = "HMAC_224"
	KeySpecHMAC256          KeySpec = "HMAC_256"
	KeySpecHMAC384          KeySpec = "HMAC_384"
	KeySpecHMAC512          KeySpec = "HMAC_512"
	KeySpecRSA2048          KeySpec = "RSA_2048"
	KeySpecRSA3072          KeySpec = "RSA_3072"
	KeySpecRSA4096          KeySpec = "RSA_4096"
	KeySpecECCNISTP256      KeySpec = "ECC_NIST_P256"
	KeySpecECCNISTP384      KeySpec = "ECC_NIST_P384"
	KeySpecECCNISTP521      KeySpec = "ECC_NIST_P521"
	KeySpecECCSECGP256K1    KeySpec = "ECC_SECG_P256K1"
	KeySpecSM2              KeySpec = "SM2"
)

// OriginType represents the origin of key material for a KMS key.
type OriginType string

// OriginType constants define the possible origins for key material.
const (
	OriginTypeAWSKMS      OriginType = "AWS_KMS"
	OriginTypeExternal    OriginType = "EXTERNAL"
	OriginTypeAWSCloudHSM OriginType = "AWS_CLOUDHSM"
)

// KeyManagerType represents the type of key manager for a KMS key.
type KeyManagerType string

// KeyManagerType constants define the possible key manager types.
const (
	KeyManagerTypeCustomer KeyManagerType = "CUSTOMER"
	KeyManagerTypeAWS      KeyManagerType = "AWS"
)

// Key represents a KMS key.
type Key struct {
	KeyID               string         `json:"key_id"`
	Arn                 string         `json:"arn"`
	KeyState            KeyState       `json:"key_state"`
	KeyUsage            KeyUsage       `json:"key_usage"`
	KeySpec             KeySpec        `json:"key_spec"`
	Description         string         `json:"description,omitempty"`
	Enabled             bool           `json:"enabled"`
	CreationDate        time.Time      `json:"creation_date"`
	DeletionDate        *time.Time     `json:"deletion_date,omitempty"`
	PendingWindowInDays int            `json:"pending_window_in_days,omitempty"`
	ValidTo             *time.Time     `json:"valid_to,omitempty"`
	Origin              OriginType     `json:"origin"`
	KeyManager          KeyManagerType `json:"key_manager"`
	KeyRotationEnabled  bool           `json:"key_rotation_enabled"`
	// KeyRotationEnabledAt records when automatic rotation was most
	// recently enabled. Used by GetKeyRotationStatus to compute
	// NextRotationDate as max(KeyRotationEnabledAt, lastRotationDate)
	// + RotationPeriodInDays, matching the AWS documented behaviour.
	// Reset to zero when rotation is disabled so that a subsequent
	// re-enable starts a fresh window.
	KeyRotationEnabledAt           time.Time                 `json:"key_rotation_enabled_at,omitempty"`
	RotationPeriodInDays           int32                     `json:"rotation_period_in_days,omitempty"`
	OnDemandRotationStartDate      *time.Time                `json:"on_demand_rotation_start_date,omitempty"`
	CustomKeyStoreID               string                    `json:"custom_key_store_id,omitempty"`
	CloudHsmClusterID              string                    `json:"cloud_hsm_cluster_id,omitempty"`
	BypassPolicyLockoutSafetyCheck bool                      `json:"bypass_policy_lockout_safety_check"`
	MultiRegion                    bool                      `json:"multi_region"`
	MultiRegionConfiguration       *MultiRegionConfiguration `json:"multi_region_configuration,omitempty"`
	EncryptionAlgorithms           []string                  `json:"encryption_algorithms,omitempty"`
	SigningAlgorithms              []string                  `json:"signing_algorithms,omitempty"`
	MacAlgorithms                  []string                  `json:"mac_algorithms,omitempty"`
	PreDeletionEnabled             *bool                     `json:"pre_deletion_enabled,omitempty"`
	PreDeletionRotationEnabled     bool                      `json:"pre_deletion_rotation_enabled,omitempty"`
	Tags                           []types.Tag               `json:"tags,omitempty"`
	ImportToken                    string                    `json:"import_token,omitempty"`
	// ImportTokenValidTo records the expiry timestamp of the current
	// import token. Set by GetParametersForImport (24-hour window) and
	// checked by ImportKeyMaterial to reject expired tokens per the
	// AWS ExpiredImportTokenException contract.
	ImportTokenValidTo *time.Time `json:"import_token_valid_to,omitempty"`
	WrappingPrivateKey []byte     `json:"wrapping_private_key,omitempty"`
	WrappingAlgorithm  string     `json:"wrapping_algorithm,omitempty"`
	WrappingKeySpec    string     `json:"wrapping_key_spec,omitempty"`
	// ExpirationModel indicates whether imported key material expires.
	// Valid values: KEY_MATERIAL_EXPIRES (default when ValidTo is set),
	// KEY_MATERIAL_DOES_NOT_EXPIRE. Populated only for keys with Origin
	// EXTERNAL. Mirrors the Smithy ExpirationModelType enum.
	ExpirationModel string `json:"expiration_model,omitempty"`
	// RotationHistory tracks each rotation event (automatic and on-demand).
	// Used by ListKeyRotations to return past rotation details.
	RotationHistory []RotationEntry `json:"rotation_history,omitempty"`
	// OnDemandRotationCount tracks the number of on-demand rotations
	// performed. AWS enforces a limit of 25 on-demand rotations per key.
	OnDemandRotationCount int32 `json:"on_demand_rotation_count,omitempty"`
	// LastUsedAt records the timestamp of the most recent cryptographic
	// operation using this key. Updated by Encrypt, Decrypt, Sign, Verify,
	// GenerateMac, VerifyMac, GenerateDataKey, and ReEncrypt.
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	// LastUsageOperation records the KMS operation name (e.g. "Encrypt",
	// "Decrypt") of the most recent cryptographic operation. Returned in
	// the KeyLastUsage response per the Smithy KeyLastUsageData shape.
	LastUsageOperation string `json:"last_usage_operation,omitempty"`
}

// RotationEntry represents a single key rotation event.
type RotationEntry struct {
	RotationDate  time.Time `json:"rotation_date"`
	RotationType  string    `json:"rotation_type"`   // AUTOMATIC or ON_DEMAND
	KeyMaterialId string    `json:"key_material_id"` // version identifier
}

// MultiRegionConfiguration represents the configuration for a multi-region KMS key.
type MultiRegionConfiguration struct {
	MultiRegionKeyType string           `json:"multi_region_key_type"`
	PrimaryKey         *PrimaryKeyInfo  `json:"primary_key,omitempty"`
	ReplicaKeys        []ReplicaKeyInfo `json:"replica_keys,omitempty"`
}

// PrimaryKeyInfo represents information about the primary key in a multi-region key pair.
type PrimaryKeyInfo struct {
	Arn    string `json:"arn"`
	Region string `json:"region"`
}

// ReplicaKeyInfo represents information about a replica key in a multi-region key pair.
type ReplicaKeyInfo struct {
	Arn    string `json:"arn"`
	Region string `json:"region"`
}

// Alias represents a KMS key alias.
type Alias struct {
	AliasName       string    `json:"alias_name"`
	AliasArn        string    `json:"alias_arn"`
	TargetKeyID     string    `json:"target_key_id"`
	TargetKeyArn    string    `json:"target_key_arn"`
	CreationDate    time.Time `json:"creation_date"`
	LastUpdatedDate time.Time `json:"last_updated_date"`
}

// Grant represents a grant on a KMS key.
type Grant struct {
	GrantID           string            `json:"grant_id"`
	GrantToken        string            `json:"grant_token,omitempty"`
	KeyID             string            `json:"key_id"`
	GranteePrincipal  string            `json:"grantee_principal"`
	RetiringPrincipal string            `json:"retiring_principal,omitempty"`
	Operations        []string          `json:"operations"`
	Name              string            `json:"name,omitempty"`
	IssuingAccount    string            `json:"issuing_account"`
	CreationDate      time.Time         `json:"creation_date"`
	Constraints       *GrantConstraints `json:"constraints,omitempty"`
}

// GrantConstraints represents constraints on the cryptographic operations that a grant can be used for.
// Per Smithy com.amazonaws.kms#GrantConstraints, the shape has three members:
// EncryptionContextEquals, EncryptionContextSubset, and SourceArn.
type GrantConstraints struct {
	EncryptionContextEquals map[string]string `json:"encryption_context_equals,omitempty"`
	EncryptionContextSubset map[string]string `json:"encryption_context_subset,omitempty"`
	// SourceArn constrains the grant to requests made on behalf of the
	// specified resource ARN (effectively the aws:SourceArn IAM condition
	// key). Wildcards are permitted.
	SourceArn string `json:"source_arn,omitempty"`
}

// KeyPolicy represents a key policy for a KMS key.
type KeyPolicy struct {
	KeyID          string    `json:"key_id"`
	PolicyName     string    `json:"policy_name"`
	PolicyDocument string    `json:"policy_document"`
	CreateDate     time.Time `json:"create_date"`
}

// KeyListResult represents the result of listing KMS keys.
type KeyListResult struct {
	Keys        []*KeyListItem
	IsTruncated bool
	NextMarker  string
}

// KeyListItem represents a single KMS key in a list result.
type KeyListItem struct {
	KeyID    string   `json:"key_id"`
	KeyArn   string   `json:"key_arn"`
	Enabled  bool     `json:"enabled"`
	KeyState KeyState `json:"key_state"`
}

// AliasListResult represents the result of listing KMS key aliases.
type AliasListResult struct {
	Aliases     []*Alias
	IsTruncated bool
	NextMarker  string
}

// GrantListResult represents the result of listing grants for a KMS key.
type GrantListResult struct {
	Grants      []*Grant
	IsTruncated bool
	NextMarker  string
}

// ImportKeyMaterial represents the key material to be imported into a KMS key.
type ImportKeyMaterial struct {
	KeyID                string     `json:"key_id"`
	ImportToken          string     `json:"import_token"`
	EncryptedKeyMaterial []byte     `json:"encrypted_key_material"`
	ExpiryDate           *time.Time `json:"expiry_date,omitempty"`
	ValidTo              *time.Time `json:"valid_to,omitempty"`
}

// CustomKeyStore represents a custom key store in KMS.
type CustomKeyStore struct {
	CustomKeyStoreID       string    `json:"custom_key_store_id"`
	CustomKeyStoreName     string    `json:"custom_key_store_name"`
	CloudHsmClusterID      string    `json:"cloud_hsm_cluster_id"`
	TrustAnchorCertificate string    `json:"trust_anchor_certificate,omitempty"`
	CreationDate           time.Time `json:"creation_date"`
	ConnectionState        string    `json:"connection_state"`
}
