package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTDomainConfigTests covers the DomainConfiguration lifecycle
// (Create/Describe/Update/List/Delete) with real assertions and a NotFound
// negative path.
func (r *TestRunner) runIoTDomainConfigTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	domainName := uniqueName("domain")

	defer tc.client.DeleteDomainConfiguration(tc.ctx, &iot.DeleteDomainConfigurationInput{DomainConfigurationName: aws.String(domainName)})

	results = append(results, r.RunTest("iot", "DomainConfig_CreateDomainConfiguration", func() error {
		out, err := tc.client.CreateDomainConfiguration(tc.ctx, &iot.CreateDomainConfigurationInput{
			DomainConfigurationName: aws.String(domainName),
			DomainName:              aws.String(domainName + ".example.com"),
			ServiceType:             "DATA",
			Tags: []types.Tag{
				{Key: aws.String("purpose"), Value: aws.String("sdk-test")},
			},
		})
		if err != nil {
			return fmt.Errorf("CreateDomainConfiguration failed: %w", err)
		}
		if out.DomainConfigurationName == nil || *out.DomainConfigurationName != domainName {
			return fmt.Errorf("expected domainConfigurationName=%s, got %v", domainName, out.DomainConfigurationName)
		}
		if out.DomainConfigurationArn == nil || *out.DomainConfigurationArn == "" {
			return fmt.Errorf("expected non-empty domainConfigurationArn")
		}
		// Create-time tags must be visible through ListTagsForResource.
		found, err := tc.resourceHasTag(out.DomainConfigurationArn, "purpose", "sdk-test")
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", *out.DomainConfigurationArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_DescribeDomainConfiguration", func() error {
		out, err := tc.client.DescribeDomainConfiguration(tc.ctx, &iot.DescribeDomainConfigurationInput{DomainConfigurationName: aws.String(domainName)})
		if err != nil {
			return fmt.Errorf("DescribeDomainConfiguration failed: %w", err)
		}
		if out.DomainConfigurationName == nil || *out.DomainConfigurationName != domainName {
			return fmt.Errorf("expected domainConfigurationName=%s, got %v", domainName, out.DomainConfigurationName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_UpdateDomainConfiguration", func() error {
		authName := uniqueName("default-auth")
		out, err := tc.client.UpdateDomainConfiguration(tc.ctx, &iot.UpdateDomainConfigurationInput{
			DomainConfigurationName: aws.String(domainName),
			AuthorizerConfig: &types.AuthorizerConfig{
				DefaultAuthorizerName:   aws.String(authName),
				AllowAuthorizerOverride: aws.Bool(true),
			},
			DomainConfigurationStatus: types.DomainConfigurationStatusDisabled,
		})
		if err != nil {
			return fmt.Errorf("UpdateDomainConfiguration failed: %w", err)
		}
		if aws.ToString(out.DomainConfigurationName) != domainName {
			return fmt.Errorf("expected domainConfigurationName=%s in update response", domainName)
		}
		if aws.ToString(out.DomainConfigurationArn) == "" {
			return fmt.Errorf("expected non-empty domainConfigurationArn in update response")
		}
		desc, err := tc.client.DescribeDomainConfiguration(tc.ctx, &iot.DescribeDomainConfigurationInput{DomainConfigurationName: aws.String(domainName)})
		if err != nil {
			return fmt.Errorf("DescribeDomainConfiguration after update failed: %w", err)
		}
		if desc.AuthorizerConfig == nil {
			return fmt.Errorf("expected authorizerConfig persisted by update")
		}
		if aws.ToString(desc.AuthorizerConfig.DefaultAuthorizerName) != authName {
			return fmt.Errorf("expected defaultAuthorizerName=%s, got %v", authName, desc.AuthorizerConfig.DefaultAuthorizerName)
		}
		if aws.ToBool(desc.AuthorizerConfig.AllowAuthorizerOverride) != true {
			return fmt.Errorf("expected allowAuthorizerOverride=true, got %v", desc.AuthorizerConfig.AllowAuthorizerOverride)
		}
		if desc.DomainConfigurationStatus != types.DomainConfigurationStatusDisabled {
			return fmt.Errorf("expected domainConfigurationStatus=DISABLED, got %v", desc.DomainConfigurationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_CreateDomainConfiguration_InvalidServiceTypeRejected", func() error {
		_, err := tc.client.CreateDomainConfiguration(tc.ctx, &iot.CreateDomainConfigurationInput{
			DomainConfigurationName: aws.String(uniqueName("domain-bad-type")),
			ServiceType:             "GARBAGE",
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_MultipleServerCertificateArnsRejected", func() error {
		_, err := tc.client.CreateDomainConfiguration(tc.ctx, &iot.CreateDomainConfigurationInput{
			DomainConfigurationName: aws.String(uniqueName("domain-two-certs")),
			ServiceType:             "DATA",
			ServerCertificateArns: []string{
				"arn:aws:acm:us-east-1:123456789012:certificate/first",
				"arn:aws:acm:us-east-1:123456789012:certificate/second",
			},
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_ListDomainConfigurations_IncludesCreated", func() error {
		out, err := tc.client.ListDomainConfigurations(tc.ctx, &iot.ListDomainConfigurationsInput{})
		if err != nil {
			return fmt.Errorf("ListDomainConfigurations failed: %w", err)
		}
		for _, d := range out.DomainConfigurations {
			if d.DomainConfigurationName != nil && *d.DomainConfigurationName == domainName {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d domain configs", domainName, len(out.DomainConfigurations))
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_DescribeDomainConfiguration_NotFound", func() error {
		_, err := tc.client.DescribeDomainConfiguration(tc.ctx, &iot.DescribeDomainConfigurationInput{DomainConfigurationName: aws.String(uniqueName("nope-domain"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "DomainConfig_DeleteDomainConfiguration", func() error {
		_, err := tc.client.DeleteDomainConfiguration(tc.ctx, &iot.DeleteDomainConfigurationInput{DomainConfigurationName: aws.String(domainName)})
		return err
	}))

	return results
}
