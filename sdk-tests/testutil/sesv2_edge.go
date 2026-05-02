package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func (r *TestRunner) runSESv2EdgeTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("sesv2", "ListConfigurationSets_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		prefix := fmt.Sprintf("PagCS-%s-", pgTs)
		var pgConfigSets []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("%s%d", prefix, i)
			_, err := tc.client.CreateConfigurationSet(tc.ctx, &sesv2.CreateConfigurationSetInput{
				ConfigurationSetName: aws.String(name),
			})
			if err != nil {
				for _, cn := range pgConfigSets {
					tc.deleteConfigSet(cn)
				}
				return fmt.Errorf("create configuration set %s: %v", name, err)
			}
			pgConfigSets = append(pgConfigSets, name)
		}

		var allConfigSets []string
		var nextToken *string
		for {
			resp, err := tc.client.ListConfigurationSets(tc.ctx, &sesv2.ListConfigurationSetsInput{
				PageSize:  aws.Int32(2),
				NextToken: nextToken,
			})
			if err != nil {
				for _, cn := range pgConfigSets {
					tc.deleteConfigSet(cn)
				}
				return fmt.Errorf("list configuration sets page: %v", err)
			}
			for _, cs := range resp.ConfigurationSets {
				if strings.HasPrefix(cs, prefix) {
					allConfigSets = append(allConfigSets, cs)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, cn := range pgConfigSets {
			tc.deleteConfigSet(cn)
		}
		if len(allConfigSets) != 5 {
			return fmt.Errorf("expected 5 paginated configuration sets, got %d", len(allConfigSets))
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetConfigurationSet_NonExistent", func() error {
		_, err := tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String("nonexistent-cs-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSet_NonExistent", func() error {
		_, err := tc.client.DeleteConfigurationSet(tc.ctx, &sesv2.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String("nonexistent-cs-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentity_NonExistent", func() error {
		_, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String("nonexistent@example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_NonExistent", func() error {
		_, err := tc.client.DeleteEmailIdentity(tc.ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String("nonexistent@example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailTemplate_NonExistent", func() error {
		_, err := tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String("nonexistent-template-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailTemplate_NonExistent", func() error {
		_, err := tc.client.DeleteEmailTemplate(tc.ctx, &sesv2.DeleteEmailTemplateInput{
			TemplateName: aws.String("nonexistent-template-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetDedicatedIpPool_NonExistent", func() error {
		_, err := tc.client.GetDedicatedIpPool(tc.ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String("nonexistent-pool-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteDedicatedIpPool_NonExistent", func() error {
		_, err := tc.client.DeleteDedicatedIpPool(tc.ctx, &sesv2.DeleteDedicatedIpPoolInput{
			PoolName: aws.String("nonexistent-pool-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetContactList_NonExistent", func() error {
		_, err := tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String("nonexistent-cl-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetContact_NonExistent", func() error {
		_, err := tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
			ContactListName: aws.String("nonexistent-cl-xyz"),
			EmailAddress:    aws.String("nonexistent@example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContactList_NonExistent", func() error {
		_, err := tc.client.DeleteContactList(tc.ctx, &sesv2.DeleteContactListInput{
			ContactListName: aws.String("nonexistent-cl-xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetSuppressedDestination_NonExistent", func() error {
		_, err := tc.client.GetSuppressedDestination(tc.ctx, &sesv2.GetSuppressedDestinationInput{
			EmailAddress: aws.String("nonexistent@example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentityPolicies_NonExistent", func() error {
		_, err := tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String("nonexistent@example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
