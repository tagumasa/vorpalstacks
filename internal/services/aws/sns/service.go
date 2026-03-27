// Package sns provides SNS (Simple Notification Service) operations for vorpalstacks.
package sns

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// SNSService provides SNS topic and subscription operations.
type SNSService struct {
	storageManager *storage.RegionStorageManager
	sqsStore       sqsstore.SQSStoreInterface
	lambdaInvoker  common.LambdaInvoker
	accountID      string
	httpClient     *http.Client
	stores         sync.Map

	signingKeyOnce sync.Once
	signingKey     *rsa.PrivateKey
	signingCertPEM []byte
}

func (s *SNSService) store(reqCtx *request.RequestContext) (snsstore.SNSStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (snsstore.SNSStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return snsstore.NewSNSStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// NewSNSService creates a new SNS service instance.
// Optional cross-service dependencies (SQS store, Lambda invoker) should be
// injected via setter methods before registering handlers.
func NewSNSService(storageMgr *storage.RegionStorageManager, accountID string) *SNSService {
	return &SNSService{
		storageManager: storageMgr,
		accountID:      accountID,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}
}

// SetSQSStore injects an SQS store for cross-service topic-to-queue delivery.
func (s *SNSService) SetSQSStore(store sqsstore.SQSStoreInterface) {
	s.sqsStore = store
}

// SetLambdaInvoker injects a Lambda invoker for cross-service topic-to-function delivery.
func (s *SNSService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// SetSNSStore pre-populates the regional store cache with an existing SNS store instance.
func (s *SNSService) SetSNSStore(region string, snsStore *snsstore.SNSStore) {
	if snsStore != nil {
		s.stores.Store(region, snsStore)
	}
}

func (s *SNSService) initSigningKey() {
	s.signingKeyOnce.Do(func() {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}
		s.signingKey = privateKey

		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"Vorpalstacks"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(365 * 24 * time.Hour),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return
		}

		s.signingCertPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certDER,
		})
	})
}

func (s *SNSService) RegisterHandlers(d *dispatcher.Dispatcher) {
	s.initSigningKey()
	d.RegisterHandlerForService("sns", "CreateTopic", s.CreateTopic)
	d.RegisterHandlerForService("sns", "DeleteTopic", s.DeleteTopic)
	d.RegisterHandlerForService("sns", "GetTopicAttributes", s.GetTopicAttributes)
	d.RegisterHandlerForService("sns", "SetTopicAttributes", s.SetTopicAttributes)
	d.RegisterHandlerForService("sns", "ListTopics", s.ListTopics)

	d.RegisterHandlerForService("sns", "Subscribe", s.Subscribe)
	d.RegisterHandlerForService("sns", "Unsubscribe", s.Unsubscribe)
	d.RegisterHandlerForService("sns", "ConfirmSubscription", s.ConfirmSubscription)
	d.RegisterHandlerForService("sns", "GetSubscriptionAttributes", s.GetSubscriptionAttributes)
	d.RegisterHandlerForService("sns", "SetSubscriptionAttributes", s.SetSubscriptionAttributes)
	d.RegisterHandlerForService("sns", "ListSubscriptions", s.ListSubscriptions)
	d.RegisterHandlerForService("sns", "ListSubscriptionsByTopic", s.ListSubscriptionsByTopic)

	d.RegisterHandlerForService("sns", "Publish", s.Publish)
	d.RegisterHandlerForService("sns", "PublishBatch", s.PublishBatch)

	d.RegisterHandlerForService("sns", "TagResource", s.TagResource)
	d.RegisterHandlerForService("sns", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("sns", "ListTagsForResource", s.ListTagsForResource)

	d.RegisterHandlerForService("sns", "CreatePlatformApplication", s.CreatePlatformApplication)
	d.RegisterHandlerForService("sns", "DeletePlatformApplication", s.DeletePlatformApplication)
	d.RegisterHandlerForService("sns", "GetPlatformApplicationAttributes", s.GetPlatformApplicationAttributes)
	d.RegisterHandlerForService("sns", "SetPlatformApplicationAttributes", s.SetPlatformApplicationAttributes)
	d.RegisterHandlerForService("sns", "ListPlatformApplications", s.ListPlatformApplications)

	d.RegisterHandlerForService("sns", "CreatePlatformEndpoint", s.CreatePlatformEndpoint)
	d.RegisterHandlerForService("sns", "DeleteEndpoint", s.DeleteEndpoint)
	d.RegisterHandlerForService("sns", "GetEndpointAttributes", s.GetEndpointAttributes)
	d.RegisterHandlerForService("sns", "SetEndpointAttributes", s.SetEndpointAttributes)
	d.RegisterHandlerForService("sns", "ListEndpointsByPlatformApplication", s.ListEndpointsByPlatformApplication)

	d.RegisterHandlerForService("sns", "GetDataProtectionPolicy", s.GetDataProtectionPolicy)
	d.RegisterHandlerForService("sns", "PutDataProtectionPolicy", s.PutDataProtectionPolicy)
	d.RegisterHandlerForService("sns", "AddPermission", s.AddPermission)
	d.RegisterHandlerForService("sns", "RemovePermission", s.RemovePermission)
}

// PublishToTopic publishes a message to an SNS topic.
func (s *SNSService) PublishToTopic(ctx context.Context, accountID, region, topicArn, message string) error {
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}

	store := snsstore.NewSNSStore(storage, accountID, region)

	topic, err := store.GetTopic(topicArn)
	if err != nil {
		return fmt.Errorf("topic not found: %s", topicArn)
	}

	subscriptions, err := store.ListSubscriptionsByTopic(topicArn, storecommon.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}

	if len(subscriptions.Items) > 0 {
		msg := &snsstore.Message{
			MessageId:          fmt.Sprintf("%d", time.Now().UnixNano()),
			TopicArn:           topic.Arn,
			Message:            message,
			PublishedTimestamp: time.Now().UTC(),
			ReceivedTimestamp:  time.Now().UTC(),
		}

		subsCopy := make([]*snsstore.Subscription, len(subscriptions.Items))
		for i, sub := range subscriptions.Items {
			subCopy := *sub
			subsCopy[i] = &subCopy
		}

		go s.deliverToSubscriptions(msg, subsCopy, region)
	}

	return nil
}
