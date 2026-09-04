package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaESMTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	esmFuncName, cleanupEsmFunc, err := tc.setupFunction("EsmFunc", "exports.handler = async () => { return 1; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to set up function: %v", err)}}
	}
	defer cleanupEsmFunc()

	esmEventSourceArn := fmt.Sprintf("arn:aws:sqs:%s:%s:test-queue", tc.r.region, tc.r.accountID)

	var esmUUID string

	results = append(results, tc.r.RunTest("lambda", "CreateEventSourceMapping", func() error {
		resp, err := tc.client.CreateEventSourceMapping(tc.ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(esmFuncName),
			EventSourceArn: aws.String(esmEventSourceArn),
			BatchSize:      aws.Int32(10),
			Enabled:        aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.UUID == nil || *resp.UUID == "" {
			return fmt.Errorf("UUID is nil or empty")
		}
		if resp.EventSourceArn == nil || *resp.EventSourceArn != esmEventSourceArn {
			return fmt.Errorf("EventSourceArn mismatch, got %v", resp.EventSourceArn)
		}
		if resp.BatchSize == nil || *resp.BatchSize != 10 {
			return fmt.Errorf("BatchSize mismatch, expected 10, got %v", resp.BatchSize)
		}
		if resp.State == nil || *resp.State == "" {
			return fmt.Errorf("State is nil or empty")
		}
		esmUUID = *resp.UUID
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		resp, err := tc.client.GetEventSourceMapping(tc.ctx, &lambda.GetEventSourceMappingInput{
			UUID: aws.String(esmUUID),
		})
		if err != nil {
			return err
		}
		if resp.FunctionArn == nil {
			return fmt.Errorf("FunctionArn is nil")
		}
		if resp.EventSourceArn == nil || *resp.EventSourceArn != esmEventSourceArn {
			return fmt.Errorf("EventSourceArn mismatch, got %v", resp.EventSourceArn)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		resp, err := tc.client.UpdateEventSourceMapping(tc.ctx, &lambda.UpdateEventSourceMappingInput{
			UUID:                           aws.String(esmUUID),
			BatchSize:                      aws.Int32(50),
			MaximumBatchingWindowInSeconds: aws.Int32(1),
			Enabled:                        aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.BatchSize == nil || *resp.BatchSize != 50 {
			return fmt.Errorf("BatchSize not updated, got %v", resp.BatchSize)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListEventSourceMappings", func() error {
		resp, err := tc.client.ListEventSourceMappings(tc.ctx, &lambda.ListEventSourceMappingsInput{
			FunctionName: aws.String(esmFuncName),
		})
		if err != nil {
			return err
		}
		if resp.EventSourceMappings == nil {
			return fmt.Errorf("event source mappings list is nil")
		}
		if len(resp.EventSourceMappings) == 0 {
			return fmt.Errorf("expected at least 1 event source mapping")
		}
		found := false
		for _, m := range resp.EventSourceMappings {
			if m.UUID != nil && *m.UUID == esmUUID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created ESM %s not found in ListEventSourceMappings", esmUUID)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		_, err := tc.client.DeleteEventSourceMapping(tc.ctx, &lambda.DeleteEventSourceMappingInput{
			UUID: aws.String(esmUUID),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetEventSourceMapping(tc.ctx, &lambda.GetEventSourceMappingInput{
			UUID: aws.String(esmUUID),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetEventSourceMapping_NonExistent", func() error {
		_, err := tc.client.GetEventSourceMapping(tc.ctx, &lambda.GetEventSourceMappingInput{
			UUID: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateEventSourceMapping_SQSDefaultBatchSize", func() error {
		resp, err := tc.client.CreateEventSourceMapping(tc.ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(esmFuncName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:default-batch-queue", tc.r.region, tc.r.accountID)),
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteEventSourceMapping(tc.ctx, &lambda.DeleteEventSourceMappingInput{UUID: resp.UUID})
		// SQS event sources default to a batch size of 10.
		if resp.BatchSize == nil || *resp.BatchSize != 10 {
			return fmt.Errorf("default SQS BatchSize should be 10, got %v", resp.BatchSize)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateEventSourceMapping_FIFOBatchSizeCap", func() error {
		_, err := tc.client.CreateEventSourceMapping(tc.ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(esmFuncName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:orders.fifo", tc.r.region, tc.r.accountID)),
			BatchSize:      aws.Int32(11),
		})
		if err := expectAWSErrorCode(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateEventSourceMapping_StartingPositionTimestamp", func() error {
		ts := time.Unix(time.Now().Add(-1*time.Hour).Unix(), 0).UTC()
		resp, err := tc.client.CreateEventSourceMapping(tc.ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:              aws.String(esmFuncName),
			EventSourceArn:            aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/esm-ts-stream", tc.r.region, tc.r.accountID)),
			StartingPosition:          types.EventSourcePositionAtTimestamp,
			StartingPositionTimestamp: aws.Time(ts),
		})
		if err != nil {
			return err
		}
		defer tc.client.DeleteEventSourceMapping(tc.ctx, &lambda.DeleteEventSourceMappingInput{UUID: resp.UUID})
		if resp.StartingPosition != types.EventSourcePositionAtTimestamp {
			return fmt.Errorf("StartingPosition mismatch, got %v", resp.StartingPosition)
		}
		if resp.StartingPositionTimestamp == nil || !resp.StartingPositionTimestamp.Equal(ts) {
			return fmt.Errorf("StartingPositionTimestamp not echoed, got %v", resp.StartingPositionTimestamp)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListEventSourceMappings_MaxItemsCap", func() error {
		// ListEventSourceMappings allows up to 100 items per response.
		before, err := tc.client.ListEventSourceMappings(tc.ctx, &lambda.ListEventSourceMappingsInput{
			FunctionName: aws.String(esmFuncName),
			MaxItems:     aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("list before: %v", err)
		}

		const mappingCount = 55
		var uuids []string
		for i := 0; i < mappingCount; i++ {
			created, err := tc.client.CreateEventSourceMapping(tc.ctx, &lambda.CreateEventSourceMappingInput{
				FunctionName:   aws.String(esmFuncName),
				EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:cap-queue-%d", tc.r.region, tc.r.accountID, i)),
			})
			if err != nil {
				return fmt.Errorf("create mapping %d: %v", i, err)
			}
			uuids = append(uuids, *created.UUID)
		}
		defer func() {
			for _, id := range uuids {
				_, _ = tc.client.DeleteEventSourceMapping(tc.ctx, &lambda.DeleteEventSourceMappingInput{UUID: aws.String(id)})
			}
		}()

		resp, err := tc.client.ListEventSourceMappings(tc.ctx, &lambda.ListEventSourceMappingsInput{
			FunctionName: aws.String(esmFuncName),
			MaxItems:     aws.Int32(100),
		})
		if err != nil {
			return err
		}
		want := len(before.EventSourceMappings) + mappingCount
		if len(resp.EventSourceMappings) != want {
			return fmt.Errorf("expected all %d mappings in one response, got %d", want, len(resp.EventSourceMappings))
		}
		return nil
	}))

	return results
}
