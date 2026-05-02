package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerRotationTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "RotateSecret_Basic", func() error {
		name := tc.uniqueName("RotateTest")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("rotate-me"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}
		if resp.VersionId == nil {
			return fmt.Errorf("VersionId is nil after rotation")
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch")
		}

		descResp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.LastRotatedDate == nil {
			return fmt.Errorf("LastRotatedDate should be set after rotation")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CancelRotateSecret_Basic", func() error {
		name := tc.uniqueName("CancelRot")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("cancel-rotate"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}

		cancelResp, err := tc.client.CancelRotateSecret(tc.ctx, &secretsmanager.CancelRotateSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("cancel rotate: %v", err)
		}
		if cancelResp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if cancelResp.Name == nil || *cancelResp.Name != name {
			return fmt.Errorf("name mismatch")
		}

		descResp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if descResp.RotationEnabled != nil && *descResp.RotationEnabled {
			return fmt.Errorf("rotation should be disabled after cancel")
		}
		return nil
	}))

	return results
}
