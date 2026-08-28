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
	caPem := iotCACertPEM
	var caID string

	defer func() {
		if caID != "" {
			tc.client.DeleteCACertificate(tc.ctx, &iot.DeleteCACertificateInput{CertificateId: aws.String(caID)})
		}
	}()

	results = append(results, r.RunTest("iot", "CACert_RegisterCACertificate", func() error {
		out, err := tc.client.RegisterCACertificate(tc.ctx, &iot.RegisterCACertificateInput{
			CaCertificate:           aws.String(caPem),
			VerificationCertificate: aws.String(iotCAVerificationCertPEM),
			SetAsActive:             true,
			AllowAutoRegistration:   true,
			CertificateMode:         iottypes.CertificateModeDefault,
			RegistrationConfig: &iottypes.RegistrationConfig{
				TemplateName: aws.String("test-provisioning-template"),
				RoleArn:      aws.String("arn:aws:iam::123456789012:role/test-jitp"),
			},
			Tags: []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("ca-sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("RegisterCACertificate failed: %w", err)
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("expected non-empty certificateId")
		}
		caID = *out.CertificateId
		// setAsActive/allowAutoRegistration/certificateMode take effect.
		desc, err := tc.client.DescribeCACertificate(tc.ctx, &iot.DescribeCACertificateInput{CertificateId: aws.String(caID)})
		if err != nil {
			return fmt.Errorf("DescribeCACertificate after register failed: %w", err)
		}
		if desc.CertificateDescription == nil {
			return fmt.Errorf("expected certificateDescription after register")
		}
		if desc.CertificateDescription.Status != iottypes.CACertificateStatusActive {
			return fmt.Errorf("expected status ACTIVE from setAsActive=true, got %v", desc.CertificateDescription.Status)
		}
		if desc.CertificateDescription.AutoRegistrationStatus != iottypes.AutoRegistrationStatusEnable {
			return fmt.Errorf("expected autoRegistrationStatus=ENABLE, got %v", desc.CertificateDescription.AutoRegistrationStatus)
		}
		if desc.CertificateDescription.CertificateMode != iottypes.CertificateModeDefault {
			return fmt.Errorf("expected certificateMode=DEFAULT, got %v", desc.CertificateDescription.CertificateMode)
		}
		// Registration-time tags are visible through ListTagsForResource.
		tags, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: out.CertificateArn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, t := range tags.Tags {
			if aws.ToString(t.Key) == "purpose" && aws.ToString(t.Value) == "ca-sdk-test" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("registration tag purpose=ca-sdk-test not found on %s", aws.ToString(out.CertificateArn))
		}
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
		if out.CertificateDescription.CertificateArn == nil || *out.CertificateDescription.CertificateArn == "" {
			return fmt.Errorf("expected non-empty certificateArn in description")
		}
		if out.CertificateDescription.CertificatePem == nil || *out.CertificateDescription.CertificatePem == "" {
			return fmt.Errorf("expected non-empty certificatePem in description")
		}
		if out.RegistrationConfig == nil || aws.ToString(out.RegistrationConfig.RoleArn) != "arn:aws:iam::123456789012:role/test-jitp" {
			return fmt.Errorf("expected registrationConfig with the registered roleArn, got %v", out.RegistrationConfig)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "CACert_UpdateCACertificate", func() error {
		if _, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId:             aws.String(caID),
			NewStatus:                 iottypes.CACertificateStatusInactive,
			NewAutoRegistrationStatus: iottypes.AutoRegistrationStatusDisable,
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
		if out.CertificateDescription.AutoRegistrationStatus != iottypes.AutoRegistrationStatusDisable {
			return fmt.Errorf("expected autoRegistrationStatus=DISABLE, got %s", out.CertificateDescription.AutoRegistrationStatus)
		}
		// removeAutoRegistration clears auto registration after re-enabling it.
		if _, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId:             aws.String(caID),
			NewAutoRegistrationStatus: iottypes.AutoRegistrationStatusEnable,
		}); err != nil {
			return fmt.Errorf("re-enabling auto registration failed: %w", err)
		}
		if _, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId:          aws.String(caID),
			RemoveAutoRegistration: true,
		}); err != nil {
			return fmt.Errorf("UpdateCACertificate with removeAutoRegistration failed: %w", err)
		}
		out, err = tc.client.DescribeCACertificate(tc.ctx, &iot.DescribeCACertificateInput{CertificateId: aws.String(caID)})
		if err != nil {
			return fmt.Errorf("DescribeCACertificate after removeAutoRegistration failed: %w", err)
		}
		if out.CertificateDescription.AutoRegistrationStatus != iottypes.AutoRegistrationStatusDisable {
			return fmt.Errorf("expected autoRegistrationStatus=DISABLE after removeAutoRegistration, got %s", out.CertificateDescription.AutoRegistrationStatus)
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
				if c.CreationDate == nil {
					return fmt.Errorf("expected non-nil creationDate in list entry")
				}
				return nil
			}
		}
		return fmt.Errorf("CA certificate %s not found in list of %d", caID, len(out.Certificates))
	}))

	results = append(results, r.RunTest("iot", "CACert_UpdateCACertificate_InvalidStatusRejected", func() error {
		if _, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId: aws.String(caID),
			NewStatus:     iottypes.CACertificateStatus("GARBAGE"),
		}); err == nil {
			return fmt.Errorf("expected invalid newStatus to be rejected")
		} else if ve := expectValidationError(err); ve != nil {
			return ve
		}
		_, err := tc.client.UpdateCACertificate(tc.ctx, &iot.UpdateCACertificateInput{
			CertificateId:             aws.String(caID),
			NewAutoRegistrationStatus: iottypes.AutoRegistrationStatus("GARBAGE"),
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "CACert_RegisterCACertificate_ModePairingRejected", func() error {
		sniPem := iotSniCACertPEM
		// DEFAULT (or omitted mode) requires a verification certificate.
		if _, err := tc.client.RegisterCACertificate(tc.ctx, &iot.RegisterCACertificateInput{
			CaCertificate: aws.String(sniPem),
		}); err == nil {
			return fmt.Errorf("expected omitted verificationCertificate with DEFAULT mode to be rejected")
		} else if ve := expectValidationError(err); ve != nil {
			return ve
		}
		// SNI_ONLY forbids a verification certificate.
		if _, err := tc.client.RegisterCACertificate(tc.ctx, &iot.RegisterCACertificateInput{
			CaCertificate:           aws.String(sniPem),
			VerificationCertificate: aws.String(iotCAVerificationCertPEM),
			CertificateMode:         iottypes.CertificateModeSniOnly,
		}); err == nil {
			return fmt.Errorf("expected SNI_ONLY with verificationCertificate to be rejected")
		} else if ve := expectValidationError(err); ve != nil {
			return ve
		}
		// SNI_ONLY without a verification certificate is the valid pairing.
		out, err := tc.client.RegisterCACertificate(tc.ctx, &iot.RegisterCACertificateInput{
			CaCertificate:   aws.String(sniPem),
			CertificateMode: iottypes.CertificateModeSniOnly,
		})
		if err != nil {
			return fmt.Errorf("SNI_ONLY registration failed: %w", err)
		}
		// Registered INACTIVE by default, so it can be deleted right away.
		_, _ = tc.client.DeleteCACertificate(tc.ctx, &iot.DeleteCACertificateInput{CertificateId: out.CertificateId})
		return nil
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
