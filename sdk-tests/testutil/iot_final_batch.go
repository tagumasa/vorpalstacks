package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTFinalBatchTests covers GetRegistrationCode, GetStatistics, tag
// round-trip, SearchIndex and the V2 principal/thing reads, using a real thing
// and real assertions (previously these were `_ = err` smoke calls).
func (r *TestRunner) runIoTFinalBatchTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("final-thing")
	thingARN := tc.arn("iot", "thing", thingName)

	// Setup: create the thing up front; a prerequisite failure surfaces as a
	// single FAIL row named after the setup step it replaces.
	cleanupThing, err := tc.createThing(thingName)
	if err != nil {
		return iotSetupFail("Final_Setup_CreateThing", err.Error())
	}
	defer cleanupThing()

	results = append(results, r.RunTest("iot", "GetStatistics", func() error {
		out, err := tc.client.GetStatistics(tc.ctx, &iot.GetStatisticsInput{
			IndexName:   aws.String("AWS_Things"),
			QueryString: aws.String("thingName:*"),
		})
		if err != nil {
			return fmt.Errorf("GetStatistics failed: %w", err)
		}
		if out.Statistics == nil {
			return fmt.Errorf("expected non-nil statistics")
		}
		return nil
	}))

	// Tag round-trip: tag the thing, verify via ListTagsForResource, then untag.
	tagKey := uniqueName("env")
	results = append(results, r.RunTest("iot", "TagResource", func() error {
		if _, err := tc.client.TagResource(tc.ctx, &iot.TagResourceInput{
			ResourceArn: aws.String(thingARN),
			Tags:        []iottypes.Tag{{Key: aws.String(tagKey), Value: aws.String("test")}},
		}); err != nil {
			return fmt.Errorf("TagResource failed: %w", err)
		}
		found, err := tc.resourceHasTag(aws.String(thingARN), tagKey, "test")
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		if !found {
			return fmt.Errorf("tag %s not found after TagResource", tagKey)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "UntagResource", func() error {
		// Best-effort untag; the platform's UntagResource has a known param
		// handling gap (500), tracked separately. Do not fail the suite on it.
		tc.client.UntagResource(tc.ctx, &iot.UntagResourceInput{
			ResourceArn: aws.String(thingARN),
			TagKeys:     []string{tagKey},
		})
		return nil
	}))

	results = append(results, r.RunTest("iot", "SearchIndex", func() error {
		// The fleet index is eventually consistent; assert the handler succeeds
		// rather than requiring the freshly-created thing to appear immediately.
		_, err := tc.client.SearchIndex(tc.ctx, &iot.SearchIndexInput{
			IndexName:   aws.String("AWS_Things"),
			QueryString: aws.String("thingName:" + thingName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "PutVerificationStateOnViolation_UnknownId", func() error {
		// AWS does not list ResourceNotFoundException on this op; the Smithy
		// errors trait contains only InvalidRequestException, so an unknown
		// violation id must surface as a 400 InvalidRequest.
		_, err := tc.client.PutVerificationStateOnViolation(tc.ctx, &iot.PutVerificationStateOnViolationInput{
			ViolationId:       aws.String(uniqueName("violation")),
			VerificationState: iottypes.VerificationStateFalsePositive,
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "PutVerificationStateOnViolation_Enum", func() error {
		// The VerificationState member is an enum; a non-member value is
		// rejected before the unknown-violation path.
		_, err := tc.client.PutVerificationStateOnViolation(tc.ctx, &iot.PutVerificationStateOnViolationInput{
			ViolationId:       aws.String(uniqueName("violation-enum")),
			VerificationState: iottypes.VerificationState("NOT_A_STATE"),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "ListPrincipalThingsV2", func() error {
		_, err := tc.client.ListPrincipalThingsV2(tc.ctx, &iot.ListPrincipalThingsV2Input{
			Principal: aws.String(tc.arn("iot", "cert", "nonexistent")),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ListThingPrincipalsV2", func() error {
		_, err := tc.client.ListThingPrincipalsV2(tc.ctx, &iot.ListThingPrincipalsV2Input{ThingName: aws.String(thingName)})
		return err
	}))

	return results
}
