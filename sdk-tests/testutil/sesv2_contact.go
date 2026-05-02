package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2ContactTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	contactListName := fmt.Sprintf("test-cl-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateContactList", func() error {
		_, err := tc.client.CreateContactList(tc.ctx, &sesv2.CreateContactListInput{
			ContactListName: aws.String(contactListName),
			Topics: []types.Topic{
				{
					TopicName:                 aws.String("Updates"),
					DefaultSubscriptionStatus: types.SubscriptionStatusOptIn,
					Description:               aws.String("Product updates"),
					DisplayName:               aws.String("Updates"),
				},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return fmt.Errorf("get after create: %v", err)
		}
		if resp.ContactListName == nil || *resp.ContactListName != contactListName {
			return fmt.Errorf("expected name %s, got %v", contactListName, resp.ContactListName)
		}
		if len(resp.Topics) != 1 {
			return fmt.Errorf("expected 1 topic, got %d", len(resp.Topics))
		}
		if resp.Topics[0].TopicName == nil || *resp.Topics[0].TopicName != "Updates" {
			return fmt.Errorf("expected topic name Updates, got %v", resp.Topics[0].TopicName)
		}
		if resp.Topics[0].DefaultSubscriptionStatus != types.SubscriptionStatusOptIn {
			return fmt.Errorf("expected OptIn, got %s", resp.Topics[0].DefaultSubscriptionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetContactList", func() error {
		resp, err := tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return err
		}
		if resp.ContactListName == nil || *resp.ContactListName != contactListName {
			return fmt.Errorf("expected %s, got %v", contactListName, resp.ContactListName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListContactLists", func() error {
		all, err := tc.listAllContactLists()
		if err != nil {
			return err
		}
		if !containsContactListName(all, contactListName) {
			return fmt.Errorf("contact list %s not found", contactListName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateContactList", func() error {
		_, err := tc.client.UpdateContactList(tc.ctx, &sesv2.UpdateContactListInput{
			ContactListName: aws.String(contactListName),
			Topics: []types.Topic{
				{
					TopicName:                 aws.String("Newsletter"),
					DefaultSubscriptionStatus: types.SubscriptionStatusOptIn,
					Description:               aws.String("Weekly newsletter"),
					DisplayName:               aws.String("Newsletter"),
				},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return fmt.Errorf("get after update: %v", err)
		}
		if len(resp.Topics) != 1 || resp.Topics[0].TopicName == nil || *resp.Topics[0].TopicName != "Newsletter" {
			return fmt.Errorf("topic not updated to Newsletter")
		}
		return nil
	}))

	contactEmail := fmt.Sprintf("contact-%d@example.com", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateContact", func() error {
		_, err := tc.client.CreateContact(tc.ctx, &sesv2.CreateContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
			TopicPreferences: []types.TopicPreference{
				{
					TopicName:          aws.String("Newsletter"),
					SubscriptionStatus: types.SubscriptionStatusOptIn,
				},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		if err != nil {
			return fmt.Errorf("get after create contact: %v", err)
		}
		if resp.EmailAddress == nil || *resp.EmailAddress != contactEmail {
			return fmt.Errorf("expected email %s, got %v", contactEmail, resp.EmailAddress)
		}
		if len(resp.TopicPreferences) != 1 {
			return fmt.Errorf("expected 1 topic preference, got %d", len(resp.TopicPreferences))
		}
		if resp.TopicPreferences[0].TopicName == nil || *resp.TopicPreferences[0].TopicName != "Newsletter" {
			return fmt.Errorf("expected topic preference Newsletter, got %v", resp.TopicPreferences[0].TopicName)
		}
		if resp.TopicPreferences[0].SubscriptionStatus != types.SubscriptionStatusOptIn {
			return fmt.Errorf("expected OptIn, got %s", resp.TopicPreferences[0].SubscriptionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetContact", func() error {
		resp, err := tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		if err != nil {
			return err
		}
		if resp.EmailAddress == nil || *resp.EmailAddress != contactEmail {
			return fmt.Errorf("expected %s, got %v", contactEmail, resp.EmailAddress)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListContacts", func() error {
		resp, err := tc.client.ListContacts(tc.ctx, &sesv2.ListContactsInput{
			ContactListName: aws.String(contactListName),
			PageSize:        aws.Int32(100),
		})
		if err != nil {
			return err
		}
		found := false
		for _, c := range resp.Contacts {
			if c.EmailAddress != nil && *c.EmailAddress == contactEmail {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("contact %s not found in list", contactEmail)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateContact", func() error {
		_, err := tc.client.UpdateContact(tc.ctx, &sesv2.UpdateContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
			TopicPreferences: []types.TopicPreference{
				{
					TopicName:          aws.String("Newsletter"),
					SubscriptionStatus: types.SubscriptionStatusOptOut,
				},
			},
			UnsubscribeAll: false,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		if err != nil {
			return fmt.Errorf("get after update contact: %v", err)
		}
		if len(resp.TopicPreferences) != 1 {
			return fmt.Errorf("expected 1 topic preference, got %d", len(resp.TopicPreferences))
		}
		if resp.TopicPreferences[0].SubscriptionStatus != types.SubscriptionStatusOptOut {
			return fmt.Errorf("expected OptOut after update, got %s", resp.TopicPreferences[0].SubscriptionStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContact", func() error {
		_, err := tc.client.DeleteContact(tc.ctx, &sesv2.DeleteContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetContact(tc.ctx, &sesv2.GetContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted contact")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContactList", func() error {
		_, err := tc.client.DeleteContactList(tc.ctx, &sesv2.DeleteContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetContactList(tc.ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted contact list")
		}
		return nil
	}))

	return results
}
