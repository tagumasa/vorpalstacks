package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2SuppressionTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	suppressedEmail := fmt.Sprintf("suppressed-%d@example.com", tc.uid)

	results = append(results, r.RunTest("sesv2", "PutSuppressedDestination", func() error {
		_, err := tc.client.PutSuppressedDestination(tc.ctx, &sesv2.PutSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
			Reason:       types.SuppressionListReasonBounce,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetSuppressedDestination", func() error {
		resp, err := tc.client.GetSuppressedDestination(tc.ctx, &sesv2.GetSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
		})
		if err != nil {
			return err
		}
		if resp.SuppressedDestination == nil {
			return fmt.Errorf("suppressed destination is nil")
		}
		if resp.SuppressedDestination.EmailAddress == nil || *resp.SuppressedDestination.EmailAddress != suppressedEmail {
			return fmt.Errorf("expected email %s, got %v", suppressedEmail, resp.SuppressedDestination.EmailAddress)
		}
		if resp.SuppressedDestination.Reason != types.SuppressionListReasonBounce {
			return fmt.Errorf("expected Reason=Bounce, got %s", resp.SuppressedDestination.Reason)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListSuppressedDestinations", func() error {
		all, err := tc.listAllSuppressedDestinations()
		if err != nil {
			return err
		}
		if !containsSuppressedEmail(all, suppressedEmail) {
			return fmt.Errorf("suppressed email %s not found in list", suppressedEmail)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteSuppressedDestination", func() error {
		_, err := tc.client.DeleteSuppressedDestination(tc.ctx, &sesv2.DeleteSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetSuppressedDestination(tc.ctx, &sesv2.GetSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting suppressed destination")
		}
		return nil
	}))

	return results
}
