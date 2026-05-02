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
			resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
				DomainName:       aws.String(fmt.Sprintf("page-%d-%d.com", time.Now().UnixNano(), i)),
				ValidationMethod: types.ValidationMethodDns,
			})
			if err != nil {
				return err
			}
			arns = append(arns, aws.ToString(resp.CertificateArn))
		}
		defer func() {
			for _, arn := range arns {
				tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{CertificateArn: aws.String(arn)})
			}
		}()

		var allArns []string
		var nextToken *string
		for {
			input := &acm.ListCertificatesInput{
				MaxItems: aws.Int32(2),
			}
			if nextToken != nil {
				input.NextToken = nextToken
			}
			resp, err := tc.client.ListCertificates(tc.ctx, input)
			if err != nil {
				return fmt.Errorf("list page: %v", err)
			}
			for _, s := range resp.CertificateSummaryList {
				allArns = append(allArns, aws.ToString(s.CertificateArn))
			}
			if resp.NextToken == nil || aws.ToString(resp.NextToken) == "" {
				break
			}
			nextToken = resp.NextToken
		}

		foundCount := 0
		arnSet := make(map[string]bool)
		for _, a := range arns {
			arnSet[a] = true
		}
		for _, a := range allArns {
			if arnSet[a] {
				foundCount++
			}
		}
		if foundCount != 3 {
			return fmt.Errorf("expected to find all 3 certs across pages, found %d", foundCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_CertificateStatusesFilter", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		importArn := aws.ToString(importResp.CertificateArn)
		defer tc.deleteCert(importArn)

		resp, err := tc.client.ListCertificates(tc.ctx, &acm.ListCertificatesInput{
			CertificateStatuses: []types.CertificateStatus{types.CertificateStatusIssued},
		})
		if err != nil {
			return err
		}
		found := false
		for _, s := range resp.CertificateSummaryList {
			if aws.ToString(s.CertificateArn) == importArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("imported ISSUED cert not found in filtered list")
		}
		for _, s := range resp.CertificateSummaryList {
			if s.Status != types.CertificateStatusIssued {
				return fmt.Errorf("found non-ISSUED cert in filtered list: %s", s.Status)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ListCertificates_SummaryFields", func() error {
		domain := acmUniqueDomain("summary")
		reqResp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		arn := aws.ToString(reqResp.CertificateArn)
		defer tc.deleteCert(arn)

		var found *types.CertificateSummary
		var nextToken *string
		for {
			input := &acm.ListCertificatesInput{MaxItems: aws.Int32(100)}
			if nextToken != nil {
				input.NextToken = nextToken
			}
			listResp, err := tc.client.ListCertificates(tc.ctx, input)
			if err != nil {
				return err
			}
			for _, s := range listResp.CertificateSummaryList {
				if aws.ToString(s.CertificateArn) == arn {
					cs := s
					found = &cs
					break
				}
			}
			if found != nil || listResp.NextToken == nil {
				break
			}
			nextToken = listResp.NextToken
		}
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
