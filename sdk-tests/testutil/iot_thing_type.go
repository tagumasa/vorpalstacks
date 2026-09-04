package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTThingTypeTests covers the ThingType lifecycle: Create/Describe/Update/
// List/Delete, nested properties (description, searchable attributes), and the
// DeprecateThingType / UndoDeprecate toggle, with negative NotFound paths.
func (r *TestRunner) runIoTThingTypeTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	// ── Basic CRUD with nested description ──
	ttName := uniqueName("thing-type")
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttName)})
	desc := "type with description"

	results = append(results, r.RunTest("iot", "ThingType_CreateThingType", func() error {
		out, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String(ttName),
			ThingTypeProperties: &types.ThingTypeProperties{
				ThingTypeDescription: aws.String(desc),
			},
		})
		if err != nil {
			return fmt.Errorf("CreateThingType failed: %w", err)
		}
		if out.ThingTypeName == nil || *out.ThingTypeName != ttName {
			return fmt.Errorf("expected thingTypeName=%s, got %v", ttName, out.ThingTypeName)
		}
		if out.ThingTypeArn == nil || *out.ThingTypeArn == "" {
			return fmt.Errorf("expected non-empty thingTypeArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_DescribeThingType", func() error {
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(ttName)})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeName == nil || *out.ThingTypeName != ttName {
			return fmt.Errorf("expected thingTypeName=%s, got %v", ttName, out.ThingTypeName)
		}
		if out.ThingTypeProperties == nil {
			return fmt.Errorf("expected thingTypeProperties to be non-nil")
		}
		if out.ThingTypeProperties.ThingTypeDescription == nil || *out.ThingTypeProperties.ThingTypeDescription != desc {
			return fmt.Errorf("expected thingTypeDescription=%q, got %v", desc, out.ThingTypeProperties.ThingTypeDescription)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_ListThingTypes_IncludesCreated", func() error {
		out, err := tc.client.ListThingTypes(tc.ctx, &iot.ListThingTypesInput{})
		if err != nil {
			return fmt.Errorf("ListThingTypes failed: %w", err)
		}
		for _, t := range out.ThingTypes {
			if t.ThingTypeName != nil && *t.ThingTypeName == ttName {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d thing types", ttName, len(out.ThingTypes))
	}))

	results = append(results, r.RunTest("iot", "ThingType_DescribeThingType_NotFound", func() error {
		_, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(uniqueName("nope-type"))})
		return expectNotFound(err)
	}))

	// ── Searchable attributes update ──
	ttSearch := uniqueName("thing-type-search")
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttSearch)})

	results = append(results, r.RunTest("iot", "ThingType_SearchableAttributes", func() error {
		if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String(ttSearch),
			ThingTypeProperties: &types.ThingTypeProperties{
				SearchableAttributes: []string{"model", "serial"},
			},
		}); err != nil {
			return fmt.Errorf("CreateThingType failed: %w", err)
		}
		d, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(ttSearch)})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if d.ThingTypeProperties == nil || len(d.ThingTypeProperties.SearchableAttributes) != 2 {
			return fmt.Errorf("expected 2 searchableAttributes, got %v", d.ThingTypeProperties)
		}
		if _, err := tc.client.UpdateThingType(tc.ctx, &iot.UpdateThingTypeInput{
			ThingTypeName: aws.String(ttSearch),
			ThingTypeProperties: &types.ThingTypeProperties{
				SearchableAttributes: []string{"location"},
			},
		}); err != nil {
			return fmt.Errorf("UpdateThingType failed: %w", err)
		}
		d2, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(ttSearch)})
		if err != nil {
			return fmt.Errorf("DescribeThingType after update failed: %w", err)
		}
		if len(d2.ThingTypeProperties.SearchableAttributes) != 1 || d2.ThingTypeProperties.SearchableAttributes[0] != "location" {
			return fmt.Errorf("expected single searchableAttribute=location, got %v", d2.ThingTypeProperties.SearchableAttributes)
		}
		return nil
	}))

	// ── Deprecate / UndoDeprecate lifecycle ──
	// The type is created in the first test and deleted at group end, because
	// the UndoDeprecate test below reuses it.
	ttDep := uniqueName("thing-type-dep")
	var cleanupDep func()
	defer func() {
		if cleanupDep != nil {
			cleanupDep()
		}
	}()

	results = append(results, r.RunTest("iot", "ThingType_Deprecate_Persist", func() error {
		cleanup, err := tc.createThingType(ttDep)
		if err != nil {
			return fmt.Errorf("CreateThingType failed: %w", err)
		}
		cleanupDep = cleanup
		if _, err := tc.client.DeprecateThingType(tc.ctx, &iot.DeprecateThingTypeInput{ThingTypeName: aws.String(ttDep)}); err != nil {
			return fmt.Errorf("DeprecateThingType failed: %w", err)
		}
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(ttDep)})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeMetadata == nil || !out.ThingTypeMetadata.Deprecated {
			return fmt.Errorf("expected deprecated=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_UndoDeprecate", func() error {
		if _, err := tc.client.DeprecateThingType(tc.ctx, &iot.DeprecateThingTypeInput{
			ThingTypeName: aws.String(ttDep),
			UndoDeprecate: true,
		}); err != nil {
			return fmt.Errorf("UndoDeprecate failed: %w", err)
		}
		out, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{ThingTypeName: aws.String(ttDep)})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if out.ThingTypeMetadata == nil || out.ThingTypeMetadata.Deprecated {
			return fmt.Errorf("expected deprecated=false after undo")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ThingType_DeleteThingType", func() error {
		_, err := tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttName)})
		return err
	}))

	return results
}
