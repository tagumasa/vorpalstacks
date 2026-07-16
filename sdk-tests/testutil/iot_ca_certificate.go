package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTCACertificateTests covers the CA certificate lifecycle
// (Register/Describe/Update/List/Delete) using the id returned by Register so
// the assertions target a known resource instead of ListCACertificates[0].
func (r *TestRunner) runIoTCACertificateTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	const caPem = "-----BEGIN CERTIFICATE-----\nMIICXTCCAgegAwIBAgIJANtestcacetemplate\n-----END CERTIFICATE-----"
	var caID string

	defer func() {
		if caID != "" {
			tc.client.DeleteCACertificate(tc.ctx, &iot.DeleteCACertificateInput{CertificateId: aws.String(caID)})
		}
	}()

	results = append(results, r.RunTest("iot", "CACert_RegisterCACertificate", func() error {
		out, err := tc.client.RegisterCACertificate(tc.ctx, &iot.RegisterCACertificateInput{
			CaCertificate:           aws.String(caPem),
			VerificationCertificate: aws.String("verification-cert"),
			SetAsActive:             true,
			AllowAutoRegistration:   false,
		})
		if err != nil {
			return fmt.Errorf("RegisterCACertificate failed: %w", err)
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("expected non-empty certificateId")
		}
		caID = *out.CertificateId
		return nil
	}))

	results = append(results, r.RunTest("iot", "CACert_DescribeCACertificate", func() error {
		out, err := tc.client.DescribeCACertificate(tc.ctx, &iot.DescribeCACertificateInput{CertificateId: aws.String(caID)})
		if err != nil {
			return fmt.Errorf("DescribeCACertificate failed: %w", err)
		}
		if out.CertificateDescription == nil {
			return fmt.Errorf("expected certificateDescription")
		}
		if out.CertificateDescription.CertificateId == nil || *out.CertificateDescription.CertificateId != caID {
			return fmt.Errorf("expected certificateId=%s, got %v", caID, out.CertificateDescription.CertificateId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "CACert_UpdateCACertificate", func() error {
		if _, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId: aws.String(caID),
			NewStatus:     iottypes.CACertificateStatusInactive,
		}); err != nil {
			return fmt.Errorf("UpdateCACertificate failed: %w", err)
		}
		out, err := tc.client.DescribeCACertificate(tc.ctx, &iot.DescribeCACertificateInput{CertificateId: aws.String(caID)})
		if err != nil {
			return fmt.Errorf("DescribeCACertificate after update failed: %w", err)
		}
		if out.CertificateDescription.Status != iottypes.CACertificateStatusInactive {
			return fmt.Errorf("expected status INACTIVE, got %s", out.CertificateDescription.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "CACert_ListCACertificates_IncludesCreated", func() error {
		out, err := tc.client.ListCACertificates(tc.ctx, &iot.ListCACertificatesInput{})
		if err != nil {
			return fmt.Errorf("ListCACertificates failed: %w", err)
		}
		for _, c := range out.Certificates {
			if c.CertificateId != nil && *c.CertificateId == caID {
				return nil
			}
		}
		return fmt.Errorf("CA certificate %s not found in list of %d", caID, len(out.Certificates))
	}))

	results = append(results, r.RunTest("iot", "CACert_DescribeCACertificate_NotFound", func() error {
		_, err := tc.client.DescribeCACertificate(tc.ctx, &iot.DescribeCACertificateInput{
			CertificateId: aws.String("0123456789012345678901234567890123456789012345678901234567890123"),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "CACert_DeleteCACertificate", func() error {
		if _, err := tc.client.DeleteCACertificate(tc.ctx, &iot.DeleteCACertificateInput{CertificateId: aws.String(caID)}); err != nil {
			return fmt.Errorf("DeleteCACertificate failed: %w", err)
		}
		caID = "" // prevent the deferred delete from erroring
		return nil
	}))

	return results
}
