package iam

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	commoniam "vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const accessKeyBucketName = "iam_access_keys"

// accessKeySecretIndexBucketName holds the secret→accessKeyId index that
// serves GetBySecretKey — the signature-verification path — without the
// full-bucket scan. Keys are SHA-256 hex digests of the secret, so the
// index bucket never duplicates secret material in the clear.
const accessKeySecretIndexBucketName = "iam_access_key_secret_index"

// accessKeySecretIndexKey derives the index key of a secret access key.
func accessKeySecretIndexKey(secretAccessKey string) string {
	sum := sha256.Sum256([]byte(secretAccessKey))
	return hex.EncodeToString(sum[:])
}

// RootUserName is the canonical constant for root user access keys,
// sourced from common/iam for cross-package consistency.
const RootUserName = commoniam.RootUserName

// AccessKeyStore manages IAM access key data in persistent storage.
type AccessKeyStore struct {
	*common.BaseStore
	kl          common.KeyLocker
	secretIndex *common.BaseStore
}

// NewAccessKeyStore creates a new store for IAM access keys.
func NewAccessKeyStore(store storage.BasicStorage) *AccessKeyStore {
	return &AccessKeyStore{
		BaseStore:   common.NewBaseStore(store.Bucket(accessKeyBucketName), "iam"),
		secretIndex: common.NewBaseStore(store.Bucket(accessKeySecretIndexBucketName), "iam"),
	}
}

// Get retrieves an access key by its ID.
func (s *AccessKeyStore) Get(accessKeyId string) (*AccessKey, error) {
	var key AccessKey
	if err := s.BaseStore.Get(accessKeyId, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

// GetBySecretKey retrieves an access key by its secret key value. The
// secret-hash index serves the hot path; an index miss or a stale entry
// falls back to the authoritative record scan, which also repairs the
// index entry — the record and index writes are not a transaction, so a
// crash or a failed index write between them leaves the pair inconsistent
// until the next lookup heals it.
func (s *AccessKeyStore) GetBySecretKey(secretAccessKey string) (*AccessKey, error) {
	indexKey := accessKeySecretIndexKey(secretAccessKey)
	if id, err := s.secretIndex.GetRaw(indexKey); err != nil {
		logs.Warn("iam: failed to read the access key secret index", logs.String("indexKey", indexKey), logs.Err(err))
	} else if len(id) > 0 {
		key, err := s.Get(string(id))
		if err == nil && key.SecretAccessKey == secretAccessKey {
			return key, nil
		}
		if err := s.secretIndex.Delete(indexKey); err != nil {
			logs.Warn("iam: failed to drop a stale access key secret index entry", logs.String("accessKeyId", string(id)), logs.Err(err))
		}
	}
	var found *AccessKey
	err := s.ForEach(func(k string, v []byte) error {
		var key AccessKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.SecretAccessKey == secretAccessKey && found == nil {
			found = &key
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_access_key_by_secret", err)
	}
	if found == nil {
		return nil, NewStoreError("get_access_key_by_secret", ErrAccessKeyNotFound)
	}
	if err := s.secretIndex.PutRaw(indexKey, []byte(found.AccessKeyId)); err != nil {
		logs.Warn("iam: failed to write the access key secret index entry", logs.String("accessKeyId", found.AccessKeyId), logs.Err(err))
	}
	return found, nil
}

// ListByUserName lists all access keys for a user, hiding the secret key.
func (s *AccessKeyStore) ListByUserName(userName string) ([]*AccessKey, error) {
	var keys []*AccessKey
	err := s.ForEach(func(k string, v []byte) error {
		var key AccessKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			keyCopy := key
			keyCopy.SecretAccessKey = ""
			keys = append(keys, &keyCopy)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_access_keys", err)
	}
	return keys, nil
}

// ListByUserNameWithSecret lists all access keys for a user, including the secret key.
func (s *AccessKeyStore) ListByUserNameWithSecret(userName string) ([]*AccessKey, error) {
	var keys []*AccessKey
	err := s.ForEach(func(k string, v []byte) error {
		var key AccessKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			keys = append(keys, &key)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_access_keys_with_secret", err)
	}
	return keys, nil
}

// Put stores an access key. The record write is authoritative; a failed
// index write only costs the fast path — GetBySecretKey's scan backstop
// repairs the entry on the next lookup.
func (s *AccessKeyStore) Put(key *AccessKey) error {
	if key.CreateDate.IsZero() {
		key.CreateDate = time.Now().UTC()
	}
	if err := s.BaseStore.Put(key.AccessKeyId, key); err != nil {
		return err
	}
	if key.SecretAccessKey != "" {
		if err := s.secretIndex.PutRaw(accessKeySecretIndexKey(key.SecretAccessKey), []byte(key.AccessKeyId)); err != nil {
			logs.Warn("iam: failed to write the access key secret index entry", logs.String("accessKeyId", key.AccessKeyId), logs.Err(err))
		}
	}
	return nil
}

// Delete removes an access key by its ID together with its secret index
// entry. A read failure before the delete only skips the index cleanup —
// a stale entry is verified and dropped on the next GetBySecretKey.
func (s *AccessKeyStore) Delete(accessKeyId string) error {
	if key, err := s.Get(accessKeyId); err == nil && key.SecretAccessKey != "" {
		if err := s.secretIndex.Delete(accessKeySecretIndexKey(key.SecretAccessKey)); err != nil {
			logs.Warn("iam: failed to delete the access key secret index entry", logs.String("accessKeyId", accessKeyId), logs.Err(err))
		}
	}
	return s.BaseStore.Delete(accessKeyId)
}

// UpdateStatus changes the status of an access key.
func (s *AccessKeyStore) UpdateStatus(accessKeyId string, status AccessKeyStatus) error {
	return s.kl.WithLock(accessKeyId, func() error {
		key, err := s.Get(accessKeyId)
		if err != nil {
			return err
		}

		if status != AccessKeyStatusActive && status != AccessKeyStatusInactive {
			return NewStoreError("update_access_key_status", ErrInvalidAccessKeyStatus)
		}

		key.Status = status
		return s.Put(key)
	})
}

// UpdateLastUsed updates the last used timestamp and location for an access key.
func (s *AccessKeyStore) UpdateLastUsed(accessKeyId, region, service string) error {
	return s.kl.WithLock(accessKeyId, func() error {
		key, err := s.Get(accessKeyId)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		key.LastUsedDate = &now
		key.LastUsedRegion = region
		key.LastUsedService = service

		return s.Put(key)
	})
}

// Create generates a new access key for a user.
func (s *AccessKeyStore) Create(userName string) (*AccessKey, error) {
	accessKeyID, err := GenerateAccessKeyID()
	if err != nil {
		return nil, NewStoreError("generate_access_key_id", err)
	}

	secretAccessKey, err := GenerateSecretAccessKey()
	if err != nil {
		return nil, NewStoreError("generate_secret_access_key", err)
	}

	key := &AccessKey{
		AccessKeyId:     accessKeyID,
		UserName:        userName,
		Status:          AccessKeyStatusActive,
		SecretAccessKey: secretAccessKey,
		CreateDate:      time.Now().UTC(),
	}

	if err := s.Put(key); err != nil {
		return nil, err
	}
	return key, nil
}

// CreateWithLimit atomically checks the per-user access-key quota and
// creates a new key inside a single lock scope, preventing the
// race condition where concurrent requests could both observe a count
// below the limit and both succeed.
func (s *AccessKeyStore) CreateWithLimit(userName string, maxKeys int) (*AccessKey, error) {
	var newKey *AccessKey
	err := s.kl.WithLock("user:"+userName, func() error {
		count, err := s.CountByUserName(userName)
		if err != nil {
			return err
		}
		if count >= maxKeys {
			return NewStoreError("create_access_key", ErrAccessKeyLimitExceeded)
		}

		accessKeyID, err := GenerateAccessKeyID()
		if err != nil {
			return NewStoreError("generate_access_key_id", err)
		}

		secretAccessKey, err := GenerateSecretAccessKey()
		if err != nil {
			return NewStoreError("generate_secret_access_key", err)
		}

		newKey = &AccessKey{
			AccessKeyId:     accessKeyID,
			UserName:        userName,
			Status:          AccessKeyStatusActive,
			SecretAccessKey: secretAccessKey,
			CreateDate:      time.Now().UTC(),
		}

		return s.Put(newKey)
	})
	if err != nil {
		return nil, err
	}
	return newKey, nil
}

// Exists checks whether an access key exists.
func (s *AccessKeyStore) Exists(accessKeyId string) bool {
	return s.BaseStore.Exists(accessKeyId)
}

// DeleteByUserName removes all access keys for a user.
func (s *AccessKeyStore) DeleteByUserName(userName string) error {
	keys, err := s.ListByUserNameWithSecret(userName)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := s.Delete(key.AccessKeyId); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of access keys.
func (s *AccessKeyStore) Count() int {
	return s.BaseStore.Count()
}

// CountByUserName returns the number of access keys for a specific user.
func (s *AccessKeyStore) CountByUserName(userName string) (int, error) {
	count := 0
	err := s.ForEach(func(k string, v []byte) error {
		var key AccessKey
		if err := json.Unmarshal(v, &key); err != nil {
			return err
		}
		if key.UserName == userName {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, NewStoreError("count_access_keys", err)
	}
	return count, nil
}
