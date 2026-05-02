package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (r *TestRunner) runKMSRotationTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "EnableKeyRotation", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.EnableKeyRotation(tc.ctx, &kms.EnableKeyRotationInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_ContentVerify", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
		if len(resp.Rotations) != 1 {
			return fmt.Errorf("expected 1 rotation, got %d", len(resp.Rotations))
		}
		if resp.Rotations[0].RotationType != types.RotationTypeAutomatic {
			return fmt.Errorf("expected RotationType=AUTOMATIC, got %s", resp.Rotations[0].RotationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "DisableKeyRotation", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		_, err := tc.client.DisableKeyRotation(tc.ctx, &kms.DisableKeyRotationInput{
			KeyId: aws.String(tc.keyID),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GetKeyRotationStatus_DisabledRotation", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
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
