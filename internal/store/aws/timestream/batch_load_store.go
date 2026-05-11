package timestream

import (
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_timestream"
	"vorpalstacks/internal/store/aws/common"
)

// BatchLoadTaskStore manages batch load tasks in Timestream.
type BatchLoadTaskStore struct {
	*common.BaseStore
	tableStore *TableStore
	createMu   sync.Mutex
	region     string
}

// NewBatchLoadTaskStore creates a new BatchLoadTaskStore.
func NewBatchLoadTaskStore(store storage.BasicStorage, tableStore *TableStore, region string) *BatchLoadTaskStore {
	return &BatchLoadTaskStore{
		BaseStore:  common.NewBaseStore(store.Bucket(batchLoadTaskBucketName(region)), "timestream-batch-load"),
		tableStore: tableStore,
		region:     region,
	}
}

// CreateBatchLoadTask creates a new batch load task.
func (s *BatchLoadTaskStore) CreateBatchLoadTask(taskId, targetDatabaseName, targetTableName string, dataSourceConfig *DataSourceConfiguration, dataModelConfig *DataModelConfiguration, reportConfig *ReportConfiguration, recordVersion int64) (*BatchLoadTaskDescription, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.Exists(taskId) {
		return nil, ErrBatchLoadTaskAlreadyExists
	}

	if _, err := s.tableStore.GetTable(targetDatabaseName, targetTableName); err != nil {
		return nil, ErrTableNotFound
	}

	now := time.Now().UTC()
	resumableUntil := now.Add(7 * 24 * time.Hour)

	task := &BatchLoadTaskDescription{
		TaskId:                  taskId,
		TargetDatabaseName:      targetDatabaseName,
		TargetTableName:         targetTableName,
		TaskStatus:              BatchLoadStatusCreated,
		CreationTime:            now,
		LastUpdatedTime:         now,
		ResumableUntil:          resumableUntil,
		RecordVersion:           recordVersion,
		DataSourceConfiguration: dataSourceConfig,
		DataModelConfiguration:  dataModelConfig,
		ReportConfiguration:     reportConfig,
		ProgressReport: &BatchLoadProgressReport{
			BytesMetered:            0,
			FileFailures:            0,
			ParseFailures:           0,
			RecordIngestionFailures: 0,
			RecordsIngested:         0,
			RecordsProcessed:        0,
		},
	}

	if err := s.PutProto(taskId, BatchLoadTaskDescriptionToProto(task)); err != nil {
		return nil, err
	}

	return task, nil
}

// GetBatchLoadTask retrieves a batch load task by its ID.
func (s *BatchLoadTaskStore) GetBatchLoadTask(taskId string) (*BatchLoadTaskDescription, error) {
	var pbTask pb.BatchLoadTaskDescription
	if err := s.GetProto(taskId, &pbTask); err != nil {
		return nil, ErrBatchLoadTaskNotFound
	}
	return ProtoToBatchLoadTaskDescription(&pbTask), nil
}

// UpdateBatchLoadTaskStatus updates the status of a batch load task.
func (s *BatchLoadTaskStore) UpdateBatchLoadTaskStatus(taskId string, status BatchLoadStatus, errorMessage string) error {
	task, err := s.GetBatchLoadTask(taskId)
	if err != nil {
		return err
	}

	task.TaskStatus = status
	task.LastUpdatedTime = time.Now().UTC()
	if errorMessage != "" {
		task.ErrorMessage = errorMessage
	}

	return s.PutProto(taskId, BatchLoadTaskDescriptionToProto(task))
}

// UpdateBatchLoadTaskProgress updates the progress of a batch load task.
func (s *BatchLoadTaskStore) UpdateBatchLoadTaskProgress(taskId string, progress *BatchLoadProgressReport) error {
	task, err := s.GetBatchLoadTask(taskId)
	if err != nil {
		return err
	}

	task.ProgressReport = progress
	task.LastUpdatedTime = time.Now().UTC()

	return s.PutProto(taskId, BatchLoadTaskDescriptionToProto(task))
}

// DeleteBatchLoadTask deletes a batch load task.
func (s *BatchLoadTaskStore) DeleteBatchLoadTask(taskId string) error {
	if !s.Exists(taskId) {
		return ErrBatchLoadTaskNotFound
	}
	return s.BaseStore.Delete(taskId)
}

// ListBatchLoadTasks lists batch load tasks with optional status filter.
func (s *BatchLoadTaskStore) ListBatchLoadTasks(taskStatus BatchLoadStatus) ([]*BatchLoadTask, error) {
	result, err := common.ListProto[*pb.BatchLoadTaskDescription](s.BaseStore, common.ListOptions{}, func() *pb.BatchLoadTaskDescription { return &pb.BatchLoadTaskDescription{} }, nil)
	if err != nil {
		return nil, err
	}

	var tasks []*BatchLoadTask
	for _, pbDesc := range result.Items {
		desc := ProtoToBatchLoadTaskDescription(pbDesc)
		if taskStatus != "" && desc.TaskStatus != taskStatus {
			continue
		}

		tasks = append(tasks, &BatchLoadTask{
			TaskId:          desc.TaskId,
			DatabaseName:    desc.TargetDatabaseName,
			TableName:       desc.TargetTableName,
			TaskStatus:      desc.TaskStatus,
			CreationTime:    desc.CreationTime,
			LastUpdatedTime: desc.LastUpdatedTime,
			ResumableUntil:  desc.ResumableUntil,
		})
	}
	return tasks, nil
}

// AccountSettingsStore manages account settings for Timestream.
type AccountSettingsStore struct {
	*common.BaseStore
	accountID string
	region    string
}

// NewAccountSettingsStore creates a new AccountSettingsStore.
func NewAccountSettingsStore(store storage.BasicStorage, accountID, region string) *AccountSettingsStore {
	return &AccountSettingsStore{
		BaseStore: common.NewBaseStore(store.Bucket(accountSettingsBucketName(region)), "timestream-account-settings"),
		accountID: accountID,
		region:    region,
	}
}

// GetAccountSettings retrieves account settings.
func (s *AccountSettingsStore) GetAccountSettings() (*AccountSettings, error) {
	key := s.accountID
	var pbSettings pb.AccountSettings
	if err := s.GetProto(key, &pbSettings); err != nil {
		return &AccountSettings{
			MaxQueryTCU:      4,
			QueryPricingMode: QueryPricingModeBytesScanned,
			QueryComputeType: QueryComputeTypeOnDemand,
		}, nil
	}
	return ProtoToAccountSettings(&pbSettings), nil
}

// UpdateAccountSettings updates account settings.
func (s *AccountSettingsStore) UpdateAccountSettings(maxQueryTCU *int64, queryPricingMode, queryComputeType, encryptionConfiguration string) (*AccountSettings, error) {
	settings, err := s.GetAccountSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get account settings: %w", err)
	}

	if maxQueryTCU != nil {
		settings.MaxQueryTCU = *maxQueryTCU
	}
	if queryPricingMode != "" {
		settings.QueryPricingMode = queryPricingMode
	}
	if queryComputeType != "" {
		settings.QueryComputeType = queryComputeType
	}
	if encryptionConfiguration != "" {
		settings.EncryptionConfiguration = encryptionConfiguration
	}

	key := s.accountID
	if err := s.PutProto(key, AccountSettingsToProto(settings)); err != nil {
		return nil, err
	}

	return settings, nil
}
