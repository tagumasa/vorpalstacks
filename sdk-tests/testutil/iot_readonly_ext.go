package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// runIoTReadonlyExtTests covers additional read-only List/Get/Describe ops
// that take minimal input. Each asserts the handler succeeds and returns a
// well-formed response (previously these swallowed the error with `_ = err`,
// so a 500 or unregistered handler would still "pass").
func (r *TestRunner) runIoTReadonlyExtTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	// Audit mitigation / detect reads on a non-existent task must succeed and
	// return an empty list rather than error.
	results = append(results, r.RunTest("iot", "ListAuditMitigationActionsTasks_Validation", func() error {
		_, err := tc.client.ListAuditMitigationActionsTasks(tc.ctx, &iot.ListAuditMitigationActionsTasksInput{})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "ListAuditMitigationActionsExecutions_Validation", func() error {
		_, err := tc.client.ListAuditMitigationActionsExecutions(tc.ctx, &iot.ListAuditMitigationActionsExecutionsInput{
			TaskId:       aws.String("nonexistent"),
			ActionStatus: "PENDING",
		})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "ListDetectMitigationActionsExecutions", func() error {
		_, err := tc.client.ListDetectMitigationActionsExecutions(tc.ctx, &iot.ListDetectMitigationActionsExecutionsInput{
			TaskId: aws.String("nonexistent"),
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListDetectMitigationActionsTasks_Validation", func() error {
		_, err := tc.client.ListDetectMitigationActionsTasks(tc.ctx, &iot.ListDetectMitigationActionsTasksInput{})
		return expectValidationError(err)
	}))

	// Certificate-by-CA and job-execution reads require a real id; calling them
	// without one must be rejected by validation, proving the handler validates.
	results = append(results, r.RunTest("iot", "ListCertificatesByCA_Validation", func() error {
		_, err := tc.client.ListCertificatesByCA(tc.ctx, &iot.ListCertificatesByCAInput{})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "ListJobExecutionsForJob", func() error {
		_, err := tc.client.ListJobExecutionsForJob(tc.ctx, &iot.ListJobExecutionsForJobInput{JobId: aws.String("nonexistent")})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for non-existent job")
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "ListJobExecutionsForThing", func() error {
		_, err := tc.client.ListJobExecutionsForThing(tc.ctx, &iot.ListJobExecutionsForThingInput{ThingName: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListMetricValues_Validation", func() error {
		_, err := tc.client.ListMetricValues(tc.ctx, &iot.ListMetricValuesInput{ThingName: aws.String("nonexistent")})
		return expectValidationError(err)
	}))

	// Principal/policy relationship reads on non-existent resources.
	results = append(results, r.RunTest("iot", "ListPolicyPrincipals", func() error {
		_, err := tc.client.ListPolicyPrincipals(tc.ctx, &iot.ListPolicyPrincipalsInput{PolicyName: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListPrincipalPolicies_Validation", func() error {
		resp, err := tc.client.ListPrincipalPolicies(tc.ctx, &iot.ListPrincipalPoliciesInput{Principal: aws.String("nonexistent")})
		if err != nil {
			return fmt.Errorf("ListPrincipalPolicies failed: %w", err)
		}
		if resp == nil || len(resp.Policies) != 0 {
			return fmt.Errorf("expected empty policy list for nonexistent principal")
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "ListPrincipalThings", func() error {
		_, err := tc.client.ListPrincipalThings(tc.ctx, &iot.ListPrincipalThingsInput{Principal: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListThingGroupsForThing", func() error {
		_, err := tc.client.ListThingGroupsForThing(tc.ctx, &iot.ListThingGroupsForThingInput{ThingName: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListThingPrincipals", func() error {
		_, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String("nonexistent")})
		return err
	}))

	// Audit / security-profile relationship reads.
	results = append(results, r.RunTest("iot", "ListRelatedResourcesForAuditFinding", func() error {
		_, err := tc.client.ListRelatedResourcesForAuditFinding(tc.ctx, &iot.ListRelatedResourcesForAuditFindingInput{FindingId: aws.String("nonexistent")})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for non-existent finding")
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "ListSecurityProfilesForTarget", func() error {
		_, err := tc.client.ListSecurityProfilesForTarget(tc.ctx, &iot.ListSecurityProfilesForTargetInput{
			SecurityProfileTargetArn: aws.String(tc.arn("iot", "thinggroup", "nonexistent")),
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListTargetsForSecurityProfile", func() error {
		_, err := tc.client.ListTargetsForSecurityProfile(tc.ctx, &iot.ListTargetsForSecurityProfileInput{SecurityProfileName: aws.String("nonexistent")})
		return err
	}))

	results = append(results, r.RunTest("iot", "ListTagsForResource", func() error {
		_, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: aws.String(tc.arn("iot", "thing", "nonexistent"))})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))
	results = append(results, r.RunTest("iot", "ListTargetsForPolicy", func() error {
		_, err := tc.client.ListTargetsForPolicy(tc.ctx, &iot.ListTargetsForPolicyInput{PolicyName: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListThingRegistrationTaskReports_Validation", func() error {
		_, err := tc.client.ListThingRegistrationTaskReports(tc.ctx, &iot.ListThingRegistrationTaskReportsInput{TaskId: aws.String("nonexistent")})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "GetEffectivePolicies", func() error {
		_, err := tc.client.GetEffectivePolicies(tc.ctx, &iot.GetEffectivePoliciesInput{ThingName: aws.String("nonexistent")})
		return err
	}))
	results = append(results, r.RunTest("iot", "DescribeIndex", func() error {
		_, err := tc.client.DescribeIndex(tc.ctx, &iot.DescribeIndexInput{IndexName: aws.String("AWS_Things")})
		return err
	}))

	return results
}
