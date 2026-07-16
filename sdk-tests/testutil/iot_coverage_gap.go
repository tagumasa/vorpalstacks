package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTCoverageGapTests covers 25 registered operations that previously had
// no SDK test coverage. Deprecated operations (DescribeTopicRule -> GetTopicRule,
// ListThingsForThingType -> ListThings?thingTypeName=) are excluded as they
// share the same handler code path and are already tested via their replacements.
func (r *TestRunner) runIoTCoverageGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.runIoTCertCSRGapTests(tc)...)
	results = append(results, r.runIoTPolicyV2GapTests(tc)...)
	results = append(results, r.runIoTLoggingConfigGapTests(tc)...)
	results = append(results, r.runIoTTopicRuleJobSecGapTests(tc)...)
	results = append(results, r.runIoTProvTemplateVersionGapTests(tc)...)
	results = append(results, r.runIoTPackageExtraGapTests(tc)...)
	results = append(results, r.runIoTCommandExecutionGapTests(tc)...)
	results = append(results, r.runIoTCertProviderUpdateGapTests(tc)...)

	return results
}

// --- Certificate CSR ---
func (r *TestRunner) runIoTCertCSRGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Gap_CreateCertificateFromCsr", func() error {
		csr := `-----BEGIN CERTIFICATE REQUEST-----
MIICVDCCATwCAQAwDzENMAsGA1UEAwwEdGVzdDCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAO10oYxFUymZkylNzfIAsSZwUmrNEKGuzcTrk4Fxko2S7Jw7
28NcpeRbj9W3vNE43L9CwJOnL8j29YbcaloEarZ0ZPBk65SNTglonHRmMfcD46tu
+ZF5wbvBkxiwmtCFgYuXKJApr6aIhSNHkukI7OcJNb/aZCzo7iq7Z75ZynLr5zXc
PS3zGQv1C08xHozwRdv6OCiuBjiUP1U6Tbc5nwec8oJ4hhCziQaoEHEMvGh+Dr51
/4EqKoOau5mespY/Fi8FiebF+8Y9/1tzsDCLXiHwA1rX7J2iX1JLuZvdjKJQcLao
NlWt0i9vXwsAqulrYhn3o8ZrkjVcNkUxoYgH/58CAwEAAaAAMA0GCSqGSIb3DQEB
CwUAA4IBAQDUx2DemSBE9HbVGv4cbydNelDEH7lO2duxAnkDsS78c+3bwuXNt0B8
phM2yxEt1sPPRxaeBvkKHtAplQAeU7xNkaJ1wYmnDMsYm+qfGlcfRJL4ZMkD90r0
yb4Jz1iw3AC7C4rEy3rJd70Jtsqt9DA2+nObq0HgbDkKPE6Abj9lURpqXnnoU+eO
wsEOldC5/vHN6ezAWr9mxu6GKXP9d3tdlJDP6JXC72cKztlcyczPRYHXGZawjkfF
9oBZ2D/oxmhrE0CP/KV8GvA0usn4WYB7karC4/irCd/EnRbV0eED+4LBytntfAIf
8fT2Q3JQCWnKmEmfwtyQ/HB29OypTjKZ
-----END CERTIFICATE REQUEST-----`
		out, err := tc.client.CreateCertificateFromCsr(tc.ctx, &iot.CreateCertificateFromCsrInput{
			CertificateSigningRequest: aws.String(csr),
		})
		if err != nil {
			return err
		}
		if out.CertificateId == nil || *out.CertificateId == "" {
			return fmt.Errorf("empty certificateId")
		}
		// Cleanup
		_, _ = tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
			CertificateId: out.CertificateId,
		})
		return nil
	}))

	return results
}

// --- Policy V2 (AttachPrincipalPolicy / DetachPrincipalPolicy) ---
func (r *TestRunner) runIoTPolicyV2GapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	policyName := uniqueName("gap-policy")
	certID := ""

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		if certID != "" {
			tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
				CertificateId: aws.String(certID), NewStatus: iottypes.CertificateStatusInactive,
			})
			tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
				CertificateId: aws.String(certID),
			})
		}
		tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{PolicyName: aws.String(policyName)})
	}()

	// Setup: create policy + certificate
	results = append(results, r.RunTest("iot", "Gap_PolicyV2_Setup", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iot.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`),
		})
		if err != nil {
			return err
		}
		cert, err := tc.client.CreateKeysAndCertificate(tc.ctx, &iot.CreateKeysAndCertificateInput{
			SetAsActive: true,
		})
		if err != nil {
			return err
		}
		certID = *cert.CertificateId
		return nil
	}))

	principalARN := func() string {
		return fmt.Sprintf("arn:aws:iot:us-east-1:000000000000:cert/%s", certID)
	}

	results = append(results, r.RunTest("iot", "Gap_AttachPrincipalPolicy", func() error {
		_, err := tc.client.AttachPrincipalPolicy(tc.ctx, &iot.AttachPrincipalPolicyInput{
			PolicyName: aws.String(policyName),
			Principal:  aws.String(principalARN()),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_DetachPrincipalPolicy", func() error {
		_, err := tc.client.DetachPrincipalPolicy(tc.ctx, &iot.DetachPrincipalPolicyInput{
			PolicyName: aws.String(policyName),
			Principal:  aws.String(principalARN()),
		})
		return err
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_PolicyV2_Teardown", func() error {
		_, err := tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: aws.String(certID),
			NewStatus:     iottypes.CertificateStatusInactive,
		})
		if err != nil {
			return err
		}
		_, err = tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
			CertificateId: aws.String(certID),
		})
		return err
	}))

	return results
}

// --- Logging / Config ---
func (r *TestRunner) runIoTLoggingConfigGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Gap_GetLoggingOptions", func() error {
		_, err := tc.client.GetLoggingOptions(tc.ctx, &iot.GetLoggingOptionsInput{})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_SetLoggingOptions", func() error {
		_, err := tc.client.SetLoggingOptions(tc.ctx, &iot.SetLoggingOptionsInput{
			LoggingOptionsPayload: &iottypes.LoggingOptionsPayload{
				RoleArn:  aws.String("arn:aws:iam::000000000000:role/test"),
				LogLevel: iottypes.LogLevelInfo,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_SetV2LoggingLevel", func() error {
		_, err := tc.client.SetV2LoggingLevel(tc.ctx, &iot.SetV2LoggingLevelInput{
			LogLevel: iottypes.LogLevelInfo,
			LogTarget: &iottypes.LogTarget{
				TargetType: iottypes.LogTargetTypeDefault,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_DeleteV2LoggingLevel", func() error {
		_, err := tc.client.DeleteV2LoggingLevel(tc.ctx, &iot.DeleteV2LoggingLevelInput{
			TargetType: iottypes.LogTargetTypeDefault,
			TargetName: aws.String(""),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_UpdateIndexingConfiguration", func() error {
		_, err := tc.client.UpdateIndexingConfiguration(tc.ctx, &iot.UpdateIndexingConfigurationInput{
			ThingIndexingConfiguration: &iottypes.ThingIndexingConfiguration{
				ThingIndexingMode: iottypes.ThingIndexingModeOff,
			},
		})
		return err
	}))

	return results
}

// --- TopicRule Describe + Job Update + SecurityProfile Update/Validate ---
func (r *TestRunner) runIoTTopicRuleJobSecGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	ruleName := uniqueName("gap-rule")
	jobID := uniqueName("gap-job")
	secProfileName := uniqueName("gap-sec-profile")
	thingType := uniqueName("gap-thing-type")

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{SecurityProfileName: aws.String(secProfileName)})
		_, _ = tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{RuleName: aws.String(ruleName)})
		_, _ = tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{JobId: aws.String(jobID)})
		_, _ = tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{ThingTypeName: aws.String(thingType)})
	}()

	// Setup: create thing type
	results = append(results, r.RunTest("iot", "Gap_SecProfile_Setup", func() error {
		_, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
			ThingTypeName: aws.String(thingType),
		})
		return err
	}))

	// SecurityProfile: Create -> Update -> Validate
	results = append(results, r.RunTest("iot", "Gap_UpdateSecurityProfile", func() error {
		// Create first
		_, err := tc.client.CreateSecurityProfile(tc.ctx, &iot.CreateSecurityProfileInput{
			SecurityProfileName: aws.String(secProfileName),
			Behaviors: []iottypes.Behavior{{
				Name: aws.String("auth-failures"),
				Criteria: &iottypes.BehaviorCriteria{
					ComparisonOperator: iottypes.ComparisonOperatorLessThan,
					Value: &iottypes.MetricValue{
						Count: aws.Int64(5),
					},
				},
			}},
		})
		if err != nil {
			return err
		}
		// Now update
		_, err = tc.client.UpdateSecurityProfile(tc.ctx, &iot.UpdateSecurityProfileInput{
			SecurityProfileName:        aws.String(secProfileName),
			SecurityProfileDescription: aws.String("updated"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_ValidateSecurityProfileBehaviors", func() error {
		out, err := tc.client.ValidateSecurityProfileBehaviors(tc.ctx, &iot.ValidateSecurityProfileBehaviorsInput{
			Behaviors: []iottypes.Behavior{{
				Name: aws.String("valid-behavior"),
				Criteria: &iottypes.BehaviorCriteria{
					ComparisonOperator: iottypes.ComparisonOperatorLessThan,
					Value: &iottypes.MetricValue{
						Count: aws.Int64(5),
					},
				},
			}},
		})
		if err != nil {
			return err
		}
		if !out.Valid {
			return fmt.Errorf("expected valid=true, got false")
		}
		return nil
	}))

	// TopicRule: Create -> Describe (via GetTopicRule which is the non-deprecated alias)
	results = append(results, r.RunTest("iot", "Gap_DescribeTopicRule_GetTopicRule", func() error {
		_, err := tc.client.CreateTopicRule(tc.ctx, &iot.CreateTopicRuleInput{
			RuleName: aws.String(ruleName),
			TopicRulePayload: &iottypes.TopicRulePayload{
				Sql:          aws.String("SELECT * FROM 'test/topic'"),
				Actions:      []iottypes.Action{{}},
				RuleDisabled: aws.Bool(false),
			},
		})
		if err != nil {
			return err
		}
		// DescribeTopicRule is deprecated in SDK; GetTopicRule exercises the same handler
		out, err := tc.client.GetTopicRule(tc.ctx, &iot.GetTopicRuleInput{
			RuleName: aws.String(ruleName),
		})
		if err != nil {
			return err
		}
		if out.Rule == nil || out.Rule.RuleName == nil || *out.Rule.RuleName != ruleName {
			return fmt.Errorf("rule name mismatch in GetTopicRule")
		}
		return nil
	}))

	// Job: Create -> Update
	results = append(results, r.RunTest("iot", "Gap_UpdateJob", func() error {
		_, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:    aws.String(jobID),
			Targets:  []string{"arn:aws:iot:us-east-1:000000000000:thing/dummy"},
			Document: aws.String(`{"operation":"test"}`),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.UpdateJob(tc.ctx, &iot.UpdateJobInput{
			JobId:       aws.String(jobID),
			Description: aws.String("updated description"),
		})
		return err
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_SecProfile_Teardown", func() error {
		_, _ = tc.client.DeleteSecurityProfile(tc.ctx, &iot.DeleteSecurityProfileInput{
			SecurityProfileName: aws.String(secProfileName),
		})
		_, _ = tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{
			RuleName: aws.String(ruleName),
		})
		_, _ = tc.client.DeleteJob(tc.ctx, &iot.DeleteJobInput{
			JobId: aws.String(jobID),
		})
		_, _ = tc.client.DeleteThingType(tc.ctx, &iot.DeleteThingTypeInput{
			ThingTypeName: aws.String(thingType),
		})
		return nil
	}))

	return results
}

// --- ProvisioningTemplateVersion lifecycle ---
func (r *TestRunner) runIoTProvTemplateVersionGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	templateName := uniqueName("gap-prov-tmpl")
	templateBody := `{"Parameters":{"SerialNumber":{"Type":"String"},"DeviceLocation":{"Type":"String"}},"Resources":{"certificate":{"Properties":{"CertificateId":{"Ref":"SerialNumber"},"CertificatePem":{"Ref":"SerialNumber"}},"Type":"AWS::IoT::Certificate"},"thing":{"Properties":{"ThingName":{"Ref":"SerialNumber"},"ThingGroups":["vstest-group"]},"Type":"AWS::IoT::Thing"}}}`

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeleteProvisioningTemplate(tc.ctx, &iot.DeleteProvisioningTemplateInput{TemplateName: aws.String(templateName)})
	}()

	// Setup: create template
	results = append(results, r.RunTest("iot", "Gap_ProvTemplateVersion_Setup", func() error {
		_, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
			TemplateName:        aws.String(templateName),
			TemplateBody:        aws.String(templateBody),
			ProvisioningRoleArn: aws.String("arn:aws:iam::000000000000:role/test"),
		})
		return err
	}))

	var createdVersionID int32

	results = append(results, r.RunTest("iot", "Gap_CreateProvisioningTemplateVersion", func() error {
		out, err := tc.client.CreateProvisioningTemplateVersion(tc.ctx, &iot.CreateProvisioningTemplateVersionInput{
			TemplateName: aws.String(templateName),
			TemplateBody: aws.String(templateBody),
		})
		if err != nil {
			return err
		}
		if out.VersionId == nil {
			return fmt.Errorf("nil versionId")
		}
		createdVersionID = *out.VersionId
		return nil
	}))

	results = append(results, r.RunTest("iot", "Gap_ListProvisioningTemplateVersions", func() error {
		out, err := tc.client.ListProvisioningTemplateVersions(tc.ctx, &iot.ListProvisioningTemplateVersionsInput{
			TemplateName: aws.String(templateName),
		})
		if err != nil {
			return err
		}
		if len(out.Versions) < 1 {
			return fmt.Errorf("expected at least 1 version")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Gap_DescribeProvisioningTemplateVersion", func() error {
		out, err := tc.client.DescribeProvisioningTemplateVersion(tc.ctx, &iot.DescribeProvisioningTemplateVersionInput{
			TemplateName: aws.String(templateName),
			VersionId:    aws.Int32(createdVersionID),
		})
		if err != nil {
			return err
		}
		if out.VersionId == nil || *out.VersionId != createdVersionID {
			return fmt.Errorf("version ID mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Gap_DeleteProvisioningTemplateVersion", func() error {
		_, err := tc.client.DeleteProvisioningTemplateVersion(tc.ctx, &iot.DeleteProvisioningTemplateVersionInput{
			TemplateName: aws.String(templateName),
			VersionId:    aws.Int32(createdVersionID),
		})
		return err
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_ProvTemplateVersion_Teardown", func() error {
		_, _ = tc.client.DeleteProvisioningTemplate(tc.ctx, &iot.DeleteProvisioningTemplateInput{
			TemplateName: aws.String(templateName),
		})
		return nil
	}))

	return results
}

// --- Package extras: GetPackageConfiguration, UpdatePackageConfiguration ---
// SBOM extras: DisassociateSbomFromPackageVersion, ListSbomValidationResults
func (r *TestRunner) runIoTPackageExtraGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	pkgName := uniqueName("gap-pkg")
	versionName := "1.0.0"

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeletePackageVersion(tc.ctx, &iot.DeletePackageVersionInput{PackageName: aws.String(pkgName), VersionName: aws.String(versionName)})
		_, _ = tc.client.DeletePackage(tc.ctx, &iot.DeletePackageInput{PackageName: aws.String(pkgName)})
	}()

	results = append(results, r.RunTest("iot", "Gap_Package_Setup", func() error {
		_, err := tc.client.CreatePackage(tc.ctx, &iot.CreatePackageInput{
			PackageName: aws.String(pkgName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.CreatePackageVersion(tc.ctx, &iot.CreatePackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_GetPackageConfiguration", func() error {
		_, err := tc.client.GetPackageConfiguration(tc.ctx, &iot.GetPackageConfigurationInput{})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_UpdatePackageConfiguration", func() error {
		_, err := tc.client.UpdatePackageConfiguration(tc.ctx, &iot.UpdatePackageConfigurationInput{
			VersionUpdateByJobsConfig: &iottypes.VersionUpdateByJobsConfig{
				Enabled: aws.Bool(true),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_AssociateSbomForDisassociate", func() error {
		_, err := tc.client.AssociateSbomWithPackageVersion(tc.ctx, &iot.AssociateSbomWithPackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
			Sbom: &iottypes.Sbom{
				S3Location: &iottypes.S3Location{
					Bucket: aws.String("test-bucket"),
					Key:    aws.String("sbom.json"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_ListSbomValidationResults", func() error {
		_, err := tc.client.ListSbomValidationResults(tc.ctx, &iot.ListSbomValidationResultsInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_DisassociateSbomFromPackageVersion", func() error {
		_, err := tc.client.DisassociateSbomFromPackageVersion(tc.ctx, &iot.DisassociateSbomFromPackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		return err
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_Package_Teardown", func() error {
		_, _ = tc.client.DeletePackageVersion(tc.ctx, &iot.DeletePackageVersionInput{
			PackageName: aws.String(pkgName),
			VersionName: aws.String(versionName),
		})
		_, _ = tc.client.DeletePackage(tc.ctx, &iot.DeletePackageInput{
			PackageName: aws.String(pkgName),
		})
		return nil
	}))

	return results
}

// --- Command Execution: GetCommandExecution, DeleteCommandExecution, ListCommandExecutions ---
// These operations exercise the execution record read/delete path. Since
// command executions are triggered via MQTT data-plane (not API), we test
// the NotFound/empty cases.
func (r *TestRunner) runIoTCommandExecutionGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	cmdID := uniqueName("gap-cmd")

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{CommandId: aws.String(cmdID)})
	}()

	// Setup: create command
	results = append(results, r.RunTest("iot", "Gap_CmdExec_Setup", func() error {
		_, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
			CommandId:   aws.String(cmdID),
			DisplayName: aws.String("gap test command"),
		})
		return err
	}))

	cmdARN := fmt.Sprintf("arn:aws:iot:us-east-1:000000000000:command/%s", cmdID)

	results = append(results, r.RunTest("iot", "Gap_ListCommandExecutions_Empty", func() error {
		out, err := tc.client.ListCommandExecutions(tc.ctx, &iot.ListCommandExecutionsInput{
			CommandArn: aws.String(cmdARN),
		})
		if err != nil {
			return err
		}
		if len(out.CommandExecutions) != 0 {
			return fmt.Errorf("expected 0 executions, got %d", len(out.CommandExecutions))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Gap_GetCommandExecution_NotFound", func() error {
		_, err := tc.client.GetCommandExecution(tc.ctx, &iot.GetCommandExecutionInput{
			ExecutionId: aws.String(uniqueName("nonexistent-exec")),
			TargetArn:   aws.String("arn:aws:iot:us-east-1:000000000000:thing/dummy"),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Gap_DeleteCommandExecution_NotFound", func() error {
		_, err := tc.client.DeleteCommandExecution(tc.ctx, &iot.DeleteCommandExecutionInput{
			ExecutionId: aws.String(uniqueName("nonexistent-exec")),
			TargetArn:   aws.String("arn:aws:iot:us-east-1:000000000000:thing/dummy"),
		})
		return expectNotFound(err)
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_CmdExec_Teardown", func() error {
		_, _ = tc.client.DeleteCommand(tc.ctx, &iot.DeleteCommandInput{
			CommandId: aws.String(cmdID),
		})
		return nil
	}))

	return results
}

// --- CertificateProvider: UpdateCertificateProvider ---
func (r *TestRunner) runIoTCertProviderUpdateGapTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	providerName := uniqueName("gap-cert-provider")

	// Best-effort cleanup so a failed run never leaves resources behind.
	defer func() {
		_, _ = tc.client.DeleteCertificateProvider(tc.ctx, &iot.DeleteCertificateProviderInput{CertificateProviderName: aws.String(providerName)})
	}()

	// Setup
	results = append(results, r.RunTest("iot", "Gap_CertProvider_Setup", func() error {
		_, err := tc.client.CreateCertificateProvider(tc.ctx, &iot.CreateCertificateProviderInput{
			CertificateProviderName: aws.String(providerName),
			LambdaFunctionArn:       aws.String("arn:aws:lambda:us-east-1:000000000000:function:test-cert-signer"),
			AccountDefaultForOperations: []iottypes.CertificateProviderOperation{
				iottypes.CertificateProviderOperationCreateCertificateFromCsr,
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Gap_UpdateCertificateProvider", func() error {
		_, err := tc.client.UpdateCertificateProvider(tc.ctx, &iot.UpdateCertificateProviderInput{
			CertificateProviderName: aws.String(providerName),
			LambdaFunctionArn:       aws.String("arn:aws:lambda:us-east-1:000000000000:function:updated-cert-signer"),
		})
		return err
	}))

	// Teardown
	results = append(results, r.RunTest("iot", "Gap_CertProvider_Teardown", func() error {
		_, _ = tc.client.DeleteCertificateProvider(tc.ctx, &iot.DeleteCertificateProviderInput{
			CertificateProviderName: aws.String(providerName),
		})
		return nil
	}))

	return results
}
