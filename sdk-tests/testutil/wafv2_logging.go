package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func (r *TestRunner) runWAFv2LoggingTests(tc *wafv2TestContext) []TestResult {
	var results []TestResult

	logACLName := tc.uniqueName("log-webacl")
	var logACLID, logACLARN, logACLLockToken string

	results = append(results, r.RunTest("wafv2", "Logging_Setup", func() error {
		resp, err := tc.client.CreateWebACL(tc.ctx, &wafv2.CreateWebACLInput{
			Name:          aws.String(logACLName),
			Scope:         tc.scope,
			DefaultAction: &types.DefaultAction{Allow: &types.AllowAction{}},
			VisibilityConfig: &types.VisibilityConfig{
				SampledRequestsEnabled: true, CloudWatchMetricsEnabled: true,
				MetricName: aws.String("log-webacl-metric"),
			},
		})
		if err != nil {
			return err
		}
		logACLID = aws.ToString(resp.Summary.Id)
		logACLARN = aws.ToString(resp.Summary.ARN)
		logACLLockToken = aws.ToString(resp.Summary.LockToken)
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "PutLoggingConfiguration", func() error {
		resp, err := tc.client.PutLoggingConfiguration(tc.ctx, &wafv2.PutLoggingConfigurationInput{
			LoggingConfiguration: &types.LoggingConfiguration{
				ResourceArn: aws.String(logACLARN),
				LogDestinationConfigs: []string{
					"arn:aws:logs:us-east-1:123456789012:log-group:/aws/waf/test-log",
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.LoggingConfiguration == nil {
			return fmt.Errorf("LoggingConfiguration is nil")
		}
		if len(resp.LoggingConfiguration.LogDestinationConfigs) != 1 {
			return fmt.Errorf("expected 1 log destination, got %d", len(resp.LoggingConfiguration.LogDestinationConfigs))
		}
		if aws.ToString(resp.LoggingConfiguration.ResourceArn) != logACLARN {
			return fmt.Errorf("ResourceArn mismatch in PutLoggingConfiguration response")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetLoggingConfiguration", func() error {
		resp, err := tc.client.GetLoggingConfiguration(tc.ctx, &wafv2.GetLoggingConfigurationInput{
			ResourceArn: aws.String(logACLARN),
		})
		if err != nil {
			return err
		}
		if resp.LoggingConfiguration == nil {
			return fmt.Errorf("LoggingConfiguration is nil")
		}
		if aws.ToString(resp.LoggingConfiguration.ResourceArn) != logACLARN {
			return fmt.Errorf("ResourceArn mismatch")
		}
		if len(resp.LoggingConfiguration.LogDestinationConfigs) != 1 {
			return fmt.Errorf("expected 1 log destination, got %d", len(resp.LoggingConfiguration.LogDestinationConfigs))
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "ListLoggingConfigurations", func() error {
		resp, err := tc.client.ListLoggingConfigurations(tc.ctx, &wafv2.ListLoggingConfigurationsInput{
			Scope: tc.scope,
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.LoggingConfigurations == nil {
			return fmt.Errorf("LoggingConfigurations is nil")
		}
		found := false
		for _, lc := range resp.LoggingConfigurations {
			if aws.ToString(lc.ResourceArn) == logACLARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("logging configuration not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "DeleteLoggingConfiguration", func() error {
		_, err := tc.client.DeleteLoggingConfiguration(tc.ctx, &wafv2.DeleteLoggingConfigurationInput{
			ResourceArn: aws.String(logACLARN),
		})
		if err != nil {
			return err
		}
		_, verifyErr := tc.client.GetLoggingConfiguration(tc.ctx, &wafv2.GetLoggingConfigurationInput{
			ResourceArn: aws.String(logACLARN),
		})
		if verifyErr == nil {
			return fmt.Errorf("expected error after delete, got nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("wafv2", "GetLoggingConfiguration_AfterDelete", func() error {
		_, err := tc.client.GetLoggingConfiguration(tc.ctx, &wafv2.GetLoggingConfigurationInput{
			ResourceArn: aws.String(logACLARN),
		})
		return AssertErrorContains(err, "WAFNonexistentItemException")
	}))

	results = append(results, r.RunTest("wafv2", "Logging_Cleanup", func() error {
		return tc.deleteWebACL(logACLName, logACLID, logACLLockToken)
	}))

	return results
}
