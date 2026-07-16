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
		_, err := tc.client.UpdateDomainConfiguration(tc.ctx, &iot.UpdateDomainConfigurationInput{
			DomainConfigurationName: aws.String(domainName),
			AuthorizerConfig: &types.AuthorizerConfig{
				DefaultAuthorizerName: aws.String(uniqueName("default-auth")),
			},
		})
		if err != nil {
			return fmt.Errorf("UpdateDomainConfiguration failed: %w", err)
		}
		return nil
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
