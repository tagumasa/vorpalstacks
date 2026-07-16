package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func (r *TestRunner) runSESv2IdentityTests(tc *sesv2TestContext) []TestResult {
	var results []TestResult

	emailAddress := fmt.Sprintf("test-%d@example.com", tc.uid)
	domainIdentity := fmt.Sprintf("test-%d.example.com", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentity_Email", func() error {
		resp, err := tc.client.CreateEmailIdentity(tc.ctx, &sesv2.CreateEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.DkimAttributes == nil {
			return fmt.Errorf("DkimAttributes is nil for email identity")
		}
		if !resp.DkimAttributes.SigningEnabled {
			return fmt.Errorf("expected SigningEnabled=true")
		}
		if len(resp.DkimAttributes.Tokens) == 0 {
			return fmt.Errorf("expected DKIM tokens")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentity_Domain", func() error {
		resp, err := tc.client.CreateEmailIdentity(tc.ctx, &sesv2.CreateEmailIdentityInput{
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
		if !resp.DkimAttributes.SigningEnabled {
			return fmt.Errorf("expected SigningEnabled=true for domain identity")
		}
		if len(resp.DkimAttributes.Tokens) == 0 {
			return fmt.Errorf("expected DKIM tokens for domain identity")
		}
		if resp.DkimAttributes.CurrentSigningKeyLength == "" {
			return fmt.Errorf("expected CurrentSigningKeyLength")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentity", func() error {
		resp, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		if resp.DkimAttributes == nil {
			return fmt.Errorf("DkimAttributes is nil for email identity")
		}
		if !resp.DkimAttributes.SigningEnabled {
			return fmt.Errorf("expected SigningEnabled=true")
		}
		if resp.IdentityType != types.IdentityTypeEmailAddress {
			return fmt.Errorf("expected IdentityType=EMAIL_ADDRESS, got %s", resp.IdentityType)
		}
		if !resp.VerifiedForSendingStatus {
			return fmt.Errorf("expected VerifiedForSendingStatus=true")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "ListEmailIdentities", func() error {
		all, err := tc.listAllEmailIdentities()
		if err != nil {
			return err
		}
		if !containsIdentityName(all, emailAddress) {
			return fmt.Errorf("email identity %s not found in list", emailAddress)
		}
		if !containsIdentityName(all, domainIdentity) {
			return fmt.Errorf("domain identity %s not found in list", domainIdentity)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityDkimAttributes", func() error {
		_, err := tc.client.PutEmailIdentityDkimAttributes(tc.ctx, &sesv2.PutEmailIdentityDkimAttributesInput{
			EmailIdentity:  aws.String(domainIdentity),
			SigningEnabled: true,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		if err != nil {
			return fmt.Errorf("get identity after DKIM put: %v", err)
		}
		if resp.DkimAttributes == nil {
			return fmt.Errorf("DKIM attributes nil after put")
		}
		if !resp.DkimAttributes.SigningEnabled {
			return fmt.Errorf("expected SigningEnabled=true after put")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityDkimSigningAttributes", func() error {
		_, err := tc.client.PutEmailIdentityDkimSigningAttributes(tc.ctx, &sesv2.PutEmailIdentityDkimSigningAttributesInput{
			EmailIdentity:           aws.String(domainIdentity),
			SigningAttributesOrigin: types.DkimSigningAttributesOriginAwsSes,
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityFeedbackAttributes", func() error {
		_, err := tc.client.PutEmailIdentityFeedbackAttributes(tc.ctx, &sesv2.PutEmailIdentityFeedbackAttributesInput{
			EmailIdentity:          aws.String(emailAddress),
			EmailForwardingEnabled: true,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return fmt.Errorf("get after feedback put: %v", err)
		}
		if !resp.FeedbackForwardingStatus {
			return fmt.Errorf("expected FeedbackForwardingStatus=true after put")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityMailFromAttributes", func() error {
		mailFrom := fmt.Sprintf("mail.%s", domainIdentity)
		_, err := tc.client.PutEmailIdentityMailFromAttributes(tc.ctx, &sesv2.PutEmailIdentityMailFromAttributesInput{
			EmailIdentity:       aws.String(domainIdentity),
			MailFromDomain:      aws.String(mailFrom),
			BehaviorOnMxFailure: types.BehaviorOnMxFailureUseDefaultValue,
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		if err != nil {
			return fmt.Errorf("get after mail-from put: %v", err)
		}
		if resp.MailFromAttributes == nil {
			return fmt.Errorf("MailFromAttributes is nil after put")
		}
		if resp.MailFromAttributes.MailFromDomain == nil || *resp.MailFromAttributes.MailFromDomain != mailFrom {
			return fmt.Errorf("expected MailFromDomain=%s, got %v", mailFrom, resp.MailFromAttributes.MailFromDomain)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "PutEmailIdentityConfigurationSetAttributes", func() error {
		csName := fmt.Sprintf("test-cs2-%d", tc.uid)
		_, _ = tc.client.CreateConfigurationSet(tc.ctx, &sesv2.CreateConfigurationSetInput{
			ConfigurationSetName: aws.String(csName),
		})
		_, err := tc.client.PutEmailIdentityConfigurationSetAttributes(tc.ctx, &sesv2.PutEmailIdentityConfigurationSetAttributesInput{
			EmailIdentity:        aws.String(emailAddress),
			ConfigurationSetName: aws.String(csName),
		})
		if err != nil {
			tc.deleteConfigSet(csName)
			return err
		}
		resp, err := tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(emailAddress),
		})
		tc.deleteConfigSet(csName)
		if err != nil {
			return fmt.Errorf("get after config set attr put: %v", err)
		}
		if resp.ConfigurationSetName == nil || *resp.ConfigurationSetName != csName {
			return fmt.Errorf("expected ConfigurationSetName=%s, got %v", csName, resp.ConfigurationSetName)
		}
		return nil
	}))

	policyName := fmt.Sprintf("test-policy-%d", tc.uid)

	results = append(results, r.RunTest("sesv2", "CreateEmailIdentityPolicy", func() error {
		_, err := tc.client.CreateEmailIdentityPolicy(tc.ctx, &sesv2.CreateEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
			Policy:        aws.String(`{"Version":"2008-10-17","Statement":[]}`),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return fmt.Errorf("get policies after create: %v", err)
		}
		if _, ok := resp.Policies[policyName]; !ok {
			return fmt.Errorf("policy %s not found after creation", policyName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "GetEmailIdentityPolicies", func() error {
		resp, err := tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return err
		}
		if len(resp.Policies) == 0 {
			return fmt.Errorf("expected at least 1 policy")
		}
		policy, ok := resp.Policies[policyName]
		if !ok {
			return fmt.Errorf("policy %s not found", policyName)
		}
		if policy != `{"Version":"2008-10-17","Statement":[]}` {
			return fmt.Errorf("policy content mismatch: %s", policy)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "UpdateEmailIdentityPolicy", func() error {
		updatedPolicy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"SES:SendEmail","Resource":"*"}]}`
		_, err := tc.client.UpdateEmailIdentityPolicy(tc.ctx, &sesv2.UpdateEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
			Policy:        aws.String(updatedPolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return fmt.Errorf("get policies after update: %v", err)
		}
		if resp.Policies[policyName] != updatedPolicy {
			return fmt.Errorf("policy not updated correctly")
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentityPolicy", func() error {
		_, err := tc.client.DeleteEmailIdentityPolicy(tc.ctx, &sesv2.DeleteEmailIdentityPolicyInput{
			EmailIdentity: aws.String(emailAddress),
			PolicyName:    aws.String(policyName),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetEmailIdentityPolicies(tc.ctx, &sesv2.GetEmailIdentityPoliciesInput{
			EmailIdentity: aws.String(emailAddress),
		})
		if err != nil {
			return fmt.Errorf("get policies after delete: %v", err)
		}
		if _, ok := resp.Policies[policyName]; ok {
			return fmt.Errorf("policy %s should have been deleted", policyName)
		}
		return nil
	}))

	results = append(results, r.RunTest("sesv2", "DeleteEmailIdentity_Domain", func() error {
		_, err := tc.client.DeleteEmailIdentity(tc.ctx, &sesv2.DeleteEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetEmailIdentity(tc.ctx, &sesv2.GetEmailIdentityInput{
			EmailIdentity: aws.String(domainIdentity),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting domain identity")
		}
		return nil
	}))

	return results
}
