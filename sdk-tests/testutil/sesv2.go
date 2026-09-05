package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"vorpalstacks-sdk-tests/config"
)

type sesv2TestContext struct {
	client    *sesv2.Client
	ctx       context.Context
	uid       int64
	region    string
	accountID string
}

func newSESTestContext(endpoint, region, accountID string) (*sesv2TestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: endpoint,
		Region:   region,
	})
	if err != nil {
		return nil, fmt.Errorf("load config: %v", err)
	}
	return &sesv2TestContext{
		client:    sesv2.NewFromConfig(cfg),
		ctx:       context.Background(),
		uid:       time.Now().UnixNano(),
		region:    region,
		accountID: accountID,
	}, nil
}

func (tc *sesv2TestContext) configSetARN(name string) string {
	return fmt.Sprintf("arn:aws:ses:%s:%s:configuration-set/%s", tc.region, tc.accountID, name)
}

func (tc *sesv2TestContext) contactListARN(name string) string {
	return fmt.Sprintf("arn:aws:ses:%s:%s:contact-list/%s", tc.region, tc.accountID, name)
}

func (tc *sesv2TestContext) identityARN(name string) string {
	return fmt.Sprintf("arn:aws:ses:%s:%s:identity/%s", tc.region, tc.accountID, name)
}

// createConfigSet is the plain (ConfigurationSetName only) form of
// CreateConfigurationSet used for setup; a creation carrying delivery,
// tracking or other options is the operation under test and keeps its
// literal input.
func (tc *sesv2TestContext) createConfigSet(name string) error {
	_, err := tc.client.CreateConfigurationSet(tc.ctx, &sesv2.CreateConfigurationSetInput{
		ConfigurationSetName: aws.String(name),
	})
	return err
}

// getConfigSet is the plain (ConfigurationSetName only) form of
// GetConfigurationSet.
func (tc *sesv2TestContext) getConfigSet(name string) (*sesv2.GetConfigurationSetOutput, error) {
	return tc.client.GetConfigurationSet(tc.ctx, &sesv2.GetConfigurationSetInput{
		ConfigurationSetName: aws.String(name),
	})
}

// getEventDestinations is the plain (ConfigurationSetName only) form of
// GetConfigurationSetEventDestinations.
func (tc *sesv2TestContext) getEventDestinations(csName string) (*sesv2.GetConfigurationSetEventDestinationsOutput, error) {
	return tc.client.GetConfigurationSetEventDestinations(tc.ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
		ConfigurationSetName: aws.String(csName),
	})
}

// getEmailIdentity is the plain (EmailIdentity only) form of
// GetEmailIdentity.
func (tc *sesv2TestContext) getEmailIdentity(identity string) (*sesv2.GetEmailIdentityOutput, error) {
	return tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: aws.String(identity),
	})
}

// getIdentityPolicies is the plain (EmailIdentity only) form of
// GetEmailIdentityPolicies.
func (tc *sesv2TestContext) getIdentityPolicies(identity string) (*sesv2.GetEmailIdentityPoliciesOutput, error) {
	return tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
		EmailIdentity: aws.String(identity),
	})
}

// createContactList is the plain (ContactListName only, no topics) form of
// CreateContactList used for setup; a creation carrying topics is the
// operation under test and keeps its literal input.
func (tc *sesv2TestContext) createContactList(name string) error {
	_, err := tc.client.CreateContactList(tc.ctx, &sesv2.CreateContactListInput{
		ContactListName: aws.String(name),
	})
	return err
}

// getContactList is the plain (ContactListName only) form of GetContactList.
func (tc *sesv2TestContext) getContactList(name string) (*sesv2.GetContactListOutput, error) {
	return tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
		ContactListName: aws.String(name),
	})
}

// getContact is the plain (ContactListName + EmailAddress only) form of
// GetContact.
func (tc *sesv2TestContext) getContact(listName, email string) (*sesv2.GetContactOutput, error) {
	return tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
		ContactListName: aws.String(listName),
		EmailAddress:    aws.String(email),
	})
}

// getEmailTemplate is the plain (TemplateName only) form of
// GetEmailTemplate.
func (tc *sesv2TestContext) getEmailTemplate(name string) (*sesv2.GetEmailTemplateOutput, error) {
	return tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
		TemplateName: aws.String(name),
	})
}

// getAccount is the empty-input form of GetAccount.
func (tc *sesv2TestContext) getAccount() (*sesv2.GetAccountOutput, error) {
	return tc.client.GetAccount(tc.ctx, &sesv2.GetAccountInput{})
}

// listTags returns the tags carried by a resource ARN.
func (tc *sesv2TestContext) listTags(arn string) ([]types.Tag, error) {
	resp, err := tc.client.ListTagsForResource(tc.ctx, &sesv2.ListTagsForResourceInput{
		ResourceArn: aws.String(arn),
	})
	if err != nil {
		return nil, err
	}
	return resp.Tags, nil
}

// expectNotFound asserts the typed NotFoundException contract of a failed
// call, tagging the failure with the operation name.
func (tc *sesv2TestContext) expectNotFound(op string, err error) error {
	if aerr := expectAWSErrorCode(err, "NotFoundException"); aerr != nil {
		return fmt.Errorf("%s: %w", op, aerr)
	}
	return nil
}

// containsTag reports whether the tag list carries key with the given value.
func containsTag(tags []types.Tag, key, value string) bool {
	for _, t := range tags {
		if t.Key != nil && *t.Key == key && t.Value != nil && *t.Value == value {
			return true
		}
	}
	return false
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
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListConfigurationSets(tc.ctx, &sesv2.ListConfigurationSetsInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.ConfigurationSets, resp.NextToken, nil
	})
}

func (tc *sesv2TestContext) listAllContactLists() ([]types.ContactList, error) {
	return paginate(func(next *string) ([]types.ContactList, *string, error) {
		resp, err := tc.client.ListContactLists(tc.ctx, &sesv2.ListContactListsInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.ContactLists, resp.NextToken, nil
	})
}

func (tc *sesv2TestContext) listAllEmailIdentities() ([]types.IdentityInfo, error) {
	return paginate(func(next *string) ([]types.IdentityInfo, *string, error) {
		resp, err := tc.client.ListEmailIdentities(tc.ctx, &sesv2.ListEmailIdentitiesInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.EmailIdentities, resp.NextToken, nil
	})
}

func (tc *sesv2TestContext) listAllDedicatedIpPools() ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListDedicatedIpPools(tc.ctx, &sesv2.ListDedicatedIpPoolsInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DedicatedIpPools, resp.NextToken, nil
	})
}

func (tc *sesv2TestContext) listAllEmailTemplates() ([]types.EmailTemplateMetadata, error) {
	return paginate(func(next *string) ([]types.EmailTemplateMetadata, *string, error) {
		resp, err := tc.client.ListEmailTemplates(tc.ctx, &sesv2.ListEmailTemplatesInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.TemplatesMetadata, resp.NextToken, nil
	})
}

func (tc *sesv2TestContext) listAllSuppressedDestinations() ([]types.SuppressedDestinationSummary, error) {
	return paginate(func(next *string) ([]types.SuppressedDestinationSummary, *string, error) {
		resp, err := tc.client.ListSuppressedDestinations(tc.ctx, &sesv2.ListSuppressedDestinationsInput{
			PageSize:  aws.Int32(100),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.SuppressedDestinationSummaries, resp.NextToken, nil
	})
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

func containsSuppressedEmail(summaries []types.SuppressedDestinationSummary, email string) bool {
	for _, s := range summaries {
		if s.EmailAddress != nil && *s.EmailAddress == email {
			return true
		}
	}
	return false
}
