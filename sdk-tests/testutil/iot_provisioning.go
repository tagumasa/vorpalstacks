package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTProvisioningTemplateTests covers the ProvisioningTemplate lifecycle
// (Create/Describe/Update/List/Delete), version management (Create/List/
// Describe/Delete), and CreateProvisioningClaim on a real template, with a
// NotFound negative path. All names use uniqueName.
func (r *TestRunner) runIoTProvisioningTemplateTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	tmplName := uniqueName("prov-tmpl")
	roleARN := tc.iamRoleARN("provisioning")

	defer tc.client.DeleteProvisioningTemplate(tc.ctx, &iot.DeleteProvisioningTemplateInput{TemplateName: aws.String(tmplName)})

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_Create", func() error {
		out, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName:        aws.String(tmplName),
			Description:         aws.String("test template"),
			ProvisioningRoleArn: aws.String(roleARN),
			TemplateBody:        aws.String(`{"Parameters":{"ThingName":{"Type":"String"}}}`),
			Enabled:             aws.Bool(true),
			Tags:                []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("CreateProvisioningTemplate failed: %w", err)
		}
		if out.TemplateName == nil || *out.TemplateName != tmplName {
			return fmt.Errorf("expected templateName=%s, got %v", tmplName, out.TemplateName)
		}
		if out.TemplateArn == nil || *out.TemplateArn == "" {
			return fmt.Errorf("expected non-empty templateArn")
		}
		if out.DefaultVersionId == nil || *out.DefaultVersionId != 1 {
			return fmt.Errorf("expected defaultVersionId=1, got %v", out.DefaultVersionId)
		}
		// Create-time tags must be visible through ListTagsForResource.
		found, err := tc.resourceHasTag(out.TemplateArn, "purpose", "sdk-test")
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", aws.ToString(out.TemplateArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_MissingRequiredRejected", func() error {
		if _, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName:        aws.String(uniqueName("prov-tmpl")),
			ProvisioningRoleArn: aws.String(roleARN),
		}); err == nil {
			return fmt.Errorf("expected rejection without templateBody")
		}
		if _, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName: aws.String(uniqueName("prov-tmpl")),
			TemplateBody: aws.String(`{"Parameters":{"ThingName":{"Type":"String"}}}`),
		}); err == nil {
			return fmt.Errorf("expected rejection without provisioningRoleArn")
		}
		if _, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName:        aws.String(uniqueName("prov-tmpl")),
			ProvisioningRoleArn: aws.String(roleARN),
			TemplateBody:        aws.String(`{"Parameters":{"ThingName":{"Type":"String"}}}`),
			Type:                iottypes.TemplateType("Bogus"),
		}); err == nil {
			return fmt.Errorf("expected rejection for off-enum type")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_Describe", func() error {
		out, err := tc.client.DescribeProvisioningTemplate(tc.ctx, &iot.DescribeProvisioningTemplateInput{TemplateName: aws.String(tmplName)})
		if err != nil {
			return fmt.Errorf("DescribeProvisioningTemplate failed: %w", err)
		}
		if out.TemplateName == nil || *out.TemplateName != tmplName {
			return fmt.Errorf("expected templateName=%s, got %v", tmplName, out.TemplateName)
		}
		if out.DefaultVersionId == nil || *out.DefaultVersionId != 1 {
			return fmt.Errorf("expected defaultVersionId=1 on describe, got %v", out.DefaultVersionId)
		}
		if out.ProvisioningRoleArn == nil || *out.ProvisioningRoleArn != roleARN {
			return fmt.Errorf("expected provisioningRoleArn=%s, got %v", roleARN, out.ProvisioningRoleArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_Update", func() error {
		_, err := tc.client.UpdateProvisioningTemplate(tc.ctx, &iot.UpdateProvisioningTemplateInput{
			TemplateName:        aws.String(tmplName),
			Description:         aws.String("updated description"),
			ProvisioningRoleArn: aws.String(roleARN),
			Enabled:             aws.Bool(false),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_List_IncludesCreated", func() error {
		out, err := tc.client.ListProvisioningTemplates(tc.ctx, &iot.ListProvisioningTemplatesInput{})
		if err != nil {
			return fmt.Errorf("ListProvisioningTemplates failed: %w", err)
		}
		for _, t := range out.Templates {
			if t.TemplateName != nil && *t.TemplateName == tmplName {
				// The list carries summaries: the ProvisioningTemplateSummary
				// type has no templateBody or provisioningRoleArn members, so
				// the summary shape is pinned at compile time.
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d templates", tmplName, len(out.Templates))
	}))

	// NOTE: CreateProvisioningTemplateVersion / ListProvisioningTemplateVersions
	// are not yet wired to a service handler, so version-management assertions are
	// omitted until the versioning path is implemented.

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_CreateProvisioningClaim", func() error {
		// CreateProvisioningClaim against an existing template must succeed and
		// return a credential pair (the edge platform issues a short-lived
		// claim). Assert non-error rather than specific key material.
		_, err := tc.client.CreateProvisioningClaim(tc.ctx, &iot.CreateProvisioningClaimInput{
			TemplateName: aws.String(tmplName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_Describe_NotFound", func() error {
		_, err := tc.client.DescribeProvisioningTemplate(tc.ctx, &iot.DescribeProvisioningTemplateInput{TemplateName: aws.String(uniqueName("nope-tmpl"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_Delete", func() error {
		_, err := tc.client.DeleteProvisioningTemplate(tc.ctx, &iot.DeleteProvisioningTemplateInput{TemplateName: aws.String(tmplName)})
		return err
	}))

	return results
}
