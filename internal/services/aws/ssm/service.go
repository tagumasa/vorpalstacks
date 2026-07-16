// Package ssm provides AWS Systems Manager (SSM) Parameter Store service operations for vorpalstacks.
package ssm

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// SSMService provides AWS Systems Manager Parameter Store operations.
type SSMService struct {
	accountID      string
	kmsEncryptor   common.KMSEncryptor
	stores         sync.Map // region → ssmstore.SSMStoreInterface
	storageManager *storage.RegionStorageManager
}

// NewSSMService creates a new SSM service instance.
func NewSSMService(accountID string) *SSMService {
	return &SSMService{
		accountID: accountID,
	}
}

// NewSSMServiceWithKMS creates a new SSM service instance with KMS support.
func NewSSMServiceWithKMS(accountID string, kmsEncryptor common.KMSEncryptor) *SSMService {
	return &SSMService{
		accountID:    accountID,
		kmsEncryptor: kmsEncryptor,
	}
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *SSMService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// GetStoreForRegion returns the cached SSMStore for the given region,
// creating a new store instance if not already cached.
func (s *SSMService) GetStoreForRegion(region string) (ssmstore.SSMStoreInterface, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(ssmstore.SSMStoreInterface), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("ssm storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := ssmstore.NewStore(st, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(region, store)
	return actual.(ssmstore.SSMStoreInterface), nil
}

func (s *SSMService) store(reqCtx *request.RequestContext) (ssmstore.SSMStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (ssmstore.SSMStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return ssmstore.NewStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

func (s *SSMService) encryptValue(ctx context.Context, plaintext, keyID string) (string, error) {
	if s.kmsEncryptor == nil {
		return plaintext, nil
	}
	return s.kmsEncryptor.EncryptString(ctx, keyID, plaintext)
}

func (s *SSMService) decryptValue(ctx context.Context, ciphertext, keyID string) (string, error) {
	if s.kmsEncryptor == nil {
		return ciphertext, nil
	}
	return s.kmsEncryptor.DecryptString(ctx, keyID, ciphertext)
}

// RegisterHandlers registers all SSM operation handlers with the dispatcher.
func (s *SSMService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("ssm", "PutParameter", s.PutParameter)
	d.RegisterHandlerForService("ssm", "GetParameter", s.GetParameter)
	d.RegisterHandlerForService("ssm", "GetParameters", s.GetParameters)
	d.RegisterHandlerForService("ssm", "GetParametersByPath", s.GetParametersByPath)
	d.RegisterHandlerForService("ssm", "DeleteParameter", s.DeleteParameter)
	d.RegisterHandlerForService("ssm", "DeleteParameters", s.DeleteParameters)
	d.RegisterHandlerForService("ssm", "DescribeParameters", s.DescribeParameters)
	d.RegisterHandlerForService("ssm", "GetParameterHistory", s.GetParameterHistory)
	d.RegisterHandlerForService("ssm", "LabelParameterVersion", s.LabelParameterVersion)
	d.RegisterHandlerForService("ssm", "UnlabelParameterVersion", s.UnlabelParameterVersion)

	d.RegisterHandlerForService("ssm", "AddTagsToResource", s.AddTagsToResource)
	d.RegisterHandlerForService("ssm", "RemoveTagsFromResource", s.RemoveTagsFromResource)
	d.RegisterHandlerForService("ssm", "ListTagsForResource", s.ListTagsForResource)
}
