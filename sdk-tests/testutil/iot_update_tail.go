package testutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// certIDFromARN extracts the certificate ID from an IoT certificate ARN
// (arn:aws:iot:<region>:<account>:cert/<id>).
func certIDFromARN(arn string) string {
	if i := strings.LastIndex(arn, "cert/"); i >= 0 {
		return arn[i+len("cert/"):]
	}
	return arn
}

// runIoTUpdateAndTailTests covers Update operations on security/audit/index
// resources plus account logging options. Each resource is created first so
// the update targets a real record (previously these were `_ = err` smoke
// calls against hard-coded names).
func (r *TestRunner) runIoTUpdateAndTailTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	dimName := uniqueName("dim")
	mitName := uniqueName("mit")
	fmName := uniqueName("fleet-metric")
	saName := uniqueName("scheduled-audit")
	streamID := uniqueName("stream")
	thingName := uniqueName("upd-thing")
	thingGroup := uniqueName("upd-group")

	defer tc.client.DeleteDimension(tc.ctx, &iot.DeleteDimensionInput{Name: aws.String(dimName)})
	defer tc.client.DeleteMitigationAction(tc.ctx, &iot.DeleteMitigationActionInput{ActionName: aws.String(mitName)})
	defer tc.client.DeleteFleetMetric(tc.ctx, &iot.DeleteFleetMetricInput{MetricName: aws.String(fmName)})
	defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(saName)})
	defer tc.client.DeleteStream(tc.ctx, &iot.DeleteStreamInput{StreamId: aws.String(streamID)})
	defer tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{ThingGroupName: aws.String(thingGroup)})

	// Setup: create all resources in one step so updates have real targets; a
	// prerequisite failure surfaces as a single FAIL row named after the setup
	// step it replaces.
	setupFail := func(step string, err error) []TestResult {
		return []TestResult{{Service: "iot", TestName: "Update_Setup", Status: "FAIL", Error: fmt.Sprintf("%s failed: %v", step, err)}}
	}
	if _, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
		Name: aws.String(dimName), Type: iottypes.DimensionTypeTopicFilter, StringValues: []string{"x"},
	}); err != nil {
		return setupFail("CreateDimension", err)
	}
	if _, err := tc.client.CreateMitigationAction(tc.ctx, &iot.CreateMitigationActionInput{
		ActionName: aws.String(mitName), RoleArn: aws.String(tc.iamRoleARN("test")), ActionParams: &iottypes.MitigationActionParams{},
	}); err != nil {
		return setupFail("CreateMitigationAction", err)
	}
	if _, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
		MetricName: aws.String(fmName), QueryString: aws.String("*"), IndexName: aws.String("AWS_Things"),
		AggregationField: aws.String("thingName"), Period: aws.Int32(60),
		AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
	}); err != nil {
		return setupFail("CreateFleetMetric", err)
	}
	if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
		ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyDaily,
		TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
	}); err != nil {
		return setupFail("CreateScheduledAudit", err)
	}
	if _, err := tc.client.CreateStream(tc.ctx, &iot.CreateStreamInput{
		StreamId: aws.String(streamID), Files: []iottypes.StreamFile{}, RoleArn: aws.String(tc.iamRoleARN("test")),
	}); err != nil {
		return setupFail("CreateStream", err)
	}
	cleanupThing, err := tc.createThing(thingName)
	if err != nil {
		return setupFail("CreateThing", err)
	}
	defer cleanupThing()
	if _, err := tc.client.CreateThingGroup(tc.ctx, &iot.CreateThingGroupInput{ThingGroupName: aws.String(thingGroup)}); err != nil {
		return setupFail("CreateThingGroup", err)
	}

	results = append(results, r.RunTest("iot", "UpdateDimension", func() error {
		// The stringValues member is required: an update without it must be
		// rejected rather than silently skipped.
		if _, err := tc.client.UpdateDimension(tc.ctx, &iot.UpdateDimensionInput{
			Name: aws.String(dimName),
		}); err == nil {
			return fmt.Errorf("expected rejection without stringValues")
		}
		_, err := tc.client.UpdateDimension(tc.ctx, &iot.UpdateDimensionInput{
			Name: aws.String(dimName), StringValues: []string{"updated"},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeDimension(tc.ctx, &iot.DescribeDimensionInput{Name: aws.String(dimName)})
		if err != nil {
			return fmt.Errorf("DescribeDimension after update: %w", err)
		}
		if len(desc.StringValues) != 1 || desc.StringValues[0] != "updated" {
			return fmt.Errorf("expected StringValues=[updated], got %v", desc.StringValues)
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "UpdateMitigationAction", func() error {
		_, err := tc.client.UpdateMitigationAction(tc.ctx, &iot.UpdateMitigationActionInput{
			ActionName: aws.String(mitName), RoleArn: aws.String(tc.iamRoleARN("test")), ActionParams: &iottypes.MitigationActionParams{},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeMitigationAction(tc.ctx, &iot.DescribeMitigationActionInput{ActionName: aws.String(mitName)})
		if err != nil {
			return fmt.Errorf("DescribeMitigationAction after update: %w", err)
		}
		if desc.ActionName == nil || *desc.ActionName != mitName {
			return fmt.Errorf("expected ActionName=%s, got %v", mitName, desc.ActionName)
		}
		// A partial update without roleArn must keep the stored role.
		if _, err := tc.client.UpdateMitigationAction(tc.ctx, &iot.UpdateMitigationActionInput{
			ActionName: aws.String(mitName), ActionParams: &iottypes.MitigationActionParams{},
		}); err != nil {
			return fmt.Errorf("UpdateMitigationAction without roleArn: %w", err)
		}
		desc2, err := tc.client.DescribeMitigationAction(tc.ctx, &iot.DescribeMitigationActionInput{ActionName: aws.String(mitName)})
		if err != nil {
			return fmt.Errorf("DescribeMitigationAction after partial update: %w", err)
		}
		if desc2.RoleArn == nil || *desc2.RoleArn != tc.iamRoleARN("test") {
			return fmt.Errorf("expected roleArn preserved after partial update, got %v", desc2.RoleArn)
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "UpdateFleetMetric", func() error {
		_, err := tc.client.UpdateFleetMetric(tc.ctx, &iot.UpdateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(120),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeFleetMetric(tc.ctx, &iot.DescribeFleetMetricInput{MetricName: aws.String(fmName)})
		if err != nil {
			return fmt.Errorf("DescribeFleetMetric after update: %w", err)
		}
		if desc.Period == nil || *desc.Period != 120 {
			return fmt.Errorf("expected Period=120 after update, got %v", desc.Period)
		}
		if desc.Version != 2 {
			return fmt.Errorf("expected Version=2 after the first update, got %d", desc.Version)
		}
		// Omitted optional members are preserved: an update without period
		// must leave the stored period and queryString intact.
		if _, err := tc.client.UpdateFleetMetric(tc.ctx, &iot.UpdateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
		}); err != nil {
			return fmt.Errorf("period-omitting update failed: %w", err)
		}
		desc, err = tc.client.DescribeFleetMetric(tc.ctx, &iot.DescribeFleetMetricInput{MetricName: aws.String(fmName)})
		if err != nil {
			return fmt.Errorf("DescribeFleetMetric after omission update: %w", err)
		}
		if desc.Period == nil || *desc.Period != 120 {
			return fmt.Errorf("omitted period must be preserved, got %v", desc.Period)
		}
		if aws.ToString(desc.QueryString) != "thingName:*" {
			return fmt.Errorf("expected queryString from the omission update, got %v", desc.QueryString)
		}
		if desc.Version != 3 {
			return fmt.Errorf("expected Version=3 after the second update, got %d", desc.Version)
		}
		// A mismatched expectedVersion is rejected with the documented
		// VersionConflictException.
		_, err = tc.client.UpdateFleetMetric(tc.ctx, &iot.UpdateFleetMetricInput{
			MetricName: aws.String(fmName), IndexName: aws.String("AWS_Things"), ExpectedVersion: aws.Int64(99),
		})
		if vErr := expectAWSErrorCode(err, "VersionConflictException"); vErr != nil {
			return vErr
		}
		if _, err := tc.client.UpdateFleetMetric(tc.ctx, &iot.UpdateFleetMetricInput{
			MetricName: aws.String(fmName), IndexName: aws.String("AWS_Things"),
			ExpectedVersion: aws.Int64(3), Period: aws.Int32(300),
		}); err != nil {
			return fmt.Errorf("expectedVersion=3 update failed: %w", err)
		}
		desc, err = tc.client.DescribeFleetMetric(tc.ctx, &iot.DescribeFleetMetricInput{MetricName: aws.String(fmName)})
		if err != nil {
			return fmt.Errorf("DescribeFleetMetric after CAS update: %w", err)
		}
		if desc.Period == nil || *desc.Period != 300 {
			return fmt.Errorf("expected Period=300 after the CAS update, got %v", desc.Period)
		}
		if desc.Version != 4 {
			return fmt.Errorf("expected Version=4 after the CAS update, got %d", desc.Version)
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "UpdateScheduledAudit", func() error {
		_, err := tc.client.UpdateScheduledAudit(tc.ctx, &iot.UpdateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyWeekly,
			DayOfWeek:        iottypes.DayOfWeek("MON"),
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{ScheduledAuditName: aws.String(saName)})
		if err != nil {
			return fmt.Errorf("DescribeScheduledAudit after update: %w", err)
		}
		if desc.Frequency != iottypes.AuditFrequencyWeekly {
			return fmt.Errorf("expected Frequency=Weekly after update, got %v", desc.Frequency)
		}
		// A partial update supplying only the check list must keep the
		// frequency and day of week.
		if _, err := tc.client.UpdateScheduledAudit(tc.ctx, &iot.UpdateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("UpdateScheduledAudit partial update: %w", err)
		}
		desc2, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{ScheduledAuditName: aws.String(saName)})
		if err != nil {
			return fmt.Errorf("DescribeScheduledAudit after partial update: %w", err)
		}
		if desc2.Frequency != iottypes.AuditFrequencyWeekly {
			return fmt.Errorf("expected Frequency=Weekly preserved, got %v", desc2.Frequency)
		}
		if desc2.DayOfWeek != iottypes.DayOfWeek("MON") {
			return fmt.Errorf("expected DayOfWeek=MON preserved, got %v", desc2.DayOfWeek)
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "UpdateStream", func() error {
		_, err := tc.client.UpdateStream(tc.ctx, &iot.UpdateStreamInput{StreamId: aws.String(streamID), Files: []iottypes.StreamFile{}})
		if err != nil {
			return err
		}
		desc, err := tc.client.DescribeStream(tc.ctx, &iot.DescribeStreamInput{StreamId: aws.String(streamID)})
		if err != nil {
			return fmt.Errorf("DescribeStream after update: %w", err)
		}
		if desc.StreamInfo == nil || desc.StreamInfo.StreamId == nil || *desc.StreamInfo.StreamId != streamID {
			return fmt.Errorf("expected StreamId=%s after update", streamID)
		}
		return nil
	}))
	results = append(results, r.RunTest("iot", "UpdateThingGroupsForThing", func() error {
		_, err := tc.client.UpdateThingGroupsForThing(tc.ctx, &iot.UpdateThingGroupsForThingInput{ThingName: aws.String(thingName)})
		return err
	}))

	// Account logging options and index reads.
	results = append(results, r.RunTest("iot", "SetV2LoggingOptions", func() error {
		_, err := tc.client.SetV2LoggingOptions(tc.ctx, &iot.SetV2LoggingOptionsInput{})
		return err
	}))
	results = append(results, r.RunTest("iot", "ListV2LoggingLevels", func() error {
		_, err := tc.client.ListV2LoggingLevels(tc.ctx, &iot.ListV2LoggingLevelsInput{})
		return err
	}))
	results = append(results, r.RunTest("iot", "GetPercentiles", func() error {
		_, err := tc.client.GetPercentiles(tc.ctx, &iot.GetPercentilesInput{IndexName: aws.String("AWS_Things"), QueryString: aws.String("*")})
		return err
	}))
	// RegisterThing with an invalid template body must be rejected.
	results = append(results, r.RunTest("iot", "RegisterThing_InvalidTemplate", func() error {
		_, err := tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{
			TemplateBody: aws.String("{}"),
			Parameters:   map[string]string{"x": "y"},
		})
		return expectValidationError(err)
	}))

	// RegisterThing with a minimal valid provisioning template must
	// provision a thing; the parameters map carries the thing name.
	registeredName := uniqueName("registered-thing")
	results = append(results, r.RunTest("iot", "RegisterThing_MinimalTemplate", func() error {
		templateBody := `{"Parameters":{"ThingName":{"Type":"String"}},"Resources":{"Thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName": {"Ref": "{{ThingName}}"}}}}}`
		out, err := tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{
			TemplateBody: aws.String(templateBody),
			Parameters:   map[string]string{"ThingName": registeredName},
		})
		if err != nil {
			return fmt.Errorf("RegisterThing failed: %w", err)
		}
		thingARN, ok := out.ResourceArns["Thing"]
		if !ok || thingARN == "" {
			return fmt.Errorf("expected resourceArns to carry the Thing ARN, got %v", out.ResourceArns)
		}
		if !strings.Contains(thingARN, registeredName) {
			return fmt.Errorf("expected Thing ARN to name %s, got %s", registeredName, thingARN)
		}
		// The provisioned thing must be visible to DescribeThing.
		if _, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(registeredName)}); err != nil {
			return fmt.Errorf("DescribeThing on the registered thing failed: %w", err)
		}
		return nil
	}))
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(registeredName)})

	// A provisioning template identifies its resources by Type; the logical
	// name is caller-chosen (the AWS guide's own examples use "thing").
	results = append(results, r.RunTest("iot", "RegisterThing_LogicalNameArbitrary", func() error {
		lowerName := uniqueName("lower-thing")
		templateBody := `{"Parameters":{"ThingName":{"Type":"String"}},"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"}}}}}`
		out, err := tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{
			TemplateBody: aws.String(templateBody),
			Parameters:   map[string]string{"ThingName": lowerName},
		})
		if err != nil {
			return fmt.Errorf("RegisterThing with lowercase logical name failed: %w", err)
		}
		thingARN, ok := out.ResourceArns["thing"]
		if !ok || thingARN == "" {
			return fmt.Errorf("expected resourceArns keyed by the logical name \"thing\", got %v", out.ResourceArns)
		}
		if !strings.Contains(thingARN, lowerName) {
			return fmt.Errorf("expected thing ARN to name %s, got %s", lowerName, thingARN)
		}
		defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(lowerName)})
		if _, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(lowerName)}); err != nil {
			return fmt.Errorf("DescribeThing on the registered thing failed: %w", err)
		}
		return nil
	}))

	// The full template provisions a thing plus a certificate signed from
	// the CSR parameter and a policy created from the inline document; the
	// response carries certificatePem and ARNs for all three resources
	// under their logical names.
	results = append(results, r.RunTest("iot", "RegisterThing_FullTemplateWithCSR", func() error {
		fullName := uniqueName("full-registered")
		policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`
		templateBody := fmt.Sprintf(`{"Parameters":{"ThingName":{"Type":"String"},"CSR":{"Type":"String"}},"Resources":{`+
			`"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":{"Ref":"{{ThingName}}"},"AttributePayload":{"provisioned":"engine"}}},`+
			`"certificate":{"Type":"AWS::IoT::Certificate","Properties":{"CertificateSigningRequest":{"Ref":"CSR"},"Status":"ACTIVE","ThingPrincipalType":"NON_EXCLUSIVE_THING"}},`+
			`"policy":{"Type":"AWS::IoT::Policy","Properties":{"PolicyDocument":%q}}}}`, policyDoc)
		out, err := tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{
			TemplateBody: aws.String(templateBody),
			Parameters:   map[string]string{"ThingName": fullName, "CSR": testCSRPEM},
		})
		if err != nil {
			return fmt.Errorf("RegisterThing with a full template failed: %w", err)
		}
		certARN := out.ResourceArns["certificate"]
		policyARN := out.ResourceArns["policy"]
		if out.CertificatePem == nil || *out.CertificatePem == "" {
			return fmt.Errorf("expected certificatePem in the response")
		}
		if certARN == "" || policyARN == "" || out.ResourceArns["thing"] == "" {
			return fmt.Errorf("expected resourceArns for thing/certificate/policy, got %v", out.ResourceArns)
		}
		// Cleanup: detach both attachments, then remove the resources.
		sum := sha256.Sum256([]byte(policyDoc))
		policyName := hex.EncodeToString(sum[:])[:32]
		defer func() {
			_, _ = tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{PolicyName: aws.String(policyName)})
		}()
		defer func() { _, _ = tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(fullName)}) }()
		defer func() {
			_, _ = tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
				CertificateId: aws.String(certIDFromARN(certARN)), NewStatus: iottypes.CertificateStatusInactive,
			})
			_, _ = tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{CertificateId: aws.String(certIDFromARN(certARN))})
		}()
		defer func() {
			_, _ = tc.client.DetachPolicy(tc.ctx, &iot.DetachPolicyInput{PolicyName: aws.String(policyName), Target: aws.String(certARN)})
		}()
		defer func() {
			_, _ = tc.client.DetachThingPrincipal(tc.ctx, &iot.DetachThingPrincipalInput{ThingName: aws.String(fullName), Principal: aws.String(certARN)})
		}()

		certID := certIDFromARN(certARN)
		desc, err := tc.client.DescribeCertificate(tc.ctx, &iot.DescribeCertificateInput{CertificateId: aws.String(certID)})
		if err != nil {
			return fmt.Errorf("DescribeCertificate on the provisioned certificate failed: %w", err)
		}
		if desc.CertificateDescription == nil || desc.CertificateDescription.Status != iottypes.CertificateStatusActive {
			return fmt.Errorf("expected provisioned certificate ACTIVE, got %+v", desc.CertificateDescription)
		}
		lp, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(fullName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals failed: %w", err)
		}
		if !containsString(lp.Principals, certARN) {
			return fmt.Errorf("expected the provisioned certificate among the thing principals, got %v", lp.Principals)
		}
		pp, err := tc.client.ListPrincipalPolicies(tc.ctx, &iot.ListPrincipalPoliciesInput{Principal: aws.String(certARN)})
		if err != nil {
			return fmt.Errorf("ListPrincipalPolicies failed: %w", err)
		}
		found := false
		for _, p := range pp.Policies {
			if aws.ToString(p.PolicyName) == policyName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected the hash-named policy %s attached to the certificate", policyName)
		}
		thingDesc, err := tc.client.DescribeThing(tc.ctx, &iot.DescribeThingInput{ThingName: aws.String(fullName)})
		if err != nil {
			return fmt.Errorf("DescribeThing failed: %w", err)
		}
		if thingDesc.Attributes["provisioned"] != "engine" {
			return fmt.Errorf("expected AttributePayload from the template, got %v", thingDesc.Attributes)
		}
		return nil
	}))

	// A CertificateId declaration reuses an existing certificate and
	// attaches it to the provisioned thing.
	results = append(results, r.RunTest("iot", "RegisterThing_ExistingCertificateId", func() error {
		cert, cleanup, err := tc.createCertificate(true)
		if err != nil {
			return err
		}
		defer cleanup()
		idName := uniqueName("id-registered")
		templateBody := fmt.Sprintf(`{"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":"%s"}},"certificate":{"Type":"AWS::IoT::Certificate","Properties":{"CertificateId":"%s"}}}}`, idName, cert.ID)
		out, err := tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{TemplateBody: aws.String(templateBody)})
		if err != nil {
			return fmt.Errorf("RegisterThing with a CertificateId declaration failed: %w", err)
		}
		if out.ResourceArns["certificate"] != cert.ARN {
			return fmt.Errorf("expected the existing certificate ARN %s, got %v", cert.ARN, out.ResourceArns)
		}
		defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(idName)})
		lp, err := tc.client.ListThingPrincipals(tc.ctx, &iot.ListThingPrincipalsInput{ThingName: aws.String(idName)})
		if err != nil {
			return fmt.Errorf("ListThingPrincipals failed: %w", err)
		}
		if !containsString(lp.Principals, cert.ARN) {
			return fmt.Errorf("expected the reused certificate attached to the thing, got %v", lp.Principals)
		}
		return nil
	}))

	// An existing thing without OverrideSettings fails the registration
	// with ResourceConflictsException (the documented FAIL default).
	results = append(results, r.RunTest("iot", "RegisterThing_ConflictingThingRejected", func() error {
		conflictName := uniqueName("conflict-thing")
		cleanup, err := tc.createThing(conflictName)
		if err != nil {
			return err
		}
		defer cleanup()
		templateBody := fmt.Sprintf(`{"Resources":{"thing":{"Type":"AWS::IoT::Thing","Properties":{"ThingName":"%s"}}}}`, conflictName)
		_, err = tc.client.RegisterThing(tc.ctx, &iot.RegisterThingInput{TemplateBody: aws.String(templateBody)})
		return expectAWSErrorCode(err, "ResourceConflictsException")
	}))

	return results
}
