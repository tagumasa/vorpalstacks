package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerBatchTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_Basic", func() error {
		sec1 := tc.uniqueName("Batch1")
		sec2 := tc.uniqueName("Batch2")
		for _, name := range []string{sec1, sec2} {
			_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(name),
				SecretString: aws.String("batch-value-" + name),
			})
			if err != nil {
				return fmt.Errorf("create %s: %v", name, err)
			}
			defer tc.forceDeleteSecret(name)
		}

		resp, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: []string{sec1, sec2},
		})
		if err != nil {
			return fmt.Errorf("batch get: %v", err)
		}
		if len(resp.SecretValues) != 2 {
			return fmt.Errorf("expected 2 secret values, got %d", len(resp.SecretValues))
		}

		foundNames := make(map[string]bool)
		for _, sv := range resp.SecretValues {
			if sv.Name != nil {
				foundNames[*sv.Name] = true
			}
		}
		if !foundNames[sec1] || !foundNames[sec2] {
			return fmt.Errorf("missing secrets in batch response: found %v", foundNames)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_NonExistent", func() error {
		secName := tc.uniqueName("BatchNE")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(secName),
			SecretString: aws.String("exists"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(secName)

		resp, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: []string{secName, "nonexistent-batch-secret"},
		})
		if err != nil {
			return fmt.Errorf("batch get: %v", err)
		}
		if len(resp.SecretValues) != 1 {
			return fmt.Errorf("expected 1 secret value, got %d", len(resp.SecretValues))
		}
		if len(resp.Errors) != 1 {
			return fmt.Errorf("expected 1 error, got %d", len(resp.Errors))
		}
		return nil
	}))

	return results
}
