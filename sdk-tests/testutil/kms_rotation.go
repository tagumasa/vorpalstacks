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
		// Two on-demand rotations back to back land within the same second,
		// so a rotation-date-keyed marker would rewind onto the first
		// page's entry; the walk must deliver each rotation exactly once.
		for i := 0; i < 2; i++ {
			if _, err := tc.client.RotateKeyOnDemand(tc.ctx, &kms.RotateKeyOnDemandInput{
				KeyId: aws.String(tc.keyID),
			}); err != nil {
				return err
			}
		}
		seen := map[string]int{}
		var marker *string
		for pages := 0; ; pages++ {
			if pages > 4 {
				return fmt.Errorf("rotation pagination did not terminate")
			}
			resp, err := tc.client.ListKeyRotations(tc.ctx, &kms.ListKeyRotationsInput{
				KeyId:  aws.String(tc.keyID),
				Limit:  aws.Int32(1),
				Marker: marker,
			})
			if err != nil {
				return err
			}
			if resp.Rotations == nil {
				return fmt.Errorf("rotations is nil")
			}
			for _, rotation := range resp.Rotations {
				seen[aws.ToString(rotation.KeyMaterialId)]++
			}
			if !resp.Truncated || resp.NextMarker == nil {
				break
			}
			marker = resp.NextMarker
		}
		if len(seen) != 2 {
			return fmt.Errorf("expected exactly 2 rotations, got %d", len(seen))
		}
		for material, count := range seen {
			if count != 1 {
				return fmt.Errorf("rotation %s delivered %d times, want exactly once", material, count)
			}
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
