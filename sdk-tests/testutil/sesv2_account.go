package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2AccountTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sesv2", "GetAccount", func() error {
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if !resp.SendingEnabled {
			return fmt.Errorf("expected SendingEnabled=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountSendingAttributes", func() error {
		_, err := tc.client.PutAccountSendingAttributes(tc.ctx, &sesv2.PutAccountSendingAttributesInput{
			SendingEnabled: true,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return fmt.Errorf("get account after put: %v", err)
		}
		if !resp.SendingEnabled {
			return fmt.Errorf("expected SendingEnabled=true after put")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountSuppressionAttributes", func() error {
		reasons := []types.SuppressionListReason{types.SuppressionListReasonBounce, types.SuppressionListReasonComplaint}
		_, err := tc.client.PutAccountSuppressionAttributes(tc.ctx, &sesv2.PutAccountSuppressionAttributesInput{
			SuppressedReasons: reasons,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return fmt.Errorf("get account after suppression put: %v", err)
		}
		if resp.SuppressionAttributes == nil {
			return fmt.Errorf("SuppressionAttributes is nil")
		}
		if len(resp.SuppressionAttributes.SuppressedReasons) != 2 {
			return fmt.Errorf("expected 2 suppressed reasons, got %d", len(resp.SuppressionAttributes.SuppressedReasons))
		}
		foundBounce, foundComplaint := false, false
		for _, r := range resp.SuppressionAttributes.SuppressedReasons {
			if r == types.SuppressionListReasonBounce {
				foundBounce = true
			}
			if r == types.SuppressionListReasonComplaint {
				foundComplaint = true
			}
		}
		if !foundBounce || !foundComplaint {
			return fmt.Errorf("expected Bounce and Complaint reasons, got %v", resp.SuppressionAttributes.SuppressedReasons)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountDetails", func() error {
		_, err := tc.client.PutAccountDetails(tc.ctx, &sesv2.PutAccountDetailsInput{
			MailType:           types.MailTypeTransactional,
			UseCaseDescription: aws.String("SDK test"),
			WebsiteURL:         aws.String("https://example.com"),
			ContactLanguage:    types.ContactLanguageEn,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return fmt.Errorf("get account after PutAccountDetails: %v", err)
		}
		if resp.Details == nil {
			return fmt.Errorf("AccountDetails is nil after PutAccountDetails")
		}
		if resp.Details.MailType != types.MailTypeTransactional {
			return fmt.Errorf("expected MailType=Transactional, got %s", resp.Details.MailType)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountVdmAttributes", func() error {
		_, err := tc.client.PutAccountVdmAttributes(tc.ctx, &sesv2.PutAccountVdmAttributesInput{
			VdmAttributes: &types.VdmAttributes{
				VdmEnabled: types.FeatureStatusEnabled,
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return fmt.Errorf("get account after VDM put: %v", err)
		}
		if resp.VdmAttributes == nil {
			return fmt.Errorf("VdmAttributes is nil after PutAccountVdmAttributes")
		}
		if resp.VdmAttributes.VdmEnabled != types.FeatureStatusEnabled {
			return fmt.Errorf("expected VdmEnabled=ENABLED, got %s", resp.VdmAttributes.VdmEnabled)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountDedicatedIpWarmupAttributes", func() error {
		_, err := tc.client.PutAccountDedicatedIpWarmupAttributes(tc.ctx, &sesv2.PutAccountDedicatedIpWarmupAttributesInput{
			AutoWarmupEnabled: true,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return fmt.Errorf("get account after warmup put: %v", err)
		}
		if !resp.DedicatedIpAutoWarmupEnabled {
			return fmt.Errorf("expected DedicatedIpAutoWarmupEnabled=true after put")
		}
		return nil
	}))

	return results
}
