package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func pgTimestamp() int64 {
	return time.Now().UnixNano()
}

func (r *TestRunner) runSecretsManagerEdgeTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "GetSecretValue_NonExistent", func() error {
		_, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String("nonexistent-secret-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DescribeSecret_NonExistent", func() error {
		_, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
			SecretId: aws.String("nonexistent-secret-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteSecret_NonExistent", func() error {
		_, err := tc.client.DeleteSecret(tc.ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String("nonexistent-delete-xyz"),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "RestoreSecret_NonExistent", func() error {
		_, err := tc.client.RestoreSecret(tc.ctx, &secretsmanager.RestoreSecretInput{
			SecretId: aws.String("nonexistent-restore-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Filters", func() error {
		prefix := tc.uniqueName("FilterTest")
		for _, suffix := range []string{"alpha", "beta"} {
			_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(prefix + "-" + suffix),
				SecretString: aws.String(suffix),
			})
			if err != nil {
				return fmt.Errorf("create %s: %v", suffix, err)
			}
			defer tc.forceDeleteSecret(prefix + "-" + suffix)
		}

		resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{Key: types.FilterNameStringTypeName, Values: []string{prefix + "-alpha"}},
			},
		})
		if err != nil {
			return fmt.Errorf("list with filter: %v", err)
		}
		if len(resp.SecretList) != 1 {
			return fmt.Errorf("expected 1 secret, got %d", len(resp.SecretList))
		}
		if resp.SecretList[0].Name == nil || *resp.SecretList[0].Name != prefix+"-alpha" {
			return fmt.Errorf("wrong secret returned: %q", aws.ToString(resp.SecretList[0].Name))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ListSecrets_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", pgTimestamp())
		var pgSecrets []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagSecret-%s-%d", pgTs, i)
			_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(name),
				SecretString: aws.String("pagval"),
			})
			if err != nil {
				for _, sn := range pgSecrets {
					tc.forceDeleteSecret(sn)
				}
				return fmt.Errorf("create secret %s: %v", name, err)
			}
			pgSecrets = append(pgSecrets, name)
		}

		var allSecrets []string
		var nextToken *string
		for {
			resp, err := tc.client.ListSecrets(tc.ctx, &secretsmanager.ListSecretsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, sn := range pgSecrets {
					tc.forceDeleteSecret(sn)
				}
				return fmt.Errorf("list secrets page: %v", err)
			}
			for _, s := range resp.SecretList {
				if s.Name != nil && strings.Contains(*s.Name, "PagSecret-"+pgTs) {
					allSecrets = append(allSecrets, *s.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, sn := range pgSecrets {
			tc.forceDeleteSecret(sn)
		}
		if len(allSecrets) != 5 {
			return fmt.Errorf("expected 5 paginated secrets, got %d", len(allSecrets))
		}
		return nil
	}))

	return results
}
