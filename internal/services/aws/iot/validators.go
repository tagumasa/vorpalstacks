package iot

import (
	"regexp"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// thingNamePattern enforces the Smithy ThingName shape constraints:
// pattern ^[a-zA-Z0-9:_-]+$, length 1–128.
var thingNamePattern = regexp.MustCompile(`^[a-zA-Z0-9:_-]{1,128}$`)

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
