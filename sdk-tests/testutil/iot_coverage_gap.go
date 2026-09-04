package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTCoverageGapTests covers 25 tc.regionistered operations that previously had
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
		csr := testCreateCertificateCSRPEM
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

	// Setup: create policy + certificate through the shared cleanup-returning
	// helpers so a mid-suite failure never strands resources. A prerequisite
	// failure surfaces as one FAIL row named after the setup step it replaces.
	const gapPolicyDocument = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`
	policyCleanup, err := tc.createPolicy(policyName, gapPolicyDocument)
	if err != nil {
		return iotSetupFail("Gap_PolicyV2_Setup", err.Error())
	}
	defer policyCleanup()
	cert, certCleanup, err := tc.createCertificate(true)
	if err != nil {
		return iotSetupFail("Gap_PolicyV2_Setup", err.Error())
	}
	defer certCleanup()
	certID := cert.ID

	principalARN := func() string {
		return fmt.Sprintf("arn:aws:iot:%s:%s:cert/%s", tc.region, tc.accountID, certID)
	}

	results = append(results, r.RunTest("iot", "Gap_AttachPrincipalPolicy", func() error {
		_, err := tc.client.AttachPrincipalPolicy(tc.ctx, &iot.AttachPrincipalPolicyInput{
			PolicyName: aws.String(policyName),
			Principal:  aws.String(principalARN()),
		})
		if err != nil {
			return err
		}
		// Verify the policy is still valid after the attach operation.
		getResp, err := tc.client.GetPolicy(tc.ctx, &iot.GetPolicyInput{
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy after attach: %w", err)
		}
		if getResp.PolicyName == nil || *getResp.PolicyName != policyName {
			return fmt.Errorf("expected PolicyName=%s, got %v", policyName, getResp.PolicyName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Gap_DetachPrincipalPolicy", func() error {
		_, err := tc.client.DetachPrincipalPolicy(tc.ctx, &iot.DetachPrincipalPolicyInput{
			PolicyName: aws.String(policyName),
			Principal:  aws.String(principalARN()),
		})
		if err != nil {
			return err
		}
		listResp, err := tc.client.ListTargetsForPolicy(tc.ctx, &iot.ListTargetsForPolicyInput{
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return fmt.Errorf("ListTargetsForPolicy after detach: %w", err)
		}
		for _, t := range listResp.Targets {
			if t == principalARN() {
				return fmt.Errorf("principal %s still attached after detach", principalARN())
			}
		}
		return nil
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
				RoleArn:  aws.String(fmt.Sprintf("arn:aws:iam::%s:role/test", tc.accountID)),
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
			TargetName: aws.String("some-target"),
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

	// Setup: create thing type up front; a prerequisite failure surfaces as a
	// single FAIL row named after the setup step it replaces.
	if _, err := tc.client.CreateThingType(tc.ctx, &iot.CreateThingTypeInput{
		ThingTypeName: aws.String(thingType),
	}); err != nil {
		return iotSetupFail("Gap_SecProfile_Setup", err.Error())
	}

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
		out, err := tc.client.UpdateSecurityProfile(tc.ctx, &iot.UpdateSecurityProfileInput{
			SecurityProfileName:        aws.String(secProfileName),
			SecurityProfileDescription: aws.String("updated"),
			AdditionalMetricsToRetainV2: []iottypes.MetricToRetain{{
				Metric: aws.String("aws:num-messages-received"),
			}},
		})
		if err != nil {
			return err
		}
		// UpdateSecurityProfileResponse carries the updated profile.
		if aws.ToString(out.SecurityProfileName) != secProfileName {
			return fmt.Errorf("expected securityProfileName=%s on update response", secProfileName)
		}
		if out.Version != 2 {
			return fmt.Errorf("expected version=2 after first update, got %d", out.Version)
		}
		if len(out.AdditionalMetricsToRetainV2) != 1 || aws.ToString(out.AdditionalMetricsToRetainV2[0].Metric) != "aws:num-messages-received" {
			return fmt.Errorf("expected updated retained metric on update response, got %v", out.AdditionalMetricsToRetainV2)
		}
		return nil
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

	results = append(results, r.RunTest("iot", "SecProfile_ValidateErrorsShape", func() error {
		// An invalid behavior reports through validationErrors (a
		// ValidationError carrying errorMessage), which the SDK can only
		// decode when the documented member is used.
		out, err := tc.client.ValidateSecurityProfileBehaviors(tc.ctx, &iot.ValidateSecurityProfileBehaviorsInput{
			Behaviors: []iottypes.Behavior{{
				Name: aws.String("criteria-less"),
			}},
		})
		if err != nil {
			return err
		}
		if out.Valid {
			return fmt.Errorf("expected valid=false for a behavior without criteria")
		}
		if len(out.ValidationErrors) == 0 {
			return fmt.Errorf("expected non-empty validationErrors")
		}
		if aws.ToString(out.ValidationErrors[0].ErrorMessage) == "" {
			return fmt.Errorf("expected non-empty validationErrors[0].errorMessage")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "SecProfile_DeleteBehaviorsFlag", func() error {
		// A delete flag plus a replacement value in the same invocation is
		// rejected; the flag alone clears the stored behaviors.
		_, err := tc.client.UpdateSecurityProfile(tc.ctx, &iot.UpdateSecurityProfileInput{
			SecurityProfileName: aws.String(secProfileName),
			DeleteBehaviors:     true,
			Behaviors: []iottypes.Behavior{{
				Name: aws.String("conflicting"),
				Criteria: &iottypes.BehaviorCriteria{
					ComparisonOperator: iottypes.ComparisonOperatorLessThan,
					Value:              &iottypes.MetricValue{Count: aws.Int64(1)},
				},
			}},
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("delete+replace conflict: %w", err)
		}
		if _, err := tc.client.UpdateSecurityProfile(tc.ctx, &iot.UpdateSecurityProfileInput{
			SecurityProfileName: aws.String(secProfileName),
			DeleteBehaviors:     true,
		}); err != nil {
			return fmt.Errorf("deleteBehaviors alone failed: %w", err)
		}
		desc, err := tc.client.DescribeSecurityProfile(tc.ctx, &iot.DescribeSecurityProfileInput{
			SecurityProfileName: aws.String(secProfileName),
		})
		if err != nil {
			return err
		}
		if len(desc.Behaviors) != 0 {
			return fmt.Errorf("expected behaviors cleared, got %d", len(desc.Behaviors))
		}
		return nil
	}))

	// Job: Create -> Update
	results = append(results, r.RunTest("iot", "Gap_UpdateJob", func() error {
		_, err := tc.client.CreateJob(tc.ctx, &iot.CreateJobInput{
			JobId:    aws.String(jobID),
			Targets:  []string{fmt.Sprintf("arn:aws:iot:%s:%s:thing/dummy", tc.region, tc.accountID)},
			Document: aws.String(`{"operation":"test"}`),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.UpdateJob(tc.ctx, &iot.UpdateJobInput{
			JobId:       aws.String(jobID),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeJob(tc.ctx, &iot.DescribeJobInput{JobId: aws.String(jobID)})
		if err != nil {
			return err
		}
		if desc.Job == nil || aws.ToString(desc.Job.Description) != "updated description" {
			return fmt.Errorf("expected updated description to persist, got %v", desc.Job)
		}
		return nil
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

	// Setup: create template up front; a prerequisite failure surfaces as a
	// single FAIL row named after the setup step it replaces.
	if _, err := tc.client.CreateProvisioningTemplate(tc.ctx, &iot.CreateProvisioningTemplateInput{
		TemplateName:        aws.String(templateName),
		TemplateBody:        aws.String(templateBody),
		ProvisioningRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/test", tc.accountID)),
	}); err != nil {
		return iotSetupFail("Gap_ProvTemplateVersion_Setup", err.Error())
	}

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
		for _, v := range out.Versions {
			if v.CreationDate == nil || v.CreationDate.IsZero() {
				return fmt.Errorf("expected non-zero creationDate on version summary")
			}
			if v.IsDefaultVersion && (v.VersionId == nil || *v.VersionId != 1) {
				return fmt.Errorf("expected version 1 to be the default, got %v", v.VersionId)
			}
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

	// Setup: create package + version up front; a prerequisite failure
	// surfaces as a single FAIL row named after the setup step it replaces.
	if _, err := tc.client.CreatePackage(tc.ctx, &iot.CreatePackageInput{
		PackageName: aws.String(pkgName),
	}); err != nil {
		return iotSetupFail("Gap_Package_Setup", err.Error())
	}
	if _, err := tc.client.CreatePackageVersion(tc.ctx, &iot.CreatePackageVersionInput{
		PackageName: aws.String(pkgName),
		VersionName: aws.String(versionName),
	}); err != nil {
		return iotSetupFail("Gap_Package_Setup", err.Error())
	}

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

	// Setup: create command up front; a prerequisite failure surfaces as a
	// single FAIL row named after the setup step it replaces.
	if _, err := tc.client.CreateCommand(tc.ctx, &iot.CreateCommandInput{
		CommandId:   aws.String(cmdID),
		DisplayName: aws.String("gap test command"),
	}); err != nil {
		return iotSetupFail("Gap_CmdExec_Setup", err.Error())
	}

	cmdARN := fmt.Sprintf("arn:aws:iot:%s:%s:command/%s", tc.region, tc.accountID, cmdID)

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
			TargetArn:   aws.String(fmt.Sprintf("arn:aws:iot:%s:%s:thing/dummy", tc.region, tc.accountID)),
		})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Gap_DeleteCommandExecution_NotFound", func() error {
		_, err := tc.client.DeleteCommandExecution(tc.ctx, &iot.DeleteCommandExecutionInput{
			ExecutionId: aws.String(uniqueName("nonexistent-exec")),
			TargetArn:   aws.String(fmt.Sprintf("arn:aws:iot:%s:%s:thing/dummy", tc.region, tc.accountID)),
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

	// Setup: create the provider up front; a prerequisite failure surfaces as
	// a single FAIL row named after the setup step it replaces.
	if _, err := tc.client.CreateCertificateProvider(tc.ctx, &iot.CreateCertificateProviderInput{
		CertificateProviderName: aws.String(providerName),
		LambdaFunctionArn:       aws.String(tc.lambdaARN("test-cert-signer")),
		AccountDefaultForOperations: []iottypes.CertificateProviderOperation{
			iottypes.CertificateProviderOperationCreateCertificateFromCsr,
		},
	}); err != nil {
		return iotSetupFail("Gap_CertProvider_Setup", err.Error())
	}

	results = append(results, r.RunTest("iot", "Gap_UpdateCertificateProvider", func() error {
		_, err := tc.client.UpdateCertificateProvider(tc.ctx, &iot.UpdateCertificateProviderInput{
			CertificateProviderName: aws.String(providerName),
			LambdaFunctionArn:       aws.String(tc.lambdaARN("updated-cert-signer")),
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
