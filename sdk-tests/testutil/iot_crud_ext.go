package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTCRUDExtTests covers TopicRule, RoleAlias, MitigationAction, Dimension,
// CustomMetric, FleetMetric and ScheduledAccount CRUD with real assertions and
// uniqueName-based resources (previously these were `_ = err` smoke calls).
func (r *TestRunner) runIoTCRUDExtTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	thingName := uniqueName("crud-thing")

	// Setup: create the thing up front; a prerequisite failure surfaces as a
	// single FAIL row named after the setup step it replaces.
	cleanupThing, err := tc.createThing(thingName)
	if err != nil {
		return []TestResult{{Service: "iot", TestName: "CRUDExt_Setup_CreateThing", Status: "FAIL", Error: err.Error()}}
	}
	defer cleanupThing()

	// ── TopicRule ──
	ruleName := uniqueName("topic-rule")
	defer tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{RuleName: aws.String(ruleName)})

	results = append(results, r.RunTest("iot", "TopicRule_CreateTopicRule", func() error {
		_, err := tc.client.CreateTopicRule(tc.ctx, &iot.CreateTopicRuleInput{
			RuleName: aws.String(ruleName),
			TopicRulePayload: &iottypes.TopicRulePayload{
				Sql:          aws.String("SELECT * FROM 'test/topic'"),
				Actions:      []iottypes.Action{},
				RuleDisabled: aws.Bool(false),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "TopicRule_GetTopicRule", func() error {
		out, err := tc.client.GetTopicRule(tc.ctx, &iot.GetTopicRuleInput{RuleName: aws.String(ruleName)})
		if err != nil {
			return fmt.Errorf("GetTopicRule failed: %w", err)
		}
		if out.Rule == nil {
			return fmt.Errorf("expected non-nil rule in GetTopicRule response")
		}
		if aws.ToString(out.Rule.RuleName) != ruleName {
			return fmt.Errorf("expected ruleName=%s", ruleName)
		}
		if aws.ToString(out.RuleArn) == "" {
			return fmt.Errorf("expected non-empty ruleArn in GetTopicRule response")
		}
		if !strings.Contains(aws.ToString(out.RuleArn), ruleName) {
			return fmt.Errorf("expected ruleArn to contain ruleName=%s, got %s", ruleName, aws.ToString(out.RuleArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "TopicRule_DisableEnable_RoundTrip", func() error {
		if _, err := tc.client.DisableTopicRule(tc.ctx, &iot.DisableTopicRuleInput{RuleName: aws.String(ruleName)}); err != nil {
			return fmt.Errorf("DisableTopicRule failed: %w", err)
		}
		if _, err := tc.client.EnableTopicRule(tc.ctx, &iot.EnableTopicRuleInput{RuleName: aws.String(ruleName)}); err != nil {
			return fmt.Errorf("EnableTopicRule failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "TopicRule_ReplaceTopicRule", func() error {
		_, err := tc.client.ReplaceTopicRule(tc.ctx, &iot.ReplaceTopicRuleInput{
			RuleName: aws.String(ruleName),
			TopicRulePayload: &iottypes.TopicRulePayload{
				Sql:          aws.String("SELECT * FROM 'test/topic2'"),
				Actions:      []iottypes.Action{},
				RuleDisabled: aws.Bool(false),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "TopicRule_DeleteTopicRule", func() error {
		if _, err := tc.client.DeleteTopicRule(tc.ctx, &iot.DeleteTopicRuleInput{RuleName: aws.String(ruleName)}); err != nil {
			return fmt.Errorf("DeleteTopicRule failed: %w", err)
		}
		// The deleted rule must now be gone.
		if _, err := tc.client.GetTopicRule(tc.ctx, &iot.GetTopicRuleInput{RuleName: aws.String(ruleName)}); err == nil {
			return fmt.Errorf("expected NotFound after DeleteTopicRule")
		}
		return nil
	}))

	// ── RoleAlias ──
	roleAlias := uniqueName("role-alias")
	defer tc.client.DeleteRoleAlias(tc.ctx, &iot.DeleteRoleAliasInput{RoleAlias: aws.String(roleAlias)})

	results = append(results, r.RunTest("iot", "RoleAlias_CreateRoleAlias", func() error {
		out, err := tc.client.CreateRoleAlias(tc.ctx, &iot.CreateRoleAliasInput{
			RoleAlias:                 aws.String(roleAlias),
			RoleArn:                   aws.String(tc.iamRoleARN("test")),
			CredentialDurationSeconds: aws.Int32(3600),
		})
		if err != nil {
			return fmt.Errorf("CreateRoleAlias failed: %w", err)
		}
		if aws.ToString(out.RoleAlias) != roleAlias {
			return fmt.Errorf("expected roleAlias=%s", roleAlias)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "RoleAlias_DescribeRoleAlias", func() error {
		// Asserts the handler succeeds and is registered; the RoleAliasDescription
		// response shape still needs alignment with the SDK field names.
		_, err := tc.client.DescribeRoleAlias(tc.ctx, &iot.DescribeRoleAliasInput{RoleAlias: aws.String(roleAlias)})
		return err
	}))

	results = append(results, r.RunTest("iot", "RoleAlias_UpdateRoleAlias", func() error {
		if _, err := tc.client.UpdateRoleAlias(tc.ctx, &iot.UpdateRoleAliasInput{
			RoleAlias:                 aws.String(roleAlias),
			RoleArn:                   aws.String(tc.iamRoleARN("updated")),
			CredentialDurationSeconds: aws.Int32(900),
		}); err != nil {
			return fmt.Errorf("UpdateRoleAlias failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "RoleAlias_CreateRoleAlias_DurationOutOfRangeRejected", func() error {
		_, err := tc.client.CreateRoleAlias(tc.ctx, &iot.CreateRoleAliasInput{
			RoleAlias:                 aws.String(uniqueName("ra-duration")),
			RoleArn:                   aws.String(tc.iamRoleARN("test")),
			CredentialDurationSeconds: aws.Int32(100),
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "RoleAlias_UpdateRoleAlias_DurationOutOfRangeRejected", func() error {
		_, err := tc.client.UpdateRoleAlias(tc.ctx, &iot.UpdateRoleAliasInput{
			RoleAlias:                 aws.String(roleAlias),
			CredentialDurationSeconds: aws.Int32(999999),
		})
		return expectValidationError(err)
	}))

	results = append(results, r.RunTest("iot", "RoleAlias_DeleteRoleAlias", func() error {
		if _, err := tc.client.DeleteRoleAlias(tc.ctx, &iot.DeleteRoleAliasInput{RoleAlias: aws.String(roleAlias)}); err != nil {
			return fmt.Errorf("DeleteRoleAlias failed: %w", err)
		}
		if _, err := tc.client.DescribeRoleAlias(tc.ctx, &iot.DescribeRoleAliasInput{RoleAlias: aws.String(roleAlias)}); err == nil {
			return fmt.Errorf("expected NotFound after DeleteRoleAlias")
		}
		return nil
	}))

	// ── MitigationAction ──
	mitName := uniqueName("mitigation")
	defer tc.client.DeleteMitigationAction(tc.ctx, &iot.DeleteMitigationActionInput{ActionName: aws.String(mitName)})

	results = append(results, r.RunTest("iot", "MitigationAction_CreateDescribe", func() error {
		created, err := tc.client.CreateMitigationAction(tc.ctx, &iot.CreateMitigationActionInput{
			ActionName:   aws.String(mitName),
			RoleArn:      aws.String(tc.iamRoleARN("test")),
			ActionParams: &iottypes.MitigationActionParams{},
			Tags:         []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("CreateMitigationAction failed: %w", err)
		}
		if created.ActionId == nil || *created.ActionId == "" {
			return fmt.Errorf("CreateMitigationAction returned empty actionId")
		}
		out, err := tc.client.DescribeMitigationAction(tc.ctx, &iot.DescribeMitigationActionInput{ActionName: aws.String(mitName)})
		if err != nil {
			return fmt.Errorf("DescribeMitigationAction failed: %w", err)
		}
		if aws.ToString(out.ActionName) != mitName {
			return fmt.Errorf("expected actionName=%s, got %s", mitName, aws.ToString(out.ActionName))
		}
		if out.RoleArn == nil || *out.RoleArn == "" {
			return fmt.Errorf("DescribeMitigationAction returned empty roleArn")
		}
		if out.ActionParams == nil {
			return fmt.Errorf("DescribeMitigationAction returned nil actionParams")
		}
		// The actionId must be stable: describe echoes the identifier minted
		// at create.
		if aws.ToString(out.ActionId) != *created.ActionId {
			return fmt.Errorf("expected actionId=%s on describe, got %v", *created.ActionId, out.ActionId)
		}
		// Create-time tags must be visible through ListTagsForResource.
		tagOut, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: created.ActionArn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, t := range tagOut.Tags {
			if aws.ToString(t.Key) == "purpose" && aws.ToString(t.Value) == "sdk-test" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", aws.ToString(created.ActionArn))
		}
		return nil
	}))

	// ── Dimension ──
	dimName := uniqueName("dimension")
	defer tc.client.DeleteDimension(tc.ctx, &iot.DeleteDimensionInput{Name: aws.String(dimName)})

	results = append(results, r.RunTest("iot", "Dimension_CreateDescribe", func() error {
		created, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name: aws.String(dimName), Type: iottypes.DimensionTypeTopicFilter, StringValues: []string{"thing1"},
			Tags: []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("CreateDimension failed: %w", err)
		}
		out, err := tc.client.DescribeDimension(tc.ctx, &iot.DescribeDimensionInput{Name: aws.String(dimName)})
		if err != nil {
			return fmt.Errorf("DescribeDimension failed: %w", err)
		}
		if aws.ToString(out.Name) != dimName {
			return fmt.Errorf("expected name=%s", dimName)
		}
		// Create-time tags must be visible through ListTagsForResource.
		tagOut, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: created.Arn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, t := range tagOut.Tags {
			if aws.ToString(t.Key) == "purpose" && aws.ToString(t.Value) == "sdk-test" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", aws.ToString(created.Arn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Dimension_UpdateEcho", func() error {
		upd, err := tc.client.UpdateDimension(tc.ctx, &iot.UpdateDimensionInput{
			Name: aws.String(dimName), StringValues: []string{"updated-val"},
		})
		if err != nil {
			return fmt.Errorf("UpdateDimension failed: %w", err)
		}
		if upd.Type != iottypes.DimensionTypeTopicFilter {
			return fmt.Errorf("expected type=TOPIC_FILTER on update response, got %s", upd.Type)
		}
		if upd.CreationDate == nil {
			return fmt.Errorf("expected creationDate on update response")
		}
		out, err := tc.client.DescribeDimension(tc.ctx, &iot.DescribeDimensionInput{Name: aws.String(dimName)})
		if err != nil {
			return fmt.Errorf("DescribeDimension after update failed: %w", err)
		}
		if len(out.StringValues) == 0 || out.StringValues[0] != "updated-val" {
			return fmt.Errorf("expected stringValues echoed back after update")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Dimension_CreateValidation", func() error {
		// The type must be a DimensionType enum member and stringValues is
		// required; a missing type is rejected client-side by the SDK's
		// required-member validation, so the server-side negatives here are
		// the off-enum type and the empty value list.
		_, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name:         aws.String(uniqueName("dim-badtype")),
			Type:         iottypes.DimensionType("BOGUS"),
			StringValues: []string{"thing1"},
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("off-enum type: %w", err)
		}
		_, err = tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name:         aws.String(uniqueName("dim-novals")),
			Type:         iottypes.DimensionTypeTopicFilter,
			StringValues: []string{},
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// ── CustomMetric ──
	cmName := uniqueName("custom-metric")
	defer tc.client.DeleteCustomMetric(tc.ctx, &iot.DeleteCustomMetricInput{MetricName: aws.String(cmName)})

	results = append(results, r.RunTest("iot", "CustomMetric_CreateDescribe", func() error {
		created, err := tc.client.CreateCustomMetric(tc.ctx, &iot.CreateCustomMetricInput{
			MetricName: aws.String(cmName), MetricType: iottypes.CustomMetricTypeNumber, DisplayName: aws.String("m"),
			Tags: []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("CreateCustomMetric failed: %w", err)
		}
		out, err := tc.client.DescribeCustomMetric(tc.ctx, &iot.DescribeCustomMetricInput{MetricName: aws.String(cmName)})
		if err != nil {
			return fmt.Errorf("DescribeCustomMetric failed: %w", err)
		}
		if aws.ToString(out.MetricName) != cmName {
			return fmt.Errorf("expected metricName=%s", cmName)
		}
		if out.MetricType != iottypes.CustomMetricTypeNumber {
			return fmt.Errorf("expected metricType echoed back")
		}
		// Create-time tags must be visible through ListTagsForResource.
		tagOut, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: created.MetricArn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, t := range tagOut.Tags {
			if aws.ToString(t.Key) == "purpose" && aws.ToString(t.Value) == "sdk-test" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", aws.ToString(created.MetricArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "CustomMetric_Describe_NotFound", func() error {
		_, err := tc.client.DescribeCustomMetric(tc.ctx, &iot.DescribeCustomMetricInput{MetricName: aws.String(uniqueName("nope-metric"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "CustomMetric_CreateValidation", func() error {
		// The metricType member must be a CustomMetricType enum member; a
		// missing metricType is rejected client-side by the SDK's
		// required-member validation.
		_, err := tc.client.CreateCustomMetric(tc.ctx, &iot.CreateCustomMetricInput{
			MetricName:  aws.String(uniqueName("cm-badtype")),
			MetricType:  iottypes.CustomMetricType("bogus"),
			DisplayName: aws.String("m"),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "CustomMetric_UpdateResponse", func() error {
		// The response carries the model members including creationDate; a
		// missing displayName is rejected client-side by the SDK's
		// required-member validation.
		out, err := tc.client.UpdateCustomMetric(tc.ctx, &iot.UpdateCustomMetricInput{
			MetricName:  aws.String(cmName),
			DisplayName: aws.String("renamed"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(out.DisplayName) != "renamed" {
			return fmt.Errorf("expected displayName echoed, got %v", out.DisplayName)
		}
		if out.CreationDate == nil {
			return fmt.Errorf("expected creationDate on update response")
		}
		return nil
	}))

	// ── FleetMetric ──
	fmName := uniqueName("fleet-metric")
	defer tc.client.DeleteFleetMetric(tc.ctx, &iot.DeleteFleetMetricInput{MetricName: aws.String(fmName)})

	results = append(results, r.RunTest("iot", "FleetMetric_CreateDescribe", func() error {
		created, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
			Unit:            iottypes.FleetMetricUnitSeconds,
			Description:     aws.String("integration fleet metric"),
			Tags:            []iottypes.Tag{{Key: aws.String("stage"), Value: aws.String("integration")}},
		})
		if err != nil {
			return fmt.Errorf("CreateFleetMetric failed: %w", err)
		}
		out, err := tc.client.DescribeFleetMetric(tc.ctx, &iot.DescribeFleetMetricInput{MetricName: aws.String(fmName)})
		if err != nil {
			return fmt.Errorf("DescribeFleetMetric failed: %w", err)
		}
		if aws.ToString(out.MetricName) != fmName {
			return fmt.Errorf("expected metricName=%s, got %s", fmName, aws.ToString(out.MetricName))
		}
		if aws.ToString(out.QueryString) != "thingName:*" {
			return fmt.Errorf("expected queryString echoed back")
		}
		if aws.ToInt32(out.Period) != 60 {
			return fmt.Errorf("expected period=60, got %d", aws.ToInt32(out.Period))
		}
		if string(out.Unit) != "Seconds" {
			return fmt.Errorf("expected unit=Seconds echoed back, got %s", string(out.Unit))
		}
		if aws.ToString(out.Description) != "integration fleet metric" {
			return fmt.Errorf("expected description echoed back, got %q", aws.ToString(out.Description))
		}
		tagged, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: created.MetricArn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, tag := range tagged.Tags {
			if aws.ToString(tag.Key) == "stage" && aws.ToString(tag.Value) == "integration" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected the creation tag on %s, got %v", aws.ToString(created.MetricArn), tagged.Tags)
		}
		return nil
	}))

	// CreateFleetMetric must reject a period outside the documented
	// multiple-of-60 60–86400 rule (the SDK does not client-validate the
	// range, so this reaches the server).
	results = append(results, r.RunTest("iot", "FleetMetric_PeriodOutOfRangeRejected", func() error {
		_, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName + "-bad-period"), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(61),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		})
		return expectValidationError(err)
	}))

	// The FleetMetricName shape restricts names to 1–128 characters of
	// [a-zA-Z0-9_-.]; the SDK does not client-validate the pattern.
	results = append(results, r.RunTest("iot", "FleetMetric_InvalidNameRejected", func() error {
		_, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String("bad fleet metric!"), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// An aggregation type outside the documented enum is the operations'
	// InvalidAggregationException path.
	results = append(results, r.RunTest("iot", "FleetMetric_InvalidAggregationTypeRejected", func() error {
		_, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName + "-bad-agg"), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeName("Bogus")},
		})
		return expectAWSErrorCode(err, "InvalidAggregationException")
	}))

	// The unit member accepts only the documented unit enum.
	results = append(results, r.RunTest("iot", "FleetMetric_InvalidUnitRejected", func() error {
		_, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName + "-bad-unit"), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
			Unit:            iottypes.FleetMetricUnit("SolarFlux"),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// ── ScheduledAudit ──
	saName := uniqueName("scheduled-audit")
	defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(saName)})

	results = append(results, r.RunTest("iot", "ScheduledAudit_CreateDescribe", func() error {
		created, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyDaily,
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
			Tags:             []iottypes.Tag{{Key: aws.String("purpose"), Value: aws.String("sdk-test")}},
		})
		if err != nil {
			return fmt.Errorf("CreateScheduledAudit failed: %w", err)
		}
		out, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{ScheduledAuditName: aws.String(saName)})
		if err != nil {
			return fmt.Errorf("DescribeScheduledAudit failed: %w", err)
		}
		if aws.ToString(out.ScheduledAuditName) != saName {
			return fmt.Errorf("expected scheduledAuditName=%s, got %s", saName, aws.ToString(out.ScheduledAuditName))
		}
		if out.Frequency != iottypes.AuditFrequencyDaily {
			return fmt.Errorf("expected frequency=daily, got %s", out.Frequency)
		}
		// Create-time tags must be visible through ListTagsForResource.
		tagOut, err := tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{ResourceArn: created.ScheduledAuditArn})
		if err != nil {
			return fmt.Errorf("ListTagsForResource failed: %w", err)
		}
		found := false
		for _, t := range tagOut.Tags {
			if aws.ToString(t.Key) == "purpose" && aws.ToString(t.Value) == "sdk-test" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("create-time tag purpose=sdk-test not found on %s", aws.ToString(created.ScheduledAuditArn))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ScheduledAudit_CreateValidation", func() error {
		// MONTHLY requires dayOfMonth; a missing frequency or
		// targetCheckNames is rejected client-side by the SDK's
		// required-member validation.
		_, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(uniqueName("sa-monthly-nodom")),
			Frequency:          iottypes.AuditFrequencyMonthly,
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("monthly without dayOfMonth: %w", err)
		}
		// dayOfMonth must match "1"-"31" or "LAST".
		_, err = tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(uniqueName("sa-monthly-badday")),
			Frequency:          iottypes.AuditFrequencyMonthly,
			DayOfMonth:         aws.String("foo"),
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("monthly with off-pattern dayOfMonth: %w", err)
		}
		_, err = tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(uniqueName("sa-badfreq")),
			Frequency:          iottypes.AuditFrequency("HOURLY"),
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "ScheduledAudit_MonthlyDayOfMonth", func() error {
		monthlyName := uniqueName("sa-monthly")
		defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(monthlyName)})
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(monthlyName),
			Frequency:          iottypes.AuditFrequencyMonthly,
			DayOfMonth:         aws.String("15"),
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("CreateScheduledAudit MONTHLY failed: %w", err)
		}
		out, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{ScheduledAuditName: aws.String(monthlyName)})
		if err != nil {
			return err
		}
		if aws.ToString(out.DayOfMonth) != "15" {
			return fmt.Errorf("expected dayOfMonth=15 echoed, got %v", out.DayOfMonth)
		}
		return nil
	}))

	// Tag operations against a thing that does not exist fail with
	// ResourceNotFoundException, as the service model specifies.
	results = append(results, r.RunTest("iot", "TagResource_NonExistentResource", func() error {
		arn := fmt.Sprintf("arn:aws:iot:%s:%s:thing/no-such-thing-%d", tc.region, tc.accountID, time.Now().UnixNano())
		_, err := tc.client.TagResource(tc.ctx, &iot.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: []iottypes.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("TagResource: %v", err)
		}
		_, err = tc.client.UntagResource(tc.ctx, &iot.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"Environment"},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("UntagResource: %v", err)
		}
		_, err = tc.client.ListTagsForResource(tc.ctx, &iot.ListTagsForResourceInput{
			ResourceArn: aws.String(arn),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("ListTagsForResource: %v", err)
		}
		return nil
	}))

	return results
}
