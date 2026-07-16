package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2EmailTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	emailAddress := fmt.Sprintf("test-%d@example.com", tc.uid)

	results = append(results, r.RunTest("sesv2", "SendEmail", func() error {
		resp, err := tc.client.SendEmail(tc.ctx, &sesv2.SendEmailInput{
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
		if resp.MessageId == nil || *resp.MessageId == "" {
			return fmt.Errorf("MessageId is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "SendBulkEmail", func() error {
		resp, err := tc.client.SendBulkEmail(tc.ctx, &sesv2.SendBulkEmailInput{
			FromEmailAddress: aws.String(emailAddress),
			DefaultContent: &types.BulkEmailContent{
				Template: &types.Template{
					TemplateName: aws.String("test-bulk-template"),
					TemplateData: aws.String(`{"name":"World"}`),
				},
			},
			BulkEmailEntries: []types.BulkEmailEntry{
				{
					Destination: &types.Destination{
						ToAddresses: []string{fmt.Sprintf("bulk1-%d@example.com", tc.uid)},
					},
				},
				{
					Destination: &types.Destination{
						ToAddresses: []string{fmt.Sprintf("bulk2-%d@example.com", tc.uid)},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if len(resp.BulkEmailEntryResults) != 2 {
			return fmt.Errorf("expected 2 bulk entry results, got %d", len(resp.BulkEmailEntryResults))
		}
		for i, entry := range resp.BulkEmailEntryResults {
			if entry.MessageId == nil || *entry.MessageId == "" {
				return fmt.Errorf("entry %d: MessageId is nil or empty", i)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_Email", func() error {
		_, err := tc.client.DeleteEmailIdentity(tc.ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted identity")
		}
		return nil
	}))

	return results
}
