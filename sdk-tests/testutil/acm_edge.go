package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMEdgeTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "DescribeCertificate_NonExistent", func() error {
		_, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "GetCertificate_NonExistent", func() error {
		_, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_NonExistent", func() error {
		_, err := tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate_NonExistent", func() error {
		_, err := tc.client.RemoveTagsFromCertificate(tc.ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Tags:           []types.Tag{{Key: aws.String("X")}},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "ListTagsForCertificate_NonExistent", func() error {
		_, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "ExportCertificate_AMAZON_ISSUED_Error", func() error {
		domain := acmUniqueDomain("export-ai")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: resp.CertificateArn,
			Passphrase:     []byte("test-passphrase"),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// B3: ImportCertificate without PrivateKey (initial import) should fail.
	// The AWS SDK enforces PrivateKey as client-side required (nil check),
	// but an empty byte slice bypasses it. The server should still reject.
	results = append(results, r.RunTest("acm", "ImportCertificate_EmptyPrivateKey", func() error {
		_, err := tc.client.ImportCertificate(tc.ctx, &acm.ImportCertificateInput{
			Certificate: testCertPEM,
			PrivateKey:  []byte{},
		})
		if err == nil {
			return fmt.Errorf("expected error for empty PrivateKey, got nil")
		}
		return nil
	}))

	// B4: Invalid domain name should fail
	results = append(results, r.RunTest("acm", "RequestCertificate_InvalidDomainName", func() error {
		_, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String("not-a-valid-domain"),
			ValidationMethod: types.ValidationMethodDns,
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// B7: Passphrase too short (< 4 bytes)
	results = append(results, r.RunTest("acm", "ExportCertificate_PassphraseTooShort", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: importResp.CertificateArn,
			Passphrase:     []byte("ab"),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	return results
}
