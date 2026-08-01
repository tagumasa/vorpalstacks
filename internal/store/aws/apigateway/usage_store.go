package apigateway

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

// UsageStore provides storage operations for API keys and usage plans.
type UsageStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	accountId  string
	region     string
	keyLocker  common.KeyLocker
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

// deleteUsagePlanKeyLocked removes a usage plan key without acquiring the
// store lock. The caller must hold s.mu.
func (s *UsageStore) deleteUsagePlanKeyLocked(usagePlanId, keyId string) error {
	keyKey := "usageplankey#" + usagePlanId + "#" + keyId
	if !s.Exists(keyKey) {
		return ErrUsagePlanKeyNotFound
	}
	if err := s.BaseStore.DeleteByPrefix("usage#" + usagePlanId + "#" + keyId + "#"); err != nil {
		logs.Warn("failed to delete usage records for usage plan key", logs.String("usagePlanId", usagePlanId), logs.String("keyId", keyId), logs.Err(err))
	}
	// Delete reverse index entry (apikeyplan#apiKeyId#planId).
	_ = s.BaseStore.Delete("apikeyplan#" + keyId + "#" + usagePlanId)
	return s.BaseStore.Delete(keyKey)
}

// CreateApiKey creates a new API key.
func (s *UsageStore) CreateApiKey(apiKey *ApiKey) (*ApiKey, error) {
	if apiKey.Id == "" {
		apiKey.Id = s.arnBuilder.GenerateApiKeyId()
	}
	if apiKey.Name == "" {
		apiKey.Name = generateId("auto_", 16)
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

	// Delete usage records across all plans that include this API key.
	// Usage record keys have the format: usage#planId#apiKeyId#date
	plans, err := common.ListMatching[UsagePlan](s.BaseStore, "usageplan#", nil)
	if err != nil {
		logs.Warn("failed to list usage plans for API key cleanup", logs.String("apiKeyId", apiKeyId), logs.Err(err))
	} else {
		for _, plan := range plans {
			if err := s.BaseStore.DeleteByPrefix("usage#" + plan.Id + "#" + apiKeyId + "#"); err != nil {
				logs.Warn("failed to delete usage records", logs.String("apiKeyId", apiKeyId), logs.String("planId", plan.Id), logs.Err(err))
			}
		}
	}

	// Delete reverse index entries (apikeyplan#apiKeyId#planId).
	if err := s.BaseStore.DeleteByPrefix("apikeyplan#" + apiKeyId + "#"); err != nil {
		logs.Warn("failed to delete reverse index for API key", logs.String("apiKeyId", apiKeyId), logs.Err(err))
	}

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
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("usageplan#" + usagePlan.Id) {
		return ErrUsagePlanNotFound
	}
	return s.Put("usageplan#"+usagePlan.Id, usagePlan)
}

// DeleteUsagePlan deletes a usage plan and its associated keys.
func (s *UsageStore) DeleteUsagePlan(usagePlanId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists("usageplan#" + usagePlanId) {
		return ErrUsagePlanNotFound
	}

	if err := common.ForEachAll[UsagePlanKey](s.BaseStore, "usageplankey#"+usagePlanId+"#", nil, func(k *UsagePlanKey) error {
		if delErr := s.deleteUsagePlanKeyLocked(usagePlanId, k.Id); delErr != nil {
			logs.Error("Failed to delete usage plan key", logs.String("usagePlanId", usagePlanId), logs.String("keyId", k.Id), logs.Err(delErr))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to clean up usage plan keys: %w", err)
	}

	if err := s.BaseStore.DeleteByPrefix("usage#" + usagePlanId + "#"); err != nil {
		logs.Warn("failed to delete usage records for usage plan", logs.String("usagePlanId", usagePlanId), logs.Err(err))
	}
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

	// Maintain reverse index: apiKeyId → usagePlanId for O(1) lookup in
	// ListUsagePlansForAPIKey (called on every authenticated request).
	if err := s.Put("apikeyplan#"+key.Id+"#"+usagePlanId, true); err != nil {
		logs.Warn("failed to write reverse index for usage plan key", logs.String("apiKeyId", key.Id), logs.String("usagePlanId", usagePlanId), logs.Err(err))
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
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.deleteUsagePlanKeyLocked(usagePlanId, keyId)
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
	keys, err := common.ListMatching[ApiKey](s.BaseStore, "apikey#", func(k *ApiKey) bool {
		return k.Value == value
	})
	if err != nil {
		return nil, err
	}
	if len(keys) > 0 {
		return keys[0], nil
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

// ListUsagePlansForAPIKey returns all usage plans associated with an API key.
// Uses the reverse index (apikeyplan#apiKeyId#planId) for O(keys) lookup
// instead of the previous O(plans × keys/plan) nested scan.
func (s *UsageStore) ListUsagePlansForAPIKey(apiKeyId string) ([]*UsagePlan, error) {
	prefix := "apikeyplan#" + apiKeyId + "#"

	var plans []*UsagePlan
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		planId := strings.TrimPrefix(key, prefix)
		if planId == "" {
			return nil
		}
		plan, err := s.GetUsagePlan(planId)
		if err != nil {
			logs.Warn("failed to get usage plan from reverse index", logs.String("planId", planId), logs.Err(err))
			return nil
		}
		plans = append(plans, plan)
		return nil
	})
	if err != nil {
		return nil, err
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

	s.keyLocker.Lock(key)
	defer s.keyLocker.Unlock(key)

	existing, err := s.GetUsage(record.UsagePlanID, record.APIKeyID, record.Date)
	if err != nil {
		record.RequestCount = 1
	} else {
		record.RequestCount = existing.RequestCount + 1
	}

	return s.Put(key, record)
}

// TagApiKey adds or updates tags on an API key.
func (s *UsageStore) TagApiKey(apiKeyId string, inputTags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	apiKey, err := s.GetApiKey(apiKeyId)
	if err != nil {
		return err
	}
	if apiKey.Tags == nil {
		apiKey.Tags = []types.Tag{}
	}
	apiKey.Tags = tags.Apply(apiKey.Tags, tags.MapToTags(inputTags))
	return s.Put("apikey#"+apiKeyId, apiKey)
}

// UntagApiKey removes tags from an API key.
func (s *UsageStore) UntagApiKey(apiKeyId string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	apiKey, err := s.GetApiKey(apiKeyId)
	if err != nil {
		return err
	}
	apiKey.Tags = tags.RemoveByTagKeys(apiKey.Tags, tagKeys)
	return s.Put("apikey#"+apiKeyId, apiKey)
}

// GetApiKeyTags returns tags for an API key.
func (s *UsageStore) GetApiKeyTags(apiKeyId string) ([]types.Tag, error) {
	apiKey, err := s.GetApiKey(apiKeyId)
	if err != nil {
		return nil, err
	}
	return apiKey.Tags, nil
}

// TagUsagePlan adds or updates tags on a usage plan.
func (s *UsageStore) TagUsagePlan(usagePlanId string, inputTags map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	usagePlan, err := s.GetUsagePlan(usagePlanId)
	if err != nil {
		return err
	}
	if usagePlan.Tags == nil {
		usagePlan.Tags = []types.Tag{}
	}
	usagePlan.Tags = tags.Apply(usagePlan.Tags, tags.MapToTags(inputTags))
	return s.Put("usageplan#"+usagePlanId, usagePlan)
}

// UntagUsagePlan removes tags from a usage plan.
func (s *UsageStore) UntagUsagePlan(usagePlanId string, tagKeys []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	usagePlan, err := s.GetUsagePlan(usagePlanId)
	if err != nil {
		return err
	}
	usagePlan.Tags = tags.RemoveByTagKeys(usagePlan.Tags, tagKeys)
	return s.Put("usageplan#"+usagePlanId, usagePlan)
}

// GetUsagePlanTags returns tags for a usage plan.
func (s *UsageStore) GetUsagePlanTags(usagePlanId string) ([]types.Tag, error) {
	usagePlan, err := s.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, err
	}
	return usagePlan.Tags, nil
}
