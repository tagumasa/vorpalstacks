package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaESMTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	cwlClient *cloudwatchlogs.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	esmFuncName := fmt.Sprintf("EsmFunc-%d", time.Now().UnixNano())
	esmRoleName := fmt.Sprintf("EsmRole-%d", time.Now().UnixNano())
	esmRole := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, esmRoleName)
	esmCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to zip lambda code: %v", err)}}
	}
	esmEventSourceArn := fmt.Sprintf("arn:aws:sqs:%s:%s:test-queue", r.region, r.accountID)

	if err := createIAMRole(esmRoleName); err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer deleteIAMRole(esmRoleName)

	_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(esmFuncName),
		Runtime:      types.RuntimeNodejs22x,
		Role:         aws.String(esmRole),
		Handler:      aws.String("index.handler"),
		Code:         &types.FunctionCode{ZipFile: esmCode},
	})
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateEventSourceMapping_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(esmFuncName)})
	defer deleteLambdaLogGroup(cwlClient, ctx, esmFuncName)

	var esmUUID string

	results = append(results, r.RunTest("lambda", "CreateEventSourceMapping", func() error {
		resp, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
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

	results = append(results, r.RunTest("lambda", "GetEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		resp, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
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

	results = append(results, r.RunTest("lambda", "UpdateEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		resp, err := client.UpdateEventSourceMapping(ctx, &lambda.UpdateEventSourceMappingInput{
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

	results = append(results, r.RunTest("lambda", "ListEventSourceMappings", func() error {
		resp, err := client.ListEventSourceMappings(ctx, &lambda.ListEventSourceMappingsInput{
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

	results = append(results, r.RunTest("lambda", "DeleteEventSourceMapping", func() error {
		if esmUUID == "" {
			return fmt.Errorf("no UUID from CreateEventSourceMapping")
		}
		_, err := client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{
			UUID: aws.String(esmUUID),
		})
		if err != nil {
			return err
		}
		_, err = client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
			UUID: aws.String(esmUUID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetEventSourceMapping_NonExistent", func() error {
		_, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{
			UUID: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateEventSourceMapping_SQSDefaultBatchSize", func() error {
		resp, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(esmFuncName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:default-batch-queue", r.region, r.accountID)),
		})
		if err != nil {
			return err
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: resp.UUID})
		// SQS event sources default to a batch size of 10.
		if resp.BatchSize == nil || *resp.BatchSize != 10 {
			return fmt.Errorf("default SQS BatchSize should be 10, got %v", resp.BatchSize)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateEventSourceMapping_FIFOBatchSizeCap", func() error {
		_, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(esmFuncName),
			EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:orders.fifo", r.region, r.accountID)),
			BatchSize:      aws.Int32(11),
		})
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "CreateEventSourceMapping_StartingPositionTimestamp", func() error {
		ts := time.Unix(time.Now().Add(-1*time.Hour).Unix(), 0).UTC()
		resp, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
			FunctionName:              aws.String(esmFuncName),
			EventSourceArn:            aws.String(fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/esm-ts-stream", r.region, r.accountID)),
			StartingPosition:          types.EventSourcePositionAtTimestamp,
			StartingPositionTimestamp: aws.Time(ts),
		})
		if err != nil {
			return err
		}
		defer client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: resp.UUID})
		if resp.StartingPosition != types.EventSourcePositionAtTimestamp {
			return fmt.Errorf("StartingPosition mismatch, got %v", resp.StartingPosition)
		}
		if resp.StartingPositionTimestamp == nil || !resp.StartingPositionTimestamp.Equal(ts) {
			return fmt.Errorf("StartingPositionTimestamp not echoed, got %v", resp.StartingPositionTimestamp)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListEventSourceMappings_MaxItemsCap", func() error {
		// ListEventSourceMappings allows up to 100 items per response.
		before, err := client.ListEventSourceMappings(ctx, &lambda.ListEventSourceMappingsInput{
			FunctionName: aws.String(esmFuncName),
			MaxItems:     aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("list before: %v", err)
		}

		const mappingCount = 55
		var uuids []string
		for i := 0; i < mappingCount; i++ {
			created, err := client.CreateEventSourceMapping(ctx, &lambda.CreateEventSourceMappingInput{
				FunctionName:   aws.String(esmFuncName),
				EventSourceArn: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:cap-queue-%d", r.region, r.accountID, i)),
			})
			if err != nil {
				return fmt.Errorf("create mapping %d: %v", i, err)
			}
			uuids = append(uuids, *created.UUID)
		}
		defer func() {
			for _, id := range uuids {
				_, _ = client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: aws.String(id)})
			}
		}()

		resp, err := client.ListEventSourceMappings(ctx, &lambda.ListEventSourceMappingsInput{
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
