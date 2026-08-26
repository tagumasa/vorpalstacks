package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
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

	results = append(results, r.RunTest("secretsmanager", "RotateSecret_ScheduleExpression_SetsNextRotationDate", func() error {
		name := tc.uniqueName("RotSched")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("rotate-me"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId:          aws.String(name),
			RotateImmediately: boolPtr(false),
			RotationRules: &types.RotationRulesType{
				ScheduleExpression: aws.String("rate(10 days)"),
				Duration:           aws.String("2h"),
			},
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}

		desc, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if desc.NextRotationDate == nil {
			return fmt.Errorf("NextRotationDate not set from ScheduleExpression")
		}
		// A fresh rate() schedule first fires one full period after
		// configuration.
		if desc.NextRotationDate.Before(time.Now().AddDate(0, 0, 9)) {
			return fmt.Errorf("NextRotationDate %v is not one rate period ahead", desc.NextRotationDate)
		}
		if desc.RotationRules == nil || aws.ToString(desc.RotationRules.ScheduleExpression) != "rate(10 days)" {
			return fmt.Errorf("ScheduleExpression not echoed in RotationRules")
		}
		if aws.ToString(desc.RotationRules.Duration) != "2h" {
			return fmt.Errorf("Duration not echoed in RotationRules")
		}
		if desc.LastRotatedDate != nil {
			return fmt.Errorf("LastRotatedDate set although RotateImmediately=false")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "RotateSecret_ExternalSecretRotationMembers", func() error {
		name := tc.uniqueName("RotExt")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("rotate-me"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		roleArn := "arn:aws:iam::123456789012:role/external-rotation"
		_, err = tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId:                       aws.String(name),
			RotateImmediately:              boolPtr(false),
			ExternalSecretRotationMetadata: []types.ExternalSecretRotationMetadataItem{{Key: aws.String("partner-id"), Value: aws.String("abc123")}},
			ExternalSecretRotationRoleArn:  aws.String(roleArn),
		})
		if err != nil {
			return fmt.Errorf("rotate: %v", err)
		}

		desc, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(desc.ExternalSecretRotationMetadata) != 1 {
			return fmt.Errorf("ExternalSecretRotationMetadata not echoed, got %d items", len(desc.ExternalSecretRotationMetadata))
		}
		if aws.ToString(desc.ExternalSecretRotationMetadata[0].Key) != "partner-id" ||
			aws.ToString(desc.ExternalSecretRotationMetadata[0].Value) != "abc123" {
			return fmt.Errorf("ExternalSecretRotationMetadata item mismatch: %+v", desc.ExternalSecretRotationMetadata[0])
		}
		if aws.ToString(desc.ExternalSecretRotationRoleArn) != roleArn {
			return fmt.Errorf("ExternalSecretRotationRoleArn not echoed: %q", aws.ToString(desc.ExternalSecretRotationRoleArn))
		}

		// Out-of-constraint members are rejected: the role ARN must be
		// 20-2048 characters.
		_, err = tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId:                      aws.String(name),
			RotateImmediately:             boolPtr(false),
			ExternalSecretRotationRoleArn: aws.String("too-short"),
		})
		if err == nil {
			return fmt.Errorf("short ExternalSecretRotationRoleArn should be rejected")
		}
		return expectAWSErrorCode(err, "InvalidParameterException")
	}))

	results = append(results, r.RunTest("secretsmanager", "RotateSecret_RotationRules_BothScheduleFields_Rejected", func() error {
		name := tc.uniqueName("RotBoth")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("rotate-me"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.RotateSecret(tc.ctx, &secretsmanager.RotateSecretInput{
			SecretId: aws.String(name),
			RotationRules: &types.RotationRulesType{
				AutomaticallyAfterDays: aws.Int64(30),
				ScheduleExpression:     aws.String("rate(30 days)"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error when both AutomaticallyAfterDays and ScheduleExpression are set")
		}
		if assertErr := AssertErrorContains(err, "InvalidParameterException"); assertErr != nil {
			return assertErr
		}
		return nil
	}))

	return results
}
