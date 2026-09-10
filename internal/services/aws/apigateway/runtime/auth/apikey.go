// Package auth provides API Gateway authentication functionality for vorpalstacks.
package auth

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/apigateway/runtime/ratelimit"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	storagecommon "vorpalstacks/internal/store/aws/common"
)

// usageRecordTimeout bounds a single usage-recording write so the detached
// goroutine cannot linger indefinitely.
const usageRecordTimeout = 10 * time.Second

// AuthError represents an authentication error with HTTP details.
type AuthError struct {
	Message  string
	Type     string
	HTTPCode int
}

// Error returns the error message for the authentication error.
func (e *AuthError) Error() string {
	return e.Message
}

// APIKeyAuthenticator handles API key authentication for API Gateway.
type APIKeyAuthenticator struct {
	usageStore   *apigatewaystore.UsageStore
	rateLimiters sync.Map
}

// NewAPIKeyAuthenticator creates a new API key authenticator instance.
func NewAPIKeyAuthenticator(usageStore *apigatewaystore.UsageStore) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{
		usageStore: usageStore,
	}
}

// RemoveApiKey cleans up the rate limiter associated with the given API key.
// Call this when an API key is deleted to prevent unbounded growth of the
// rate limiter cache.
func (a *APIKeyAuthenticator) RemoveApiKey(apiKeyId string) {
	a.rateLimiters.Delete(apiKeyId)
}

// Authenticate validates an API key for the given method and stage.
func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, apiKeyValue string, method *apigatewaystore.Method, restAPIID, stageName string) error {
	if method == nil || !method.ApiKeyRequired {
		if apiKeyValue != "" {
			logs.Debug("API key supplied for a method that does not require one; ignoring",
				logs.String("restApiId", restAPIID), logs.String("stageName", stageName))
		}
		return nil
	}

	if apiKeyValue == "" {
		return &AuthError{
			Message:  "Missing API Key",
			Type:     "ForbiddenException",
			HTTPCode: http.StatusForbidden,
		}
	}

	apiKey, err := a.usageStore.GetApiKeyByValue(apiKeyValue)
	if err != nil {
		return &AuthError{
			Message:  "Invalid API Key",
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
		}
	}

	if !apiKey.Enabled {
		return &AuthError{
			Message:  "API Key is disabled",
			Type:     "ForbiddenException",
			HTTPCode: http.StatusForbidden,
		}
	}

	stageKey := apigatewaystore.StageKey(restAPIID, stageName)
	if !slices.Contains(apiKey.StageKeys, stageKey) {
		return &AuthError{
			Message:  "API Key is not authorized for this stage",
			Type:     "ForbiddenException",
			HTTPCode: http.StatusForbidden,
		}
	}

	if err := a.checkUsageQuota(ctx, apiKey); err != nil {
		return err
	}

	// Usage recording must not race the request's completion: it runs on
	// its own bounded context so a cancelled or disconnected request
	// still leaves its usage accounted.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("Panic in API key usage recording",
					logs.String("apiKeyId", apiKey.Id),
					logs.Any("panic", r))
			}
		}()
		recordCtx, cancel := context.WithTimeout(context.Background(), usageRecordTimeout)
		defer cancel()
		a.recordUsage(recordCtx, apiKey, restAPIID, stageName)
	}()

	return nil
}

func (a *APIKeyAuthenticator) checkUsageQuota(ctx context.Context, apiKey *apigatewaystore.ApiKey) error {
	usagePlans, err := a.usageStore.ListUsagePlansForAPIKey(apiKey.Id)
	if err != nil {
		logs.Warn("failed to list usage plans for API key", logs.Err(err), logs.String("apiKeyId", apiKey.Id))
		return nil
	}

	for _, plan := range usagePlans {
		if plan.Quota != nil {
			totalCount, err := a.getQuotaUsageCount(plan.Id, apiKey.Id, plan.Quota.Period)
			if err != nil {
				logs.Warn("failed to get quota usage count, denying request (fail-closed)",
					logs.Err(err), logs.String("apiKeyId", apiKey.Id), logs.String("planId", plan.Id))
				return &AuthError{
					Message:  "API Key quota check failed",
					Type:     "TooManyRequestsException",
					HTTPCode: http.StatusTooManyRequests,
				}
			}
			// A remaining-quota override granted by UpdateUsage raises the
			// effective limit until its quota period lapses.
			effectiveLimit := plan.Quota.Limit + a.usageStore.ActiveQuotaExtension(
				plan.Id, apiKey.Id, time.Now().Format("2006-01-02"))
			if totalCount >= effectiveLimit {
				return &AuthError{
					Message:  "API Key quota exceeded",
					Type:     "TooManyRequestsException",
					HTTPCode: http.StatusTooManyRequests,
				}
			}
		}

		if plan.Throttle != nil {
			limiter := a.getRateLimiter(apiKey.Id, plan.Throttle.RateLimit, plan.Throttle.BurstLimit)
			if !limiter.Allow() {
				return &AuthError{
					Message:  "Rate limit exceeded",
					Type:     "TooManyRequestsException",
					HTTPCode: http.StatusTooManyRequests,
				}
			}
		}
	}

	return nil
}

func (a *APIKeyAuthenticator) getQuotaUsageCount(planId, apiKeyId, period string) (int64, error) {
	now := time.Now()
	var dates []string
	switch period {
	case "DAY":
		dates = []string{now.Format("2006-01-02")}
	case "WEEK":
		for i := 0; i < 7; i++ {
			dates = append(dates, now.AddDate(0, 0, -i).Format("2006-01-02"))
		}
	case "MONTH":
		year, month, _ := now.Date()
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for day := 1; day <= daysInMonth; day++ {
			dates = append(dates, time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02"))
		}
	default:
		return 0, fmt.Errorf("unsupported quota period: %s", period)
	}

	// Days with no usage record simply contribute nothing to the count.
	var totalCount int64
	for _, date := range dates {
		usage, err := a.usageStore.GetUsage(planId, apiKeyId, date)
		if err != nil {
			if !storagecommon.IsNotFound(err) {
				return 0, err
			}
			continue
		}
		totalCount += usage.RequestCount
	}
	return totalCount, nil
}

func (a *APIKeyAuthenticator) getRateLimiter(apiKeyId string, rateLimit float64, burstLimit int64) *ratelimit.TokenBucket {
	if actual, loaded := a.rateLimiters.Load(apiKeyId); loaded {
		if typed, ok := actual.(*ratelimit.TokenBucket); ok {
			return typed
		}
	}
	limiter := ratelimit.New(rateLimit, float64(burstLimit))
	if actual, loaded := a.rateLimiters.LoadOrStore(apiKeyId, limiter); loaded {
		if typed, ok := actual.(*ratelimit.TokenBucket); ok {
			return typed
		}
	}
	return limiter
}

func (a *APIKeyAuthenticator) recordUsage(ctx context.Context, apiKey *apigatewaystore.ApiKey, restAPIID, stage string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	usagePlans, err := a.usageStore.ListUsagePlansForAPIKey(apiKey.Id)
	if err != nil {
		return
	}

	today := time.Now().Format("2006-01-02")
	for _, plan := range usagePlans {
		record := &apigatewaystore.UsageRecord{
			UsagePlanID:  plan.Id,
			APIKeyID:     apiKey.Id,
			Date:         today,
			RequestCount: 1,
		}
		if err := a.usageStore.RecordUsage(record); err != nil {
			logs.Debug("failed to record usage", logs.String("usagePlanId", plan.Id), logs.Err(err))
		}
	}
}
