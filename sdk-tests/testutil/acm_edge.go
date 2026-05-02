package testutil

import (
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

	return results
}
