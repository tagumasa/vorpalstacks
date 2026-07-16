package chunk

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

// CloudWatchMetricEntry represents a metric entry for CloudWatch storage.
// This structure is used to store and retrieve metric data in the chunk storage system.
type CloudWatchMetricEntry struct {
	Namespace         string                  `json:"namespace"`
	MetricName        string                  `json:"metricName"`
	Dimensions        []CloudWatchDimension   `json:"dimensions,omitempty"`
	Time              time.Time               `json:"timestamp"`
	Value             float64                 `json:"value,omitempty"`
	Values            []float64               `json:"values,omitempty"`
	Counts            []float64               `json:"counts,omitempty"`
	Unit              string                  `json:"unit,omitempty"`
	StorageResolution int32                   `json:"storageResolution,omitempty"`
	StatisticValues   *CloudWatchStatisticSet `json:"statisticValues,omitempty"`
}

// CloudWatchDimension represents a dimension pair for CloudWatch metrics.
// Dimensions provide additional context for metric data, allowing for more granular grouping.
type CloudWatchDimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CloudWatchStatisticSet contains statistical values for a CloudWatch metric.
// This includes the maximum, minimum, sample count, and sum of the metric values.
type CloudWatchStatisticSet struct {
	Maximum     float64 `json:"maximum"`
	Minimum     float64 `json:"minimum"`
	SampleCount float64 `json:"sampleCount"`
	Sum         float64 `json:"sum"`
}

// Timestamp returns the Unix timestamp in nanoseconds for the metric entry.
// This is used as the key for storing and retrieving metric data.
func (e *CloudWatchMetricEntry) Timestamp() int64 {
	return e.Time.UnixNano()
}

// Message serialises the metric entry to JSON bytes for storage.
// This is used as the value when persisting the metric to chunk storage.
func (e *CloudWatchMetricEntry) Message() []byte {
	data, _ := json.Marshal(e)
	return data
}

// BuildDimensionsKey constructs a sorted, comma-separated key from the metric dimensions.
// This key is used for grouping and querying metrics with specific dimension values.
// If no dimensions are present, it returns "_empty" as a placeholder.
func (e *CloudWatchMetricEntry) BuildDimensionsKey() string {
	if len(e.Dimensions) == 0 {
		return "_empty"
	}
	parts := make([]string, len(e.Dimensions))
	for i, d := range e.Dimensions {
		parts[i] = d.Name + "=" + d.Value
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// NamespaceKey returns a composite key combining the namespace and metric name.
// This is used as part of the lookup key for retrieving metric data.
func (e *CloudWatchMetricEntry) NamespaceKey() string {
	return e.Namespace + "#" + e.MetricName
}
