package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2TemplateTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	templateName := fmt.Sprintf("test-tmpl-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailTemplate", func() error {
		_, err := tc.client.CreateEmailTemplate(tc.ctx, &sesv2.CreateEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateContent: &types.EmailTemplateContent{
				Subject: aws.String("Hello {{name}}"),
				Html:    aws.String("<h1>Hello {{name}}</h1>"),
				Text:    aws.String("Hello {{name}}"),
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return fmt.Errorf("get after create: %v", err)
		}
		if resp.TemplateContent == nil {
			return fmt.Errorf("template content is nil")
		}
		if resp.TemplateContent.Subject == nil || *resp.TemplateContent.Subject != "Hello {{name}}" {
			return fmt.Errorf("expected subject 'Hello {{name}}', got %v", resp.TemplateContent.Subject)
		}
		if resp.TemplateContent.Html == nil || *resp.TemplateContent.Html != "<h1>Hello {{name}}</h1>" {
			return fmt.Errorf("html content mismatch")
		}
		if resp.TemplateContent.Text == nil || *resp.TemplateContent.Text != "Hello {{name}}" {
			return fmt.Errorf("text content mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailTemplate", func() error {
		resp, err := tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return err
		}
		if resp.TemplateContent == nil {
			return fmt.Errorf("template content is nil")
		}
		if resp.TemplateContent.Subject == nil || *resp.TemplateContent.Subject != "Hello {{name}}" {
			return fmt.Errorf("expected subject 'Hello {{name}}', got %v", resp.TemplateContent.Subject)
		}
		if resp.TemplateContent.Html == nil || *resp.TemplateContent.Html != "<h1>Hello {{name}}</h1>" {
			return fmt.Errorf("html content mismatch")
		}
		if resp.TemplateContent.Text == nil || *resp.TemplateContent.Text != "Hello {{name}}" {
			return fmt.Errorf("text content mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListEmailTemplates", func() error {
		all, err := tc.listAllEmailTemplates()
		if err != nil {
			return err
		}
		if !containsTemplateName(all, templateName) {
			return fmt.Errorf("template %s not found in list", templateName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateEmailTemplate", func() error {
		_, err := tc.client.UpdateEmailTemplate(tc.ctx, &sesv2.UpdateEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateContent: &types.EmailTemplateContent{
				Subject: aws.String("Updated {{name}}"),
				Html:    aws.String("<h1>Updated {{name}}</h1>"),
				Text:    aws.String("Updated {{name}}"),
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return fmt.Errorf("get after update: %v", err)
		}
		if resp.TemplateContent.Subject == nil || *resp.TemplateContent.Subject != "Updated {{name}}" {
			return fmt.Errorf("subject not updated")
		}
		if resp.TemplateContent.Html == nil || *resp.TemplateContent.Html != "<h1>Updated {{name}}</h1>" {
			return fmt.Errorf("html not updated")
		}
		if resp.TemplateContent.Text == nil || *resp.TemplateContent.Text != "Updated {{name}}" {
			return fmt.Errorf("text not updated")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "TestRenderEmailTemplate", func() error {
		resp, err := tc.client.TestRenderEmailTemplate(tc.ctx, &sesv2.TestRenderEmailTemplateInput{
			TemplateName: aws.String(templateName),
			TemplateData: aws.String(`{"name":"World"}`),
		})
		if err != nil {
			return err
		}
		if resp.RenderedTemplate == nil || *resp.RenderedTemplate == "" {
			return fmt.Errorf("rendered template is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailTemplate", func() error {
		_, err := tc.client.DeleteEmailTemplate(tc.ctx, &sesv2.DeleteEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetEmailTemplate(tc.ctx, &sesv2.GetEmailTemplateInput{
			TemplateName: aws.String(templateName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting template")
		}
		return nil
	}))

	return results
}
