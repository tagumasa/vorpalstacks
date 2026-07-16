package dynamodb

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// DynamoDBStoreFactory creates and manages DynamoDB stores per region.
type DynamoDBStoreFactory struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	stores         sync.Map
}

// NewDynamoDBStoreFactory creates a new DynamoDB store factory.
func NewDynamoDBStoreFactory(storageManager *storage.RegionStorageManager, accountID string) *DynamoDBStoreFactory {
	return &DynamoDBStoreFactory{
		storageManager: storageManager,
		accountID:      accountID,
	}
}

// GetStore returns the DynamoDB store for the given region.
func (f *DynamoDBStoreFactory) GetStore(region string) (DynamoDBStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&f.stores, region, func() (DynamoDBStoreInterface, error) {
		basicStore, err := f.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}

		tstore, ok := basicStore.(storage.TransactionalStorageWith2PC)
		if !ok {
			return nil, fmt.Errorf("storage does not support 2PC transactions")
		}

		return NewDynamoDBStore(tstore, f.accountID, region), nil
	})
}
