// Package sns provides SNS (Simple Notification Service) operations for vorpalstacks.
package sns

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"

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

	// deliveryCtx ends at Close: a delivery retry ladder mid-sleep abandons
	// its remaining attempts instead of holding its goroutine (and Close's
	// wait) through the scheduled delays.
	deliveryCtx    context.Context
	deliveryCancel context.CancelFunc

	// throttleNext paces HTTP/S deliveries per subscription under the
	// delivery policy's maxReceivesPerSecond; throttleMu guards the map.
	throttleMu   sync.Mutex
	throttleNext map[string]time.Time

	// deliverySpawnMu and deliveryClosed close the goroutine spawn path
	// against Close's WaitGroup drain: an Add racing a Wait that observed
	// a zero counter is the misuse the WaitGroup spec forbids.
	deliverySpawnMu sync.RWMutex
	deliveryClosed  bool

	signingKeyMu   sync.Mutex
	signingKey     *rsa.PrivateKey
	signingCertPEM []byte
}

// store resolves the request region's SNS store through the shared
// get-or-create helper: two racing requests cannot leak a loser store
// (the helper closes it, stopping its background workers).
func (s *SNSService) store(reqCtx *request.RequestContext) (snsstore.SNSStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (snsstore.SNSStoreInterface, error) {
		basicStore, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return newStoreFromStorage(basicStore, s.accountID, reqCtx.GetRegion())
	})
}

// defaultDeliveryHTTPTimeout bounds every notification and confirmation
// POST the delivery engine and the confirmation sender issue.
const defaultDeliveryHTTPTimeout = 30 * time.Second

// NewSNSService creates a new SNS service instance.
// Cross-service delivery is routed through the event bus, which must be
// injected via SetEventBus before registering handlers.
func NewSNSService(storageMgr *storage.RegionStorageManager, accountID, region string) *SNSService {
	ctx, cancel := context.WithCancel(context.Background())
	return &SNSService{
		storageManager: storageMgr,
		accountID:      accountID,
		defaultRegion:  region,
		httpClient:     &http.Client{Timeout: defaultDeliveryHTTPTimeout},
		deliveryCtx:    ctx,
		deliveryCancel: cancel,
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

// Close ends the delivery lifecycle (retry ladders abandon their pending
// attempts) and waits for all in-flight delivery goroutines to complete.
// Every regional store the resolver cached — injected or lazily created —
// carries a background deduplication sweeper only the store's own Close
// stops, so the cached stores close here too; the shutdown wiring closing
// the default-region instance afterwards is harmless, because the store's
// Close is idempotent (a second cancel and the drained WaitGroup both
// tolerate repeats).
func (s *SNSService) Close() {
	// The spawn path closes before the Wait drains it: the Go spec
	// requires a positive-delta Add that starts when the counter is zero
	// to happen before a Wait, so a publish racing shutdown must be
	// excluded under the lock, not left to the race.
	s.deliverySpawnMu.Lock()
	s.deliveryClosed = true
	s.deliverySpawnMu.Unlock()

	s.deliveryCancel()
	s.deliveryWg.Wait()
	s.stores.Range(func(_, v any) bool {
		if store, ok := v.(snsstore.SNSStoreInterface); ok {
			store.Close()
		}
		return true
	})
}

// deliverAsync spawns a tracked goroutine to deliver a message to subscriptions.
func (s *SNSService) deliverAsync(msg *snsstore.Message, subs []*snsstore.Subscription, region string) {
	s.runTracked(fmt.Sprintf("message %s on topic %s", msg.MessageId, msg.TopicArn), func() {
		s.deliverToSubscriptions(msg, subs, region)
	})
}

// runTracked runs fn on a goroutine the delivery WaitGroup tracks, with the
// panic isolation the synchronous FIFO path also carries: a panicking
// delivery is logged and abandoned, never propagated into the request that
// spawned it and never left running past Close. Once Close has marked the
// spawn path closed no new goroutine joins the WaitGroup — a delivery
// racing shutdown is dropped with a warning rather than racing Close's
// Wait on a zero counter.
func (s *SNSService) runTracked(what string, fn func()) {
	s.deliverySpawnMu.RLock()
	defer s.deliverySpawnMu.RUnlock()
	if s.deliveryClosed {
		logs.Warn("SNS delivery not spawned — the service is closing", logs.String("delivery", what))
		return
	}
	s.deliveryWg.Add(1)
	go func() {
		defer s.deliveryWg.Done()
		defer recoverDeliveryPanic(what)
		fn()
	}()
}

// recoverDeliveryPanic logs a recovered delivery panic. what names the
// delivery that panicked (message and topic, or the confirmation target).
func recoverDeliveryPanic(what string) {
	if r := recover(); r != nil {
		logs.Error("SNS delivery panicked",
			logs.String("delivery", what),
			logs.Any("panic", r))
	}
}

// httpDeliveryClient returns the service's shared HTTP client for outbound
// notification and confirmation POSTs, constructing a default-timeout
// client when the service was built without one (test harnesses construct
// SNSService directly).
func (s *SNSService) httpDeliveryClient() *http.Client {
	if s.httpClient == nil {
		return &http.Client{Timeout: defaultDeliveryHTTPTimeout}
	}
	return s.httpClient
}

// deliverWithRecover runs one delivery pass, isolating a panic so it
// neither kills the delivering goroutine (the async path) nor propagates
// through the synchronous FIFO delivery path into the publish request.
func (s *SNSService) deliverWithRecover(msg *snsstore.Message, subs []*snsstore.Subscription, region string) {
	defer recoverDeliveryPanic(fmt.Sprintf("message %s on topic %s", msg.MessageId, msg.TopicArn))
	s.deliverToSubscriptions(msg, subs, region)
}

func (s *SNSService) handleBusDelivery(ctx context.Context, evt *eventbus.SNSDeliveryEvent) eventbus.HandlerResult {
	store, err := s.getSNSStoreByRegion(evt.Region)
	if err != nil {
		logs.Warn("SNS bus delivery: failed to get store", logs.String("region", evt.Region), logs.Err(err))
		return eventbus.HandlerResult{Error: err}
	}

	subscriptions, err := store.ListAllSubscriptionsByTopic(evt.TopicARN)
	if err != nil {
		// Surface, never swallow-and-ACK: a failed listing returned as an
		// empty handler result would ACK the bus event and permanently
		// lose the message; an Error result keeps the outbox entry
		// retryable instead.
		logs.Error("SNS bus delivery: listing subscriptions failed",
			logs.String("topicArn", evt.TopicARN),
			logs.String("messageId", evt.MessageID),
			logs.Err(err))
		return eventbus.HandlerResult{Error: err}
	}
	if len(subscriptions) == 0 {
		return eventbus.HandlerResult{}
	}

	msg := &snsstore.Message{
		MessageId:              evt.MessageID,
		TopicArn:               evt.TopicARN,
		Message:                evt.Message,
		Subject:                evt.Subject,
		MessageStructure:       evt.MessageStructure,
		MessageGroupId:         evt.MessageGroupId,
		MessageDeduplicationId: evt.MessageDeduplicationID,
		PublishedTimestamp:     evt.EventTimestamp(),
		MessageAttributes:      make(map[string]*snsstore.MessageAttribute),
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

	s.deliverToSubscriptions(msg, subscriptions, evt.Region)
	return eventbus.HandlerResult{}
}

// getSNSStoreByRegion resolves a region's SNS store through the shared
// get-or-create helper — the same loser-closing path s.store uses — instead
// of a hand-rolled Load/LoadOrStore whose racing loser leaked its store
// (and its background cleanup goroutine).
func (s *SNSService) getSNSStoreByRegion(region string) (snsstore.SNSStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (snsstore.SNSStoreInterface, error) {
		basicStore, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return newStoreFromStorage(basicStore, s.accountID, region)
	})
}

// newStoreFromStorage asserts the storage handle carries the transactional
// interface the store's multi-record commits need and builds the region's
// SNS store.
func newStoreFromStorage(basicStore storage.BasicStorage, accountID, region string) (snsstore.SNSStoreInterface, error) {
	tstore, ok := basicStore.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not support TransactionalStorageWith2PC")
	}
	return snsstore.NewSNSStore(tstore, accountID, region), nil
}

// GetSNSStoreForRegion resolves the SNS store of the given region — the
// region-resolving provider cross-service consumers address instead of a
// single concrete default-region store instance.
func (s *SNSService) GetSNSStoreForRegion(region string) (snsstore.SNSStoreInterface, error) {
	return s.getSNSStoreByRegion(region)
}

// initSigningKey loads or generates the SNS signing key with its
// certificate. A failed attempt leaves both fields unset so the next call
// retries — every envelope signs through this initialiser, and a transient
// storage or key-generation fault must not strip signatures from every
// later delivery for the process lifetime.
func (s *SNSService) initSigningKey() {
	s.signingKeyMu.Lock()
	defer s.signingKeyMu.Unlock()
	if s.signingKey != nil && s.signingCertPEM != nil {
		return
	}
	rs, err := s.storageManager.GetStorage(s.defaultRegion)
	if err != nil {
		return
	}
	bucket := rs.Bucket("sns-signing")

	// A persisted key/certificate pair initialises directly; a persisted
	// key without its certificate regenerates — the pair is the unit the
	// envelope needs.
	if keyPEM, err := bucket.Get([]byte("signing_key")); err == nil && keyPEM != nil {
		if certPEM, err := bucket.Get([]byte("signing_cert")); err == nil && certPEM != nil {
			if parsedKey, err := vcrypto.ParsePrivateKeyPEM(keyPEM); err == nil {
				if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
					s.signingKey = rsaKey
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
	s.signingKey = privateKey
	s.signingCertPEM = certBytes

	if err := bucket.Put([]byte("signing_key"), []byte(keyPEM)); err != nil {
		logs.Warn("failed to persist SNS signing key; key regenerated on restart will invalidate existing message signatures", logs.Err(err))
	}
	if err := bucket.Put([]byte("signing_cert"), certBytes); err != nil {
		logs.Warn("failed to persist SNS signing certificate; certificate regenerated on restart will invalidate existing message signatures", logs.Err(err))
	}
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
// subscriptions. It is the cross-service invoker entry, and it delegates to
// publishCore so an internal publish carries the same validation (message
// size, subject bounds, MessageStructure checks), FIFO semantics (group
// requirement, deduplication, sequence numbers) and the fan-out seam as the
// Publish API — the invoker plane can no longer publish past validation.
// Returns the generated message ID.
func (s *SNSService) PublishToTopic(ctx context.Context, region, topicArn, message, subject string, messageAttributes map[string]invokers.SQSMessageAttribute) (string, error) {
	store, err := s.getSNSStoreByRegion(region)
	if err != nil {
		return "", fmt.Errorf("failed to get SNS store for region %s: %w", region, err)
	}

	in := PublishInput{
		TopicArn: topicArn,
		Message:  message,
		Subject:  subject,
	}
	// Typed attributes ride the wire-attribute form so the same
	// parseMessageAttributes validation path applies to them.
	if len(messageAttributes) > 0 {
		attrs := make(map[string]interface{}, len(messageAttributes))
		for k, v := range messageAttributes {
			entry := map[string]interface{}{"DataType": v.DataType}
			if len(v.BinaryValue) > 0 {
				entry["BinaryValue"] = base64.StdEncoding.EncodeToString(v.BinaryValue)
			} else {
				entry["StringValue"] = v.StringValue
			}
			attrs[k] = entry
		}
		in.Parameters = map[string]interface{}{"MessageAttributes": attrs}
	}

	result, err := s.publishCore(store, region, in)
	if err != nil {
		return "", err
	}
	if messageID, ok := result.(map[string]interface{})["MessageId"].(string); ok {
		return messageID, nil
	}
	return "", fmt.Errorf("sns: publish result carries no MessageId")
}
