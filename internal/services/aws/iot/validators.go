package iot

import (
	"regexp"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// thingNamePattern enforces the Smithy ThingName shape constraints:
// pattern ^[a-zA-Z0-9:_-]+$, length 1–128.
var thingNamePattern = regexp.MustCompile(`^[a-zA-Z0-9:_-]{1,128}$`)

// MaxTemplateBodyLength is the Smithy TemplateBody shape's documented
// maximum length (provisioning template bodies on RegisterThing,
// StartThingRegistrationTask and the provisioning-template family).
const MaxTemplateBodyLength = 10240

// Bounds and patterns documented on the API_StartThingRegistrationTask
// members (inputFileBucket, inputFileKey, roleArn). templateBody shares the
// TemplateBody shape bound above.
const (
	MinInputFileBucketLength = 3
	MaxInputFileBucketLength = 256
	MinInputFileKeyLength    = 1
	MaxInputFileKeyLength    = 1024
	MinRoleArnLength         = 20
	MaxRoleArnLength         = 2048
)

var (
	inputFileBucketPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	inputFileKeyPattern    = regexp.MustCompile(`^[a-zA-Z0-9!_.*'()/-]+$`)
)

// Fleet metric period bounds: the Smithy FleetMetricPeriod range trait
// (60–86400 seconds); AWS additionally requires the period to be a
// multiple of 60.
const (
	MinFleetMetricPeriod = 60
	MaxFleetMetricPeriod = 86400
)

// Role-alias credential duration bounds: the Smithy
// CredentialDurationSeconds range trait (900–43200 seconds).
const (
	MinRoleAliasCredentialDuration int64 = 900
	MaxRoleAliasCredentialDuration int64 = 43200
)

// Fleet-metric member bounds from the Smithy shapes FleetMetricName
// (1–128), FleetMetricDescription (0–1024), IndexName (1–128) and
// AggregationTypeValue (1–12).
const (
	MinFleetMetricNameLength        = 1
	MaxFleetMetricNameLength        = 128
	MaxFleetMetricDescriptionLength = 1024
	MinFleetMetricIndexNameLength   = 1
	MaxFleetMetricIndexNameLength   = 128
	MaxAggregationValueLength       = 12
)

var (
	fleetMetricNamePattern        = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	fleetMetricDescriptionPattern = regexp.MustCompile(`^[[:graph:]\x20]*$`)
	fleetMetricIndexNamePattern   = regexp.MustCompile(`^[a-zA-Z0-9:_-]+$`)
	aggregationValuePattern       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

// fleetMetricUnits is the FleetMetricUnit enum's value set.
var fleetMetricUnits = map[string]struct{}{
	"Bits": {}, "Bits/Second": {}, "Bytes": {}, "Bytes/Second": {},
	"Count": {}, "Count/Second": {}, "Gigabits": {}, "Gigabits/Second": {},
	"Gigabytes": {}, "Gigabytes/Second": {}, "Kilobits": {},
	"Kilobits/Second": {}, "Kilobytes": {}, "Kilobytes/Second": {},
	"Megabits": {}, "Megabits/Second": {}, "Megabytes": {},
	"Megabytes/Second": {}, "Microseconds": {}, "Milliseconds": {},
	"None": {}, "Percent": {}, "Seconds": {}, "Terabits": {},
	"Terabits/Second": {}, "Terabytes": {}, "Terabytes/Second": {},
}

// isValidAggregationTypeName checks the AggregationTypeName enum.
func isValidAggregationTypeName(name string) bool {
	switch name {
	case "Statistics", "Percentiles", "Cardinality":
		return true
	}
	return false
}

// ValidateThingName checks that the name matches the AWS IoT ThingName
// shape (^[a-zA-Z0-9:_-]{1,128}$).
func ValidateThingName(name string) error {
	if !thingNamePattern.MatchString(name) {
		return iotstore.ErrValidation
	}
	return nil
}

// IsValidCertUpdateStatus checks whether the given status is one that
// a caller may set via UpdateCertificate.  AWS only allows ACTIVE,
// INACTIVE, and REVOKED as user-settable statuses.  Other statuses
// (PENDING_ACTIVATION, PENDING_TRANSFER, REGISTER_INACTIVE) are set
// internally by the system.
func IsValidCertUpdateStatus(status string) bool {
	switch status {
	case "ACTIVE", "INACTIVE", "REVOKED":
		return true
	}
	return false
}

// IsValidCertStateTransition enforces the AWS IoT certificate status
// state machine.  The key restriction is that a REVOKED certificate
// cannot be transitioned to any other status.
func IsValidCertStateTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case "REVOKED":
		return false
	default:
		return true
	}
}
