package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMCertificateTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "RequestCertificate_WithSubjectAlternativeNames", func() error {
		domain := acmUniqueDomain("san-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:              aws.String(domain),
			ValidationMethod:        types.ValidationMethodDns,
			SubjectAlternativeNames: []string{fmt.Sprintf("www.%s", domain), fmt.Sprintf("api.%s", domain)},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(desc.Certificate.SubjectAlternativeNames) != 2 {
			return fmt.Errorf("expected 2 SANs, got %d", len(desc.Certificate.SubjectAlternativeNames))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithOptions", func() error {
		domain := acmUniqueDomain("opts-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.Options == nil {
			return fmt.Errorf("options is nil")
		}
		if desc.Certificate.Options.CertificateTransparencyLoggingPreference != types.CertificateTransparencyLoggingPreferenceDisabled {
			return fmt.Errorf("expected DISABLED, got %s", desc.Certificate.Options.CertificateTransparencyLoggingPreference)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithEmailValidation", func() error {
		domain := acmUniqueDomain("email-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodEmail,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(desc.Certificate.DomainValidationOptions) != 1 {
			return fmt.Errorf("expected 1 DVO, got %d", len(desc.Certificate.DomainValidationOptions))
		}
		dvo := desc.Certificate.DomainValidationOptions[0]
		if dvo.ValidationMethod != types.ValidationMethodEmail {
			return fmt.Errorf("expected EMAIL validation method, got %s", dvo.ValidationMethod)
		}
		if dvo.ResourceRecord != nil {
			return fmt.Errorf("EMAIL validation should not have ResourceRecord")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_VerifyArnFormat", func() error {
		domain := acmUniqueDomain("arn-test")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		arn := aws.ToString(resp.CertificateArn)
		if !strings.HasPrefix(arn, "arn:aws:acm:") {
			return fmt.Errorf("ARN should start with arn:aws:acm:, got %s", arn)
		}
		if !strings.Contains(arn, "certificate/") {
			return fmt.Errorf("ARN should contain certificate/, got %s", arn)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_AMAZON_ISSUED_Fields", func() error {
		domain := acmUniqueDomain("desc-ai")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		c := desc.Certificate
		if c.Status != types.CertificateStatusIssued {
			return fmt.Errorf("expected ISSUED, got %s", c.Status)
		}
		if c.Type != types.CertificateTypeAmazonIssued {
			return fmt.Errorf("expected AMAZON_ISSUED, got %s", c.Type)
		}
		if c.RenewalEligibility != types.RenewalEligibilityEligible {
			return fmt.Errorf("expected ELIGIBLE, got %s", c.RenewalEligibility)
		}
		if aws.ToString(c.DomainName) != domain {
			return fmt.Errorf("expected domain %s, got %s", domain, aws.ToString(c.DomainName))
		}
		if c.KeyAlgorithm != types.KeyAlgorithmRsa2048 {
			return fmt.Errorf("expected RSA_2048, got %s", c.KeyAlgorithm)
		}
		if c.CertificateArn == nil || !strings.Contains(aws.ToString(c.CertificateArn), "acm") {
			return fmt.Errorf("CertificateArn missing or malformed: %s", aws.ToString(c.CertificateArn))
		}
		if c.CreatedAt == nil {
			return fmt.Errorf("CreatedAt is nil")
		}
		if c.Serial == nil || aws.ToString(c.Serial) == "" {
			return fmt.Errorf("Serial is nil or empty")
		}
		if c.Subject == nil || aws.ToString(c.Subject) == "" {
			return fmt.Errorf("Subject is nil or empty")
		}
		if c.Issuer == nil || aws.ToString(c.Issuer) == "" {
			return fmt.Errorf("Issuer is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_DomainValidationOptions_DNS", func() error {
		domain := acmUniqueDomain("dv-dns")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if len(desc.Certificate.DomainValidationOptions) != 1 {
			return fmt.Errorf("expected 1 DVO, got %d", len(desc.Certificate.DomainValidationOptions))
		}
		dvo := desc.Certificate.DomainValidationOptions[0]
		if dvo.ValidationMethod != types.ValidationMethodDns {
			return fmt.Errorf("expected DNS, got %s", dvo.ValidationMethod)
		}
		if dvo.ResourceRecord == nil {
			return fmt.Errorf("ResourceRecord is nil")
		}
		if dvo.ResourceRecord.Type != "CNAME" {
			return fmt.Errorf("expected CNAME, got %s", dvo.ResourceRecord.Type)
		}
		if !strings.Contains(aws.ToString(dvo.ResourceRecord.Name), domain) {
			return fmt.Errorf("ResourceRecord.Name should contain domain, got %s", aws.ToString(dvo.ResourceRecord.Name))
		}
		if aws.ToString(dvo.ResourceRecord.Value) == "" {
			return fmt.Errorf("ResourceRecord.Value is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_IMPORTED_Fields", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		c := desc.Certificate
		if c.Status != types.CertificateStatusIssued {
			return fmt.Errorf("expected ISSUED, got %s", c.Status)
		}
		if c.Type != types.CertificateTypeImported {
			return fmt.Errorf("expected IMPORTED, got %s", c.Type)
		}
		if c.RenewalEligibility != types.RenewalEligibilityIneligible {
			return fmt.Errorf("expected INELIGIBLE, got %s", c.RenewalEligibility)
		}
		if c.ImportedAt == nil {
			return fmt.Errorf("ImportedAt is nil")
		}
		if c.NotBefore == nil || c.NotAfter == nil {
			return fmt.Errorf("NotBefore or NotAfter is nil")
		}
		if c.KeyAlgorithm != types.KeyAlgorithmRsa2048 {
			return fmt.Errorf("expected RSA_2048, got %s", c.KeyAlgorithm)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "GetCertificate_ImportedCert", func() error {
		importResp, err := tc.importCertWithChain()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.Certificate) == "" {
			return fmt.Errorf("certificate is empty")
		}
		if !strings.Contains(aws.ToString(getResp.Certificate), "-----BEGIN CERTIFICATE-----") {
			return fmt.Errorf("certificate should be PEM encoded")
		}
		if aws.ToString(getResp.CertificateChain) == "" {
			return fmt.Errorf("CertificateChain is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ImportCertificate", func() error {
		resp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))
		if resp.CertificateArn == nil {
			return fmt.Errorf("CertificateArn is nil")
		}
		if !strings.HasPrefix(aws.ToString(resp.CertificateArn), "arn:aws:acm:") {
			return fmt.Errorf("ARN format incorrect: %s", aws.ToString(resp.CertificateArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ImportCertificate_WithCertificateChain", func() error {
		resp, err := tc.importCertWithChain()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{CertificateArn: resp.CertificateArn})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.CertificateChain) == "" {
			return fmt.Errorf("CertificateChain is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ExportCertificate", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		exportResp, err := tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: importResp.CertificateArn,
			Passphrase:     []byte("test-passphrase"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(exportResp.Certificate) == "" {
			return fmt.Errorf("certificate is empty")
		}
		if exportResp.PrivateKey == nil || aws.ToString(exportResp.PrivateKey) == "" {
			return fmt.Errorf("PrivateKey is empty")
		}
		return nil
	}))

	return results
}
