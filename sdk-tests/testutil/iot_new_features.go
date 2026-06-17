package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func (r *TestRunner) runIoTAuthorizerTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Authorizer_DeleteAuthorizer", func() error {
		_, err := tc.client.DeleteAuthorizer(tc.ctx, &iot.DeleteAuthorizerInput{
			AuthorizerName: aws.String("test-authorizer-1"),
		})
		_ = err
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_CreateAuthorizer", func() error {
		out, err := tc.client.CreateAuthorizer(tc.ctx, &iot.CreateAuthorizerInput{
			AuthorizerName:        aws.String("test-authorizer-1"),
			AuthorizerFunctionArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:auth-fn"),
			TokenKeyName:          aws.String("token"),
			Status:                "ACTIVE",
		})
		if err != nil {
			return fmt.Errorf("CreateAuthorizer failed: %w", err)
		}
		if out.AuthorizerName == nil || *out.AuthorizerName != "test-authorizer-1" {
			return fmt.Errorf("expected authorizerName=test-authorizer-1, got %v", out.AuthorizerName)
		}
		if out.AuthorizerArn == nil || *out.AuthorizerArn == "" {
			return fmt.Errorf("expected non-empty authorizerArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DescribeAuthorizer", func() error {
		out, err := tc.client.DescribeAuthorizer(tc.ctx, &iot.DescribeAuthorizerInput{
			AuthorizerName: aws.String("test-authorizer-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeAuthorizer failed: %w", err)
		}
		if out.AuthorizerDescription == nil {
			return fmt.Errorf("expected non-nil authorizerDescription")
		}
		d := out.AuthorizerDescription
		if d.AuthorizerName == nil || *d.AuthorizerName != "test-authorizer-1" {
			return fmt.Errorf("expected authorizerName=test-authorizer-1, got %v", d.AuthorizerName)
		}
		if d.AuthorizerFunctionArn == nil || *d.AuthorizerFunctionArn != "arn:aws:lambda:us-east-1:000000000000:function:auth-fn" {
			return fmt.Errorf("expected function ARN, got %s", aws.ToString(d.AuthorizerFunctionArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_UpdateAuthorizer", func() error {
		_, err := tc.client.UpdateAuthorizer(tc.ctx, &iot.UpdateAuthorizerInput{
			AuthorizerName:        aws.String("test-authorizer-1"),
			AuthorizerFunctionArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:auth-fn-v2"),
		})
		if err != nil {
			return fmt.Errorf("UpdateAuthorizer failed: %w", err)
		}
		out, err := tc.client.DescribeAuthorizer(tc.ctx, &iot.DescribeAuthorizerInput{
			AuthorizerName: aws.String("test-authorizer-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeAuthorizer after update failed: %w", err)
		}
		if out.AuthorizerDescription.AuthorizerFunctionArn == nil || *out.AuthorizerDescription.AuthorizerFunctionArn != "arn:aws:lambda:us-east-1:000000000000:function:auth-fn-v2" {
			return fmt.Errorf("expected updated function ARN, got %v", out.AuthorizerDescription.AuthorizerFunctionArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_ListAuthorizers", func() error {
		out, err := tc.client.ListAuthorizers(tc.ctx, &iot.ListAuthorizersInput{})
		if err != nil {
			return fmt.Errorf("ListAuthorizers failed: %w", err)
		}
		if out.Authorizers == nil || len(out.Authorizers) == 0 {
			return fmt.Errorf("expected at least 1 authorizer, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DeleteAuthorizer", func() error {
		_, err := tc.client.DeleteAuthorizer(tc.ctx, &iot.DeleteAuthorizerInput{
			AuthorizerName: aws.String("test-authorizer-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteAuthorizer failed: %w", err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runIoTSecurityProfileTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "SecurityProfile_DeleteSecurityProfile", func() error {
		_, err := tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{
			SecurityProfileName: aws.String("test-sec-profile-1"),
		})
		_ = err
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_CreateSecurityProfile", func() error {
		out, err := tc.client.CreateSecurityProfile(tc.ctx, &iot.CreateSecurityProfileInput{
			SecurityProfileName:       aws.String("test-sec-profile-1"),
			SecurityProfileDescription: aws.String("test security profile"),
		})
		if err != nil {
			return fmt.Errorf("CreateSecurityProfile failed: %w", err)
		}
		if out.SecurityProfileName == nil || *out.SecurityProfileName != "test-sec-profile-1" {
			return fmt.Errorf("expected securityProfileName, got %v", out.SecurityProfileName)
		}
		if out.SecurityProfileArn == nil || *out.SecurityProfileArn == "" {
			return fmt.Errorf("expected non-empty securityProfileArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_DescribeSecurityProfile", func() error {
		out, err := tc.client.DescribeSecurityProfile(tc.ctx, &iot.DescribeSecurityProfileInput{
			SecurityProfileName: aws.String("test-sec-profile-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeSecurityProfile failed: %w", err)
		}
		if out.SecurityProfileName == nil || *out.SecurityProfileName != "test-sec-profile-1" {
			return fmt.Errorf("expected securityProfileName, got %v", out.SecurityProfileName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_ListSecurityProfiles", func() error {
		out, err := tc.client.ListSecurityProfiles(tc.ctx, &iot.ListSecurityProfilesInput{})
		if err != nil {
			return fmt.Errorf("ListSecurityProfiles failed: %w", err)
		}
		if out.SecurityProfileIdentifiers == nil || len(out.SecurityProfileIdentifiers) == 0 {
			return fmt.Errorf("expected at least 1 profile, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecurityProfile_DeleteSecurityProfile", func() error {
		_, err := tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{
			SecurityProfileName: aws.String("test-sec-profile-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteSecurityProfile failed: %w", err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runIoTDomainConfigTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "DomainConfig_CreateDomainConfiguration", func() error {
		out, err := tc.client.CreateDomainConfiguration(tc.ctx, &iot.CreateDomainConfigurationInput{
			DomainConfigurationName: aws.String("test-domain-1"),
			DomainName:              aws.String("test.iot.example.com"),
			ServiceType:             "DATA",
		})
		if err != nil {
			return fmt.Errorf("CreateDomainConfiguration failed: %w", err)
		}
		if out.DomainConfigurationName == nil || *out.DomainConfigurationName != "test-domain-1" {
			return fmt.Errorf("expected domainConfigurationName, got %v", out.DomainConfigurationName)
		}
		if out.DomainConfigurationArn == nil || *out.DomainConfigurationArn == "" {
			return fmt.Errorf("expected non-empty domainConfigurationArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_DescribeDomainConfiguration", func() error {
		out, err := tc.client.DescribeDomainConfiguration(tc.ctx, &iot.DescribeDomainConfigurationInput{
			DomainConfigurationName: aws.String("test-domain-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeDomainConfiguration failed: %w", err)
		}
		if out.DomainConfigurationName == nil || *out.DomainConfigurationName != "test-domain-1" {
			return fmt.Errorf("expected domainConfigurationName, got %v", out.DomainConfigurationName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_UpdateDomainConfiguration", func() error {
		_, err := tc.client.UpdateDomainConfiguration(tc.ctx, &iot.UpdateDomainConfigurationInput{
			DomainConfigurationName: aws.String("test-domain-1"),
			AuthorizerConfig: &types.AuthorizerConfig{
				DefaultAuthorizerName: aws.String("test-authorizer"),
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateDomainConfiguration failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_ListDomainConfigurations", func() error {
		out, err := tc.client.ListDomainConfigurations(tc.ctx, &iot.ListDomainConfigurationsInput{})
		if err != nil {
			return fmt.Errorf("ListDomainConfigurations failed: %w", err)
		}
		if out.DomainConfigurations == nil || len(out.DomainConfigurations) == 0 {
			return fmt.Errorf("expected at least 1 domain config, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_DeleteDomainConfiguration", func() error {
		_, err := tc.client.DeleteDomainConfiguration(tc.ctx, &iot.DeleteDomainConfigurationInput{
			DomainConfigurationName: aws.String("test-domain-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteDomainConfiguration failed: %w", err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runIoTProvisioningTemplateTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_CreateProvisioningTemplate", func() error {
		out, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName:       aws.String("test-prov-tmpl-1"),
			Description:         aws.String("test template"),
			ProvisioningRoleArn: aws.String("arn:aws:iam::000000000000:role/provisioning"),
			TemplateBody:        aws.String(`{"Parameters":{"ThingName":{"Type":"String"}}}`),
			Enabled:             aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("CreateProvisioningTemplate failed: %w", err)
		}
		if out.TemplateName == nil || *out.TemplateName != "test-prov-tmpl-1" {
			return fmt.Errorf("expected templateName, got %v", out.TemplateName)
		}
		if out.TemplateArn == nil || *out.TemplateArn == "" {
			return fmt.Errorf("expected non-empty templateArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_DescribeProvisioningTemplate", func() error {
		out, err := tc.client.DescribeProvisioningTemplate(tc.ctx, &iot.DescribeProvisioningTemplateInput{
			TemplateName: aws.String("test-prov-tmpl-1"),
		})
		if err != nil {
			return fmt.Errorf("DescribeProvisioningTemplate failed: %w", err)
		}
		if out.TemplateName == nil || *out.TemplateName != "test-prov-tmpl-1" {
			return fmt.Errorf("expected templateName, got %v", out.TemplateName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_UpdateProvisioningTemplate", func() error {
		_, err := tc.client.UpdateProvisioningTemplate(tc.ctx, &iot.UpdateProvisioningTemplateInput{
			TemplateName:       aws.String("test-prov-tmpl-1"),
			Description:         aws.String("updated description"),
			ProvisioningRoleArn: aws.String("arn:aws:iam::000000000000:role/provisioning-v2"),
			Enabled:             aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("UpdateProvisioningTemplate failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_ListProvisioningTemplates", func() error {
		out, err := tc.client.ListProvisioningTemplates(tc.ctx, &iot.ListProvisioningTemplatesInput{})
		if err != nil {
			return fmt.Errorf("ListProvisioningTemplates failed: %w", err)
		}
		if out.Templates == nil || len(out.Templates) == 0 {
			return fmt.Errorf("expected at least 1 template, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ProvisioningTemplate_DeleteProvisioningTemplate", func() error {
		_, err := tc.client.DeleteProvisioningTemplate(tc.ctx, &iot.DeleteProvisioningTemplateInput{
			TemplateName: aws.String("test-prov-tmpl-1"),
		})
		if err != nil {
			return fmt.Errorf("DeleteProvisioningTemplate failed: %w", err)
		}
		return nil
	}))

	return results
}
