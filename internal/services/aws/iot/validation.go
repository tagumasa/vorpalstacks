package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ValidateTargetSelection checks that the value is a valid AWS IoT
// job target selection enum.
func ValidateTargetSelection(v string) error {
	switch v {
	case "CONTINUOUS", "SNAPSHOT":
		return nil
	}
	return iotstore.ErrInvalidRequest
}

// ValidateAuthorizerStatus checks that the value is a valid authorizer status.
func ValidateAuthorizerStatus(v string) error {
	switch v {
	case "ACTIVE", "INACTIVE":
		return nil
	}
	return iotstore.ErrInvalidRequest
}
