// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const accessKeyBucketName = "iam_access_keys"

// AccessKeyStore manages IAM access key data in persistent storage.
type AccessKeyStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewAccessKeyStore creates a new store for IAM access keys.
func NewAccessKeyStore(store storage.BasicStorage) *AccessKeyStore {
	return &AccessKeyStore{
		BaseStore: common.NewBaseStore(store.Bucket(accessKeyBucketName), "iam"),
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

// GetBySecretKey retrieves an access key by its secret key value.
func (s *AccessKeyStore) GetBySecretKey(secretAccessKey string) (*AccessKey, error) {
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

// Put stores an access key.
func (s *AccessKeyStore) Put(key *AccessKey) error {
	if key.CreateDate.IsZero() {
		key.CreateDate = time.Now().UTC()
	}
	return s.BaseStore.Put(key.AccessKeyId, key)
}

// Delete removes an access key by its ID.
func (s *AccessKeyStore) Delete(accessKeyId string) error {
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
