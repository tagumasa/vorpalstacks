package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func (r *TestRunner) runIoTThingTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "ThingType_CreateThingType", func() error {
		out, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String("test-thing-type-1"),
		})
		if err != nil {
			return fmt.Errorf("CreateThingType failed: %w", err)
		}
		if out.ThingTypeName == nil || *out.ThingTypeName != "test-thing-type-1" {
			return fmt.Errorf("expected thingTypeName=test-thing-type-1, got %v", out.ThingTypeName)
		}
		if out.ThingTypeArn == nil || *out.ThingTypeArn == "" {
			return fmt.Errorf("expected non-empty thingTypeArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_DescribeThingType", func() error {
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{
			ThingTypeName: aws.String("test-thing-type-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeName == nil || *out.ThingTypeName != "test-thing-type-1" {
			return fmt.Errorf("expected thingTypeName=test-thing-type-1, got %v", out.ThingTypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_DeleteThingType", func() error {
		_, err := tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{
			ThingTypeName: aws.String("test-thing-type-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteThingType failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_CreateThingGroup", func() error {
		out, err := tc.client.CreateThingGroup(tc.ctx, &iot.CreateThingGroupInput{
			ThingGroupName: aws.String("test-thing-group-1"),
		})
		if err != nil {
			return fmt.Errorf("CreateThingGroup failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != "test-thing-group-1" {
			return fmt.Errorf("expected thingGroupName=test-thing-group-1, got %v", out.ThingGroupName)
		}
		if out.ThingGroupArn == nil || *out.ThingGroupArn == "" {
			return fmt.Errorf("expected non-empty thingGroupArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DescribeThingGroup", func() error {
		out, err := tc.client.DescribeThingGroup(tc.ctx, &iot.DescribeThingGroupInput{
			ThingGroupName: aws.String("test-thing-group-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingGroup failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != "test-thing-group-1" {
			return fmt.Errorf("expected thingGroupName=test-thing-group-1, got %v", out.ThingGroupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DeleteThingGroup", func() error {
		_, err := tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{
			ThingGroupName: aws.String("test-thing-group-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteThingGroup failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_CreateThing", func() error {
		out, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName: aws.String("test-thing-1"),
			AttributePayload: &types.AttributePayload{
				Attributes: map[string]string{
					"version":  "1.0",
					"location": "us-east-1",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != "test-thing-1" {
			return fmt.Errorf("expected thingName=test-thing-1, got %v", out.ThingName)
		}
		if out.ThingArn == nil || *out.ThingArn == "" {
			return fmt.Errorf("expected non-empty thingArn")
		}
		if out.ThingId == nil || *out.ThingId == "" {
			return fmt.Errorf("expected non-empty thingId")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DescribeThing", func() error {
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{
			ThingName: aws.String("test-thing-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != "test-thing-1" {
			return fmt.Errorf("expected thingName=test-thing-1, got %v", out.ThingName)
		}
		if out.Attributes == nil {
			return fmt.Errorf("expected attributes to be non-nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_UpdateThing", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName: aws.String("test-thing-1"),
			AttributePayload: &types.AttributePayload{
				Attributes: map[string]string{
					"version": "2.0",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateThing failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_ListThings", func() error {
		out, err := tc.client.ListThings(tc.ctx, &iot.ListThingsInput{})
		if err != nil {
			return fmt.Errorf("ListThings failed: %w", err)
		}
		if out.Things == nil {
			return fmt.Errorf("expected non-nil things list")
		}
		found := false
		for _, t := range out.Things {
			if t.ThingName != nil && *t.ThingName == "test-thing-1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test-thing-1 not found in list of %d things", len(out.Things))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_CreateThing_Duplicate", func() error {
		_, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName: aws.String("test-thing-1"),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate thing creation")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DeleteThing", func() error {
		_, err := tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{
			ThingName: aws.String("test-thing-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteThing failed: %w", err)
		}
		return nil
	}))

	return results
}
