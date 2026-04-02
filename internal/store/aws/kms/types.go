package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// KeyState represents the state of a KMS key.
type KeyState string

// KeyState constants define the possible states of a KMS key.
const (
	KeyStateEnabled         KeyState = "Enabled"
	KeyStateDisabled        KeyState = "Disabled"
	KeyStatePendingDeletion KeyState = "PendingDeletion"
	KeyStatePendingImport   KeyState = "PendingImport"
	KeyStateUnavailable     KeyState = "Unavailable"
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
	KeyID                          string                    `json:"key_id"`
	Arn                            string                    `json:"arn"`
	KeyState                       KeyState                  `json:"key_state"`
	KeyUsage                       KeyUsage                  `json:"key_usage"`
	KeySpec                        KeySpec                   `json:"key_spec"`
	Description                    string                    `json:"description,omitempty"`
	Enabled                        bool                      `json:"enabled"`
	CreationDate                   time.Time                 `json:"creation_date"`
	DeletionDate                   *time.Time                `json:"deletion_date,omitempty"`
	PendingWindowInDays            int                       `json:"pending_window_in_days,omitempty"`
	ValidTo                        *time.Time                `json:"valid_to,omitempty"`
	Origin                         OriginType                `json:"origin"`
	KeyManager                     KeyManagerType            `json:"key_manager"`
	KeyRotationEnabled             bool                      `json:"key_rotation_enabled"`
	CustomKeyStoreID               string                    `json:"custom_key_store_id,omitempty"`
	CloudHsmClusterID              string                    `json:"cloud_hsm_cluster_id,omitempty"`
	BypassPolicyLockoutSafetyCheck bool                      `json:"bypass_policy_lockout_safety_check"`
	MultiRegion                    bool                      `json:"multi_region"`
	MultiRegionConfiguration       *MultiRegionConfiguration `json:"multi_region_configuration,omitempty"`
	EncryptionAlgorithms           []string                  `json:"encryption_algorithms,omitempty"`
	SigningAlgorithms              []string                  `json:"signing_algorithms,omitempty"`
	PreDeletionEnabled             *bool                     `json:"pre_deletion_enabled,omitempty"`
	Tags                           []common.Tag              `json:"tags,omitempty"`
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
type GrantConstraints struct {
	EncryptionContextEquals map[string]string `json:"encryption_context_equals,omitempty"`
	EncryptionContextSubset map[string]string `json:"encryption_context_subset,omitempty"`
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
