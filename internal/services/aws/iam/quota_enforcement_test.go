package iam

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// quotaTestStore builds a store backed by a fresh temp directory.
func quotaTestStore(t *testing.T) *iamstore.IAMStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return iamstore.NewIAMStore(st, "123456789012")
}

// totpCodeAt computes the expected six-digit code for the seed at the
// current time step plus the offset, mirroring the TOTP derivation the
// Core validates against.
func totpCodeAt(base32Seed string, stepOffset int64) string {
	secret := strings.ToUpper(strings.ReplaceAll(base32Seed, " ", ""))
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		panic(fmt.Sprintf("base32 decode failed: %v", err))
	}
	step := time.Now().Unix()/30 + stepOffset
	counter := make([]byte, 8)
	binary.BigEndian.PutUint64(counter, uint64(step))
	h := hmac.New(sha1.New, decoded)
	h.Write(counter)
	sum := h.Sum(nil)
	offset := sum[len(sum)-1] & 0x0F
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7FFFFFFF
	return fmt.Sprintf("%06d", code%1000000)
}

// The per-user group quota is enforced on the add path: memberships one
// to MaxIAMGroupsPerUser succeed, the next distinct group is rejected with
// the documented LimitExceeded fault, and an already-held membership stays
// an idempotent success at the quota.
func TestAddUserToGroupCoreQuota(t *testing.T) {
	store := quotaTestStore(t)
	s := NewIAMService("123456789012")

	if _, err := store.Users().Create("quota-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}
	for i := 0; i < iamstore.MaxIAMGroupsPerUser; i++ {
		groupName := fmt.Sprintf("quota-group-%02d", i)
		if _, err := store.Groups().Create(groupName, "/", "123456789012"); err != nil {
			t.Fatalf("create group %s: %v", groupName, err)
		}
		if err := s.addUserToGroupCore(store, &UserGroupMembershipInput{UserName: "quota-user", GroupName: groupName}); err != nil {
			t.Fatalf("add membership %d: %v", i+1, err)
		}
	}

	store.Groups().Create("quota-group-overflow", "/", "123456789012")
	err := s.addUserToGroupCore(store, &UserGroupMembershipInput{UserName: "quota-user", GroupName: "quota-group-overflow"})
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("quota overflow: expected *awserrors.AWSError, got %T (%v)", err, err)
	}
	if awsErr.GetCode() != "LimitExceeded" {
		t.Fatalf("quota overflow code: got %s, want LimitExceeded", awsErr.GetCode())
	}
	if !strings.Contains(awsErr.Error(), fmt.Sprintf("GroupsPerUser: %d", iamstore.MaxIAMGroupsPerUser)) {
		t.Fatalf("quota overflow message: got %s", awsErr.Error())
	}
	if awsErr.GetHTTPStatusCode() != http.StatusConflict {
		t.Fatalf("quota overflow status: got %d, want 409", awsErr.GetHTTPStatusCode())
	}

	if err := s.addUserToGroupCore(store, &UserGroupMembershipInput{UserName: "quota-user", GroupName: "quota-group-00"}); err != nil {
		t.Fatalf("already-held membership at quota: got %v, want idempotent success", err)
	}
}

// The per-user MFA-device quota is enforced on the enable path: enables
// one to MaxMFADevicesPerUser succeed, the next device is rejected with
// the documented LimitExceeded fault, and deactivating one device frees
// quota for a new enable.
func TestEnableMFADeviceCoreQuota(t *testing.T) {
	store := quotaTestStore(t)
	s := NewIAMService("123456789012")

	if _, err := store.Users().Create("mfa-quota-user", "/", "123456789012", nil); err != nil {
		t.Fatalf("create user: %v", err)
	}

	var firstSerial string
	for i := 0; i < iamstore.MaxMFADevicesPerUser; i++ {
		device, err := store.MFADevices().Create("123456789012", fmt.Sprintf("quota-mfa-%02d", i), nil)
		if err != nil {
			t.Fatalf("create device %d: %v", i, err)
		}
		if i == 0 {
			firstSerial = device.SerialNumber
		}
		if err := s.enableMFADeviceCore(store, &EnableMFADeviceInput{
			UserName:            "mfa-quota-user",
			SerialNumber:        device.SerialNumber,
			AuthenticationCode1: totpCodeAt(device.Base32StringSeed, 0),
			AuthenticationCode2: totpCodeAt(device.Base32StringSeed, 1),
		}); err != nil {
			t.Fatalf("enable device %d: %v", i+1, err)
		}
	}

	device, err := store.MFADevices().Create("123456789012", "quota-mfa-overflow", nil)
	if err != nil {
		t.Fatalf("create overflow device: %v", err)
	}
	err = s.enableMFADeviceCore(store, &EnableMFADeviceInput{
		UserName:            "mfa-quota-user",
		SerialNumber:        device.SerialNumber,
		AuthenticationCode1: totpCodeAt(device.Base32StringSeed, 0),
		AuthenticationCode2: totpCodeAt(device.Base32StringSeed, 1),
	})
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("quota overflow: expected *awserrors.AWSError, got %T (%v)", err, err)
	}
	if awsErr.GetCode() != "LimitExceeded" {
		t.Fatalf("quota overflow code: got %s, want LimitExceeded", awsErr.GetCode())
	}
	if !strings.Contains(awsErr.Error(), fmt.Sprintf("MFADevicesPerUser: %d", iamstore.MaxMFADevicesPerUser)) {
		t.Fatalf("quota overflow message: got %s", awsErr.Error())
	}
	if awsErr.GetHTTPStatusCode() != http.StatusConflict {
		t.Fatalf("quota overflow status: got %d, want 409", awsErr.GetHTTPStatusCode())
	}

	if err := s.deactivateMFADeviceCore(store, "mfa-quota-user", firstSerial); err != nil {
		t.Fatalf("deactivate first device: %v", err)
	}
	if err := s.enableMFADeviceCore(store, &EnableMFADeviceInput{
		UserName:            "mfa-quota-user",
		SerialNumber:        device.SerialNumber,
		AuthenticationCode1: totpCodeAt(device.Base32StringSeed, 0),
		AuthenticationCode2: totpCodeAt(device.Base32StringSeed, 1),
	}); err != nil {
		t.Fatalf("enable after deactivating one device: %v", err)
	}
}
