package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerSecretTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "CreateSecret", func() error {
		name := tc.uniqueName("TestSecret")
		value := "my-secret-value"
		resp, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String(value),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.Name), name)
		}
		if resp.VersionId == nil {
			return fmt.Errorf("VersionId is nil")
		}
		tc.forceDeleteSecret(name)
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret", func() error {
		name := tc.uniqueName("DescSecret")
		createResp, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("describe-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.Name), name)
		}
		if resp.ARN == nil || *resp.ARN != *createResp.ARN {
			return fmt.Errorf("ARN mismatch: got %q, want %q", aws.ToString(resp.ARN), *createResp.ARN)
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("CreatedDate is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue", func() error {
		name := tc.uniqueName("GetSecret")
		value := "my-secret-value"
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String(value),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp.SecretString == nil {
			return fmt.Errorf("SecretString is nil")
		}
		if *resp.SecretString != value {
			return fmt.Errorf("value mismatch: got %q, want %q", *resp.SecretString, value)
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch: got %q", aws.ToString(resp.Name))
		}
		if resp.VersionId == nil {
			return fmt.Errorf("VersionId is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret", func() error {
		name := tc.uniqueName("UpdateSecret")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("original"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		newValue := "updated-secret-value"
		resp, err := tc.client.UpdateSecret(tc.ctx, &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(name),
			SecretString: aws.String(newValue),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch")
		}

		getResp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get after update: %v", err)
		}
		if getResp.SecretString == nil || *getResp.SecretString != newValue {
			return fmt.Errorf("value not updated: got %q, want %q", aws.ToString(getResp.SecretString), newValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret", func() error {
		name := tc.uniqueName("DelSecret")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("delete-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		resp, err := tc.client.DeleteSecret(tc.ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(name),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if resp.Name == nil || *resp.Name != name {
			return fmt.Errorf("name mismatch")
		}
		if resp.DeletionDate == nil {
			return fmt.Errorf("DeletionDate is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets", func() error {
		name := tc.uniqueName("ListSecret")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("list-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return err
		}
		if resp.SecretList == nil {
			return fmt.Errorf("SecretList is nil")
		}
		found := false
		for _, s := range resp.SecretList {
			if s.Name != nil && *s.Name == name {
				found = true
				if s.ARN == nil {
					return fmt.Errorf("secret ARN is nil in list")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created secret not found in ListSecrets")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_WithDescription", func() error {
		name := tc.uniqueName("DescTest")
		desc := "My test secret description"
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("desc-value"),
			Description:  aws.String(desc),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.Description == nil || *resp.Description != desc {
			return fmt.Errorf("description mismatch: got %q, want %q", aws.ToString(resp.Description), desc)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_Binary", func() error {
		name := tc.uniqueName("BinaryTest")
		binaryData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretBinary: binaryData,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.SecretBinary == nil {
			return fmt.Errorf("SecretBinary is nil")
		}
		if string(resp.SecretBinary) != string(binaryData) {
			return fmt.Errorf("binary data mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_Duplicate", func() error {
		dupName := tc.uniqueName("DupSecret")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("initial-value"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.forceDeleteSecret(dupName)

		_, err = tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(dupName),
			SecretString: aws.String("duplicate-value"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate secret name")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UpdateSecret_ClearDescription", func() error {
		name := tc.uniqueName("ClearDesc")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("clear-desc"),
			Description:  aws.String("initial description"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.UpdateSecret(tc.ctx, &secretsmanager.UpdateSecretInput{
			SecretId:    aws.String(name),
			Description: aws.String(""),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		resp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if resp.Description != nil && *resp.Description != "" {
			return fmt.Errorf("description should be cleared, got %q", aws.ToString(resp.Description))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "RestoreSecret_Basic", func() error {
		name := tc.uniqueName("RestoreSec")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("restore-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}

		_, err = tc.client.DeleteSecret(tc.ctx, &secretsmanager.DeleteSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		listResp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		for _, s := range listResp.SecretList {
			if s.Name != nil && *s.Name == name {
				return fmt.Errorf("soft-deleted secret should not appear in ListSecrets")
			}
		}

		restoreResp, err := tc.client.RestoreSecret(tc.ctx, &secretsmanager.RestoreSecretInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("restore: %v", err)
		}
		if restoreResp.ARN == nil {
			return fmt.Errorf("ARN is nil after restore")
		}
		if restoreResp.Name == nil || *restoreResp.Name != name {
			return fmt.Errorf("name mismatch after restore")
		}

		getResp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get after restore: %v", err)
		}
		if getResp.SecretString == nil || *getResp.SecretString != "restore-test" {
			return fmt.Errorf("value mismatch after restore")
		}

		tc.forceDeleteSecret(name)
		return nil
	}))

	return results
}
