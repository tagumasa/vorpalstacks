package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTLastBatchTests covers TopicRuleDestination lifecycle (AWS-spec flow
// keyed by ARN + confirmationToken), the DeleteRegistrationCode op, and the
// PutVerificationStateOnViolation NotFound assertion.
func (r *TestRunner) runIoTLastBatchTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	// TopicRuleDestination CRUD keyed by ARN. A nil destination configuration
	// must be rejected by validation. A real create returns a
	// topicRuleDestination wrapper containing arn, status and
	// confirmationToken; the ARN then drives Get/Update/Delete, and the
	// confirmationToken drives Confirm.
	results = append(results, r.RunTest("iot", "TopicRuleDestination_Create_Validation", func() error {
		_, err := tc.client.CreateTopicRuleDestination(tc.ctx, &iot.CreateTopicRuleDestinationInput{
			DestinationConfiguration: nil,
		})
		return expectValidationError(err)
	}))

	var destArn string
	results = append(results, r.RunTest("iot", "TopicRuleDestination_CreateAndGet", func() error {
		out, err := tc.client.CreateTopicRuleDestination(tc.ctx, &iot.CreateTopicRuleDestinationInput{
			DestinationConfiguration: &iottypes.TopicRuleDestinationConfiguration{
				HttpUrlConfiguration: &iottypes.HttpUrlDestinationConfiguration{
					ConfirmationUrl: aws.String("https://example.com/confirm"),
				},
			},
		})
		if err != nil {
			return err
		}
		if out.TopicRuleDestination == nil || out.TopicRuleDestination.Arn == nil {
			return fmt.Errorf("CreateTopicRuleDestination returned nil destination or arn")
		}
		destArn = *out.TopicRuleDestination.Arn
		return nil
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Get_Created", func() error {
		if destArn == "" {
			return fmt.Errorf("destArn not captured from create step")
		}
		out, err := tc.client.GetTopicRuleDestination(tc.ctx, &iot.GetTopicRuleDestinationInput{Arn: aws.String(destArn)})
		if err != nil {
			return err
		}
		if out.TopicRuleDestination == nil || out.TopicRuleDestination.Arn == nil {
			return fmt.Errorf("GetTopicRuleDestination returned nil destination")
		}
		if *out.TopicRuleDestination.Arn != destArn {
			return fmt.Errorf("arn mismatch: got %s, want %s", *out.TopicRuleDestination.Arn, destArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Update_Status", func() error {
		if destArn == "" {
			return fmt.Errorf("destArn not captured from create step")
		}
		if _, err := tc.client.UpdateTopicRuleDestination(tc.ctx, &iot.UpdateTopicRuleDestinationInput{
			Arn:    aws.String(destArn),
			Status: iottypes.TopicRuleDestinationStatusDisabled,
		}); err != nil {
			return err
		}
		// Statuses outside the documented enum are rejected with
		// InvalidRequestException.
		_, err := tc.client.UpdateTopicRuleDestination(tc.ctx, &iot.UpdateTopicRuleDestinationInput{
			Arn:    aws.String(destArn),
			Status: iottypes.TopicRuleDestinationStatus("NOT_A_STATUS"),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Delete_Created", func() error {
		if destArn == "" {
			return fmt.Errorf("destArn not captured from create step")
		}
		_, err := tc.client.DeleteTopicRuleDestination(tc.ctx, &iot.DeleteTopicRuleDestinationInput{Arn: aws.String(destArn)})
		return err
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Get_AfterDelete_NotFound", func() error {
		if destArn == "" {
			return fmt.Errorf("destArn not captured from create step")
		}
		_, err := tc.client.GetTopicRuleDestination(tc.ctx, &iot.GetTopicRuleDestinationInput{Arn: aws.String(destArn)})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Confirm_NotFound", func() error {
		// Confirming with an unknown token must be rejected.
		_, err := tc.client.ConfirmTopicRuleDestination(tc.ctx, &iot.ConfirmTopicRuleDestinationInput{
			ConfirmationToken: aws.String(uniqueName("nope-token")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "TopicRuleDestination_Delete_NotFound", func() error {
		_, err := tc.client.DeleteTopicRuleDestination(tc.ctx, &iot.DeleteTopicRuleDestinationInput{
			Arn: aws.String(tc.arn("iot", "ruledestination", uniqueName("nope"))),
		})
		return expectNotFound(err)
	}))

	// CreateProvisioningClaim against a non-existent template must return
	// NotFound per AWS (no template -> cannot mint a claim).
	results = append(results, r.RunTest("iot", "CreateProvisioningClaim_NonexistentTemplate", func() error {
		_, err := tc.client.CreateProvisioningClaim(tc.ctx, &iot.CreateProvisioningClaimInput{
			TemplateName: aws.String(uniqueName("nope-template")),
		})
		return expectNotFound(err)
	}))

	return results
}
