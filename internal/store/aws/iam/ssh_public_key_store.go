package iam

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const sshPublicKeyBucketName = "iam_ssh_public_keys"

// SSHPublicKeyStore provides storage operations for IAM SSH public keys.
type SSHPublicKeyStore struct {
	bucket storage.Bucket
	kl     common.KeyLocker
}

// NewSSHPublicKeyStore creates a new SSHPublicKeyStore instance.
func NewSSHPublicKeyStore(store storage.BasicStorage) *SSHPublicKeyStore {
	return &SSHPublicKeyStore{
		bucket: store.Bucket(sshPublicKeyBucketName),
	}
}

// Get retrieves an SSH public key by its key ID.
func (s *SSHPublicKeyStore) Get(keyId string) (*SSHPublicKey, error) {
	data, err := s.bucket.Get([]byte(keyId))
	if err != nil {
		return nil, NewStoreError("get_ssh_public_key", err)
	}
	if data == nil {
		return nil, NewStoreError("get_ssh_public_key", ErrSSHPublicKeyNotFound)
	}
	var key SSHPublicKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, NewStoreError("get_ssh_public_key", err)
	}
	return &key, nil
}

// Put stores an SSH public key, keyed by its key ID.
func (s *SSHPublicKeyStore) Put(key *SSHPublicKey) error {
	data, err := json.Marshal(key)
	if err != nil {
		return NewStoreError("put_ssh_public_key", err)
	}
	if err := s.bucket.Put([]byte(key.SSHPublicKeyId), data); err != nil {
		return NewStoreError("put_ssh_public_key", err)
	}
	return nil
}

// Delete removes an SSH public key by its key ID.
func (s *SSHPublicKeyStore) Delete(keyId string) error {
	if err := s.bucket.Delete([]byte(keyId)); err != nil {
		return NewStoreError("delete_ssh_public_key", err)
	}
	return nil
}

// Exists reports whether an SSH public key exists with the given key ID.
func (s *SSHPublicKeyStore) Exists(keyId string) bool {
	return s.bucket.Has([]byte(keyId))
}

// Upload uploads a new SSH public key for the given user.
func (s *SSHPublicKeyStore) Upload(userName, sshPublicKeyBody string) (*SSHPublicKey, error) {
	id, err := generateSSHPublicKeyID()
	if err != nil {
		return nil, err
	}
	fingerprint := computeFingerprint(sshPublicKeyBody)
	key := &SSHPublicKey{
		SSHPublicKeyId:   id,
		UserName:         userName,
		SSHPublicKeyBody: sshPublicKeyBody,
		Fingerprint:      fingerprint,
		Status:           "Active",
		UploadDate:       time.Now().UTC(),
	}
	if err := s.Put(key); err != nil {
		return nil, err
	}
	return key, nil
}

// UpdateStatus changes the status of an SSH public key (e.g. Active/Inactive).
func (s *SSHPublicKeyStore) UpdateStatus(keyId, status string) error {
	return s.kl.WithLock(keyId, func() error {
		key, err := s.Get(keyId)
		if err != nil {
			return err
		}
		key.Status = status
		return s.Put(key)
	})
}

// ListByUserName returns all SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) ListByUserName(userName string) ([]*SSHPublicKey, error) {
	var keys []*SSHPublicKey
	err := s.bucket.ForEach(func(k, v []byte) error {
		var key SSHPublicKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			keys = append(keys, &key)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_ssh_public_keys", err)
	}
	return keys, nil
}

// DeleteAllForUser removes all SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) DeleteAllForUser(userName string) error {
	var toDelete []string
	err := s.bucket.ForEach(func(k, v []byte) error {
		var key SSHPublicKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			toDelete = append(toDelete, key.SSHPublicKeyId)
		}
		return nil
	})
	if err != nil {
		return NewStoreError("delete_user_ssh_public_keys", err)
	}
	for _, id := range toDelete {
		if err := s.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of SSH public keys.
func (s *SSHPublicKeyStore) Count() int {
	return s.bucket.Count()
}

// CountByUserName returns the number of SSH public keys belonging to the given user.
func (s *SSHPublicKeyStore) CountByUserName(userName string) (int, error) {
	count := 0
	err := s.bucket.ForEach(func(k, v []byte) error {
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

func generateSSHPublicKeyID() (string, error) {
	id, err := GenerateServerCertificateID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("APKA%s", id[4:]), nil
}

func generateSigningCertificateID() (string, error) {
	id, err := GenerateServerCertificateID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CERT%s", id[4:]), nil
}
