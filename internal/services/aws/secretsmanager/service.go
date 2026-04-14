// Package secretsmanager provides AWS Secrets Manager operations for vorpalstacks.
package secretsmanager

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// SecretsManagerService provides AWS Secrets Manager operations.
type SecretsManagerService struct {
	accountID      string
	region         string
	storageManager *storage.RegionStorageManager
	bus            eventbus.Bus
	lambdaInvoker  common.LambdaInvoker
	logger         logs.Logger
	rotChecker     *rotationChecker
	stores         sync.Map // region → secretsmanagerstore.SecretStoreInterface
}

// NewSecretsManagerService creates a new Secrets Manager service instance.
func NewSecretsManagerService(accountID string) *SecretsManagerService {
	return &SecretsManagerService{
		accountID: accountID,
	}
}

// SetRegion sets the default AWS region for this service instance.
func (s *SecretsManagerService) SetRegion(region string) {
	s.region = region
}

// SetStorageManager injects a region storage manager for cross-region operations.
func (s *SecretsManagerService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetEventBus registers the Secrets Manager rotation handler on the event bus.
func (s *SecretsManagerService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// SetLambdaInvoker injects the Lambda invoker for rotation Lambda calls.
func (s *SecretsManagerService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// SetLogger injects a structured logger for rotation diagnostics.
func (s *SecretsManagerService) SetLogger(logger logs.Logger) {
	s.logger = logger
}

// StartRotationChecker launches the background goroutine that triggers
// automatic secret rotation based on NextRotationDate.
func (s *SecretsManagerService) StartRotationChecker(ctx context.Context) {
	if s.storageManager == nil {
		return
	}
	s.rotChecker = newRotationChecker(s, s.logger)
	s.rotChecker.start(ctx)
}

// StopRotationChecker signals the background rotation checker to stop.
func (s *SecretsManagerService) StopRotationChecker() {
	if s.rotChecker != nil {
		s.rotChecker.stop()
	}
}

func (s *SecretsManagerService) store(reqCtx *request.RequestContext) (secretsmanagerstore.SecretStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (secretsmanagerstore.SecretStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return secretsmanagerstore.NewSecretStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// RegisterHandlers registers the Secrets Manager handlers with the dispatcher.
func (s *SecretsManagerService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("secretsmanager", "CreateSecret", s.CreateSecret)
	d.RegisterHandlerForService("secretsmanager", "GetSecretValue", s.GetSecretValue)
	d.RegisterHandlerForService("secretsmanager", "UpdateSecret", s.UpdateSecret)
	d.RegisterHandlerForService("secretsmanager", "DeleteSecret", s.DeleteSecret)
	d.RegisterHandlerForService("secretsmanager", "ListSecrets", s.ListSecrets)
	d.RegisterHandlerForService("secretsmanager", "DescribeSecret", s.DescribeSecret)
	d.RegisterHandlerForService("secretsmanager", "PutSecretValue", s.PutSecretValue)
	d.RegisterHandlerForService("secretsmanager", "TagResource", s.TagResource)
	d.RegisterHandlerForService("secretsmanager", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("secretsmanager", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("secretsmanager", "RestoreSecret", s.RestoreSecret)
	d.RegisterHandlerForService("secretsmanager", "RotateSecret", s.RotateSecret)
	d.RegisterHandlerForService("secretsmanager", "CancelRotateSecret", s.CancelRotateSecret)
	d.RegisterHandlerForService("secretsmanager", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "ValidateResourcePolicy", s.ValidateResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "GetRandomPassword", s.GetRandomPassword)
	d.RegisterHandlerForService("secretsmanager", "ListSecretVersionIds", s.ListSecretVersionIds)
	d.RegisterHandlerForService("secretsmanager", "UpdateSecretVersionStage", s.UpdateSecretVersionStage)
	d.RegisterHandlerForService("secretsmanager", "BatchGetSecretValue", s.BatchGetSecretValue)
	d.RegisterHandlerForService("secretsmanager", "ReplicateSecretToRegions", s.ReplicateSecretToRegions)
	d.RegisterHandlerForService("secretsmanager", "RemoveRegionsFromReplication", s.RemoveRegionsFromReplication)
	d.RegisterHandlerForService("secretsmanager", "StopReplicationToReplica", s.StopReplicationToReplica)
}
