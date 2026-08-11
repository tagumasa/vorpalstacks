package apigateway

import (
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

// ApiKeyInput is the transport-agnostic input for creating an API key.
type ApiKeyInput struct {
	Name        string
	Description string
	Enabled     bool
	CustomerId  string
	Value       string
	Id          string
	StageKeys   []string
	Tags        []types.Tag
}

// UsagePlanInput is the transport-agnostic input for creating a usage plan.
type UsagePlanInput struct {
	Name        string
	Description string
	ApiStages   []ApiStageInput
	Quota       *QuotaInput
	Throttle    *ThrottleInput
	Tags        []types.Tag
}

// ApiStageInput describes a single api-stage attachment on a usage plan.
type ApiStageInput struct {
	ApiId    string
	Stage    string
	Throttle map[string]*apigateway.Throttle
}

// QuotaInput is the transport-agnostic quota settings for a usage plan.
type QuotaInput struct {
	Limit  int64
	Offset int64
	Period string
}

// ThrottleInput is the transport-agnostic throttle settings for a usage
// plan.
type ThrottleInput struct {
	BurstLimit int64
	RateLimit  float64
}

// UsagePlanKeyInput is the transport-agnostic input for attaching an API
// key to a usage plan.
type UsagePlanKeyInput struct {
	KeyId   string
	KeyType string
}

// createApiKeyCore persists an API key. The handler is responsible for
// transport-specific name validation (e.g. the admin handler rejects an
// empty name up-front, while the data-plane path accepts it as before);
// the core here preserves the store-level behaviour.
func (s *APIGatewayService) createApiKeyCore(stores *apiGatewayStores, in *ApiKeyInput) (*apigateway.ApiKey, error) {
	apiKey := &apigateway.ApiKey{
		Name:        in.Name,
		Description: in.Description,
		Enabled:     in.Enabled,
		CustomerId:  in.CustomerId,
		Value:       in.Value,
		Id:          in.Id,
		StageKeys:   in.StageKeys,
		Tags:        in.Tags,
	}
	if !apiKey.Enabled {
		// Default to enabled unless caller explicitly disabled it.
		apiKey.Enabled = true
	}

	created, err := stores.usage.CreateApiKey(apiKey)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getApiKeyCore retrieves a single API key.
func (s *APIGatewayService) getApiKeyCore(stores *apiGatewayStores, apiKeyId string) (*apigateway.ApiKey, error) {
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}
	return stores.usage.GetApiKey(apiKeyId)
}

// deleteApiKeyCore removes an API key and asks the runtime server to drop
// any cached state keyed by the api key id.
func (s *APIGatewayService) deleteApiKeyCore(stores *apiGatewayStores, apiKeyId string) error {
	if apiKeyId == "" {
		return NewBadRequestException("apiKey is required")
	}
	if err := stores.usage.DeleteApiKey(apiKeyId); err != nil {
		return toApiGatewayError(err)
	}
	if s.runtimeServer != nil {
		s.runtimeServer.RemoveApiKey(apiKeyId)
	}
	return nil
}

// listApiKeysCore returns a page of API keys.
func (s *APIGatewayService) listApiKeysCore(stores *apiGatewayStores, limit int, marker string) (*common.ListResult[apigateway.ApiKey], error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListApiKeys(common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// updateApiKeyCore applies patch operations to an API key under the
// per-key key locker and asks the runtime server to refresh its cache.
func (s *APIGatewayService) updateApiKeyCore(
	stores *apiGatewayStores,
	apiKeyId string,
	patches []PatchOperation,
) (*apigateway.ApiKey, error) {
	if apiKeyId == "" {
		return nil, NewBadRequestException("apiKey is required")
	}

	stores.keyLocker.Lock(apiKeyId)
	defer stores.keyLocker.Unlock(apiKeyId)

	apiKey, err := stores.usage.GetApiKey(apiKeyId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		switch {
		case po.Path == "/name":
			apiKey.Name = po.Value
		case po.Path == "/description":
			apiKey.Description = po.Value
		case po.Path == "/customerId":
			apiKey.CustomerId = po.Value
		case po.Path == "/enabled":
			apiKey.Enabled = po.Value == "true"
		case strings.HasPrefix(po.Path, "/stageKeys/"):
			if apiKey.StageKeys == nil {
				apiKey.StageKeys = []string{}
			}
			stageKey := strings.TrimPrefix(po.Path, "/stageKeys/")
			if (po.Op == "add" || po.Op == "replace") && !validateStageKey(stageKey) {
				return nil, NewBadRequestException("invalid stageKey format, expected restApiId/stageName: " + stageKey)
			}
			if po.Op == "add" || po.Op == "replace" {
				if !sliceContains(apiKey.StageKeys, stageKey) {
					apiKey.StageKeys = append(apiKey.StageKeys, stageKey)
				}
			} else if po.Op == "remove" {
				apiKey.StageKeys = removeString(apiKey.StageKeys, stageKey)
			}
		}
	}
	if err := stores.usage.UpdateApiKey(apiKey); err != nil {
		return nil, toApiGatewayError(err)
	}

	if s.runtimeServer != nil {
		s.runtimeServer.RemoveApiKey(apiKeyId)
	}
	return apiKey, nil
}

// createUsagePlanCore persists a usage plan.
func (s *APIGatewayService) createUsagePlanCore(stores *apiGatewayStores, in *UsagePlanInput) (*apigateway.UsagePlan, error) {
	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if !validateUsagePlanNameLen(in.Name) {
		return nil, NewBadRequestException("name must be between 1 and 255 characters")
	}

	plan := &apigateway.UsagePlan{
		Name:        in.Name,
		Description: in.Description,
		Tags:        in.Tags,
	}

	for _, as := range in.ApiStages {
		stage := apigateway.ApiStage{
			ApiId:    as.ApiId,
			Stage:    as.Stage,
			Throttle: as.Throttle,
		}
		plan.ApiStages = append(plan.ApiStages, stage)
	}

	if in.Quota != nil {
		if !validateQuotaPeriod(in.Quota.Period) {
			return nil, NewBadRequestException("Invalid quota period: must be DAY, WEEK, or MONTH")
		}
		plan.Quota = &apigateway.Quota{
			Limit:  in.Quota.Limit,
			Offset: in.Quota.Offset,
			Period: in.Quota.Period,
		}
	}

	if in.Throttle != nil {
		if !validateThrottleBurstLimit(in.Throttle.BurstLimit) {
			return nil, NewBadRequestException("throttle burstLimit must be between 0 and 10000")
		}
		if !validateThrottleRateLimit(in.Throttle.RateLimit) {
			return nil, NewBadRequestException("throttle rateLimit must be between 0 and 10000")
		}
		plan.Throttle = &apigateway.Throttle{
			BurstLimit: in.Throttle.BurstLimit,
			RateLimit:  in.Throttle.RateLimit,
		}
	}

	created, err := stores.usage.CreateUsagePlan(plan)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getUsagePlanCore retrieves a single usage plan.
func (s *APIGatewayService) getUsagePlanCore(stores *apiGatewayStores, usagePlanId string) (*apigateway.UsagePlan, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	plan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return plan, nil
}

// deleteUsagePlanCore removes a usage plan.
func (s *APIGatewayService) deleteUsagePlanCore(stores *apiGatewayStores, usagePlanId string) error {
	if usagePlanId == "" {
		return NewBadRequestException("usagePlanId is required")
	}
	if err := stores.usage.DeleteUsagePlan(usagePlanId); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// listUsagePlansCore returns a page of usage plans.
func (s *APIGatewayService) listUsagePlansCore(stores *apiGatewayStores, limit int, marker string) (*common.ListResult[apigateway.UsagePlan], error) {
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListUsagePlans(common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// updateUsagePlanCore applies patch operations to a usage plan under the
// per-plan key locker.
func (s *APIGatewayService) updateUsagePlanCore(
	stores *apiGatewayStores,
	usagePlanId string,
	patches []PatchOperation,
) (*apigateway.UsagePlan, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}

	stores.keyLocker.Lock(usagePlanId)
	defer stores.keyLocker.Unlock(usagePlanId)

	usagePlan, err := stores.usage.GetUsagePlan(usagePlanId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range patches {
		switch {
		case po.Path == "/name":
			usagePlan.Name = po.Value
		case po.Path == "/description":
			usagePlan.Description = po.Value
		case po.Path == "/productCode":
			usagePlan.ProductCode = po.Value
		case po.Path == "/quota/limit":
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid quota limit: not a number")
			} else {
				usagePlan.Quota.Limit = v
			}
		case po.Path == "/quota/period":
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if !validateQuotaPeriod(po.Value) {
				return nil, NewBadRequestException("Invalid quota period: must be DAY, WEEK, or MONTH")
			}
			usagePlan.Quota.Period = po.Value
		case po.Path == "/quota/offset":
			if usagePlan.Quota == nil {
				usagePlan.Quota = &apigateway.Quota{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid quota offset: not a number")
			} else {
				usagePlan.Quota.Offset = v
			}
		case po.Path == "/throttle/burstLimit":
			if usagePlan.Throttle == nil {
				usagePlan.Throttle = &apigateway.Throttle{}
			}
			if v, err := parseInt64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid throttle burstLimit: not a number")
			} else if !validateThrottleBurstLimit(v) {
				return nil, NewBadRequestException("throttle burstLimit must be between 0 and 10000")
			} else {
				usagePlan.Throttle.BurstLimit = v
			}
		case po.Path == "/throttle/rateLimit":
			if usagePlan.Throttle == nil {
				usagePlan.Throttle = &apigateway.Throttle{}
			}
			if v, err := parseFloat64(po.Value); err != nil {
				return nil, NewBadRequestException("invalid throttle rateLimit: not a number")
			} else if !validateThrottleRateLimit(v) {
				return nil, NewBadRequestException("throttle rateLimit must be between 0 and 10000")
			} else {
				usagePlan.Throttle.RateLimit = v
			}
		case strings.HasPrefix(po.Path, "/apiStages/"):
			rest := strings.TrimPrefix(po.Path, "/apiStages/")
			parts := strings.SplitN(rest, "/", 2)
			idxStr := parts[0]
			var idx int
			if idxStr == "-" || idxStr == "" {
				if po.Op != "remove" {
					usagePlan.ApiStages = append(usagePlan.ApiStages, apigateway.ApiStage{})
					idx = len(usagePlan.ApiStages) - 1
				} else {
					continue
				}
			} else {
				var err error
				idx, err = strconv.Atoi(idxStr)
				if err != nil || idx < 0 {
					continue
				}
				if idx >= len(usagePlan.ApiStages) {
					if po.Op != "remove" {
						usagePlan.ApiStages = append(usagePlan.ApiStages, apigateway.ApiStage{})
						idx = len(usagePlan.ApiStages) - 1
					} else {
						continue
					}
				}
			}
			if len(parts) < 2 {
				continue
			}
			field := parts[1]
			switch {
			case field == "apiId":
				usagePlan.ApiStages[idx].ApiId = po.Value
			case field == "stage":
				usagePlan.ApiStages[idx].Stage = po.Value
			}
		}
	}

	if usagePlan.Quota != nil && !validateQuotaPeriod(usagePlan.Quota.Period) {
		return nil, NewBadRequestException("quota period is required when quota is set")
	}

	if err := stores.usage.UpdateUsagePlan(usagePlan); err != nil {
		return nil, toApiGatewayError(err)
	}
	return usagePlan, nil
}

// createUsagePlanKeyCore attaches an API key to a usage plan. The API key
// must already exist; the key's value and name are copied into the new
// usage plan key entry.
func (s *APIGatewayService) createUsagePlanKeyCore(
	stores *apiGatewayStores,
	usagePlanId string,
	in *UsagePlanKeyInput,
) (*apigateway.UsagePlanKey, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if in.KeyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}
	if in.KeyType == "" {
		return nil, NewBadRequestException("keyType is required")
	}
	if in.KeyType != "API_KEY" {
		return nil, NewBadRequestException("keyType must be API_KEY")
	}

	apiKey, err := stores.usage.GetApiKey(in.KeyId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	key := &apigateway.UsagePlanKey{
		Id:    in.KeyId,
		Type:  in.KeyType,
		Value: apiKey.Value,
		Name:  apiKey.Name,
	}

	created, err := stores.usage.CreateUsagePlanKey(usagePlanId, key)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getUsagePlanKeyCore retrieves a single usage plan key.
func (s *APIGatewayService) getUsagePlanKeyCore(stores *apiGatewayStores, usagePlanId, keyId string) (*apigateway.UsagePlanKey, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if keyId == "" {
		return nil, NewBadRequestException("keyId is required")
	}
	key, err := stores.usage.GetUsagePlanKey(usagePlanId, keyId)
	if err != nil {
		return nil, ErrNotFoundException
	}
	return key, nil
}

// deleteUsagePlanKeyCore removes an API key association from a usage plan.
func (s *APIGatewayService) deleteUsagePlanKeyCore(stores *apiGatewayStores, usagePlanId, keyId string) error {
	if usagePlanId == "" {
		return NewBadRequestException("usagePlanId is required")
	}
	if keyId == "" {
		return NewBadRequestException("keyId is required")
	}
	if err := stores.usage.DeleteUsagePlanKey(usagePlanId, keyId); err != nil {
		return ErrNotFoundException
	}
	return nil
}

// listUsagePlanKeysCore returns a page of usage plan keys.
func (s *APIGatewayService) listUsagePlanKeysCore(stores *apiGatewayStores, usagePlanId string, limit int, marker string) (*common.ListResult[apigateway.UsagePlanKey], error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return stores.usage.ListUsagePlanKeys(usagePlanId, common.ListOptions{
		Marker:   marker,
		MaxItems: resolved,
	})
}

// getUsageCore returns aggregated usage data for a usage plan over a date
// range. The caller pre-validates the date format and the 90-day range
// because that is transport-specific.
func (s *APIGatewayService) getUsageCore(
	stores *apiGatewayStores,
	usagePlanId, keyId, startDate, endDate string,
) (map[string]interface{}, error) {
	if usagePlanId == "" {
		return nil, NewBadRequestException("usagePlanId is required")
	}
	if startDate == "" || endDate == "" {
		return nil, NewBadRequestException("startDate and endDate are required")
	}
	if !validateUsageDateFormat(startDate) || !validateUsageDateFormat(endDate) {
		return nil, NewBadRequestException("startDate and endDate must be in YYYY-MM-DD format")
	}

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)
	if startTime.After(endTime) {
		return nil, NewBadRequestException("startDate must not be after endDate")
	}
	if endTime.Sub(startTime).Hours()/24 > 90 {
		return nil, NewBadRequestException("The date range must not exceed 90 days")
	}

	if _, err := stores.usage.GetUsagePlan(usagePlanId); err != nil {
		return nil, ErrNotFoundException
	}

	var apiKeys []string
	if keyId != "" {
		apiKeys = []string{keyId}
	} else {
		allKeys, err := common.ListMatching[apigateway.UsagePlanKey](stores.usage.BaseStore, "usageplankey#"+usagePlanId+"#", nil)
		if err != nil {
			return nil, toApiGatewayError(err)
		}
		for _, key := range allKeys {
			apiKeys = append(apiKeys, key.Id)
		}
	}

	usageCounts := make(map[string]int64)
	for _, keyId := range apiKeys {
		records, err := stores.usage.ListUsageRecordsForAPIKey(usagePlanId, keyId, startDate, endDate)
		if err != nil {
			logs.Warn("GetUsage: failed to list usage records for key", logs.String("keyId", keyId), logs.Err(err))
			continue
		}
		for _, record := range records {
			usageCounts[record.Date] += record.RequestCount
		}
	}

	items := make([]interface{}, 0, len(usageCounts))
	for date, count := range usageCounts {
		items = append(items, map[string]interface{}{
			"date":         date,
			"requestCount": count,
		})
	}

	return map[string]interface{}{
		"usagePlanId": usagePlanId,
		"startDate":   startDate,
		"endDate":     endDate,
		"items":       items,
	}, nil
}
