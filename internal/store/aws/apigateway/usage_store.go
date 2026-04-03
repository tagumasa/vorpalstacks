package apigateway

import (
	"encoding/json"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// UsageStore provides storage operations for API keys and usage plans.
type UsageStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	accountId  string
	region     string
	usageMutex sync.Map
	mu         sync.Mutex
}

func usageBucketName(region string) string {
	return "apigateway-usage-" + region
}

// NewUsageStore creates a new UsageStore instance.
func NewUsageStore(store storage.BasicStorage, accountId, region string) *UsageStore {
	bucket := store.Bucket(usageBucketName(region))
	return &UsageStore{
		BaseStore:  common.NewBaseStore(bucket, "apigateway-usage"),
		arnBuilder: NewARNBuilder(accountId, region),
		accountId:  accountId,
		region:     region,
	}
}

func (s *UsageStore) getUsageLock(key string) *sync.Mutex {
	v, _ := s.usageMutex.LoadOrStore(key, &sync.Mutex{})
	if typed, ok := v.(*sync.Mutex); ok {
		return typed
	}
	return &sync.Mutex{}
}

// CreateApiKey creates a new API key.
func (s *UsageStore) CreateApiKey(apiKey *ApiKey) (*ApiKey, error) {
	if apiKey.Id == "" {
		apiKey.Id = s.arnBuilder.GenerateApiKeyId()
	}
	if apiKey.Name == "" {
		return nil, ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists("apikey#" + apiKey.Id) {
		return nil, ErrApiKeyAlreadyExists
	}

	now := time.Now().UTC()
	apiKey.CreatedDate = now
	apiKey.LastUpdatedDate = now
	if apiKey.Value == "" {
		apiKey.Value = generateApiKeyValue()
	}

	if err := s.Put("apikey#"+apiKey.Id, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// GetApiKey retrieves an API key by its ID.
func (s *UsageStore) GetApiKey(apiKeyId string) (*ApiKey, error) {
	var apiKey ApiKey
	if err := s.BaseStore.Get("apikey#"+apiKeyId, &apiKey); err != nil {
		return nil, ErrApiKeyNotFound
	}
	return &apiKey, nil
}

// UpdateApiKey updates an existing API key.
func (s *UsageStore) UpdateApiKey(apiKey *ApiKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("apikey#" + apiKey.Id) {
		return ErrApiKeyNotFound
	}
	apiKey.LastUpdatedDate = time.Now().UTC()
	return s.Put("apikey#"+apiKey.Id, apiKey)
}

// DeleteApiKey deletes an API key.
func (s *UsageStore) DeleteApiKey(apiKeyId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("apikey#" + apiKeyId) {
		return ErrApiKeyNotFound
	}
	s.BaseStore.DeleteByPrefix("usage##" + apiKeyId + "#")
	return s.BaseStore.Delete("apikey#" + apiKeyId)
}

// ListApiKeys returns all API keys.
func (s *UsageStore) ListApiKeys(opts common.ListOptions) (*common.ListResult[ApiKey], error) {
	return common.List[ApiKey](s.BaseStore, common.ListOptions{
		Prefix:   "apikey#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}

// CreateUsagePlan creates a new usage plan.
func (s *UsageStore) CreateUsagePlan(usagePlan *UsagePlan) (*UsagePlan, error) {
	if usagePlan.Id == "" {
		usagePlan.Id = s.arnBuilder.GenerateUsagePlanId()
	}
	if usagePlan.Name == "" {
		return nil, ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists("usageplan#" + usagePlan.Id) {
		return nil, ErrUsagePlanAlreadyExists
	}

	if err := s.Put("usageplan#"+usagePlan.Id, usagePlan); err != nil {
		return nil, err
	}

	return usagePlan, nil
}

// GetUsagePlan retrieves a usage plan by its ID.
func (s *UsageStore) GetUsagePlan(usagePlanId string) (*UsagePlan, error) {
	var usagePlan UsagePlan
	if err := s.BaseStore.Get("usageplan#"+usagePlanId, &usagePlan); err != nil {
		return nil, ErrUsagePlanNotFound
	}
	return &usagePlan, nil
}

// UpdateUsagePlan updates an existing usage plan.
func (s *UsageStore) UpdateUsagePlan(usagePlan *UsagePlan) error {
	if !s.Exists("usageplan#" + usagePlan.Id) {
		return ErrUsagePlanNotFound
	}
	return s.Put("usageplan#"+usagePlan.Id, usagePlan)
}

// DeleteUsagePlan deletes a usage plan.
func (s *UsageStore) DeleteUsagePlan(usagePlanId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("usageplan#" + usagePlanId) {
		return ErrUsagePlanNotFound
	}

	keys, err := s.ListUsagePlanKeys(usagePlanId, common.ListOptions{MaxItems: 1000})
	if err == nil {
		for _, key := range keys.Items {
			if delErr := s.DeleteUsagePlanKey(usagePlanId, key.Id); delErr != nil {
				logs.Error("Failed to delete usage plan key", logs.String("usagePlanId", usagePlanId), logs.String("keyId", key.Id), logs.Err(delErr))
			}
		}
	}

	s.BaseStore.DeleteByPrefix("usage#" + usagePlanId + "#")
	return s.BaseStore.Delete("usageplan#" + usagePlanId)
}

// ListUsagePlans returns all usage plans.
func (s *UsageStore) ListUsagePlans(opts common.ListOptions) (*common.ListResult[UsagePlan], error) {
	return common.List[UsagePlan](s.BaseStore, common.ListOptions{
		Prefix:   "usageplan#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}

// CreateUsagePlanKey creates a new usage plan key.
func (s *UsageStore) CreateUsagePlanKey(usagePlanId string, key *UsagePlanKey) (*UsagePlanKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, err
	}

	if _, err := s.GetApiKey(key.Id); err != nil {
		return nil, ErrApiKeyNotFound
	}

	keyKey := "usageplankey#" + usagePlanId + "#" + key.Id
	if s.Exists(keyKey) {
		return nil, ErrUsagePlanKeyAlreadyExists
	}

	key.Type = "API_KEY"

	if err := s.Put(keyKey, key); err != nil {
		return nil, err
	}

	return key, nil
}

// GetUsagePlanKey retrieves a usage plan key by its ID.
func (s *UsageStore) GetUsagePlanKey(usagePlanId, keyId string) (*UsagePlanKey, error) {
	var key UsagePlanKey
	if err := s.BaseStore.Get("usageplankey#"+usagePlanId+"#"+keyId, &key); err != nil {
		return nil, ErrUsagePlanKeyNotFound
	}
	return &key, nil
}

// DeleteUsagePlanKey deletes a usage plan key.
func (s *UsageStore) DeleteUsagePlanKey(usagePlanId, keyId string) error {
	keyKey := "usageplankey#" + usagePlanId + "#" + keyId
	if !s.Exists(keyKey) {
		return ErrUsagePlanKeyNotFound
	}
	s.BaseStore.DeleteByPrefix("usage#" + usagePlanId + "#" + keyId + "#")
	return s.BaseStore.Delete(keyKey)
}

// ListUsagePlanKeys returns all usage plan keys for a usage plan.
func (s *UsageStore) ListUsagePlanKeys(usagePlanId string, opts common.ListOptions) (*common.ListResult[UsagePlanKey], error) {
	return common.List[UsagePlanKey](s.BaseStore, common.ListOptions{
		Prefix:   "usageplankey#" + usagePlanId + "#",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}, nil)
}

// GetApiKeyByValue retrieves an API key by its value.
func (s *UsageStore) GetApiKeyByValue(value string) (*ApiKey, error) {
	keys, err := s.ListApiKeys(common.ListOptions{MaxItems: 1000})
	if err != nil {
		return nil, err
	}
	for _, key := range keys.Items {
		if key.Value == value {
			return key, nil
		}
	}
	return nil, ErrApiKeyNotFound
}

// UsageRecord represents a usage record for API Gateway.
type UsageRecord struct {
	UsagePlanID  string `json:"usagePlanId"`
	APIKeyID     string `json:"apiKeyId"`
	Date         string `json:"date"`
	RequestCount int64  `json:"requestCount"`
}

// ListUsagePlansForAPIKey returns all usage plans for an API key.
func (s *UsageStore) ListUsagePlansForAPIKey(apiKeyId string) ([]*UsagePlan, error) {
	var plans []*UsagePlan

	allPlans, err := s.ListUsagePlans(common.ListOptions{MaxItems: 1000})
	if err != nil {
		return nil, err
	}

	for _, plan := range allPlans.Items {
		keys, err := s.ListUsagePlanKeys(plan.Id, common.ListOptions{MaxItems: 1000})
		if err != nil {
			logs.Error("Failed to list usage plan keys", logs.String("planId", plan.Id), logs.Err(err))
			continue
		}
		for _, key := range keys.Items {
			if key.Id == apiKeyId {
				plans = append(plans, plan)
				break
			}
		}
	}

	return plans, nil
}

// ListUsageRecordsForAPIKey returns usage records for an API key.
func (s *UsageStore) ListUsageRecordsForAPIKey(usagePlanId, apiKeyId, startDate, endDate string) ([]*UsageRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := "usage#" + usagePlanId + "#" + apiKeyId + "#"
	var records []*UsageRecord

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var record UsageRecord
		if err := json.Unmarshal(value, &record); err != nil {
			return err
		}
		if startDate != "" && record.Date < startDate {
			return nil
		}
		if endDate != "" && record.Date > endDate {
			return nil
		}
		records = append(records, &record)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return records, nil
}

// GetUsage retrieves usage for an API key on a specific date.
func (s *UsageStore) GetUsage(usagePlanId, apiKeyId, date string) (*UsageRecord, error) {
	key := "usage#" + usagePlanId + "#" + apiKeyId + "#" + date
	var record UsageRecord
	if err := s.BaseStore.Get(key, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// RecordUsage records usage for an API key.
func (s *UsageStore) RecordUsage(record *UsageRecord) error {
	key := "usage#" + record.UsagePlanID + "#" + record.APIKeyID + "#" + record.Date

	mutex := s.getUsageLock(key)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		s.usageMutex.Delete(key)
	}()

	existing, err := s.GetUsage(record.UsagePlanID, record.APIKeyID, record.Date)
	if err != nil {
		record.RequestCount = 1
	} else {
		record.RequestCount = existing.RequestCount + 1
	}

	return s.Put(key, record)
}
