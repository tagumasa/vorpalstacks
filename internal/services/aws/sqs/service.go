// Package sqs provides SQS service operations for vorpalstacks.
package sqs

import (
	"fmt"
	"sync"

	appconfig "vorpalstacks/internal/config"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// SQSService provides SQS operations for managing queues and messages.
type SQSService struct {
	storage        storage.BasicStorage
	storageManager *storage.RegionStorageManager
	accountID      string
	stores         sync.Map
}

// NewSQSService creates a new SQS service instance.
func NewSQSService(store storage.BasicStorage, accountID string) *SQSService {
	return &SQSService{
		storage:   store,
		accountID: accountID,
	}
}

// SetStorageManager injects the region storage manager for admin console access.
func (s *SQSService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// GetStoreForRegion returns the cached SQS store for the given region,
// creating one if not already cached. Used by the admin console to ensure a
// single store instance per region.
func (s *SQSService) GetStoreForRegion(region string) (sqsstore.SQSStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (sqsstore.SQSStoreInterface, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("storage manager not set")
		}
		rs, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return sqsstore.NewSQSStore(rs, s.accountID, region, appconfig.BaseURL()), nil
	})
}

// NewSQSServiceWithStore creates a new SQS service instance with a pre-configured store.
func NewSQSServiceWithStore(store *sqsstore.SQSStore) *SQSService {
	return &SQSService{
		storage:   store.Storage(),
		accountID: store.GetAccountID(),
	}
}

func (s *SQSService) store(reqCtx *request.RequestContext) (sqsstore.SQSStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (sqsstore.SQSStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return sqsstore.NewSQSStore(storage, s.accountID, reqCtx.GetRegion(), appconfig.BaseURL()), nil
	})
}

// RegisterHandlers registers all SQS operation handlers with the dispatcher.
func (s *SQSService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("sqs", "CreateQueue", s.CreateQueue)
	d.RegisterHandlerForService("sqs", "DeleteQueue", s.DeleteQueue)
	d.RegisterHandlerForService("sqs", "GetQueueUrl", s.GetQueueUrl)
	d.RegisterHandlerForService("sqs", "ListQueues", s.ListQueues)
	d.RegisterHandlerForService("sqs", "GetQueueAttributes", s.GetQueueAttributes)
	d.RegisterHandlerForService("sqs", "SetQueueAttributes", s.SetQueueAttributes)
	d.RegisterHandlerForService("sqs", "PurgeQueue", s.PurgeQueue)
	d.RegisterHandlerForService("sqs", "ListDeadLetterSourceQueues", s.ListDeadLetterSourceQueues)
	d.RegisterHandlerForService("sqs", "StartMessageMoveTask", s.StartMessageMoveTask)
	d.RegisterHandlerForService("sqs", "CancelMessageMoveTask", s.CancelMessageMoveTask)
	d.RegisterHandlerForService("sqs", "ListMessageMoveTasks", s.ListMessageMoveTasks)

	d.RegisterHandlerForService("sqs", "SendMessage", s.SendMessage)
	d.RegisterHandlerForService("sqs", "SendMessageBatch", s.SendMessageBatch)
	d.RegisterHandlerForService("sqs", "ReceiveMessage", s.ReceiveMessage)
	d.RegisterHandlerForService("sqs", "DeleteMessage", s.DeleteMessage)
	d.RegisterHandlerForService("sqs", "DeleteMessageBatch", s.DeleteMessageBatch)
	d.RegisterHandlerForService("sqs", "ChangeMessageVisibility", s.ChangeMessageVisibility)
	d.RegisterHandlerForService("sqs", "ChangeMessageVisibilityBatch", s.ChangeMessageVisibilityBatch)

	d.RegisterHandlerForService("sqs", "AddPermission", s.AddPermission)
	d.RegisterHandlerForService("sqs", "RemovePermission", s.RemovePermission)

	d.RegisterHandlerForService("sqs", "TagQueue", s.TagQueue)
	d.RegisterHandlerForService("sqs", "UntagQueue", s.UntagQueue)
	d.RegisterHandlerForService("sqs", "ListQueueTags", s.ListQueueTags)
}
