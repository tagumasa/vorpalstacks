package testutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamMFATests(tc *iamTestContext) []TestResult {
	var results []TestResult

	mfaDeviceName := fmt.Sprintf("TestMFA-%s", tc.ts)

	results = append(results, r.RunTest("iam", "CreateVirtualMFADevice", func() error {
		resp, err := tc.client.CreateVirtualMFADevice(tc.ctx, &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: aws.String(mfaDeviceName),
			Tags: []types.Tag{
				{Key: aws.String("Purpose"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.VirtualMFADevice == nil {
			return fmt.Errorf("virtual mfa device is nil")
		}
		if resp.VirtualMFADevice.SerialNumber == nil {
			return fmt.Errorf("serial number is nil")
		}
		// Base32StringSeed is a Smithy blob whose content is the UTF-8
		// bytes of the base32 seed string.  Verifying that the decoded
		// bytes form a valid base32 string (RFC 3548 alphabet) catches
		// the regression where raw 20 bytes were re-encoded instead.
		seed := resp.VirtualMFADevice.Base32StringSeed
		if len(seed) == 0 {
			return fmt.Errorf("base32 seed is empty")
		}
		for _, c := range seed {
			if !((c >= 'A' && c <= 'Z') || (c >= '2' && c <= '7')) {
				return fmt.Errorf("base32 seed contains non-alphabet byte %q (seed=%q)", c, string(seed))
			}
		}
		tc.virtualMFASerial = *resp.VirtualMFADevice.SerialNumber
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListVirtualMFADevices", func() error {
		var found bool
		var marker *string
		for {
			resp, err := tc.client.ListVirtualMFADevices(tc.ctx, &iam.ListVirtualMFADevicesInput{
				Marker: marker,
			})
			if err != nil {
				return err
			}
			for _, d := range resp.VirtualMFADevices {
				if aws.ToString(d.SerialNumber) == tc.virtualMFASerial {
					found = true
					break
				}
			}
			if found || !resp.IsTruncated || resp.Marker == nil {
				break
			}
			marker = resp.Marker
		}
		if !found {
			return fmt.Errorf("virtual mfa device %s not found", tc.virtualMFASerial)
		}
		return nil
	}))

	// MFA device tags
	results = append(results, r.RunTest("iam", "TagMFADevice", func() error {
		_, err := tc.client.TagMFADevice(tc.ctx, &iam.TagMFADeviceInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
			Tags: []types.Tag{
				{Key: aws.String("Env"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListMFADeviceTags(tc.ctx, &iam.ListMFADeviceTagsInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
		})
		if err != nil {
			return fmt.Errorf("ListMFADeviceTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env=test tag not found after TagMFADevice")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListMFADeviceTags", func() error {
		resp, err := tc.client.ListMFADeviceTags(tc.ctx, &iam.ListMFADeviceTagsInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env=test tag not found")
		}
		if !iamTagPresent(resp.Tags, "Purpose", "test") {
			return fmt.Errorf("Purpose=test tag not found (from create)")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagMFADevice", func() error {
		_, err := tc.client.UntagMFADevice(tc.ctx, &iam.UntagMFADeviceInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
			TagKeys:      []string{"Env"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListMFADeviceTags(tc.ctx, &iam.ListMFADeviceTagsInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Env", "test") {
			return fmt.Errorf("Env tag should be removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteVirtualMFADevice", func() error {
		_, err := tc.client.DeleteVirtualMFADevice(tc.ctx, &iam.DeleteVirtualMFADeviceInput{
			SerialNumber: aws.String(tc.virtualMFASerial),
		})
		if err != nil {
			return err
		}
		var marker *string
		for {
			resp, err := tc.client.ListVirtualMFADevices(tc.ctx, &iam.ListVirtualMFADevicesInput{Marker: marker})
			if err != nil {
				return fmt.Errorf("ListVirtualMFADevices after delete: %w", err)
			}
			for _, d := range resp.VirtualMFADevices {
				if aws.ToString(d.SerialNumber) == tc.virtualMFASerial {
					return fmt.Errorf("MFA device %s still present after delete", tc.virtualMFASerial)
				}
			}
			if !resp.IsTruncated || resp.Marker == nil {
				break
			}
			marker = resp.Marker
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListVirtualMFADevices_AssignmentStatus", func() error {
		// Create a fresh device so the test does not depend on the
		// lifecycle of the earlier device tests.
		deviceName := fmt.Sprintf("AssignMFA-%s", tc.ts)
		created, err := tc.client.CreateVirtualMFADevice(tc.ctx, &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: aws.String(deviceName),
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteVirtualMFADevice(tc.ctx, &iam.DeleteVirtualMFADeviceInput{
			SerialNumber: created.VirtualMFADevice.SerialNumber,
		})

		resp, err := tc.client.ListVirtualMFADevices(tc.ctx, &iam.ListVirtualMFADevicesInput{
			AssignmentStatus: types.AssignmentStatusTypeUnassigned,
			MaxItems:         aws.Int32(1000),
		})
		if err != nil {
			return err
		}
		found := false
		for _, d := range resp.VirtualMFADevices {
			if aws.ToString(d.SerialNumber) == aws.ToString(created.VirtualMFADevice.SerialNumber) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unassigned device not listed under AssignmentStatus=Unassigned")
		}

		assigned, err := tc.client.ListVirtualMFADevices(tc.ctx, &iam.ListVirtualMFADevicesInput{
			AssignmentStatus: types.AssignmentStatusTypeAssigned,
			MaxItems:         aws.Int32(1000),
		})
		if err != nil {
			return err
		}
		for _, d := range assigned.VirtualMFADevices {
			if aws.ToString(d.SerialNumber) == aws.ToString(created.VirtualMFADevice.SerialNumber) {
				return fmt.Errorf("unassigned device must not appear under AssignmentStatus=Assigned")
			}
		}

		_, err = tc.client.ListVirtualMFADevices(tc.ctx, &iam.ListVirtualMFADevicesInput{
			AssignmentStatus: "Bogus",
		})
		if err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("invalid AssignmentStatus: got %v, want InvalidInput", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "MFA_DeviceLifecycle", func() error {
		user := fmt.Sprintf("MFALife-%s", tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})

		device, err := tc.client.CreateVirtualMFADevice(tc.ctx, &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: aws.String(fmt.Sprintf("MFALife-Dev-%s", tc.ts)),
		})
		if err != nil {
			return err
		}
		serial := device.VirtualMFADevice.SerialNumber
		defer tc.client.DeleteVirtualMFADevice(tc.ctx, &iam.DeleteVirtualMFADeviceInput{SerialNumber: serial})

		// EnableMFADevice requires two consecutive codes.
		if _, err := tc.client.EnableMFADevice(tc.ctx, &iam.EnableMFADeviceInput{
			UserName:            aws.String(user),
			SerialNumber:        serial,
			AuthenticationCode1: aws.String(totpCodeAt(device.VirtualMFADevice.Base32StringSeed, 0)),
			AuthenticationCode2: aws.String(totpCodeAt(device.VirtualMFADevice.Base32StringSeed, 1)),
		}); err != nil {
			return fmt.Errorf("EnableMFADevice: %w", err)
		}

		list, err := tc.client.ListMFADevices(tc.ctx, &iam.ListMFADevicesInput{UserName: aws.String(user)})
		if err != nil {
			return fmt.Errorf("ListMFADevices: %w", err)
		}
		found := false
		for _, d := range list.MFADevices {
			if aws.ToString(d.SerialNumber) == aws.ToString(serial) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("enabled device not listed by ListMFADevices")
		}

		get, err := tc.client.GetMFADevice(tc.ctx, &iam.GetMFADeviceInput{
			UserName:     aws.String(user),
			SerialNumber: serial,
		})
		if err != nil {
			return fmt.Errorf("GetMFADevice: %w", err)
		}
		if aws.ToString(get.SerialNumber) == "" {
			return fmt.Errorf("GetMFADevice returned an empty serial number")
		}

		// ResyncMFADevice accepts a fresh pair of consecutive codes.
		if _, err := tc.client.ResyncMFADevice(tc.ctx, &iam.ResyncMFADeviceInput{
			UserName:            aws.String(user),
			SerialNumber:        serial,
			AuthenticationCode1: aws.String(totpCodeAt(device.VirtualMFADevice.Base32StringSeed, 0)),
			AuthenticationCode2: aws.String(totpCodeAt(device.VirtualMFADevice.Base32StringSeed, 1)),
		}); err != nil {
			return fmt.Errorf("ResyncMFADevice: %w", err)
		}

		if _, err := tc.client.DeactivateMFADevice(tc.ctx, &iam.DeactivateMFADeviceInput{
			UserName:     aws.String(user),
			SerialNumber: serial,
		}); err != nil {
			return fmt.Errorf("DeactivateMFADevice: %w", err)
		}

		list, err = tc.client.ListMFADevices(tc.ctx, &iam.ListMFADevicesInput{UserName: aws.String(user)})
		if err != nil {
			return err
		}
		for _, d := range list.MFADevices {
			if aws.ToString(d.SerialNumber) == aws.ToString(serial) {
				return fmt.Errorf("device still listed after DeactivateMFADevice")
			}
		}
		return nil
	}))

	return results
}

// totpCodeAt derives the RFC 6238 six-digit code for the given base32
// seed at a step offset from the current time, mirroring the server-side
// validator.
func totpCodeAt(base32Seed []byte, stepOffset int64) string {
	secret := strings.ToUpper(strings.ReplaceAll(string(base32Seed), " ", ""))
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
