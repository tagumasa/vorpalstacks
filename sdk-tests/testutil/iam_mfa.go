package testutil

import (
	"fmt"

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

	return results
}
