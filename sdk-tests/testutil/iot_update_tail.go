package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

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
	defer tc.client.DeleteThing(tc.ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})
	defer tc.client.DeleteThingGroup(tc.ctx, &iot.DeleteThingGroupInput{ThingGroupName: aws.String(thingGroup)})

	// Setup: create all resources in one step so updates have real targets.
	results = append(results, r.RunTest("iot", "Update_Setup", func() error {
		if _, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name: aws.String(dimName), Type: iottypes.DimensionTypeTopicFilter, StringValues: []string{"x"},
		}); err != nil {
			return fmt.Errorf("CreateDimension failed: %w", err)
		}
		if _, err := tc.client.CreateMitigationAction(tc.ctx, &iot.CreateMitigationActionInput{
			ActionName: aws.String(mitName), RoleArn: aws.String(tc.iamRoleARN("test")), ActionParams: &iottypes.MitigationActionParams{},
		}); err != nil {
			return fmt.Errorf("CreateMitigationAction failed: %w", err)
		}
		if _, err := tc.client.CreateFleetMetric(tc.ctx, &iot.CreateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(60),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		}); err != nil {
			return fmt.Errorf("CreateFleetMetric failed: %w", err)
		}
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyDaily,
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("CreateScheduledAudit failed: %w", err)
		}
		if _, err := tc.client.CreateStream(tc.ctx, &iot.CreateStreamInput{
			StreamId: aws.String(streamID), Files: []iottypes.StreamFile{}, RoleArn: aws.String(tc.iamRoleARN("test")),
		}); err != nil {
			return fmt.Errorf("CreateStream failed: %w", err)
		}
		if _, err := tc.client.CreateThing(tc.ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)}); err != nil {
			return fmt.Errorf("CreateThing failed: %w", err)
		}
		if _, err := tc.client.CreateThingGroup(tc.ctx, &iot.CreateThingGroupInput{ThingGroupName: aws.String(thingGroup)}); err != nil {
			return fmt.Errorf("CreateThingGroup failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "UpdateDimension", func() error {
		_, err := tc.client.UpdateDimension(tc.ctx, &iot.UpdateDimensionInput{
			Name: aws.String(dimName), StringValues: []string{"updated"},
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "UpdateMitigationAction", func() error {
		_, err := tc.client.UpdateMitigationAction(tc.ctx, &iot.UpdateMitigationActionInput{
			ActionName: aws.String(mitName), RoleArn: aws.String(tc.iamRoleARN("test")), ActionParams: &iottypes.MitigationActionParams{},
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "UpdateFleetMetric", func() error {
		_, err := tc.client.UpdateFleetMetric(tc.ctx, &iot.UpdateFleetMetricInput{
			MetricName: aws.String(fmName), QueryString: aws.String("*"), IndexName: aws.String("AWS_Things"),
			AggregationField: aws.String("thingName"), Period: aws.Int32(120),
			AggregationType: &iottypes.AggregationType{Name: iottypes.AggregationTypeNameStatistics},
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "UpdateScheduledAudit", func() error {
		_, err := tc.client.UpdateScheduledAudit(tc.ctx, &iot.UpdateScheduledAuditInput{
			ScheduledAuditName: aws.String(saName), Frequency: iottypes.AuditFrequencyWeekly,
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		return err
	}))
	results = append(results, r.RunTest("iot", "UpdateStream", func() error {
		_, err := tc.client.UpdateStream(tc.ctx, &iot.UpdateStreamInput{StreamId: aws.String(streamID), Files: []iottypes.StreamFile{}})
		return err
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

	return results
}
