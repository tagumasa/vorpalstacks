package testutil

import (
	"fmt"

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
		return iotSetupFail("FinalRem_Setup", err.Error())
	}
	defer cleanupThing()
	cleanupPolicy, err := tc.createPolicy(policyName, `{"Version":"2012-10-17","Statement":[]}`)
	if err != nil {
		return iotSetupFail("FinalRem_Setup", err.Error())
	}
	defer cleanupPolicy()

	results = append(results, r.RunTest("iot", "UpdateEncryptionConfiguration_Validation", func() error {
		_, err := tc.client.UpdateEncryptionConfiguration(tc.ctx, &iot.UpdateEncryptionConfigurationInput{})
		return expectValidationError(err)
	}))
	results = append(results, r.RunTest("iot", "Encryption_EnumAndDefault", func() error {
		// encryptionType is required and must be an EncryptionType enum
		// member (CUSTOMER_MANAGED_KMS_KEY | AWS_OWNED_KMS_KEY); the
		// account default is AWS_OWNED_KMS_KEY.
		_, err := tc.client.UpdateEncryptionConfiguration(tc.ctx, &iot.UpdateEncryptionConfigurationInput{
			EncryptionType: iottypes.EncryptionType("TLS"),
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("off-enum encryptionType: %w", err)
		}
		if _, err := tc.client.UpdateEncryptionConfiguration(tc.ctx, &iot.UpdateEncryptionConfigurationInput{
			EncryptionType: iottypes.EncryptionTypeAwsOwnedKmsKey,
		}); err != nil {
			return fmt.Errorf("valid update rejected: %w", err)
		}
		out, err := tc.client.DescribeEncryptionConfiguration(tc.ctx, &iot.DescribeEncryptionConfigurationInput{})
		if err != nil {
			return err
		}
		if out.EncryptionType != iottypes.EncryptionTypeAwsOwnedKmsKey {
			return fmt.Errorf("expected encryptionType=AWS_OWNED_KMS_KEY echoed, got %s", out.EncryptionType)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "V2LoggingLevel_Enum", func() error {
		// logLevel must be a LogLevel enum member.
		_, err := tc.client.SetV2LoggingLevel(tc.ctx, &iot.SetV2LoggingLevelInput{
			LogTarget: &iottypes.LogTarget{TargetType: iottypes.LogTargetType("DEFAULT")},
			LogLevel:  iottypes.LogLevel("VERBOSE"),
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("off-enum logLevel: %w", err)
		}
		return nil
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
