package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
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
		_, err := tc.getConfigSet("nonexistent-cs-xyz")
		return tc.expectNotFound("GetConfigurationSet", err)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSet_NonExistent", func() error {
		_, err := tc.client.DeleteConfigurationSet(tc.ctx, &sesv2.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String("nonexistent-cs-xyz"),
		})
		return tc.expectNotFound("DeleteConfigurationSet", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentity_NonExistent", func() error {
		_, err := tc.getEmailIdentity("nonexistent@example.com")
		return tc.expectNotFound("GetEmailIdentity", err)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_NonExistent", func() error {
		_, err := tc.client.DeleteEmailIdentity(tc.ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String("nonexistent@example.com"),
		})
		return tc.expectNotFound("DeleteEmailIdentity", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailTemplate_NonExistent", func() error {
		_, err := tc.getEmailTemplate("nonexistent-template-xyz")
		return tc.expectNotFound("GetEmailTemplate", err)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailTemplate_NonExistent", func() error {
		_, err := tc.client.DeleteEmailTemplate(tc.ctx, &sesv2.DeleteEmailTemplateInput{
			TemplateName: aws.String("nonexistent-template-xyz"),
		})
		return tc.expectNotFound("DeleteEmailTemplate", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetDedicatedIpPool_NonExistent", func() error {
		_, err := tc.client.GetDedicatedIpPool(tc.ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String("nonexistent-pool-xyz"),
		})
		return tc.expectNotFound("GetDedicatedIpPool", err)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteDedicatedIpPool_NonExistent", func() error {
		_, err := tc.client.DeleteDedicatedIpPool(tc.ctx, &sesv2.DeleteDedicatedIpPoolInput{
			PoolName: aws.String("nonexistent-pool-xyz"),
		})
		return tc.expectNotFound("DeleteDedicatedIpPool", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetContactList_NonExistent", func() error {
		_, err := tc.getContactList("nonexistent-cl-xyz")
		return tc.expectNotFound("GetContactList", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetContact_NonExistent", func() error {
		_, err := tc.getContact("nonexistent-cl-xyz", "nonexistent@example.com")
		return tc.expectNotFound("GetContact", err)
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContactList_NonExistent", func() error {
		_, err := tc.client.DeleteContactList(tc.ctx, &sesv2.DeleteContactListInput{
			ContactListName: aws.String("nonexistent-cl-xyz"),
		})
		return tc.expectNotFound("DeleteContactList", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetSuppressedDestination_NonExistent", func() error {
		_, err := tc.client.GetSuppressedDestination(tc.ctx, &sesv2.GetSuppressedDestinationInput{
			EmailAddress: aws.String("nonexistent@example.com"),
		})
		return tc.expectNotFound("GetSuppressedDestination", err)
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentityPolicies_NonExistent", func() error {
		_, err := tc.getIdentityPolicies("nonexistent@example.com")
		return tc.expectNotFound("GetEmailIdentityPolicies", err)
	}))

	// --- Validation error-case tests ---

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetDeliveryOptions_InvalidMaxDeliverySeconds", func() error {
		csName := fmt.Sprintf("edge-cs-%d", tc.uid)
		if err := tc.createConfigSet(csName); err != nil {
			return fmt.Errorf("setup create config set: %v", err)
		}
		defer tc.deleteConfigSet(csName)

		_, err := tc.client.PutConfigurationSetDeliveryOptions(tc.ctx, &sesv2.PutConfigurationSetDeliveryOptionsInput{
			ConfigurationSetName: aws.String(csName),
			MaxDeliverySeconds:   aws.Int64(100),
		})
		if err == nil {
			return fmt.Errorf("expected error for MaxDeliverySeconds=100 (min=300)")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSet_InvalidTlsPolicy", func() error {
		csName := fmt.Sprintf("edge-tls-%d", tc.uid)
		_, err := tc.client.CreateConfigurationSet(tc.ctx, &sesv2.CreateConfigurationSetInput{
			ConfigurationSetName: aws.String(csName),
			DeliveryOptions: &types.DeliveryOptions{
				TlsPolicy: types.TlsPolicy("INVALID"),
			},
		})
		if err == nil {
			defer tc.deleteConfigSet(csName)
			return fmt.Errorf("expected error for invalid TlsPolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountDetails_MissingMailType", func() error {
		_, err := tc.client.PutAccountDetails(tc.ctx, &sesv2.PutAccountDetailsInput{
			WebsiteURL: aws.String("https://example.com"),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing MailType")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSetEventDestination_MultipleDestinations", func() error {
		csName := fmt.Sprintf("edge-ed-%d", tc.uid)
		if err := tc.createConfigSet(csName); err != nil {
			return fmt.Errorf("setup create config set: %v", err)
		}
		defer tc.deleteConfigSet(csName)

		_, err := tc.client.CreateConfigurationSetEventDestination(tc.ctx, &sesv2.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(csName),
			EventDestinationName: aws.String("multi-dest"),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend},
				SnsDestination: &types.SnsDestination{
					TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test"),
				},
				EventBridgeDestination: &types.EventBridgeDestination{
					EventBusArn: aws.String("arn:aws:events:us-east-1:123456789012:test"),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for multiple destinations (union violation)")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "CreateEmailTemplate_InvalidName", func() error {
		_, err := tc.client.CreateEmailTemplate(tc.ctx, &sesv2.CreateEmailTemplateInput{
			TemplateName: aws.String("invalid name with spaces"),
			TemplateContent: &types.EmailTemplateContent{
				Subject: aws.String("Test"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid template name")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "CreateDedicatedIpPool_InvalidScalingMode", func() error {
		_, err := tc.client.CreateDedicatedIpPool(tc.ctx, &sesv2.CreateDedicatedIpPoolInput{
			PoolName:    aws.String(fmt.Sprintf("edge-pool-%d", tc.uid)),
			ScalingMode: types.ScalingMode("INVALID"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid ScalingMode")
		}
		return nil
	}))

	// --- UpdateEventDestination preserves destination type ---

	results = append(results, r.RunTest("sesv2", "UpdateEventDestination_PreserveDestinationType", func() error {
		csName := fmt.Sprintf("edge-update-preserve-%d", tc.uid)
		if err := tc.createConfigSet(csName); err != nil {
			return fmt.Errorf("setup create config set: %v", err)
		}
		defer tc.deleteConfigSet(csName)

		_, err := tc.client.CreateConfigurationSetEventDestination(tc.ctx, &sesv2.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(csName),
			EventDestinationName: aws.String("update-preserve-dest"),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend},
				SnsDestination: &types.SnsDestination{
					TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:update-preserve-sns"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create event destination: %v", err)
		}

		// Update only Enabled — destination type must be preserved.
		_, err = tc.client.UpdateConfigurationSetEventDestination(tc.ctx, &sesv2.UpdateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(csName),
			EventDestinationName: aws.String("update-preserve-dest"),
			EventDestination: &types.EventDestinationDefinition{
				Enabled: false,
			},
		})
		if err != nil {
			return fmt.Errorf("update event destination: %v", err)
		}

		// Verify SNS destination is still present.
		resp, err := tc.getEventDestinations(csName)
		if err != nil {
			return fmt.Errorf("get event destinations: %v", err)
		}
		for _, ed := range resp.EventDestinations {
			if *ed.Name == "update-preserve-dest" {
				if ed.SnsDestination == nil || ed.SnsDestination.TopicArn == nil {
					return fmt.Errorf("SNS destination was wiped by partial update")
				}
				if *ed.SnsDestination.TopicArn != "arn:aws:sns:us-east-1:123456789012:update-preserve-sns" {
					return fmt.Errorf("SNS TopicArn changed: got %s", *ed.SnsDestination.TopicArn)
				}
			}
		}
		return nil
	}))

	// --- Enabled defaults to true when omitted ---

	results = append(results, r.RunTest("sesv2", "CreateEventDestination_DefaultEnabled", func() error {
		csName := fmt.Sprintf("edge-default-enabled-%d", tc.uid)
		if err := tc.createConfigSet(csName); err != nil {
			return fmt.Errorf("setup create config set: %v", err)
		}
		defer tc.deleteConfigSet(csName)

		_, err := tc.client.CreateConfigurationSetEventDestination(tc.ctx, &sesv2.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(csName),
			EventDestinationName: aws.String("default-enabled-dest"),
			EventDestination: &types.EventDestinationDefinition{
				MatchingEventTypes: []types.EventType{types.EventTypeSend},
				SnsDestination: &types.SnsDestination{
					TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:default-enabled-sns"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create event destination: %v", err)
		}

		resp, err := tc.getEventDestinations(csName)
		if err != nil {
			return fmt.Errorf("get event destinations: %v", err)
		}
		for _, ed := range resp.EventDestinations {
			if *ed.Name == "default-enabled-dest" {
				if !ed.Enabled {
					return fmt.Errorf("Enabled should default to true when omitted")
				}
			}
		}
		return nil
	}))

	// --- Empty content rejected ---

	results = append(results, r.RunTest("sesv2", "SendEmail_EmptyContent", func() error {
		_, err := tc.client.SendEmail(tc.ctx, &sesv2.SendEmailInput{
			FromEmailAddress: aws.String("sender@example.com"),
			Destination: &types.Destination{
				ToAddresses: []string{"recipient@example.com"},
			},
			Content: &types.EmailContent{
				Simple: &types.Message{
					Body: &types.Body{},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for empty Simple body")
		}
		return nil
	}))

	// --- Missing Destination rejected ---

	results = append(results, r.RunTest("sesv2", "SendEmail_MissingDestination", func() error {
		_, err := tc.client.SendEmail(tc.ctx, &sesv2.SendEmailInput{
			FromEmailAddress: aws.String("sender@example.com"),
			Content: &types.EmailContent{
				Simple: &types.Message{
					Body: &types.Body{
						Text: &types.Content{Data: aws.String("hello")},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for missing Destination")
		}
		return nil
	}))

	// --- CreateContact with non-existent list ---

	results = append(results, r.RunTest("sesv2", "CreateContact_NonExistentList", func() error {
		_, err := tc.client.CreateContact(tc.ctx, &sesv2.CreateContactInput{
			ContactListName: aws.String(fmt.Sprintf("nonexistent-list-%d", tc.uid)),
			EmailAddress:    aws.String("contact@example.com"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent contact list")
		}
		return nil
	}))

	// --- Invalid FilteredStatus ---

	results = append(results, r.RunTest("sesv2", "ListContacts_InvalidFilteredStatus", func() error {
		clName := fmt.Sprintf("filter-status-%d", tc.uid)
		if err := tc.createContactList(clName); err != nil {
			return fmt.Errorf("setup: %v", err)
		}
		defer tc.deleteContactList(clName)

		_, err := tc.client.ListContacts(tc.ctx, &sesv2.ListContactsInput{
			ContactListName: aws.String(clName),
			Filter: &types.ListContactsFilter{
				FilteredStatus: types.SubscriptionStatus("INVALID"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid FilteredStatus")
		}
		return nil
	}))

	// --- Non-existent ConfigurationSet on SendEmail ---

	results = append(results, r.RunTest("sesv2", "SendEmail_NonExistentConfigSet", func() error {
		_, err := tc.client.SendEmail(tc.ctx, &sesv2.SendEmailInput{
			FromEmailAddress:     aws.String("sender@example.com"),
			ConfigurationSetName: aws.String("nonexistent-config-set-" + fmt.Sprintf("%d", tc.uid)),
			Destination: &types.Destination{
				ToAddresses: []string{"recipient@example.com"},
			},
			Content: &types.EmailContent{
				Simple: &types.Message{
					Body: &types.Body{
						Text: &types.Content{Data: aws.String("hello")},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent ConfigurationSetName")
		}
		return nil
	}))

	return results
}
