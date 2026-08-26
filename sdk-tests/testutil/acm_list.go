package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMListTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "ListCertificates", func() error {
		resp, err := tc.client.ListCertificates(tc.ctx, &acm.ListCertificatesInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.CertificateSummaryList == nil {
			return fmt.Errorf("certificate summary list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_Pagination", func() error {
		var arns []string
		for i := 0; i < 3; i++ {
			arn, err := tc.requestDNSCert(fmt.Sprintf("page-%d-%d.com", time.Now().UnixNano(), i))
			if err != nil {
				return err
			}
			arns = append(arns, arn)
		}
		defer func() {
			for _, arn := range arns {
				tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{CertificateArn: aws.String(arn)})
			}
		}()

		allCerts, err := tc.allCertificates(nil)
		if err != nil {
			return err
		}
		for _, want := range arns {
			if containsID(allCerts, func(s *types.CertificateSummary) bool {
				return aws.ToString(s.CertificateArn) == want
			}) == nil {
				return fmt.Errorf("requested cert %s not found across pages", want)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_CertificateStatusesFilter", func() error {
		importArn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(importArn)

		certs, err := tc.allCertificates([]types.CertificateStatus{types.CertificateStatusIssued})
		if err != nil {
			return err
		}
		if containsID(certs, func(s *types.CertificateSummary) bool {
			return aws.ToString(s.CertificateArn) == importArn
		}) == nil {
			return fmt.Errorf("imported ISSUED cert not found in filtered list")
		}
		for _, s := range certs {
			if s.Status != types.CertificateStatusIssued {
				return fmt.Errorf("found non-ISSUED cert in filtered list: %s", s.Status)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_SummaryFields", func() error {
		domain := acmUniqueDomain("summary")
		arn, err := tc.requestDNSCert(domain)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		certs, err := tc.allCertificates(nil)
		if err != nil {
			return err
		}
		found := containsID(certs, func(s *types.CertificateSummary) bool {
			return aws.ToString(s.CertificateArn) == arn
		})
		if found == nil {
			return fmt.Errorf("certificate not found in list")
		}
		if aws.ToString(found.DomainName) != domain {
			return fmt.Errorf("expected DomainName %s, got %s", domain, aws.ToString(found.DomainName))
		}
		if found.Status != types.CertificateStatusIssued {
			return fmt.Errorf("expected ISSUED status, got %s", found.Status)
		}
		if found.Type != types.CertificateTypeAmazonIssued {
			return fmt.Errorf("expected AMAZON_ISSUED type, got %s", found.Type)
		}
		return nil
	}))

	return results
}
