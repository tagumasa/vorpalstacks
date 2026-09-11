package iam

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const sshPublicKeyBucketName = "iam_ssh_public_keys"

// SSHPublicKeyStore provides storage operations for IAM SSH public keys.
type SSHPublicKeyStore struct {
	uk userKeyed[SSHPublicKey]
}

// NewSSHPublicKeyStore creates a new SSHPublicKeyStore instance.
func NewSSHPublicKeyStore(store storage.BasicStorage) *SSHPublicKeyStore {
	return &SSHPublicKeyStore{
		uk: newUserKeyed[SSHPublicKey](
			common.NewBaseStore(store.Bucket(sshPublicKeyBucketName), "iam"),
			func(k *SSHPublicKey) string { return k.SSHPublicKeyId },
			func(k *SSHPublicKey) string { return k.UserName },
		),
	}
}

// Get retrieves an SSH public key by its key ID.
func (s *SSHPublicKeyStore) Get(keyId string) (*SSHPublicKey, error) {
	return getByKey[SSHPublicKey](s.uk.BaseStore, keyId, "get_ssh_public_key", ErrSSHPublicKeyNotFound)
}

// Put stores an SSH public key, keyed by its key ID.
func (s *SSHPublicKeyStore) Put(key *SSHPublicKey) error {
	return s.uk.BaseStore.Put(key.SSHPublicKeyId, key)
}

// Delete removes an SSH public key by its key ID.
func (s *SSHPublicKeyStore) Delete(keyId string) error {
	return s.uk.BaseStore.Delete(keyId)
}

// Exists reports whether an SSH public key exists with the given key ID.
func (s *SSHPublicKeyStore) Exists(keyId string) bool {
	return s.uk.BaseStore.Exists(keyId)
}

// MaxSSHPublicKeysPerUser is the AWS-enforced quota of SSH public keys per
// IAM user.
const MaxSSHPublicKeysPerUser = 5

// UploadWithGuards stores a canonicalised SSH public key after checking,
// inside a single lock scope, that the same key material is not already
// registered for the user and that the per-user quota is not exceeded.
// The fingerprint is computed over the key blob, so differing comments do
// not hide a duplicate.
func (s *SSHPublicKeyStore) UploadWithGuards(userName, sshPublicKeyBody string) (*SSHPublicKey, error) {
	var created *SSHPublicKey
	err := s.uk.kl.WithLock("ssh-key:"+userName, func() error {
		existing, err := s.ListByUserName(userName)
		if err != nil {
			return err
		}
		fingerprint := computeFingerprint(sshPublicKeyBody)
		for _, key := range existing {
			if key.Fingerprint == fingerprint {
				return NewStoreError("upload_ssh_public_key", ErrDuplicateSSHPublicKey)
			}
		}
		if len(existing) >= MaxSSHPublicKeysPerUser {
			return NewStoreError("upload_ssh_public_key", ErrSSHPublicKeyLimitExceeded)
		}

		id, err := GenerateSSHPublicKeyID()
		if err != nil {
			return NewStoreError("generate_ssh_public_key_id", err)
		}
		created = &SSHPublicKey{
			SSHPublicKeyId:   id,
			UserName:         userName,
			SSHPublicKeyBody: sshPublicKeyBody,
			Fingerprint:      fingerprint,
			Status:           "Active",
			UploadDate:       time.Now().UTC(),
		}
		return s.Put(created)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// UpdateStatus changes the status of an SSH public key (e.g. Active/Inactive).
func (s *SSHPublicKeyStore) UpdateStatus(keyId, status string) error {
	return s.uk.updateStatus(keyId, status, s.Get, func(k *SSHPublicKey, status string) { k.Status = status })
}

// ListByUserName returns all SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) ListByUserName(userName string) ([]*SSHPublicKey, error) {
	return s.uk.listByUserName(userName, "list_ssh_public_keys")
}

// DeleteAllForUser removes all SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) DeleteAllForUser(userName string) error {
	return s.uk.deleteAllForUser(userName, "delete_user_ssh_public_keys")
}

// MigrateUser updates the UserName field on all SSH public keys from
// oldUserName to newUserName. Called during IAM user rename operations.
func (s *SSHPublicKeyStore) MigrateUser(oldUserName, newUserName string) error {
	return s.uk.migrateUser(oldUserName, newUserName, "migrate_ssh_public_keys",
		func(k *SSHPublicKey, newName string) { k.UserName = newName })
}

// Count returns the total number of SSH public keys.
func (s *SSHPublicKeyStore) Count() int {
	return s.uk.BaseStore.Count()
}

// CountByUserName returns the number of SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) CountByUserName(userName string) (int, error) {
	count := 0
	err := s.uk.ForEach(func(k string, v []byte) error {
		var key SSHPublicKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, NewStoreError("count_ssh_public_keys", err)
	}
	return count, nil
}

func computeFingerprint(publicKeyBody string) string {
	parts := strings.Fields(publicKeyBody)
	var keyData []byte
	if len(parts) >= 2 {
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err == nil {
			keyData = decoded
		}
	}
	if keyData == nil {
		keyData = []byte(publicKeyBody)
	}
	h := sha256.Sum256(keyData)
	return base64.StdEncoding.EncodeToString(h[:])
}
