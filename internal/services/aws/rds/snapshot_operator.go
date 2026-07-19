package rds

// SnapshotOperator abstracts the row-level snapshot lifecycle for
// engines that support it (currently vmysql). The admin handler and the
// Neptune SDK handler both reference this interface so that manual and
// final snapshots capture actual row data — not just metadata — and
// restores can reconstruct user tables.
//
// When nil, snapshots store metadata only (legacy behaviour).
type SnapshotOperator interface {
	// SnapshotData copies every database, table, schema, row, and
	// index from the source instance into a snapshot store keyed by
	// snapshotID.
	SnapshotData(instanceID, snapshotID string) error

	// RestoreData copies every database, table, schema, row, and
	// index from the snapshot store into the destination instance.
	// Existing rows with the same PK are left untouched (idempotent).
	RestoreData(snapshotID, instanceID string) error

	// DeleteSnapshotData removes all row data associated with the
	// snapshotID. Called when a DBSnapshot is deleted so row data
	// does not leak indefinitely.
	DeleteSnapshotData(snapshotID string) error
}
