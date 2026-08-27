package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the fleet-metric family. Records live in the
// generic-KV namespace under the fleetMetric category. Handlers on both
// protocol planes are thin adapters; validation and persistence live here
// only.
//
// The write operations share one DTO; every member the model marks
// Required: No is a pointer so an omitted member is distinguishable from an
// explicitly provided value, and the update core merges only the provided
// members — API_UpdateFleetMetric preserves omitted optional members
// (period, queryString, aggregationType, aggregationField, description,
// queryVersion, unit), with indexName the sole required member alongside
// metricName.

// fleetMetricCategory is the generic-KV category prefix for fleet metrics.
const fleetMetricCategory = "fleetMetric"

// FleetMetricInput is the service-layer DTO for the fleet-metric write
// operations. Tags apply on create only; ExpectedVersion drives the
// documented optimistic-locking check on update.
type FleetMetricInput struct {
	MetricName       string
	IndexName        string
	QueryString      *string
	AggregationType  map[string]interface{}
	Period           *int64
	AggregationField *string
	Unit             *string
	Description      *string
	QueryVersion     *string
	ExpectedVersion  *int64
	Tags             map[string]string
}

// FleetMetricResult carries the stored record and the metric ARN.
type FleetMetricResult struct {
	Record map[string]interface{}
	ARN    string
}

// FleetMetricSummary is one ListFleetMetrics entry.
type FleetMetricSummary struct {
	Name string
	ARN  string
}

// fleetMetricFields builds the persistence map from the provided members
// only, so a bulk update merge never overwrites an omitted member with a
// zero value.
func fleetMetricFields(in FleetMetricInput) map[string]interface{} {
	fields := map[string]interface{}{}
	if in.QueryString != nil {
		fields["queryString"] = *in.QueryString
	}
	if in.AggregationType != nil {
		fields["aggregationType"] = in.AggregationType
	}
	if in.Period != nil {
		fields["period"] = *in.Period
	}
	if in.AggregationField != nil {
		fields["aggregationField"] = *in.AggregationField
	}
	if in.Unit != nil {
		fields["unit"] = *in.Unit
	}
	if in.Description != nil {
		fields["description"] = *in.Description
	}
	if in.QueryVersion != nil {
		fields["queryVersion"] = *in.QueryVersion
	}
	if in.IndexName != "" {
		fields["indexName"] = in.IndexName
	}
	return fields
}

// validateFleetMetricCreate enforces the documented required members of
// API_CreateFleetMetric: metricName, queryString, aggregationType (with its
// required name), period, and aggregationField.
func validateFleetMetricCreate(in FleetMetricInput) error {
	if in.MetricName == "" {
		return iotstore.ErrMissingParam
	}
	if in.QueryString == nil || *in.QueryString == "" {
		return iotstore.ErrMissingParam
	}
	if in.AggregationField == nil || *in.AggregationField == "" {
		return iotstore.ErrMissingParam
	}
	if in.Period == nil {
		return iotstore.ErrMissingParam
	}
	if name, _ := in.AggregationType["name"].(string); in.AggregationType == nil || name == "" {
		return iotstore.ErrMissingParam
	}
	if err := validateFleetMetricMembers(in); err != nil {
		return err
	}
	return validateFleetMetricPeriod(*in.Period)
}

func validateFleetMetricPeriod(period int64) error {
	if period < MinFleetMetricPeriod || period > MaxFleetMetricPeriod {
		return iotstore.ErrInvalidRequest
	}
	if period%60 != 0 {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// validateFleetMetricMembers enforces the Smithy shape constraints on
// every provided member of either write operation: metricName 1–128
// matching [a-zA-Z0-9_\-\.]+, queryString and aggregationField minimum
// length 1, description 0–1024 printable characters, indexName 1–128
// matching [a-zA-Z0-9:_-]+, and the FleetMetricUnit enum. Aggregation
// violations — a name outside the AggregationTypeName enum or a value
// outside 1–12 [a-zA-Z0-9]+ — surface as InvalidAggregationException,
// the error both write operations document; every other violation is
// InvalidRequestException.
func validateFleetMetricMembers(in FleetMetricInput) error {
	if len(in.MetricName) < MinFleetMetricNameLength || len(in.MetricName) > MaxFleetMetricNameLength ||
		!fleetMetricNamePattern.MatchString(in.MetricName) {
		return iotstore.ErrInvalidRequest
	}
	if in.QueryString != nil && *in.QueryString == "" {
		return iotstore.ErrInvalidRequest
	}
	if in.AggregationField != nil && *in.AggregationField == "" {
		return iotstore.ErrInvalidRequest
	}
	if in.Description != nil {
		d := *in.Description
		if len(d) > MaxFleetMetricDescriptionLength || !fleetMetricDescriptionPattern.MatchString(d) {
			return iotstore.ErrInvalidRequest
		}
	}
	if in.IndexName != "" &&
		(len(in.IndexName) < MinFleetMetricIndexNameLength || len(in.IndexName) > MaxFleetMetricIndexNameLength ||
			!fleetMetricIndexNamePattern.MatchString(in.IndexName)) {
		return iotstore.ErrInvalidRequest
	}
	if in.Unit != nil {
		if _, ok := fleetMetricUnits[*in.Unit]; !ok {
			return iotstore.ErrInvalidRequest
		}
	}
	if in.AggregationType != nil {
		name, _ := in.AggregationType["name"].(string)
		if !isValidAggregationTypeName(name) {
			return iotstore.ErrInvalidAggregation
		}
		if values, ok := in.AggregationType["values"].([]interface{}); ok {
			for _, v := range values {
				s, _ := v.(string)
				if len(s) < 1 || len(s) > MaxAggregationValueLength || !aggregationValuePattern.MatchString(s) {
					return iotstore.ErrInvalidAggregation
				}
			}
		}
	}
	return nil
}

// fleetMetricRecordVersion reads the record's version counter. Versions
// start at 1 on creation and increment on every update.
func fleetMetricRecordVersion(rec map[string]interface{}) int64 {
	switch v := rec["version"].(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func (s *IoTService) createFleetMetricCore(store iotstore.IotStoreInterface, in FleetMetricInput) (*FleetMetricResult, error) {
	if err := validateFleetMetricCreate(in); err != nil {
		return nil, err
	}
	fields := fleetMetricFields(in)
	fields["version"] = int64(1)
	rec, err := s.bulkCreateCore(store, fleetMetricCategory, in.MetricName, fields)
	if err != nil {
		return nil, err
	}
	name := bulkName(rec)
	arn := iotstore.BuildFleetMetricARN(store.GetAccountID(), store.GetRegion(), name)
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return &FleetMetricResult{
		Record: rec,
		ARN:    arn,
	}, nil
}

func (s *IoTService) describeFleetMetricCore(store iotstore.IotStoreInterface, metricName string) (*FleetMetricResult, error) {
	rec, exists, err := s.bulkGetCore(store, fleetMetricCategory, metricName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return &FleetMetricResult{
		Record: rec,
		ARN:    iotstore.BuildFleetMetricARN(store.GetAccountID(), store.GetRegion(), bulkName(rec)),
	}, nil
}

func (s *IoTService) updateFleetMetricCore(store iotstore.IotStoreInterface, in FleetMetricInput) (*FleetMetricResult, error) {
	if in.MetricName == "" || in.IndexName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if err := validateFleetMetricMembers(in); err != nil {
		return nil, err
	}
	if in.Period != nil {
		if err := validateFleetMetricPeriod(*in.Period); err != nil {
			return nil, err
		}
	}
	rec, exists, err := s.bulkGetCore(store, fleetMetricCategory, in.MetricName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	if in.ExpectedVersion != nil && fleetMetricRecordVersion(rec) != *in.ExpectedVersion {
		return nil, iotstore.ErrVersionConflict
	}
	merge := fleetMetricFields(in)
	merge["version"] = fleetMetricRecordVersion(rec) + 1
	updated, exists, err := s.bulkUpdateCore(store, fleetMetricCategory, in.MetricName, merge)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}
	return &FleetMetricResult{
		Record: updated,
		ARN:    iotstore.BuildFleetMetricARN(store.GetAccountID(), store.GetRegion(), bulkName(updated)),
	}, nil
}

func (s *IoTService) deleteFleetMetricCore(store iotstore.IotStoreInterface, metricName string) error {
	if err := s.bulkDeleteCore(store, fleetMetricCategory, metricName); err != nil {
		return err
	}
	arn := iotstore.BuildFleetMetricARN(store.GetAccountID(), store.GetRegion(), metricName)
	return store.DeleteAllTags(arn)
}

func (s *IoTService) listFleetMetricsCore(store iotstore.IotStoreInterface) ([]FleetMetricSummary, error) {
	items, err := s.bulkListCore(store, fleetMetricCategory)
	if err != nil {
		return nil, err
	}
	result := make([]FleetMetricSummary, 0, len(items))
	for _, item := range items {
		name := bulkName(item)
		result = append(result, FleetMetricSummary{
			Name: name,
			ARN:  iotstore.BuildFleetMetricARN(store.GetAccountID(), store.GetRegion(), name),
		})
	}
	return result, nil
}
