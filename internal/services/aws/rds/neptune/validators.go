package neptune

import (
	"fmt"

	rdssvc "vorpalstacks/internal/services/aws/rds"
)

// validatePort validates that the given port falls within the AWS Neptune
// allowed range (1150-65535). Delegates to the shared RDS validator.
func validatePort(v int) error {
	return rdssvc.ValidatePort(int32(v))
}

// validateBackupRetentionPeriod validates that the backup retention period
// falls within the AWS Neptune allowed range (1-35 days). Delegates to the
// shared RDS validator.
func validateBackupRetentionPeriod(v int) error {
	return rdssvc.ValidateBackupRetentionPeriod(int32(v))
}

// maxSubnetsPerGroup is the AWS RDS/Neptune quota for subnets per DB subnet
// group. See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RDS_Fea_Regions.html
const maxSubnetsPerGroup = 20

// validateSubnetCount returns an error if the number of subnets exceeds the
// AWS Neptune quota of 20 subnets per DB subnet group.
func validateSubnetCount(n int) error {
	if n > maxSubnetsPerGroup {
		return fmt.Errorf("Cannot assign more than %d subnets to a DB subnet group", maxSubnetsPerGroup)
	}
	return nil
}

// validApplyActions are the permitted values for the ApplyAction parameter
// of ApplyPendingMaintenanceAction per AWS RDS API Reference:
// https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_ApplyPendingMaintenanceAction.html
var validApplyActions = map[string]bool{
	"ca-certificate-rotation": true,
	"db-upgrade":              true,
	"hardware-maintenance":    true,
	"os-upgrade":              true,
	"system-update":           true,
}

// validateApplyAction returns an error if the given action is not one of the
// permitted ApplyAction values.
func validateApplyAction(action string) error {
	if !validApplyActions[action] {
		return fmt.Errorf("Invalid ApplyAction: %s. Valid values: ca-certificate-rotation, db-upgrade, hardware-maintenance, os-upgrade, system-update", action)
	}
	return nil
}

// validOptInTypes are the permitted values for the OptInType parameter
// of ApplyPendingMaintenanceAction per AWS RDS API Reference:
// https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_ApplyPendingMaintenanceAction.html
var validOptInTypes = map[string]bool{
	"immediate":        true,
	"next-maintenance": true,
	"undo-opt-in":      true,
}

// validateOptInType returns an error if the given opt-in type is not one of
// the permitted OptInType values.
func validateOptInType(optIn string) error {
	if !validOptInTypes[optIn] {
		return fmt.Errorf("Invalid OptInType: %s. Valid values: immediate, next-maintenance, undo-opt-in", optIn)
	}
	return nil
}
