package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTThingTests covers the Thing resource lifecycle (Create/Describe/Update/
// List/Delete) plus negative paths (Duplicate, NotFound, RemoveThingType) and
// the StartThingRegistrationTask bulk path. All resources use uniqueName so the
// suite is safe under parallel runs and re-runs.
func (r *TestRunner) runIoTThingTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("thing")
	attrs := map[string]string{"version": "1.0", "location": tc.region}

	// Best-effort cleanup so a failed run never leaves the thing behind.
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	results = append(results, r.RunTest("iot", "Thing_CreateThing", func() error {
		out, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName:        aws.String(thingName),
			AttributePayload: &types.AttributePayload{Attributes: attrs},
		})
		if err != nil {
			return fmt.Errorf("CreateThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != thingName {
			return fmt.Errorf("expected thingName=%s, got %v", thingName, out.ThingName)
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
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if out.ThingName == nil || *out.ThingName != thingName {
			return fmt.Errorf("expected thingName=%s, got %v", thingName, out.ThingName)
		}
		if out.Attributes == nil {
			return fmt.Errorf("expected attributes to be non-nil")
		}
		if v, ok := out.Attributes["version"]; !ok || v != "1.0" {
			return fmt.Errorf("expected attribute version=1.0, got %v", out.Attributes["version"])
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_UpdateThing", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName: aws.String(thingName),
			AttributePayload: &types.AttributePayload{
				Attributes: map[string]string{"version": "2.0"},
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateThing failed: %w", err)
		}
		out, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DescribeThing after update failed: %w", err)
		}
		if v, ok := out.Attributes["version"]; !ok || v != "2.0" {
			return fmt.Errorf("expected attribute version=2.0 after update, got %v", out.Attributes["version"])
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_ListThings_IncludesCreated", func() error {
		found, err := tc.thingExists(thingName)
		if err != nil {
			return fmt.Errorf("ListThings failed: %w", err)
		}
		if !found {
			return fmt.Errorf("%s not found in thing list", thingName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_CreateThing_Duplicate_Conflict", func() error {
		_, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)})
		return expectConflict(err)
	}))

	results = append(results, r.RunTest("iot", "Thing_DescribeThing_NotFound", func() error {
		_, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(uniqueName("nope"))})
		return expectNotFound(err)
	}))

	// RemoveThingType: create a thing type, attach it to a fresh thing, then
	// remove it and verify DescribeThing no longer reports the type.
	ttName := uniqueName("thing-type")
	thingWithType := uniqueName("thing-tt")
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingWithType)})
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttName)})

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_RemoveThingType", func() error {
		if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{ThingTypeName: aws.String(ttName)}); err != nil {
			return fmt.Errorf("CreateThingType prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{
			ThingName:     aws.String(thingWithType),
			ThingTypeName: aws.String(ttName),
		}); err != nil {
			return fmt.Errorf("CreateThing with type failed: %w", err)
		}
		before, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingWithType)})
		if err != nil {
			return fmt.Errorf("DescribeThing before remove failed: %w", err)
		}
		if before.ThingTypeName == nil || *before.ThingTypeName != ttName {
			return fmt.Errorf("expected thingTypeName=%s before removal, got %v", ttName, before.ThingTypeName)
		}
		if _, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String(thingWithType),
			RemoveThingType: true,
		}); err != nil {
			return fmt.Errorf("UpdateThing RemoveThingType failed: %w", err)
		}
		after, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingWithType)})
		if err != nil {
			return fmt.Errorf("DescribeThing after remove failed: %w", err)
		}
		if after.ThingTypeName != nil && *after.ThingTypeName != "" {
			return fmt.Errorf("expected empty thingTypeName after removal, got %v", after.ThingTypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_NonexistentThingType_NotFound", func() error {
		_, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:     aws.String(thingName),
			ThingTypeName: aws.String(uniqueName("nope-type")),
		})
		return expectNotFound(err)
	}))

	// Set a thing type on a previously typeless thing (RemoveThingType=false +
	// ThingTypeName path).
	thingForSet := uniqueName("thing-set")
	ttForSet := uniqueName("thing-type-set")
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingForSet)})
	defer tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(ttForSet)})

	results = append(results, r.RunTest("iot", "Thing_UpdateThing_SetThingType", func() error {
		if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{ThingTypeName: aws.String(ttForSet)}); err != nil {
			return fmt.Errorf("CreateThingType prerequisite failed: %w", err)
		}
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingForSet)}); err != nil {
			return fmt.Errorf("CreateThing prerequisite failed: %w", err)
		}
		if _, err := tc.client.UpdateThing(tc.ctx, &iot.UpdateThingInput{
			ThingName:       aws.String(thingForSet),
			RemoveThingType: false,
			ThingTypeName:   aws.String(ttForSet),
		}); err != nil {
			return fmt.Errorf("UpdateThing set ThingTypeName failed: %w", err)
		}
		desc, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(thingForSet)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if desc.ThingTypeName == nil || *desc.ThingTypeName != ttForSet {
			return fmt.Errorf("expected thingTypeName=%s, got %v", ttForSet, desc.ThingTypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_StartThingRegistrationTask_Validation", func() error {
		// StartThingRegistrationTask requires a registered ProvisioningTemplate
		// and a valid template body; a minimal body must be rejected with a
		// validation error, proving the handler is registered and validates.
		_, err := tc.client.StartThingRegistrationTask(tc.ctx, &iot.StartThingRegistrationTaskInput{
			TemplateBody: aws.String("{}"),
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "Thing_AttachDetachThingPrincipal", func() error {
		// Create a certificate to use as the principal.
		cert, err := tc.client.CreateKeysAndCertificate(tc.ctx, &iot.CreateKeysAndCertificateInput{SetAsActive: true})
		if err != nil {
			return fmt.Errorf("CreateKeysAndCertificate failed: %w", err)
		}
		certARN := aws.ToString(cert.CertificateArn)
		defer tc.client.DetachThingPrincipal(tc.ctx, &iot.DetachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		})
		defer tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: cert.CertificateId, NewStatus: types.CertificateStatusInactive,
		})
		defer tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: cert.CertificateId})

		if _, err := tc.client.AttachThingPrincipal(tc.ctx, &iot.AttachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("AttachThingPrincipal failed: %w", err)
		}
		// The principal must appear in ListThingPrincipals.
		lp, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals failed: %w", err)
		}
		found := false
		for _, p := range lp.Principals {
			if p == certARN {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("attached principal %s not found in ListThingPrincipals", certARN)
		}
		if _, err := tc.client.DetachThingPrincipal(tc.ctx, &iot.DetachThingPrincipalInput{
			ThingName: aws.String(thingName), Principal: aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("DetachThingPrincipal failed: %w", err)
		}
		lp2, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals after detach failed: %w", err)
		}
		for _, p := range lp2.Principals {
			if p == certARN {
				return fmt.Errorf("principal still attached after DetachThingPrincipal")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DeleteThing", func() error {
		_, err := tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})
		if err != nil {
			return fmt.Errorf("DeleteThing failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Thing_DeleteThing_NotFound", func() error {
		_, err := tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(uniqueName("nope"))})
		return expectNotFound(err)
	}))

	return results
}
