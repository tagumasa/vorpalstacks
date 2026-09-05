// Package sns provides SNS (Simple Notification Service) operations for vorpalstacks.
package sns

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	awserrors "vorpalstacks/internal/common/errors"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"

	"github.com/google/uuid"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// SNSService provides SNS topic and subscription operations.
type SNSService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	defaultRegion  string
	httpClient     *http.Client
	bus            eventbus.ServiceBus
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
func (s *SNSService) SetEventBus(bus eventbus.ServiceBus) error {
	s.bus = bus
	if _, err := eventbus.SubscribeTyped[*eventbus.SNSDeliveryEvent](bus, s.handleBusDelivery, eventbus.WithAsync()); err != nil {
		return fmt.Errorf("sns: subscribe SNSDeliveryEvent: %w", err)
	}
	return nil
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
		defer func() {
			if r := recover(); r != nil {
				logs.Error("SNS delivery panicked",
					logs.String("messageId", msg.MessageId),
					logs.String("topicArn", msg.TopicArn),
					logs.Any("panic", r))
			}
		}()
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
		MessageStructure:   evt.MessageStructure,
		MessageGroupId:     evt.MessageGroupId,
		PublishedTimestamp: evt.EventTimestamp(),
		ReceivedTimestamp:  evt.EventTimestamp(),
		MessageAttributes:  make(map[string]*snsstore.MessageAttribute),
	}
	// Deserialise message attributes from raw JSON transport format.
	for k, raw := range evt.MessageAttributes {
		attr := &snsstore.MessageAttribute{}
		if err := json.Unmarshal(raw, attr); err != nil {
			logs.Warn("SNS bus delivery: failed to deserialise message attribute, skipping entry",
				logs.String("topicArn", evt.TopicARN),
				logs.String("messageId", evt.MessageID),
				logs.String("attrKey", k),
				logs.Err(err))
			continue
		}
		msg.MessageAttributes[k] = attr
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
			if parsedKey, err := vcrypto.ParsePrivateKeyPEM(keyPEM); err == nil {
				if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
					s.signingKey = rsaKey
					if certPEM, err := bucket.Get([]byte("signing_cert")); err == nil && certPEM != nil {
						s.signingCertPEM = certPEM
						return
					}
				}
			}
		}

		privateKey, err := vcrypto.GenerateRSAKey(2048)
		if err != nil {
			return
		}
		s.signingKey = privateKey

		serial, err := vcrypto.GenerateSerialNumber()
		if err != nil {
			return
		}
		template := x509.Certificate{
			SerialNumber: serial,
			Subject: pkix.Name{
				Organization: []string{"Vorpalstacks"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(365 * 24 * time.Hour),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		certDER, err := vcrypto.CreateCertificate(&template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return
		}

		keyPEM, err := vcrypto.EncodePrivateKeyPEM(privateKey)
		if err != nil {
			return
		}
		certBytes := []byte(vcrypto.EncodeCertificatePEM(certDER))
		s.signingCertPEM = certBytes

		if err := bucket.Put([]byte("signing_key"), []byte(keyPEM)); err != nil {
			logs.Warn("failed to persist SNS signing key; key regenerated on restart will invalidate existing message signatures", logs.Err(err))
		}
		if err := bucket.Put([]byte("signing_cert"), certBytes); err != nil {
			logs.Warn("failed to persist SNS signing certificate; certificate regenerated on restart will invalidate existing message signatures", logs.Err(err))
		}
	})
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

// PublishToTopic publishes a message to an SNS topic and delivers it to all
// subscriptions. Returns the generated message ID.
func (s *SNSService) PublishToTopic(ctx context.Context, accountID, region, topicArn, message, subject string, messageAttributes map[string]string) (string, error) {
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return "", fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}

	store, _ := storecommon.GetOrCreateStoreE(&s.stores, region, func() (snsstore.SNSStoreInterface, error) {
		return snsstore.NewSNSStore(storage, accountID, region), nil
	})

	topic, err := store.GetTopic(topicArn)
	if err != nil {
		return "", awserrors.NewAWSError("NotFound", fmt.Sprintf("topic not found: %s", topicArn), http.StatusNotFound)
	}

	subscriptions, err := store.ListSubscriptionsByTopic(topicArn, storecommon.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list subscriptions: %w", err)
	}

	messageID := uuid.New().String()

	if len(subscriptions.Items) > 0 {
		msg := &snsstore.Message{
			MessageId:          messageID,
			TopicArn:           topic.Arn,
			Subject:            subject,
			Message:            message,
			PublishedTimestamp: time.Now().UTC(),
			ReceivedTimestamp:  time.Now().UTC(),
		}
		if len(messageAttributes) > 0 {
			msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(messageAttributes))
			for k, v := range messageAttributes {
				msg.MessageAttributes[k] = &snsstore.MessageAttribute{Type: "String", StringValue: v}
			}
		}

		if s.bus != nil {
			var msgAttrs map[string]json.RawMessage
			if len(msg.MessageAttributes) > 0 {
				msgAttrs = make(map[string]json.RawMessage, len(msg.MessageAttributes))
				for k, v := range msg.MessageAttributes {
					raw, err := json.Marshal(v)
					if err == nil {
						msgAttrs[k] = raw
					}
				}
			}
			snsEvt := &eventbus.SNSDeliveryEvent{
				TopicARN:          topic.Arn,
				MessageID:         msg.MessageId,
				Message:           message,
				Subject:           subject,
				MessageStructure:  msg.MessageStructure,
				MessageAttributes: msgAttrs,
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

	return messageID, nil
}
