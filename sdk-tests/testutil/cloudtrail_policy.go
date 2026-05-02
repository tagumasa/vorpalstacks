package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailPolicyTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_GetResourcePolicy", func() error {
		name := fmt.Sprintf("policy-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("policy-bucket"),
		})
		if err != nil {
			return err
		}

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetTrail","Resource":"*"}]}`
		putResp, err := tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.TrailARN,
			ResourcePolicy: &policyDoc,
		})
		if err != nil {
			return fmt.Errorf("put resource policy: %v", err)
		}
		if putResp.ResourceArn == nil || *putResp.ResourceArn != *createResp.TrailARN {
			return fmt.Errorf("resource ARN mismatch in put response")
		}

		getResp, err := tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get resource policy: %v", err)
		}
		if getResp.ResourcePolicy == nil || *getResp.ResourcePolicy != policyDoc {
			return fmt.Errorf("policy content mismatch, got: %v", getResp.ResourcePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetResourcePolicy_NotFound", func() error {
		name := fmt.Sprintf("nopolicy-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("nopolicy-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get resource policy (no policy): %v", err)
		}
		if resp.ResourcePolicy != nil && *resp.ResourcePolicy != "" {
			return fmt.Errorf("expected empty policy, got: %s", *resp.ResourcePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteResourcePolicy", func() error {
		name := fmt.Sprintf("delpolicy-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		createResp, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("delpolicy-bucket"),
		})
		if err != nil {
			return err
		}

		policyDoc := `{"Version":"2012-10-17","Statement":[]}`
		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.TrailARN,
			ResourcePolicy: &policyDoc,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = tc.client.DeleteResourcePolicy(tc.ctx, &cloudtrail.DeleteResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("delete resource policy: %v", err)
		}

		resp, err := tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err != nil {
			return fmt.Errorf("get after delete: %v", err)
		}
		if resp.ResourcePolicy != nil && *resp.ResourcePolicy != "" {
			return fmt.Errorf("expected empty policy after delete")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_NonExistentTrail", func() error {
		fakeARN := "arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-policy-trail"
		_, err := tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    aws.String(fakeARN),
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
