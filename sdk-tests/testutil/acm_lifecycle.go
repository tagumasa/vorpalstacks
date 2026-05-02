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
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "DeleteCertificate_VerifyGone", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		arn := aws.ToString(importResp.CertificateArn)

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
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodEmail,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.ResendValidationEmail(tc.ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   resp.CertificateArn,
			Domain:           aws.String(domain),
			ValidationDomain: aws.String(domain),
		})
		return err
	}))

	results = append(results, r.RunTest("acm", "ResendValidationEmail_IssuedStatus", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.ResendValidationEmail(tc.ctx, &acm.ResendValidationEmailInput{
			CertificateArn:   importResp.CertificateArn,
			Domain:           aws.String("example.com"),
			ValidationDomain: aws.String("example.com"),
		})
		return AssertErrorContains(err, "InvalidStateException")
	}))

	results = append(results, r.RunTest("acm", "UpdateCertificateOptions_VerifyInDescribe", func() error {
		domain := acmUniqueDomain("updopt")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.UpdateCertificateOptions(tc.ctx, &acm.UpdateCertificateOptionsInput{
			CertificateArn: resp.CertificateArn,
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceDisabled,
			},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: resp.CertificateArn})
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
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			Options: &types.CertificateOptions{
				CertificateTransparencyLoggingPreference: types.CertificateTransparencyLoggingPreferenceEnabled,
			},
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RenewCertificate_ImportedCert_Error", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.RenewCertificate(tc.ctx, &acm.RenewCertificateInput{
			CertificateArn: importResp.CertificateArn,
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("acm", "RenewCertificate_NonExistent", func() error {
		_, err := tc.client.RenewCertificate(tc.ctx, &acm.RenewCertificateInput{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_ImportedCert", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.Status != types.CertificateStatusRevoked {
			return fmt.Errorf("expected REVOKED, got %s", desc.Certificate.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_VerifyRevokedAt", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonSuperseded,
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.RevokedAt == nil {
			return fmt.Errorf("RevokedAt is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_VerifyRevocationReason", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{CertificateArn: importResp.CertificateArn})
		if err != nil {
			return err
		}
		if desc.Certificate.RevocationReason != types.RevocationReasonKeyCompromise {
			return fmt.Errorf("expected KEY_COMPROMISE, got %s", desc.Certificate.RevocationReason)
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_AlreadyRevoked", func() error {
		importResp, err := tc.importDefaultCert()
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(importResp.CertificateArn))

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("first revoke failed: %v", err)
		}
		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   importResp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err == nil {
			return fmt.Errorf("expected error for already revoked cert")
		}
		return nil
	}))

	results = append(results, r.RunTest("acm", "RevokeCertificate_PendingValidation", func() error {
		domain := acmUniqueDomain("revoke-pv")
		resp, err := tc.client.RequestCertificate(tc.ctx, &acm.RequestCertificateInput{
			DomainName:       aws.String(domain),
			ValidationMethod: types.ValidationMethodDns,
		})
		if err != nil {
			return err
		}
		defer tc.deleteCert(aws.ToString(resp.CertificateArn))

		_, err = tc.client.RevokeCertificate(tc.ctx, &acm.RevokeCertificateInput{
			CertificateArn:   resp.CertificateArn,
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		if err != nil {
			return fmt.Errorf("RevokeCertificate failed: %v", err)
		}
		desc, err := tc.client.DescribeCertificate(tc.ctx, &acm.DescribeCertificateInput{
			CertificateArn: resp.CertificateArn,
		})
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
			CertificateArn:   aws.String("arn:aws:acm:us-east-1:123456789012:certificate/nonexistent"),
			RevocationReason: types.RevocationReasonKeyCompromise,
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	return results
}
