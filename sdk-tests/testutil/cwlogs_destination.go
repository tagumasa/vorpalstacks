package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) destinationTests() []TestResult {
	var results []TestResult
	acct := tc.runner.AccountID()

	results = append(results, tc.runner.RunTest("logs", "PutDestination_VerifyFields", func() error {
		destName := tc.uniquePrefix("TestDest")
		resp, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(destName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/dest-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/test-stream", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("put destination: %v", err)
		}
		if resp.Destination == nil {
			return fmt.Errorf("destination is nil")
		}
		if resp.Destination.DestinationName == nil || *resp.Destination.DestinationName != destName {
			return fmt.Errorf("destinationName mismatch: got %q", aws.ToString(resp.Destination.DestinationName))
		}
		if resp.Destination.TargetArn == nil || *resp.Destination.TargetArn != fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/test-stream", tc.region, acct) {
			return fmt.Errorf("targetArn mismatch: got %q", aws.ToString(resp.Destination.TargetArn))
		}
		if resp.Destination.RoleArn == nil || *resp.Destination.RoleArn != fmt.Sprintf("arn:aws:iam::%s:role/dest-role", acct) {
			return fmt.Errorf("roleArn mismatch: got %q", aws.ToString(resp.Destination.RoleArn))
		}
		if resp.Destination.Arn == nil || *resp.Destination.Arn == "" {
			return fmt.Errorf("arn is nil or empty")
		}
		tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(destName),
		})
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeDestinations_Basic", func() error {
		ddName := tc.uniquePrefix("DescDest")
		_, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(ddName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/dd-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/dd-stream", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(ddName),
		})

		resp, err := tc.client.DescribeDestinations(tc.ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(ddName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) != 1 {
			return fmt.Errorf("expected 1 destination, got %d", len(resp.Destinations))
		}
		if *resp.Destinations[0].DestinationName != ddName {
			return fmt.Errorf("name mismatch: got %q", *resp.Destinations[0].DestinationName)
		}
		if resp.Destinations[0].Arn == nil || *resp.Destinations[0].Arn == "" {
			return fmt.Errorf("ARN is empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutDestinationPolicy_Basic", func() error {
		pdpName := tc.uniquePrefix("PdpDest")
		_, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(pdpName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/pdp-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/pdp-stream", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(pdpName),
		})

		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"logs:PutSubscriptionFilter"}]}`
		_, err = tc.client.PutDestinationPolicy(tc.ctx, &cloudwatchlogs.PutDestinationPolicyInput{
			DestinationName: aws.String(pdpName),
			AccessPolicy:    aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}

		resp, err := tc.client.DescribeDestinations(tc.ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(pdpName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) == 0 {
			return fmt.Errorf("destination not found")
		}
		if resp.Destinations[0].AccessPolicy == nil || *resp.Destinations[0].AccessPolicy != policy {
			return fmt.Errorf("access policy mismatch")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutDestination_UpdateInPlace", func() error {
		udpName := tc.uniquePrefix("UdpDest")
		_, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(udpName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/original-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/original", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("initial put: %v", err)
		}
		defer tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(udpName),
		})

		_, err = tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(udpName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/updated-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/updated", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("update put: %v", err)
		}

		resp, err := tc.client.DescribeDestinations(tc.ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(udpName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) == 0 {
			return fmt.Errorf("destination not found")
		}
		if *resp.Destinations[0].RoleArn != fmt.Sprintf("arn:aws:iam::%s:role/updated-role", acct) {
			return fmt.Errorf("roleArn not updated: got %q", *resp.Destinations[0].RoleArn)
		}
		if *resp.Destinations[0].TargetArn != fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/updated", tc.region, acct) {
			return fmt.Errorf("targetArn not updated: got %q", *resp.Destinations[0].TargetArn)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteDestination_Basic", func() error {
		ddelName := tc.uniquePrefix("DelDest")
		_, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(ddelName),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/ddel-role", acct)),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/ddel-stream", tc.region, acct)),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(ddelName),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		resp, err := tc.client.DescribeDestinations(tc.ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(ddelName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Destinations) != 0 {
			return fmt.Errorf("expected 0 destinations after delete, got %d", len(resp.Destinations))
		}
		return nil
	}))

	return results
}
