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
)

const mfaDeviceBucketName = "iam_mfa_devices"

// MFADeviceStore manages IAM MFA device data in persistent storage.
type MFADeviceStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
}

// NewMFADeviceStore creates a new store for IAM MFA devices.
func NewMFADeviceStore(store storage.BasicStorage, accountId string) *MFADeviceStore {
	return &MFADeviceStore{
		bucket:     store.Bucket(mfaDeviceBucketName),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an MFA device by its serial number.
func (s *MFADeviceStore) Get(serialNumber string) (*VirtualMFADevice, error) {
	data, err := s.bucket.Get([]byte(serialNumber))
	if err != nil {
		return nil, NewStoreError("get_mfa_device", err)
	}
	if data == nil {
		return nil, NewStoreError("get_mfa_device", ErrMFADeviceNotFound)
	}
	var device VirtualMFADevice
	if err := json.Unmarshal(data, &device); err != nil {
		return nil, NewStoreError("get_mfa_device", err)
	}
	return &device, nil
}

// Put stores an MFA device.
func (s *MFADeviceStore) Put(device *VirtualMFADevice) error {
	if device.Tags == nil {
		device.Tags = []Tag{}
	}

	data, err := json.Marshal(device)
	if err != nil {
		return NewStoreError("put_mfa_device", err)
	}

	if err := s.bucket.Put([]byte(device.SerialNumber), data); err != nil {
		return NewStoreError("put_mfa_device", err)
	}
	return nil
}

// Delete removes an MFA device by its serial number.
func (s *MFADeviceStore) Delete(serialNumber string) error {
	if err := s.bucket.Delete([]byte(serialNumber)); err != nil {
		return NewStoreError("delete_mfa_device", err)
	}
	return nil
}

// Exists checks whether an MFA device exists.
func (s *MFADeviceStore) Exists(serialNumber string) bool {
	return s.bucket.Has([]byte(serialNumber))
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
}

// Deactivate removes the user assignment from an MFA device.
func (s *MFADeviceStore) Deactivate(serialNumber string) error {
	device, err := s.Get(serialNumber)
	if err != nil {
		return err
	}

	device.UserAssignment = nil
	device.EnableDate = nil

	return s.Put(device)
}

// ListForUser returns all MFA devices assigned to a user.
func (s *MFADeviceStore) ListForUser(userName string, marker string, maxItems int) (*MFADeviceListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var devices []*VirtualMFADevice
	count := 0
	started := marker == ""
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var device VirtualMFADevice
		if err := json.Unmarshal(v, &device); err != nil {
			return err
		}

		if device.UserAssignment == nil || device.UserAssignment.UserName != userName {
			return nil
		}

		if !started {
			if device.SerialNumber == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			devices = append(devices, &device)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_mfa_devices", err)
	}

	result := &MFADeviceListResult{
		MFADevices:  devices,
		IsTruncated: hasMore,
	}
	if len(devices) > 0 {
		result.Marker = devices[len(devices)-1].SerialNumber
	}

	return result, nil
}

// ListVirtual returns all virtual MFA devices.
func (s *MFADeviceStore) ListVirtual(marker string, maxItems int) (*MFADeviceListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var devices []*VirtualMFADevice
	count := 0
	started := marker == ""
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var device VirtualMFADevice
		if err := json.Unmarshal(v, &device); err != nil {
			return err
		}

		if !started {
			if device.SerialNumber == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			devices = append(devices, &device)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_virtual_mfa_devices", err)
	}

	result := &MFADeviceListResult{
		MFADevices:  devices,
		IsTruncated: hasMore,
	}
	if len(devices) > 0 {
		result.Marker = devices[len(devices)-1].SerialNumber
	}

	return result, nil
}

// Resync regenerates the base32 seed for an MFA device.
func (s *MFADeviceStore) Resync(serialNumber string) error {
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
}

// Count returns the total number of MFA devices.
func (s *MFADeviceStore) Count() int {
	return s.bucket.Count()
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
	var devicesToUpdate []*VirtualMFADevice

	err := s.bucket.ForEach(func(k, v []byte) error {
		var device VirtualMFADevice
		if err := json.Unmarshal(v, &device); err != nil {
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
}

// CountForUser returns the number of MFA devices assigned to a user.
func (s *MFADeviceStore) CountForUser(userName string) (int, error) {
	count := 0
	err := s.bucket.ForEach(func(k, v []byte) error {
		var device VirtualMFADevice
		if err := json.Unmarshal(v, &device); err != nil {
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
