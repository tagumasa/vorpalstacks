package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSecretsManagerTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "secretsmanager",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := secretsmanager.NewFromConfig(cfg)
	ctx := context.Background()

	secretName := fmt.Sprintf("TestSecret-%d", time.Now().UnixNano())
	secretValue := "my-secret-value"

	results = append(results, r.RunTest("secretsmanager", "CreateSecret", func() error {
		resp, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secretName),
			SecretString: aws.String(secretValue),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("secret ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret", func() error {
		resp, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil {
			return fmt.Errorf("secret name is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue", func() error {
		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.SecretString == nil && resp.SecretBinary == nil {
			return fmt.Errorf("secret value is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets", func() error {
		resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return err
		}
		if resp.SecretList == nil {
			return fmt.Errorf("secrets list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret", func() error {
		resp, err := client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(secretName),
			SecretString: aws.String("updated-secret-value"),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("secret ARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "TagResource", func() error {
		resp, err := client.TagResource(ctx, &secretsmanager.TagResourceInput{
			SecretId: aws.String(secretName),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("sdk-tests")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecretVersionIds", func() error {
		resp, err := client.ListSecretVersionIds(ctx, &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("versions list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret", func() error {
		resp, err := client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secretName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_NonExistent", func() error {
		_, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-secret-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent secret")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret_NonExistent", func() error {
		_, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-secret-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent secret")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_Duplicate", func() error {
		dupName := fmt.Sprintf("DupSecret-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("initial-value"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(dupName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		_, err = client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("duplicate-value"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate secret name")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_ContentVerify", func() error {
		secName := fmt.Sprintf("VerifySecret-%d", time.Now().UnixNano())
		secValue := "my-verified-secret-123"
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String(secValue),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretString == nil || *resp.SecretString != secValue {
			return fmt.Errorf("secret value mismatch: got %q, want %q", aws.ToString(resp.SecretString), secValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret_ContentVerify", func() error {
		secName := fmt.Sprintf("UpdateVerify-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("original-value"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		updatedValue := "updated-secret-value-456"
		_, err = client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(secName),
			SecretString: aws.String(updatedValue),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretString == nil || *resp.SecretString != updatedValue {
			return fmt.Errorf("secret value not updated: got %q, want %q", aws.ToString(resp.SecretString), updatedValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_ContainsCreated", func() error {
		secName := fmt.Sprintf("ListVerify-%d", time.Now().UnixNano())
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("list-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(secName),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})

		resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.SecretList {
			if s.Name != nil && *s.Name == secName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created secret not found in ListSecrets")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret_NonExistent", func() error {
		_, err := client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String("arn:aws:secretsmanager:us-east-1:000000000000:secret:nonexistent-xyz"),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent secret")
		}
		return nil
	}))

	return results
}
