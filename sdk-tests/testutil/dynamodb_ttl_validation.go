package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

// createTTLTestTable creates a simple on-demand table for TTL tests.
func createTTLTestTable(ctx context.Context, client *dynamodb.Client, name string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	})
	if err != nil {
		return err
	}
	return waitKinesisDestTableActive(ctx, client, name)
}

// ttlUpdateIsValidation reports whether err is a ValidationException.
func ttlUpdateIsValidation(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.ErrorCode() == "ValidationException"
}

// dynamoDBTTLValidationTests pins the documented UpdateTimeToLive contract:
// the response carries only Enabled and AttributeName, enabling TTL on a
// table that already has TTL enabled (same or different attribute) is
// rejected, renaming requires a disable in between, and an empty attribute
// name is rejected because the member is required.
func (r *TestRunner) dynamoDBTTLValidationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	// Table A: re-enable and rename while enabled are both rejected.
	tableA := fmt.Sprintf("ttl-validate-a-%d", suffix)
	if err := createTTLTestTable(ctx, client, tableA); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "UpdateTimeToLive_ReEnableWhileEnabled_Rejected",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create table %s: %v", tableA, err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableA)})

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive_ReEnableWhileEnabled_Rejected", func() error {
		_, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableA),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		_, err = client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableA),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for re-enabling an enabled TTL")
		}
		if !ttlUpdateIsValidation(err) {
			return fmt.Errorf("expected ValidationException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive_RenameWhileEnabled_Rejected", func() error {
		_, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableA),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("expiry"),
				Enabled:       aws.Bool(true),
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for renaming an enabled TTL attribute")
		}
		if !ttlUpdateIsValidation(err) {
			return fmt.Errorf("expected ValidationException, got %v", err)
		}
		return nil
	}))

	// Table B: the documented rename path — disable, then re-enable with
	// the new attribute name — succeeds and the response echoes only the
	// requested specification members.
	tableB := fmt.Sprintf("ttl-validate-b-%d", suffix)
	if err := createTTLTestTable(ctx, client, tableB); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "UpdateTimeToLive_RenameViaDisable",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create table %s: %v", tableB, err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableB)})

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive_RenameViaDisable", func() error {
		if _, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableB),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		}); err != nil {
			return err
		}
		disableResp, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableB),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(false),
			},
		})
		if err != nil {
			return err
		}
		if disableResp.TimeToLiveSpecification == nil {
			return fmt.Errorf("disable response has no TimeToLiveSpecification")
		}
		if *disableResp.TimeToLiveSpecification.Enabled {
			return fmt.Errorf("expected Enabled=false on disable")
		}
		if aws.ToString(disableResp.TimeToLiveSpecification.AttributeName) != "ttl" {
			return fmt.Errorf("expected AttributeName=ttl, got %v", disableResp.TimeToLiveSpecification.AttributeName)
		}

		enableResp, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableB),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("expiry"),
				Enabled:       aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		if enableResp.TimeToLiveSpecification == nil {
			return fmt.Errorf("enable response has no TimeToLiveSpecification")
		}
		if !*enableResp.TimeToLiveSpecification.Enabled {
			return fmt.Errorf("expected Enabled=true after rename")
		}
		if aws.ToString(enableResp.TimeToLiveSpecification.AttributeName) != "expiry" {
			return fmt.Errorf("expected AttributeName=expiry, got %v", enableResp.TimeToLiveSpecification.AttributeName)
		}

		desc, err := client.DescribeTimeToLive(ctx, &dynamodb.DescribeTimeToLiveInput{
			TableName: aws.String(tableB),
		})
		if err != nil {
			return err
		}
		if desc.TimeToLiveDescription == nil {
			return fmt.Errorf("DescribeTimeToLive returned no description")
		}
		if desc.TimeToLiveDescription.TimeToLiveStatus != dynamodbtypes.TimeToLiveStatusEnabled {
			return fmt.Errorf("expected TimeToLiveStatus=ENABLED, got %v", desc.TimeToLiveDescription.TimeToLiveStatus)
		}
		if aws.ToString(desc.TimeToLiveDescription.AttributeName) != "expiry" {
			return fmt.Errorf("expected described AttributeName=expiry, got %v", desc.TimeToLiveDescription.AttributeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive_EmptyAttributeName_Rejected", func() error {
		// Disable first so the rejection can only come from the empty
		// attribute name, not from the enable-while-enabled rule.
		if _, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableB),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String("expiry"),
				Enabled:       aws.Bool(false),
			},
		}); err != nil {
			return err
		}
		_, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableB),
			TimeToLiveSpecification: &dynamodbtypes.TimeToLiveSpecification{
				AttributeName: aws.String(""),
				Enabled:       aws.Bool(true),
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for an empty attribute name")
		}
		if !ttlUpdateIsValidation(err) {
			return fmt.Errorf("expected ValidationException, got %v", err)
		}
		return nil
	}))

	return results
}
