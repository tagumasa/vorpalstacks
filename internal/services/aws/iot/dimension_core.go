package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Dimension Core (Device Defender dimension management). Dimensions scope
// metric behaviours to a value list (for example a set of MQTT topic
// filters); records live under the generic-KV "dimension/" prefix.
// ---------------------------------------------------------------------------

// CreateDimensionInput carries the parsed CreateDimension request.
type CreateDimensionInput struct {
	Name         string
	Type         string
	StringValues []string
	Tags         map[string]string
}

// dimensionTypes is the DimensionType enum member set.
var dimensionTypes = map[string]bool{"TOPIC_FILTER": true}

// CreateDimensionResult is the transport-agnostic result of CreateDimension.
type CreateDimensionResult struct {
	Name interface{}
	Arn  string
}

// createDimensionCore validates and persists a dimension record. The type
// and stringValues members are required by the model.
func (s *IoTService) createDimensionCore(store iotstore.IotStoreInterface, in CreateDimensionInput) (*CreateDimensionResult, error) {
	if !dimensionTypes[in.Type] {
		return nil, iotstore.ErrInvalidRequest
	}
	if len(in.StringValues) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	rec, err := s.bulkCreateCore(store, "dimension", in.Name, map[string]interface{}{
		"type":         in.Type,
		"stringValues": in.StringValues,
	})
	if err != nil {
		return nil, err
	}
	arn := iotstore.BuildDimensionARN(store.GetAccountID(), store.GetRegion(), bulkName(rec))
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return &CreateDimensionResult{
		Name: rec["name"],
		Arn:  arn,
	}, nil
}

// deleteDimensionCore removes a dimension record and its tags.
func (s *IoTService) deleteDimensionCore(store iotstore.IotStoreInterface, name string) error {
	arn := iotstore.BuildDimensionARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return s.bulkDeleteCore(store, "dimension", name)
}

// DimensionRecord is the persisted dimension record plus its ARN.
type DimensionRecord struct {
	Rec map[string]interface{}
	Arn string
}

// describeDimensionCore loads a dimension record and computes its ARN.
// An unknown name yields ErrDimensionNotFound.
func (s *IoTService) describeDimensionCore(store iotstore.IotStoreInterface, name string) (*DimensionRecord, error) {
	rec, exists, err := s.bulkGetCore(store, "dimension", name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return &DimensionRecord{
		Rec: rec,
		Arn: iotstore.BuildDimensionARN(store.GetAccountID(), store.GetRegion(), name),
	}, nil
}

// UpdateDimensionInput carries the parsed UpdateDimension request.
type UpdateDimensionInput struct {
	Name         string
	StringValues []string
}

// UpdateDimensionResult is the transport-agnostic result of UpdateDimension.
type UpdateDimensionResult struct {
	Rec map[string]interface{}
	Arn string
}

// updateDimensionCore merges the value list into an existing dimension.
// The stringValues member is required and carries at least one entry.
func (s *IoTService) updateDimensionCore(store iotstore.IotStoreInterface, in UpdateDimensionInput) (*UpdateDimensionResult, error) {
	if len(in.StringValues) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	rec, exists, err := s.bulkUpdateCore(store, "dimension", in.Name, map[string]interface{}{
		"stringValues": in.StringValues,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return &UpdateDimensionResult{
		Rec: rec,
		Arn: iotstore.BuildDimensionARN(store.GetAccountID(), store.GetRegion(), in.Name),
	}, nil
}
