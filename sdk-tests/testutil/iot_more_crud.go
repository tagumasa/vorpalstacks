package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// runIoTMoreCRUDTests covers SecurityProfile attach/detach on a real thing
// group, plus cancel/delete-job negative paths. Previously these were
// `_ = err` smoke calls against hard-coded names.
func (r *TestRunner) runIoTMoreCRUDTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	profileName := uniqueName("sec-profile")
	groupName := uniqueName("more-group")
	groupARN := tc.arn("iot", "thinggroup", groupName)

	defer tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{SecurityProfileName: aws.String(profileName)})

	// Setup: create the group and profile up front; a prerequisite failure
	// surfaces as a single FAIL row named after the setup step it replaces.
	cleanupGroup, err := tc.createThingGroup(groupName)
	if err != nil {
		return iotSetupFail("More_Setup", fmt.Sprintf("CreateThingGroup failed: %v", err))
	}
	defer cleanupGroup()
	cleanupProfile, err := tc.createSecurityProfile(profileName, nil)
	if err != nil {
		return iotSetupFail("More_Setup", fmt.Sprintf("CreateSecurityProfile failed: %v", err))
	}
	defer cleanupProfile()

	results = append(results, r.RunTest("iot", "AttachSecurityProfile", func() error {
		// An unknown profile must be rejected before the association is
		// recorded.
		_, err := tc.client.AttachSecurityProfile(tc.ctx, &iot.AttachSecurityProfileInput{
			SecurityProfileName:      aws.String(uniqueName("nope-profile")),
			SecurityProfileTargetArn: aws.String(groupARN),
		})
		if err := expectNotFound(err); err != nil {
			return fmt.Errorf("unknown profile attach: %w", err)
		}
		_, err = tc.client.AttachSecurityProfile(tc.ctx, &iot.AttachSecurityProfileInput{
			SecurityProfileName:      aws.String(profileName),
			SecurityProfileTargetArn: aws.String(groupARN),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ListSecurityProfilesForTarget_Echo", func() error {
		out, err := tc.client.ListSecurityProfilesForTarget(tc.ctx, &iot.ListSecurityProfilesForTargetInput{
			SecurityProfileTargetArn: aws.String(groupARN),
		})
		if err != nil {
			return err
		}
		found := false
		for _, m := range out.SecurityProfileTargetMappings {
			if m.SecurityProfileIdentifier != nil &&
				m.SecurityProfileIdentifier.Name != nil &&
				*m.SecurityProfileIdentifier.Name == profileName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("ListSecurityProfilesForTarget did not return %s", profileName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ListTargetsForSecurityProfile_Echo", func() error {
		out, err := tc.client.ListTargetsForSecurityProfile(tc.ctx, &iot.ListTargetsForSecurityProfileInput{
			SecurityProfileName: aws.String(profileName),
		})
		if err != nil {
			return err
		}
		found := false
		for _, t := range out.SecurityProfileTargets {
			if t.Arn != nil && *t.Arn == groupARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("ListTargetsForSecurityProfile did not return %s", groupARN)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DetachSecurityProfile", func() error {
		_, err := tc.client.DetachSecurityProfile(tc.ctx, &iot.DetachSecurityProfileInput{
			SecurityProfileName:      aws.String(profileName),
			SecurityProfileTargetArn: aws.String(groupARN),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "DetachSecurityProfile_AfterDetach_NotFound", func() error {
		_, err := tc.client.DetachSecurityProfile(tc.ctx, &iot.DetachSecurityProfileInput{
			SecurityProfileName:      aws.String(profileName),
			SecurityProfileTargetArn: aws.String(groupARN),
		})
		return expectNotFound(err)
	}))

	// Job cancel/delete handlers are stubs that accept any id; assert they
	// respond without a server error.
	results = append(results, r.RunTest("iot", "CancelJob_NotFound", func() error {
		_, err := tc.client.CancelJob(tc.ctx, &iot.CancelJobInput{JobId: aws.String(uniqueName("nope-job")), ReasonCode: aws.String("test")})
		return expectNotFound(err)
	}))
	results = append(results, r.RunTest("iot", "CancelJobExecution", func() error {
		_, err := tc.client.CancelJobExecution(tc.ctx, &iot.CancelJobExecutionInput{JobId: aws.String(uniqueName("nope-job")), ThingName: aws.String("nonexistent")})
		return expectNotFound(err)
	}))
	results = append(results, r.RunTest("iot", "DeleteJobExecution", func() error {
		_, err := tc.client.DeleteJobExecution(tc.ctx, &iot.DeleteJobExecutionInput{JobId: aws.String(uniqueName("nope-job")), ThingName: aws.String("nonexistent"), ExecutionNumber: aws.Int64(1)})
		return expectNotFound(err)
	}))

	return results
}
