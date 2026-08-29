package testutil

import (
	"fmt"
	"time"

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

	results = append(results, r.RunTest("iot", "Package_UpdatePackage_UnsetDefaultVersion", func() error {
		out, err := tc.client.GetPackage(tc.ctx, &iot.GetPackageInput{PackageName: aws.String(pkgName)})
		if err != nil {
			return err
		}
		if aws.ToString(out.DefaultVersionName) != versionName {
			return fmt.Errorf("expected defaultVersionName=%s before unset, got %v", versionName, out.DefaultVersionName)
		}
		if _, err := tc.client.UpdatePackage(tc.ctx, &iot.UpdatePackageInput{
			PackageName:         aws.String(pkgName),
			UnsetDefaultVersion: aws.Bool(true),
		}); err != nil {
			return fmt.Errorf("UpdatePackage unsetDefaultVersion failed: %w", err)
		}
		after, err := tc.client.GetPackage(tc.ctx, &iot.GetPackageInput{PackageName: aws.String(pkgName)})
		if err != nil {
			return err
		}
		if after.DefaultVersionName != nil {
			return fmt.Errorf("expected defaultVersionName cleared, got %v", *after.DefaultVersionName)
		}
		// Setting and unsetting the default version at once is rejected.
		_, err = tc.client.UpdatePackage(tc.ctx, &iot.UpdatePackageInput{
			PackageName:         aws.String(pkgName),
			DefaultVersionName:  aws.String(versionName),
			UnsetDefaultVersion: aws.Bool(true),
		})
		if err == nil {
			return fmt.Errorf("expected set+unset to be rejected")
		}
		return expectValidationError(err)
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

	results = append(results, r.RunTest("iot", "Command_ListCommands_DefaultDescendingOrder", func() error {
		earlyID := uniqueName("test-cmd-early")
		lateID := uniqueName("test-cmd-late")
		defer tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{CommandId: aws.String(earlyID)})
		defer tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{CommandId: aws.String(lateID)})
		if _, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
			CommandId:   aws.String(earlyID),
			DisplayName: aws.String("default order early command"),
		}); err != nil {
			return fmt.Errorf("CreateCommand early failed: %w", err)
		}
		// createdAt carries second granularity, so keep the two records
		// in distinct seconds before comparing their default order.
		time.Sleep(1100 * time.Millisecond)
		if _, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
			CommandId:   aws.String(lateID),
			DisplayName: aws.String("default order late command"),
		}); err != nil {
			return fmt.Errorf("CreateCommand late failed: %w", err)
		}
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
		earlyIdx, lateIdx := -1, -1
		for i, c := range cmds {
			switch aws.ToString(c.CommandId) {
			case earlyID:
				earlyIdx = i
			case lateID:
				lateIdx = i
			}
		}
		if earlyIdx == -1 || lateIdx == -1 {
			return fmt.Errorf("both commands expected in the default list (early=%d late=%d)", earlyIdx, lateIdx)
		}
		// The API documents that, without sortOrder, commands are listed
		// in descending order of creation time.
		if lateIdx > earlyIdx {
			return fmt.Errorf("expected the later-created command first, early=%d late=%d", earlyIdx, lateIdx)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Command_ListCommands_SortOrderAndParameterFilter", func() error {
		paramCmdID := uniqueName("test-cmd-param")
		defer tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{CommandId: aws.String(paramCmdID)})
		if _, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
			CommandId:   aws.String(paramCmdID),
			DisplayName: aws.String("parameterised command"),
			MandatoryParameters: []iottypes.CommandParameter{{
				Name: aws.String("p1"),
			}},
		}); err != nil {
			return fmt.Errorf("CreateCommand with a mandatory parameter failed: %w", err)
		}

		// The parameter filter matches commands declaring the parameter
		// and hides commands without it.
		filtered, err := tc.client.ListCommands(tc.ctx, &iot.ListCommandsInput{
			CommandParameterName: aws.String("p1"),
		})
		if err != nil {
			return fmt.Errorf("ListCommands commandParameterName failed: %w", err)
		}
		hasParam, leakedPlain := false, false
		for _, c := range filtered.Commands {
			if aws.ToString(c.CommandId) == paramCmdID {
				hasParam = true
			}
			if aws.ToString(c.CommandId) == cmdID {
				leakedPlain = true
			}
		}
		if !hasParam || leakedPlain {
			return fmt.Errorf("commandParameterName filter: hasParam=%v leakedPlain=%v", hasParam, leakedPlain)
		}

		// Descending order returns the later-created command first.
		descending, err := tc.client.ListCommands(tc.ctx, &iot.ListCommandsInput{
			SortOrder: iottypes.SortOrderDescending,
		})
		if err != nil {
			return fmt.Errorf("ListCommands sortOrder failed: %w", err)
		}
		paramIdx, plainIdx := -1, -1
		for i, c := range descending.Commands {
			if aws.ToString(c.CommandId) == paramCmdID {
				paramIdx = i
			}
			if aws.ToString(c.CommandId) == cmdID {
				plainIdx = i
			}
		}
		if paramIdx == -1 || plainIdx == -1 {
			return fmt.Errorf("both commands expected in the descending list (param=%d plain=%d)", paramIdx, plainIdx)
		}
		if paramIdx > plainIdx {
			return fmt.Errorf("expected the later-created command first, paramIdx=%d plainIdx=%d", paramIdx, plainIdx)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Command_ListCommandExecutions_FiltersAccepted", func() error {
		// No command executions exist on this platform, so every accepted
		// filter shape must return an empty page. The documented date and
		// time format is yyyy-MM-dd'T'HH:mm, without seconds or a zone.
		out, err := tc.client.ListCommandExecutions(tc.ctx, &iot.ListCommandExecutionsInput{
			Namespace:         iottypes.CommandNamespaceAWSIoT,
			SortOrder:         iottypes.SortOrderDescending,
			StartedTimeFilter: &iottypes.TimeFilter{After: aws.String("2026-01-15T10:30")},
		})
		if err != nil {
			return fmt.Errorf("ListCommandExecutions with the documented time format failed: %w", err)
		}
		if len(out.CommandExecutions) != 0 {
			return fmt.Errorf("expected no executions, got %d", len(out.CommandExecutions))
		}
		// Providing both time filters is documented to generate an error.
		_, err = tc.client.ListCommandExecutions(tc.ctx, &iot.ListCommandExecutionsInput{
			StartedTimeFilter:   &iottypes.TimeFilter{After: aws.String("2026-01-01T00:00")},
			CompletedTimeFilter: &iottypes.TimeFilter{Before: aws.String("2030-01-01T00:00")},
		})
		if err == nil {
			return fmt.Errorf("expected both time filters to be rejected")
		}
		if ve := expectValidationError(err); ve != nil {
			return ve
		}
		// Providing both the command ARN and the target ARN is documented
		// to generate an error as well.
		_, err = tc.client.ListCommandExecutions(tc.ctx, &iot.ListCommandExecutionsInput{
			CommandArn:        aws.String(tc.arn("iot", "command", uniqueName("cmd"))),
			TargetArn:         aws.String(tc.arn("iot", "thing", uniqueName("thing"))),
			StartedTimeFilter: &iottypes.TimeFilter{After: aws.String("2026-01-01T00:00")},
		})
		if err == nil {
			return fmt.Errorf("expected commandArn plus targetArn to be rejected")
		}
		if ve := expectValidationError(err); ve != nil {
			return ve
		}
		_, err = tc.client.ListCommandExecutions(tc.ctx, &iot.ListCommandExecutionsInput{
			StartedTimeFilter: &iottypes.TimeFilter{After: aws.String("not-a-timestamp")},
		})
		if err == nil {
			return fmt.Errorf("expected an invalid time filter to be rejected")
		}
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "Command_DeleteCommand", func() error {
		_, err := tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{
			CommandId: aws.String(cmdID),
		})
		return err
	}))

	// --- CertificateProvider lifecycle ---
	results = append(results, r.RunTest("iot", "CertProvider_InvalidOperationsRejected", func() error {
		badLambda := fmt.Sprintf("arn:aws:lambda:%s:%s:function:test-cert-signer", reg, acct)
		if _, err := tc.client.CreateCertificateProvider(tc.ctx, &iot.CreateCertificateProviderInput{
			CertificateProviderName: aws.String(uniqueName("bad-provider-two-ops")),
			LambdaFunctionArn:       aws.String(badLambda),
			AccountDefaultForOperations: []iottypes.CertificateProviderOperation{
				iottypes.CertificateProviderOperationCreateCertificateFromCsr,
				iottypes.CertificateProviderOperationCreateCertificateFromCsr,
			},
		}); err == nil {
			return fmt.Errorf("expected two-entry operations list to be rejected")
		} else if ve := expectValidationError(err); ve != nil {
			return ve
		}
		if _, err := tc.client.CreateCertificateProvider(tc.ctx, &iot.CreateCertificateProviderInput{
			CertificateProviderName:     aws.String(uniqueName("bad-provider-bogus-op")),
			LambdaFunctionArn:           aws.String(badLambda),
			AccountDefaultForOperations: []iottypes.CertificateProviderOperation{iottypes.CertificateProviderOperation("Bogus")},
		}); err == nil {
			return fmt.Errorf("expected invalid operations value to be rejected")
		} else if ve := expectValidationError(err); ve != nil {
			return ve
		}
		return nil
	}))

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
