package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r *TestRunner) runKMSRotationTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "EnableKeyRotation", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.EnableKeyRotation(tc.ctx, &kms.EnableKeyRotationInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_ContentVerify", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		resp, err := tc.client.GetKeyRotationStatus(tc.ctx, &kms.GetKeyRotationStatusInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if !resp.KeyRotationEnabled {
			return fmt.Errorf("expected KeyRotationEnabled=true")
		}
		if resp.RotationPeriodInDays == nil || *resp.RotationPeriodInDays != 365 {
			return fmt.Errorf("expected RotationPeriodInDays=365, got %d", aws.ToInt32(resp.RotationPeriodInDays))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeyRotations", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		resp, err := tc.client.ListKeyRotations(tc.ctx, &kms.ListKeyRotationsInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.Rotations == nil {
			return fmt.Errorf("rotations is nil")
		}
		// No actual rotation has been performed on the test key, so the
		// rotation history should be empty. AWS returns an empty list for
		// keys that have never been rotated.
		if len(resp.Rotations) != 0 {
			return fmt.Errorf("expected 0 rotations, got %d", len(resp.Rotations))
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DisableKeyRotation", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		_, err := tc.client.DisableKeyRotation(tc.ctx, &kms.DisableKeyRotationInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_DisabledRotation", func() error {
		if err := tc.requireKeyID(); err != nil {
			return err
		}
		resp, err := tc.client.GetKeyRotationStatus(tc.ctx, &kms.GetKeyRotationStatusInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.KeyRotationEnabled {
			return fmt.Errorf("expected KeyRotationEnabled=false")
		}
		return nil
	}))

	return results
}
