package dynamodb

import (
	"time"

	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// systemBackupSweepInterval controls how often the system-backup pruner
// runs. A deleted table's system backup only leaves the store through this
// sweep (or an explicit DeleteBackup), so without it every recovery-enabled
// deletion would accumulate for the lifetime of the process.
const systemBackupSweepInterval = time.Minute

// deletedTableBackupRetentionDays is the retention the developer guide
// documents for the delete-time system backup: "DynamoDB automatically
// creates a backup snapshot called a system backup and retains it for
// 35 days (at no additional cost)."
const deletedTableBackupRetentionDays = 35

// ensureSystemBackupSweeper starts the background pruner that keeps
// delete-time system backups inside their retention window. On-demand
// backups are outside the sweep: they leave the store only through
// DeleteBackup.
func (s *DynamoDBService) ensureSystemBackupSweeper() {
	s.startIntervalSweeper(&s.systemBackupSweepOnce, systemBackupSweepInterval, "dynamodb system-backup sweep", s.sweepStoreSystemBackups)
}

// sweepStoreSystemBackups deletes the expired system backups of one
// regional store, record and item snapshot together (Backups().Delete
// removes both).
func (s *DynamoDBService) sweepStoreSystemBackups(store dbstore.DynamoDBStoreInterface) {
	cutoff := time.Now().Add(-deletedTableBackupRetentionDays * 24 * time.Hour)
	marker := ""
	for {
		backups, next, err := store.Backups().List(marker, listBackupsMinFetchSize, "")
		if err != nil {
			logs.Error("Failed to list backups for system-backup sweep", logs.Err(err))
			return
		}
		for _, b := range backups {
			if b.BackupType != dbstore.BackupTypeSystem {
				continue
			}
			if b.BackupCreationDateTime.After(cutoff) {
				continue
			}
			if err := store.Backups().Delete(b.BackupArn); err != nil {
				logs.Error("Failed to delete expired system backup",
					logs.String("backup", b.BackupName), logs.Err(err))
			}
		}
		if next == "" {
			return
		}
		marker = next
	}
}
