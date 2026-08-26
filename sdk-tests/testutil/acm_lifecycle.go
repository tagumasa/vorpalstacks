package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func (r *TestRunner) runACMLifecycleTests(tc *acmTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("acm", "DeleteCertificate_NonExistent", func() error {
		_, err := tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(tc.nonExistentArn()),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "DeleteCertificate_VerifyGone", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}

		_, err = tc.client.DeleteCertificate(tc.ctx, &acm.DeleteCertificateInput{
			CertificateArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("delete failed: %v", err)
		}

		_, err = tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
			CertificateArn: aws.String(arn),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("expected ResourceNotFoundException after delete: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "ResendValidationEmail", func() error {
		domain := acmUniqueDomain("resend")
		arn, err := tc.requestEmailCert(domain)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		// Our edge implementation immediately issues certificates (no real
		// validation phase), so the cert is in ISSUED status by the time
		// ResendValidationEmail is called. AWS only allows resend on
		// PENDING_VALIDATION certs, so this correctly returns
		// InvalidStateException.
		_, err = tc.client.ResendValidationEmail(tc.ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   aws.String(arn),
			Domain:           aws.String(domain),
			ValidationDomain: aws.String(domain),
		})
		if err := AssertErrorContains(err, "InvalidStateException"); err != nil {
			return err
		}

		// Imported certs are born ISSUED without ever entering
		// PENDING_VALIDATION, so the same InvalidStateException applies.
		importArn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(importArn)
		_, err = tc.client.ResendValidationEmail(tc.ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   aws.String(importArn),
			Domain:           aws.String("example.com"),
			ValidationDomain: aws.String("example.com"),
		})
		return AssertErrorContains(err, "InvalidStateException")
	}))

	results = append(results, r.RunTest("acm", "UpdateCertificateOptions_VerifyInDescribe", func() error {
		domain := acmUniqueDomain("updopt")
		arn, err := tc.requestDNSCert(domain)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.UpdateCertificateOptions(tc.ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: aws.String(arn),
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		if desc.Certificate.Options == nil {
			return fmt.Errorf("options is nil after update")
		}
		if desc.Certificate.Options.CertificateTransparencyLoggingPreference != types.CertificateTransparencyLoggingPreferenceDisabled {
			return fmt.Errorf("expected DISABLED, got %s", desc.Certificate.Options.CertificateTransparencyLoggingPreference)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "UpdateCertificateOptions_NonExistent", func() error {
		_, err := tc.client.UpdateCertificateOptions(tc.ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: aws.String(tc.nonExistentArn()),
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled,
			},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RenewCertificate_ImportedCert_Error", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.RenewCertificate(tc.ctx, &acm.RenewCertificateInput{
			CertificateArn: aws.String(arn),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("acm", "RenewCertificate_NonExistent", func() error {
		_, err := tc.client.RenewCertificate(tc.ctx, &acm.RenewCertificateInput{
			CertificateArn: aws.String(tc.nonExistentArn()),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_ImportedCert", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(arn),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		if desc.Certificate.Status != types.CertificateStatusRevoked {
			return fmt.Errorf("expected REVOKED, got %s", desc.Certificate.Status)
		}
		if desc.Certificate.RevokedAt == nil {
			return fmt.Errorf("RevokedAt is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_VerifyRevocationReason", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(arn),
			RevocationReason: types.RevocationReasonSuperseded,
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		if desc.Certificate.RevocationReason != types.RevocationReasonSuperseded {
			return fmt.Errorf("expected SUPERSEDED, got %s", desc.Certificate.RevocationReason)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_AlreadyRevoked", func() error {
		arn, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(arn),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("first revoke failed: %v", err)
		}
		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(arn),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err == nil {
			return fmt.Errorf("expected error for already revoked cert")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_PendingValidation", func() error {
		domain := acmUniqueDomain("revoke-pv")
		arn, err := tc.requestDNSCert(domain)
		if err != nil {
			return err
		}
		defer tc.deleteCert(arn)

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(arn),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("RevokeCertificate failed: %v", err)
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: aws.String(arn)})
		if err != nil {
			return err
		}
		if desc.Certificate.Status != types.CertificateStatusRevoked {
			return fmt.Errorf("expected REVOKED, got %s", desc.Certificate.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_NonExistent", func() error {
		_, err := tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   aws.String(tc.nonExistentArn()),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	return results
}
