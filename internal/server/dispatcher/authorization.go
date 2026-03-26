package dispatcher

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/iam/policy"
	"vorpalstacks/internal/store/aws/iam"
)

// AuthorizationResult contains the result of an authorization decision.
// It indicates whether the request is allowed, denied, or has no opinion.
// The Reason field provides additional context for the decision.
type AuthorizationResult struct {
	Allowed  bool
	Denied   bool
	Reason   string
	Decision *policy.Decision
}

// Authorizer handles IAM-based authorization for AWS service requests.
// It evaluates IAM policies to determine whether a request should be allowed or denied.
// The authoriser uses a policy cache to improve performance for repeated requests.
type Authorizer struct {
	iamStore          *iam.IAMStore
	policyEvaluator   *policy.PolicyEvaluator
	resourceExtractor *ResourceExtractor
	actionMapper      *ActionMapper

	rootAccessKeys     map[string]bool
	defaultAccessKeyID string
	failureMode        string
	cacheTTL           time.Duration
	policyCache        sync.Map
	maxCacheSize       int

	stopCleanup chan struct{}
	wg          sync.WaitGroup
}

// cachedPolicies holds IAM policies in the cache with their timestamp.
type cachedPolicies struct {
	policies []*policy.Document
	cachedAt time.Time
}

// NewAuthorizer creates a new Authorizer instance with the given IAM store.
// It reads configuration from environment variables:
// - AUTHORIZATION_ROOT_ACCESS_KEYS: Comma-separated list of root access key IDs
// - AUTHORIZATION_DEFAULT_ACCESS_KEY_ID: Default access key ID to use when no signature is provided
// - AUTHORIZATION_FAILURE_MODE: "permissive" (default) or "strict" - how to handle policy fetch errors
// - AUTHORIZATION_CACHE_TTL_SECONDS: Cache TTL in seconds (default 300)
// - AUTHORIZATION_CACHE_MAX_SIZE: Maximum number of cached entries (default 1000)
func NewAuthorizer(iamStore *iam.IAMStore) *Authorizer {
	rootKeys := make(map[string]bool)
	if keys := os.Getenv("AUTHORIZATION_ROOT_ACCESS_KEYS"); keys != "" {
		for _, key := range strings.Split(keys, ",") {
			rootKeys[strings.TrimSpace(key)] = true
		}
	}

	defaultAccessKeyID := os.Getenv("AUTHORIZATION_DEFAULT_ACCESS_KEY_ID")

	failureMode := os.Getenv("AUTHORIZATION_FAILURE_MODE")
	if failureMode == "" {
		failureMode = "strict"
	}

	cacheTTL := 300 * time.Second
	if ttl := os.Getenv("AUTHORIZATION_CACHE_TTL_SECONDS"); ttl != "" {
		if seconds, err := strconv.Atoi(ttl); err == nil && seconds > 0 {
			cacheTTL = time.Duration(seconds) * time.Second
		}
	}

	maxCacheSize := 1000
	if size := os.Getenv("AUTHORIZATION_CACHE_MAX_SIZE"); size != "" {
		if s, err := strconv.Atoi(size); err == nil && s > 0 {
			maxCacheSize = s
		}
	}

	a := &Authorizer{
		iamStore:           iamStore,
		policyEvaluator:    policy.NewPolicyEvaluator(),
		resourceExtractor:  NewResourceExtractor(),
		actionMapper:       NewActionMapper(),
		rootAccessKeys:     rootKeys,
		defaultAccessKeyID: defaultAccessKeyID,
		failureMode:        failureMode,
		cacheTTL:           cacheTTL,
		maxCacheSize:       maxCacheSize,
		stopCleanup:        make(chan struct{}),
	}

	a.startCleanupRoutine()

	return a
}

// Authorize evaluates IAM policies to determine whether a request should be allowed.
// It checks the access key ID against IAM users and their attached policies.
// The method returns an AuthorizationResult indicating whether the request is allowed or denied.
func (a *Authorizer) Authorize(
	ctx context.Context,
	reqCtx *request.RequestContext,
	parsedReq *request.ParsedRequest,
	serviceName string,
	r *http.Request,
) (*AuthorizationResult, error) {
	accessKeyID := parsedReq.AccessKeyID

	if accessKeyID == "" && os.Getenv("TEST_MODE") == "true" && r != nil {
		if testKeyID := r.Header.Get("X-Test-Access-Key-ID"); testKeyID != "" {
			accessKeyID = testKeyID
		}
	}

	// If no access key from signature verification, use default if configured
	if accessKeyID == "" && a.defaultAccessKeyID != "" {
		accessKeyID = a.defaultAccessKeyID
	}

	if a.isRootAccessKey(accessKeyID) {
		return &AuthorizationResult{
			Allowed: true,
			Reason:  "Root access key",
		}, nil
	}

	accessKey, err := a.iamStore.AccessKeys().Get(accessKeyID)
	if err != nil {
		return &AuthorizationResult{
			Allowed: false,
			Denied:  true,
			Reason:  "Invalid access key",
		}, nil
	}

	if accessKey.Status != iam.AccessKeyStatusActive {
		return &AuthorizationResult{
			Allowed: false,
			Denied:  true,
			Reason:  "Access key is inactive",
		}, nil
	}

	user, err := a.iamStore.Users().Get(accessKey.UserName)
	if err != nil {
		return &AuthorizationResult{
			Allowed: false,
			Denied:  true,
			Reason:  "User not found",
		}, nil
	}

	policies, err := a.getEffectivePolicies(ctx, user.UserName)
	if err != nil {
		if a.failureMode == "strict" {
			return nil, fmt.Errorf("failed to get effective policies: %w", err)
		}
		return &AuthorizationResult{
			Allowed: true,
			Reason:  "Policy fetch error (permissive mode)",
		}, nil
	}

	if len(policies) == 0 {
		return &AuthorizationResult{
			Allowed: false,
			Denied:  false,
			Reason:  "No policies attached to user",
		}, nil
	}

	evalCtx := a.buildEvaluationContext(parsedReq, serviceName, user, r)

	log.Printf("[AUTHZ] Evaluating %d policies for user=%s action=%s resource=%s",
		len(policies), user.UserName, evalCtx.Action, evalCtx.Resource)

	decision := a.policyEvaluator.Evaluate(evalCtx, policies)

	log.Printf("[AUTHZ] Decision: effect=%s reason=%s", decision.Effect, decision.Reason)

	switch decision.Effect {
	case policy.DecisionEffectAllow:
		return &AuthorizationResult{
			Allowed:  true,
			Reason:   decision.Reason,
			Decision: decision,
		}, nil
	case policy.DecisionEffectDeny:
		return &AuthorizationResult{
			Allowed:  false,
			Denied:   true,
			Reason:   decision.Reason,
			Decision: decision,
		}, nil
	default:
		return &AuthorizationResult{
			Allowed:  false,
			Denied:   false,
			Reason:   decision.Reason,
			Decision: decision,
		}, nil
	}
}

func (a *Authorizer) isRootAccessKey(accessKeyID string) bool {
	return a.rootAccessKeys[accessKeyID]
}

func (a *Authorizer) getEffectivePolicies(ctx context.Context, userName string) ([]*policy.Document, error) {
	cacheKey := userName

	if cached, ok := a.policyCache.Load(cacheKey); ok {
		cp := cached.(*cachedPolicies)
		if time.Since(cp.cachedAt) < a.cacheTTL {
			return cp.policies, nil
		}
	}

	policies, err := a.fetchEffectivePolicies(userName)
	if err != nil {
		return nil, err
	}

	a.policyCache.Store(cacheKey, &cachedPolicies{
		policies: policies,
		cachedAt: time.Now(),
	})

	return policies, nil
}

func (a *Authorizer) fetchEffectivePolicies(userName string) ([]*policy.Document, error) {
	var documents []*policy.Document

	log.Printf("[AUTHZ] FetchEffectivePolicies for user=%s", userName)

	inlineNames, err := a.iamStore.InlinePolicies().List("user", userName)
	if err != nil {
		log.Printf("[AUTHZ] Warning: failed to list inline policies for user %s: %v", userName, err)
	}
	for _, name := range inlineNames {
		inline, err := a.iamStore.InlinePolicies().Get("user", userName, name)
		if err != nil {
			continue
		}
		doc, err := policy.ParseDocument(inline.PolicyDocument)
		if err != nil {
			continue
		}
		documents = append(documents, doc)
	}

	attachedARNs, err := a.iamStore.AttachedPolicies().ListAttachedPolicies("user", userName)
	if err != nil {
		return nil, fmt.Errorf("list user attached policies: %w", err)
	}
	for _, arn := range attachedARNs {
		version, err := a.iamStore.Policies().GetDefaultVersion(arn)
		if err != nil {
			continue
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil {
			continue
		}
		documents = append(documents, doc)
	}

	groups, err := a.iamStore.UserGroups().ListGroupsForUser(userName)
	if err != nil {
		return nil, fmt.Errorf("list groups for user: %w", err)
	}
	for _, group := range groups {
		groupInlineNames, _ := a.iamStore.InlinePolicies().List("group", group)
		for _, name := range groupInlineNames {
			inline, _ := a.iamStore.InlinePolicies().Get("group", group, name)
			if inline == nil {
				continue
			}
			doc, _ := policy.ParseDocument(inline.PolicyDocument)
			if doc != nil {
				documents = append(documents, doc)
			}
		}

		groupARNs, _ := a.iamStore.AttachedPolicies().ListAttachedPolicies("group", group)
		for _, arn := range groupARNs {
			version, _ := a.iamStore.Policies().GetDefaultVersion(arn)
			if version == nil {
				continue
			}
			doc, _ := policy.ParseDocument(version.Document)
			if doc != nil {
				documents = append(documents, doc)
			}
		}
	}

	log.Printf("[AUTHZ] Found %d policies for user %s", len(documents), userName)
	return documents, nil
}

func (a *Authorizer) buildEvaluationContext(
	parsedReq *request.ParsedRequest,
	serviceName string,
	user *iam.User,
	r *http.Request,
) *policy.EvaluationContext {
	action := a.actionMapper.Map(serviceName, parsedReq.Operation)
	resource := a.resourceExtractor.Extract(
		serviceName,
		parsedReq.Operation,
		parsedReq.Parameters,
		user.AccountId,
		parsedReq.GetRegion(),
	)

	return &policy.EvaluationContext{
		Principal:        user.Arn,
		PrincipalAccount: user.AccountId,
		Action:           action,
		Resource:         resource,
		RequestTime:      time.Now(),
		SourceIP:         extractSourceIP(r),
		UserAgent:        r.UserAgent(),
		UserID:           user.ID,
		UserName:         user.UserName,
		ServiceContext: map[string]string{
			"region":  parsedReq.GetRegion(),
			"service": serviceName,
		},
	}
}

// InvalidateCache removes the cached policies for a specific user.
// This should be called when a user's IAM policies are modified.
func (a *Authorizer) InvalidateCache(userName string) {
	a.policyCache.Delete(userName)
}

// InvalidateAllCache clears all cached policies from the authoriser's cache.
// This is useful when a large number of policy changes occur.
func (a *Authorizer) InvalidateAllCache() {
	a.policyCache.Range(func(key, value interface{}) bool {
		a.policyCache.Delete(key)
		return true
	})
}

// Stop stops the the cleanup goroutine.
func (a *Authorizer) Stop() {
	close(a.stopCleanup)
	a.wg.Wait()
}

func (a *Authorizer) startCleanupRoutine() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		ticker := time.NewTicker(a.cacheTTL)
		defer ticker.Stop()

		for {
			select {
			case <-a.stopCleanup:
				return
			case <-ticker.C:
				a.cleanupExpiredEntries()
			}
		}
	}()
}

func (a *Authorizer) cleanupExpiredEntries() {
	now := time.Now()
	var toDelete []string
	var count int

	a.policyCache.Range(func(key, value interface{}) bool {
		count++
		if cp, ok := value.(*cachedPolicies); ok {
			if now.Sub(cp.cachedAt) > a.cacheTTL {
				toDelete = append(toDelete, key.(string))
			}
		}
		return true
	})

	for _, key := range toDelete {
		a.policyCache.Delete(key)
	}

	if len(toDelete) > 0 {
		log.Printf("[AUTHZ] Cleaned up %d expired cache entries (total: %d)", len(toDelete), count-len(toDelete))
	}

	if count > a.maxCacheSize {
		a.evictOldestEntries(count - a.maxCacheSize)
	}
}

func (a *Authorizer) evictOldestEntries(toEvict int) {
	type entry struct {
		key      string
		cachedAt time.Time
	}

	var entries []entry
	a.policyCache.Range(func(key, value interface{}) bool {
		if cp, ok := value.(*cachedPolicies); ok {
			entries = append(entries, entry{key: key.(string), cachedAt: cp.cachedAt})
		}
		return true
	})

	for i := 0; i < len(entries) && i < toEvict; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].cachedAt.Before(entries[i].cachedAt) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	evicted := 0
	for i := 0; i < toEvict && i < len(entries); i++ {
		a.policyCache.Delete(entries[i].key)
		evicted++
	}

	if evicted > 0 {
		log.Printf("[AUTHZ] Evicted %d oldest cache entries to maintain max size", evicted)
	}
}

func extractSourceIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
