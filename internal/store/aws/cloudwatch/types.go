// Package cloudwatch provides CloudWatch storage functionality for vorpalstacks.
package cloudwatch

import (
	"time"
)

// Dimension represents a dimension for a CloudWatch metric.
type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// StandardUnit represents the unit of measurement for a CloudWatch metric.
type StandardUnit string

const (
	// UnitSeconds represents seconds.
	UnitSeconds StandardUnit = "Seconds"
	// UnitMicroseconds represents microseconds.
	UnitMicroseconds StandardUnit = "Microseconds"
	// UnitMilliseconds represents milliseconds.
	UnitMilliseconds StandardUnit = "Milliseconds"
	// UnitBytes represents bytes.
	UnitBytes StandardUnit = "Bytes"
	// UnitKilobytes represents kilobytes.
	UnitKilobytes StandardUnit = "Kilobytes"
	// UnitMegabytes represents megabytes.
	UnitMegabytes StandardUnit = "Megabytes"
	// UnitGigabytes represents gigabytes.
	UnitGigabytes StandardUnit = "Gigabytes"
	// UnitTerabytes represents terabytes.
	UnitTerabytes StandardUnit = "Terabytes"
	// UnitBits represents bits.
	UnitBits StandardUnit = "Bits"
	// UnitKilobits represents kilobits.
	UnitKilobits StandardUnit = "Kilobits"
	// UnitMegabits represents megabits.
	UnitMegabits StandardUnit = "Megabits"
	// UnitGigabits represents gigabits.
	UnitGigabits StandardUnit = "Gigabits"
	// UnitTerabits represents terabits.
	UnitTerabits StandardUnit = "Terabits"
	// UnitPercent represents percent.
	UnitPercent StandardUnit = "Percent"
	// UnitCount represents count.
	UnitCount StandardUnit = "Count"
	// UnitBytesPerSecond represents bytes per second.
	UnitBytesPerSecond StandardUnit = "Bytes/Second"
	// UnitKilobytesPerSecond represents kilobytes per second.
	UnitKilobytesPerSecond StandardUnit = "Kilobytes/Second"
	// UnitMegabytesPerSecond represents megabytes per second.
	UnitMegabytesPerSecond StandardUnit = "Megabytes/Second"
	// UnitGigabytesPerSecond represents gigabytes per second.
	UnitGigabytesPerSecond StandardUnit = "Gigabytes/Second"
	// UnitTerabytesPerSecond represents terabytes per second.
	UnitTerabytesPerSecond StandardUnit = "Terabytes/Second"
	// UnitBitsPerSecond represents bits per second.
	UnitBitsPerSecond StandardUnit = "Bits/Second"
	// UnitKilobitsPerSecond represents kilobits per second.
	UnitKilobitsPerSecond StandardUnit = "Kilobits/Second"
	// UnitMegabitsPerSecond represents megabits per second.
	UnitMegabitsPerSecond StandardUnit = "Megabits/Second"
	// UnitGigabitsPerSecond represents gigabits per second.
	UnitGigabitsPerSecond StandardUnit = "Gigabits/Second"
	// UnitTerabitsPerSecond represents terabits per second.
	UnitTerabitsPerSecond StandardUnit = "Terabits/Second"
	// UnitCountPerSecond represents count per second.
	UnitCountPerSecond StandardUnit = "Count/Second"
	// UnitNone represents no unit.
	UnitNone StandardUnit = "None"
)

// MetricDatum represents a metric data point for CloudWatch.
type MetricDatum struct {
	Namespace         string        `json:"namespace,omitempty"`
	MetricName        string        `json:"metricName"`
	Value             float64       `json:"value,omitempty"`
	HasValue          bool          `json:"-"`
	Values            []float64     `json:"values,omitempty"`
	Counts            []float64     `json:"counts,omitempty"`
	Timestamp         time.Time     `json:"timestamp"`
	Unit              StandardUnit  `json:"unit,omitempty"`
	Dimensions        []Dimension   `json:"dimensions,omitempty"`
	StorageResolution int32         `json:"storageResolution,omitempty"`
	StatisticValues   *StatisticSet `json:"statisticValues,omitempty"`
}

// StatisticSet represents a set of statistical values for a metric.
type StatisticSet struct {
	SampleCount float64 `json:"sampleCount"`
	Sum         float64 `json:"sum"`
	Minimum     float64 `json:"minimum"`
	Maximum     float64 `json:"maximum"`
}

// Alarm represents a CloudWatch alarm.
type Alarm struct {
	Name                        string            `json:"name"`
	ARN                         string            `json:"arn"`
	Namespace                   string            `json:"namespace"`
	MetricName                  string            `json:"metricName"`
	Dimensions                  []Dimension       `json:"dimensions,omitempty"`
	ComparisonOperator          string            `json:"comparisonOperator"`
	Threshold                   float64           `json:"threshold"`
	EvaluationPeriods           int32             `json:"evaluationPeriods"`
	DatapointsToAlarm           int32             `json:"datapointsToAlarm,omitempty"`
	Period                      int32             `json:"period"`
	Statistic                   string            `json:"statistic"`
	ExtendedStatistic           string            `json:"extendedStatistic,omitempty"`
	TreatMissingData            string            `json:"treatMissingData,omitempty"`
	AlarmDescription            string            `json:"alarmDescription,omitempty"`
	Unit                        StandardUnit      `json:"unit,omitempty"`
	ActionsEnabled              bool              `json:"actionsEnabled"`
	AlarmActions                []string          `json:"alarmActions,omitempty"`
	OKActions                   []string          `json:"okActions,omitempty"`
	InsufficientDataActions     []string          `json:"insufficientDataActions,omitempty"`
	State                       string            `json:"state"`
	StateReason                 string            `json:"stateReason,omitempty"`
	StateReasonData             string            `json:"stateReasonData,omitempty"`
	StateUpdatedTimestamp       time.Time         `json:"stateUpdatedTimestamp"`
	CreatedAt                   time.Time         `json:"createdAt"`
	Tags                        map[string]string `json:"tags,omitempty"`
	AlarmRule                   string            `json:"alarmRule,omitempty"`
	ActionsSuppressedBy         string            `json:"actionsSuppressedBy,omitempty"`
	ActionsSuppressedReason     string            `json:"actionsSuppressedReason,omitempty"`
	ActionsSuppressor           string            `json:"actionsSuppressor,omitempty"`
	ActionsSuppressorWaitPeriod int32             `json:"actionsSuppressorWaitPeriod,omitempty"`
	ActionsSuppressorExtPeriod  int32             `json:"actionsSuppressorExtPeriod,omitempty"`
	Metrics                     []MetricDataQuery `json:"metrics,omitempty"`
	ThresholdMetricID           string            `json:"thresholdMetricId,omitempty"`
	AlarmType                   string            `json:"alarmType,omitempty"`
	// PromQL evaluation criteria (EvaluationCriteria.PromQLCriteria).
	EvaluationCriteria *EvaluationCriteria `json:"evaluationCriteria,omitempty"`
	// Log alarm fields (PutLogAlarm).
	ActionLogLineCount          int32                 `json:"actionLogLineCount,omitempty"`
	ActionLogLineRoleArn        string                `json:"actionLogLineRoleArn,omitempty"`
	QueryResultsToEvaluate      int32                 `json:"queryResultsToEvaluate,omitempty"`
	QueryResultsToAlarm         int32                 `json:"queryResultsToAlarm,omitempty"`
	ScheduledQueryConfiguration *ScheduledQueryConfig `json:"scheduledQueryConfiguration,omitempty"`
}

// EvaluationCriteria holds the criteria for evaluating an alarm.
// Currently only PromQLCriteria is supported.
type EvaluationCriteria struct {
	PromQLCriteria *AlarmPromQLCriteria `json:"promqlCriteria,omitempty"`
}

// AlarmPromQLCriteria holds PromQL-based alarm evaluation parameters.
type AlarmPromQLCriteria struct {
	Query          string `json:"query,omitempty"`
	PendingPeriod  int32  `json:"pendingPeriod,omitempty"`
	RecoveryPeriod int32  `json:"recoveryPeriod,omitempty"`
}

// ScheduledQueryConfig holds the CloudWatch Logs scheduled query
// configuration that backs a log alarm (Smithy ScheduledQueryConfiguration).
type ScheduledQueryConfig struct {
	QueryString           string            `json:"queryString,omitempty"`
	LogGroupIdentifiers   []string          `json:"logGroupIdentifiers,omitempty"`
	QueryARN              string            `json:"queryArn,omitempty"`
	ScheduledQueryRoleARN string            `json:"scheduledQueryRoleArn,omitempty"`
	ScheduleConfiguration *ScheduleConfig   `json:"scheduleConfiguration,omitempty"`
	AggregationExpression string            `json:"aggregationExpression,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
}

// ScheduleConfig holds the schedule expression and time-range offsets
// for a scheduled query (Smithy ScheduleConfiguration).
type ScheduleConfig struct {
	ScheduleExpression string `json:"scheduleExpression,omitempty"`
	StartTimeOffset    int64  `json:"startTimeOffset,omitempty"`
	EndTimeOffset      int64  `json:"endTimeOffset,omitempty"`
}

// AlarmType constants
const (
	AlarmTypeMetricAlarm    = "MetricAlarm"
	AlarmTypeCompositeAlarm = "CompositeAlarm"
	AlarmTypeLogAlarm       = "LogAlarm"
)

// HistoryItemType constants
const (
	HistoryItemTypeConfigurationUpdate = "ConfigurationUpdate"
	HistoryItemTypeStateUpdate         = "StateUpdate"
	HistoryItemTypeAction              = "Action"
)

// AlarmHistoryEntry represents an entry in an alarm's history.
type AlarmHistoryEntry struct {
	AlarmName       string `json:"alarmName"`
	AlarmType       string `json:"alarmType"`
	Timestamp       int64  `json:"timestamp"`
	HistoryItemType string `json:"historyItemType"`
	HistorySummary  string `json:"historySummary"`
	HistoryData     string `json:"historyData,omitempty"`
}

// Dashboard represents a CloudWatch dashboard.
type Dashboard struct {
	Name      string            `json:"name"`
	ARN       string            `json:"arn"`
	Body      string            `json:"body"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// NewAlarm creates a new Alarm with the specified name, namespace, and metric name.
func NewAlarm(name, namespace, metricName string) *Alarm {
	now := time.Now().UTC()
	return &Alarm{
		Name:                  name,
		Namespace:             namespace,
		MetricName:            metricName,
		Dimensions:            []Dimension{},
		ComparisonOperator:    "GreaterThanOrEqualToThreshold",
		EvaluationPeriods:     1,
		Period:                60,
		Statistic:             "Average",
		TreatMissingData:      "missing",
		ActionsEnabled:        true,
		State:                 "INSUFFICIENT_DATA",
		StateUpdatedTimestamp: now,
		CreatedAt:             now,
		Tags:                  make(map[string]string),
	}
}

// MetricDataQuery represents a metric data query for CloudWatch.
type MetricDataQuery struct {
	Id         string      `json:"id"`
	MetricStat *MetricStat `json:"metricStat,omitempty"`
	Expression string      `json:"expression,omitempty"`
	Label      string      `json:"label,omitempty"`
	ReturnData bool        `json:"returnData,omitempty"`
	Period     int32       `json:"period,omitempty"`
}

// MetricStat represents a metric statistic specification.
type MetricStat struct {
	Metric MetricRef    `json:"metric"`
	Period int32        `json:"period"`
	Stat   string       `json:"stat"`
	Unit   StandardUnit `json:"unit,omitempty"`
}

// MetricRef represents a reference to a CloudWatch metric.
type MetricRef struct {
	Namespace  string      `json:"namespace"`
	MetricName string      `json:"metricName"`
	Dimensions []Dimension `json:"dimensions,omitempty"`
}

// AnomalyDetectorType enumerates the kinds of anomaly detector.
const (
	AnomalyDetectorTypeSingleMetric = "SINGLE_METRIC"
	AnomalyDetectorTypeMetricMath   = "METRIC_MATH"
)

// AnomalyDetector represents a CloudWatch anomaly detection model.
// A single-metric detector is keyed on Namespace, MetricName, Dimensions
// and Stat. A metric-math detector is keyed on its MetricDataQueries.
type AnomalyDetector struct {
	ID                  string `json:"id,omitempty"`
	ARN                 string `json:"arn,omitempty"`
	AnomalyDetectorType string `json:"anomalyDetectorType"`
	// Single-metric fields
	Namespace  string      `json:"namespace,omitempty"`
	MetricName string      `json:"metricName,omitempty"`
	Dimensions []Dimension `json:"dimensions,omitempty"`
	Stat       string      `json:"stat,omitempty"`
	AccountID  string      `json:"accountId,omitempty"`
	// Metric characteristics
	AnomalyDetectorConfiguration *AnomalyDetectorConfiguration `json:"configuration,omitempty"`
	MetricCharacteristics        *MetricCharacteristics        `json:"metricCharacteristics,omitempty"`
	// Metric math
	MetricDataQueries []MetricDataQuery `json:"metricDataQueries,omitempty"`
	// State
	State     string    `json:"state,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// Range represents a time range with start and end timestamps.
type Range struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// AnomalyDetectorConfiguration holds configuration for an anomaly detector.
type AnomalyDetectorConfiguration struct {
	ExcludedTimeRanges []Range `json:"excludedTimeRanges,omitempty"`
	MetricTimezone     string  `json:"metricTimezone,omitempty"`
}

// MetricCharacteristics describes characteristics of the metric data
// that affect anomaly detection.
type MetricCharacteristics struct {
	PeriodicSpikes bool `json:"periodicSpikes,omitempty"`
}

// InsightRule represents a CloudWatch Contributor Insights rule.
type InsightRule struct {
	Name                   string            `json:"name"`
	State                  string            `json:"state"`
	Schema                 string            `json:"schema,omitempty"`
	Definition             string            `json:"definition,omitempty"`
	ManagedRule            bool              `json:"managedRule"`
	ApplyOnTransformedLogs bool              `json:"applyOnTransformedLogs,omitempty"`
	Tags                   map[string]string `json:"tags,omitempty"`
	// Managed rule fields
	TemplateName string    `json:"templateName,omitempty"`
	ResourceARN  string    `json:"resourceArn,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

// AlarmMuteRule represents a CloudWatch alarm mute rule that suppresses
// alarm actions for specified alarms during a scheduled time window.
type AlarmMuteRule struct {
	Name            string            `json:"name"`
	ARN             string            `json:"arn,omitempty"`
	Description     string            `json:"description,omitempty"`
	ScheduleExpr    string            `json:"scheduleExpr,omitempty"`
	MutedAlarmNames []string          `json:"mutedAlarmNames,omitempty"`
	StartDate       time.Time         `json:"startDate,omitempty"`
	ExpireDate      time.Time         `json:"expireDate,omitempty"`
	Status          string            `json:"status,omitempty"`
	MuteType        string            `json:"muteType,omitempty"`
	Tags            map[string]string `json:"tags,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty"`
	UpdatedAt       time.Time         `json:"updatedAt,omitempty"`
}
