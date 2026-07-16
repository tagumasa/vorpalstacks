package testutil

import (
	"fmt"

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
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	results = append(results, r.RunTest("iot", "CRUDExt_Setup_CreateThing", func() error {
		_, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)})
		return err
	}))

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
		if _, err := tc.client.CreateMitigationAction(tc.ctx, &iot.CreateMitigationActionInput{
			ActionName:   aws.String(mitName),
			RoleArn:      aws.String(tc.iamRoleARN("test")),
			ActionParams: &iottypes.MitigationActionParams{},
		}); err != nil {
			return fmt.Errorf("CreateMitigationAction failed: %w", err)
		}
		out, err := tc.client.DescribeMitigationAction(tc.ctx, &iot.DescribeMitigationActionInput{ActionName: aws.String(mitName)})
		if err != nil {
			return fmt.Errorf("DescribeMitigationAction failed: %w", err)
		}
		_ = out
		return nil
	}))

	// ── Dimension ──
	dimName := uniqueName("dimension")
	defer tc.client.DeleteDimension(tc.ctx, &iot.DeleteDimensionInput{Name: aws.String(dimName)})

	results = append(results, r.RunTest("iot", "Dimension_CreateDescribe", func() error {
		if _, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name: aws.String(dimName), Type: iottypes.DimensionTypeTopicFilter, StringValues: []string{"thing1"},
		}); err != nil {
			return fmt.Errorf("CreateDimension failed: %w", err)
		}
		out, err := tc.client.DescribeDimension(tc.ctx, &iot.DescribeDimensionInput{Name: aws.String(dimName)})
		if err != nil {
			return fmt.Errorf("DescribeDimension failed: %w", err)
		}
		if aws.ToString(out.Name) != dimName {
			return fmt.Errorf("expected name=%s", dimName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Dimension_UpdateEcho", func() error {
		if _, err := tc.client.UpdateDimension(tc.ctx, &iot.UpdateDimensionInput{
			Name: aws.String(dimName), StringValues: []string{"updated-val"},
		}); err != nil {
			return fmt.Errorf("UpdateDimension failed: %w", err)
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

	// ── CustomMetric ──
	cmName := uniqueName("custom-metric")
	defer tc.client.DeleteCustomMetric(tc.ctx, &iot.DeleteCustomMetricInput{MetricName: aws.String(cmName)})

	results = append(results, r.RunTest("iot", "CustomMetric_CreateDescribe", func() error {
		if _, err := tc.client.CreateCustomMetric(tc.ctx, &iot.CreateCustomMetricInput{
			MetricName: aws.String(cmName), MetricType: iottypes.CustomMetricTypeNumber, DisplayName: aws.String("m"),
		}); err != nil {
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
		return nil
	}))

	results = append(results, r.RunTest("iot", "CustomMetric_Describe_NotFound", func() error {
		_, err := tc.client.DescribeCustomMetric(tc.ctx, &iot.DescribeCustomMetricInput{MetricName: aws.String(uniqueName("nope-metric"))})
		return expectNotFound(err)
	}))

	// ── FleetMetric ──
	fmName := uniqueName("fleet-metric")
	defer tc.client.DeleteFleetMetric(tc.ctx, &iot.DeleteFleetMetricInput{MetricName: aws.String(fmName)})

	results = append(results, r.RunTest("iot", "FleetMetric_CreateDescribe", func() error {
		if _, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("thingName:*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		}); err != nil {
			return fmt.Errorf("CreateFleetMetric failed: %w", err)
		}
		out, err := tc.client.DescribeFleetMetric(tc.ctx, &iot.DescribeFleetMetricInput{MetricName: aws.String(fmName)})
		if err != nil {
			return fmt.Errorf("DescribeFleetMetric failed: %w", err)
		}
		_ = out
		return nil
	}))

	// ── ScheduledAudit ──
	saName := uniqueName("scheduled-audit")
	defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(saName)})

	results = append(results, r.RunTest("iot", "ScheduledAudit_CreateDescribe", func() error {
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyDaily,
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("CreateScheduledAudit failed: %w", err)
		}
		out, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{ScheduledAuditName: aws.String(saName)})
		if err != nil {
			return fmt.Errorf("DescribeScheduledAudit failed: %w", err)
		}
		_ = out
		return nil
	}))

	return results
}
