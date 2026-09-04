package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailPolicyTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_GetResourcePolicy", func() error {
		name := tc.uniqueName("policy")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "policy-bucket")
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
		name := tc.uniqueName("nopolicy")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "nopolicy-bucket")
		if err != nil {
			return err
		}

		_, err = tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err == nil {
			return fmt.Errorf("expected ResourcePolicyNotFoundException for trail with no policy")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteResourcePolicy", func() error {
		name := tc.uniqueName("delpolicy")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "delpolicy-bucket")
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
			return nil
		}
		if resp.ResourcePolicy != nil && *resp.ResourcePolicy != "" {
			return fmt.Errorf("expected empty policy after delete, got: %s", *resp.ResourcePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_NonExistentTrail", func() error {
		fakeARN := tc.trailARN("nonexistent-policy-trail")
		_, err := tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    aws.String(fakeARN),
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// The resource policy operations are documented for event data stores
	// and channels as well as trails, so the Put/Get round trip must work
	// on an event data store ARN too.
	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_GetResourcePolicy_EventDataStore", func() error {
		createResp, err := tc.createEventDataStore("ct-eds-policy", aws.Int32(90))
		if err != nil {
			return fmt.Errorf("CreateEventDataStore failed: %w", err)
		}
		if createResp.EventDataStoreArn == nil {
			return fmt.Errorf("EventDataStoreArn is nil")
		}
		defer tc.deleteEventDataStore(aws.ToString(createResp.EventDataStoreArn))

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`
		if _, err := tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.EventDataStoreArn,
			ResourcePolicy: aws.String(policyDoc),
		}); err != nil {
			return fmt.Errorf("put resource policy: %v", err)
		}

		getResp, err := tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.EventDataStoreArn,
		})
		if err != nil {
			return fmt.Errorf("get resource policy: %v", err)
		}
		if getResp.ResourcePolicy == nil || *getResp.ResourcePolicy != policyDoc {
			return fmt.Errorf("policy content mismatch, got: %v", getResp.ResourcePolicy)
		}
		return nil
	}))

	return results
}
