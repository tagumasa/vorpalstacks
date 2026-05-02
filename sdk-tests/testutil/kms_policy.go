package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r *TestRunner) runKMSPolicyTests(tc *kmsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kms", "PutKeyPolicy", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": "kms:*",
				"Resource": "*"
			}]
		}`
		_, err := tc.client.PutKeyPolicy(tc.ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(tc.keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String(policy),
		})
		return err
	}))

	results = append(results, r.RunTest("kms", "GetKeyPolicy", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.GetKeyPolicy(tc.ctx, &kms.GetKeyPolicyInput{
			KeyId:      aws.String(tc.keyID),
			PolicyName: aws.String("default"),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "GetKeyPolicy_ContentVerify", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":"kms:*","Resource":"*"}]}`
		_, err := tc.client.PutKeyPolicy(tc.ctx, &kms.PutKeyPolicyInput{
			KeyId:      aws.String(tc.keyID),
			PolicyName: aws.String("default"),
			Policy:     aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}
		resp, err := tc.client.GetKeyPolicy(tc.ctx, &kms.GetKeyPolicyInput{
			KeyId:      aws.String(tc.keyID),
			PolicyName: aws.String("default"),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if resp.Policy == nil || *resp.Policy == "" {
			return fmt.Errorf("policy is empty")
		}
		if *resp.Policy != policy {
			return fmt.Errorf("policy content mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("kms", "ListKeyPolicies", func() error {
		if tc.keyID == "" {
			return fmt.Errorf("key ID not available")
		}
		resp, err := tc.client.ListKeyPolicies(tc.ctx, &kms.ListKeyPoliciesInput{
			KeyId: aws.String(tc.keyID),
		})
		if err != nil {
			return err
		}
		if resp.PolicyNames == nil {
			return fmt.Errorf("policy names list is nil")
		}
		return nil
	}))

	return results
}
