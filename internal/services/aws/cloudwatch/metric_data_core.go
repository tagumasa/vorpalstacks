package cloudwatch

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/pkg/metricmath"
)

// GetMetricWidgetImageInput holds parameters for GetMetricWidgetImage.
type GetMetricWidgetImageInput struct {
	MetricWidget string
	OutputFormat string
}

// getMetricWidgetImageCore validates input and retrieves the metric
// widget image. Returns the parsed widget definition, validated format,
// and any validation error.
func (s *CloudWatchService) getMetricWidgetImageCore(input *GetMetricWidgetImageInput) (*widgetDef, string, error) {
	if input.MetricWidget == "" {
		return nil, "", ErrInvalidParameter
	}
	if !json.Valid([]byte(input.MetricWidget)) {
		return nil, "", ErrInvalidParameter
	}

	format := input.OutputFormat
	if format == "" {
		format = "png"
	}
	if err := validateOutputFormat(format); err != nil {
		return nil, "", err
	}

	def, err := parseWidgetDefinition(input.MetricWidget)
	if err != nil {
		return nil, "", ErrInvalidParameter
	}

	return def, format, nil
}

// GetMetricDataInput holds validated parameters for GetMetricData.
type GetMetricDataInput struct {
	StartTime         time.Time
	EndTime           time.Time
	MetricDataQueries []cwstore.MetricDataQuery
	ScanBy            string
}

// MetricDataQueryOutcome holds the evaluated outcome of one GetMetricData
// query: the query definition, the resolved data-point series (NaN grid
// periods already filtered), or the evaluation error message. Only the
// queries the response includes are present, in request order.
type MetricDataQueryOutcome struct {
	Query      cwstore.MetricDataQuery
	Points     []metricmath.DataPoint
	ErrMessage string
}

// getMetricDataCore validates GetMetricData input and evaluates the query
// set: every MetricStat query is resolved against the metric store, then
// expression queries are evaluated in dependency order on the resolved
// series.
func (s *CloudWatchService) getMetricDataCore(stores *cloudwatchStores, input *GetMetricDataInput) ([]MetricDataQueryOutcome, error) {
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return nil, awserrors.NewMissingParameter("StartTime and EndTime are required")
	}
	if !input.StartTime.Before(input.EndTime) {
		return nil, awserrors.NewInvalidParameterValueException("StartTime must be before EndTime")
	}
	if len(input.MetricDataQueries) == 0 {
		return nil, awserrors.NewMissingParameter("MetricDataQueries is required")
	}
	if input.ScanBy != "" && input.ScanBy != "TimestampAscending" && input.ScanBy != "TimestampDescending" {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid ScanBy: %s", input.ScanBy))
	}
	metricStore := stores.metrics
	startTime := input.StartTime
	endTime := input.EndTime
	metricDataQueries := input.MetricDataQueries

	// Determine whether any expression queries exist. If so, we need
	// to generate a period-aligned grid for MetricStat results so that
	// FILL() and other functions can operate on missing data points.
	hasExpressions := false
	for i := range metricDataQueries {
		if metricDataQueries[i].Expression != "" {
			hasExpressions = true
			break
		}
	}

	// Phase 1: Evaluate all MetricStat queries and store results by ID
	// for use by Expression queries that reference them.
	queryResults := make(map[string][]metricmath.DataPoint, len(metricDataQueries))
	queryErrors := make(map[string]string)
	queryOrder := make([]string, 0, len(metricDataQueries))
	queryLookup := make(map[string]*cwstore.MetricDataQuery, len(metricDataQueries))

	for i := range metricDataQueries {
		query := &metricDataQueries[i]
		queryLookup[query.Id] = query
		queryOrder = append(queryOrder, query.Id)

		if query.MetricStat != nil {
			statLower := strings.ToLower(query.MetricStat.Stat)
			isExtended := !isBasicStatistic(statLower)

			mq := cwstore.MetricQuery{
				Namespace:  query.MetricStat.Metric.Namespace,
				MetricName: query.MetricStat.Metric.MetricName,
				Dimensions: query.MetricStat.Metric.Dimensions,
				StartTime:  startTime,
				EndTime:    endTime,
				Period:     query.MetricStat.Period,
			}
			if isExtended {
				mq.ExtendedStatistics = []string{query.MetricStat.Stat}
			} else {
				mq.Statistics = []string{query.MetricStat.Stat}
			}

			stats, err := metricStore.GetMetricStatistics(mq)
			if err != nil {
				queryResults[query.Id] = nil
				queryErrors[query.Id] = fmt.Sprintf("failed to query metric statistics: %v", err)
				continue
			}

			// Convert MetricStatistics to DataPoint series for the
			// expression evaluator. Extract the requested stat value
			// rather than defaulting to Average.
			dataPoints := make([]metricmath.DataPoint, 0, len(stats))
			for _, s := range stats {
				val := extractStatValue(s, statLower, query.MetricStat.Stat, isExtended)
				dataPoints = append(dataPoints, metricmath.DataPoint{
					Timestamp: s.Timestamp,
					Value:     val,
				})
			}

			// When expressions are present, expand the series onto a
			// period-aligned grid with NaN for missing periods. This
			// enables FILL() and other functions to operate on gaps.
			if hasExpressions && query.MetricStat.Period > 0 {
				dataPoints = expandToPeriodGrid(dataPoints, startTime, endTime, query.MetricStat.Period)
			}

			queryResults[query.Id] = dataPoints
		}
	}

	// Phase 2: Evaluate Expression queries in dependency order.
	// Expressions can reference other queries by ID (e.g., m1 + m2).
	// We iterate until all expressions are resolved or no progress is
	// made (indicating a circular or unresolved reference).
	exprPending := make(map[string]string)
	for _, id := range queryOrder {
		q := queryLookup[id]
		if q.Expression != "" {
			// ANOMALY_DETECTION_BAND is handled specially because it
			// produces two series (upper and lower bounds) and needs
			// the referenced metric's data already computed.
			if containsAnomalyDetectionBand(q.Expression) {
				resolveAnomalyDetectionBand(q, queryResults)
				continue
			}
			exprPending[id] = q.Expression
		}
	}

	for len(exprPending) > 0 {
		progress := false
		for id, expr := range exprPending {
			ast, err := metricmath.Parse(expr)
			if err != nil {
				delete(exprPending, id)
				queryErrors[id] = fmt.Sprintf("failed to parse expression: %v", err)
				continue
			}

			// Check if all referenced variables are available.
			refs := ast.References()
			ready := true
			for _, ref := range refs {
				if _, ok := queryResults[ref]; !ok {
					if _, isExpr := exprPending[ref]; isExpr {
						ready = false
						break
					}
				}
			}
			if !ready {
				continue
			}

			result, err := ast.Eval(queryResults)
			if err != nil {
				queryErrors[id] = fmt.Sprintf("failed to evaluate expression: %v", err)
			} else {
				queryResults[id] = result
			}
			delete(exprPending, id)
			progress = true
		}
		if !progress {
			break
		}
	}

	// Phase 3: Select the queries with ReturnData=true (default) for the
	// response. Queries with ReturnData=false are intermediate values used
	// only for expression resolution and must not appear in the response.
	// NaN values (from period grid expansion) are filtered out so that the
	// response only contains periods with actual data.
	outcomes := make([]MetricDataQueryOutcome, 0, len(queryOrder))
	for _, id := range queryOrder {
		query := queryLookup[id]
		if !query.ReturnData {
			continue
		}
		if errMsg, hasErr := queryErrors[id]; hasErr {
			outcomes = append(outcomes, MetricDataQueryOutcome{
				Query:      *query,
				ErrMessage: errMsg,
			})
			continue
		}
		if query.Expression != "" && query.MetricStat == nil {
			dataPoints, ok := queryResults[id]
			if !ok {
				continue
			}
			outcomes = append(outcomes, MetricDataQueryOutcome{
				Query:  *query,
				Points: filterNaN(dataPoints),
			})
		} else if query.MetricStat != nil {
			outcomes = append(outcomes, MetricDataQueryOutcome{
				Query:  *query,
				Points: filterNaN(queryResults[id]),
			})
		}
	}
	return outcomes, nil
}

// WidgetMetricData holds the queried statistics for one resolvable widget
// metric, together with its position in the widget definition (the chart
// palette colour derives from that position).
type WidgetMetricData struct {
	MetricIndex int
	Metric      widgetMetric
	Stats       []*cwstore.MetricStatistics
}

// queryWidgetMetricsCore queries the metric store for every metric of the
// parsed widget definition. Metrics without a namespace or metric name and
// metrics whose query fails are skipped — the widget renders only the
// resolvable series.
func (s *CloudWatchService) queryWidgetMetricsCore(stores *cloudwatchStores, def *widgetDef) ([]WidgetMetricData, error) {
	var data []WidgetMetricData
	for i, m := range def.Metrics {
		if m.Namespace == "" || m.MetricName == "" {
			continue
		}

		mq := cwstore.MetricQuery{
			Namespace:  m.Namespace,
			MetricName: m.MetricName,
			Dimensions: m.Dimensions,
			StartTime:  def.Start,
			EndTime:    def.End,
			Period:     m.Period,
		}
		statLower := strings.ToLower(m.Stat)
		if isBasicStatistic(statLower) {
			mq.Statistics = []string{m.Stat}
		} else {
			mq.ExtendedStatistics = []string{m.Stat}
		}

		stats, err := stores.metrics.GetMetricStatistics(mq)
		if err != nil {
			continue
		}

		data = append(data, WidgetMetricData{MetricIndex: i, Metric: m, Stats: stats})
	}
	return data, nil
}

// describeAlarmContributorsCore validates input, resolves the alarm, and
// computes the anomaly-band contributors from the alarm's referenced
// metric history.
func (s *CloudWatchService) describeAlarmContributorsCore(stores *cloudwatchStores, alarmName string) ([]alarmContributor, error) {
	if alarmName == "" {
		return nil, awserrors.NewMissingParameter("AlarmName is required")
	}

	alarm, err := stores.alarms.GetAlarm(alarmName)
	if err != nil || alarm == nil {
		return nil, awserrors.NewResourceNotFoundException("Alarm", alarmName)
	}

	if len(alarm.Metrics) == 0 || !hasAnomalyDetectionBand(alarm.Metrics) {
		return nil, awserrors.NewInvalidParameterValueException(
			"The specified alarm does not use anomaly detection")
	}

	contributors, err := computeAlarmContributors(alarm, stores.metrics)
	if err != nil {
		return nil, awserrors.NewInternalFailureException(
			fmt.Sprintf("failed to compute alarm contributors: %v", err))
	}
	return contributors, nil
}

// alarmContributor represents a single contributor to an anomaly
// detection alarm evaluation.
type alarmContributor struct {
	Timestamp   string  `json:"timestamp"`
	MetricValue float64 `json:"metricValue"`
	BandUpper   float64 `json:"bandUpper"`
	BandLower   float64 `json:"bandLower"`
}

// computeAlarmContributors computes the anomaly band for the given
// alarm and returns data points that fall outside the band as
// contributors.
func computeAlarmContributors(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) ([]alarmContributor, error) {
	ctx := prepareAnomalyBand(alarm, metricStore)
	if ctx == nil {
		return nil, fmt.Errorf("no anomaly detection band or metric data available")
	}

	var contributors []alarmContributor
	for i, s := range ctx.stats {
		if s.Timestamp.Before(ctx.startTime) || s.Timestamp.After(ctx.endTime) {
			continue
		}
		val := statisticValue(s, ctx.statLower)
		if val < ctx.lower[i] || val > ctx.upper[i] {
			contributors = append(contributors, alarmContributor{
				Timestamp:   s.Timestamp.Format(time.RFC3339),
				MetricValue: val,
				BandUpper:   ctx.upper[i],
				BandLower:   ctx.lower[i],
			})
		}
	}

	return contributors, nil
}
