package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func (r *TestRunner) runSecretsManagerTagTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "CreateSecret_WithTags", func() error {
		name := tc.uniqueName("TagCreate")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("tagged"),
			Tags: []types.Tag{
				{Key: aws.String("Owner"), Value: aws.String("test-suite")},
				{Key: aws.String("Env"), Value: aws.String("dev")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		resp, err := tc.describeSecret(name)
		if err != nil {
			return err
		}
		tagMap := make(map[string]string)
		for _, t := range resp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Owner"] != "test-suite" {
			return fmt.Errorf("Owner tag mismatch: %v", tagMap)
		}
		if tagMap["Env"] != "dev" {
			return fmt.Errorf("Env tag mismatch: %v", tagMap)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "TagResource_ReservedPrefixRejected", func() error {
		name := tc.uniqueName("TagReserved")
		_, err := tc.createSecret(name, "tag-reserved-test")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.TagResource(tc.ctx, &secretsmanager.TagResourceInput{
			SecretId: aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("aws:reserved"), Value: aws.String("v")},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "TagResource_Basic", func() error {
		name := tc.uniqueName("TagTest")
		_, err := tc.createSecret(name, "tag-test")
		if err != nil {
			return err
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.TagResource(tc.ctx, &secretsmanager.TagResourceInput{
			SecretId: aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Project"), Value: aws.String("sdk-tests")},
			},
		})
		if err != nil {
			return err
		}

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return fmt.Errorf("describe after tag: %w", err)
		}
		tagMap := make(map[string]string)
		for _, t := range descResp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("Environment tag not found: %v", tagMap)
		}
		if tagMap["Project"] != "sdk-tests" {
			return fmt.Errorf("Project tag not found: %v", tagMap)
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "UntagResource_Basic", func() error {
		name := tc.uniqueName("UntagTest")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name: aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("dev")},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		_, err = tc.client.UntagResource(tc.ctx, &secretsmanager.UntagResourceInput{
			SecretId: aws.String(name),
			TagKeys:  []string{"env"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		descResp, err := tc.describeSecret(name)
		if err != nil {
			return err
		}
		for _, t := range descResp.Tags {
			if t.Key != nil && *t.Key == "env" {
				return fmt.Errorf("env tag should have been removed")
			}
		}

		tagMap := make(map[string]string)
		for _, t := range descResp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["team"] != "dev" {
			return fmt.Errorf("team tag should still exist")
		}
		return nil
	}))

	return results
}
