package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
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

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_MaxResultsRejected", func() error {
		prefix := tc.uniqueName("BatchRange")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(prefix),
			SecretString: aws.String("range"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(prefix)

		for _, bad := range []int32{0, 21} {
			_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
				Filters:    []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{prefix}}},
				MaxResults: aws.Int32(bad),
			})
			if err == nil {
				return fmt.Errorf("MaxResults=%d should be rejected", bad)
			}
			if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
				return fmt.Errorf("MaxResults=%d: %v", bad, e)
			}
		}

		// MaxResults is documented as requiring Filters; pairing it with
		// SecretIdList is rejected rather than silently ignored.
		_, err = tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: []string{prefix},
			MaxResults:   aws.Int32(5),
		})
		if err == nil {
			return fmt.Errorf("MaxResults with SecretIdList should be rejected")
		}
		if e := expectAWSErrorCode(err, "InvalidParameterException"); e != nil {
			return fmt.Errorf("MaxResults with SecretIdList: %v", e)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_InvalidNextTokenRejected", func() error {
		prefix := tc.uniqueName("BatchBadToken")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(prefix),
			SecretString: aws.String("bad-token"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(prefix)

		for _, bad := range []string{"9999", "-1", "abc", "1x"} {
			_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
				Filters:   []types.Filter{{Key: types.FilterNameStringTypeName, Values: []string{prefix}}},
				NextToken: aws.String(bad),
			})
			if err == nil {
				return fmt.Errorf("NextToken=%q should be rejected", bad)
			}
			if e := expectAWSErrorCode(err, "InvalidNextTokenException"); e != nil {
				return fmt.Errorf("NextToken=%q: %v", bad, e)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "BatchGetSecretValue_TooManyIdsRejected", func() error {
		ids := make([]string, 21)
		for i := range ids {
			ids[i] = fmt.Sprintf("batch-id-%d", i)
		}
		_, err := tc.client.BatchGetSecretValue(tc.ctx, &secretsmanager.BatchGetSecretValueInput{
			SecretIdList: ids,
		})
		if err == nil {
			return fmt.Errorf("21 SecretIds should be rejected")
		}
		// The service model has no ValidationException shape; the batch
		// limit reports InvalidParameterException.
		return expectAWSErrorCode(err, "InvalidParameterException")
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
