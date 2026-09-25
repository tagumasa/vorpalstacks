// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func backupBucketName(region string) string {
	return "dynamodb_backups-" + region
}

// BackupStore manages DynamoDB table backups in persistent storage.
type BackupStore struct {
	*common.BaseStore
	arnBuilder *svcarn.DynamoDBBuilder
}

// NewBackupStore creates a new backup store for DynamoDB backups.
func NewBackupStore(store storage.BasicStorage, accountId, region string) *BackupStore {
	return &BackupStore{
		BaseStore:  common.NewBaseStore(store.Bucket(backupBucketName(region)), "dynamodb_backups"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).DynamoDB(),
	}
}

// Get retrieves a backup by its ARN. The ARN's identity segment is the
// storage key, so the lookup needs no name resolution.
func (s *BackupStore) Get(backupArn string) (*Backup, error) {
	id := svcarn.ExtractBackupIdFromARN(backupArn)
	if id == "" {
		return nil, ErrBackupNotFound
	}
	var pbBackup pb.Backup
	if err := s.BaseStore.GetProto(id, &pbBackup); err != nil {
		return nil, err
	}
	return ProtoToBackup(&pbBackup), nil
}

// Create creates a new backup record for a DynamoDB table. The record's
// identity is a freshly minted id — the ARN's final segment and the
// storage key — so two backups may share a name: the documented
// CreateBackup error set carries no duplicate-name rejection, and the
// name is a label, never the key. The record is born CREATING: the
// snapshot it describes is taken after the record exists, and the
// creating operation writes AVAILABLE once the snapshot lands — a
// concurrent describe or restore must never see a snapshot-less backup
// presented as restorable. (The deleted-table system backup inside one
// transaction is the exception: its record and snapshot commit together,
// so it is born AVAILABLE.)
func (s *BackupStore) Create(backupName, tableName, tableArn string, tableSize int64) (*Backup, error) {
	now := time.Now().UTC()
	id, err := tableScopedId(now)
	if err != nil {
		return nil, err
	}
	backupArn := s.arnBuilder.Backup(tableName, id)

	backup := &Backup{
		BackupName:             backupName,
		BackupArn:              backupArn,
		SourceTableName:        tableName,
		SourceTableArn:         tableArn,
		BackupStatus:           BackupStatusCreating,
		BackupType:             BackupTypeUser,
		BackupCreationDateTime: now,
		BackupSizeBytes:        tableSize,
	}

	if err := s.BaseStore.PutProto(id, BackupToProto(backup)); err != nil {
		return nil, err
	}

	return backup, nil
}

// Put stores a backup under its own identity segment, taken from the
// record's ARN.
func (s *BackupStore) Put(backup *Backup) error {
	return s.BaseStore.PutProto(svcarn.ExtractBackupIdFromARN(backup.BackupArn), BackupToProto(backup))
}

// Delete deletes a backup by ARN and removes any associated item snapshot.
// The snapshot goes first: a snapshot delete that fails leaves the record in
// place so the deletion is retried as a whole, rather than deleting the
// record and orphaning snapshot bytes no later path can reach.
func (s *BackupStore) Delete(backupArn string) error {
	if err := s.DeleteSnapshot(backupArn); err != nil {
		return fmt.Errorf("failed to delete backup snapshot: %w", err)
	}
	return s.BaseStore.Delete(svcarn.ExtractBackupIdFromARN(backupArn))
}

// List lists backups with optional table name filter.
func (s *BackupStore) List(marker string, limit int, tableName string) ([]*Backup, string, error) {
	filter := func(pbBackup *pb.Backup) bool {
		if tableName == "" {
			return true
		}
		return pbBackup.SourceTableName == tableName
	}
	return listProtoConverted(s.BaseStore, marker, limit, func() *pb.Backup { return &pb.Backup{} }, ProtoToBackup, filter)
}

// snapshotKeyForId returns the storage key for the item snapshot of the
// backup identified by id. The "#" prefix ensures common.ListProto skips
// this key (it only unmarshals keys that do not start with "#"),
// preventing the snapshot record from breaking ListBackups which expects
// Backup records.
func snapshotKeyForId(id string) string {
	return "#" + id + "__snapshot"
}

// snapshotKey derives the snapshot key from a backup ARN's identity
// segment.
func snapshotKey(backupArn string) string {
	return snapshotKeyForId(svcarn.ExtractBackupIdFromARN(backupArn))
}

// SaveSnapshot stores all items for a backup as a single protobuf record.
func (s *BackupStore) SaveSnapshot(backupArn string, items []*Item) error {
	return s.BaseStore.PutProto(snapshotKey(backupArn), backupSnapshotToProto(items))
}

// GetSnapshot loads the item snapshot for a backup.
func (s *BackupStore) GetSnapshot(backupArn string) ([]*Item, error) {
	var pbSnapshot pb.BackupSnapshot
	if err := s.BaseStore.GetProto(snapshotKey(backupArn), &pbSnapshot); err != nil {
		return nil, err
	}
	return protoToBackupSnapshot(&pbSnapshot), nil
}

// DeleteSnapshot removes the item snapshot for a backup.
func (s *BackupStore) DeleteSnapshot(backupArn string) error {
	return s.BaseStore.Delete(snapshotKey(backupArn))
}
