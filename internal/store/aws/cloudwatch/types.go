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
	Name                    string            `json:"name"`
	ARN                     string            `json:"arn"`
	Namespace               string            `json:"namespace"`
	MetricName              string            `json:"metricName"`
	Dimensions              []Dimension       `json:"dimensions,omitempty"`
	ComparisonOperator      string            `json:"comparisonOperator"`
	Threshold               float64           `json:"threshold"`
	EvaluationPeriods       int32             `json:"evaluationPeriods"`
	DatapointsToAlarm       int32             `json:"datapointsToAlarm,omitempty"`
	Period                  int32             `json:"period"`
	Statistic               string            `json:"statistic"`
	TreatMissingData        string            `json:"treatMissingData,omitempty"`
	AlarmDescription        string            `json:"alarmDescription,omitempty"`
	ActionsEnabled          bool              `json:"actionsEnabled"`
	AlarmActions            []string          `json:"alarmActions,omitempty"`
	OKActions               []string          `json:"okActions,omitempty"`
	InsufficientDataActions []string          `json:"insufficientDataActions,omitempty"`
	State                   string            `json:"state"`
	StateReason             string            `json:"stateReason,omitempty"`
	StateUpdatedTimestamp   time.Time         `json:"stateUpdatedTimestamp"`
	CreatedAt               time.Time         `json:"createdAt"`
	Tags                    map[string]string `json:"tags,omitempty"`
	AlarmRule               string            `json:"alarmRule,omitempty"`
	AlarmType               string            `json:"alarmType,omitempty"`
}

// AlarmType constants
const (
	AlarmTypeMetricAlarm    = "MetricAlarm"
	AlarmTypeCompositeAlarm = "CompositeAlarm"
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
	Name      string    `json:"name"`
	ARN       string    `json:"arn"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
