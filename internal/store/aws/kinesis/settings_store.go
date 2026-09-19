package kinesis

import (
	"context"

	"vorpalstacks/internal/core/storage"
)

// MTBCStatusDisabled and MTBCStatusEnabled are the
// MinimumThroughputBillingCommitmentInputStatus enum's two values — the
// persisted account setting's whole vocabulary.
const (
	MTBCStatusDisabled = "DISABLED"
	MTBCStatusEnabled  = "ENABLED"
)

// accountSettingsKey is the single key the account-level settings live
// under; the region-qualified bucket keeps the settings at the scope
// Kinesis documents them (per region, per account).
const accountSettingsKey = "account-settings"

func (s *KinesisStore) settingsBucketName() string {
	return "kinesis-account-settings-" + s.region
}

// MinimumThroughputBillingCommitmentStatus returns the persisted status
// of the account's minimum-throughput billing commitment — DISABLED until
// an UpdateAccountSettings call changes it, the state an account without
// a commitment sits in.
func (s *KinesisStore) MinimumThroughputBillingCommitmentStatus() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	status := MTBCStatusDisabled
	err := s.storage.View(context.Background(), func(txn storage.Transaction) error {
		data, err := txn.Bucket(s.settingsBucketName()).Get([]byte(accountSettingsKey))
		if err != nil {
			return err
		}
		if data != nil {
			status = string(data)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return status, nil
}

// SetMinimumThroughputBillingCommitmentStatus persists the commitment
// status; the caller has already validated it against the enum.
func (s *KinesisStore) SetMinimumThroughputBillingCommitmentStatus(status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return txn.Bucket(s.settingsBucketName()).Put([]byte(accountSettingsKey), []byte(status))
	})
}
