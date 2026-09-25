package dynamodb

import (
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"google.golang.org/protobuf/proto"
)

// ---------------------------------------------------------------------------
// Transaction cascade — the per-table state sweep a table deletion runs in
// its single transaction (items, index entries, stream records, journal,
// contributor counters, exports, imports, the table record, auto-scaling
// settings, tags) and the delete-time system backup written from the items
// the sweep is about to remove. The bounded prefix-delete machinery the
// sweep builds on stays in store.go, shared with the per-index entry
// sweeps.
// ---------------------------------------------------------------------------

// DeleteTableCascade removes the per-table state a deleted table owns
// within a single transaction: items, GSI/LSI/vector index entries, stream
// records, the PITR journal, contributor counters, the table record
// itself, exports, imports, and tags. The global-table record is not the
// cascade's to sweep: it lives in the account-global storage outside the
// regional transaction, and the membership cleanup a member table's
// deletion owes runs in the service layer after this cascade commits.
// It does NOT check DeletionProtectionEnabled — that is the caller's responsibility.
//
// Backups are deliberately absent from the sweep: a backup is an
// independently ARN-addressed resource removed through its own DeleteBackup
// operation, and the documented contract of DeleteTable names no side
// effect on them ("The DeleteTable operation deletes a table and all of its
// items"). A deleted table's recovery runs through the delete-time system
// backup the caller writes before invoking this cascade (and through any
// on-demand backups, which equally outlive the table).
func (t *DynamoDBTxn) DeleteTableCascade(name string) error {
	if err := t.deleteAllByPrefix(itemBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete items for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(gsiIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete GSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(lsiIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete LSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(vectorIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete vector index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(streamBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete stream records for table %s: %w", name, err)
	}

	// The journal is keyed by the bare table name, so it cannot outlive the
	// table: a same-name successor table must never replay a dead table's
	// mutations. The documented recovery path for a deleted table is the
	// delete-time system backup, not point-in-time restore by name (which
	// requires a live source table).
	if err := t.deleteAllByPrefix(journalBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete journal records for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(contributorBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete contributor counters for table %s: %w", name, err)
	}

	if err := t.deleteExportsForTable(name); err != nil {
		return fmt.Errorf("delete exports for table %s: %w", name, err)
	}

	if err := t.deleteImportsForTable(name); err != nil {
		return fmt.Errorf("delete imports for table %s: %w", name, err)
	}

	tableBucket := t.txn.Bucket(tableBucketName(t.region()))
	if err := tableBucket.Delete([]byte(name)); err != nil {
		return fmt.Errorf("delete table record %s: %w", name, err)
	}

	// The auto-scaling settings record is keyed by the bare table name in
	// its own bucket; without this sweep it outlives the table and a
	// same-name re-created table inherits the dead table's settings.
	autoScalingBucket := t.txn.Bucket(autoScalingBucketName(t.region()))
	if autoScalingBucket != nil {
		if err := autoScalingBucket.Delete([]byte(name)); err != nil {
			return fmt.Errorf("delete auto-scaling settings for table %s: %w", name, err)
		}
	}

	tagMainBucket := t.txn.Bucket(tagMainBucketName(t.region()))
	if tagMainBucket != nil {
		if err := tagMainBucket.Delete([]byte(name)); err != nil {
			return fmt.Errorf("delete tags for table %s: %w", name, err)
		}
	}

	tagIdxBucket := t.txn.Bucket(tagIndexBucketName(t.region()))
	if tagIdxBucket != nil {
		suffix := "\x00" + name
		iter := tagIdxBucket.ScanPrefix(nil)
		defer iter.Close()
		var idxKeysToDelete []string
		for iter.Next() {
			k := string(iter.Key())
			if strings.HasSuffix(k, suffix) {
				idxKeysToDelete = append(idxKeysToDelete, k)
			}
		}
		if err := iter.Error(); err != nil {
			return fmt.Errorf("scan tag index for table %s: %w", name, err)
		}
		for _, k := range idxKeysToDelete {
			if err := tagIdxBucket.Delete([]byte(k)); err != nil {
				return fmt.Errorf("delete tag index entry %s: %w", k, err)
			}
		}
	}

	return nil
}

// CreateDeletedTableSystemBackup writes the table-name$DeletedTableBackup
// system backup inside the deletion transaction: the backup record plus an
// item snapshot scanned from the items the cascade is about to sweep, so
// the snapshot is the table exactly as it stood before deletion. The
// developer guide ("Delete a table with PITR enabled"): "DynamoDB
// automatically creates a backup snapshot called a system backup and
// retains it for 35 days (at no additional cost). You can use the system
// backup to restore the deleted table to the state it was in before
// deletion. All system backups follow a standard naming convention of
// table-name$DeletedTableBackup." The name is the documented convention, a
// label; the identity is a freshly minted id, so a same-name recreate and
// re-delete inside the window leaves BOTH generations retained — each
// addressed by its own ARN, each living its own 35 days. The record enters
// AVAILABLE state directly — the snapshot completes within this
// transaction.
func (t *DynamoDBTxn) CreateDeletedTableSystemBackup(table *Table) error {
	backupName := table.Name + DeletedTableBackupSuffix

	itemBucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := itemBucket.ScanPrefix([]byte(table.Name + KeySep))
	defer iter.Close()
	var items []*Item
	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			return fmt.Errorf("unmarshal item during system backup scan for table %s: %w", table.Name, err)
		}
		items = append(items, itemFromProto(&pbItem))
	}
	if err := iter.Error(); err != nil {
		return fmt.Errorf("scan items for system backup of table %s: %w", table.Name, err)
	}

	now := time.Now().UTC()
	backupId, idErr := tableScopedId(now)
	if idErr != nil {
		return fmt.Errorf("mint system backup id for table %s: %w", table.Name, idErr)
	}
	backup := &Backup{
		BackupName:              backupName,
		BackupArn:               t.tableStore.arnBuilder.Backup(table.Name, backupId),
		SourceTableName:         table.Name,
		SourceTableArn:          table.ARN,
		SourceTableId:           table.TableId,
		BackupStatus:            BackupStatusAvailable,
		BackupType:              BackupTypeSystem,
		BackupCreationDateTime:  now,
		BackupSizeBytes:         table.TableSizeBytes,
		KeySchema:               table.KeySchema,
		AttributeDefinitions:    table.AttributeDefinitions,
		BillingMode:             table.BillingMode,
		ProvisionedThroughput:   table.ProvisionedThroughput,
		GlobalSecondaryIndexes:  table.GlobalSecondaryIndexes,
		LocalSecondaryIndexes:   table.LocalSecondaryIndexes,
		VectorIndexes:           table.VectorIndexes,
		SourceTableCreationTime: table.CreationDateTime,
		SourceTableSizeBytes:    table.TableSizeBytes,
		SourceTableItemCount:    table.ItemCount,
	}

	backupBucket := t.txn.Bucket(backupBucketName(t.region()))
	recordBytes, err := proto.Marshal(BackupToProto(backup))
	if err != nil {
		return fmt.Errorf("marshal system backup record for table %s: %w", table.Name, err)
	}
	if err := backupBucket.Put([]byte(backupId), recordBytes); err != nil {
		return fmt.Errorf("put system backup record for table %s: %w", table.Name, err)
	}
	snapshotBytes, err := proto.Marshal(backupSnapshotToProto(items))
	if err != nil {
		return fmt.Errorf("marshal system backup snapshot for table %s: %w", table.Name, err)
	}
	if err := backupBucket.Put([]byte(snapshotKeyForId(backupId)), snapshotBytes); err != nil {
		return fmt.Errorf("put system backup snapshot for table %s: %w", table.Name, err)
	}
	return nil
}

// deleteJobRecordsForTable removes every record of one job family whose
// TableArn addresses the given table, inside the caller's transaction.
// decode reads the record's TableArn; the match goes through the ARN
// parser (the table name the ARN's resource names), never a substring
// probe — the ARN grammar places a colon before the resource, so a
// suffix probe built from "/table/"+name cannot match a real table ARN.
// A record that fails to decode is removed as well — its carrying table
// is going away, so a record that cannot even be read cannot outlive it.
func (t *DynamoDBTxn) deleteJobRecordsForTable(bucketName, family, tableName string, decode func([]byte) (string, error)) error {
	bucket := t.txn.Bucket(bucketName)
	if bucket == nil {
		return nil
	}
	iter := bucket.ScanPrefix(nil)
	defer iter.Close()
	var keysToDelete []string
	for iter.Next() {
		tableArn, err := decode(iter.Value())
		if err != nil {
			logs.Warn("cascade delete: corrupted "+family+" record, will delete", logs.String("key", string(iter.Key())), logs.Err(err))
			keysToDelete = append(keysToDelete, string(iter.Key()))
			continue
		}
		if svcarn.ParseTableARN(tableArn) == tableName {
			keysToDelete = append(keysToDelete, string(iter.Key()))
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	for _, k := range keysToDelete {
		if err := bucket.Delete([]byte(k)); err != nil {
			return err
		}
	}
	return nil
}

func (t *DynamoDBTxn) deleteExportsForTable(tableName string) error {
	return t.deleteJobRecordsForTable(exportBucketName(t.region()), "export", tableName, func(value []byte) (string, error) {
		var export pb.ExportDescription
		if err := proto.Unmarshal(value, &export); err != nil {
			return "", err
		}
		return export.TableArn, nil
	})
}

func (t *DynamoDBTxn) deleteImportsForTable(tableName string) error {
	return t.deleteJobRecordsForTable(importBucketName(t.region()), "import", tableName, func(value []byte) (string, error) {
		var imp pb.ImportTableDescription
		if err := proto.Unmarshal(value, &imp); err != nil {
			return "", err
		}
		return imp.TableArn, nil
	})
}
