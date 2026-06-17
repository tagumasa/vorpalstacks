package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func (r *TestRunner) runIoTGroupTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "BillingGroup_Create", func() error {
		out, err := tc.client.CreateBillingGroup(tc.ctx, &iot.CreateBillingGroupInput{
			BillingGroupName: aws.String("test-bg-1"),
		})
		if err != nil {
			return fmt.Errorf("CreateBillingGroup failed: %w", err)
		}
		if out.BillingGroupName == nil || *out.BillingGroupName != "test-bg-1" {
			return fmt.Errorf("expected billingGroupName=test-bg-1, got %v", out.BillingGroupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_Describe_ReturnsVersion", func() error {
		out, err := tc.client.DescribeBillingGroup(tc.ctx, &iot.DescribeBillingGroupInput{
			BillingGroupName: aws.String("test-bg-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeBillingGroup failed: %w", err)
		}
		if out.Version != 1 {
			return fmt.Errorf("expected version=1, got %d", out.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_Update_CAS_Reject", func() error {
		_, err := tc.client.UpdateBillingGroup(tc.ctx, &iot.UpdateBillingGroupInput{
			BillingGroupName: aws.String("test-bg-1"),
			ExpectedVersion:  aws.Int64(999),
			BillingGroupProperties: &types.BillingGroupProperties{
				BillingGroupDescription: aws.String("should fail"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for wrong expectedVersion")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_Update_CAS_Accept", func() error {
		out, err := tc.client.UpdateBillingGroup(tc.ctx, &iot.UpdateBillingGroupInput{
			BillingGroupName: aws.String("test-bg-1"),
			ExpectedVersion:  aws.Int64(1),
			BillingGroupProperties: &types.BillingGroupProperties{
				BillingGroupDescription: aws.String("updated"),
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateBillingGroup failed: %w", err)
		}
		if out.Version != 2 {
			return fmt.Errorf("expected version=2, got %d", out.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_Delete", func() error {
		_, err := tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{
			BillingGroupName: aws.String("test-bg-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteBillingGroup failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_Deprecate_Persist", func() error {
		tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		_, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		if err != nil {
			return fmt.Errorf("CreateThingType failed: %w", err)
		}
		_, err = tc.client.DeprecateThingType(tc.ctx, &iot.DeprecateThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		if err != nil {
			return fmt.Errorf("DeprecateThingType failed: %w", err)
		}
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeMetadata == nil {
			return fmt.Errorf("expected non-nil thingTypeMetadata")
		}
		if !out.ThingTypeMetadata.Deprecated {
			return fmt.Errorf("expected deprecated=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_UndoDeprecate", func() error {
		_, err := tc.client.DeprecateThingType(tc.ctx, &iot.DeprecateThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
			UndoDeprecate: true,
		})
		if err != nil {
			return fmt.Errorf("UndoDeprecate failed: %w", err)
		}
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeMetadata == nil {
			return fmt.Errorf("expected non-nil thingTypeMetadata")
		}
		if out.ThingTypeMetadata.Deprecated {
			return fmt.Errorf("expected deprecated=false after undo")
		}
		_, err = tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{
			ThingTypeName: aws.String("test-type-dep"),
		})
		return nil
	}))

	return results
}
