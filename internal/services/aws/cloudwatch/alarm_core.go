package cloudwatch

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// --- Transport-agnostic Input DTOs ---

// PutMetricAlarmInput holds all fields for PutMetricAlarm in a
// transport-agnostic format. Both the HTTP API handler and the gRPC-Web
// admin handler build this struct and delegate to putMetricAlarmCore.
type PutMetricAlarmInput struct {
	AlarmName               string
	Namespace               string
	MetricName              string
	Dimensions              []cwstore.Dimension
	ComparisonOperator      string
	Threshold               float64
	EvaluationPeriods       int32
	Period                  int32
	Statistic               string
	ExtendedStatistic       string
	TreatMissingData        string
	AlarmDescription        string
	ActionsEnabled          bool
	DatapointsToAlarm       int32
	AlarmActions            []string
	OKActions               []string
	InsufficientDataActions []string
	Unit                    string
	ThresholdMetricID       string
	Metrics                 []cwstore.MetricDataQuery
	EvaluationCriteria      *cwstore.EvaluationCriteria
	Tags                    map[string]string
}

// DeleteAlarmsInput holds the alarm names to delete.
type DeleteAlarmsInput struct {
	AlarmNames []string
}

// ListMetricsInput holds the filter parameters for ListMetrics.
type ListMetricsInput struct {
	Namespace      string
	MetricName     string
	Dimensions     []cwstore.Dimension
	NextToken      string
	MaxResults     int
	RecentlyActive string
}

// ListMetricsResult holds the paginated result of ListMetrics.
type ListMetricsResult struct {
	Metrics     []cwstore.MetricDatum
	NextToken   string
	IsTruncated bool
}

// DescribeAlarmsInput holds the filter parameters for DescribeAlarms.
type DescribeAlarmsInput struct {
	AlarmNamePrefix     string
	AlarmNames          []string
	StateValueFilter    string
	ActionPrefix        string
	AlarmTypes          []string
	ChildrenOfAlarmName string
	ParentsOfAlarmName  string
	NextToken           string
	MaxRecords          int
}

// --- Core methods ---

// putMetricAlarmCore validates the input, builds a store Alarm struct,
// applies AWS-specified defaults, and upserts the alarm. Both the HTTP
// API handler and the admin handler call this method so that validation
// has a single source of truth.
func (s *CloudWatchService) putMetricAlarmCore(stores *cloudwatchStores, input *PutMetricAlarmInput) (string, error) {
	if input.AlarmName == "" {
		return "", ErrInvalidParameter
	}
	if len(input.AlarmName) > maxAlarmNameLen {
		return "", awserrors.NewInvalidParameterValueException(
			"AlarmName must not exceed 255 characters")
	}
	if err := validateAlarmDescription(input.AlarmDescription); err != nil {
		return "", err
	}

	if input.ComparisonOperator == "" {
		input.ComparisonOperator = "GreaterThanOrEqualToThreshold"
	}
	if !validComparisonOperators[input.ComparisonOperator] {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid ComparisonOperator: %s", input.ComparisonOperator))
	}

	if input.EvaluationPeriods == 0 {
		input.EvaluationPeriods = 1
	}
	if input.Period == 0 {
		input.Period = 60
	}
	if err := validatePeriod(input.Period); err != nil {
		return "", err
	}
	if int64(input.Period)*int64(input.EvaluationPeriods) > maxEvaluationWindow {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Period * EvaluationPeriods must not exceed %d seconds (7 days)", maxEvaluationWindow))
	}

	if input.DatapointsToAlarm == 0 {
		input.DatapointsToAlarm = input.EvaluationPeriods
	}
	if input.DatapointsToAlarm > input.EvaluationPeriods {
		return "", awserrors.NewInvalidParameterValueException(
			"DatapointsToAlarm must be less than or equal to EvaluationPeriods")
	}
	if input.DatapointsToAlarm != 0 {
		if err := validateDatapointsToAlarm(input.DatapointsToAlarm); err != nil {
			return "", err
		}
	}

	if input.Statistic == "" && input.ExtendedStatistic == "" && len(input.Metrics) == 0 {
		input.Statistic = "Average"
	}
	if input.Statistic != "" && !validStatistics[input.Statistic] {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid Statistic: %s. Must be one of SampleCount, Average, Sum, Minimum, Maximum", input.Statistic))
	}
	if input.Statistic != "" && input.ExtendedStatistic != "" {
		return "", awserrors.NewInvalidParameterValueException(
			"Statistic and ExtendedStatistic are mutually exclusive")
	}
	if input.ExtendedStatistic != "" {
		if err := validateExtendedStatistic(input.ExtendedStatistic); err != nil {
			return "", err
		}
	}

	if input.TreatMissingData == "" {
		input.TreatMissingData = "missing"
	}
	if !validTreatMissingData[input.TreatMissingData] {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid TreatMissingData: %s. Must be one of breaching, notBreaching, ignore, missing", input.TreatMissingData))
	}

	if len(input.Dimensions) > maxDimensions {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of Dimensions must not exceed %d", maxDimensions))
	}

	if err := validateAlarmActions(input.AlarmActions, input.OKActions, input.InsufficientDataActions); err != nil {
		return "", err
	}

	if err := validateUnit(input.Unit); err != nil {
		return "", err
	}
	if err := validateThresholdMetricId(input.ThresholdMetricID); err != nil {
		return "", err
	}

	alarm := cwstore.NewAlarm(input.AlarmName, input.Namespace, input.MetricName)
	alarm.Dimensions = input.Dimensions
	alarm.ComparisonOperator = input.ComparisonOperator
	alarm.Threshold = input.Threshold
	alarm.EvaluationPeriods = input.EvaluationPeriods
	alarm.DatapointsToAlarm = input.DatapointsToAlarm
	alarm.Period = input.Period
	alarm.Statistic = input.Statistic
	alarm.ExtendedStatistic = input.ExtendedStatistic
	alarm.TreatMissingData = input.TreatMissingData
	alarm.AlarmDescription = input.AlarmDescription
	alarm.ActionsEnabled = input.ActionsEnabled
	alarm.AlarmActions = input.AlarmActions
	alarm.OKActions = input.OKActions
	alarm.InsufficientDataActions = input.InsufficientDataActions
	alarm.Unit = cwstore.StandardUnit(input.Unit)
	alarm.ThresholdMetricID = input.ThresholdMetricID
	alarm.Metrics = input.Metrics
	alarm.EvaluationCriteria = input.EvaluationCriteria
	alarm.Tags = input.Tags

	alarm, err := s.upsertAlarm(stores.alarms, alarm, cwstore.AlarmTypeMetricAlarm)
	if err != nil {
		return "", err
	}

	return alarm.ARN, nil
}

// deleteAlarmsCore validates alarm names, checks for composite alarm
// references, and deletes each alarm. Returns ResourceConflict when a
// composite alarm references one of the alarms being deleted.
func (s *CloudWatchService) deleteAlarmsCore(stores *cloudwatchStores, input *DeleteAlarmsInput) error {
	if len(input.AlarmNames) == 0 {
		return ErrInvalidParameter
	}

	allAlarms, err := stores.alarms.ListAlarms("")
	if err != nil {
		return err
	}
	deleteSet := make(map[string]bool, len(input.AlarmNames))
	for _, n := range input.AlarmNames {
		deleteSet[n] = true
	}
	for _, a := range allAlarms {
		if a.AlarmType != cwstore.AlarmTypeCompositeAlarm || a.AlarmRule == "" {
			continue
		}
		rule, err := parseAlarmRule(a.AlarmRule)
		if err != nil {
			continue
		}
		for _, childName := range rule.childAlarmNames() {
			if deleteSet[childName] {
				return awserrors.NewAWSError("ResourceConflict",
					fmt.Sprintf("Alarm %s is referenced by composite alarm %s", childName, a.Name),
					http.StatusConflict)
			}
		}
	}

	for _, name := range input.AlarmNames {
		if err := stores.alarms.DeleteAlarm(name); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return ErrAlarmNotFound
			}
			return err
		}
	}

	return nil
}

// listMetricsCore queries the metric store with the given filters and
// returns paginated results. Both the HTTP API handler and the admin
// handler call this method.
func (s *CloudWatchService) listMetricsCore(stores *cloudwatchStores, input *ListMetricsInput) (*ListMetricsResult, error) {
	if input.RecentlyActive != "" && input.RecentlyActive != "PT3H" {
		return nil, awserrors.NewInvalidParameterValueException(
			"RecentlyActive must be PT3H")
	}

	if len(input.Dimensions) > maxListMetricsDimensions {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of Dimensions must not exceed %d", maxListMetricsDimensions))
	}

	maxResults := 500
	if input.MaxResults > 0 {
		maxResults = input.MaxResults
	}

	metrics, nextMarker, isTruncated, err := stores.metrics.ListMetricsPaginated(
		input.Namespace, input.MetricName, input.Dimensions,
		input.NextToken, maxResults, input.RecentlyActive == "PT3H")
	if err != nil {
		return nil, err
	}

	return &ListMetricsResult{
		Metrics:     metrics,
		NextToken:   nextMarker,
		IsTruncated: isTruncated,
	}, nil
}

// describeAlarmsCore queries the alarm store with the given filters
// and returns matching alarms. The admin handler calls a simplified
// version (prefix-only) while the HTTP API handler passes full filters.
func (s *CloudWatchService) describeAlarmsCore(stores *cloudwatchStores, input *DescribeAlarmsInput) ([]*cwstore.Alarm, string, error) {
	if err := validateStateValueFilter(input.StateValueFilter); err != nil {
		return nil, "", err
	}
	if err := validateAlarmTypeFilters(input.AlarmTypes); err != nil {
		return nil, "", err
	}

	alarmNamePrefix := input.AlarmNamePrefix
	childrenOfAlarmName := input.ChildrenOfAlarmName
	parentsOfAlarmName := input.ParentsOfAlarmName
	stateValueFilter := input.StateValueFilter
	actionPrefix := input.ActionPrefix
	alarmNames := input.AlarmNames
	alarmTypes := input.AlarmTypes

	var typeFilter func(*cwstore.Alarm) bool
	if len(alarmTypes) > 0 {
		typeSet := make(map[string]bool)
		for _, t := range alarmTypes {
			typeSet[t] = true
		}
		typeFilter = func(a *cwstore.Alarm) bool {
			at := a.AlarmType
			if at == "" {
				at = cwstore.AlarmTypeMetricAlarm
			}
			return typeSet[at]
		}
	}
	if stateValueFilter != "" {
		svf := typeFilter
		typeFilter = func(a *cwstore.Alarm) bool {
			if svf != nil && !svf(a) {
				return false
			}
			return a.State == stateValueFilter
		}
	}
	if actionPrefix != "" {
		apf := typeFilter
		typeFilter = func(a *cwstore.Alarm) bool {
			if apf != nil && !apf(a) {
				return false
			}
			return alarmMatchesActionPrefix(a, actionPrefix)
		}
	}

	if childrenOfAlarmName != "" {
		parent, err := stores.alarms.GetAlarm(childrenOfAlarmName)
		if err != nil {
			return nil, "", nil
		}
		rule, err := parseAlarmRule(parent.AlarmRule)
		if err != nil {
			return nil, "", nil
		}
		childNames := rule.childAlarmNames()
		var children []*cwstore.Alarm
		for _, name := range childNames {
			if a, err := stores.alarms.GetAlarm(name); err == nil {
				if typeFilter == nil || typeFilter(a) {
					children = append(children, a)
				}
			}
		}
		return children, "", nil
	}

	if parentsOfAlarmName != "" {
		allAlarms, err := stores.alarms.ListAlarms("")
		if err != nil {
			return nil, "", err
		}
		var parents []*cwstore.Alarm
		for _, a := range allAlarms {
			if a.AlarmType != cwstore.AlarmTypeCompositeAlarm || a.AlarmRule == "" {
				continue
			}
			rule, err := parseAlarmRule(a.AlarmRule)
			if err != nil {
				continue
			}
			for _, childName := range rule.childAlarmNames() {
				if childName == parentsOfAlarmName {
					if typeFilter == nil || typeFilter(a) {
						parents = append(parents, a)
					}
					break
				}
			}
		}
		return parents, "", nil
	}

	if len(alarmNames) > 0 {
		var alarms []*cwstore.Alarm
		for _, name := range alarmNames {
			alarm, err := stores.alarms.GetAlarm(name)
			if err == nil {
				if typeFilter == nil || typeFilter(alarm) {
					alarms = append(alarms, alarm)
				}
			}
		}
		return alarms, "", nil
	}

	maxRecords := capMaxRecords(input.MaxRecords, 100, maxDescribeAlarmsRecords)
	opts := common.ListOptions{Marker: input.NextToken, MaxItems: maxRecords}
	result, err := stores.alarms.ListAlarmsPaginated(alarmNamePrefix, opts, typeFilter)
	if err != nil {
		return nil, "", err
	}

	return result.Items, result.NextMarker, nil
}

// DescribeAlarmsForMetricInput holds parameters for DescribeAlarmsForMetric.
type DescribeAlarmsForMetricInput struct {
	Namespace         string
	MetricName        string
	Dimensions        []cwstore.Dimension
	Statistic         string
	ExtendedStatistic string
	Period            int32
	Unit              string
	NextToken         string
	MaxRecords        int
}

// describeAlarmsForMetricCore validates input and retrieves alarms for
// a specific metric.
func (s *CloudWatchService) describeAlarmsForMetricCore(stores *cloudwatchStores, input *DescribeAlarmsForMetricInput) ([]*cwstore.Alarm, string, error) {
	if err := validateNamespace(input.Namespace); err != nil {
		return nil, "", err
	}
	if err := validateMetricName(input.MetricName); err != nil {
		return nil, "", err
	}
	if input.Statistic != "" && !validStatistics[input.Statistic] {
		return nil, "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid Statistic: %s", input.Statistic))
	}
	if input.ExtendedStatistic != "" {
		if err := validateExtendedStatistic(input.ExtendedStatistic); err != nil {
			return nil, "", err
		}
	}
	if err := validateUnit(input.Unit); err != nil {
		return nil, "", err
	}

	maxRecords := capMaxRecords(input.MaxRecords, 100, maxDescribeAlarmsRecords)

	filter := func(a *cwstore.Alarm) bool {
		if a.Namespace != input.Namespace || a.MetricName != input.MetricName {
			return false
		}
		if len(input.Dimensions) > 0 && len(a.Dimensions) > 0 {
			if !dimensionsMatch(a.Dimensions, input.Dimensions) {
				return false
			}
		}
		if input.Statistic != "" && a.Statistic != input.Statistic {
			return false
		}
		if input.ExtendedStatistic != "" && a.ExtendedStatistic != input.ExtendedStatistic {
			return false
		}
		if input.Period > 0 && a.Period != input.Period {
			return false
		}
		if input.Unit != "" && string(a.Unit) != input.Unit {
			return false
		}
		return true
	}

	opts := common.ListOptions{Marker: input.NextToken, MaxItems: maxRecords}
	result, err := stores.alarms.ListAlarmsPaginated("", opts, filter)
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// PutCompositeAlarmInput holds parameters for PutCompositeAlarm.
type PutCompositeAlarmInput struct {
	AlarmName                   string
	AlarmRule                   string
	AlarmDescription            string
	ActionsEnabled              bool
	AlarmActions                []string
	OKActions                   []string
	InsufficientDataActions     []string
	ActionsSuppressor           string
	ActionsSuppressorWaitPeriod int32
	ActionsSuppressorExtPeriod  int32
	Tags                        map[string]string
}

// putCompositeAlarmCore validates input and creates or updates a
// composite alarm.
func (s *CloudWatchService) putCompositeAlarmCore(stores *cloudwatchStores, input *PutCompositeAlarmInput) (string, error) {
	if input.AlarmName == "" {
		return "", ErrInvalidParameter
	}
	if len(input.AlarmName) > maxAlarmNameLen {
		return "", awserrors.NewInvalidParameterValueException(
			"AlarmName must not exceed 255 characters")
	}
	if err := validateAlarmDescription(input.AlarmDescription); err != nil {
		return "", err
	}
	if input.AlarmRule == "" {
		return "", ErrInvalidParameter
	}
	rule, err := parseAlarmRule(input.AlarmRule)
	if err != nil {
		return "", awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid AlarmRule: %v", err))
	}
	if err := validateCompositeAcyclic(stores.alarms, input.AlarmName, rule); err != nil {
		return "", err
	}
	if err := validateAlarmActions(input.AlarmActions, input.OKActions, input.InsufficientDataActions); err != nil {
		return "", err
	}

	alarm := cwstore.NewAlarm(input.AlarmName, "", "")
	alarm.AlarmRule = input.AlarmRule
	alarm.AlarmType = cwstore.AlarmTypeCompositeAlarm
	alarm.AlarmDescription = input.AlarmDescription
	alarm.ActionsEnabled = input.ActionsEnabled
	alarm.AlarmActions = input.AlarmActions
	alarm.OKActions = input.OKActions
	alarm.InsufficientDataActions = input.InsufficientDataActions
	alarm.ActionsSuppressor = input.ActionsSuppressor
	alarm.ActionsSuppressorWaitPeriod = input.ActionsSuppressorWaitPeriod
	alarm.ActionsSuppressorExtPeriod = input.ActionsSuppressorExtPeriod
	alarm.Tags = input.Tags

	result, err := s.upsertAlarm(stores.alarms, alarm, cwstore.AlarmTypeCompositeAlarm)
	if err != nil {
		return "", err
	}
	return result.ARN, nil
}

// SetAlarmStateInput holds parameters for SetAlarmState.
type SetAlarmStateInput struct {
	AlarmName       string
	StateValue      string
	StateReason     string
	StateReasonData string
}

// setAlarmStateCore validates input and sets the alarm state.
func (s *CloudWatchService) setAlarmStateCore(stores *cloudwatchStores, input *SetAlarmStateInput) (*cwstore.Alarm, error) {
	if input.AlarmName == "" {
		return nil, ErrInvalidParameter
	}
	if err := validateStateValueFilter(input.StateValue); err != nil {
		return nil, err
	}
	if err := validateStateReason(input.StateReason); err != nil {
		return nil, err
	}

	if err := stores.alarms.SetAlarmState(input.AlarmName, input.StateValue, input.StateReason, input.StateReasonData); err != nil {
		if errors.Is(err, cwstore.ErrAlarmNotFound) {
			return nil, ErrAlarmNotFound
		}
		return nil, err
	}

	alarm, _ := stores.alarms.GetAlarm(input.AlarmName)
	if err := stores.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       input.AlarmName,
		AlarmType:       resolveAlarmType(alarm),
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  input.StateReason,
	}); err != nil {
		logs.Warn("failed to add alarm state history", logs.String("alarm", input.AlarmName), logs.Err(err))
	}

	return alarm, nil
}

// DescribeAlarmHistoryInput holds parameters for DescribeAlarmHistory.
type DescribeAlarmHistoryInput struct {
	AlarmName       string
	AlarmTypes      []string
	HistoryItemType string
	StartDate       time.Time
	EndDate         time.Time
	ScanBy          string
	NextToken       string
	MaxRecords      int
}

// describeAlarmHistoryCore validates input and retrieves alarm history.
func (s *CloudWatchService) describeAlarmHistoryCore(stores *cloudwatchStores, input *DescribeAlarmHistoryInput) ([]*cwstore.AlarmHistoryEntry, string, error) {
	if err := validateHistoryItemType(input.HistoryItemType); err != nil {
		return nil, "", err
	}

	maxRecords := capMaxRecords(input.MaxRecords, 100, maxDescribeAlarmsRecords)

	var startDateMs, endDateMs int64
	if !input.StartDate.IsZero() {
		startDateMs = input.StartDate.UnixMilli()
	}
	if !input.EndDate.IsZero() {
		endDateMs = input.EndDate.UnixMilli()
	}

	opts := cwstore.ListAlarmHistoryOpts{
		AlarmName:       input.AlarmName,
		HistoryItemType: input.HistoryItemType,
		StartDate:       startDateMs,
		EndDate:         endDateMs,
		ScanBy:          input.ScanBy,
		ListOpts:        common.ListOptions{Marker: input.NextToken, MaxItems: maxRecords},
	}
	if len(input.AlarmTypes) > 0 {
		opts.AlarmTypes = make(map[string]bool, len(input.AlarmTypes))
		for _, t := range input.AlarmTypes {
			opts.AlarmTypes[t] = true
		}
	}

	result, err := stores.alarms.ListAlarmHistoryPaginated(opts)
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// EnableAlarmActionsInput holds parameters for EnableAlarmActions.
type EnableAlarmActionsInput struct {
	AlarmNames []string
}

// enableAlarmActionsCore validates input and enables alarm actions.
func (s *CloudWatchService) enableAlarmActionsCore(stores *cloudwatchStores, input *EnableAlarmActionsInput) error {
	if len(input.AlarmNames) == 0 {
		return ErrInvalidParameter
	}
	for _, name := range input.AlarmNames {
		if err := stores.alarms.SetAlarmActionsEnabled(name, true); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return ErrAlarmNotFound
			}
			return err
		}
		addAlarmActionHistory(stores.alarms, name, "Alarm actions enabled")
	}
	return nil
}

// disableAlarmActionsCore validates input and disables alarm actions.
func (s *CloudWatchService) disableAlarmActionsCore(stores *cloudwatchStores, input *EnableAlarmActionsInput) error {
	if len(input.AlarmNames) == 0 {
		return ErrInvalidParameter
	}
	for _, name := range input.AlarmNames {
		if err := stores.alarms.SetAlarmActionsEnabled(name, false); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return ErrAlarmNotFound
			}
			return err
		}
		addAlarmActionHistory(stores.alarms, name, "Alarm actions disabled")
	}
	return nil
}

// PutLogAlarmInput holds parameters for PutLogAlarm.
type PutLogAlarmInput struct {
	AlarmName               string
	AlarmDescription        string
	ComparisonOperator      string
	Threshold               float64
	TreatMissingData        string
	ActionsEnabled          bool
	ScheduledQueryConfig    *cwstore.ScheduledQueryConfig
	AlarmActions            []string
	OKActions               []string
	InsufficientDataActions []string
	ActionLogLineCount      int32
	ActionLogLineRoleArn    string
	QueryResultsToEvaluate  int32
	QueryResultsToAlarm     int32
	Tags                    map[string]string
}

// putLogAlarmCore validates input and creates or updates a log alarm.
func (s *CloudWatchService) putLogAlarmCore(stores *cloudwatchStores, input *PutLogAlarmInput) (string, error) {
	if input.AlarmName == "" {
		return "", awserrors.NewMissingParameter("AlarmName is required")
	}
	if err := validateAlarmDescription(input.AlarmDescription); err != nil {
		return "", err
	}

	alarm := cwstore.NewAlarm(input.AlarmName, "AWS/Logs", "LogQueryResult")
	alarm.AlarmDescription = input.AlarmDescription
	alarm.ComparisonOperator = input.ComparisonOperator
	alarm.Threshold = input.Threshold
	alarm.TreatMissingData = input.TreatMissingData
	alarm.ActionsEnabled = input.ActionsEnabled
	alarm.AlarmActions = input.AlarmActions
	alarm.OKActions = input.OKActions
	alarm.InsufficientDataActions = input.InsufficientDataActions

	if input.ScheduledQueryConfig != nil {
		alarm.ScheduledQueryConfiguration = input.ScheduledQueryConfig
	}
	alarm.ActionLogLineCount = input.ActionLogLineCount
	alarm.ActionLogLineRoleArn = input.ActionLogLineRoleArn
	alarm.QueryResultsToEvaluate = input.QueryResultsToEvaluate
	alarm.QueryResultsToAlarm = input.QueryResultsToAlarm
	alarm.Tags = input.Tags

	result, err := s.upsertAlarm(stores.alarms, alarm, cwstore.AlarmTypeLogAlarm)
	if err != nil {
		return "", err
	}
	return result.ARN, nil
}

// AssociateDatasetKmsKeyInput holds parameters for AssociateDatasetKmsKey.
type AssociateDatasetKmsKeyInput struct {
	DatasetIdentifier string
	KmsKeyArn         string
}

// associateDatasetKmsKeyCore validates input and associates a KMS key.
func (s *CloudWatchService) associateDatasetKmsKeyCore(input *AssociateDatasetKmsKeyInput) error {
	if input.DatasetIdentifier == "" {
		return awserrors.NewMissingParameter("DatasetIdentifier is required")
	}
	if err := validateKmsKeyArn(input.KmsKeyArn); err != nil {
		return err
	}

	s.globalMu.Lock()
	s.datasetKMS[input.DatasetIdentifier] = input.KmsKeyArn
	s.globalMu.Unlock()
	return nil
}

// disassociateDatasetKmsKeyCore validates input and removes KMS key.
func (s *CloudWatchService) disassociateDatasetKmsKeyCore(datasetIdentifier string) error {
	if datasetIdentifier == "" {
		return awserrors.NewMissingParameter("DatasetIdentifier is required")
	}

	s.globalMu.Lock()
	delete(s.datasetKMS, datasetIdentifier)
	s.globalMu.Unlock()
	return nil
}
