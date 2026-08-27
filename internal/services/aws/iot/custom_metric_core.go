package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Custom Metric Core (Device Defender custom metric management). Custom
// metrics feed behaviour criteria from device metric reports; records live
// under the generic-KV "customMetric/" prefix.
// ---------------------------------------------------------------------------

// CreateCustomMetricInput carries the parsed CreateCustomMetric request.
type CreateCustomMetricInput struct {
	MetricName         string
	MetricType         string
	DisplayName        string
	ClientRequestToken string
	Tags               map[string]string
}

// CreateCustomMetricResult is the transport-agnostic result of
// CreateCustomMetric.
type CreateCustomMetricResult struct {
	MetricName interface{}
	MetricArn  string
}

// customMetricTypes is the CustomMetricType enum member set (wire values).
var customMetricTypes = map[string]bool{
	"string-list": true, "ip-address-list": true, "number-list": true, "number": true,
}

// createCustomMetricCore validates and persists a custom metric record.
func (s *IoTService) createCustomMetricCore(store iotstore.IotStoreInterface, in CreateCustomMetricInput) (*CreateCustomMetricResult, error) {
	if !customMetricTypes[in.MetricType] {
		return nil, iotstore.ErrInvalidRequest
	}
	rec, err := s.bulkCreateCore(store, "customMetric", in.MetricName, map[string]interface{}{
		"metricType":         in.MetricType,
		"displayName":        in.DisplayName,
		"clientRequestToken": in.ClientRequestToken,
	})
	if err != nil {
		return nil, err
	}
	arn := iotstore.BuildCustomMetricARN(store.GetAccountID(), store.GetRegion(), bulkName(rec))
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return nil, err
		}
	}
	return &CreateCustomMetricResult{
		MetricName: rec["name"],
		MetricArn:  arn,
	}, nil
}

// deleteCustomMetricCore removes a custom metric record and its tags.
func (s *IoTService) deleteCustomMetricCore(store iotstore.IotStoreInterface, name string) error {
	arn := iotstore.BuildCustomMetricARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return s.bulkDeleteCore(store, "customMetric", name)
}

// CustomMetricRecord is the persisted custom metric record plus its ARN.
type CustomMetricRecord struct {
	Rec map[string]interface{}
	Arn string
}

// describeCustomMetricCore loads a custom metric record and computes its
// ARN. An unknown name yields ErrCustomMetricNotFound.
func (s *IoTService) describeCustomMetricCore(store iotstore.IotStoreInterface, name string) (*CustomMetricRecord, error) {
	rec, exists, err := s.bulkGetCore(store, "customMetric", name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return &CustomMetricRecord{
		Rec: rec,
		Arn: iotstore.BuildCustomMetricARN(store.GetAccountID(), store.GetRegion(), name),
	}, nil
}

// UpdateCustomMetricInput carries the parsed UpdateCustomMetric request.
type UpdateCustomMetricInput struct {
	MetricName  string
	DisplayName string
}

// UpdateCustomMetricResult is the transport-agnostic result of
// UpdateCustomMetric.
type UpdateCustomMetricResult struct {
	Rec map[string]interface{}
	Arn string
}

// updateCustomMetricCore merges the display name into an existing custom
// metric. The displayName member is required by the model.
func (s *IoTService) updateCustomMetricCore(store iotstore.IotStoreInterface, in UpdateCustomMetricInput) (*UpdateCustomMetricResult, error) {
	if in.DisplayName == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec, exists, err := s.bulkUpdateCore(store, "customMetric", in.MetricName, map[string]interface{}{
		"displayName": in.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCustomMetricNotFound
	}
	return &UpdateCustomMetricResult{
		Rec: rec,
		Arn: iotstore.BuildCustomMetricARN(store.GetAccountID(), store.GetRegion(), in.MetricName),
	}, nil
}
