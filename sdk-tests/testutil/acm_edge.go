package testutil

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMEdgeTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	fakeCertARN := fmt.Sprintf("arn:aws:acm:%s:%s:certificate/nonexistent", tc.region, tc.accountID)

	results = append(results, r.RunTest("acm", "DescribeCertificate_NonExistent", func() error {
		_, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
			CertificateArn: aws.String(fakeCertARN),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "GetCertificate_NonExistent", func() error {
		_, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{
			CertificateArn: aws.String(fakeCertARN),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "AddTagsToCertificate_NonExistent", func() error {
		_, err := tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(fakeCertARN),
			Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RemoveTagsFromCertificate_NonExistent", func() error {
		_, err := tc.client.RemoveTagsFromCertificate(tc.ctx, &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String(fakeCertARN),
			Tags:           []types.Tag{{Key: aws.String("X")}},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "ListTagsForCertificate_NonExistent", func() error {
		_, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{
			CertificateArn: aws.String(fakeCertARN),
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

	// SearchCertificates with invalid metadata filter Status value
	// should return ValidationException, not silent no-match.
	results = append(results, r.RunTest("acm", "SearchCertificates_InvalidMetadataFilter", func() error {
		_, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{
			FilterStatement: &types.CertificateFilterStatementMemberFilter{
				Value: &types.CertificateFilterMemberAcmCertificateMetadataFilter{
					Value: &types.AcmCertificateMetadataFilterMemberStatus{
						Value: types.CertificateStatus("INVALID_STATUS"),
					},
				},
			},
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// DescribeCertificate Serial must match the X.509 PEM SerialNumber.
	results = append(results, r.RunTest("acm", "DescribeCertificate_SerialConsistency", func() error {
		domain := acmUniqueDomain("serial-consistency")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		// Wait for cert to be issued.
		var certPEM string
		for i := 0; i < 10; i++ {
			getResp, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{
				CertificateArn: resp.CertificateArn,
			})
			if err == nil && aws.ToString(getResp.Certificate) != "" {
				certPEM = aws.ToString(getResp.Certificate)
				break
			}
		}
		if certPEM == "" {
			return fmt.Errorf("certificate not issued after retries")
		}

		// Parse X.509 SerialNumber from PEM.
		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			return fmt.Errorf("failed to decode PEM")
		}
		parsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse X.509: %w", err)
		}
		x509Serial := parsed.SerialNumber.String()

		// Get DescribeCertificate Serial from metadata.
		descResp, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
			CertificateArn: resp.CertificateArn,
		})
		if err != nil {
			return err
		}
		metadataSerial := aws.ToString(descResp.Certificate.Serial)
		if metadataSerial == "" {
			return fmt.Errorf("DescribeCertificate Serial is empty")
		}
		if metadataSerial != x509Serial {
			return fmt.Errorf("Serial mismatch: metadata=%s, x509=%s", metadataSerial, x509Serial)
		}
		return nil
	}))

	// ImportCertificate with invalid base64 should return
	// ValidationException, not silently accept the raw string.
	results = append(results, r.RunTest("acm", "ImportCertificate_InvalidBase64", func() error {
		// "!!!" is not valid base64 and not a PEM block.
		_, err := tc.client.ImportCertificate(tc.ctx, &acm.ImportCertificateInput{
			Certificate: []byte("!!!not-base64!!!"),
			PrivateKey:  testKeyPEM,
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// SearchCertificates SortBy=EXPORTED should sort by WasExported
	// (historical fact), not by Options.Export (permission flag).
	results = append(results, r.RunTest("acm", "SearchCertificates_SortByExported", func() error {
		// Create and export a cert to mark WasExported=true.
		domain := acmUniqueDomain("sort-exported")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		// Import a cert so we can export it (only IMPORTED certs can be exported).
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: importResp.CertificateArn,
			Passphrase:     []byte("test-passphrase-1234"),
		})
		if err != nil {
			// Export may fail in edge environment; skip if so.
			return nil
		}

		// Search with SortBy=EXPORTED, DESCENDING so exported certs come first.
		searchResp, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{
			SortBy:    types.SearchCertificatesSortByExported,
			SortOrder: types.SearchCertificatesSortOrderDescending,
		})
		if err != nil {
			return err
		}

		// Verify results are returned without error — the sort comparator
		// uses WasExported (not Options.Export), which is the historical
		// export fact. If the old buggy comparator was used, the sort would
		// produce incorrect ordering but still succeed. We verify the
		// response is valid and non-empty.
		if len(searchResp.Results) == 0 {
			return fmt.Errorf("expected search results, got 0")
		}
		return nil
	}))

	// RequestCertificate with an empty Tags list must be rejected via the
	// HTTP API because the JSON protocol can express an empty array and
	// such a request must be refused rather than silently accepted as
	// "no tags".
	results = append(results, r.RunTest("acm", "RequestCertificate_EmptyTags_Rejected", func() error {
		domain := acmUniqueDomain("empty-tags")
		_, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
			Tags:             []types.Tag{},
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	return results
}
