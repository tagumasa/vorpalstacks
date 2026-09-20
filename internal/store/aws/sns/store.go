package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	types "vorpalstacks/internal/common/tags"
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
	// storage is the transactional handle the multi-record sequences commit
	// through; the pre-opened BaseStore and Bucket handles below serve the
	// single-record paths. Both forms resolve the same buckets — the
	// bucket-name methods are the single definition sites of the names.
	storage                   storage.TransactionalStorageWith2PC
	topicSubscriptionsStore   *common.BaseStore
	platformApplicationsStore *common.BaseStore
	platformEndpointsStore    *common.BaseStore
	platformAppEndpointsIndex storage.Bucket
	endpointTokenIndex        storage.Bucket
	// fifoDedupStore and fifoSequenceStore hold the durable FIFO state: the
	// deduplication-window entries and the per-message-group sequence
	// counters live in their own buckets so both survive a process restart —
	// a publisher retry inside the window still hits, and SequenceNumber
	// keeps increasing per message group. deduplicationMu and sequenceMu
	// serialise each bucket's check-and-record / read-modify-write section
	// in-process; the group lock table serialises publishing.
	fifoDedupStore    storage.Bucket
	fifoSequenceStore storage.Bucket
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string

	// Lock matrix — the acquisition discipline every store method follows.
	// Family order: topicMu → subscriptionMu and platformAppMu →
	// platformEndpointMu. No section holds mutexes of both families.
	//
	//   topicMu             the topic record: CreateTopic, DeleteTopic,
	//                       SetTopicAttributes, PutDataProtectionPolicy,
	//                       AddPermission, RemovePermission
	//   subscriptionMu      subscription records. Sections that also read
	//                       or mutate the topic record (its counters) hold
	//                       topicMu FIRST:
	//                       CreateSubscription, DeleteSubscription,
	//                       AutoConfirmSubscription, ConfirmSubscription.
	//                       Subscription-record-only sections hold
	//                       subscriptionMu alone: SetSubscriptionAttributes
	//                       (single-record scope cannot invert the order).
	//   platformAppMu       platform application records:
	//                       CreatePlatformApplication,
	//                       SetPlatformApplicationAttributes. Sections that
	//                       also touch endpoints hold platformAppMu FIRST:
	//                       DeletePlatformApplication,
	//                       CreatePlatformEndpoint,
	//                       ListEndpointsByPlatformApplication.
	//   platformEndpointMu  endpoint records and both reverse indexes:
	//                       DeleteEndpoint, SetEndpointAttributes.
	//   deduplicationMu,    the durable FIFO state's single-key sections
	//   sequenceMu          (check-and-record, sequence read-modify-write);
	//                       leaf mutexes, never taken with a family lock.
	//   fifoGroupMu         guards the fifoGroupLocks table alone; it is
	//                       released before the acquired entry lock, so it
	//                       never nests.
	//   fifoGroupLocks      one lock per topic message group. The service
	//                       holds the entry across sequence allocation and
	//                       synchronous delivery, making sequence order the
	//                       delivery order — the per-group FIFO guarantee.
	//                       Entries are never purged: a deleted topic's
	//                       entry is inert, and a re-created topic name
	//                       safely reuses its key. Entries are acquired
	//                       with no other lock held.
	//
	// Single-record reads ride storage-layer atomicity; only the
	// multi-record sections listed above take family locks.
	deduplicationMu    sync.Mutex
	sequenceMu         sync.Mutex
	fifoGroupMu        sync.Mutex
	fifoGroupLocks     map[string]*sync.Mutex
	topicMu            sync.RWMutex
	subscriptionMu     sync.RWMutex
	platformAppMu      sync.RWMutex
	platformEndpointMu sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
}

// fifoDedupRecord is the durable form of one deduplication-window entry:
// the original publish's identifiers — returned verbatim on a dedup hit —
// and the instant the window closes.
type fifoDedupRecord struct {
	MessageID      string    `json:"messageId"`
	SequenceNumber string    `json:"sequenceNumber"`
	ExpiresAt      time.Time `json:"expiresAt"`
}

// fifoSequenceRecord is the durable form of one message group's sequence
// counter. Durability keeps SequenceNumber increasing across process
// restarts, as the Publish documentation states it does per MessageGroupId.
type fifoSequenceRecord struct {
	Value int64 `json:"value"`
}

// NewSNSStore creates a new SNSStore instance with the specified storage, account ID, and region.
func NewSNSStore(store storage.TransactionalStorageWith2PC, accountID, region string) *SNSStore {
	ctx, cancel := context.WithCancel(context.Background())
	s := &SNSStore{
		storage:        store,
		TagStore:       common.NewTagStoreWithRegion(store, "sns", region, common.StandardTagBudget("TagLimitExceeded")),
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		accountID:      accountID,
		region:         region,
		fifoGroupLocks: make(map[string]*sync.Mutex),
		ctx:            ctx,
		cancel:         cancel,
	}
	s.BaseStore = common.NewBaseStore(store.Bucket(s.topicsBucketName()), "sns-topics")
	s.topicSubscriptionsStore = common.NewBaseStore(store.Bucket(s.subscriptionsBucketName()), "sns-subscriptions")
	s.platformApplicationsStore = common.NewBaseStore(store.Bucket(s.platformAppsBucketName()), "sns-platform-apps")
	s.platformEndpointsStore = common.NewBaseStore(store.Bucket(s.platformEndpointsBucketName()), "sns-platform-endpoints")
	s.platformAppEndpointsIndex = store.Bucket(s.appEndpointsIndexBucketName())
	s.endpointTokenIndex = store.Bucket(s.endpointTokenIndexBucketName())
	s.fifoDedupStore = store.Bucket(s.fifoDedupBucketName())
	s.fifoSequenceStore = store.Bucket(s.fifoSequenceBucketName())
	s.wg.Add(1)
	go s.cleanupExpiredDeduplicationEntries()
	return s
}

// Bucket names — the single definition sites. The transactional paths
// resolve buckets by name inside their transactions, so the names travel
// as methods rather than pre-opened handles.
func (s *SNSStore) topicsBucketName() string            { return "sns-topics-" + s.region }
func (s *SNSStore) subscriptionsBucketName() string     { return "sns-subscriptions-" + s.region }
func (s *SNSStore) platformAppsBucketName() string      { return "sns-platform-apps-" + s.region }
func (s *SNSStore) platformEndpointsBucketName() string { return "sns-platform-endpoints-" + s.region }
func (s *SNSStore) appEndpointsIndexBucketName() string { return "sns-app-endpoints-index-" + s.region }
func (s *SNSStore) endpointTokenIndexBucketName() string {
	return "sns-endpoint-token-index-" + s.region
}
func (s *SNSStore) fifoDedupBucketName() string    { return "sns-fifo-dedup-" + s.region }
func (s *SNSStore) fifoSequenceBucketName() string { return "sns-fifo-sequence-" + s.region }

// cleanupExpiredDeduplicationEntries periodically removes deduplication
// entries whose window has closed, bounding the durable dedup bucket. An
// expired entry already reads as a miss, so the sweep is reclamation, not
// correctness.
func (s *SNSStore) cleanupExpiredDeduplicationEntries() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.reclaimDeduplicationEntries(s.scanExpiredDeduplicationEntries())
		}
	}
}

// scanExpiredDeduplicationEntries collects the deduplication keys whose
// window had already closed when the scan read them. The scan runs without
// the deduplication mutex — its result is advisory input to the reclaim
// step, which re-checks every key before deleting.
func (s *SNSStore) scanExpiredDeduplicationEntries() [][]byte {
	now := time.Now()
	var expired [][]byte
	err := s.fifoDedupStore.ForEach(func(k, v []byte) error {
		var record fifoDedupRecord
		if json.Unmarshal(v, &record) != nil || now.After(record.ExpiresAt) {
			expired = append(expired, append([]byte(nil), k...))
		}
		return nil
	})
	if err != nil {
		logs.Warn("SNS deduplication sweep scan failed", logs.Err(err))
		return nil
	}
	return expired
}

// reclaimDeduplicationEntries deletes the scanned keys that are STILL
// expired, re-reading each under the deduplication mutex: a publisher may
// have re-recorded a scanned key with a fresh window between the scan and
// the delete (CheckAndRecordFifoDeduplication runs under the same mutex),
// and deleting that fresh entry would breach the deduplication window the
// sweep exists to bound — a duplicate publish would deliver twice.
func (s *SNSStore) reclaimDeduplicationEntries(keys [][]byte) {
	for _, key := range keys {
		s.deduplicationMu.Lock()
		raw, err := s.fifoDedupStore.Get(key)
		if err == nil && raw != nil {
			var record fifoDedupRecord
			if json.Unmarshal(raw, &record) == nil && time.Now().After(record.ExpiresAt) {
				err = s.fifoDedupStore.Delete(key)
			}
		}
		s.deduplicationMu.Unlock()
		if err != nil {
			logs.Warn("SNS deduplication sweep reclamation failed", logs.Err(err))
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
	return s.arnBuilder.SNS().Subscription(topicArn, subscriptionID)
}

// CreateTopic creates a new SNS topic in the store, committing the record
// and its initial tag set as one storage transaction. The create is
// idempotent: an existing name returns the stored topic with the requested
// tags merged, so a retry whose first attempt failed between the record
// write and the tag write still lands its tags. tags ride the input only —
// the Topic record carries no tag state; the TagStore owns tags.
func (s *SNSStore) CreateTopic(topic *Topic, tags map[string]string) (*Topic, error) {
	topicArn := s.buildTopicArn(topic.Name)

	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	if s.Exists(topicArn) {
		var existingTopic Topic
		if err := s.BaseStore.Get(topicArn, &existingTopic); err != nil {
			return nil, err
		}
		if len(tags) > 0 {
			if err := s.TagStore.TagFromSlice(topicArn, tagsToSlice(tags)); err != nil {
				return nil, err
			}
		}
		return &existingTopic, nil
	}

	// Per-account topic quota for the topic's type, counted under the
	// create lock so a concurrent creator cannot slip past the bound. The
	// idempotent retry above does not consume quota — it adds no topic.
	// Counting walks the topics bucket; on this platform's scale that is a
	// bounded in-memory iteration.
	if err := s.checkTopicQuotaLocked(strings.HasSuffix(topic.Name, ".fifo")); err != nil {
		return nil, err
	}

	topic.Arn = topicArn
	topic.Owner = s.accountID
	if topic.Attributes == nil {
		topic.Attributes = make(map[string]string)
	}

	topicBytes, err := json.Marshal(topic)
	if err != nil {
		return nil, err
	}
	if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.topicsBucketName()).Put([]byte(topicArn), topicBytes); err != nil {
			return err
		}
		// The initial tag set commits with the topic itself: a tag failure
		// leaves no half-created topic for a retry to ride the idempotent
		// short-circuit against, and the create-time cap is the fresh
		// set's own size.
		return s.TagStore.TagInTxn(txn, topicArn, tags, common.MaxTagsPerResource)
	}); err != nil {
		return nil, err
	}

	return topic, nil
}

// checkTopicQuotaLocked enforces the documented per-account topic quota for
// the type being created — "Standard: 100,000 per account; FIFO: 1,000 per
// account" (SNS endpoints and quotas) — against the live topic count of this
// Region store. The caller must hold topicMu.
func (s *SNSStore) checkTopicQuotaLocked(fifoTopic bool) error {
	limit := MaxStandardTopicsPerAccount
	if fifoTopic {
		limit = MaxFifoTopicsPerAccount
	}
	count := 0
	marker := ""
	for {
		result, err := s.ListTopics(common.ListOptions{Marker: marker})
		if err != nil {
			return err
		}
		for _, t := range result.Items {
			if t.IsFifoTopic() == fifoTopic {
				count++
			}
		}
		if !result.IsTruncated || result.NextMarker == "" {
			break
		}
		marker = result.NextMarker
	}
	if count >= limit {
		return ErrTopicLimitExceeded
	}
	return nil
}

// tagsToSlice converts a tag map to the TagStore's slice form.
func tagsToSlice(tags map[string]string) []types.Tag {
	tagSlice := make([]types.Tag, 0, len(tags))
	for k, v := range tags {
		tagSlice = append(tagSlice, types.Tag{Key: k, Value: v})
	}
	return tagSlice
}

// GetTopic retrieves an SNS topic by its ARN from the store.
// Returns the topic or an error if not found.
func (s *SNSStore) GetTopic(topicArn string) (*Topic, error) {
	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return nil, notFoundOr(err, ErrTopicNotFound)
	}
	return &topic, nil
}

// DeleteTopic deletes an SNS topic by its ARN from the store.
// Returns an error if the topic does not exist.
func (s *SNSStore) DeleteTopic(topicArn string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	// The existence read is a Get, not an Exists probe: a storage failure
	// must surface as itself, not join the miss branch's 404.
	var existing Topic
	if err := s.BaseStore.Get(topicArn, &existing); err != nil {
		return notFoundOr(err, ErrTopicNotFound)
	}

	// Every subscription leaves with the topic, through the all-pages
	// walk the fan-out seam owns (a single list page holds at most
	// DefaultMaxItems) — the 101st subscriber's records cannot survive
	// the delete as orphans.
	subs, err := s.ListAllSubscriptionsByTopic(topicArn)
	if err != nil {
		return fmt.Errorf("listing topic subscriptions: %w", err)
	}

	// The subscriptions, the topic's tags, the topic record and its FIFO
	// state — the topic's deduplication-window entries and message-group
	// sequence counters — leave storage as one transaction, so a
	// reported-successful delete leaves no state a re-created topic name
	// would inherit and no orphan a retry must self-heal. A publish racing
	// the delete can still write a dedup entry after this transaction
	// commits; such an orphan self-heals (expired entries read as misses
	// and the sweep reclaims them), and an orphaned sequence counter for
	// the same key is the monotonic continuation a re-created topic wants.
	// The sweep prefix carries the same NUL separator the keys join with.
	prefix := []byte(topicArn + "\x00")
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		for _, sub := range subs {
			if err := txn.Bucket(s.subscriptionsBucketName()).Delete([]byte(sub.SubscriptionArn)); err != nil {
				return err
			}
		}
		if err := s.TagStore.DeleteInTxn(txn, topicArn); err != nil {
			return err
		}
		if err := deletePrefixed(txn.Bucket(s.fifoDedupBucketName()), prefix); err != nil {
			return err
		}
		if err := deletePrefixed(txn.Bucket(s.fifoSequenceBucketName()), prefix); err != nil {
			return err
		}
		return txn.Bucket(s.topicsBucketName()).Delete([]byte(topicArn))
	})
}

// deletePrefixed deletes every key carrying the prefix inside the
// transaction, collecting the keys before deleting them so the iterator is
// not mutated mid-walk.
func deletePrefixed(bucket storage.Bucket, prefix []byte) error {
	var keys [][]byte
	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		keys = append(keys, append([]byte(nil), iter.Key()...))
	}
	err := iter.Error()
	iter.Close()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := bucket.Delete(key); err != nil {
			return err
		}
	}
	return nil
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
		return notFoundOr(err, ErrTopicNotFound)
	}

	topic.Attributes = mergeAttrs(topic.Attributes, attributes)

	return s.Put(topic.Arn, &topic)
}

// markSubscriptionPending records the create transition on the topic's
// counters: one more subscription awaits confirmation.
func markSubscriptionPending(topic *Topic) {
	topic.SubscriptionsPending++
}

// markSubscriptionConfirmed records the confirm transition on the topic's
// counters: the subscription leaves pending and joins confirmed, clamped at
// zero so a record confirmed twice cannot drive the pending count negative.
func markSubscriptionConfirmed(topic *Topic) {
	if topic.SubscriptionsPending > 0 {
		topic.SubscriptionsPending--
	}
	topic.SubscriptionsConfirmed++
}

// markSubscriptionRemoved records the delete transition on the topic's
// counters: the subscription leaves whichever live counter held it (pending
// or confirmed, clamped at zero) and joins the deleted count.
func markSubscriptionRemoved(topic *Topic, wasPending bool) {
	if wasPending {
		if topic.SubscriptionsPending > 0 {
			topic.SubscriptionsPending--
		}
	} else if topic.SubscriptionsConfirmed > 0 {
		topic.SubscriptionsConfirmed--
	}
	topic.SubscriptionsDeleted++
}

// CreateSubscription creates a new SNS subscription in the store. The call
// is idempotent on the subscription's natural key: when a live subscription
// with the same TopicArn+Protocol+Endpoint+Owner already exists it is
// returned unchanged with created=false, so a repeated Subscribe never
// accumulates duplicate deliveries. The AWS API reference documents no
// duplicate behaviour for Subscribe; the return-existing outcome follows
// the model, which declares no already-exists error for the operation (its
// error set carries only authorisation, not-found, validation, internal and
// limit errors) and the documented CreateTopic idempotency pattern
// ("This action is idempotent, so if the requester already owns a topic
// with the specified name, that topic's ARN is returned without creating a
// new topic"). Returns the subscription, whether it was newly created, or
// an error.
func (s *SNSStore) CreateSubscription(subscription *Subscription) (*Subscription, bool, error) {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err != nil {
		return nil, false, notFoundOr(err, ErrTopicNotFound)
	}

	existing, err := s.findSubscriptionByKeyLocked(topic.Arn, subscription.Protocol, subscription.Endpoint, subscription.Owner)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}

	if err := s.checkSubscriptionQuotaLocked(&topic, subscription); err != nil {
		return nil, false, err
	}

	subscriptionID := uuid.New().String()
	subscriptionArn := s.buildSubscriptionArn(subscription.TopicArn, subscriptionID)
	subscription.SubscriptionArn = subscriptionArn
	subscription.TopicArn = topic.Arn
	subscription.PendingConfirmation = true
	subscription.ConfirmationToken = uuid.New().String()
	if subscription.Attributes == nil {
		subscription.Attributes = make(map[string]string)
	}

	subscriptionBytes, err := json.Marshal(subscription)
	if err != nil {
		return nil, false, err
	}
	markSubscriptionPending(&topic)
	topicBytes, err := json.Marshal(&topic)
	if err != nil {
		return nil, false, err
	}
	// The subscription record and the topic's pending counter commit as
	// one transaction — a mid-sequence storage failure leaves no orphaned
	// pending subscription for a retry to duplicate.
	if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.subscriptionsBucketName()).Put([]byte(subscriptionArn), subscriptionBytes); err != nil {
			return err
		}
		return txn.Bucket(s.topicsBucketName()).Put([]byte(topic.Arn), topicBytes)
	}); err != nil {
		return nil, false, err
	}

	return subscription, true, nil
}

// findSubscriptionByKeyLocked walks the topic's subscriptions for a live
// record matching the subscription natural key (protocol, endpoint, owner).
// Returns nil when no such subscription exists. The caller must hold
// topicMu and subscriptionMu.
func (s *SNSStore) findSubscriptionByKeyLocked(topicArn, protocol, endpoint, owner string) (*Subscription, error) {
	subs, err := s.ListAllSubscriptionsByTopic(topicArn)
	if err != nil {
		return nil, err
	}
	for _, sub := range subs {
		if sub.Protocol == protocol && sub.Endpoint == endpoint && sub.Owner == owner {
			return sub, nil
		}
	}
	return nil, nil
}

// checkSubscriptionQuotaLocked enforces the documented per-topic quotas at
// subscription creation: "Subscriptions — Standard: 12,500,000 per topic;
// FIFO: 100 per topic" and "Subscription filter policies — 200 filter
// policies per topic" (SNS endpoints and quotas), the latter counted only
// when the new subscription itself carries a filter policy. The caller must
// hold topicMu and subscriptionMu.
func (s *SNSStore) checkSubscriptionQuotaLocked(topic *Topic, subscription *Subscription) error {
	limit := MaxStandardSubscriptionsPerTopic
	if topic.IsFifoTopic() {
		limit = MaxFifoSubscriptionsPerTopic
	}
	subs, err := s.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		return err
	}
	if len(subs) >= limit {
		return ErrSubscriptionLimitExceeded
	}
	if strings.TrimSpace(subscription.Attributes[AttrFilterPolicy]) != "" {
		filtered := 0
		for _, sub := range subs {
			if strings.TrimSpace(sub.Attributes[AttrFilterPolicy]) != "" {
				filtered++
			}
		}
		if filtered >= MaxFilterPoliciesPerTopic {
			return ErrFilterPolicyLimitExceeded
		}
	}
	// "For Firehose delivery streams, 5 per topic, per subscription
	// owner" (SNS endpoints and quotas) — counted only when the new
	// subscription itself rides the firehose protocol; the platform is
	// single-account, so the per-topic count is the per-owner count.
	if subscription.Protocol == "firehose" {
		firehose := 0
		for _, sub := range subs {
			if sub.Protocol == "firehose" {
				firehose++
			}
		}
		if firehose >= MaxFirehoseSubscriptionsPerTopic {
			return ErrSubscriptionLimitExceeded
		}
	}
	return nil
}

// GetSubscription retrieves an SNS subscription by its ARN from the store.
// Returns the subscription or an error if not found.
func (s *SNSStore) GetSubscription(subscriptionArn string) (*Subscription, error) {
	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return nil, notFoundOr(err, ErrSubscriptionNotFound)
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
		return notFoundOr(err, ErrSubscriptionNotFound)
	}

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err != nil {
		logs.Warn("sns: topic not found during subscription deletion (concurrently deleted?)",
			logs.String("topicArn", subscription.TopicArn),
			logs.String("subscriptionArn", subscriptionArn))
		// The topic is gone; deleting the subscription record alone is the
		// complete remaining work.
		return s.topicSubscriptionsStore.Delete(subscriptionArn)
	}

	markSubscriptionRemoved(&topic, subscription.PendingConfirmation)
	topicBytes, err := json.Marshal(&topic)
	if err != nil {
		return err
	}
	// The counter update and the record delete commit as one transaction;
	// a failure fails the whole call instead of leaving a skewed counter
	// behind a deleted subscription.
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.topicsBucketName()).Put([]byte(topic.Arn), topicBytes); err != nil {
			return err
		}
		return txn.Bucket(s.subscriptionsBucketName()).Delete([]byte(subscriptionArn))
	})
}

// AutoConfirmSubscription confirms a subscription without requiring a
// token. It applies to the protocols the service's protocol registry
// auto-confirms at Subscribe time (sqs and lambda); http/https instead go
// through sendSubscriptionConfirmation and ConfirmSubscription. It updates
// the subscription status and adjusts the topic counters.
func (s *SNSStore) AutoConfirmSubscription(subscription *Subscription) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	subscription.PendingConfirmation = false
	subscription.ConfirmationWasAuthenticated = true
	// The token stays recorded: "Confirmation tokens are valid for two
	// days" (Subscribe), so a repeated ConfirmSubscription with the same
	// token within the window succeeds — confirmation is idempotent, and
	// clearing the token here would turn the documented retry into an
	// InvalidParameter.

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err != nil {
		logs.Warn("sns: topic not found during auto-confirm (concurrently deleted?)",
			logs.String("topicArn", subscription.TopicArn),
			logs.String("subscriptionArn", subscription.SubscriptionArn))
		// The topic is gone; the subscription state alone is the complete
		// remaining work.
		return s.topicSubscriptionsStore.Put(subscription.SubscriptionArn, subscription)
	}

	markSubscriptionConfirmed(&topic)
	subscriptionBytes, err := json.Marshal(subscription)
	if err != nil {
		return err
	}
	topicBytes, err := json.Marshal(&topic)
	if err != nil {
		return err
	}
	// The confirmed state and the topic counters commit as one
	// transaction; a failure fails the whole call instead of leaving a
	// confirmed subscription behind a stale counter.
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.subscriptionsBucketName()).Put([]byte(subscription.SubscriptionArn), subscriptionBytes); err != nil {
			return err
		}
		return txn.Bucket(s.topicsBucketName()).Put([]byte(topic.Arn), topicBytes)
	})
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

// ListAllSubscriptionsByTopic returns every subscription of a topic,
// walking ALL pages. ListSubscriptionsByTopic caps a single page at
// DefaultMaxItems (100); delivery fan-out and inventory paths that must
// see the complete set walk pages through this method so subscription
// 101+ is never silently dropped.
func (s *SNSStore) ListAllSubscriptionsByTopic(topicArn string) ([]*Subscription, error) {
	var all []*Subscription
	marker := ""
	for {
		result, err := s.ListSubscriptionsByTopic(topicArn, common.ListOptions{Marker: marker})
		if err != nil {
			return nil, err
		}
		all = append(all, result.Items...)
		if !result.IsTruncated || result.NextMarker == "" {
			return all, nil
		}
		marker = result.NextMarker
	}
}

// FindSubscriptionByToken finds a subscription for the given topic whose
// confirmation token matches. Unlike ListSubscriptionsByTopic (which is
// capped at DefaultMaxItems), this scans all subscriptions to ensure tokens
// are found even when a topic has more than 100 pending subscriptions.
func (s *SNSStore) FindSubscriptionByToken(topicArn, token string) (*Subscription, error) {
	sub, err := common.FindFirst[Subscription](s.topicSubscriptionsStore, func(sub *Subscription) bool {
		return sub.TopicArn == topicArn && sub.ConfirmationToken == token
	})
	if err != nil {
		// The scan answers a miss and a failing read with distinct
		// wrappers; only the miss is the not-found sentinel, so a storage
		// failure surfaces as itself instead of a phantom miss the
		// service layer would answer with a client error.
		return nil, notFoundOr(err, ErrSubscriptionNotFound)
	}
	return sub, nil
}

// ConfirmSubscription confirms an SNS subscription using the subscription ARN
// and confirmation token. The ConfirmSubscription API call is an AWS-signed
// request, so the confirmed subscription's ConfirmationWasAuthenticated
// attribute is always true (it reports whether the confirmation request was
// authenticated, per the GetSubscriptionAttributes documentation).
// authenticateOnUnsubscribe carries the separate ConfirmSubscription input
// flag of the same name: when non-nil it is persisted as the
// "AuthenticateOnUnsubscribe" subscription attribute, which Unsubscribe
// enforces by requiring an authenticated request from the topic or
// subscription owner.
func (s *SNSStore) ConfirmSubscription(subscriptionArn, token string, authenticateOnUnsubscribe *bool) (*Subscription, error) {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return nil, notFoundOr(err, ErrSubscriptionNotFound)
	}

	if subscription.ConfirmationToken != token {
		return nil, ErrInvalidToken
	}

	subscription.PendingConfirmation = false
	subscription.ConfirmationWasAuthenticated = true
	// The token stays recorded: "Confirmation tokens are valid for two
	// days" (Subscribe), so a repeated ConfirmSubscription with the same
	// token within the window succeeds — confirmation is idempotent, and
	// clearing the token here would turn the documented retry into an
	// InvalidParameter.

	if authenticateOnUnsubscribe != nil {
		if subscription.Attributes == nil {
			subscription.Attributes = make(map[string]string)
		}
		if *authenticateOnUnsubscribe {
			subscription.Attributes[AttrAuthenticateOnUnsubscribe] = "true"
		} else {
			subscription.Attributes[AttrAuthenticateOnUnsubscribe] = "false"
		}
	}

	var topic Topic
	if err := s.BaseStore.Get(subscription.TopicArn, &topic); err != nil {
		logs.Warn("sns: topic not found during confirm (concurrently deleted?)",
			logs.String("topicArn", subscription.TopicArn),
			logs.String("subscriptionArn", subscriptionArn))
		// The topic is gone; the confirmed subscription state alone is the
		// complete remaining work.
		return &subscription, s.topicSubscriptionsStore.Put(subscriptionArn, &subscription)
	}

	markSubscriptionConfirmed(&topic)
	subscriptionBytes, err := json.Marshal(&subscription)
	if err != nil {
		return nil, err
	}
	topicBytes, err := json.Marshal(&topic)
	if err != nil {
		return nil, err
	}
	// The confirmed state and the topic counters commit as one
	// transaction; a failure fails the whole call instead of leaving a
	// confirmed subscription behind a stale counter.
	if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.subscriptionsBucketName()).Put([]byte(subscriptionArn), subscriptionBytes); err != nil {
			return err
		}
		return txn.Bucket(s.topicsBucketName()).Put([]byte(topic.Arn), topicBytes)
	}); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// SetSubscriptionAttributes sets the attributes for an SNS subscription.
func (s *SNSStore) SetSubscriptionAttributes(subscriptionArn string, attributes map[string]string) error {
	s.subscriptionMu.Lock()
	defer s.subscriptionMu.Unlock()

	var subscription Subscription
	if err := s.topicSubscriptionsStore.Get(subscriptionArn, &subscription); err != nil {
		return notFoundOr(err, ErrSubscriptionNotFound)
	}

	if v, ok := attributes[AttrFilterPolicy]; ok && strings.TrimSpace(v) != "" {
		// Attaching a filter policy to a subscription that carries none
		// consumes one of the topic's "200 filter policies per topic"
		// (SNS endpoints and quotas); updating an existing policy does
		// not. The count runs under the same lock family as
		// CreateSubscription's quota check, so the two paths cannot admit
		// past the bound concurrently.
		if strings.TrimSpace(subscription.Attributes[AttrFilterPolicy]) == "" {
			subs, err := s.ListAllSubscriptionsByTopic(subscription.TopicArn)
			if err != nil {
				return err
			}
			filtered := 0
			for _, sub := range subs {
				if strings.TrimSpace(sub.Attributes[AttrFilterPolicy]) != "" {
					filtered++
				}
			}
			if filtered >= MaxFilterPoliciesPerTopic {
				return ErrFilterPolicyLimitExceeded
			}
		}
	}

	subscription.Attributes = mergeAttrs(subscription.Attributes, attributes)

	return s.topicSubscriptionsStore.Put(subscriptionArn, &subscription)
}

// GetSubscriptionAttributes retrieves the attributes for an SNS subscription.
// The response carries the stored attribute map plus the synthesised
// documented keys (SubscriptionArn, TopicArn, Owner, PendingConfirmation,
// ConfirmationWasAuthenticated); Protocol and Endpoint are not attribute-plane
// keys — the documented GetSubscriptionAttributes enumeration does not list
// them, and the subscription's identity already reaches clients through
// ListSubscriptions and the envelope.
func (s *SNSStore) GetSubscriptionAttributes(subscriptionArn string) (map[string]string, error) {
	subscription, err := s.GetSubscription(subscriptionArn)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	for k, v := range subscription.Attributes {
		// AuthenticateOnUnsubscribe is an internal enforcement flag set
		// through the ConfirmSubscription input parameter; it is not part
		// of the documented GetSubscriptionAttributes response set.
		if k == AttrAuthenticateOnUnsubscribe {
			continue
		}
		attrs[k] = v
	}
	attrs[AttrSubscriptionArn] = subscription.SubscriptionArn
	attrs[AttrTopicArn] = subscription.TopicArn
	attrs[AttrOwner] = subscription.Owner
	attrs[AttrPendingConfirmation] = "false"
	if subscription.PendingConfirmation {
		attrs[AttrPendingConfirmation] = "true"
	}
	attrs[AttrConfirmationWasAuthenticated] = "false"
	if subscription.ConfirmationWasAuthenticated {
		attrs[AttrConfirmationWasAuthenticated] = "true"
	}

	return attrs, nil
}

// ListTagsForResource retrieves the tags associated with an SNS resource.
func (s *SNSStore) ListTagsForResource(resourceArn string) ([]types.Tag, error) {
	return s.TagStore.ListAsSlice(resourceArn)
}

// TagResource adds tags to an SNS resource.
func (s *SNSStore) Tag(resourceArn string, tags []types.Tag) error {
	return s.TagStore.TagFromSlice(resourceArn, tags)
}

// UntagResource removes tags from an SNS resource.
func (s *SNSStore) Untag(resourceArn string, tagKeys []string) error {
	return s.TagStore.Untag(resourceArn, tagKeys)
}

// AcquireFifoGroupLock serialises FIFO publishing per topic message group
// and returns the release function. The caller holds this lock across
// sequence allocation and synchronous delivery, which makes sequence order
// and delivery order the same order — the per-group FIFO guarantee.
func (s *SNSStore) AcquireFifoGroupLock(topicArn, messageGroupId string) func() {
	key := fifoStateKey(topicArn, messageGroupId)

	s.fifoGroupMu.Lock()
	entry, ok := s.fifoGroupLocks[key]
	if !ok {
		entry = &sync.Mutex{}
		s.fifoGroupLocks[key] = entry
	}
	s.fifoGroupMu.Unlock()

	entry.Lock()
	return entry.Unlock
}

// fifoStateKey joins the FIFO state's key components (deduplication
// windows, sequence counters, group locks). The separator is NUL: the
// identifiers may contain colons (their documented punctuation set
// includes one), so a flat colon concatenation would conflate distinct
// component tuples into one key — a shared window or counter.
func fifoStateKey(parts ...string) string {
	return strings.Join(parts, "\x00")
}

// fifoDedupKey builds the deduplication-window key. The documented scope
// rules the key encodes: "Message deduplication applies to an entire Amazon
// SNS FIFO topic when the topic attribute FifoThroughputScope is set to
// Topic. When the topic attribute FifoThroughputScope is set to
// MessageGroup, message deduplication applies to each individual message
// group" (message deduplication for FIFO topics).
func fifoDedupKey(topicArn, messageGroupId, messageDeduplicationId string, perGroupDedup bool) string {
	if perGroupDedup {
		return fifoStateKey(topicArn, messageGroupId, messageDeduplicationId)
	}
	return fifoStateKey(topicArn, messageDeduplicationId)
}

// dedupRecord reads one deduplication-window entry from the durable store;
// a missing, expired or empty identifier is a miss (nil, nil). The key's
// NUL separator makes the empty-key check the only shape guard — a
// trailing-colon check would wrongly drop the legal identifiers that end
// in a colon.
func (s *SNSStore) dedupRecord(dedupKey string) (*fifoDedupRecord, error) {
	if dedupKey == "" {
		return nil, nil
	}
	raw, err := s.fifoDedupStore.Get([]byte(dedupKey))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	var record fifoDedupRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, err
	}
	if !time.Now().Before(record.ExpiresAt) {
		return nil, nil
	}
	return &record, nil
}

// CheckFifoDeduplication reports whether an unexpired deduplication-window
// entry exists for the identifier at the topic's deduplication scope,
// returning the original publish's MessageId and SequenceNumber when it
// does.
func (s *SNSStore) CheckFifoDeduplication(topicArn, messageGroupId, messageDeduplicationId string, perGroupDedup bool) (messageID, sequenceNumber string, hit bool, err error) {
	if messageDeduplicationId == "" {
		return "", "", false, nil
	}
	record, err := s.dedupRecord(fifoDedupKey(topicArn, messageGroupId, messageDeduplicationId, perGroupDedup))
	if err != nil || record == nil {
		return "", "", false, err
	}
	return record.MessageID, record.SequenceNumber, true, nil
}

// CheckAndRecordFifoDeduplication atomically records the deduplication
// entry — the publish's identifiers and the window expiry — or returns the
// existing entry's identifiers when an unexpired one is already recorded.
// The mutex makes check-and-record one critical section, eliminating the
// TOCTOU race between separate check and record calls; publishes in
// different message groups can share a deduplication key (topic scope) and
// so cannot rely on the group lock for this atomicity.
func (s *SNSStore) CheckAndRecordFifoDeduplication(topicArn, messageGroupId, messageDeduplicationId, messageID, sequenceNumber string, perGroupDedup bool) (string, string, bool, error) {
	if messageDeduplicationId == "" {
		return "", "", false, nil
	}
	key := []byte(fifoDedupKey(topicArn, messageGroupId, messageDeduplicationId, perGroupDedup))

	s.deduplicationMu.Lock()
	defer s.deduplicationMu.Unlock()

	if record, err := s.dedupRecord(string(key)); err != nil {
		return "", "", false, err
	} else if record != nil {
		return record.MessageID, record.SequenceNumber, true, nil
	}

	value, err := json.Marshal(&fifoDedupRecord{
		MessageID:      messageID,
		SequenceNumber: sequenceNumber,
		ExpiresAt:      time.Now().Add(deduplicationWindow),
	})
	if err != nil {
		return "", "", false, err
	}
	if err := s.fifoDedupStore.Put(key, value); err != nil {
		return "", "", false, err
	}
	return "", "", false, nil
}

// AllocateFifoSequence returns the message group's next sequence number,
// persisting the counter so SequenceNumber keeps increasing across process
// restarts. AWS SNS assigns large, non-consecutive numbers: the first
// allocation for a group starts from a large random base and every
// allocation adds a random increment, producing numbers that increase
// monotonically within a group without being sequential.
func (s *SNSStore) AllocateFifoSequence(topicArn, messageGroupId string) (string, error) {
	key := []byte(fifoStateKey(topicArn, messageGroupId))

	s.sequenceMu.Lock()
	defer s.sequenceMu.Unlock()

	record := fifoSequenceRecord{
		// Start from a large base to emulate AWS's 19-digit numbers.
		Value: 1000000000000000000 + int64(rand.Intn(89999999999999999)),
	}
	if raw, err := s.fifoSequenceStore.Get(key); err != nil {
		return "", err
	} else if raw != nil {
		if err := json.Unmarshal(raw, &record); err != nil {
			return "", err
		}
	}
	// A random increment (1-1000) keeps the numbers non-consecutive.
	record.Value += int64(rand.Intn(1000)) + 1

	value, err := json.Marshal(&record)
	if err != nil {
		return "", err
	}
	if err := s.fifoSequenceStore.Put(key, value); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", record.Value), nil
}

// GetDataProtectionPolicy retrieves the data protection policy for a topic.
func (s *SNSStore) GetDataProtectionPolicy(topicArn string) (string, error) {
	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return "", notFoundOr(err, ErrTopicNotFound)
	}
	return topic.GetDataProtectionPolicy(), nil
}

// PutDataProtectionPolicy sets the data protection policy for a topic.
func (s *SNSStore) PutDataProtectionPolicy(topicArn, policy string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return notFoundOr(err, ErrTopicNotFound)
	}

	if topic.Attributes == nil {
		topic.Attributes = make(map[string]string)
	}
	topic.Attributes[AttrDataProtectionPolicy] = policy
	return s.Put(topicArn, &topic)
}

// AddPermission adds a permission to a topic.
func (s *SNSStore) AddPermission(topicArn string, permission *Permission) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return notFoundOr(err, ErrTopicNotFound)
	}

	if topic.Permissions == nil {
		topic.Permissions = []Permission{}
	}

	for _, p := range topic.Permissions {
		if p.Label == permission.Label {
			// Label is "A unique identifier for the new policy statement"
			// (AddPermission model documentation) — a duplicate violates
			// the parameter's constraint instead of overwriting the
			// existing statement.
			return ErrPermissionLabelExists
		}
	}

	topic.Permissions = append(topic.Permissions, *permission)
	return s.Put(topicArn, &topic)
}

// RemovePermission removes a permission from a topic by label.
func (s *SNSStore) RemovePermission(topicArn, label string) error {
	s.topicMu.Lock()
	defer s.topicMu.Unlock()

	var topic Topic
	if err := s.BaseStore.Get(topicArn, &topic); err != nil {
		return notFoundOr(err, ErrTopicNotFound)
	}

	for i, p := range topic.Permissions {
		if p.Label == label {
			topic.Permissions = append(topic.Permissions[:i], topic.Permissions[i+1:]...)
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
		return nil, notFoundOr(err, ErrPlatformApplicationNotFound)
	}
	return &app, nil
}

// DeletePlatformApplication deletes an SNS platform application and all its endpoints.
func (s *SNSStore) DeletePlatformApplication(arn string) error {
	s.platformAppMu.Lock()
	defer s.platformAppMu.Unlock()

	// The existence read is a Get, not an Exists probe: a storage failure
	// must surface as itself, not join the miss branch's 404.
	var existingApp PlatformApplication
	if err := s.platformApplicationsStore.Get(arn, &existingApp); err != nil {
		return notFoundOr(err, ErrPlatformApplicationNotFound)
	}

	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	// Collect endpoint keys from the reverse index before mutating.
	// Mutating a lazy iterator's backing store during iteration is unsafe;
	// collect-then-delete follows the BaseStore.DeleteByPrefix pattern.
	prefix := []byte(arn + "\x00")
	var idxKeys [][]byte
	var endpointArns []string
	var tokenKeys [][]byte
	iter := s.platformAppEndpointsIndex.ScanPrefix(prefix)
	for iter.Next() {
		raw := iter.Key()
		copied := make([]byte, len(raw))
		copy(copied, raw)
		idxKeys = append(idxKeys, copied)
		epArn := string(raw[len(prefix):])
		endpointArns = append(endpointArns, epArn)
		// Load the endpoint to obtain its token for token-index cleanup; a
		// nil entry marks an endpoint with no token to clean.
		var ep PlatformEndpoint
		if err := s.platformEndpointsStore.Get(epArn, &ep); err == nil && ep.Token != "" {
			tokenKeys = append(tokenKeys, []byte(arn+"\x00"+ep.Token))
		} else {
			tokenKeys = append(tokenKeys, nil)
		}
	}
	if err := iter.Error(); err != nil {
		iter.Close()
		return fmt.Errorf("scanning endpoint index: %w", err)
	}
	iter.Close()

	// Every endpoint record, every index entry, the application's tags and
	// the application record delete as one transaction — a failure deletes
	// nothing and the call stays retryable.
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := s.TagStore.DeleteInTxn(txn, arn); err != nil {
			return err
		}
		endpointsBucket := txn.Bucket(s.platformEndpointsBucketName())
		appIndexBucket := txn.Bucket(s.appEndpointsIndexBucketName())
		tokenIndexBucket := txn.Bucket(s.endpointTokenIndexBucketName())
		for i, epArn := range endpointArns {
			if err := endpointsBucket.Delete([]byte(epArn)); err != nil {
				return err
			}
			if err := appIndexBucket.Delete(idxKeys[i]); err != nil {
				return err
			}
			if tokenKeys[i] != nil {
				if err := tokenIndexBucket.Delete(tokenKeys[i]); err != nil {
					return err
				}
			}
		}
		return txn.Bucket(s.platformAppsBucketName()).Delete([]byte(arn))
	})
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
		return notFoundOr(err, ErrPlatformApplicationNotFound)
	}

	app.Attributes = mergeAttrs(app.Attributes, attributes)

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

// CreatePlatformEndpoint creates a new SNS platform endpoint, or merges
// into the existing endpoint when the application already holds one with
// the same token. The application check and every write run under one
// platformAppMu+platformEndpointMu section, so a concurrent
// DeletePlatformApplication cannot interleave between the check and the
// writes. The endpoint record and both index entries commit as one
// transaction — a failed index write leaves no endpoint record behind.
// The store owns endpoint ARN generation; no caller supplies an ARN.
func (s *SNSStore) CreatePlatformEndpoint(endpoint *PlatformEndpoint) (*PlatformEndpoint, error) {
	s.platformAppMu.RLock()
	defer s.platformAppMu.RUnlock()
	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	// The existence read is a Get, not an Exists probe: a storage failure
	// must surface as itself, not join the miss branch's 404.
	var parentApp PlatformApplication
	if err := s.platformApplicationsStore.Get(endpoint.PlatformApplicationArn, &parentApp); err != nil {
		return nil, notFoundOr(err, ErrPlatformApplicationNotFound)
	}

	tokenKey := endpoint.PlatformApplicationArn + "\x00" + endpoint.Token
	// An index read failure is not an absent token: proceeding would mint
	// a second endpoint for a token the index still maps to the first.
	existingArnBytes, err := s.endpointTokenIndex.Get([]byte(tokenKey))
	if err != nil {
		return nil, err
	}
	if len(existingArnBytes) > 0 {
		existingArn := string(existingArnBytes)
		var existing PlatformEndpoint
		if err := s.platformEndpointsStore.Get(existingArn, &existing); err == nil {
			mergeEndpointAttrs(&existing, endpoint)
			existingBytes, merr := json.Marshal(&existing)
			if merr != nil {
				return nil, merr
			}
			idxKey := endpoint.PlatformApplicationArn + "\x00" + existingArn
			// Merge-time self-heal: an app-index entry missing here is
			// backfilled on the merge — the single index-reconciliation
			// mechanism (construction-time backfill is prohibited as
			// migration code).
			needBackfill := !s.platformAppEndpointsIndex.Has([]byte(idxKey))
			if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
				if err := txn.Bucket(s.platformEndpointsBucketName()).Put([]byte(existingArn), existingBytes); err != nil {
					return err
				}
				if needBackfill {
					return txn.Bucket(s.appEndpointsIndexBucketName()).Put([]byte(idxKey), []byte("1"))
				}
				return nil
			}); err != nil {
				return nil, err
			}
			return &existing, nil
		}
		// A stale token-index entry (its endpoint record is gone) falls
		// through to the create path, whose token-index write overwrites it.
	}

	endpointID := uuid.New().String()
	// The Core's ARN validation guarantees the app/<platform>/<name>
	// resource form, so the extraction cannot fail on a production path.
	platform, appName := svcarn.ExtractPlatformApplicationFromARN(endpoint.PlatformApplicationArn)
	endpoint.EndpointArn = s.arnBuilder.SNS().PlatformEndpoint(platform, appName, endpointID)

	if endpoint.Attributes == nil {
		endpoint.Attributes = make(map[string]string)
	}
	if _, ok := endpoint.Attributes[AttrEnabled]; !ok {
		endpoint.Attributes[AttrEnabled] = "true"
	}

	endpointBytes, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}
	idxKey := endpoint.PlatformApplicationArn + "\x00" + endpoint.EndpointArn
	if err := s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.platformEndpointsBucketName()).Put([]byte(endpoint.EndpointArn), endpointBytes); err != nil {
			return err
		}
		if err := txn.Bucket(s.appEndpointsIndexBucketName()).Put([]byte(idxKey), []byte("1")); err != nil {
			return err
		}
		return txn.Bucket(s.endpointTokenIndexBucketName()).Put([]byte(tokenKey), []byte(endpoint.EndpointArn))
	}); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// mergeEndpointAttrs merges CustomUserData and Attributes from src into dst.
func mergeEndpointAttrs(dst, src *PlatformEndpoint) {
	if src.CustomUserData != "" {
		dst.CustomUserData = src.CustomUserData
	}
	dst.Attributes = mergeAttrs(dst.Attributes, src.Attributes)
}

// GetEndpoint retrieves an SNS platform endpoint by its ARN.
func (s *SNSStore) GetEndpoint(arn string) (*PlatformEndpoint, error) {
	var endpoint PlatformEndpoint
	if err := s.platformEndpointsStore.Get(arn, &endpoint); err != nil {
		return nil, notFoundOr(err, ErrEndpointNotFound)
	}
	return &endpoint, nil
}

// DeleteEndpoint deletes an SNS platform endpoint by its ARN.
func (s *SNSStore) DeleteEndpoint(arn string) error {
	s.platformEndpointMu.Lock()
	defer s.platformEndpointMu.Unlock()

	var endpoint PlatformEndpoint
	if err := s.platformEndpointsStore.Get(arn, &endpoint); err != nil {
		return notFoundOr(err, ErrEndpointNotFound)
	}

	idxKey := endpoint.PlatformApplicationArn + "\x00" + arn
	tokenKey := endpoint.PlatformApplicationArn + "\x00" + endpoint.Token
	// The record and both index entries delete as one transaction — no
	// stale index entry survives a reported success, and a failure deletes
	// nothing, keeping the call retryable.
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(s.platformEndpointsBucketName()).Delete([]byte(arn)); err != nil {
			return err
		}
		if err := txn.Bucket(s.appEndpointsIndexBucketName()).Delete([]byte(idxKey)); err != nil {
			return err
		}
		return txn.Bucket(s.endpointTokenIndexBucketName()).Delete([]byte(tokenKey))
	})
}

// GetEndpointAttributes retrieves the attributes of an SNS platform endpoint.
func (s *SNSStore) GetEndpointAttributes(arn string) (map[string]string, error) {
	endpoint, err := s.GetEndpoint(arn)
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]string)
	attrs[AttrToken] = endpoint.Token
	if endpoint.CustomUserData != "" {
		// CustomUserData set through CreatePlatformEndpoint's member (not
		// the Attributes map) must still answer GetEndpointAttributes —
		// "Attributes include the following: CustomUserData – arbitrary
		// user data to associate with the endpoint"
		// (GetEndpointAttributes). A value later written through
		// SetEndpointAttributes lives in the Attributes map and overwrites
		// this one below.
		attrs[AttrCustomUserData] = endpoint.CustomUserData
	}
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
		return notFoundOr(err, ErrEndpointNotFound)
	}

	endpoint.Attributes = mergeAttrs(endpoint.Attributes, attributes)

	return s.platformEndpointsStore.Put(arn, &endpoint)
}

// ListEndpointsByPlatformApplication returns all endpoints for an SNS platform application.
// Uses the reverse index bucket for O(k) lookup where k is the number of endpoints
// for this application, instead of scanning all endpoints.
func (s *SNSStore) ListEndpointsByPlatformApplication(platformAppArn string, opts common.ListOptions) (*common.ListResult[PlatformEndpoint], error) {
	// The page-size clamp is the shared one; the marker walk itself is
	// substrate-bound (the reverse index bucket yields raw keys, not the
	// JSON records common.List pages over).
	opts = common.NormalizeListOptions(opts)

	s.platformAppMu.RLock()
	defer s.platformAppMu.RUnlock()
	s.platformEndpointMu.RLock()
	defer s.platformEndpointMu.RUnlock()

	prefix := []byte(platformAppArn + "\x00")
	var items []*PlatformEndpoint
	var lastEpArn string
	count := 0
	hasMore := false
	started := opts.Marker == ""

	iter := s.platformAppEndpointsIndex.ScanPrefix(prefix)
	defer iter.Close()
	for iter.Next() {
		epArn := string(iter.Key()[len(prefix):])

		if !started {
			if epArn == opts.Marker {
				started = true
				continue
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
			logs.Warn("sns: stale endpoint index entry found during list",
				logs.String("endpointArn", epArn),
				logs.String("appArn", platformAppArn))
			continue
		}

		items = append(items, &ep)
		lastEpArn = epArn
		count++
	}
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("scanning endpoint index: %w", err)
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
