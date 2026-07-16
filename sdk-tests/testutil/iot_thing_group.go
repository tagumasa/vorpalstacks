package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTThingGroupTests covers static and dynamic ThingGroup lifecycle:
// Create/Describe/Update/List/Delete, nested properties, Add/Remove thing
// membership, and DynamicThingGroup CRUD. All names are uniqueName-based.
func (r *TestRunner) runIoTThingGroupTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	// ── Static ThingGroup CRUD ──
	groupName := uniqueName("thing-group")
	thingName := uniqueName("thing-grp")
	defer tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{ThingGroupName: aws.String(groupName)})
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})
	gDesc := "group description"

	results = append(results, r.RunTest("iot", "ThingGroup_CreateThingGroup", func() error {
		out, err := tc.client.CreateThingGroup(tc.ctx, &iot.CreateThingGroupInput{
			ThingGroupName: aws.String(groupName),
			ThingGroupProperties: &types.ThingGroupProperties{
				ThingGroupDescription: aws.String(gDesc),
				AttributePayload: &types.AttributePayload{
					Attributes: map[string]string{"env": "test"},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateThingGroup failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != groupName {
			return fmt.Errorf("expected thingGroupName=%s, got %v", groupName, out.ThingGroupName)
		}
		if out.ThingGroupArn == nil || *out.ThingGroupArn == "" {
			return fmt.Errorf("expected non-empty thingGroupArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DescribeThingGroup", func() error {
		out, err := tc.client.DescribeThingGroup(tc.ctx, &iot.DescribeThingGroupInput{ThingGroupName: aws.String(groupName)})
		if err != nil {
			return fmt.Errorf("DescribeThingGroup failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != groupName {
			return fmt.Errorf("expected thingGroupName=%s, got %v", groupName, out.ThingGroupName)
		}
		if out.ThingGroupProperties == nil {
			return fmt.Errorf("expected thingGroupProperties to be non-nil")
		}
		if out.ThingGroupProperties.ThingGroupDescription == nil || *out.ThingGroupProperties.ThingGroupDescription != gDesc {
			return fmt.Errorf("expected thingGroupDescription=%q, got %v", gDesc, out.ThingGroupProperties.ThingGroupDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_UpdateThingGroup_CAS", func() error {
		// Wrong expectedVersion must be rejected.
		_, err := tc.client.UpdateThingGroup(tc.ctx, &iot.UpdateThingGroupInput{
			ThingGroupName:  aws.String(groupName),
			ExpectedVersion: aws.Int64(999),
			ThingGroupProperties: &types.ThingGroupProperties{
				ThingGroupDescription: aws.String("should fail"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for wrong expectedVersion")
		}
		// Update without a CAS token (nil ExpectedVersion) must succeed and
		// persist the new description.
		_, err = tc.client.UpdateThingGroup(tc.ctx, &iot.UpdateThingGroupInput{
			ThingGroupName: aws.String(groupName),
			ThingGroupProperties: &types.ThingGroupProperties{
				ThingGroupDescription: aws.String("updated"),
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateThingGroup failed: %w", err)
		}
		d, err := tc.client.DescribeThingGroup(tc.ctx, &iot.DescribeThingGroupInput{ThingGroupName: aws.String(groupName)})
		if err != nil {
			return fmt.Errorf("DescribeThingGroup after update failed: %w", err)
		}
		if d.ThingGroupProperties == nil || d.ThingGroupProperties.ThingGroupDescription == nil || *d.ThingGroupProperties.ThingGroupDescription != "updated" {
			return fmt.Errorf("expected thingGroupDescription=updated after CAS update")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_ListThingGroups_IncludesCreated", func() error {
		found, err := tc.thingGroupExists(groupName)
		if err != nil {
			return fmt.Errorf("ListThingGroups failed: %w", err)
		}
		if !found {
			return fmt.Errorf("%s not found in thing groups", groupName)
		}
		return nil
	}))

	// ── Thing membership ──
	results = append(results, r.RunTest("iot", "ThingGroup_AddThingToThingGroup", func() error {
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		_, err := tc.client.AddThingToThingGroup(tc.ctx, &iot.AddThingToThingGroupInput{
			ThingGroupName: aws.String(groupName),
			ThingName:      aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("AddThingToThingGroup failed: %w", err)
		}
		found, err := tc.thingInGroupExists(groupName, thingName)
		if err != nil {
			return fmt.Errorf("ListThingsInThingGroup failed: %w", err)
		}
		if !found {
			return fmt.Errorf("expected %s in thing group %s", thingName, groupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_RemoveThingFromThingGroup", func() error {
		_, err := tc.client.RemoveThingFromThingGroup(tc.ctx, &iot.RemoveThingFromThingGroupInput{
			ThingGroupName: aws.String(groupName),
			ThingName:      aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("RemoveThingFromThingGroup failed: %w", err)
		}
		found, err := tc.thingInGroupExists(groupName, thingName)
		if err != nil {
			return fmt.Errorf("ListThingsInThingGroup failed: %w", err)
		}
		if found {
			return fmt.Errorf("expected %s to be removed from group", thingName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DescribeThingGroup_NotFound", func() error {
		_, err := tc.client.DescribeThingGroup(tc.ctx, &iot.DescribeThingGroupInput{ThingGroupName: aws.String(uniqueName("nope-group"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DeleteThingGroup", func() error {
		_, err := tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{ThingGroupName: aws.String(groupName)})
		return err
	}))

	// ── Dynamic ThingGroup ──
	dynName := uniqueName("dyn-group")
	defer tc.client.DeleteDynamicThingGroup(tc.ctx, &iot.DeleteDynamicThingGroupInput{ThingGroupName: aws.String(dynName)})

	results = append(results, r.RunTest("iot", "ThingGroup_CreateDynamicThingGroup", func() error {
		out, err := tc.client.CreateDynamicThingGroup(tc.ctx, &iot.CreateDynamicThingGroupInput{
			ThingGroupName:       aws.String(dynName),
			ThingGroupProperties: &types.ThingGroupProperties{},
			QueryString:          aws.String("thingName:*"),
			QueryVersion:         aws.String("2016-11-30"),
			IndexName:            aws.String("AWS_Things"),
		})
		if err != nil {
			return fmt.Errorf("CreateDynamicThingGroup failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != dynName {
			return fmt.Errorf("expected thingGroupName=%s, got %v", dynName, out.ThingGroupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_UpdateDynamicThingGroup", func() error {
		_, err := tc.client.UpdateDynamicThingGroup(tc.ctx, &iot.UpdateDynamicThingGroupInput{
			ThingGroupName: aws.String(dynName),
			ThingGroupProperties: &types.ThingGroupProperties{
				ThingGroupDescription: aws.String("updated"),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ThingGroup_DeleteDynamicThingGroup", func() error {
		_, err := tc.client.DeleteDynamicThingGroup(tc.ctx, &iot.DeleteDynamicThingGroupInput{ThingGroupName: aws.String(dynName)})
		return err
	}))

	return results
}
