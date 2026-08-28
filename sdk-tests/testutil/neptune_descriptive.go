package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

func (r *TestRunner) runNeptuneDescriptiveTests(tc *neptuneContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptune", "DescribeEventCategories", func() error {
		resp, err := tc.client.DescribeEventCategories(tc.ctx, &neptune.DescribeEventCategoriesInput{})
		if err != nil {
			return err
		}
		if len(resp.EventCategoriesMapList) == 0 {
			return fmt.Errorf("expected at least one event category map entry")
		}
		found := false
		for _, m := range resp.EventCategoriesMapList {
			if m.SourceType != nil && *m.SourceType == "db-instance" {
				found = true
				if len(m.EventCategories) == 0 {
					return fmt.Errorf("expected non-empty EventCategories for db-instance")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("expected db-instance event category in response")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeEvents", func() error {
		resp, err := tc.client.DescribeEvents(tc.ctx, &neptune.DescribeEventsInput{})
		if err != nil {
			return err
		}
		if resp.Events == nil {
			return fmt.Errorf("Events field is nil")
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("expected at least one event after resource creation")
		}
		ev := resp.Events[0]
		if ev.SourceType == "" {
			return fmt.Errorf("expected non-empty SourceType in first event")
		}

		// A client-supplied MaxRecords must bound the page (documented range
		// 20-100); it must not fall back to the 100 default.
		paged, err := tc.client.DescribeEvents(tc.ctx, &neptune.DescribeEventsInput{
			MaxRecords: aws.Int32(20),
		})
		if err != nil {
			return err
		}
		if len(paged.Events) > 20 {
			return fmt.Errorf("expected at most 20 events with MaxRecords=20, got %d", len(paged.Events))
		}

		// A MaxRecords below the documented minimum is clamped up to 20 and
		// accepted, not rejected.
		clamped, err := tc.client.DescribeEvents(tc.ctx, &neptune.DescribeEventsInput{
			MaxRecords: aws.Int32(5),
		})
		if err != nil {
			return err
		}
		if len(clamped.Events) > 20 {
			return fmt.Errorf("expected at most 20 events with MaxRecords=5 clamped to the minimum, got %d", len(clamped.Events))
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribePendingMaintenanceActions", func() error {
		resp, err := tc.client.DescribePendingMaintenanceActions(tc.ctx, &neptune.DescribePendingMaintenanceActionsInput{})
		if err != nil {
			return err
		}
		if resp.PendingMaintenanceActions == nil {
			return fmt.Errorf("PendingMaintenanceActions field is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeOrderableDBInstanceOptions", func() error {
		resp, err := tc.client.DescribeOrderableDBInstanceOptions(tc.ctx, &neptune.DescribeOrderableDBInstanceOptionsInput{
			Engine: aws.String("neptune"),
		})
		if err != nil {
			return err
		}
		if len(resp.OrderableDBInstanceOptions) == 0 {
			return fmt.Errorf("expected at least one orderable DB instance option")
		}
		opt := resp.OrderableDBInstanceOptions[0]
		if opt.Engine == nil || *opt.Engine != "neptune" {
			return fmt.Errorf("expected engine=neptune in orderable options, got %v", opt.Engine)
		}
		if opt.DBInstanceClass == nil || *opt.DBInstanceClass == "" {
			return fmt.Errorf("expected non-empty DBInstanceClass in orderable options")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "DescribeEngineDefaultParameters", func() error {
		resp, err := tc.client.DescribeEngineDefaultParameters(tc.ctx, &neptune.DescribeEngineDefaultParametersInput{
			DBParameterGroupFamily: aws.String("neptune1"),
		})
		if err != nil {
			return err
		}
		if resp.EngineDefaults.Parameters == nil || len(resp.EngineDefaults.Parameters) == 0 {
			return fmt.Errorf("expected non-empty default parameters list")
		}
		foundModifiable := false
		for _, p := range resp.EngineDefaults.Parameters {
			if p.ParameterName != nil && *p.ParameterName != "" {
				foundModifiable = true
				break
			}
		}
		if !foundModifiable {
			return fmt.Errorf("expected at least one parameter with non-empty ParameterName")
		}
		return nil
	}))

	return results
}
