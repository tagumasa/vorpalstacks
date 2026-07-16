package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// runIoTSecurityProfileTests covers the SecurityProfile lifecycle. It is
// retained here temporarily and will be expanded (Dimension/CustomMetric,
// violation reads) in the security-profile consolidation phase. All names use
// uniqueName.
func (r *TestRunner) runIoTSecurityProfileTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	profileName := uniqueName("sec-profile")

	defer tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{SecurityProfileName: aws.String(profileName)})

	results = append(results, r.RunTest("iot", "SecurityProfile_Create", func() error {
		out, err := tc.client.CreateSecurityProfile(tc.ctx, &iot.CreateSecurityProfileInput{
			SecurityProfileName:        aws.String(profileName),
			SecurityProfileDescription: aws.String("test security profile"),
		})
		if err != nil {
			return fmt.Errorf("CreateSecurityProfile failed: %w", err)
		}
		if out.SecurityProfileName == nil || *out.SecurityProfileName != profileName {
			return fmt.Errorf("expected securityProfileName=%s, got %v", profileName, out.SecurityProfileName)
		}
		if out.SecurityProfileArn == nil || *out.SecurityProfileArn == "" {
			return fmt.Errorf("expected non-empty securityProfileArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_Describe", func() error {
		out, err := tc.client.DescribeSecurityProfile(tc.ctx, &iot.DescribeSecurityProfileInput{SecurityProfileName: aws.String(profileName)})
		if err != nil {
			return fmt.Errorf("DescribeSecurityProfile failed: %w", err)
		}
		if out.SecurityProfileName == nil || *out.SecurityProfileName != profileName {
			return fmt.Errorf("expected securityProfileName=%s, got %v", profileName, out.SecurityProfileName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_List_IncludesCreated", func() error {
		out, err := tc.client.ListSecurityProfiles(tc.ctx, &iot.ListSecurityProfilesInput{})
		if err != nil {
			return fmt.Errorf("ListSecurityProfiles failed: %w", err)
		}
		for _, s := range out.SecurityProfileIdentifiers {
			if s.Name != nil && *s.Name == profileName {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d profiles", profileName, len(out.SecurityProfileIdentifiers))
	}))

	now := time.Now()
	results = append(results, r.RunTest("iot", "SecurityProfile_ListActiveViolations", func() error {
		out, err := tc.client.ListActiveViolations(tc.ctx, &iot.ListActiveViolationsInput{SecurityProfileName: aws.String(profileName)})
		if err != nil {
			return fmt.Errorf("ListActiveViolations failed: %w", err)
		}
		if out.ActiveViolations == nil {
			return fmt.Errorf("expected non-nil activeViolations list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_ListViolationEvents", func() error {
		out, err := tc.client.ListViolationEvents(tc.ctx, &iot.ListViolationEventsInput{
			SecurityProfileName: aws.String(profileName),
			StartTime:           aws.Time(now.Add(-24 * time.Hour)),
			EndTime:             aws.Time(now),
		})
		if err != nil {
			return fmt.Errorf("ListViolationEvents failed: %w", err)
		}
		if out.ViolationEvents == nil {
			return fmt.Errorf("expected non-nil violationEvents list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_GetBehaviorModelTrainingSummaries", func() error {
		out, err := tc.client.GetBehaviorModelTrainingSummaries(tc.ctx, &iot.GetBehaviorModelTrainingSummariesInput{})
		if err != nil {
			return fmt.Errorf("GetBehaviorModelTrainingSummaries failed: %w", err)
		}
		if out.Summaries == nil {
			return fmt.Errorf("expected non-nil summaries list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_Describe_NotFound", func() error {
		_, err := tc.client.DescribeSecurityProfile(tc.ctx, &iot.DescribeSecurityProfileInput{SecurityProfileName: aws.String(uniqueName("nope-profile"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_Delete", func() error {
		_, err := tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{SecurityProfileName: aws.String(profileName)})
		return err
	}))

	return results
}
