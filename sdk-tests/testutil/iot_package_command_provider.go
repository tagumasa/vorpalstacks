package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTPackageCommandProviderTests covers Package/SBOM, CertificateProvider,
// and Command operations that were previously unregistered stubs.
func (r *TestRunner) runIoTPackageCommandProviderTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	reg := tc.region
	acct := tc.accountID
	pkgName := uniqueName("test-pkg")
	versionName := "1.0.0"
	cmdID := uniqueName("test-cmd")
	certProviderName := uniqueName("test-cert-provider")

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeleteCertificateProvider(tc.ctx, &iot.DeleteCertificateProviderInput{CertificateProviderName: aws.String(certProviderName)})
		_, _ = tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{CommandId: aws.String(cmdID)})
		_, _ = tc.client.DeletePackageVersion(tc.ctx, &iot.DeletePackageVersionInput{PackageName: aws.String(pkgName), VersionName: aws.String(versionName)})
		_, _ = tc.client.DeletePackage(tc.ctx, &iot.DeletePackageInput{PackageName: aws.String(pkgName)})
	}()

	// --- Package lifecycle ---
	results = append(results, r.RunTest("iot", "Package_CreatePackage", func() error {
		_, err := tc.client.CreatePackage(tc.ctx, &iot.CreatePackageInput{
			PackageName: aws.String(pkgName),
			Description: aws.String("test package"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_GetPackage", func() error {
		out, err := tc.client.GetPackage(tc.ctx, &iot.GetPackageInput{
			PackageName: aws.String(pkgName),
		})
		if err != nil {
			return err
		}
		if out.PackageName == nil || *out.PackageName != pkgName {
			return fmt.Errorf("packageName mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Package_GetPackage_NotFound", func() error {
		_, err := tc.client.GetPackage(tc.ctx, &iot.GetPackageInput{
			PackageName: aws.String(uniqueName("nope-pkg")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Package_UpdatePackage", func() error {
		_, err := tc.client.UpdatePackage(tc.ctx, &iot.UpdatePackageInput{
			PackageName:        aws.String(pkgName),
			Description:        aws.String("updated"),
			DefaultVersionName: aws.String(versionName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_ListPackages_IncludesCreated", func() error {
		pkgs, err := paginate(func(next *string) ([]iottypes.PackageSummary, *string, error) {
			out, err := tc.client.ListPackages(tc.ctx, &iot.ListPackagesInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.PackageSummaries, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, p := range pkgs {
			if p.PackageName != nil && *p.PackageName == pkgName {
				return nil
			}
		}
		return fmt.Errorf("created package not found in ListPackages")
	}))

	// --- PackageVersion lifecycle ---
	results = append(results, r.RunTest("iot", "Package_CreatePackageVersion", func() error {
		_, err := tc.client.CreatePackageVersion(tc.ctx, &iot.CreatePackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
			Description: aws.String("v1"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_GetPackageVersion", func() error {
		out, err := tc.client.GetPackageVersion(tc.ctx, &iot.GetPackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		if err != nil {
			return err
		}
		if out.Status != iottypes.PackageVersionStatusDraft {
			return fmt.Errorf("expected DRAFT status, got %s", out.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Package_UpdatePackageVersion_Publish", func() error {
		_, err := tc.client.UpdatePackageVersion(tc.ctx, &iot.UpdatePackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
			Action:      iottypes.PackageVersionActionPublish,
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_ListPackageVersions", func() error {
		out, err := tc.client.ListPackageVersions(tc.ctx, &iot.ListPackageVersionsInput{
			PackageName: aws.String(pkgName),
		})
		if err != nil {
			return err
		}
		if len(out.PackageVersionSummaries) == 0 {
			return fmt.Errorf("no versions returned")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Package_AssociateSbom", func() error {
		_, err := tc.client.AssociateSbomWithPackageVersion(tc.ctx, &iot.AssociateSbomWithPackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
			Sbom: &iottypes.Sbom{
				S3Location: &iottypes.S3Location{Bucket: aws.String("test-bucket"), Key: aws.String("sbom.json")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_DeletePackageVersion", func() error {
		_, err := tc.client.DeletePackageVersion(tc.ctx, &iot.DeletePackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Package_DeletePackage", func() error {
		_, err := tc.client.DeletePackage(tc.ctx, &iot.DeletePackageInput{
			PackageName: aws.String(pkgName),
		})
		return err
	}))

	// --- Command lifecycle ---
	results = append(results, r.RunTest("iot", "Command_CreateCommand", func() error {
		_, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
			CommandId:   aws.String(cmdID),
			DisplayName: aws.String("test command"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Command_GetCommand", func() error {
		out, err := tc.client.GetCommand(tc.ctx, &iot.GetCommandInput{
			CommandId: aws.String(cmdID),
		})
		if err != nil {
			return err
		}
		if out.CommandId == nil || *out.CommandId != cmdID {
			return fmt.Errorf("commandId mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Command_GetCommand_NotFound", func() error {
		_, err := tc.client.GetCommand(tc.ctx, &iot.GetCommandInput{
			CommandId: aws.String(uniqueName("nope-cmd")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Command_UpdateCommand", func() error {
		_, err := tc.client.UpdateCommand(tc.ctx, &iot.UpdateCommandInput{
			CommandId:   aws.String(cmdID),
			Description: aws.String("updated"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Command_ListCommands_IncludesCreated", func() error {
		cmds, err := paginate(func(next *string) ([]iottypes.CommandSummary, *string, error) {
			out, err := tc.client.ListCommands(tc.ctx, &iot.ListCommandsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.Commands, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, c := range cmds {
			if c.CommandId != nil && *c.CommandId == cmdID {
				return nil
			}
		}
		return fmt.Errorf("created command not found in ListCommands")
	}))

	results = append(results, r.RunTest("iot", "Command_DeleteCommand", func() error {
		_, err := tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{
			CommandId: aws.String(cmdID),
		})
		return err
	}))

	// --- CertificateProvider lifecycle ---
	results = append(results, r.RunTest("iot", "CertProvider_Create", func() error {
		_, err := tc.client.CreateCertificateProvider(tc.ctx, &iot.CreateCertificateProviderInput{
			CertificateProviderName: aws.String(certProviderName),
			LambdaFunctionArn:       aws.String(fmt.Sprintf("arn:aws:lambda:%s:%s:function:test-cert-signer", reg, acct)),
			AccountDefaultForOperations: []iottypes.CertificateProviderOperation{
				iottypes.CertificateProviderOperationCreateCertificateFromCsr,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "CertProvider_Describe", func() error {
		out, err := tc.client.DescribeCertificateProvider(tc.ctx, &iot.DescribeCertificateProviderInput{
			CertificateProviderName: aws.String(certProviderName),
		})
		if err != nil {
			return err
		}
		if out.CertificateProviderName == nil || *out.CertificateProviderName != certProviderName {
			return fmt.Errorf("provider name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "CertProvider_Describe_NotFound", func() error {
		_, err := tc.client.DescribeCertificateProvider(tc.ctx, &iot.DescribeCertificateProviderInput{
			CertificateProviderName: aws.String(uniqueName("nope-provider")),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "CertProvider_List_IncludesCreated", func() error {
		providers, err := paginate(func(next *string) ([]iottypes.CertificateProviderSummary, *string, error) {
			out, err := tc.client.ListCertificateProviders(tc.ctx, &iot.ListCertificateProvidersInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.CertificateProviders, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, p := range providers {
			if p.CertificateProviderName != nil && *p.CertificateProviderName == certProviderName {
				return nil
			}
		}
		return fmt.Errorf("created provider not found in ListCertificateProviders")
	}))

	results = append(results, r.RunTest("iot", "CertProvider_Delete", func() error {
		_, err := tc.client.DeleteCertificateProvider(tc.ctx, &iot.DeleteCertificateProviderInput{
			CertificateProviderName: aws.String(certProviderName),
		})
		return err
	}))

	return results
}
