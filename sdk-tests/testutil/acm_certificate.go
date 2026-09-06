package testutil

import (
	"crypto/x509"
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
		arn, err := tc.requestCertSANs(domain, fmt.Sprintf("www.%s", domain), fmt.Sprintf("api.%s", domain))
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if len(cert.SubjectAlternativeNames) != 2 {
			return fmt.Errorf("expected 2 SANs, got %d", len(cert.SubjectAlternativeNames))
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithOptions", func() error {
		domain := acmUniqueDomain("opts-test")
		arn, err := tc.requestCert(&acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if cert.Options == nil {
			return fmt.Errorf("options is nil")
		}
		if cert.Options.CertificateTransparencyLoggingPreference != types.CertificateTransparencyLoggingPreferenceDisabled {
			return fmt.Errorf("expected DISABLED, got %s", cert.Options.CertificateTransparencyLoggingPreference)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_WithEmailValidation", func() error {
		domain := acmUniqueDomain("email-test")
		arn, err := tc.requestEmailCert(domain)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if len(cert.DomainValidationOptions) != 1 {
			return fmt.Errorf("expected 1 DVO, got %d", len(cert.DomainValidationOptions))
		}
		dvo := cert.DomainValidationOptions[0]
		if dvo.ValidationMethod != types.ValidationMethodEmail {
			return fmt.Errorf("expected EMAIL validation method, got %s", dvo.ValidationMethod)
		}
		if dvo.ResourceRecord != nil {
			return fmt.Errorf("EMAIL validation should not have ResourceRecord")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RequestCertificate_VerifyArnFormat", func() error {
		arn, _, err := tc.requestOwnDNSCert("arn-test")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		if !strings.HasPrefix(arn, "arn:aws:acm:") {
			return fmt.Errorf("ARN should start with arn:aws:acm:, got %s", arn)
		}
		if !strings.Contains(arn, "certificate/") {
			return fmt.Errorf("ARN should contain certificate/, got %s", arn)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "DescribeCertificate_AMAZON_ISSUED_Fields", func() error {
		arn, domain, err := tc.requestOwnDNSCert("desc-ai")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		c, err := tc.describeCert(arn)
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
		arn, domain, err := tc.requestOwnDNSCert("dv-dns")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if len(cert.DomainValidationOptions) != 1 {
			return fmt.Errorf("expected 1 DVO, got %d", len(cert.DomainValidationOptions))
		}
		dvo := cert.DomainValidationOptions[0]
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
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		c, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
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
		arn, err := tc.importCertWithChain()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{CertificateArn: aws.String(arn)})
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
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)
		if !strings.HasPrefix(arn, "arn:aws:acm:") {
			return fmt.Errorf("ARN format incorrect: %s", arn)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ExportCertificate", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		exportResp, err := tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: aws.String(arn),
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

	// Verify SANs are embedded in the actual x509 certificate
	results = append(results, r.RunTest("acm", "RequestCertificate_SANsEmbeddedInCert", func() error {
		domain := acmUniqueDomain("san-embed")
		san1 := fmt.Sprintf("www.%s", domain)
		san2 := fmt.Sprintf("api.%s", domain)
		arn, err := tc.requestCertSANs(domain, san1, san2)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		parsed, err := parsePEMCertificate(aws.ToString(getResp.Certificate))
		if err != nil {
			return err
		}
		dnsSet := make(map[string]bool)
		for _, d := range parsed.DNSNames {
			dnsSet[d] = true
		}
		if !dnsSet[domain] {
			return fmt.Errorf("primary domain %s not in cert DNSNames: %v", domain, parsed.DNSNames)
		}
		if !dnsSet[san1] {
			return fmt.Errorf("SAN %s not in cert DNSNames: %v", san1, parsed.DNSNames)
		}
		if !dnsSet[san2] {
			return fmt.Errorf("SAN %s not in cert DNSNames: %v", san2, parsed.DNSNames)
		}
		return nil
	}))

	// KeyAlgorithm EC_prime256v1
	results = append(results, r.RunTest("acm", "RequestCertificate_KeyAlgorithm_EC_prime256v1", func() error {
		domain := acmUniqueDomain("ec256")
		arn, err := tc.requestCert(&acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			KeyAlgorithm:     types.KeyAlgorithmEcPrime256v1,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if cert.KeyAlgorithm != types.KeyAlgorithmEcPrime256v1 {
			return fmt.Errorf("expected EC_prime256v1, got %s", cert.KeyAlgorithm)
		}
		// Verify the actual cert uses ECDSA
		getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		parsed, err := parsePEMCertificate(aws.ToString(getResp.Certificate))
		if err != nil {
			return err
		}
		if parsed.PublicKeyAlgorithm != x509.ECDSA {
			return fmt.Errorf("expected ECDSA public key, got %s", parsed.PublicKeyAlgorithm)
		}
		return nil
	}))

	// KeyAlgorithm RSA_4096
	results = append(results, r.RunTest("acm", "RequestCertificate_KeyAlgorithm_RSA_4096", func() error {
		domain := acmUniqueDomain("rsa4096")
		arn, err := tc.requestCert(&acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			KeyAlgorithm:     types.KeyAlgorithmRsa4096,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		if cert.KeyAlgorithm != types.KeyAlgorithmRsa4096 {
			return fmt.Errorf("expected RSA_4096, got %s", cert.KeyAlgorithm)
		}
		return nil
	}))

	return results
}
