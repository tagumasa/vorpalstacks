package appsync

import "testing"

// modelAssociationStatuses enumerates the SourceApiAssociationStatus enum of
// the AppSync API model (appsync 2017-07-25): MERGE_SCHEDULED, MERGE_FAILED,
// MERGE_SUCCESS, MERGE_IN_PROGRESS, AUTO_MERGE_SCHEDULE_FAILED,
// DELETION_SCHEDULED, DELETION_IN_PROGRESS, DELETION_FAILED.
var modelAssociationStatuses = map[string]bool{
	"MERGE_SCHEDULED":            true,
	"MERGE_FAILED":               true,
	"MERGE_SUCCESS":              true,
	"MERGE_IN_PROGRESS":          true,
	"AUTO_MERGE_SCHEDULE_FAILED": true,
	"DELETION_SCHEDULED":         true,
	"DELETION_IN_PROGRESS":       true,
	"DELETION_FAILED":            true,
}

// TestAssociationStatusesStayWithinTheModelEnum pins every association-status
// constant the service writes to a member of the SourceApiAssociationStatus
// enum, so a hand-typed status cannot drift outside the model's wire values.
func TestAssociationStatusesStayWithinTheModelEnum(t *testing.T) {
	for _, status := range []string{
		assocStatusMergeScheduled,
		assocStatusMergeInProgress,
		assocStatusMergeSuccess,
		assocStatusDeletionScheduled,
		assocStatusDeletionFailed,
	} {
		if !modelAssociationStatuses[status] {
			t.Errorf("association status %q is not a member of the SourceApiAssociationStatus enum", status)
		}
	}
}
