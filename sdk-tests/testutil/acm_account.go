package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMAccountTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "GetAccountConfiguration_DefaultValues", func() error {
		_, _ = tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String(fmt.Sprintf("reset-%d", time.Now().UnixNano())),
			ExpiryEvents:     &types.ExpiryEventsConfiguration{DaysBeforeExpiry: aws.Int32(45)},
		})
		resp, err := tc.client.GetAccountConfiguration(tc.ctx, &acm.GetAccountConfigurationInput{})
		if err != nil {
			return err
		}
		if resp.ExpiryEvents == nil {
			return fmt.Errorf("ExpiryEvents is nil")
		}
		if resp.ExpiryEvents.DaysBeforeExpiry == nil {
			return fmt.Errorf("DaysBeforeExpiry is nil")
		}
		if aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry) != 45 {
			return fmt.Errorf("expected default 45, got %d", aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "GetAccountConfiguration_RoundTrip", func() error {
		_, err := tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String(fmt.Sprintf("rt-%d", time.Now().UnixNano())),
			ExpiryEvents: &types.ExpiryEventsConfiguration{
				DaysBeforeExpiry: aws.Int32(30),
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccountConfiguration(tc.ctx, &acm.GetAccountConfigurationInput{})
		if err != nil {
			return err
		}
		if resp.ExpiryEvents == nil {
			return fmt.Errorf("ExpiryEvents is nil")
		}
		if aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry) != 30 {
			return fmt.Errorf("expected 30, got %d", aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "PutAccountConfiguration_VerifyUpdate", func() error {
		token := fmt.Sprintf("put-%d", time.Now().UnixNano())
		_, err := tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
			IdempotencyToken: aws.String(token),
			ExpiryEvents: &types.ExpiryEventsConfiguration{
				DaysBeforeExpiry: aws.Int32(60),
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetAccountConfiguration(tc.ctx, &acm.GetAccountConfigurationInput{})
		if err != nil {
			return err
		}
		if aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry) != 60 {
			return fmt.Errorf("expected 60 after PutAccountConfiguration, got %d", aws.ToInt32(resp.ExpiryEvents.DaysBeforeExpiry))
		}
		return nil
	}))

	return results
}
