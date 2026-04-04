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

func (r *TestRunner) RunSESv2Tests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "sesv2",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := sesv2.NewFromConfig(cfg)
	ctx := context.Background()
	uid := time.Now().UnixNano()

	// --- Account Operations ---

	results = append(results, r.RunTest("sesv2", "GetAccount", func() error {
		resp, err := client.GetAccount(ctx, &sesv2.GetAccountInput{})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountSendingAttributes", func() error {
		_, err := client.PutAccountSendingAttributes(ctx, &sesv2.PutAccountSendingAttributesInput{
			SendingEnabled: true,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountSuppressionAttributes", func() error {
		_, err := client.PutAccountSuppressionAttributes(ctx, &sesv2.PutAccountSuppressionAttributesInput{
			SuppressedReasons: []types.SuppressionListReason{types.SuppressionListReasonBounce, types.SuppressionListReasonComplaint},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountDetails", func() error {
		_, err := client.PutAccountDetails(ctx, &sesv2.PutAccountDetailsInput{
			MailType:           types.MailTypeTransactional,
			UseCaseDescription: aws.String("SDK test"),
			WebsiteURL:         aws.String("https://example.com"),
			ContactLanguage:    types.ContactLanguageEn,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountVdmAttributes", func() error {
		_, err := client.PutAccountVdmAttributes(ctx, &sesv2.PutAccountVdmAttributesInput{
			VdmAttributes: &types.VdmAttributes{
				VdmEnabled: types.FeatureStatusEnabled,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutAccountDedicatedIpWarmupAttributes", func() error {
		_, err := client.PutAccountDedicatedIpWarmupAttributes(ctx, &sesv2.PutAccountDedicatedIpWarmupAttributesInput{
			AutoWarmupEnabled: false,
		})
		return err
	}))

	// --- Configuration Set Operations ---

	configSetName := fmt.Sprintf("test-cs-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSet", func() error {
		_, err := client.CreateConfigurationSet(ctx, &sesv2.CreateConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetConfigurationSet", func() error {
		resp, err := client.GetConfigurationSet(ctx, &sesv2.GetConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		if resp.ConfigurationSetName == nil || *resp.ConfigurationSetName != configSetName {
			return fmt.Errorf("expected %s, got %v", configSetName, resp.ConfigurationSetName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListConfigurationSets", func() error {
		resp, err := client.ListConfigurationSets(ctx, &sesv2.ListConfigurationSetsInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.ConfigurationSets == nil {
			return fmt.Errorf("configuration sets list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetSendingOptions", func() error {
		_, err := client.PutConfigurationSetSendingOptions(ctx, &sesv2.PutConfigurationSetSendingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			SendingEnabled:       true,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetReputationOptions", func() error {
		_, err := client.PutConfigurationSetReputationOptions(ctx, &sesv2.PutConfigurationSetReputationOptionsInput{
			ConfigurationSetName:     aws.String(configSetName),
			ReputationMetricsEnabled: true,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetDeliveryOptions", func() error {
		_, err := client.PutConfigurationSetDeliveryOptions(ctx, &sesv2.PutConfigurationSetDeliveryOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			MaxDeliverySeconds:   aws.Int64(30),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetSuppressionOptions", func() error {
		_, err := client.PutConfigurationSetSuppressionOptions(ctx, &sesv2.PutConfigurationSetSuppressionOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			SuppressedReasons: []types.SuppressionListReason{
				types.SuppressionListReasonBounce,
				types.SuppressionListReasonComplaint,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetTrackingOptions", func() error {
		_, err := client.PutConfigurationSetTrackingOptions(ctx, &sesv2.PutConfigurationSetTrackingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			HttpsPolicy:          types.HttpsPolicyRequire,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetVdmOptions", func() error {
		_, err := client.PutConfigurationSetVdmOptions(ctx, &sesv2.PutConfigurationSetVdmOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			VdmOptions: &types.VdmOptions{
				DashboardOptions: &types.DashboardOptions{
					EngagementMetrics: types.FeatureStatusEnabled,
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutConfigurationSetArchivingOptions", func() error {
		_, err := client.PutConfigurationSetArchivingOptions(ctx, &sesv2.PutConfigurationSetArchivingOptionsInput{
			ConfigurationSetName: aws.String(configSetName),
			ArchiveArn:           aws.String("arn:aws:mailmanager:us-east-1:000000000000:archive/test"),
		})
		return err
	}))

	// Event destination operations
	eventDestName := fmt.Sprintf("test-dest-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateConfigurationSetEventDestination", func() error {
		_, err := client.CreateConfigurationSetEventDestination(ctx, &sesv2.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetConfigurationSetEventDestinations", func() error {
		resp, err := client.GetConfigurationSetEventDestinations(ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		if err != nil {
			return err
		}
		if resp.EventDestinations == nil {
			return fmt.Errorf("event destinations is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateConfigurationSetEventDestination", func() error {
		_, err := client.UpdateConfigurationSetEventDestination(ctx, &sesv2.UpdateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
			EventDestination: &types.EventDestinationDefinition{
				Enabled:            true,
				MatchingEventTypes: []types.EventType{types.EventTypeSend, types.EventTypeDelivery},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSetEventDestination", func() error {
		_, err := client.DeleteConfigurationSetEventDestination(ctx, &sesv2.DeleteConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(configSetName),
			EventDestinationName: aws.String(eventDestName),
		})
		return err
	}))

	// Tag operations on configuration set
	results = append(results, r.RunTest("sesv2", "TagResource_ConfigSet", func() error {
		_, err := client.TagResource(ctx, &sesv2.TagResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:ses:us-east-1:000000000000:configuration-set/%s", configSetName)),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "ListTagsForResource_ConfigSet", func() error {
		resp, err := client.ListTagsForResource(ctx, &sesv2.ListTagsForResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:ses:us-east-1:000000000000:configuration-set/%s", configSetName)),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UntagResource_ConfigSet", func() error {
		_, err := client.UntagResource(ctx, &sesv2.UntagResourceInput{
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:ses:us-east-1:000000000000:configuration-set/%s", configSetName)),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteConfigurationSet", func() error {
		_, err := client.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String(configSetName),
		})
		return err
	}))

	// --- Email Identity Operations ---

	emailAddress := fmt.Sprintf("test-%d@example.com", uid)
	domainIdentity := fmt.Sprintf("test-%d.example.com", uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentity_Email", func() error {
		_, err := client.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentity_Domain", func() error {
		resp, err := client.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.DkimAttributes == nil {
			return fmt.Errorf("DKIM attributes nil for domain identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentity", func() error {
		resp, err := client.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListEmailIdentities", func() error {
		resp, err := client.ListEmailIdentities(ctx, &sesv2.ListEmailIdentitiesInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.EmailIdentities == nil {
			return fmt.Errorf("email identities list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityDkimAttributes", func() error {
		_, err := client.PutEmailIdentityDkimAttributes(ctx, &sesv2.PutEmailIdentityDkimAttributesInput{
			EmailIdentity:  aws.String(domainIdentity),
			SigningEnabled: true,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityDkimSigningAttributes", func() error {
		_, err := client.PutEmailIdentityDkimSigningAttributes(ctx, &sesv2.PutEmailIdentityDkimSigningAttributesInput{
			EmailIdentity:           aws.String(domainIdentity),
			SigningAttributesOrigin: types.DkimSigningAttributesOriginAwsSes,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityFeedbackAttributes", func() error {
		_, err := client.PutEmailIdentityFeedbackAttributes(ctx, &sesv2.PutEmailIdentityFeedbackAttributesInput{
			EmailIdentity:          aws.String(emailAddress),
			EmailForwardingEnabled: true,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityMailFromAttributes", func() error {
		_, err := client.PutEmailIdentityMailFromAttributes(ctx, &sesv2.PutEmailIdentityMailFromAttributesInput{
			EmailIdentity:       aws.String(domainIdentity),
			MailFromDomain:      aws.String(fmt.Sprintf("mail.%s", domainIdentity)),
			BehaviorOnMxFailure: types.BehaviorOnMxFailureUseDefaultValue,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityConfigurationSetAttributes", func() error {
		csName := fmt.Sprintf("test-cs2-%d", uid)
		_, _ = client.CreateConfigurationSet(ctx, &sesv2.CreateConfigurationSetInput{
			ConfigurationSetName: aws.String(csName),
		})
		_, err := client.PutEmailIdentityConfigurationSetAttributes(ctx, &sesv2.PutEmailIdentityConfigurationSetAttributesInput{
			EmailIdentity:        aws.String(emailAddress),
			ConfigurationSetName: aws.String(csName),
		})
		if err != nil {
			return err
		}
		client.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String(csName),
		})
		return nil
	}))

	policyName := fmt.Sprintf("test-policy-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentityPolicy", func() error {
		_, err := client.CreateEmailIdentityPolicy(ctx, &sesv2.CreateEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
			Policy:        aws.String(`{"Version":"2008-10-17","Statement":[]}`),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentityPolicies", func() error {
		resp, err := client.GetEmailIdentityPolicies(ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		if resp.Policies == nil {
			return fmt.Errorf("policies map is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateEmailIdentityPolicy", func() error {
		_, err := client.UpdateEmailIdentityPolicy(ctx, &sesv2.UpdateEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
			Policy:        aws.String(`{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"SES:SendEmail","Resource":"*"}]}`),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentityPolicy", func() error {
		_, err := client.DeleteEmailIdentityPolicy(ctx, &sesv2.DeleteEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_Domain", func() error {
		_, err := client.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		return err
	}))

	// --- Email Sending ---

	results = append(results, r.RunTest("sesv2", "SendEmail", func() error {
		resp, err := client.SendEmail(ctx, &sesv2.SendEmailInput{
			FromEmailAddress: aws.String(emailAddress),
			Destination: &types.Destination{
				ToAddresses: []string{emailAddress},
			},
			Content: &types.EmailContent{
				Simple: &types.Message{
					Subject: &types.Content{
						Data: aws.String("Test Subject"),
					},
					Body: &types.Body{
						Text: &types.Content{
							Data: aws.String("Test Body"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.MessageId == nil {
			return fmt.Errorf("message ID is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_Email", func() error {
		_, err := client.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		return err
	}))

	// --- Template Operations ---

	templateName := fmt.Sprintf("test-tmpl-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailTemplate", func() error {
		_, err := client.CreateEmailTemplate(ctx, &sesv2.CreateEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateContent: &types.EmailTemplateContent{
				Subject: aws.String("Hello {{name}}"),
				Html:    aws.String("<h1>Hello {{name}}</h1>"),
				Text:    aws.String("Hello {{name}}"),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailTemplate", func() error {
		resp, err := client.GetEmailTemplate(ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return err
		}
		if resp.TemplateContent == nil {
			return fmt.Errorf("template content is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListEmailTemplates", func() error {
		resp, err := client.ListEmailTemplates(ctx, &sesv2.ListEmailTemplatesInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.TemplatesMetadata == nil {
			return fmt.Errorf("templates metadata is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateEmailTemplate", func() error {
		_, err := client.UpdateEmailTemplate(ctx, &sesv2.UpdateEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateContent: &types.EmailTemplateContent{
				Subject: aws.String("Updated {{name}}"),
				Html:    aws.String("<h1>Updated {{name}}</h1>"),
				Text:    aws.String("Updated {{name}}"),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "TestRenderEmailTemplate", func() error {
		resp, err := client.TestRenderEmailTemplate(ctx, &sesv2.TestRenderEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateData: aws.String(`{"name":"World"}`),
		})
		if err != nil {
			return err
		}
		if resp.RenderedTemplate == nil {
			return fmt.Errorf("rendered template is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailTemplate", func() error {
		_, err := client.DeleteEmailTemplate(ctx, &sesv2.DeleteEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		return err
	}))

	// --- Dedicated IP Pool Operations ---

	poolName := fmt.Sprintf("test-pool-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateDedicatedIpPool", func() error {
		_, err := client.CreateDedicatedIpPool(ctx, &sesv2.CreateDedicatedIpPoolInput{
			PoolName:    aws.String(poolName),
			ScalingMode: types.ScalingModeStandard,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetDedicatedIpPool", func() error {
		resp, err := client.GetDedicatedIpPool(ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return err
		}
		if resp.DedicatedIpPool == nil || resp.DedicatedIpPool.PoolName == nil || *resp.DedicatedIpPool.PoolName != poolName {
			return fmt.Errorf("expected pool name %s, got %v", poolName, resp.DedicatedIpPool)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetDedicatedIpPool", func() error {
		resp, err := client.GetDedicatedIpPool(ctx, &sesv2.GetDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return err
		}
		if resp.DedicatedIpPool == nil || resp.DedicatedIpPool.PoolName == nil || *resp.DedicatedIpPool.PoolName != poolName {
			return fmt.Errorf("expected pool name %s, got %v", poolName, resp.DedicatedIpPool)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListDedicatedIpPools", func() error {
		resp, err := client.ListDedicatedIpPools(ctx, &sesv2.ListDedicatedIpPoolsInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.DedicatedIpPools == nil {
			return fmt.Errorf("dedicated IP pools list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteDedicatedIpPool", func() error {
		_, err := client.DeleteDedicatedIpPool(ctx, &sesv2.DeleteDedicatedIpPoolInput{
			PoolName: aws.String(poolName),
		})
		return err
	}))

	// --- Suppression List Operations ---

	suppressedEmail := fmt.Sprintf("suppressed-%d@example.com", uid)

	results = append(results, r.RunTest("sesv2", "PutSuppressedDestination", func() error {
		_, err := client.PutSuppressedDestination(ctx, &sesv2.PutSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
			Reason:       types.SuppressionListReasonBounce,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetSuppressedDestination", func() error {
		resp, err := client.GetSuppressedDestination(ctx, &sesv2.GetSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
		})
		if err != nil {
			return err
		}
		if resp.SuppressedDestination == nil {
			return fmt.Errorf("suppressed destination is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListSuppressedDestinations", func() error {
		resp, err := client.ListSuppressedDestinations(ctx, &sesv2.ListSuppressedDestinationsInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.SuppressedDestinationSummaries == nil {
			return fmt.Errorf("suppressed destination summaries is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteSuppressedDestination", func() error {
		_, err := client.DeleteSuppressedDestination(ctx, &sesv2.DeleteSuppressedDestinationInput{
			EmailAddress: aws.String(suppressedEmail),
		})
		return err
	}))

	// --- Contact List Operations ---

	contactListName := fmt.Sprintf("test-cl-%d", uid)

	results = append(results, r.RunTest("sesv2", "CreateContactList", func() error {
		_, err := client.CreateContactList(ctx, &sesv2.CreateContactListInput{
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
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetContactList", func() error {
		resp, err := client.GetContactList(ctx, &sesv2.GetContactListInput{
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
		resp, err := client.ListContactLists(ctx, &sesv2.ListContactListsInput{
			PageSize: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.ContactLists == nil {
			return fmt.Errorf("contact lists is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateContactList", func() error {
		_, err := client.UpdateContactList(ctx, &sesv2.UpdateContactListInput{
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
		return err
	}))

	// --- Contact Operations ---

	contactEmail := fmt.Sprintf("contact-%d@example.com", uid)

	results = append(results, r.RunTest("sesv2", "CreateContact", func() error {
		_, err := client.CreateContact(ctx, &sesv2.CreateContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
			TopicPreferences: []types.TopicPreference{
				{
					TopicName:          aws.String("Updates"),
					SubscriptionStatus: types.SubscriptionStatusOptIn,
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "GetContact", func() error {
		resp, err := client.GetContact(ctx, &sesv2.GetContactInput{
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
		resp, err := client.ListContacts(ctx, &sesv2.ListContactsInput{
			ContactListName: aws.String(contactListName),
			PageSize:        aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.Contacts == nil {
			return fmt.Errorf("contacts list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateContact", func() error {
		_, err := client.UpdateContact(ctx, &sesv2.UpdateContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
			TopicPreferences: []types.TopicPreference{
				{
					TopicName:          aws.String("Newsletter"),
					SubscriptionStatus: types.SubscriptionStatusOptIn,
				},
			},
			UnsubscribeAll: false,
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContact", func() error {
		_, err := client.DeleteContact(ctx, &sesv2.DeleteContactInput{
			ContactListName: aws.String(contactListName),
			EmailAddress:    aws.String(contactEmail),
		})
		return err
	}))

	results = append(results, r.RunTest("sesv2", "DeleteContactList", func() error {
		_, err := client.DeleteContactList(ctx, &sesv2.DeleteContactListInput{
			ContactListName: aws.String(contactListName),
		})
		return err
	}))

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("sesv2", "ListConfigurationSets_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgConfigSets []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagCS-%s-%d", pgTs, i)
			_, err := client.CreateConfigurationSet(ctx, &sesv2.CreateConfigurationSetInput{
				ConfigurationSetName: aws.String(name),
			})
			if err != nil {
				for _, cn := range pgConfigSets {
					client.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{ConfigurationSetName: aws.String(cn)})
				}
				return fmt.Errorf("create configuration set %s: %v", name, err)
			}
			pgConfigSets = append(pgConfigSets, name)
		}

		var allConfigSets []string
		var nextToken *string
		for {
			resp, err := client.ListConfigurationSets(ctx, &sesv2.ListConfigurationSetsInput{
				PageSize:  aws.Int32(2),
				NextToken: nextToken,
			})
			if err != nil {
				for _, cn := range pgConfigSets {
					client.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{ConfigurationSetName: aws.String(cn)})
				}
				return fmt.Errorf("list configuration sets page: %v", err)
			}
			for _, cs := range resp.ConfigurationSets {
				if strings.Contains(cs, "PagCS-"+pgTs) {
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
			client.DeleteConfigurationSet(ctx, &sesv2.DeleteConfigurationSetInput{ConfigurationSetName: aws.String(cn)})
		}
		if len(allConfigSets) != 5 {
			return fmt.Errorf("expected 5 paginated configuration sets, got %d", len(allConfigSets))
		}
		return nil
	}))

	return results
}
