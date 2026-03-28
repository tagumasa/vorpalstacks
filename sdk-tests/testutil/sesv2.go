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

	// Test 1: Get Account
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

	// Test 2: Create Email Identity
	emailAddress := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	results = append(results, r.RunTest("sesv2", "CreateEmailIdentity", func() error {
		resp, err := client.CreateEmailIdentity(ctx, &sesv2.CreateEmailIdentityInput{
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

	// Test 3: Get Email Identity
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

	// Test 4: List Email Identities
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

	// Test 5: Create Email Identity Policy
	policyName := fmt.Sprintf("test-policy-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("sesv2", "CreateEmailIdentityPolicy", func() error {
		resp, err := client.CreateEmailIdentityPolicy(ctx, &sesv2.CreateEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
			Policy:        aws.String(`{"Version":"2008-10-17","Statement":[]}`),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 6: Get Email Identity Policies
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

	// Test 7: Delete Email Identity Policy (cleanup the policy created in Test 5)
	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentityPolicy", func() error {
		resp, err := client.DeleteEmailIdentityPolicy(ctx, &sesv2.DeleteEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 8: Put Email Identity Feedback Attributes
	results = append(results, r.RunTest("sesv2", "PutEmailIdentityFeedbackAttributes", func() error {
		resp, err := client.PutEmailIdentityFeedbackAttributes(ctx, &sesv2.PutEmailIdentityFeedbackAttributesInput{
			EmailIdentity:          aws.String(emailAddress),
			EmailForwardingEnabled: true,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 9: Send Email
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

	// Test 10: Create Contact List
	contactListName := fmt.Sprintf("test-contactlist-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("sesv2", "CreateContactList", func() error {
		resp, err := client.CreateContactList(ctx, &sesv2.CreateContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 11: List Contact Lists
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

	// Test 12: Get Contact List
	results = append(results, r.RunTest("sesv2", "GetContactList", func() error {
		resp, err := client.GetContactList(ctx, &sesv2.GetContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return err
		}
		if resp.ContactListName == nil {
			return fmt.Errorf("contact list name is nil")
		}
		return nil
	}))

	// Test 13: Delete Contact List
	results = append(results, r.RunTest("sesv2", "DeleteContactList", func() error {
		resp, err := client.DeleteContactList(ctx, &sesv2.DeleteContactListInput{
			ContactListName: aws.String(contactListName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 14: Delete Email Identity
	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity", func() error {
		resp, err := client.DeleteEmailIdentity(ctx, &sesv2.DeleteEmailIdentityInput{
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

	return results
}
