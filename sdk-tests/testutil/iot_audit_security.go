package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTAuditSecurityTests covers the audit/security-family list filters and
// member contracts: audit-suppression filtering/ordering plus the
// clientRequestToken idempotency replay, the audit-task time-range and type
// filters, the audit-findings taskId/time-range exclusivity, the
// scheduled-audit list day members, the security-profile dimension/metric
// filters, the violation-event required time range, the
// DeleteAccountAuditConfiguration scheduled-audits flag and the V2 logging
// eventConfigurations round trip.
func (r *TestRunner) runIoTAuditSecurityTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	resourceID := &iottypes.ResourceIdentifier{
		DeviceCertificateArn: aws.String(tc.arn("iot", "cert", "1111111111111111111111111111111111111111111111111111111111111111")),
	}

	// ---- Audit suppression: token idempotency, filters, ordering ----

	supCheck := uniqueName("sup-filter")
	supToken := uniqueName("sup-token")
	results = append(results, r.RunTest("iot", "AuditSuppression_TokenIdempotentReplay", func() error {
		defer tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName: aws.String(supCheck), ResourceIdentifier: resourceID,
		})
		in := &iot.CreateAuditSuppressionInput{
			CheckName:          aws.String(supCheck),
			ResourceIdentifier: resourceID,
			ExpirationDate:     aws.Time(time.Now().Add(24 * time.Hour)),
			ClientRequestToken: aws.String(supToken),
		}
		if _, err := tc.client.CreateAuditSuppression(tc.ctx, in); err != nil {
			return fmt.Errorf("initial create failed: %w", err)
		}
		// Replaying the same token for the same tuple is idempotent.
		if _, err := tc.client.CreateAuditSuppression(tc.ctx, in); err != nil {
			return fmt.Errorf("same-token replay must succeed: %w", err)
		}
		// A different token for the same tuple is the duplicate conflict.
		_, err := tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:          aws.String(supCheck),
			ResourceIdentifier: resourceID,
			ExpirationDate:     aws.Time(time.Now().Add(24 * time.Hour)),
			ClientRequestToken: aws.String(uniqueName("sup-token")),
		})
		if err := expectAWSErrorCode(err, "ResourceAlreadyExistsException"); err != nil {
			return fmt.Errorf("different token for the same tuple: %w", err)
		}
		// Each suppression must have a unique client request token, so
		// reusing the existing token for a different tuple is rejected.
		otherCheck := uniqueName("sup-other")
		defer tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName: aws.String(otherCheck), ResourceIdentifier: resourceID,
		})
		_, err = tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:          aws.String(otherCheck),
			ResourceIdentifier: resourceID,
			ExpirationDate:     aws.Time(time.Now().Add(24 * time.Hour)),
			ClientRequestToken: aws.String(supToken),
		})
		return expectAWSErrorCode(err, "ResourceAlreadyExistsException")
	}))

	supEarly := uniqueName("sup-early")
	supLate := uniqueName("sup-late")
	results = append(results, r.RunTest("iot", "ListAuditSuppressions_FilterAndOrder", func() error {
		defer tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName: aws.String(supEarly), ResourceIdentifier: resourceID,
		})
		defer tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName: aws.String(supLate), ResourceIdentifier: resourceID,
		})
		for _, entry := range []struct {
			check string
			exp   time.Time
		}{
			{supEarly, time.Now().Add(24 * time.Hour)},
			{supLate, time.Now().Add(48 * time.Hour)},
		} {
			if _, err := tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
				CheckName:          aws.String(entry.check),
				ResourceIdentifier: resourceID,
				ExpirationDate:     aws.Time(entry.exp),
				ClientRequestToken: aws.String(uniqueName("sup-token")),
			}); err != nil {
				return fmt.Errorf("create %s failed: %w", entry.check, err)
			}
		}
		// The checkName filter returns only the matching entry.
		filtered, err := paginate(func(next *string) ([]iottypes.AuditSuppression, *string, error) {
			out, err := tc.client.ListAuditSuppressions(tc.ctx, &iot.ListAuditSuppressionsInput{
				CheckName: aws.String(supEarly), NextToken: next,
			})
			if err != nil {
				return nil, nil, err
			}
			return out.Suppressions, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		if len(filtered) != 1 || aws.ToString(filtered[0].CheckName) != supEarly {
			return fmt.Errorf("checkName filter: got %d entries, want exactly %s", len(filtered), supEarly)
		}
		// The default order is ascending by expiration date.
		both, err := paginate(func(next *string) ([]iottypes.AuditSuppression, *string, error) {
			out, err := tc.client.ListAuditSuppressions(tc.ctx, &iot.ListAuditSuppressionsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.Suppressions, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		posEarly, posLate := -1, -1
		for i, s := range both {
			switch aws.ToString(s.CheckName) {
			case supEarly:
				posEarly = i
			case supLate:
				posLate = i
			}
		}
		if posEarly < 0 || posLate < 0 {
			return fmt.Errorf("expected both suppressions in the unfiltered list, got positions %d/%d", posEarly, posLate)
		}
		if posEarly > posLate {
			return fmt.Errorf("expected ascending expiration order (%s before %s)", supEarly, supLate)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "UpdateAuditSuppression_ExpirationExclusivity", func() error {
		defer tc.client.DeleteAuditSuppression(tc.ctx, &iot.DeleteAuditSuppressionInput{
			CheckName: aws.String(supCheck), ResourceIdentifier: resourceID,
		})
		// An explicit false alongside an expiration date is the pairing
		// the documented CLI example sends on create and update.
		if _, err := tc.client.CreateAuditSuppression(tc.ctx, &iot.CreateAuditSuppressionInput{
			CheckName:            aws.String(supCheck),
			ResourceIdentifier:   resourceID,
			ExpirationDate:       aws.Time(time.Now().Add(24 * time.Hour)),
			SuppressIndefinitely: aws.Bool(false),
			ClientRequestToken:   aws.String(uniqueName("sup-token")),
		}); err != nil {
			return fmt.Errorf("create with false + expirationDate failed: %w", err)
		}
		if _, err := tc.client.UpdateAuditSuppression(tc.ctx, &iot.UpdateAuditSuppressionInput{
			CheckName:            aws.String(supCheck),
			ResourceIdentifier:   resourceID,
			ExpirationDate:       aws.Time(time.Now().Add(72 * time.Hour)),
			SuppressIndefinitely: aws.Bool(false),
		}); err != nil {
			return fmt.Errorf("update with false + expirationDate failed: %w", err)
		}
		// An indefinite suppression together with an expiration date is
		// contradictory and rejected on update as on create.
		_, err := tc.client.UpdateAuditSuppression(tc.ctx, &iot.UpdateAuditSuppressionInput{
			CheckName:            aws.String(supCheck),
			ResourceIdentifier:   resourceID,
			ExpirationDate:       aws.Time(time.Now().Add(48 * time.Hour)),
			SuppressIndefinitely: aws.Bool(true),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// ---- Audit tasks: required time range and filters ----

	results = append(results, r.RunTest("iot", "ListAuditTasks_MissingTimeRangeRejected", func() error {
		// The typed SDK validates the required time-range members
		// client-side, so a half-open range never reaches the server; the
		// server-side rejection of an absent range is unit-pinned.
		_, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			StartTime: aws.Time(time.Now().Add(-time.Hour)),
		})
		if err == nil {
			return fmt.Errorf("startTime-only request must be rejected")
		}
		_, err = tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			EndTime: aws.Time(time.Now().Add(time.Hour)),
		})
		if err == nil {
			return fmt.Errorf("endTime-only request must be rejected")
		}
		return nil
	}))

	auditTaskId := ""
	results = append(results, r.RunTest("iot", "ListAuditTasks_TimeAndTypeFilters", func() error {
		windowStart := time.Now().Add(-time.Minute)
		out, err := tc.client.StartOnDemandAuditTask(tc.ctx, &iot.StartOnDemandAuditTaskInput{
			TargetCheckNames: []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		})
		if err != nil {
			return fmt.Errorf("StartOnDemandAuditTask failed: %w", err)
		}
		auditTaskId = aws.ToString(out.TaskId)
		windowEnd := time.Now().Add(time.Minute)

		inWindow, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			StartTime: &windowStart,
			EndTime:   &windowEnd,
			TaskType:  iottypes.AuditTaskTypeOnDemandAuditTask,
		})
		if err != nil {
			return err
		}
		if !containsAuditTask(inWindow.Tasks, auditTaskId) {
			return fmt.Errorf("task %s missing from its own time window", auditTaskId)
		}
		futureStart := time.Now().Add(time.Hour)
		futureEnd := time.Now().Add(2 * time.Hour)
		outOfWindow, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			StartTime: &futureStart, EndTime: &futureEnd,
		})
		if err != nil {
			return err
		}
		if containsAuditTask(outOfWindow.Tasks, auditTaskId) {
			return fmt.Errorf("task %s leaked into a disjoint time window", auditTaskId)
		}
		scheduledOnly, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			StartTime: &windowStart,
			EndTime:   &windowEnd,
			TaskType:  iottypes.AuditTaskTypeScheduledAuditTask,
		})
		if err != nil {
			return err
		}
		if containsAuditTask(scheduledOnly.Tasks, auditTaskId) {
			return fmt.Errorf("on-demand task %s matched the SCHEDULED_AUDIT_TASK filter", auditTaskId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "ListAuditFindings_TaskIdWithTimeRangeRejected", func() error {
		_, err := tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{
			TaskId:    aws.String(auditTaskId),
			StartTime: aws.Time(time.Now().Add(-time.Hour)),
			EndTime:   aws.Time(time.Now().Add(time.Hour)),
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("taskId plus time range: %w", err)
		}
		// Either the taskId or the startTime and endTime pair must be
		// specified, so a request carrying neither is rejected — the
		// checkName filter is not a substitute for either option.
		_, err = tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{
			CheckName: aws.String("DEVICE_CERTIFICATE_EXPIRING_CHECK"),
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("checkName without taskId or time range: %w", err)
		}
		// Half of the pair does not satisfy the requirement either.
		_, err = tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{
			StartTime: aws.Time(time.Now().Add(-time.Hour)),
		})
		if err := expectAWSErrorCode(err, "InvalidRequestException"); err != nil {
			return fmt.Errorf("startTime without endTime: %w", err)
		}
		// The full time-range pair on its own is accepted and returns the
		// (empty) findings page.
		_, err = tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{
			StartTime: aws.Time(time.Now().Add(-time.Hour)),
			EndTime:   aws.Time(time.Now().Add(time.Hour)),
		})
		return err
	}))

	// ---- Scheduled audits: list day members ----

	results = append(results, r.RunTest("iot", "ListScheduledAudits_DayMembersEcho", func() error {
		weeklyName := uniqueName("sa-weekly")
		defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(weeklyName)})
		monthlyName := uniqueName("sa-monthly")
		defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(monthlyName)})
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(weeklyName),
			Frequency:          iottypes.AuditFrequencyWeekly,
			DayOfWeek:          iottypes.DayOfWeekWed,
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("create weekly failed: %w", err)
		}
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(monthlyName),
			Frequency:          iottypes.AuditFrequencyMonthly,
			DayOfMonth:         aws.String("15"),
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("create monthly failed: %w", err)
		}
		audits, err := paginate(func(next *string) ([]iottypes.ScheduledAuditMetadata, *string, error) {
			out, err := tc.client.ListScheduledAudits(tc.ctx, &iot.ListScheduledAuditsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ScheduledAudits, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		var weekly, monthly bool
		for _, a := range audits {
			switch aws.ToString(a.ScheduledAuditName) {
			case weeklyName:
				weekly = a.DayOfWeek == iottypes.DayOfWeekWed
			case monthlyName:
				monthly = aws.ToString(a.DayOfMonth) == "15"
			}
		}
		if !weekly {
			return fmt.Errorf("weekly audit %s missing dayOfWeek=WED in list output", weeklyName)
		}
		if !monthly {
			return fmt.Errorf("monthly audit %s missing dayOfMonth=15 in list output", monthlyName)
		}
		return nil
	}))

	// ---- Security profiles: dimension/metric filters ----

	profileMetric := uniqueName("sp-metric")
	profileDim := uniqueName("sp-dimension")
	results = append(results, r.RunTest("iot", "ListSecurityProfiles_DimensionMetricFilters", func() error {
		cleanupMetric, err := tc.createSecurityProfile(profileMetric, []iottypes.Behavior{{
			Name:     aws.String("count-behavior"),
			Metric:   aws.String("aws:num-connected-devices"),
			Criteria: &iottypes.BehaviorCriteria{ComparisonOperator: iottypes.ComparisonOperatorLessThan, Value: &iottypes.MetricValue{Count: aws.Int64(10)}},
		}})
		if err != nil {
			return fmt.Errorf("create metric profile failed: %w", err)
		}
		defer cleanupMetric()
		cleanupDim, err := tc.createSecurityProfile(profileDim, []iottypes.Behavior{{
			Name:            aws.String("dim-behavior"),
			MetricDimension: &iottypes.MetricDimension{DimensionName: aws.String("sdk-test-dimension"), Operator: iottypes.DimensionValueOperatorNotIn},
			Criteria:        &iottypes.BehaviorCriteria{ComparisonOperator: iottypes.ComparisonOperatorLessThan, Value: &iottypes.MetricValue{Count: aws.Int64(10)}},
		}})
		if err != nil {
			return fmt.Errorf("create dimension profile failed: %w", err)
		}
		defer cleanupDim()
		// The metricDimension operator round-trips with the dimension
		// name through describe.
		described, err := tc.client.DescribeSecurityProfile(tc.ctx, &iot.DescribeSecurityProfileInput{
			SecurityProfileName: aws.String(profileDim),
		})
		if err != nil {
			return fmt.Errorf("describe dimension profile failed: %w", err)
		}
		if len(described.Behaviors) != 1 ||
			described.Behaviors[0].MetricDimension == nil {
			return fmt.Errorf("expected one behavior with a metric dimension, got %+v", described.Behaviors)
		}
		if described.Behaviors[0].MetricDimension.Operator != iottypes.DimensionValueOperatorNotIn {
			return fmt.Errorf("expected NOT_IN operator round-trip, got %v", described.Behaviors[0].MetricDimension.Operator)
		}
		byDimension, err := paginate(func(next *string) ([]iottypes.SecurityProfileIdentifier, *string, error) {
			out, err := tc.client.ListSecurityProfiles(tc.ctx, &iot.ListSecurityProfilesInput{
				DimensionName: aws.String("sdk-test-dimension"), NextToken: next,
			})
			if err != nil {
				return nil, nil, err
			}
			return out.SecurityProfileIdentifiers, out.NextToken, nil
		})
		if err != nil {
			return fmt.Errorf("dimension filter failed: %w", err)
		}
		if !containsProfile(byDimension, profileDim) {
			return fmt.Errorf("dimension filter missing profile %s", profileDim)
		}
		if containsProfile(byDimension, profileMetric) {
			return fmt.Errorf("dimension filter wrongly matched profile %s", profileMetric)
		}
		byMetric, err := paginate(func(next *string) ([]iottypes.SecurityProfileIdentifier, *string, error) {
			out, err := tc.client.ListSecurityProfiles(tc.ctx, &iot.ListSecurityProfilesInput{
				MetricName: aws.String("aws:num-connected-devices"), NextToken: next,
			})
			if err != nil {
				return nil, nil, err
			}
			return out.SecurityProfileIdentifiers, out.NextToken, nil
		})
		if err != nil {
			return fmt.Errorf("metric filter failed: %w", err)
		}
		if !containsProfile(byMetric, profileMetric) {
			return fmt.Errorf("metric filter missing profile %s", profileMetric)
		}
		if containsProfile(byMetric, profileDim) {
			return fmt.Errorf("metric filter wrongly matched profile %s", profileDim)
		}
		// The two filters are mutually exclusive.
		_, err = tc.client.ListSecurityProfiles(tc.ctx, &iot.ListSecurityProfilesInput{
			DimensionName: aws.String("sdk-test-dimension"),
			MetricName:    aws.String("aws:num-connected-devices"),
		})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	// ---- Violations: required time range and filters ----

	results = append(results, r.RunTest("iot", "ListViolationEvents_MissingTimeRangeRejected", func() error {
		// The typed SDK validates the required time-range members
		// client-side; the server-side rejection is unit-pinned.
		_, err := tc.client.ListViolationEvents(tc.ctx, &iot.ListViolationEventsInput{})
		if err == nil {
			return fmt.Errorf("rangeless request must be rejected")
		}
		start := time.Now().Add(-time.Hour)
		end := time.Now().Add(time.Hour)
		_, err = tc.client.ListViolationEvents(tc.ctx, &iot.ListViolationEventsInput{
			StartTime:            &start,
			EndTime:              &end,
			VerificationState:    iottypes.VerificationStateTruePositive,
			ListSuppressedAlerts: aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "ListActiveViolations_FiltersAccepted", func() error {
		_, err := tc.client.ListActiveViolations(tc.ctx, &iot.ListActiveViolationsInput{
			ThingName:            aws.String(uniqueName("no-thing")),
			VerificationState:    iottypes.VerificationStateTruePositive,
			ListSuppressedAlerts: aws.Bool(true),
			BehaviorCriteriaType: iottypes.BehaviorCriteriaTypeStatic,
		})
		return err
	}))

	// ---- Account audit configuration: scheduled-audits delete flag ----

	results = append(results, r.RunTest("iot", "DeleteAccountAuditConfiguration_ScheduledAuditsFlag", func() error {
		doomed := uniqueName("sa-doomed")
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(doomed),
			Frequency:          iottypes.AuditFrequencyDaily,
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("create doomed audit failed: %w", err)
		}
		if _, err := tc.client.DeleteAccountAuditConfiguration(tc.ctx, &iot.DeleteAccountAuditConfigurationInput{
			DeleteScheduledAudits: true,
		}); err != nil {
			return fmt.Errorf("delete with flag failed: %w", err)
		}
		if _, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{
			ScheduledAuditName: aws.String(doomed),
		}); err == nil {
			return fmt.Errorf("scheduled audit %s survived DeleteScheduledAudits=true", doomed)
		}
		survivor := uniqueName("sa-survivor")
		defer tc.client.DeleteScheduledAudit(tc.ctx, &iot.DeleteScheduledAuditInput{ScheduledAuditName: aws.String(survivor)})
		if _, err := tc.client.CreateScheduledAudit(tc.ctx, &iot.CreateScheduledAuditInput{
			ScheduledAuditName: aws.String(survivor),
			Frequency:          iottypes.AuditFrequencyDaily,
			TargetCheckNames:   []string{"DEVICE_CERTIFICATE_EXPIRING_CHECK"},
		}); err != nil {
			return fmt.Errorf("create survivor audit failed: %w", err)
		}
		if _, err := tc.client.DeleteAccountAuditConfiguration(tc.ctx, &iot.DeleteAccountAuditConfigurationInput{}); err != nil {
			return fmt.Errorf("delete without flag failed: %w", err)
		}
		if _, err := tc.client.DescribeScheduledAudit(tc.ctx, &iot.DescribeScheduledAuditInput{
			ScheduledAuditName: aws.String(survivor),
		}); err != nil {
			return fmt.Errorf("scheduled audit %s deleted without the flag: %w", survivor, err)
		}
		return nil
	}))

	// ---- V2 logging: eventConfigurations round trip ----

	results = append(results, r.RunTest("iot", "V2Logging_EventConfigurationsRoundTrip", func() error {
		_, err := tc.client.SetV2LoggingOptions(tc.ctx, &iot.SetV2LoggingOptionsInput{
			DefaultLogLevel: iottypes.LogLevelError,
			EventConfigurations: []iottypes.LogEventConfiguration{
				{EventType: aws.String("Connect"), LogLevel: iottypes.LogLevelWarn},
				{EventType: aws.String("Publish"), LogLevel: iottypes.LogLevelError, LogDestination: aws.String("/aws/thing/sdk-test")},
			},
		})
		if err != nil {
			return fmt.Errorf("SetV2LoggingOptions failed: %w", err)
		}
		out, err := tc.client.GetV2LoggingOptions(tc.ctx, &iot.GetV2LoggingOptionsInput{})
		if err != nil {
			return err
		}
		if out.DefaultLogLevel != iottypes.LogLevelError {
			return fmt.Errorf("expected defaultLogLevel=ERROR, got %s", out.DefaultLogLevel)
		}
		var connect, publish bool
		for _, cfg := range out.EventConfigurations {
			switch aws.ToString(cfg.EventType) {
			case "Connect":
				connect = cfg.LogLevel == iottypes.LogLevelWarn
			case "Publish":
				publish = cfg.LogLevel == iottypes.LogLevelError && aws.ToString(cfg.LogDestination) == "/aws/thing/sdk-test"
			}
		}
		if !connect || !publish {
			return fmt.Errorf("eventConfigurations round trip incomplete: connect=%v publish=%v", connect, publish)
		}
		return nil
	}))

	// ---- Custom metrics / dimensions: token idempotency ----

	metricName := uniqueName("cm-idem")
	results = append(results, r.RunTest("iot", "CustomMetric_TokenIdempotentReplay", func() error {
		defer tc.client.DeleteCustomMetric(tc.ctx, &iot.DeleteCustomMetricInput{MetricName: aws.String(metricName)})
		in := &iot.CreateCustomMetricInput{
			MetricName:         aws.String(metricName),
			MetricType:         iottypes.CustomMetricTypeStringList,
			DisplayName:        aws.String("sdk idempotency metric"),
			ClientRequestToken: aws.String(uniqueName("cm-token")),
		}
		if _, err := tc.client.CreateCustomMetric(tc.ctx, in); err != nil {
			return fmt.Errorf("initial create failed: %w", err)
		}
		if _, err := tc.client.CreateCustomMetric(tc.ctx, in); err != nil {
			return fmt.Errorf("same-token replay must succeed: %w", err)
		}
		_, err := tc.client.CreateCustomMetric(tc.ctx, &iot.CreateCustomMetricInput{
			MetricName:         aws.String(metricName),
			MetricType:         iottypes.CustomMetricTypeStringList,
			DisplayName:        aws.String("sdk idempotency metric"),
			ClientRequestToken: aws.String(uniqueName("cm-token")),
		})
		return expectAWSErrorCode(err, "ResourceAlreadyExistsException")
	}))

	dimName := uniqueName("dim-idem")
	results = append(results, r.RunTest("iot", "Dimension_TokenIdempotentReplay", func() error {
		defer tc.client.DeleteDimension(tc.ctx, &iot.DeleteDimensionInput{Name: aws.String(dimName)})
		dimToken := uniqueName("dim-token")
		in := &iot.CreateDimensionInput{
			Name:               aws.String(dimName),
			Type:               iottypes.DimensionTypeTopicFilter,
			StringValues:       []string{"sdk/test/#"},
			ClientRequestToken: aws.String(dimToken),
		}
		if _, err := tc.client.CreateDimension(tc.ctx, in); err != nil {
			return fmt.Errorf("initial create failed: %w", err)
		}
		if _, err := tc.client.CreateDimension(tc.ctx, in); err != nil {
			return fmt.Errorf("same-token replay must succeed: %w", err)
		}
		_, err := tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name:               aws.String(dimName),
			Type:               iottypes.DimensionTypeTopicFilter,
			StringValues:       []string{"sdk/test/#"},
			ClientRequestToken: aws.String(uniqueName("dim-token")),
		})
		if err := expectAWSErrorCode(err, "ResourceAlreadyExistsException"); err != nil {
			return fmt.Errorf("different token for the same name: %w", err)
		}
		// Each dimension must have a unique client request token, so
		// reusing the existing token for a different dimension is
		// rejected.
		otherDim := uniqueName("dim-other")
		defer tc.client.DeleteDimension(tc.ctx, &iot.DeleteDimensionInput{Name: aws.String(otherDim)})
		_, err = tc.client.CreateDimension(tc.ctx, &iot.CreateDimensionInput{
			Name:               aws.String(otherDim),
			Type:               iottypes.DimensionTypeTopicFilter,
			StringValues:       []string{"sdk/test/#"},
			ClientRequestToken: aws.String(dimToken),
		})
		return expectAWSErrorCode(err, "ResourceAlreadyExistsException")
	}))

	return results
}

func containsAuditTask(tasks []iottypes.AuditTaskMetadata, taskId string) bool {
	for _, t := range tasks {
		if aws.ToString(t.TaskId) == taskId {
			return true
		}
	}
	return false
}

func containsProfile(ids []iottypes.SecurityProfileIdentifier, name string) bool {
	for _, id := range ids {
		if aws.ToString(id.Name) == name {
			return true
		}
	}
	return false
}
