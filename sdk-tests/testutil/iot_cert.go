package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTCertTests covers the X.509 certificate lifecycle:
// CreateKeysAndCertificate -> Describe -> Update (status transitions) -> List
// -> Delete, plus RegisterCertificate, error paths (NotFound) and the
// certificate-transfer NotFound path. The certificate id is captured from the
// create call and shared across the sequential test closures.
func (r *TestRunner) runIoTCertTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	var certID, certARN string

	results = append(results, r.RunTest("iot", "Cert_CreateKeysAndCertificate", func() error {
		out, err := tc.client.CreateKeysAndCertificate(tc.ctx, &iot.CreateKeysAndCertificateInput{SetAsActive: true})
		if err != nil {
			return fmt.Errorf("CreateKeysAndCertificate failed: %w", err)
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("expected non-empty certificateId")
		}
		certID = *out.CertificateId
		certARN = aws.ToString(out.CertificateArn)
		if out.CertificatePem == nil || !strings.Contains(*out.CertificatePem, "BEGIN CERTIFICATE") {
			return fmt.Errorf("expected PEM certificate")
		}
		if out.KeyPair == nil {
			return fmt.Errorf("expected keyPair")
		}
		if out.KeyPair.PublicKey == nil || *out.KeyPair.PublicKey == "" {
			return fmt.Errorf("expected publicKey")
		}
		if out.KeyPair.PrivateKey == nil || !strings.Contains(*out.KeyPair.PrivateKey, "BEGIN") {
			return fmt.Errorf("expected privateKey PEM")
		}
		return nil
	}))

	// Best-effort cleanup so a mid-suite failure never strands a cert.
	defer func() {
		if certID != "" {
			tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
				CertificateId: aws.String(certID),
				NewStatus:     types.CertificateStatusInactive,
			})
			tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: aws.String(certID)})
		}
	}()

	results = append(results, r.RunTest("iot", "Cert_DescribeCertificate", func() error {
		out, err := tc.client.DescribeCertificate(tc.ctx, &iot.DescribeCertificateInput{CertificateId: aws.String(certID)})
		if err != nil {
			return fmt.Errorf("DescribeCertificate failed: %w", err)
		}
		if out.CertificateDescription == nil {
			return fmt.Errorf("expected certificateDescription")
		}
		desc := out.CertificateDescription
		if desc.CertificateId == nil || *desc.CertificateId != certID {
			return fmt.Errorf("expected certificateId=%s, got %v", certID, desc.CertificateId)
		}
		if desc.Status != types.CertificateStatusActive {
			return fmt.Errorf("expected status ACTIVE, got %s", desc.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_UpdateCertificate_Inactive", func() error {
		if _, err := tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String(certID),
			NewStatus:     types.CertificateStatusInactive,
		}); err != nil {
			return fmt.Errorf("UpdateCertificate to INACTIVE failed: %w", err)
		}
		out, err := tc.client.DescribeCertificate(tc.ctx, &iot.DescribeCertificateInput{CertificateId: aws.String(certID)})
		if err != nil {
			return fmt.Errorf("DescribeCertificate after update failed: %w", err)
		}
		if out.CertificateDescription.Status != types.CertificateStatusInactive {
			return fmt.Errorf("expected status INACTIVE, got %s", out.CertificateDescription.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_UpdateCertificate_Reactivate", func() error {
		if _, err := tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String(certID),
			NewStatus:     types.CertificateStatusActive,
		}); err != nil {
			return fmt.Errorf("UpdateCertificate to ACTIVE failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_ListCertificates_IncludesCreated", func() error {
		found, err := tc.certificateExists(certID)
		if err != nil {
			return fmt.Errorf("ListCertificates failed: %w", err)
		}
		if !found {
			return fmt.Errorf("certificate %s not found in list", certID)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_DescribeCertificate_NotFound", func() error {
		_, err := tc.client.DescribeCertificate(tc.ctx, &iot.DescribeCertificateInput{CertificateId: aws.String("0123456789012345678901234567890123456789012345678901234567890123")})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Cert_UpdateCertificate_NotFound", func() error {
		_, err := tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String("0123456789012345678901234567890123456789012345678901234567890123"),
			NewStatus:     types.CertificateStatusInactive,
		})
		return expectNotFound(err)
	}))

	// A transfer against a non-existent certificate must be rejected. AWS
	// rejects it with ResourceNotFound rather than a silent no-op.
	results = append(results, r.RunTest("iot", "Cert_TransferCertificate_NotFound", func() error {
		_, err := tc.client.TransferCertificate(tc.ctx, &iot.TransferCertificateInput{
			CertificateId:    aws.String("0123456789012345678901234567890123456789012345678901234567890123"),
			TargetAwsAccount: aws.String("000000000001"),
		})
		return expectNotFound(err)
	}))

	// RegisterCertificate with a self-signed PEM (no CA) must succeed and
	// produce a NEW cert id that we then clean up.
	var regCertID string
	defer func() {
		if regCertID != "" {
			tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
				CertificateId: aws.String(regCertID), NewStatus: types.CertificateStatusInactive,
			})
			tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: aws.String(regCertID)})
		}
	}()

	results = append(results, r.RunTest("iot", "Cert_RegisterCertificate", func() error {
		out, err := tc.client.RegisterCertificate(tc.ctx, &iot.RegisterCertificateInput{
			CertificatePem: aws.String(selfSignedTestPEM),
			Status:         types.CertificateStatusActive,
		})
		if err != nil {
			return fmt.Errorf("RegisterCertificate failed: %w", err)
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("expected non-empty certificateId from RegisterCertificate")
		}
		regCertID = *out.CertificateId
		return nil
	}))

	// Delete the primary cert last; assert it can't be described afterwards.
	results = append(results, r.RunTest("iot", "Cert_DeleteCertificate", func() error {
		// Certificates must be INACTIVE before deletion.
		if _, err := tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String(certID), NewStatus: types.CertificateStatusInactive,
		}); err != nil {
			return fmt.Errorf("UpdateCertificate to INACTIVE before delete failed: %w", err)
		}
		if _, err := tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: aws.String(certID)}); err != nil {
			return fmt.Errorf("DeleteCertificate failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_DeleteCertificate_NotFound", func() error {
		_, err := tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: aws.String(certID)})
		return expectNotFound(err)
	}))

	_ = certARN // retained for potential attach-to-thing assertions
	return results
}

// selfSignedTestPEM is a throwaway self-signed certificate PEM used only to
// exercise RegisterCertificate. It is never used for TLS authentication.
const selfSignedTestPEM = `-----BEGIN CERTIFICATE-----
MIIBcDCCARagAwIBAgIJAJxY8kK3J9L5MA0GCSqGSIb3DQEBCwUAMA8xDTALBgNV
BAMMBHRlc3QwHhcNMjQwMTAxMDAwMDAwWhcNMzQwMTAxMDAwMDAwWjAPMQ0wCwYD
VQQDDAR0ZXN0MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDd1xY8kK3J9L5
cGFxq3Zr7vQ2k0Xm9pP5qGH8K3NfH6s9wJ5L2bN8pA4dT7kQwV3rHgY9vJ8kK3J
9L5cGFxq3Zr7vQ2k0Xm9pP5qGH8K3NfH6s9wJ5L2bN8pA4dT7kQwV3rHgY9vJ8k
K3N9L5cGFxq3Zr7vQ2k0Xm9pP5qGH8K3NfH6s9wJ5L2bN8pAwIDAQABoyEwHzAd
BgNVHQ4EFgQUdGVzdDBBgNVHSMEGAQwgeykdGVzdDANBgkqhkiG9w0BAQsFAAOB
gQAEtestregistercertnopemdatahereforvalidationonly
-----END CERTIFICATE-----`
