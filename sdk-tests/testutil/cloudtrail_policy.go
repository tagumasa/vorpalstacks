package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailPolicyTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// Trails are not in the documented ResourceArn set ("event data store,
	// dashboard, or channel"): every policy operation rejects a trail ARN
	// with ResourceTypeNotSupportedException.
	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_TrailARNRejected", func() error {
		name := tc.uniqueName("policy")
		defer tc.deleteTrail(name)

		createResp, err := tc.createTrail(name, "policy-bucket")
		if err != nil {
			return err
		}

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetTrail","Resource":"*"}]}`
		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.TrailARN,
			ResourcePolicy: &policyDoc,
		})
		if err := AssertErrorContains(err, "ResourceTypeNotSupportedException"); err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		if err := AssertErrorContains(err, "ResourceTypeNotSupportedException"); err != nil {
			return fmt.Errorf("get: %v", err)
		}

		_, _ = tc.client.DeleteResourcePolicy(tc.ctx, &cloudtrail.DeleteResourcePolicyInput{
			ResourceArn: createResp.TrailARN,
		})
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetResourcePolicy_NotFound", func() error {
		createResp, err := tc.createEventDataStore("ct-eds-nopolicy", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(aws.ToString(createResp.EventDataStoreArn))

		_, err = tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.EventDataStoreArn,
		})
		return AssertErrorContains(err, "ResourcePolicyNotFoundException")
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteResourcePolicy", func() error {
		createResp, err := tc.createEventDataStore("ct-eds-delpolicy", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(aws.ToString(createResp.EventDataStoreArn))

		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`
		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.EventDataStoreArn,
			ResourcePolicy: &policyDoc,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = tc.client.DeleteResourcePolicy(tc.ctx, &cloudtrail.DeleteResourcePolicyInput{
			ResourceArn: createResp.EventDataStoreArn,
		})
		if err != nil {
			return fmt.Errorf("delete resource policy: %v", err)
		}

		_, err = tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: createResp.EventDataStoreArn,
		})
		return AssertErrorContains(err, "ResourcePolicyNotFoundException")
	}))

	// The policy document must be a valid resource-based policy: syntax
	// errors and statements without a principal are rejected.
	results = append(results, r.RunTest("cloudtrail", "PutResourcePolicy_InvalidPolicy", func() error {
		createResp, err := tc.createEventDataStore("ct-eds-badpolicy", nil)
		if err != nil {
			return err
		}
		defer tc.deleteEventDataStore(aws.ToString(createResp.EventDataStoreArn))

		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.EventDataStoreArn,
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"`),
		})
		if err := AssertErrorContains(err, "ResourcePolicyNotValidException"); err != nil {
			return err
		}

		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudtrail.PutResourcePolicyInput{
			ResourceArn:    createResp.EventDataStoreArn,
			ResourcePolicy: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`),
		})
		return AssertErrorContains(err, "ResourcePolicyNotValidException")
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
