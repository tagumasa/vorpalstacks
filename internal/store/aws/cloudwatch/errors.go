// Package cloudwatch provides CloudWatch storage functionality for vorpalstacks.
package cloudwatch

import "vorpalstacks/internal/store/aws/common"

var (
	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = common.NewStoreError("cloudwatch", "invalid_parameter", common.ErrInvalidInput)

	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = common.NewStoreError("cloudwatch", "resource_not_found", common.ErrNotFound)

	// ErrMetricNotFound is returned when the specified metric does not exist.
	ErrMetricNotFound = common.NewStoreError("cloudwatch", "metric_not_found", common.ErrNotFound)

	// ErrAlarmNotFound is returned when the specified alarm does not exist.
	ErrAlarmNotFound = common.NewStoreError("cloudwatch", "alarm_not_found", common.ErrNotFound)

	// ErrAlarmAlreadyExists is returned when attempting to create an alarm
	// that already exists.
	ErrAlarmAlreadyExists = common.NewStoreError("cloudwatch", "alarm_already_exists", common.ErrAlreadyExists)

	// ErrAlarmLimitExceeded is returned when the alarm limit has been exceeded.
	ErrAlarmLimitExceeded = common.NewStoreError("cloudwatch", "alarm_limit_exceeded", common.ErrInvalidInput)
)
