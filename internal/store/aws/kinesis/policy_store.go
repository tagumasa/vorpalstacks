package kinesis

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// policyKey is the single key-construction site for resource policies. The
// ARN is normalised so every spelling of the same resource addresses one
// entry, and the delete-time sweep in DeleteStream removes exactly what the
// policy operations wrote.
func (s *KinesisStore) policyKey(resourceARN string) string {
	return fmt.Sprintf("policy:%s", s.normalizeARN(resourceARN))
}

// deleteResourcePolicyLocked removes the resource-policy entry of one ARN;
// the caller must hold s.mu.
func (s *KinesisStore) deleteResourcePolicyLocked(resourceARN string) error {
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket(s.policiesBucketName())
		return bucket.Delete([]byte(s.policyKey(resourceARN)))
	})
}

// PutResourcePolicy persists a resource policy for a Kinesis stream.
func (s *KinesisStore) PutResourcePolicy(resourceARN, policy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket(s.policiesBucketName())
		key := []byte(s.policyKey(resourceARN))
		return bucket.Put(key, []byte(policy))
	})
}

// GetResourcePolicy retrieves the resource policy for a Kinesis stream.
func (s *KinesisStore) GetResourcePolicy(resourceARN string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result string
	err := s.storage.View(context.Background(), func(txn storage.Transaction) error {
		bucket := txn.Bucket(s.policiesBucketName())
		key := []byte(s.policyKey(resourceARN))
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

// DeleteResourcePolicy removes the resource policy for a Kinesis stream.
func (s *KinesisStore) DeleteResourcePolicy(resourceARN string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.deleteResourcePolicyLocked(resourceARN)
}
