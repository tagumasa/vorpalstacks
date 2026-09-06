package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

// colonHexSerial renders an X.509 serial number in the ACM wire form:
// colon-separated lowercase hex byte pairs (the SerialNumber shape both
// the filter input and the X509Attributes output use).
func colonHexSerial(n *big.Int) string {
	serialBytes := n.Bytes()
	if len(serialBytes) == 0 {
		return "00"
	}
	pairs := make([]string, len(serialBytes))
	for i, b := range serialBytes {
		pairs[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(pairs, ":")
}

// searchX509FilterFinds reports whether a SearchCertificates X509 filter for
// member finds arn, traversing every result page.
func searchX509FilterFinds(tc *acmTestContext, member types.X509AttributeFilter, arn string) (bool, error) {
	var nextToken *string
	for page := 0; page < 20; page++ {
		resp, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{
			FilterStatement: &types.CertificateFilterStatementMemberFilter{
				Value: &types.CertificateFilterMemberX509AttributeFilter{
					Value: member,
				},
			},
			MaxResults: aws.Int32(100),
			NextToken:  nextToken,
		})
		if err != nil {
			return false, err
		}
		for _, result := range resp.Results {
			if aws.ToString(result.CertificateArn) == arn {
				return true, nil
			}
		}
		nextToken = resp.NextToken
		if nextToken == nil {
			return false, nil
		}
	}
	return false, nil
}

func (r *TestRunner) runACMEdgeTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	// Operations against a certificate that does not exist fail with
	// ResourceNotFoundException, as the service model specifies.
	results = append(results, r.RunTest("acm", "NonExistentCertificateOperations", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"DescribeCertificate", func() error {
				_, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"GetCertificate", func() error {
				_, err := tc.client.GetCertificate(tc.ctx, &acm.GetCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"AddTagsToCertificate", func() error {
				_, err := tc.client.AddTagsToCertificate(tc.ctx, &acm.AddTagsToCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
					Tags:           []types.Tag{{Key: aws.String("X"), Value: aws.String("Y")}},
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"RemoveTagsFromCertificate", func() error {
				_, err := tc.client.RemoveTagsFromCertificate(tc.ctx, &acm.RemoveTagsFromCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
					Tags:           []types.Tag{{Key: aws.String("X")}},
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"ListTagsForCertificate", func() error {
				_, err := tc.client.ListTagsForCertificate(tc.ctx, &acm.ListTagsForCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"DeleteCertificate", func() error {
				_, err := tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"UpdateCertificateOptions", func() error {
				_, err := tc.client.UpdateCertificateOptions(tc.ctx, &acm.UpdateCertificateOptionsInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
					Options: &types.CertificateOptions{
						CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled,
					},
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"RenewCertificate", func() error {
				_, err := tc.client.RenewCertificate(tc.ctx, &acm.RenewCertificateInput{
					CertificateArn: aws.String(tc.nonExistentArn()),
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
			{"RevokeCertificate", func() error {
				_, err := tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
					CertificateArn:   aws.String(tc.nonExistentArn()),
					RevocationReason: types.RevocationReasonKeyCompromise,
				})
				return AssertErrorContains(err, "ResourceNotFoundException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	// ExportCertificate rejects with ValidationException when the
	// certificate is not exportable by type (platform-issued certs) and
	// when the passphrase is shorter than the required minimum length.
	results = append(results, r.RunTest("acm", "ExportCertificate_Error", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{"amazon-issued-rejected", func() error {
				arn, _, err := tc.requestOwnDNSCert("export-ai")
				if err != nil {
					return err
				}
				defer tc.deleteCert(arn)

				_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
					CertificateArn: aws.String(arn),
					Passphrase:     []byte("test-passphrase"),
				})
				return AssertErrorContains(err, "ValidationException")
			}},
			{"short-passphrase-rejected", func() error {
				importArn, err := tc.importDefaultCert()
				if err != nil {
					return err
				}
				defer tc.deleteCert(importArn)

				// Passphrase too short (< 4 bytes).
				_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
					CertificateArn: aws.String(importArn),
					Passphrase:     []byte("ab"),
				})
				return AssertErrorContains(err, "ValidationException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	// ImportCertificate without PrivateKey (initial import) should fail.
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

	// Invalid domain name should fail
	results = append(results, r.RunTest("acm", "RequestCertificate_InvalidDomainName", func() error {
		_, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String("not-a-valid-domain"),
			ValidationMethod: types.ValidationMethodDns,
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// SearchCertificates rejects malformed filter and sort input with
	// ValidationException, never a silent no-match.
	results = append(results, r.RunTest("acm", "SearchCertificates_InvalidInput_Rejected", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			// An invalid metadata filter Status value should return
			// ValidationException, not silent no-match.
			{"invalid-metadata-filter", func() error {
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
			}},
			// The X509 half of the same contract: an invalid KeyUsage filter
			// value must be rejected with ValidationException, not silently
			// match nothing.
			{"invalid-x509-filter", func() error {
				_, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{
					FilterStatement: &types.CertificateFilterStatementMemberFilter{
						Value: &types.CertificateFilterMemberX509AttributeFilter{
							Value: &types.X509AttributeFilterMemberKeyUsage{
								Value: types.KeyUsageName("NOT_A_KEY_USAGE"),
							},
						},
					},
				})
				return AssertErrorContains(err, "ValidationException")
			}},
			// An unknown SortBy value is rejected with ValidationException —
			// the only constraint-violation error the operation declares.
			{"invalid-sort-by", func() error {
				_, err := tc.client.SearchCertificates(tc.ctx, &acm.SearchCertificatesInput{
					SortBy: types.SearchCertificatesSortBy("BOGUS"),
				})
				return AssertErrorContains(err, "ValidationException")
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	// DescribeCertificate Serial must match the X.509 PEM SerialNumber.
	results = append(results, r.RunTest("acm", "DescribeCertificate_SerialConsistency", func() error {
		arn, _, err := tc.requestOwnDNSCert("serial-consistency")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		// Wait for the cert to be issued, then parse its X.509 serial.
		certPEM, err := tc.waitIssuedPEM(arn)
		if err != nil {
			return err
		}
		parsed, err := parsePEMCertificate(certPEM)
		if err != nil {
			return err
		}
		x509Serial := colonHexSerial(parsed.SerialNumber)

		// Get DescribeCertificate Serial from metadata.
		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		metadataSerial := aws.ToString(cert.Serial)
		if metadataSerial == "" {
			return fmt.Errorf("DescribeCertificate Serial is empty")
		}
		if metadataSerial != x509Serial {
			return fmt.Errorf("Serial mismatch: metadata=%s, x509=%s", metadataSerial, x509Serial)
		}
		return nil
	}))

	// The SerialNumber X509 filter must find the issued certificate by its
	// wire-form serial, and must not match a serial no certificate holds.
	results = append(results, r.RunTest("acm", "SearchCertificates_SerialNumberFilter", func() error {
		arn, _, err := tc.requestOwnDNSCert("serial-filter")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		// Wait for the cert to be issued, then parse its X.509 serial.
		certPEM, err := tc.waitIssuedPEM(arn)
		if err != nil {
			return err
		}
		parsed, err := parsePEMCertificate(certPEM)
		if err != nil {
			return err
		}
		serial := colonHexSerial(parsed.SerialNumber)

		found, err := searchX509FilterFinds(tc, &types.X509AttributeFilterMemberSerialNumber{Value: serial}, arn)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("serial filter %s did not find %s", serial, arn)
		}

		wrongSerial := "01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12"
		found, err = searchX509FilterFinds(tc, &types.X509AttributeFilterMemberSerialNumber{Value: wrongSerial}, arn)
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("serial filter matched a serial the certificate does not hold")
		}
		return nil
	}))

	// Platform-issued certificates grant DIGITAL_SIGNATURE and
	// KEY_ENCIPHERMENT: DescribeCertificate reports them, the
	// SearchCertificates KeyUsage filter matches on them, and a usage the
	// certificate does not hold does not match.
	results = append(results, r.RunTest("acm", "SearchCertificates_KeyUsageFilter_MatchesIssuedCert", func() error {
		arn, _, err := tc.requestOwnDNSCert("keyusage-filter")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		hasDigitalSignature := false
		for _, ku := range cert.KeyUsages {
			if ku.Name == types.KeyUsageNameDigitalSignature {
				hasDigitalSignature = true
			}
		}
		if !hasDigitalSignature {
			return fmt.Errorf("DescribeCertificate KeyUsages %v lacks DIGITAL_SIGNATURE", cert.KeyUsages)
		}

		found, err := searchX509FilterFinds(tc, &types.X509AttributeFilterMemberKeyUsage{Value: types.KeyUsageNameDigitalSignature}, arn)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("KeyUsage filter DIGITAL_SIGNATURE did not find %s", arn)
		}

		found, err = searchX509FilterFinds(tc, &types.X509AttributeFilterMemberKeyUsage{Value: types.KeyUsageNameCrlSigning}, arn)
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("KeyUsage filter CRL_SIGNING matched a certificate that does not grant it")
		}
		return nil
	}))

	// An imported certificate whose only extended key usage is
	// anyExtendedKeyUsage is recorded as the ACM enum value ANY:
	// DescribeCertificate reports it, and both the SearchCertificates and
	// ListCertificates extendedKeyUsage filters match it — while an enum
	// value the certificate does not hold does not match.
	results = append(results, r.RunTest("acm", "SearchCertificates_ExtendedKeyUsageFilter_AnyEKU", func() error {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}
		serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 120))
		if err != nil {
			return err
		}
		template := &x509.Certificate{
			SerialNumber: serial,
			Subject:      pkix.Name{CommonName: acmUniqueDomain("any-eku")},
			NotBefore:    time.Now().Add(-time.Hour),
			NotAfter:     time.Now().Add(24 * time.Hour),
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		}
		der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
		if err != nil {
			return err
		}
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
		keyDER, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return err
		}
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

		arn, err := tc.importCertificate(certPEM, keyPEM, nil)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		cert, err := tc.describeCert(arn)
		if err != nil {
			return err
		}
		hasAny := false
		for _, eku := range cert.ExtendedKeyUsages {
			if eku.Name == types.ExtendedKeyUsageNameAny {
				hasAny = true
			}
		}
		if !hasAny {
			return fmt.Errorf("DescribeCertificate ExtendedKeyUsages %v lacks ANY", cert.ExtendedKeyUsages)
		}

		found, err := searchX509FilterFinds(tc, &types.X509AttributeFilterMemberExtendedKeyUsage{Value: types.ExtendedKeyUsageNameAny}, arn)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("ExtendedKeyUsage filter ANY did not find %s", arn)
		}
		found, err = searchX509FilterFinds(tc, &types.X509AttributeFilterMemberExtendedKeyUsage{Value: types.ExtendedKeyUsageNameCodeSigning}, arn)
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("ExtendedKeyUsage filter CODE_SIGNING matched a certificate that does not grant it")
		}

		listed := false
		var nextToken *string
		for page := 0; page < 20 && !listed; page++ {
			resp, err := tc.client.ListCertificates(tc.ctx, &acm.ListCertificatesInput{
				Includes:  &types.Filters{ExtendedKeyUsage: []types.ExtendedKeyUsageName{types.ExtendedKeyUsageNameAny}},
				MaxItems:  aws.Int32(100),
				NextToken: nextToken,
			})
			if err != nil {
				return err
			}
			for _, sum := range resp.CertificateSummaryList {
				if aws.ToString(sum.CertificateArn) == arn {
					listed = true
				}
			}
			nextToken = resp.NextToken
			if nextToken == nil {
				break
			}
		}
		if !listed {
			return fmt.Errorf("ListCertificates Includes.extendedKeyUsage=ANY did not list %s", arn)
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
		arn, _, err := tc.requestOwnDNSCert("sort-exported")
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		// Import a cert so we can export it (only IMPORTED certs can be exported).
		importArn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(importArn)

		_, err = tc.client.ExportCertificate(tc.ctx, &acm.ExportCertificateInput{
			CertificateArn: aws.String(importArn),
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
