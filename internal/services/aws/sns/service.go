// Package sns provides SNS (Simple Notification Service) operations for vorpalstacks.
package sns

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// SNSService provides SNS topic and subscription operations.
type SNSService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	defaultRegion  string
	httpClient     *http.Client
	bus            eventbus.Bus
	stores         sync.Map
	deliveryWg     sync.WaitGroup

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
// Cross-service delivery is routed through the event bus, which must be
// injected via SetEventBus before registering handlers.
func NewSNSService(storageMgr *storage.RegionStorageManager, accountID, region string) *SNSService {
	return &SNSService{
		storageManager: storageMgr,
		accountID:      accountID,
		defaultRegion:  region,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}
}

// SetSNSStore pre-populates the regional store cache with an existing SNS store instance.
func (s *SNSService) SetSNSStore(region string, snsStore *snsstore.SNSStore) {
	if snsStore != nil {
		s.stores.Store(region, snsStore)
	}
}

// SetEventBus injects the event bus and registers the SNS delivery handler.
// When the bus is set, Publish() routes delivery through the bus instead of
// spawning goroutines directly.
func (s *SNSService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	_, _ = eventbus.SubscribeTyped[*eventbus.SNSDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync())
}

// Close waits for all in-flight delivery goroutines to complete.
func (s *SNSService) Close() {
	s.deliveryWg.Wait()
}

// deliverAsync spawns a tracked goroutine to deliver a message to subscriptions.
func (s *SNSService) deliverAsync(msg *snsstore.Message, subs []*snsstore.Subscription, region string) {
	s.deliveryWg.Add(1)
	go func() {
		defer s.deliveryWg.Done()
		defer func() { resilience.RecoverPanic("SNS deliverAsync") }()
		s.deliverToSubscriptions(msg, subs, region)
	}()
}

func (s *SNSService) handleBusDelivery(ctx context.Context, evt *eventbus.SNSDeliveryEvent) eventbus.HandlerResult {
	store, err := s.getSNSStoreByRegion(evt.Region)
	if err != nil {
		logs.Warn("SNS bus delivery: failed to get store", logs.String("region", evt.Region), logs.String("error", err.Error()))
		return eventbus.HandlerResult{Error: err}
	}

	subscriptions, err := store.ListSubscriptionsByTopic(evt.TopicARN, storecommon.ListOptions{})
	if err != nil || len(subscriptions.Items) == 0 {
		return eventbus.HandlerResult{}
	}

	subsCopy := make([]*snsstore.Subscription, len(subscriptions.Items))
	for i, sub := range subscriptions.Items {
		subCopy := *sub
		subsCopy[i] = &subCopy
	}

	msg := &snsstore.Message{
		MessageId:          evt.MessageID,
		TopicArn:           evt.TopicARN,
		Message:            evt.Message,
		Subject:            evt.Subject,
		MessageGroupId:     evt.MessageGroupId,
		PublishedTimestamp: evt.EventTimestamp(),
		ReceivedTimestamp:  evt.EventTimestamp(),
		MessageAttributes:  make(map[string]*snsstore.MessageAttribute),
	}
	// Deserialise message attributes from raw JSON transport format.
	for k, raw := range evt.MessageAttributes {
		attr := &snsstore.MessageAttribute{}
		if json.Unmarshal(raw, attr) == nil {
			msg.MessageAttributes[k] = attr
		}
	}

	s.deliverToSubscriptions(msg, subsCopy, evt.Region)
	return eventbus.HandlerResult{}
}

func (s *SNSService) getSNSStoreByRegion(region string) (snsstore.SNSStoreInterface, error) {
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(snsstore.SNSStoreInterface); ok {
			return typed, nil
		}
	}
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := snsstore.NewSNSStore(storage, s.accountID, region)
	if actual, loaded := s.stores.LoadOrStore(region, store); loaded {
		return actual.(snsstore.SNSStoreInterface), nil
	}
	return store, nil
}

func (s *SNSService) initSigningKey() {
	s.signingKeyOnce.Do(func() {
		rs, err := s.storageManager.GetStorage(s.defaultRegion)
		if err != nil {
			return
		}
		bucket := rs.Bucket("sns-signing")

		if keyPEM, err := bucket.Get([]byte("signing_key")); err == nil && keyPEM != nil {
			if key, err := parseSigningKeyFromPEM(string(keyPEM)); err == nil {
				s.signingKey = key
				if certPEM, err := bucket.Get([]byte("signing_cert")); err == nil && certPEM != nil {
					s.signingCertPEM = certPEM
					return
				}
			}
		}

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

		keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
		certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
		s.signingCertPEM = certBytes

		if err := bucket.Put([]byte("signing_key"), keyBytes); err != nil {
			logs.Warn("failed to persist SNS signing key; key regenerated on restart will invalidate existing message signatures", logs.Err(err))
		}
		if err := bucket.Put([]byte("signing_cert"), certBytes); err != nil {
			logs.Warn("failed to persist SNS signing certificate; certificate regenerated on restart will invalidate existing message signatures", logs.Err(err))
		}
	})
}

func parseSigningKeyFromPEM(pemData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// RegisterHandlers registers all SNS operation handlers with the request dispatcher.
func (s *SNSService) RegisterHandlers(d handler.Registrar) {
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

// PublishToTopic publishes a message to an SNS topic and delivers it to all subscriptions.
func (s *SNSService) PublishToTopic(ctx context.Context, accountID, region, topicArn, message string) error {
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}

	store, _ := storecommon.GetOrCreateStoreE(&s.stores, region, func() (snsstore.SNSStoreInterface, error) {
		return snsstore.NewSNSStore(storage, accountID, region), nil
	})

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

		if s.bus != nil {
			snsEvt := &eventbus.SNSDeliveryEvent{
				TopicARN:  topic.Arn,
				MessageID: msg.MessageId,
				Message:   message,
			}
			snsEvt.Region = region
			if err := s.bus.Publish(context.Background(), snsEvt); err != nil {
				logs.Warn("Failed to publish SNS event", logs.Err(err))
			}
		} else {
			subsCopy := make([]*snsstore.Subscription, len(subscriptions.Items))
			for i, sub := range subscriptions.Items {
				subCopy := *sub
				subsCopy[i] = &subCopy
			}

			s.deliverAsync(msg, subsCopy, region)
		}
	}

	return nil
}
