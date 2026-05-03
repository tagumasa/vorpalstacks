package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamSAMLTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	samlProviderName := fmt.Sprintf("TestSAML-%s", tc.ts)

	results = append(results, r.RunTest("iam", "CreateSAMLProvider", func() error {
		resp, err := tc.client.CreateSAMLProvider(tc.ctx, &iam.CreateSAMLProviderInput{
			Name:                 aws.String(samlProviderName),
			SAMLMetadataDocument: aws.String(samlMetadata),
			Tags: []types.Tag{
				{Key: aws.String("Source"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		tc.samlProviderArn = *resp.SAMLProviderArn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetSAMLProvider", func() error {
		resp, err := tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if resp.SAMLMetadataDocument == nil || *resp.SAMLMetadataDocument == "" {
			return fmt.Errorf("saml metadata document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviders", func() error {
		resp, err := tc.client.ListSAMLProviders(tc.ctx, &iam.ListSAMLProvidersInput{})
		if err != nil {
			return err
		}
		if resp.SAMLProviderList == nil {
			return fmt.Errorf("saml provider list is nil")
		}
		found := false
		for _, p := range resp.SAMLProviderList {
			if aws.ToString(p.Arn) == tc.samlProviderArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("saml provider %s not found in list", tc.samlProviderArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateSAMLProvider", func() error {
		resp, err := tc.client.UpdateSAMLProvider(tc.ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn:      aws.String(tc.samlProviderArn),
			SAMLMetadataDocument: aws.String(samlMetadata),
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		return nil
	}))

	// SAML tags
	results = append(results, r.RunTest("iam", "TagSAMLProvider", func() error {
		_, err := tc.client.TagSAMLProvider(tc.ctx, &iam.TagSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return fmt.Errorf("ListSAMLProviderTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagSAMLProvider")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviderTags", func() error {
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		if !iamTagPresent(resp.Tags, "Source", "test") {
			return fmt.Errorf("Source=test tag not found (from create)")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagSAMLProvider", func() error {
		_, err := tc.client.UntagSAMLProvider(tc.ctx, &iam.UntagSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
			TagKeys:         []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListSAMLProviderTags(tc.ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment tag should be removed")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteSAMLProvider", func() error {
		_, err := tc.client.DeleteSAMLProvider(tc.ctx, &iam.DeleteSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetSAMLProvider(tc.ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(tc.samlProviderArn),
		})
		if err == nil {
			return fmt.Errorf("GetSAMLProvider should fail after DeleteSAMLProvider")
		}
		return nil
	}))

	return results
}
