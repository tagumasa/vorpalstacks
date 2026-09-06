package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMAccountTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	// GetAccountConfiguration reflects the value written by
	// PutAccountConfiguration: the seeded default and updated values
	// each survive a round trip.
	results = append(results, r.RunTest("acm", "AccountConfiguration_RoundTrip", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"default-45", func() error {
				_, _ = tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
					IdempotencyToken: aws.String(acmUniqueToken("reset")),
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
			}},
			{"roundtrip-30", func() error {
				_, err := tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
					IdempotencyToken: aws.String(acmUniqueToken("rt")),
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
			}},
			{"verify-update-60", func() error {
				_, err := tc.client.PutAccountConfiguration(tc.ctx, &acm.PutAccountConfigurationInput{
					IdempotencyToken: aws.String(acmUniqueToken("put")),
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
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	return results
}
