package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const mfaDeviceBucketName = "iam_mfa_devices"

// MFADeviceStore manages IAM MFA device data in persistent storage.
type MFADeviceStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

// NewMFADeviceStore creates a new store for IAM MFA devices.
func NewMFADeviceStore(store storage.BasicStorage, accountId string) *MFADeviceStore {
	return &MFADeviceStore{
		BaseStore:  common.NewBaseStore(store.Bucket(mfaDeviceBucketName), "iam"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an MFA device by its serial number.
func (s *MFADeviceStore) Get(serialNumber string) (*VirtualMFADevice, error) {
	var device VirtualMFADevice
	if err := s.BaseStore.Get(serialNumber, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

// Put stores an MFA device.
func (s *MFADeviceStore) Put(device *VirtualMFADevice) error {
	if device.Tags == nil {
		device.Tags = []Tag{}
	}
	return s.BaseStore.Put(device.SerialNumber, device)
}

// Delete removes an MFA device by its serial number.
func (s *MFADeviceStore) Delete(serialNumber string) error {
	return s.BaseStore.Delete(serialNumber)
}

// Exists checks whether an MFA device exists.
func (s *MFADeviceStore) Exists(serialNumber string) bool {
	return s.BaseStore.Exists(serialNumber)
}

// Create creates a new virtual MFA device.
func (s *MFADeviceStore) Create(accountId string, tags []Tag) (*VirtualMFADevice, error) {
	serialNumber, err := GenerateMFADeviceSerialNumber()
	if err != nil {
		return nil, err
	}
	serialNumberARN := s.arnBuilder.VirtualMFADeviceARN(serialNumber)

	base32Seed, err := GenerateBase32Seed()
	if err != nil {
		return nil, err
	}

	device := &VirtualMFADevice{
		SerialNumber:     serialNumberARN,
		AccountId:        accountId,
		Base32StringSeed: base32Seed,
		Tags:             tags,
	}

	if err := s.Put(device); err != nil {
		return nil, err
	}
	return device, nil
}

// EnableForUser assigns an MFA device to a user.
func (s *MFADeviceStore) EnableForUser(serialNumber, userName string) error {
	return s.kl.WithLock(serialNumber, func() error {
		device, err := s.Get(serialNumber)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		device.UserAssignment = &MFAUserAssignment{
			UserName:   userName,
			EnableDate: now,
		}
		device.EnableDate = &now

		return s.Put(device)
	})
}

// Deactivate removes the user assignment from an MFA device.
func (s *MFADeviceStore) Deactivate(serialNumber string) error {
	return s.kl.WithLock(serialNumber, func() error {
		device, err := s.Get(serialNumber)
		if err != nil {
			return err
		}

		device.UserAssignment = nil
		device.EnableDate = nil

		return s.Put(device)
	})
}

// ListForUser returns all MFA devices assigned to a user.
func (s *MFADeviceStore) ListForUser(userName string, marker string, maxItems int) (*MFADeviceListResult, error) {
	filter := func(d *VirtualMFADevice) bool {
		return d.UserAssignment != nil && d.UserAssignment.UserName == userName
	}
	result, err := common.List[VirtualMFADevice](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &MFADeviceListResult{
		MFADevices:  result.Items,
		IsTruncated: result.IsTruncated,
		Marker:      result.NextMarker,
	}, nil
}

// ListVirtual returns all virtual MFA devices.
func (s *MFADeviceStore) ListVirtual(marker string, maxItems int) (*MFADeviceListResult, error) {
	result, err := common.List[VirtualMFADevice](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, err
	}
	return &MFADeviceListResult{
		MFADevices:  result.Items,
		IsTruncated: result.IsTruncated,
		Marker:      result.NextMarker,
	}, nil
}

// Resync regenerates the base32 seed for an MFA device.
func (s *MFADeviceStore) Resync(serialNumber string) error {
	return s.kl.WithLock(serialNumber, func() error {
		device, err := s.Get(serialNumber)
		if err != nil {
			return err
		}

		base32Seed, err := GenerateBase32Seed()
		if err != nil {
			return err
		}
		device.Base32StringSeed = base32Seed
		return s.Put(device)
	})
}

// Count returns the total number of MFA devices.
func (s *MFADeviceStore) Count() int {
	return s.BaseStore.Count()
}

// GenerateMFADeviceSerialNumber generates a unique serial number for an MFA device.
func GenerateMFADeviceSerialNumber() (string, error) {
	seed, err := GenerateBase32Seed()
	if err != nil {
		return "", err
	}
	if len(seed) < 16 {
		return "", fmt.Errorf("generated seed too short: %d < 16", len(seed))
	}
	return seed[:16], nil
}

// GenerateBase32Seed generates a random base32-encoded seed for MFA.
func GenerateBase32Seed() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("crypto/rand failed: %w", err)
	}
	return strings.ToUpper(base32Encoder.EncodeToString(bytes)), nil
}

// MigrateUser reassigns all MFA devices from one user to another.
func (s *MFADeviceStore) MigrateUser(oldUserName, newUserName string) error {
	lockKey := "migrate:" + oldUserName
	return s.kl.WithLock(lockKey, func() error {
		var devicesToUpdate []*VirtualMFADevice

		err := s.ForEach(func(key string, value []byte) error {
			var device VirtualMFADevice
			if err := json.Unmarshal(value, &device); err != nil {
				return err
			}

			if device.UserAssignment != nil && device.UserAssignment.UserName == oldUserName {
				devicesToUpdate = append(devicesToUpdate, &device)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("migrate_user_mfa", err)
		}

		for _, device := range devicesToUpdate {
			device.UserAssignment.UserName = newUserName
			if err := s.Put(device); err != nil {
				return NewStoreError("migrate_user_mfa", err)
			}
		}

		return nil
	})
}

// CountForUser returns the number of MFA devices assigned to a user.
func (s *MFADeviceStore) CountForUser(userName string) (int, error) {
	count := 0
	err := s.ForEach(func(key string, value []byte) error {
		var device VirtualMFADevice
		if err := json.Unmarshal(value, &device); err != nil {
			return err
		}

		if device.UserAssignment != nil && device.UserAssignment.UserName == userName {
			count++
		}
		return nil
	})

	if err != nil {
		return 0, NewStoreError("count_mfa_devices_for_user", err)
	}
	return count, nil
}
