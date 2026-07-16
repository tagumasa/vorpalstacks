// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
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
	backupName := extractBackupNameFromARN(backupArn)
	if backupName == "" {
		return nil, ErrBackupNotFound
	}
	return s.GetByName(backupName)
}

func extractBackupNameFromARN(backupArn string) string {
	parts := strings.Split(backupArn, "/")
	if len(parts) >= 4 {
		return parts[len(parts)-1]
	}
	return ""
}

// GetByName retrieves a backup by its name.
func (s *BackupStore) GetByName(backupName string) (*Backup, error) {
	var backup Backup
	if err := s.BaseStore.Get(backupName, &backup); err != nil {
		return nil, err
	}
	return &backup, nil
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

	if err := s.BaseStore.Put(backupName, backup); err != nil {
		return nil, err
	}

	return backup, nil
}

// Put stores a backup.
func (s *BackupStore) Put(backup *Backup) error {
	return s.BaseStore.Put(backup.BackupName, backup)
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

	filter := func(b *Backup) bool {
		if tableName == "" {
			return true
		}
		return b.SourceTableName == tableName
	}

	result, err := common.List[Backup](s.BaseStore, opts, filter)
	if err != nil {
		return nil, "", err
	}

	backups := make([]*Backup, len(result.Items))
	copy(backups, result.Items)

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
// The "#" prefix ensures common.List skips this key (it only unmarshals
// keys that do not start with "#"), preventing the snapshot JSON array
// from breaking ListBackups which expects Backup structs.
func snapshotKey(backupName string) string {
	return "#" + backupName + "__snapshot"
}

// SaveSnapshot stores all items for a backup as a single JSON blob.
func (s *BackupStore) SaveSnapshot(backupName string, items []*Item) error {
	return s.BaseStore.Put(snapshotKey(backupName), items)
}

// GetSnapshot loads the item snapshot for a backup.
func (s *BackupStore) GetSnapshot(backupName string) ([]*Item, error) {
	var items []*Item
	if err := s.BaseStore.Get(snapshotKey(backupName), &items); err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteSnapshot removes the item snapshot for a backup.
func (s *BackupStore) DeleteSnapshot(backupName string) error {
	return s.BaseStore.Delete(snapshotKey(backupName))
}
