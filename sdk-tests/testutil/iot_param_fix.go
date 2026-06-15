package testutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTParamFixTests exercises the nested-parameter unwrapping fixes
// (unwrapProps helper) and related parameter-handling corrections.
func (r *TestRunner) runIoTParamFixTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	// ── ThingType with description (unwrapProps for thingTypeProperties) ──
	results = append(results, r.RunTest("iot", "ThingType_WithDescription", func() error {
		desc := "param-fix: thing type with description"
		out, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String("pf-thing-type-desc"),
			ThingTypeProperties: &types.ThingTypeProperties{
				ThingTypeDescription: aws.String(desc),
			},
		})
		if err != nil {
			return fmt.Errorf("CreateThingType with description failed: %w", err)
		}
		if out.ThingTypeName == nil || *out.ThingTypeName != "pf-thing-type-desc" {
			return fmt.Errorf("expected thingTypeName=pf-thing-type-desc, got %v", out.ThingTypeName)
		}

		// DescribeThingType should return the description via thingTypeProperties
		d, err := tc.client.DescribeThingType(tc.ctx, &iot.DescribeThingTypeInput{
			ThingTypeName: aws.String("pf-thing-type-desc"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingType failed: %w", err)
		}
		if d.ThingTypeProperties == nil {
			return fmt.Errorf("expected thingTypeProperties to be non-nil")
		}
		if d.ThingTypeProperties.ThingTypeDescription == nil || *d.ThingTypeProperties.ThingTypeDescription != desc {
			return fmt.Errorf("expected thingTypeDescription=%q, got %v", desc, d.ThingTypeProperties.ThingTypeDescription)
		}
		return nil
	}))

	// ── ThingGroup with description and attributes (unwrapProps for thingGroupProperties) ──
	results = append(results, r.RunTest("iot", "ThingGroup_WithDescriptionAndAttributes", func() error {
		desc := "param-fix: thing group desc"
		out, err := tc.client.CreateThingGroup(tc.ctx, &iot.CreateThingGroupInput{
			ThingGroupName: aws.String("pf-thing-group-desc"),
			ThingGroupProperties: &types.ThingGroupProperties{
				ThingGroupDescription: aws.String(desc),
				AttributePayload: &types.AttributePayload{
					Attributes: map[string]string{
						"env": "test",
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateThingGroup with description failed: %w", err)
		}
		if out.ThingGroupName == nil || *out.ThingGroupName != "pf-thing-group-desc" {
			return fmt.Errorf("expected thingGroupName=pf-thing-group-desc, got %v", out.ThingGroupName)
		}

		// DescribeThingGroup should return description and attributes
		d, err := tc.client.DescribeThingGroup(tc.ctx, &iot.DescribeThingGroupInput{
			ThingGroupName: aws.String("pf-thing-group-desc"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThingGroup failed: %w", err)
		}
		if d.ThingGroupProperties == nil {
			return fmt.Errorf("expected thingGroupProperties to be non-nil")
		}
		if d.ThingGroupProperties.ThingGroupDescription == nil || *d.ThingGroupProperties.ThingGroupDescription != desc {
			return fmt.Errorf("expected thingGroupDescription=%q, got %v", desc, d.ThingGroupProperties.ThingGroupDescription)
		}
		if d.ThingGroupProperties.AttributePayload == nil {
			return fmt.Errorf("expected attributePayload to be non-nil")
		}
		if d.ThingGroupProperties.AttributePayload.Attributes == nil {
			return fmt.Errorf("expected attributes map to be non-nil")
		}
		if v, ok := d.ThingGroupProperties.AttributePayload.Attributes["env"]; !ok || v != "test" {
			return fmt.Errorf("expected env=test, got %v", d.ThingGroupProperties.AttributePayload.Attributes["env"])
		}
		return nil
	}))

	// ── BillingGroup with description (unwrapProps for billingGroupProperties) ──
	results = append(results, r.RunTest("iot", "BillingGroup_WithDescription", func() error {
		desc := "param-fix: billing group desc"
		out, err := tc.client.CreateBillingGroup(tc.ctx, &iot.CreateBillingGroupInput{
			BillingGroupName: aws.String("pf-billing-group-desc"),
			BillingGroupProperties: &types.BillingGroupProperties{
				BillingGroupDescription: aws.String(desc),
			},
		})
		if err != nil {
			return fmt.Errorf("CreateBillingGroup with description failed: %w", err)
		}
		if out.BillingGroupName == nil || *out.BillingGroupName != "pf-billing-group-desc" {
			return fmt.Errorf("expected billingGroupName=pf-billing-group-desc, got %v", out.BillingGroupName)
		}

		// DescribeBillingGroup should return description
		d, err := tc.client.DescribeBillingGroup(tc.ctx, &iot.DescribeBillingGroupInput{
			BillingGroupName: aws.String("pf-billing-group-desc"),
		})
		if err != nil {
			return fmt.Errorf("DescribeBillingGroup failed: %w", err)
		}
		if d.BillingGroupProperties == nil {
			return fmt.Errorf("expected billingGroupProperties to be non-nil")
		}
		if d.BillingGroupProperties.BillingGroupDescription == nil || *d.BillingGroupProperties.BillingGroupDescription != desc {
			return fmt.Errorf("expected billingGroupDescription=%q, got %v", desc, d.BillingGroupProperties.BillingGroupDescription)
		}
		return nil
	}))

	// ── UpdateThing removeThingType (boolean type handling) ──
	results = append(results, r.RunTest("iot", "UpdateThing_RemoveThingType", func() error {
		// Create a thing type first
		_, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String("pf-remove-tt"),
		})
		if err != nil {
			var httpErr interface{ HTTPStatusCode() int }
			if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() == 409 {
				// Already exists from previous run
			} else {
				return fmt.Errorf("CreateThingType failed: %w", err)
			}
		}

		// Create thing with type
		_, err = tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName:    aws.String("pf-remove-thing"),
			ThingTypeName: aws.String("pf-remove-tt"),
		})
		if err != nil {
			var httpErr interface{ HTTPStatusCode() int }
			if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() == 409 {
				// Already exists — delete and recreate
				tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String("pf-remove-thing")})
				_, err = tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
					ThingName:    aws.String("pf-remove-thing"),
					ThingTypeName: aws.String("pf-remove-tt"),
				})
				if err != nil {
					return fmt.Errorf("CreateThing (retry) failed: %w", err)
				}
			} else {
				return fmt.Errorf("CreateThing failed: %w", err)
			}
		}

		// Verify thing has type
		desc1, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{
			ThingName: aws.String("pf-remove-thing"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThing (before remove) failed: %w", err)
		}
		if desc1.ThingTypeName == nil || *desc1.ThingTypeName != "pf-remove-tt" {
			return fmt.Errorf("expected thingTypeName=pf-remove-tt before removal, got %v", desc1.ThingTypeName)
		}

		// Remove thing type
		_, err = tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String("pf-remove-thing"),
			RemoveThingType: true,
		})
		if err != nil {
			return fmt.Errorf("UpdateThing with RemoveThingType=true failed: %w", err)
		}

		// Verify type was removed
		desc2, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{
			ThingName: aws.String("pf-remove-thing"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThing (after remove) failed: %w", err)
		}
		if desc2.ThingTypeName != nil && *desc2.ThingTypeName != "" {
			return fmt.Errorf("expected thingTypeName to be empty after RemoveThingType, got %v", desc2.ThingTypeName)
		}
		return nil
	}))

	// ── UpdateThing nonexistent thingTypeName (existence check) ──
	results = append(results, r.RunTest("iot", "UpdateThing_NonexistentThingTypeError", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:    aws.String("pf-remove-thing"),
			ThingTypeName: aws.String("pf-nonexistent-type-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for nonexistent thingTypeName")
		}
		if !strings.Contains(err.Error(), "ResourceNotFound") && !strings.Contains(err.Error(), "NotFound") {
			return fmt.Errorf("expected ResourceNotFound error, got: %v", err)
		}
		return nil
	}))

	// ── UpdateThing: removeThingType false + thingTypeName (BUG-B) ──
	results = append(results, r.RunTest("iot", "UpdateThing_RemoveFalse_ThenSetTypeName", func() error {
		_, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String("pf-rf-tt"),
		})
		if err != nil {
			var httpErr interface{ HTTPStatusCode() int }
			if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() == 409 {
			} else {
				return fmt.Errorf("CreateThingType failed: %w", err)
			}
		}

		_, err = tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName: aws.String("pf-rf-thing"),
		})
		if err != nil {
			var httpErr interface{ HTTPStatusCode() int }
			if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() == 409 {
				tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String("pf-rf-thing")})
				_, err = tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String("pf-rf-thing")})
				if err != nil {
					return fmt.Errorf("CreateThing (retry) failed: %w", err)
				}
			} else {
				return fmt.Errorf("CreateThing failed: %w", err)
			}
		}

		_, err = tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String("pf-rf-thing"),
			RemoveThingType: false,
			ThingTypeName:   aws.String("pf-rf-tt"),
		})
		if err != nil {
			return fmt.Errorf("UpdateThing with RemoveThingType=false + ThingTypeName failed: %w", err)
		}

		desc, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{
			ThingName: aws.String("pf-rf-thing"),
		})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if desc.ThingTypeName == nil || *desc.ThingTypeName != "pf-rf-tt" {
			return fmt.Errorf("expected thingTypeName=pf-rf-tt, got %v", desc.ThingTypeName)
		}
		return nil
	}))

	// ── Cleanup ──
	results = append(results, r.RunTest("iot", "ParamFix_Cleanup", func() error {
		// Thing
		tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String("pf-remove-thing")})
		tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String("pf-rf-thing")})
		// ThingType
		tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String("pf-remove-tt")})
		tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String("pf-rf-tt")})
		tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String("pf-thing-type-desc")})
		// ThingGroup
		tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{ThingGroupName: aws.String("pf-thing-group-desc")})
		// BillingGroup
		tc.client.DeleteBillingGroup(tc.ctx, &iot.DeleteBillingGroupInput{BillingGroupName: aws.String("pf-billing-group-desc")})
		return nil
	}))

	return results
}
