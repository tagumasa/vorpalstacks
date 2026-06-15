package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func (r *TestRunner) runIoTCertTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	var certID string

	results = append(results, r.RunTest("iot", "Cert_CreateKeysAndCertificate", func() error {
		out, err := tc.client.CreateKeysAndCertificate(tc.ctx, &iot.CreateKeysAndCertificateInput{
			SetAsActive: true,
		})
		if err != nil {
			return fmt.Errorf("CreateKeysAndCertificate failed: %w", err)
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("expected non-empty certificateId")
		}
		certID = *out.CertificateId
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

	results = append(results, r.RunTest("iot", "Cert_DescribeCertificate", func() error {
		if certID == "" {
			return fmt.Errorf("no certificateId from previous test")
		}
		out, err := tc.client.DescribeCertificate(tc.ctx, &iot.DescribeCertificateInput{
			CertificateId: aws.String(certID),
		})
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

	results = append(results, r.RunTest("iot", "Cert_UpdateCertificate", func() error {
		if certID == "" {
			return fmt.Errorf("no certificateId from previous test")
		}
		_, err := tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String(certID),
			NewStatus:     types.CertificateStatusInactive,
		})
		if err != nil {
			return fmt.Errorf("UpdateCertificate to INACTIVE failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_ListCertificates", func() error {
		out, err := tc.client.ListCertificates(tc.ctx, &iot.ListCertificatesInput{})
		if err != nil {
			return fmt.Errorf("ListCertificates failed: %w", err)
		}
		if out.Certificates == nil {
			return fmt.Errorf("expected non-nil certificates list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Cert_DeleteCertificate", func() error {
		if certID == "" {
			return fmt.Errorf("no certificateId from previous test")
		}
		_, err := tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
			CertificateId: aws.String(certID),
		})
		if err != nil {
			return fmt.Errorf("DeleteCertificate failed: %w", err)
		}
		return nil
	}))

	return results
}
