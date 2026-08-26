package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTFinalRemainingTests covers account-level configuration updates
// (encryption/event config) and TestAuthorization, plus policy attach/detach
// on a real thing. Each call asserts success (previously `_ = err` swallowing).
func (r *TestRunner) runIoTFinalRemainingTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("finalrem-thing")
	policyName := uniqueName("finalrem-policy")
	thingARN := tc.arn("iot", "thing", thingName)

	// Setup: create the thing and policy through the shared cleanup-returning
	// helpers; a prerequisite failure surfaces as a single FAIL row named
	// after the setup step it replaces.
	cleanupThing, err := tc.createThing(thingName)
	if err != nil {
		return []TestResult{{Service: "iot", TestName: "FinalRem_Setup", Status: "FAIL", Error: err.Error()}}
	}
	defer cleanupThing()
	cleanupPolicy, err := tc.createPolicy(policyName, `{"Version":"2012-10-17","Statement":[]}`)
	if err != nil {
		return []TestResult{{Service: "iot", TestName: "FinalRem_Setup", Status: "FAIL", Error: err.Error()}}
	}
	defer cleanupPolicy()

	results = append(results, r.RunTest("iot", "UpdateEncryptionConfiguration_Validation", func() error {
		_, err := tc.client.UpdateEncryptionConfiguration(tc.ctx, &iot.UpdateEncryptionConfigurationInput{})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "UpdateEventConfigurations", func() error {
		_, err := tc.client.UpdateEventConfigurations(tc.ctx, &iot.UpdateEventConfigurationsInput{})
		return err
	}))
	results = append(results, r.RunTest("iot", "UpdateTopicRuleDestination_Validation", func() error {
		_, err := tc.client.UpdateTopicRuleDestination(tc.ctx, &iot.UpdateTopicRuleDestinationInput{})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "AttachPolicy_ToThing", func() error {
		_, err := tc.client.AttachPolicy(tc.ctx, &iot.AttachPolicyInput{PolicyName: aws.String(policyName), Target: aws.String(thingARN)})
		return err
	}))
	results = append(results, r.RunTest("iot", "DetachPolicy_FromThing", func() error {
		_, err := tc.client.DetachPolicy(tc.ctx, &iot.DetachPolicyInput{PolicyName: aws.String(policyName), Target: aws.String(thingARN)})
		return err
	}))

	results = append(results, r.RunTest("iot", "TestAuthorization", func() error {
		_, err := tc.client.TestAuthorization(tc.ctx, &iot.TestAuthorizationInput{
			AuthInfos: []iottypes.AuthInfo{{Resources: []string{"*"}, ActionType: iottypes.ActionTypeConnect}},
		})
		return err
	}))

	return results
}
