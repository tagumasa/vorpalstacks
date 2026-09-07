// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
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

// Get retrieves a backup by its ARN.
func (s *BackupStore) Get(backupArn string) (*Backup, error) {
	backupName := svcarn.ExtractBackupNameFromARN(backupArn)
	if backupName == "" {
		return nil, ErrBackupNotFound
	}
	return s.GetByName(backupName)
}

// GetByName retrieves a backup by its name.
func (s *BackupStore) GetByName(backupName string) (*Backup, error) {
	var pbBackup pb.Backup
	if err := s.BaseStore.GetProto(backupName, &pbBackup); err != nil {
		return nil, err
	}
	return ProtoToBackup(&pbBackup), nil
}

// Create creates a new backup for a DynamoDB table.
func (s *BackupStore) Create(backupName, tableName, tableArn string, tableSize int64) (*Backup, error) {
	if s.Exists(backupName) {
		return nil, ErrBackupAlreadyExists
	}

	now := time.Now().UTC()
	backupArn := s.arnBuilder.Backup(tableName, backupName)

	backup := &Backup{
		BackupName:             backupName,
		BackupArn:              backupArn,
		SourceTableName:        tableName,
		SourceTableArn:         tableArn,
		BackupStatus:           BackupStatusAvailable,
		BackupType:             BackupTypeUser,
		BackupCreationDateTime: now,
		BackupSizeBytes:        tableSize,
	}

	if err := s.BaseStore.PutProto(backupName, BackupToProto(backup)); err != nil {
		return nil, err
	}

	return backup, nil
}

// Put stores a backup.
func (s *BackupStore) Put(backup *Backup) error {
	return s.BaseStore.PutProto(backup.BackupName, BackupToProto(backup))
}

// Delete deletes a backup by name and removes any associated item snapshot.
func (s *BackupStore) Delete(backupName string) error {
	_ = s.DeleteSnapshot(backupName)
	return s.BaseStore.Delete(backupName)
}

// Exists checks if a backup exists.
func (s *BackupStore) Exists(backupName string) bool {
	return s.BaseStore.Exists(backupName)
}

// List lists backups with optional table name filter.
func (s *BackupStore) List(marker string, limit int, tableName string) ([]*Backup, string, error) {
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: limit,
	}

	result, err := common.ListProto[*pb.Backup](s.BaseStore, opts, func() *pb.Backup { return &pb.Backup{} }, func(pbBackup *pb.Backup) bool {
		if tableName == "" {
			return true
		}
		return pbBackup.SourceTableName == tableName
	})
	if err != nil {
		return nil, "", err
	}

	backups := make([]*Backup, len(result.Items))
	for i, pbBackup := range result.Items {
		backups[i] = ProtoToBackup(pbBackup)
	}

	if !result.IsTruncated {
		return backups, "", nil
	}
	return backups, result.NextMarker, nil
}

// ARNBuilder returns the ARN builder for DynamoDB.
func (s *BackupStore) ARNBuilder() *svcarn.DynamoDBBuilder {
	return s.arnBuilder
}

// snapshotKey returns the storage key for the item snapshot of a backup.
// The "#" prefix ensures common.ListProto skips this key (it only unmarshals
// keys that do not start with "#"), preventing the snapshot record from
// breaking ListBackups which expects Backup records.
func snapshotKey(backupName string) string {
	return "#" + backupName + "__snapshot"
}

// SaveSnapshot stores all items for a backup as a single protobuf record.
func (s *BackupStore) SaveSnapshot(backupName string, items []*Item) error {
	return s.BaseStore.PutProto(snapshotKey(backupName), backupSnapshotToProto(items))
}

// GetSnapshot loads the item snapshot for a backup.
func (s *BackupStore) GetSnapshot(backupName string) ([]*Item, error) {
	var pbSnapshot pb.BackupSnapshot
	if err := s.BaseStore.GetProto(snapshotKey(backupName), &pbSnapshot); err != nil {
		return nil, err
	}
	return protoToBackupSnapshot(&pbSnapshot), nil
}

// DeleteSnapshot removes the item snapshot for a backup.
func (s *BackupStore) DeleteSnapshot(backupName string) error {
	return s.BaseStore.Delete(snapshotKey(backupName))
}
