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
		return err
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
		return err
	}))

	return results
}
