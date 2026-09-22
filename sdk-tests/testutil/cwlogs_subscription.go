package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) subscriptionTests() []TestResult {
	var results []TestResult
	acct := tc.runner.AccountID()

	results = append(results, tc.runner.RunTest("logs", "PutSubscriptionFilter_VerifyFields", func() error {
		sfName, cleanupGroup, err := tc.newLogGroupFixture("SFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		destARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:test-func", tc.region, acct)
		if err := tc.putSubscriptionFilter(sfName, "TestSub", "ERROR", destARN, roleARN); err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}

		resp, err := tc.collectAllSubscriptionFilters(sfName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp))
		}
		f := resp[0]
		if f.FilterName == nil || *f.FilterName != "TestSub" {
			return fmt.Errorf("filterName mismatch: got %q", aws.ToString(f.FilterName))
		}
		if f.FilterPattern == nil || *f.FilterPattern != "ERROR" {
			return fmt.Errorf("filterPattern mismatch: got %q", aws.ToString(f.FilterPattern))
		}
		if f.DestinationArn == nil || *f.DestinationArn != destARN {
			return fmt.Errorf("destinationArn mismatch: got %q", aws.ToString(f.DestinationArn))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeSubscriptionFilters_Basic", func() error {
		dsfName, cleanupGroup, err := tc.newLogGroupFixture("DSFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		if _, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(dsfName),
			FilterName:     aws.String("DescSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:test", tc.region, acct)),
			RoleArn:        aws.String(roleARN),
		}); err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}

		resp, err := tc.collectAllSubscriptionFilters(dsfName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp))
		}
		if *resp[0].FilterName != "DescSub" {
			return fmt.Errorf("filter name mismatch: got %q", *resp[0].FilterName)
		}
		if *resp[0].DestinationArn != fmt.Sprintf("arn:aws:lambda:%s:%s:function:test", tc.region, acct) {
			return fmt.Errorf("destination arn mismatch")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteSubscriptionFilter_Basic", func() error {
		delSFName, cleanupGroup, err := tc.newLogGroupFixture("DelSFGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		if _, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(delSFName),
			FilterName:     aws.String("DelSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:test", tc.region, acct)),
			RoleArn:        aws.String(roleARN),
		}); err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}

		_, err = tc.client.DeleteSubscriptionFilter(tc.ctx, &cloudwatchlogs.DeleteSubscriptionFilterInput{
			LogGroupName: aws.String(delSFName),
			FilterName:   aws.String("DelSub"),
		})
		if err != nil {
			return fmt.Errorf("delete subscription filter: %v", err)
		}

		resp, err := tc.collectAllSubscriptionFilters(delSFName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", len(resp))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutSubscriptionFilter_DestinationARNForms", func() error {
		dfName, cleanupGroup, err := tc.newLogGroupFixture("DestFormGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		destName := tc.uniquePrefix("central")
		destResp, err := tc.client.PutDestination(tc.ctx, &cloudwatchlogs.PutDestinationInput{
			DestinationName: aws.String(destName),
			TargetArn:       aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/dest-target", tc.region, acct)),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/cwl-deliver", acct)),
		})
		if err != nil {
			return fmt.Errorf("put destination: %v", err)
		}
		defer tc.client.DeleteDestination(tc.ctx, &cloudwatchlogs.DeleteDestinationInput{
			DestinationName: aws.String(destName),
		})

		destARN := aws.ToString(destResp.Destination.Arn)
		if _, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(dfName),
			FilterName:     aws.String("DestSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String(destARN),
		}); err != nil {
			return fmt.Errorf("CloudWatch Logs destination ARN must be accepted: %v", err)
		}

		if _, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(dfName),
			FilterName:     aws.String("FirehoseSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String(fmt.Sprintf("arn:aws:firehose:%s:%s:deliverystream/logs-dest", tc.region, acct)),
		}); err == nil {
			return fmt.Errorf("firehose destination ARN must be rejected until the platform Firehose service exists")
		}
		return nil
	}))

	// The transform-era subscription-filter members survive the full
	// Put→Describe round-trip (persistence carries them).
	results = append(results, tc.runner.RunTest("logs", "PutSubscriptionFilter_TransformSelectionMembers_Roundtrip", func() error {
		selName, cleanupGroup, err := tc.newLogGroupFixture("SfSelectionGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		if _, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:           aws.String(selName),
			FilterName:             aws.String("select"),
			FilterPattern:          aws.String("ERROR"),
			DestinationArn:         aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:test", tc.region, acct)),
			RoleArn:                aws.String(roleARN),
			ApplyOnTransformedLogs: true,
			FieldSelectionCriteria: aws.String(`@aws.region NOT IN ["cn-north-1"]`),
			EmitSystemFields:       []string{"@aws.region", "@source.log"},
		}); err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}

		resp, err := tc.collectAllSubscriptionFilters(selName, "")
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected one filter, got %d", len(resp))
		}
		f := resp[0]
		if f.ApplyOnTransformedLogs != true {
			return fmt.Errorf("applyOnTransformedLogs must survive the round-trip, got %t", f.ApplyOnTransformedLogs)
		}
		if aws.ToString(f.FieldSelectionCriteria) != `@aws.region NOT IN ["cn-north-1"]` {
			return fmt.Errorf("fieldSelectionCriteria must survive the round-trip, got %q", aws.ToString(f.FieldSelectionCriteria))
		}
		if len(f.EmitSystemFields) != 2 || f.EmitSystemFields[0] != "@aws.region" || f.EmitSystemFields[1] != "@source.log" {
			return fmt.Errorf("emitSystemFields must survive the round-trip, got %v", f.EmitSystemFields)
		}
		return nil
	}))

	return results
}
