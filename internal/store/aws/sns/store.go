package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
)

const deduplicationWindow = 5 * time.Minute

// SNSStore provides SNS topic and subscription storage functionality.
type SNSStore struct {
	*common.BaseStore
	topicSubscriptionsStore   *common.BaseStore
	platformApplicationsStore *common.BaseStore
	platformEndpointsStore    *common.BaseStore
	platformAppEndpointsIndex storage.Bucket
	*common.TagStore
	arnBuilder         *svcarn.ARNBuilder
	accountID          string
	region             string
	deduplicationCache map[string]*deduplicationEntry
	deduplicationMu    sync.RWMutex
	sequenceCounters   map[string]int64
	sequenceMu         sync.Mutex
	topicMu            sync.RWMutex
	subscriptionMu     sync.RWMutex
	platformAppMu      sync.RWMutex
	platformEndpointMu sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
}

type deduplicationEntry struct {
	messageID string
	expiresAt time.Time
}

// NewSNSStore creates a new SNSStore instance with the specified storage, account ID, and region.
func NewSNSStore(store storage.BasicStorage, accountID, region string) *SNSStore {
	ctx, cancel := context.WithCancel(context.Background())
	s := &SNSStore{
		BaseStore:                 common.NewBaseStore(store.Bucket("sns-topics-"+region), "sns-topics"),
		topicSubscriptionsStore:   common.NewBaseStore(store.Bucket("sns-subscriptions-"+region), "sns-subscriptions"),
		platformApplicationsStore: common.NewBaseStore(store.Bucket("sns-platform-apps-"+region), "sns-platform-apps"),
		platformEndpointsStore:    common.NewBaseStore(store.Bucket("sns-platform-endpoints-"+region), "sns-platform-endpoints"),
		platformAppEndpointsIndex: store.Bucket("sns-app-endpoints-index-" + region),
		TagStore:                  common.NewTagStoreWithRegion(store, "sns", region),
		arnBuilder:                svcarn.NewARNBuilder(accountID, region),
		accountID:                 accountID,
		region:                    region,
		deduplicationCache:        make(map[string]*deduplicationEntry),
		sequenceCounters:          make(map[string]int64),
		ctx:                       ctx,
		cancel:                    cancel,
	}
	s.wg.Add(1)
	go s.cleanupExpiredDeduplicationEntries()
	s.rebuildEndpointIndex()
	return s
}

// rebuildEndpointIndex populates the reverse index from existing endpoint data.
// This handles migration for data created before the index existed.
func (s *SNSStore) rebuildEndpointIndex() {
	if s.platformAppEndpointsIndex.Count() > 0 {
		return
	}
	if s.platformEndpointsStore.Count() == 0 {
		return
	}
	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()
	_ = s.platformEndpointsStore.ForEach(func(key string, value []byte) error {
		var ep PlatformEndpoint
		if err := json.Unmarshal(value, &ep); err != nil {
			return nil
		}
		if ep.PlatformApplicationArn != "" {
			idxKey := ep.PlatformApplicationArn + "\x00" + ep.EndpointArn
			_ = s.platformAppEndpointsIndex.Put([]byte(idxKey), []byte("1"))
		}
		return nil
	})
}

// cleanupExpiredDeduplicationEntries periodically removes expired deduplication cache and stale sequence counter entries.
func (s *SNSStore) cleanupExpiredDeduplicationEntries() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			s.deduplicationMu.Lock()
			for k, entry := range s.deduplicationCache {
				if now.After(entry.expiresAt) {
					delete(s.deduplicationCache, k)
				}
			}
			s.deduplicationMu.Unlock()

			s.sequenceMu.Lock()
			if len(s.sequenceCounters) > 1000 {
				for k, v := range s.sequenceCounters {
					if v == 0 {
						delete(s.sequenceCounters, k)
					}
				}
			}
			s.sequenceMu.Unlock()
		}
	}
}

// Close stops the cleanup goroutine and releases resources.
func (s *SNSStore) Close() {
	if s.cancel != nil {
		s.cancel()
		s.wg.Wait()
	}
}

func (s *SNSStore) buildTopicArn(topicName string) string {
	return s.arnBuilder.SNS().Topic(topicName)
}

func (s *SNSStore) buildSubscriptionArn(topicArn, subscriptionID string) string {
	return topicArn + ":" + subscriptionID
}

// CreateTopic creates a new SNS topic in the store.
// Returns the created topic or an error if the topic already exists or validation fails.
func (s *SNSStore) CreateTopic(topic *Topic) (*Topic, error) {
	if topic.Name == "" {
		return nil, ErrInvalidTopicName
	}

	topicArn := s.buildTopicArn(topic.Name)
	if s.Exists(topicArn) {
		var existingTopic Topic
		if err := s.BaseStore.Get(topicArn, &existingTopic); err != nil {
			return nil, err
		}
		return &existingTopic, nil
	}

	now := time.Now().UTC()
	topic.Arn = topicArn
	topic.Owner = s.accountID
	topic.CreatedDate = now
	topic.LastModifiedTime = now
	if topic.Attributes == nil {
		topic.Attributes = make(map[string]string)
	}

	if err := s.Put(topicArn, topic); err != nil {
		return nil, err
	}

	if len(topic.Tags) > 0 {
		tagSlice := make([]common.Tag, 0, len(topic.Tags))
		for k, v := range topic.Tags {
			tagSlice = append(tagSlice, common.Tag{Key: k, Value: v})
		}
		if err := s.TagStore.TagResourceFromSlice(topicArn, tagSlice); err != nil {
			return nil, err
		}
	}

	return topic, nil
}

// GetTopic retrieves an SNS topic by its ARN from the store.
// Returns the topic or an error if not found.
func (s *SNSStore) GetTopic(topicArn string) (*Topic, error) {
	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return nil, ErrTopicNotFound
	}
	return &topic, nil
}

// GetTopicByName retrieves an SNS topic by its name from the store.
// Returns the topic or an error if not found.
func (s *SNSStore) GetTopicByName(topicName string) (*Topic, error) {
	topicArn := s.buildTopicArn(topicName)
	return s.GetTopic(topicArn)
}

// UpdateTopic updates an existing SNS topic in the store.
// Returns an error if the topic does not exist.
func (s *SNSStore) UpdateTopic(topic *Topic) error {
	if !s.Exists(topic.Arn) {
		return ErrTopicNotFound
	}
	topic.LastModifiedTime = time.Now().UTC()
	return s.Put(topic.Arn, topic)
}

// DeleteTopic deletes an SNS topic by its ARN from the store.
// Returns an error if the topic does not exist.
func (s *SNSStore) DeleteTopic(topicArn string) error {
	if !s.Exists(topicArn) {
		return ErrTopicNotFound
	}

	subscriptions, err := s.ListSubscriptionsByTopic(topicArn, common.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing topic subscriptions: %w", err)
	}
	for _, sub := range subscriptions.Items {
		if err := s.topicSubscriptionsStore.Delete(sub.SubscriptionArn); err != nil {
			return fmt.Errorf("deleting subscription: %w", err)
		}
	}

	if err := s.TagStore.Delete(topicArn); err != nil {
		return fmt.Errorf("deleting topic tags: %w", err)
	}

	s.deduplicationMu.Lock()
	prefix := topicArn + ":"
	for key := range s.deduplicationCache {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			delete(s.deduplicationCache, key)
		}
	}
	s.deduplicationMu.Unlock()

	s.sequenceMu.Lock()
	for key := range s.sequenceCounters {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			delete(s.sequenceCounters, key)
		}
	}
	s.sequenceMu.Unlock()

	return s.BaseStore.Delete(topicArn)
}

// ListTopics returns a list of SNS topics from the store using the specified list options.
func (s *SNSStore) ListTopics(opts common.ListOptions) (*common.ListResult[Topic], error) {
	return common.List[Topic](s.BaseStore, opts, nil)
}

// SetTopicAttributes sets the attributes for an SNS topic.
func (s *SNSStore) SetTopicAttributes(topicArn string, attributes map[string]string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return ErrTopicNotFound
	}

	if topic.Attributes == nil {
		topic.Attributes = make(map[string]string)
	}
	for k, v := range attributes {
		topic.Attributes[k] = v
	}

	topic.LastModifiedTime = time.Now().UTC()
	return s.Put(topic.Arn, &topic)
}

// CreateSubscription creates a new SNS subscription in the store.
// Returns the created subscription or an error if validation fails.
func (s *SNSStore) CreateSubscription(subscription *Subscription) (*Subscription, error) {
	if subscription.TopicArn == "" {
		return nil, ErrInvalidParameter
	}
	if subscription.Protocol == "" {
		return nil, ErrInvalidProtocol
	}
	if subscription.Endpoint == "" {
		return nil, ErrInvalidEndpoint
	}

	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err != nil {
		return nil, ErrTopicNotFound
	}

	subscriptionID := uuid.New().String()
	subscriptionArn := s.buildSubscriptionArn(subscription.TopicArn, subscriptionID)
	subscription.SubscriptionArn = subscriptionArn
	subscription.TopicArn = topic.Arn
	subscription.PendingConfirmation = true
	subscription.ConfirmationToken = uuid.New().String()
	subscription.CreatedDate = time.Now().UTC()
	if subscription.Attributes == nil {
		subscription.Attributes = make(map[string]string)
	}

	if err := s.topicSubscriptionsStore.Put(subscriptionArn, subscription); err != nil {
		return nil, err
	}

	topic.SubscriptionsPending++
	topic.LastModifiedTime = time.Now().UTC()
	if err := s.Put(topic.Arn, &topic); err != nil {
		return nil, err
	}

	return subscription, nil
}

// GetSubscription retrieves an SNS subscription by its ARN from the store.
// Returns the subscription or an error if not found.
func (s *SNSStore) GetSubscription(subscriptionArn string) (*Subscription, error) {
	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return nil, ErrSubscriptionNotFound
	}
	return &subscription, nil
}

// DeleteSubscription deletes an SNS subscription by its ARN from the store.
// Returns an error if the subscription does not exist.
func (s *SNSStore) DeleteSubscription(subscriptionArn string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return ErrSubscriptionNotFound
	}

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err == nil {
		if subscription.PendingConfirmation {
			if topic.SubscriptionsPending > 0 {
				topic.SubscriptionsPending--
			}
		} else {
			if topic.SubscriptionsConfirmed > 0 {
				topic.SubscriptionsConfirmed--
			}
		}
		topic.SubscriptionsDeleted++
		topic.LastModifiedTime = time.Now().UTC()
		if err := s.Put(topic.Arn, &topic); err != nil {
			return fmt.Errorf("updating topic: %w", err)
		}
	}

	return s.topicSubscriptionsStore.Delete(subscriptionArn)
}

// UpdateSubscription updates an existing SNS subscription in the store.
func (s *SNSStore) UpdateSubscription(subscription *Subscription) error {
	return s.topicSubscriptionsStore.Put(subscription.SubscriptionArn, subscription)
}

// AutoConfirmSubscription confirms a subscription without requiring a token.
// This is used for sqs/lambda/http/https protocols that are auto-confirmed.
// It updates the subscription status and adjusts the topic counters.
func (s *SNSStore) AutoConfirmSubscription(subscription *Subscription) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	subscription.PendingConfirmation = false
	subscription.ConfirmationToken = ""

	if err := s.topicSubscriptionsStore.Put(subscription.SubscriptionArn, subscription); err != nil {
		return fmt.Errorf("updating subscription: %w", err)
	}

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err == nil {
		if topic.SubscriptionsPending > 0 {
			topic.SubscriptionsPending--
		}
		topic.SubscriptionsConfirmed++
		topic.LastModifiedTime = time.Now().UTC()
		if err := s.Put(topic.Arn, &topic); err != nil {
			return fmt.Errorf("updating topic: %w", err)
		}
	}

	return nil
}

// ListSubscriptions returns a list of SNS subscriptions from the store using the specified list options.
func (s *SNSStore) ListSubscriptions(opts common.ListOptions) (*common.ListResult[Subscription], error) {
	return common.List[Subscription](s.topicSubscriptionsStore, opts, nil)
}

// ListSubscriptionsByTopic returns a list of subscriptions for a specific topic ARN.
func (s *SNSStore) ListSubscriptionsByTopic(topicArn string, opts common.ListOptions) (*common.ListResult[Subscription], error) {
	return common.List[Subscription](s.topicSubscriptionsStore, opts, func(sub *Subscription) bool {
		return sub.TopicArn == topicArn
	})
}

// ConfirmSubscription confirms an SNS subscription using the subscription ARN and confirmation token.
func (s *SNSStore) ConfirmSubscription(subscriptionArn, token string) (*Subscription, error) {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return nil, ErrSubscriptionNotFound
	}

	if subscription.ConfirmationToken != token {
		return nil, ErrInvalidToken
	}

	subscription.PendingConfirmation = false
	subscription.ConfirmationToken = ""

	if err := s.topicSubscriptionsStore.Put(subscriptionArn, &subscription); err != nil {
		return nil, err
	}

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err == nil {
		if topic.SubscriptionsPending > 0 {
			topic.SubscriptionsPending--
		}
		topic.SubscriptionsConfirmed++
		topic.LastModifiedTime = time.Now().UTC()
		if err := s.Put(topic.Arn, &topic); err != nil {
			return nil, err
		}
	}

	return &subscription, nil
}

// SetSubscriptionAttributes sets the attributes for an SNS subscription.
func (s *SNSStore) SetSubscriptionAttributes(subscriptionArn string, attributes map[string]string) error {
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return ErrSubscriptionNotFound
	}

	if subscription.Attributes == nil {
		subscription.Attributes = make(map[string]string)
	}
	for k, v := range attributes {
		subscription.Attributes[k] = v
	}

	return s.topicSubscriptionsStore.Put(subscriptionArn, &subscription)
}

// GetSubscriptionAttributes retrieves the attributes for an SNS subscription.
func (s *SNSStore) GetSubscriptionAttributes(subscriptionArn string) (map[string]string, error) {
	subscription, err := s.GetSubscription(subscriptionArn)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	for k, v := range subscription.Attributes {
		attrs[k] = v
	}
	attrs["SubscriptionArn"] = subscription.SubscriptionArn
	attrs["TopicArn"] = subscription.TopicArn
	attrs["Protocol"] = subscription.Protocol
	attrs["Endpoint"] = subscription.Endpoint
	attrs["Owner"] = subscription.Owner
	attrs["PendingConfirmation"] = "false"
	if subscription.PendingConfirmation {
		attrs["PendingConfirmation"] = "true"
	}
	attrs["ConfirmationWasAuthenticated"] = "false"
	if subscription.ConfirmationWasAuthenticated {
		attrs["ConfirmationWasAuthenticated"] = "true"
	}
	if subscription.ConfirmationToken != "" {
		attrs["Token"] = subscription.ConfirmationToken
	}

	return attrs, nil
}

// ListTagsForResource retrieves the tags associated with an SNS resource.
func (s *SNSStore) ListTagsForResource(resourceArn string) ([]common.Tag, error) {
	return s.TagStore.ListTagsAsSlice(resourceArn)
}

// TagResource adds tags to an SNS resource.
func (s *SNSStore) TagResource(resourceArn string, tags []common.Tag) error {
	return s.TagStore.TagResourceFromSlice(resourceArn, tags)
}

// UntagResource removes tags from an SNS resource.
func (s *SNSStore) UntagResource(resourceArn string, tagKeys []string) error {
	return s.TagStore.UntagResource(resourceArn, tagKeys)
}

// CheckDeduplication checks if a message with the given deduplication ID was recently published.
// Returns (existingMessageID, true) if duplicate found, ("", false) otherwise.
func (s *SNSStore) CheckDeduplication(topicArn, messageDeduplicationId string) (string, bool) {
	if messageDeduplicationId == "" {
		return "", false
	}

	key := topicArn + ":" + messageDeduplicationId

	s.deduplicationMu.RLock()
	defer s.deduplicationMu.RUnlock()

	if entry, ok := s.deduplicationCache[key]; ok {
		if time.Now().Before(entry.expiresAt) {
			return entry.messageID, true
		}
	}
	return "", false
}

// RecordDeduplication records a deduplication ID for the 5-minute window.
func (s *SNSStore) RecordDeduplication(topicArn, messageDeduplicationId, messageID string) {
	if messageDeduplicationId == "" {
		return
	}

	key := topicArn + ":" + messageDeduplicationId

	s.deduplicationMu.Lock()
	defer s.deduplicationMu.Unlock()

	s.deduplicationCache[key] = &deduplicationEntry{
		messageID: messageID,
		expiresAt: time.Now().Add(deduplicationWindow),
	}

	if len(s.deduplicationCache) > 100 {
		now := time.Now()
		for k, entry := range s.deduplicationCache {
			if now.After(entry.expiresAt) {
				delete(s.deduplicationCache, k)
			}
		}
	}
}

// GetNextSequenceNumber returns the next sequence number for a FIFO topic message group.
// The sequence number is a large, non-consecutive number that Amazon SNS assigns to each message.
func (s *SNSStore) GetNextSequenceNumber(topicArn, messageGroupId string) string {
	key := topicArn + ":" + messageGroupId

	s.sequenceMu.Lock()
	defer s.sequenceMu.Unlock()

	s.sequenceCounters[key]++
	return fmt.Sprintf("%d", s.sequenceCounters[key])
}

// GetDataProtectionPolicy retrieves the data protection policy for a topic.
func (s *SNSStore) GetDataProtectionPolicy(topicArn string) (string, error) {
	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return "", ErrTopicNotFound
	}
	return topic.DataProtectionPolicy, nil
}

// PutDataProtectionPolicy sets the data protection policy for a topic.
func (s *SNSStore) PutDataProtectionPolicy(topicArn, policy string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return ErrTopicNotFound
	}

	topic.DataProtectionPolicy = policy
	topic.LastModifiedTime = time.Now().UTC()
	return s.Put(topicArn, &topic)
}

// AddPermission adds a permission to a topic.
func (s *SNSStore) AddPermission(topicArn string, permission *Permission) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return ErrTopicNotFound
	}

	if topic.Permissions == nil {
		topic.Permissions = []Permission{}
	}

	for i, p := range topic.Permissions {
		if p.Label == permission.Label {
			topic.Permissions[i] = *permission
			topic.LastModifiedTime = time.Now().UTC()
			return s.Put(topicArn, &topic)
		}
	}

	topic.Permissions = append(topic.Permissions, *permission)
	topic.LastModifiedTime = time.Now().UTC()
	return s.Put(topicArn, &topic)
}

// RemovePermission removes a permission from a topic by label.
func (s *SNSStore) RemovePermission(topicArn, label string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return ErrTopicNotFound
	}

	for i, p := range topic.Permissions {
		if p.Label == label {
			topic.Permissions = append(topic.Permissions[:i], topic.Permissions[i+1:]...)
			topic.LastModifiedTime = time.Now().UTC()
			return s.Put(topicArn, &topic)
		}
	}

	return nil
}

func (s *SNSStore) buildPlatformApplicationArn(platform, name string) string {
	return s.arnBuilder.SNS().PlatformApplication(platform, name)
}



// CreatePlatformApplication creates a new SNS platform application.
func (s *SNSStore) CreatePlatformApplication(app *PlatformApplication) (*PlatformApplication, error) {
	if app.Name == "" {
		return nil, ErrInvalidParameter
	}
	if app.Platform == "" {
		return nil, ErrInvalidParameter
	}

	platformArn := s.buildPlatformApplicationArn(app.Platform, app.Name)

	s.platformAppMu.Lock()
	defer s.platformAppMu.Unlock()

	if s.platformApplicationsStore.Exists(platformArn) {
		return nil, ErrPlatformApplicationAlreadyExists
	}

	app.PlatformApplicationArn = platformArn
	if app.Attributes == nil {
		app.Attributes = make(map[string]string)
	}

	if err := s.platformApplicationsStore.Put(platformArn, app); err != nil {
		return nil, err
	}

	return app, nil
}

// GetPlatformApplication retrieves an SNS platform application by its ARN.
func (s *SNSStore) GetPlatformApplication(arn string) (*PlatformApplication, error) {
	var app PlatformApplication
	if err := s.platformApplicationsStore.Get(arn, &app); err != nil {
		return nil, ErrPlatformApplicationNotFound
	}
	return &app, nil
}

// DeletePlatformApplication deletes an SNS platform application and all its endpoints.
func (s *SNSStore) DeletePlatformApplication(arn string) error {
	s.platformAppMu.Lock()
	defer s.platformAppMu.Unlock()

	if !s.platformApplicationsStore.Exists(arn) {
		return ErrPlatformApplicationNotFound
	}

	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	prefix := []byte(arn + "\x00")
	iter := s.platformAppEndpointsIndex.ScanPrefix(prefix)
	for iter.Next() {
		epArn := string(iter.Key()[len(prefix):])
		if delErr := s.platformEndpointsStore.Delete(epArn); delErr != nil {
			logs.Warn("sns: failed to delete endpoint during application deletion",
				logs.String("endpointArn", epArn), logs.Err(delErr))
		}
		if idxErr := s.platformAppEndpointsIndex.Delete(iter.Key()); idxErr != nil {
			logs.Warn("sns: failed to delete endpoint index entry",
				logs.String("endpointArn", epArn), logs.Err(idxErr))
		}
	}

	if err := s.TagStore.Delete(arn); err != nil {
		return fmt.Errorf("deleting platform application tags: %w", err)
	}

	return s.platformApplicationsStore.Delete(arn)
}

// ListPlatformApplications returns all SNS platform applications using the specified list options.
func (s *SNSStore) ListPlatformApplications(opts common.ListOptions) (*common.ListResult[PlatformApplication], error) {
	return common.List[PlatformApplication](s.platformApplicationsStore, opts, nil)
}

// SetPlatformApplicationAttributes sets the attributes for an SNS platform application.
func (s *SNSStore) SetPlatformApplicationAttributes(arn string, attributes map[string]string) error {
	s.platformAppMu.Lock()
	defer s.platformAppMu.Unlock()

	var app PlatformApplication
	if err := s.platformApplicationsStore.Get(arn, &app); err != nil {
		return ErrPlatformApplicationNotFound
	}

	if app.Attributes == nil {
		app.Attributes = make(map[string]string)
	}
	for k, v := range attributes {
		app.Attributes[k] = v
	}

	return s.platformApplicationsStore.Put(arn, &app)
}

// GetPlatformApplicationAttributes retrieves the attributes of an SNS platform application.
func (s *SNSStore) GetPlatformApplicationAttributes(arn string) (map[string]string, error) {
	app, err := s.GetPlatformApplication(arn)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	for k, v := range app.Attributes {
		attrs[k] = v
	}
	return attrs, nil
}

// CreatePlatformEndpoint creates a new SNS platform endpoint, or updates the existing one if the ARN matches.
func (s *SNSStore) CreatePlatformEndpoint(endpoint *PlatformEndpoint) (*PlatformEndpoint, error) {
	if endpoint.PlatformApplicationArn == "" {
		return nil, ErrInvalidParameter
	}
	if endpoint.Token == "" {
		return nil, ErrInvalidParameter
	}

	s.platformAppMu.RLock()
	if !s.platformApplicationsStore.Exists(endpoint.PlatformApplicationArn) {
		s.platformAppMu.RUnlock()
		return nil, ErrPlatformApplicationNotFound
	}
	s.platformAppMu.RUnlock()

	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	if endpoint.EndpointArn == "" {
		endpointID := uuid.New().String()
		parts := strings.SplitN(endpoint.PlatformApplicationArn, "app/", 2)
		if len(parts) != 2 {
			return nil, ErrInvalidParameter
		}
		appPart := parts[1]
		endpoint.EndpointArn = s.arnBuilder.Build("sns", "endpoint/"+appPart+"/"+endpointID)
	}

	if s.platformEndpointsStore.Exists(endpoint.EndpointArn) {
		var existing PlatformEndpoint
		if err := s.platformEndpointsStore.Get(endpoint.EndpointArn, &existing); err == nil {
			existing.Token = endpoint.Token
			if endpoint.CustomUserData != "" {
				existing.CustomUserData = endpoint.CustomUserData
			}
			if endpoint.Attributes != nil {
				if existing.Attributes == nil {
					existing.Attributes = make(map[string]string)
				}
				for k, v := range endpoint.Attributes {
					existing.Attributes[k] = v
				}
			}
			if err := s.platformEndpointsStore.Put(endpoint.EndpointArn, &existing); err != nil {
				return nil, err
			}
			return &existing, nil
		}
	}

	if endpoint.Attributes == nil {
		endpoint.Attributes = make(map[string]string)
	}
	if _, ok := endpoint.Attributes["Enabled"]; !ok {
		endpoint.Attributes["Enabled"] = "true"
	}

	if err := s.platformEndpointsStore.Put(endpoint.EndpointArn, endpoint); err != nil {
		return nil, err
	}

	idxKey := endpoint.PlatformApplicationArn + "\x00" + endpoint.EndpointArn
	if err := s.platformAppEndpointsIndex.Put([]byte(idxKey), []byte("1")); err != nil {
		logs.Warn("sns: failed to write endpoint index entry",
			logs.String("endpointArn", endpoint.EndpointArn), logs.Err(err))
	}

	return endpoint, nil
}

// GetEndpoint retrieves an SNS platform endpoint by its ARN.
func (s *SNSStore) GetEndpoint(arn string) (*PlatformEndpoint, error) {
	var endpoint PlatformEndpoint
	if err := s.platformEndpointsStore.Get(arn, &endpoint); err != nil {
		return nil, ErrEndpointNotFound
	}
	return &endpoint, nil
}

// DeleteEndpoint deletes an SNS platform endpoint by its ARN.
func (s *SNSStore) DeleteEndpoint(arn string) error {
	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	if !s.platformEndpointsStore.Exists(arn) {
		return ErrEndpointNotFound
	}

	var endpoint PlatformEndpoint
	if err := s.platformEndpointsStore.Get(arn, &endpoint); err == nil {
		idxKey := endpoint.PlatformApplicationArn + "\x00" + arn
		if idxErr := s.platformAppEndpointsIndex.Delete([]byte(idxKey)); idxErr != nil {
			logs.Warn("sns: failed to delete endpoint index entry",
				logs.String("endpointArn", arn), logs.Err(idxErr))
		}
	}

	return s.platformEndpointsStore.Delete(arn)
}

// GetEndpointAttributes retrieves the attributes of an SNS platform endpoint.
func (s *SNSStore) GetEndpointAttributes(arn string) (map[string]string, error) {
	endpoint, err := s.GetEndpoint(arn)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	attrs["EndpointArn"] = endpoint.EndpointArn
	attrs["Token"] = endpoint.Token
	for k, v := range endpoint.Attributes {
		attrs[k] = v
	}
	return attrs, nil
}

// SetEndpointAttributes sets the attributes for an SNS platform endpoint.
func (s *SNSStore) SetEndpointAttributes(arn string, attributes map[string]string) error {
	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	var endpoint PlatformEndpoint
	if err := s.platformEndpointsStore.Get(arn, &endpoint); err != nil {
		return ErrEndpointNotFound
	}

	if endpoint.Attributes == nil {
		endpoint.Attributes = make(map[string]string)
	}
	for k, v := range attributes {
		endpoint.Attributes[k] = v
	}

	return s.platformEndpointsStore.Put(arn, &endpoint)
}

// ListEndpointsByPlatformApplication returns all endpoints for an SNS platform application.
// Uses the reverse index bucket for O(k) lookup where k is the number of endpoints
// for this application, instead of scanning all endpoints.
func (s *SNSStore) ListEndpointsByPlatformApplication(platformAppArn string, opts common.ListOptions) (*common.ListResult[PlatformEndpoint], error) {
	if opts.MaxItems <= 0 {
		opts.MaxItems = common.DefaultMaxItems
	}
	if opts.MaxItems > common.AbsoluteMaxItems {
		opts.MaxItems = common.AbsoluteMaxItems
	}

	prefix := []byte(platformAppArn + "\x00")
	var items []*PlatformEndpoint
	var lastEpArn string
	count := 0
	hasMore := false
	started := opts.Marker == ""

	iter := s.platformAppEndpointsIndex.ScanPrefix(prefix)
	for iter.Next() {
		epArn := string(iter.Key()[len(prefix):])

		if !started {
			if epArn == opts.Marker {
				started = true
			} else if epArn > opts.Marker {
				started = true
			}
			if !started {
				continue
			}
		}

		if count >= opts.MaxItems {
			hasMore = true
			break
		}

		var ep PlatformEndpoint
		if err := s.platformEndpointsStore.Get(epArn, &ep); err != nil {
			logs.Warn("sns: stale endpoint index entry, cleaning up",
				logs.String("endpointArn", epArn),
				logs.String("appArn", platformAppArn))
			_ = s.platformAppEndpointsIndex.Delete(iter.Key())
			continue
		}

		items = append(items, &ep)
		lastEpArn = epArn
		count++
	}

	nextMarker := ""
	isTruncated := false
	if hasMore {
		nextMarker = lastEpArn
		isTruncated = true
	}

	return &common.ListResult[PlatformEndpoint]{
		Items:       items,
		NextMarker:  nextMarker,
		IsTruncated: isTruncated,
	}, nil
}
