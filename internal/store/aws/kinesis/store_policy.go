package kinesis

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

func (s *KinesisStore) PutResourcePolicy(resourceARN, policy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket("kinesis-policies-" + s.region)
		key := []byte(fmt.Sprintf("policy:%s", resourceARN))
		return bucket.Put(key, []byte(policy))
	})
}

func (s *KinesisStore) GetResourcePolicy(resourceARN string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result string
	err := s.storage.View(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket("kinesis-policies-" + s.region)
		key := []byte(fmt.Sprintf("policy:%s", resourceARN))
		data, err := bucket.Get(key)
		if err != nil {
			return err
		}
		if data == nil {
			return ErrResourceNotFound
		}
		result = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *KinesisStore) DeleteResourcePolicy(resourceARN string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket("kinesis-policies-" + s.region)
		key := []byte(fmt.Sprintf("policy:%s", resourceARN))
		return bucket.Delete(key)
	})
}
