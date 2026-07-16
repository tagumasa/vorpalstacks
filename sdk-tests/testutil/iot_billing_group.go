package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTBillingGroupTests covers the BillingGroup lifecycle: Create/Describe/
// Update (with optimistic-concurrency CAS)/List/Delete, nested description, and
// Add/Remove thing membership.
func (r *TestRunner) runIoTBillingGroupTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	bgName := uniqueName("billing-group")
	thingName := uniqueName("thing-bg")
	defer tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{BillingGroupName: aws.String(bgName)})
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	results = append(results, r.RunTest("iot", "BillingGroup_CreateBillingGroup", func() error {
		out, err := tc.client.CreateBillingGroup(tc.ctx, &iot.CreateBillingGroupInput{
			BillingGroupName: aws.String(bgName),
			BillingGroupProperties: &types.BillingGroupProperties{
				BillingGroupDescription: aws.String("billing desc"),
			},
		})
		if err != nil {
			return fmt.Errorf("CreateBillingGroup failed: %w", err)
		}
		if out.BillingGroupName == nil || *out.BillingGroupName != bgName {
			return fmt.Errorf("expected billingGroupName=%s, got %v", bgName, out.BillingGroupName)
		}
		if out.BillingGroupArn == nil || *out.BillingGroupArn == "" {
			return fmt.Errorf("expected non-empty billingGroupArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_DescribeBillingGroup_Version", func() error {
		out, err := tc.client.DescribeBillingGroup(tc.ctx, &iot.DescribeBillingGroupInput{BillingGroupName: aws.String(bgName)})
		if err != nil {
			return fmt.Errorf("DescribeBillingGroup failed: %w", err)
		}
		if out.Version != 1 {
			return fmt.Errorf("expected version=1, got %d", out.Version)
		}
		if out.BillingGroupProperties == nil {
			return fmt.Errorf("expected billingGroupProperties to be non-nil")
		}
		if out.BillingGroupProperties.BillingGroupDescription == nil || *out.BillingGroupProperties.BillingGroupDescription != "billing desc" {
			return fmt.Errorf("expected billingGroupDescription to round-trip")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_UpdateBillingGroup_CAS_Reject", func() error {
		// Wrong expectedVersion must be rejected.
		_, err := tc.client.UpdateBillingGroup(tc.ctx, &iot.UpdateBillingGroupInput{
			BillingGroupName: aws.String(bgName),
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

	results = append(results, r.RunTest("iot", "BillingGroup_UpdateBillingGroup_CAS_Accept", func() error {
		desc, err := tc.client.DescribeBillingGroup(tc.ctx, &iot.DescribeBillingGroupInput{BillingGroupName: aws.String(bgName)})
		if err != nil {
			return fmt.Errorf("DescribeBillingGroup failed: %w", err)
		}
		updated, err := tc.client.UpdateBillingGroup(tc.ctx, &iot.UpdateBillingGroupInput{
			BillingGroupName: aws.String(bgName),
			ExpectedVersion:  aws.Int64(desc.Version),
			BillingGroupProperties: &types.BillingGroupProperties{
				BillingGroupDescription: aws.String("updated"),
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateBillingGroup failed: %w", err)
		}
		if updated.Version != desc.Version+1 {
			return fmt.Errorf("expected version=%d, got %d", desc.Version+1, updated.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_ListBillingGroups_IncludesCreated", func() error {
		found, err := tc.billingGroupExists(bgName)
		if err != nil {
			return fmt.Errorf("ListBillingGroups failed: %w", err)
		}
		if !found {
			return fmt.Errorf("%s not found in billing groups", bgName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_AddThingToBillingGroup", func() error {
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		_, err := tc.client.AddThingToBillingGroup(tc.ctx, &iot.AddThingToBillingGroupInput{
			BillingGroupName: aws.String(bgName),
			ThingName:        aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("AddThingToBillingGroup failed: %w", err)
		}
		found, err := tc.thingInBillingGroupExists(bgName, thingName)
		if err != nil {
			return fmt.Errorf("ListThingsInBillingGroup failed: %w", err)
		}
		if !found {
			return fmt.Errorf("expected %s in billing group %s", thingName, bgName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_DescribeThing_ReturnsBillingGroup", func() error {
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if out.BillingGroupName == nil || *out.BillingGroupName != bgName {
			return fmt.Errorf("expected billingGroupName=%s, got %v", bgName, out.BillingGroupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_AddThing_SecondGroup_Rejected", func() error {
		bgName2 := uniqueName("billing-group-2")
		defer tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{BillingGroupName: aws.String(bgName2)})
		if _, err := tc.client.CreateBillingGroup(tc.ctx, &iot.CreateBillingGroupInput{
			BillingGroupName: aws.String(bgName2),
		}); err != nil {
			return fmt.Errorf("CreateBillingGroup prerequisite failed: %w", err)
		}
		_, err := tc.client.AddThingToBillingGroup(tc.ctx, &iot.AddThingToBillingGroupInput{
			BillingGroupName: aws.String(bgName2),
			ThingName:        aws.String(thingName),
		})
		if err == nil {
			return fmt.Errorf("expected error adding thing to a second billing group")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_DeleteBillingGroup_Conflict", func() error {
		_, err := tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{BillingGroupName: aws.String(bgName)})
		if err == nil {
			return fmt.Errorf("expected DeleteConflictException for non-empty billing group")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_RemoveThingFromBillingGroup", func() error {
		_, err := tc.client.RemoveThingFromBillingGroup(tc.ctx, &iot.RemoveThingFromBillingGroupInput{
			BillingGroupName: aws.String(bgName),
			ThingName:        aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("RemoveThingFromBillingGroup failed: %w", err)
		}
		found, err := tc.thingInBillingGroupExists(bgName, thingName)
		if err != nil {
			return fmt.Errorf("ListThingsInBillingGroup failed: %w", err)
		}
		if found {
			return fmt.Errorf("expected %s to be removed from billing group", thingName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_ListThings_NotFound", func() error {
		_, err := tc.client.ListThingsInBillingGroup(tc.ctx, &iot.ListThingsInBillingGroupInput{
			BillingGroupName: aws.String(uniqueName("nope-bg")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_DescribeBillingGroup_NotFound", func() error {
		_, err := tc.client.DescribeBillingGroup(tc.ctx, &iot.DescribeBillingGroupInput{BillingGroupName: aws.String(uniqueName("nope-bg"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "BillingGroup_DeleteBillingGroup", func() error {
		_, err := tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{BillingGroupName: aws.String(bgName)})
		return err
	}))

	return results
}
