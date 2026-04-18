package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// KeyStore provides KMS key storage and retrieval operations.
type KeyStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

func keyBucketName(region string) string {
	return "kms_keys-" + region
}

// NewKeyStore creates a new KeyStore instance with the specified storage, account ID, and region.
func NewKeyStore(store storage.BasicStorage, accountId, region string) *KeyStore {
	return &KeyStore{
		BaseStore:  common.NewBaseStore(store.Bucket(keyBucketName(region)), "kms"),
		arnBuilder: NewARNBuilder(accountId, region),
	}
}

// Create creates a new KMS key with the specified parameters.
// Returns the created key or an error.
func (s *KeyStore) Create(keyID string, keyUsage KeyUsage, keySpec KeySpec, description string, origin OriginType, multiRegion bool, tags []common.Tag) (*Key, error) {
	key := &Key{
		KeyID:              keyID,
		Arn:                s.arnBuilder.KeyArn(keyID),
		KeyState:           KeyStateEnabled,
		KeyUsage:           keyUsage,
		KeySpec:            keySpec,
		Description:        description,
		Enabled:            true,
		CreationDate:       time.Now(),
		Origin:             origin,
		KeyManager:         KeyManagerTypeCustomer,
		KeyRotationEnabled: false,
		MultiRegion:        multiRegion,
		Tags:               tags,
	}

	switch keyUsage {
	case KeyUsageEncryptDecrypt:
		switch keySpec {
		case KeySpecSymmetricDefault:
			key.EncryptionAlgorithms = []string{"SYMMETRIC_DEFAULT"}
		case KeySpecRSA2048, KeySpecRSA3072, KeySpecRSA4096:
			key.EncryptionAlgorithms = []string{"RSAES_OAEP_SHA_256", "RSAES_OAEP_SHA_1"}
		case KeySpecSM2:
			key.EncryptionAlgorithms = []string{"SM2PKE"}
		}
	case KeyUsageSignVerify:
		switch keySpec {
		case KeySpecRSA2048, KeySpecRSA3072, KeySpecRSA4096:
			key.SigningAlgorithms = []string{
				"RSASSA_PKCS1_V1_5_SHA_256",
				"RSASSA_PKCS1_V1_5_SHA_384",
				"RSASSA_PKCS1_V1_5_SHA_512",
				"RSASSA_PSS_SHA_256",
				"RSASSA_PSS_SHA_384",
				"RSASSA_PSS_SHA_512",
			}
		case KeySpecECCNISTP256, KeySpecECCSECGP256K1:
			key.SigningAlgorithms = []string{"ECDSA_SHA_256"}
		case KeySpecECCNISTP384:
			key.SigningAlgorithms = []string{"ECDSA_SHA_384"}
		case KeySpecECCNISTP521:
			key.SigningAlgorithms = []string{"ECDSA_SHA_512"}
		case KeySpecSM2:
			key.SigningAlgorithms = []string{"SM2DSA"}
		}
	case KeyUsageGenerateVerifyMAC:
		switch keySpec {
		case KeySpecHMAC224:
			key.SigningAlgorithms = []string{"HMAC_SHA_224"}
		case KeySpecHMAC256:
			key.SigningAlgorithms = []string{"HMAC_SHA_256"}
		case KeySpecHMAC384:
			key.SigningAlgorithms = []string{"HMAC_SHA_384"}
		case KeySpecHMAC512:
			key.SigningAlgorithms = []string{"HMAC_SHA_512"}
		}
	}

	if multiRegion {
		key.MultiRegionConfiguration = &MultiRegionConfiguration{
			MultiRegionKeyType: "PRIMARY",
			PrimaryKey: &PrimaryKeyInfo{
				Arn:    key.Arn,
				Region: s.arnBuilder.Region(),
			},
		}
	}

	if err := s.save(key); err != nil {
		return nil, err
	}

	return key, nil
}

// Get retrieves a KMS key by its ID.
// Returns the key's metadata or an error if not found.
func (s *KeyStore) Get(keyID string) (*Key, error) {
	keyID = s.arnBuilder.ParseKeyID(keyID)
	var key Key
	if err := s.BaseStore.Get(keyID, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (s *KeyStore) save(key *Key) error {
	return s.BaseStore.Put(key.KeyID, key)
}

// Update updates an existing KMS key's metadata.
func (s *KeyStore) Update(key *Key) error {
	return s.save(key)
}

// Delete removes a KMS key from the store.
// Returns an error if the key does not exist.
func (s *KeyStore) Delete(keyID string) error {
	keyID = s.arnBuilder.ParseKeyID(keyID)
	return s.kl.WithLock(keyID, func() error {
		return s.BaseStore.Delete(keyID)
	})
}

// Exists checks whether a KMS key exists in the store.
func (s *KeyStore) Exists(keyID string) bool {
	keyID = s.arnBuilder.ParseKeyID(keyID)
	return s.BaseStore.Exists(keyID)
}

// List returns a list of KMS keys from the store with pagination support.
func (s *KeyStore) List(marker string, maxItems int) (*KeyListResult, error) {
	result, err := common.List[Key](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, func(k *Key) bool {
		return k.KeyState != KeyStatePendingDeletion
	})
	if err != nil {
		return nil, err
	}
	items := make([]*KeyListItem, len(result.Items))
	for i, k := range result.Items {
		items[i] = &KeyListItem{
			KeyID:    k.KeyID,
			KeyArn:   k.Arn,
			Enabled:  k.Enabled,
			KeyState: k.KeyState,
		}
	}
	return &KeyListResult{
		Keys:        items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// atomicUpdate performs an unlocked Get-Modify-Update cycle on a KMS key.
// The modify function receives the key and may return an error to abort the update.
func (s *KeyStore) atomicUpdate(keyID string, modify func(*Key) error) error {
	return s.kl.WithLock(keyID, func() error {
		key, err := s.Get(keyID)
		if err != nil {
			return err
		}
		if err := modify(key); err != nil {
			return err
		}
		return s.Update(key)
	})
}

// Enable enables a KMS key.
// Returns an error if the key is pending deletion.
func (s *KeyStore) Enable(keyID string) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		if key.KeyState == KeyStatePendingDeletion {
			return ErrKeyPendingDeletion
		}
		key.Enabled = true
		key.KeyState = KeyStateEnabled
		return nil
	})
}

// Disable disables a KMS key.
func (s *KeyStore) Disable(keyID string) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		key.Enabled = false
		key.KeyState = KeyStateDisabled
		return nil
	})
}

// ScheduleDeletion schedules a KMS key for deletion after the specified pending window.
func (s *KeyStore) ScheduleDeletion(keyID string, pendingWindowInDays int) error {
	if pendingWindowInDays < 7 || pendingWindowInDays > 30 {
		return ErrInvalidPendingWindowDays
	}
	return s.atomicUpdate(keyID, func(key *Key) error {
		deletionDate := time.Now().AddDate(0, 0, pendingWindowInDays)
		key.DeletionDate = &deletionDate
		key.PendingWindowInDays = pendingWindowInDays
		origEnabled := key.Enabled
		key.PreDeletionEnabled = &origEnabled
		key.KeyState = KeyStatePendingDeletion
		key.Enabled = false
		return nil
	})
}

// CancelDeletion cancels a pending deletion for a KMS key.
func (s *KeyStore) CancelDeletion(keyID string) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		if key.KeyState != KeyStatePendingDeletion {
			return ErrInvalidKeyState
		}
		key.DeletionDate = nil
		key.PendingWindowInDays = 0
		if key.PreDeletionEnabled != nil && !*key.PreDeletionEnabled {
			key.KeyState = KeyStateDisabled
			key.Enabled = false
		} else {
			key.KeyState = KeyStateEnabled
			key.Enabled = true
		}
		return nil
	})
}

// UpdateDescription updates the description of a KMS key.
func (s *KeyStore) UpdateDescription(keyID, description string) error {
	return s.atomicUpdate(keyID, func(key *Key) error { key.Description = description; return nil })
}

// SetKeyRotation enables or disables automatic key rotation for a KMS key.
func (s *KeyStore) SetKeyRotation(keyID string, enabled bool) error {
	return s.atomicUpdate(keyID, func(key *Key) error { key.KeyRotationEnabled = enabled; return nil })
}

// AddTags adds tags to a KMS key.
func (s *KeyStore) AddTags(keyID string, tags []common.Tag) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		tagMap := make(map[string]string)
		for _, t := range key.Tags {
			tagMap[t.Key] = t.Value
		}
		for _, t := range tags {
			tagMap[t.Key] = t.Value
		}
		key.Tags = make([]common.Tag, 0, len(tagMap))
		for k, v := range tagMap {
			key.Tags = append(key.Tags, common.Tag{Key: k, Value: v})
		}
		return nil
	})
}

// RemoveTags removes tags from a KMS key.
func (s *KeyStore) RemoveTags(keyID string, tagKeys []string) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		keysToRemove := make(map[string]bool)
		for _, k := range tagKeys {
			keysToRemove[k] = true
		}
		var newTags []common.Tag
		for _, t := range key.Tags {
			if !keysToRemove[t.Key] {
				newTags = append(newTags, t)
			}
		}
		key.Tags = newTags
		return nil
	})
}

// ListTags retrieves the tags associated with a KMS key.
func (s *KeyStore) ListTags(keyID string) ([]common.Tag, error) {
	key, err := s.Get(keyID)
	if err != nil {
		return nil, err
	}
	return key.Tags, nil
}

// Count returns the number of KMS keys in the store (excluding keys pending deletion).
func (s *KeyStore) Count() (int, error) {
	count := 0
	err := s.ForEach(func(k string, v []byte) error {
		var key Key
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.KeyState != KeyStatePendingDeletion {
			count++
		}
		return nil
	})
	return count, err
}

// ARNBuilder returns the ARN builder for KMS keys.
func (s *KeyStore) ARNBuilder() *ARNBuilder {
	return s.arnBuilder
}

// ListByKeyUsage returns a list of KMS keys filtered by key usage.
func (s *KeyStore) ListByKeyUsage(keyUsage KeyUsage, marker string, maxItems int) (*KeyListResult, error) {
	result, err := common.List[Key](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, func(k *Key) bool {
		return k.KeyUsage == keyUsage && k.KeyState != KeyStatePendingDeletion
	})
	if err != nil {
		return nil, err
	}
	items := make([]*KeyListItem, len(result.Items))
	for i, k := range result.Items {
		items[i] = &KeyListItem{
			KeyID:    k.KeyID,
			KeyArn:   k.Arn,
			Enabled:  k.Enabled,
			KeyState: k.KeyState,
		}
	}
	return &KeyListResult{
		Keys:        items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// KeyMatchesSpec checks if a KMS key matches the given key spec.
func KeyMatchesSpec(key *Key, keySpec string) bool {
	if keySpec == "" {
		return true
	}
	return string(key.KeySpec) == keySpec
}

// KeyMatchesUsage checks if a KMS key matches the given key usage.
func KeyMatchesUsage(key *Key, keyUsage string) bool {
	if keyUsage == "" {
		return true
	}
	return string(key.KeyUsage) == keyUsage
}

// IsSymmetricKey checks if a KMS key is a symmetric key.
func IsSymmetricKey(keySpec KeySpec) bool {
	return keySpec == KeySpecSymmetricDefault
}

// IsAsymmetricKey checks if a KMS key is an asymmetric key.
func IsAsymmetricKey(keySpec KeySpec) bool {
	return strings.HasPrefix(string(keySpec), "RSA_") || strings.HasPrefix(string(keySpec), "ECC_")
}

// IsHMACKey checks if a KMS key is an HMAC key.
func IsHMACKey(keySpec KeySpec) bool {
	return strings.HasPrefix(string(keySpec), "HMAC_")
}

// GetParametersForImport returns the parameters required to import key material into a KMS key.
func (s *KeyStore) GetParametersForImport(keyID string, wrappingKeySpec string) (string, string, error) {
	key, err := s.Get(keyID)
	if err != nil {
		return "", "", err
	}

	if key.Origin != OriginTypeExternal {
		return "", "", ErrInvalidKeyOrigin
	}

	importToken, err := GenerateImportToken()
	if err != nil {
		return "", "", err
	}
	publicKey, err := GeneratePublicKey(wrappingKeySpec)
	if err != nil {
		return "", "", err
	}

	return importToken, publicKey, nil
}

// ImportKeyMaterial imports key material into a KMS key.
// The key must have an external origin.
func (s *KeyStore) ImportKeyMaterial(keyID string, importToken string, encryptedKeyMaterial []byte, validTo *time.Time) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		if key.Origin != OriginTypeExternal {
			return ErrInvalidKeyOrigin
		}
		key.KeyState = KeyStateEnabled
		key.Enabled = true
		key.ValidTo = validTo
		return nil
	})
}

// DeleteImportedKeyMaterial deletes imported key material from a KMS key.
func (s *KeyStore) DeleteImportedKeyMaterial(keyID string) error {
	return s.atomicUpdate(keyID, func(key *Key) error {
		if key.Origin != OriginTypeExternal {
			return ErrInvalidKeyOrigin
		}
		key.KeyState = KeyStatePendingImport
		key.Enabled = false
		key.ValidTo = nil
		return nil
	})
}

// GeneratePublicKey generates a public key for key material wrapping.
func GeneratePublicKey(wrappingKeySpec string) (string, error) {
	keyID, err := GenerateKeyID()
	if err != nil {
		return "", err
	}
	return "public-key-" + wrappingKeySpec + "-" + keyID[:16], nil
}

// ReplicateKey replicates a multi-region KMS key to another region.
func (s *KeyStore) ReplicateKey(keyID string, replicaRegion string, replicaKeyID string) (*Key, error) {
	var replicaKey *Key
	err := s.kl.WithLock(keyID, func() error {
		key, err := s.Get(keyID)
		if err != nil {
			return err
		}

		if !key.MultiRegion {
			return ErrNotMultiRegionKey
		}

		replicaKey = &Key{
			KeyID:              replicaKeyID,
			Arn:                s.arnBuilder.KeyArn(replicaKeyID),
			KeyState:           KeyStateEnabled,
			KeyUsage:           key.KeyUsage,
			KeySpec:            key.KeySpec,
			Description:        key.Description,
			Enabled:            true,
			CreationDate:       time.Now(),
			Origin:             key.Origin,
			KeyManager:         key.KeyManager,
			KeyRotationEnabled: key.KeyRotationEnabled,
			MultiRegion:        true,
			MultiRegionConfiguration: &MultiRegionConfiguration{
				MultiRegionKeyType: "REPLICA",
				PrimaryKey: &PrimaryKeyInfo{
					Arn:    key.Arn,
					Region: s.arnBuilder.Region(),
				},
				ReplicaKeys: nil,
			},
			Tags: key.Tags,
		}

		if key.MultiRegionConfiguration == nil {
			key.MultiRegionConfiguration = &MultiRegionConfiguration{
				MultiRegionKeyType: "PRIMARY",
				PrimaryKey: &PrimaryKeyInfo{
					Arn:    key.Arn,
					Region: s.arnBuilder.Region(),
				},
			}
		}

		key.MultiRegionConfiguration.ReplicaKeys = append(key.MultiRegionConfiguration.ReplicaKeys, ReplicaKeyInfo{
			Arn:    replicaKey.Arn,
			Region: replicaRegion,
		})

		if err := s.Update(key); err != nil {
			return err
		}

		if err := s.save(replicaKey); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return replicaKey, nil
}

// UpdatePrimaryRegion updates the primary region of a multi-region KMS key.
func (s *KeyStore) UpdatePrimaryRegion(keyID string, primaryRegion string) error {
	return s.kl.WithLock(keyID, func() error {
		key, err := s.Get(keyID)
		if err != nil {
			return err
		}

		if !key.MultiRegion {
			return ErrNotMultiRegionKey
		}

		if key.MultiRegionConfiguration == nil {
			return ErrNotMultiRegionKey
		}

		key.MultiRegionConfiguration.PrimaryKey = &PrimaryKeyInfo{
			Arn:    key.Arn,
			Region: primaryRegion,
		}

		return s.Update(key)
	})
}
