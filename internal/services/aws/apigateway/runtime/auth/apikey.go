// Package auth provides API Gateway authentication functionality for vorpalstacks.
package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

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

// Authenticate validates an API key for the given method and stage.
func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, apiKeyValue string, method *apigatewaystore.Method, restAPIID, stageName string) error {
	if method == nil || !method.ApiKeyRequired {
		return nil
	}

	if apiKeyValue == "" {
		return &AuthError{
			Message:  "Missing API Key",
			Type:     "UnauthorizedException",
			HTTPCode: http.StatusUnauthorized,
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

	stageKey := fmt.Sprintf("%s/%s", restAPIID, stageName)
	if !containsStage(apiKey.StageKeys, stageKey) {
		return &AuthError{
			Message:  "API Key is not authorized for this stage",
			Type:     "ForbiddenException",
			HTTPCode: http.StatusForbidden,
		}
	}

	if err := a.checkUsageQuota(ctx, apiKey); err != nil {
		return err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("Panic in API key usage recording",
					logs.String("apiKeyId", apiKey.Id),
					logs.Any("panic", r))
			}
		}()
		a.recordUsage(ctx, apiKey, restAPIID, stageName)
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
			if err == nil && totalCount >= plan.Quota.Limit {
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
	var totalCount int64

	switch period {
	case "DAY", "":
		today := now.Format("2006-01-02")
		usage, err := a.usageStore.GetUsage(planId, apiKeyId, today)
		if err != nil {
			return 0, nil
		}
		return usage.RequestCount, nil

	case "WEEK":
		for i := 0; i < 7; i++ {
			date := now.AddDate(0, 0, -i).Format("2006-01-02")
			usage, err := a.usageStore.GetUsage(planId, apiKeyId, date)
			if err == nil {
				totalCount += usage.RequestCount
			}
		}
		return totalCount, nil

	case "MONTH":
		year, month, _ := now.Date()
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			usage, err := a.usageStore.GetUsage(planId, apiKeyId, date)
			if err == nil {
				totalCount += usage.RequestCount
			}
		}
		return totalCount, nil

	default:
		today := now.Format("2006-01-02")
		usage, err := a.usageStore.GetUsage(planId, apiKeyId, today)
		if err != nil {
			return 0, nil
		}
		return usage.RequestCount, nil
	}
}

func (a *APIKeyAuthenticator) getRateLimiter(apiKeyId string, rateLimit float64, burstLimit int64) *rateLimiter {
	if actual, loaded := a.rateLimiters.Load(apiKeyId); loaded {
		if typed, ok := actual.(*rateLimiter); ok {
			return typed
		}
	}
	limiter := newRateLimiter(rateLimit, burstLimit)
	if actual, loaded := a.rateLimiters.LoadOrStore(apiKeyId, limiter); loaded {
		if typed, ok := actual.(*rateLimiter); ok {
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

func containsStage(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
}

func newRateLimiter(rateLimit float64, burstLimit int64) *rateLimiter {
	return &rateLimiter{
		tokens:     float64(burstLimit),
		maxTokens:  float64(burstLimit),
		refillRate: rateLimit,
		lastRefill: time.Now(),
	}
}

// Allow attempts to consume a token from the rate limiter.
// Returns true if a token was consumed, false if the rate limit has been exceeded.
func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens = min(r.maxTokens, r.tokens+elapsed*r.refillRate)
	r.lastRefill = now

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
