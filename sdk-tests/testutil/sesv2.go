package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"vorpalstacks-sdk-tests/config"
)

type sesv2TestContext struct {
	client *sesv2.Client
	ctx    context.Context
	uid    int64
}

func newSESTestContext(endpoint, region string) (*sesv2TestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: endpoint,
		Region:   region,
	})
	if err != nil {
		return nil, fmt.Errorf("load config: %v", err)
	}
	return &sesv2TestContext{
		client: sesv2.NewFromConfig(cfg),
		ctx:    context.Background(),
		uid:    time.Now().UnixNano(),
	}, nil
}

func (tc *sesv2TestContext) configSetARN(name string) string {
	return fmt.Sprintf("arn:aws:ses:us-east-1:000000000000:configuration-set/%s", name)
}

func (tc *sesv2TestContext) deleteConfigSet(name string) {
	_, _ = tc.client.DeleteConfigurationSet(tc.ctx, &sesv2.DeleteConfigurationSetInput{
		ConfigurationSetName: aws.String(name),
	})
}

func (tc *sesv2TestContext) deleteEmailIdentity(identity string) {
	_, _ = tc.client.DeleteEmailIdentity(tc.ctx, &sesv2.DeleteEmailIdentityInput{
		EmailIdentity: aws.String(identity),
	})
}

func (tc *sesv2TestContext) deleteContactList(name string) {
	_, _ = tc.client.DeleteContactList(tc.ctx, &sesv2.DeleteContactListInput{
		ContactListName: aws.String(name),
	})
}

func (tc *sesv2TestContext) deleteDedicatedIpPool(name string) {
	_, _ = tc.client.DeleteDedicatedIpPool(tc.ctx, &sesv2.DeleteDedicatedIpPoolInput{
		PoolName: aws.String(name),
	})
}

func (tc *sesv2TestContext) deleteEmailTemplate(name string) {
	_, _ = tc.client.DeleteEmailTemplate(tc.ctx, &sesv2.DeleteEmailTemplateInput{
		TemplateName: aws.String(name),
	})
}

func (tc *sesv2TestContext) listAllConfigSets() ([]string, error) {
	var all []string
	var nextToken *string
	for {
		resp, err := tc.client.ListConfigurationSets(tc.ctx, &sesv2.ListConfigurationSetsInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.ConfigurationSets...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *sesv2TestContext) listAllContactLists() ([]types.ContactList, error) {
	var all []types.ContactList
	var nextToken *string
	for {
		resp, err := tc.client.ListContactLists(tc.ctx, &sesv2.ListContactListsInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.ContactLists...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *sesv2TestContext) listAllEmailIdentities() ([]types.IdentityInfo, error) {
	var all []types.IdentityInfo
	var nextToken *string
	for {
		resp, err := tc.client.ListEmailIdentities(tc.ctx, &sesv2.ListEmailIdentitiesInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.EmailIdentities...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *sesv2TestContext) listAllDedicatedIpPools() ([]string, error) {
	var all []string
	var nextToken *string
	for {
		resp, err := tc.client.ListDedicatedIpPools(tc.ctx, &sesv2.ListDedicatedIpPoolsInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.DedicatedIpPools...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *sesv2TestContext) listAllEmailTemplates() ([]types.EmailTemplateMetadata, error) {
	var all []types.EmailTemplateMetadata
	var nextToken *string
	for {
		resp, err := tc.client.ListEmailTemplates(tc.ctx, &sesv2.ListEmailTemplatesInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.TemplatesMetadata...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *sesv2TestContext) listAllSuppressedDestinations() ([]types.SuppressedDestinationSummary, error) {
	var all []types.SuppressedDestinationSummary
	var nextToken *string
	for {
		resp, err := tc.client.ListSuppressedDestinations(tc.ctx, &sesv2.ListSuppressedDestinationsInput{
			PageSize:  aws.Int32(100),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.SuppressedDestinationSummaries...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func containsIdentityName(identities []types.IdentityInfo, name string) bool {
	for _, id := range identities {
		if id.IdentityName != nil && *id.IdentityName == name {
			return true
		}
	}
	return false
}

func containsTemplateName(templates []types.EmailTemplateMetadata, name string) bool {
	for _, t := range templates {
		if t.TemplateName != nil && *t.TemplateName == name {
			return true
		}
	}
	return false
}

func containsContactListName(lists []types.ContactList, name string) bool {
	for _, l := range lists {
		if l.ContactListName != nil && *l.ContactListName == name {
			return true
		}
	}
	return false
}

func containsPoolName(pools []string, name string) bool {
	return containsString(pools, name)
}

func containsSuppressedEmail(summaries []types.SuppressedDestinationSummary, email string) bool {
	for _, s := range summaries {
		if s.EmailAddress != nil && *s.EmailAddress == email {
			return true
		}
	}
	return false
}

func filterConfigSetsByPrefix(all []string, prefix string) []string {
	var filtered []string
	for _, cs := range all {
		if strings.HasPrefix(cs, prefix) {
			filtered = append(filtered, cs)
		}
	}
	return filtered
}
