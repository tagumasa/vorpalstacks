package authorization

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/iam"
)

// Authorizer handles IAM-based authorization for AWS service requests.
// It evaluates IAM policies to determine whether a request should be allowed or denied.
// The authoriser uses a policy cache to improve performance for repeated requests.
//
// Root user access keys (UserName == iam.RootUserName) bypass all policy evaluation
// and are granted unrestricted access, consistent with the AWS root user model.
//
// Temporary STS credentials (access keys prefixed with "ASIA") are resolved
// via the SessionResolver and authorised against the assumed role's identity-
// based policies intersected with any session policy the caller supplied.
// Session tags surface in EvaluationContext.SessionContext as
// aws:PrincipalTag/<key> values; the session's SourceIdentity surfaces as
// sts:SourceIdentity for condition-key evaluation.
type Authorizer struct {
	iamStore          iam.IAMStoreInterface
	sessionResolver   auth.SessionResolver
	policyEvaluator   *policy.PolicyEvaluator
	resourceExtractor *ResourceExtractor
	actionMapper      *ActionMapper

	defaultAccessKeyID string
	failureMode        string
	cacheTTL           time.Duration
	policyCache        sync.Map
	maxCacheSize       int

	stopCleanup chan struct{}
	stopOnce    sync.Once
	wg          sync.WaitGroup
}

// cachedPolicies holds IAM policies in the cache with their timestamp.
type cachedPolicies struct {
	policies []*policy.Document
	cachedAt time.Time
}

// NewAuthorizer creates a new Authorizer instance with the given IAM store
// and STS session resolver. The session resolver is consulted when the
// caller presents an "ASIA" temporary access key; pass nil to disable
// session-based authorisation (legacy permanent-key-only behaviour).
//
// It reads configuration from environment variables:
// - AUTHORIZATION_DEFAULT_ACCESS_KEY_ID: Default access key ID to use when no signature is provided
// - AUTHORIZATION_FAILURE_MODE: "permissive" (default) or "strict" - how to handle policy fetch errors
// - AUTHORIZATION_CACHE_TTL_SECONDS: Cache TTL in seconds (default 300)
// - AUTHORIZATION_CACHE_MAX_SIZE: Maximum number of cached entries (default 1000)
func NewAuthorizer(iamStore iam.IAMStoreInterface, sessionResolver auth.SessionResolver) *Authorizer {
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
		sessionResolver:    sessionResolver,
		policyEvaluator:    policy.NewPolicyEvaluator(),
		resourceExtractor:  NewResourceExtractor(),
		actionMapper:       NewActionMapper(),
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
func (a *Authorizer) Authorize(
	ctx context.Context,
	reqCtx *request.RequestContext,
	parsedReq *request.ParsedRequest,
	serviceName string,
	r *http.Request,
) (bool, error) {
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

	accessKey, err := a.iamStore.AccessKeys().Get(accessKeyID)
	if err != nil {
		// IAM permanent key lookup failed. If the caller presents a
		// temporary "ASIA" access key and we have a session resolver
		// configured, attempt session-based authorisation before
		// fail-closing.
		if strings.HasPrefix(accessKeyID, "ASIA") && a.sessionResolver != nil {
			return a.authorizeSession(ctx, reqCtx, parsedReq, serviceName, r, accessKeyID)
		}
		// Fail-closed: deny access on store errors, but log the underlying
		// cause so transient failures are diagnosable.
		logs.Warn("Access key lookup failed during authorization",
			logs.String("accessKeyId", accessKeyID), logs.Err(err))
		return false, nil
	}

	if accessKey.Status != iam.AccessKeyStatusActive {
		return false, nil
	}

	// Root user access keys bypass all IAM policy evaluation.
	if accessKey.UserName == iam.RootUserName {
		reqCtx.Principal = iam.RootUserName
		reqCtx.PrincipalID = iam.RootUserName
		reqCtx.PrincipalType = request.PrincipalTypeUser
		return true, nil
	}

	user, err := a.iamStore.Users().Get(accessKey.UserName)
	if err != nil {
		return false, nil
	}

	reqCtx.Principal = user.UserName
	reqCtx.PrincipalID = user.ID
	reqCtx.PrincipalType = request.PrincipalTypeUser

	policies, err := a.getEffectivePolicies(ctx, user.UserName)
	if err != nil {
		if a.failureMode == "strict" {
			return false, fmt.Errorf("failed to get effective policies: %w", err)
		}
		return true, nil
	}

	if len(policies) == 0 {
		return false, nil
	}

	evalCtx := a.buildEvaluationContext(parsedReq, serviceName, user, r)

	logs.Info("Evaluating policies",
		logs.Int("count", len(policies)),
		logs.String("user", user.UserName),
		logs.String("action", evalCtx.Action),
		logs.String("resource", evalCtx.Resource),
	)

	decision := a.policyEvaluator.Evaluate(evalCtx, policies)

	logs.Info("Authorization decision", logs.String("effect", string(decision.Effect)), logs.String("reason", decision.Reason))

	switch decision.Effect {
	case policy.DecisionEffectAllow:
		return true, nil
	case policy.DecisionEffectDeny:
		return false, nil
	default:
		return false, nil
	}
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

	a.enforceCacheSize()

	return policies, nil
}

func (a *Authorizer) fetchEffectivePolicies(userName string) ([]*policy.Document, error) {
	var documents []*policy.Document

	logs.Info("Fetching effective policies", logs.String("user", userName))

	inlineNames, err := a.iamStore.InlinePolicies().List("user", userName)
	if err != nil {
		logs.Warn("Failed to list inline policies", logs.String("user", userName), logs.Err(err))
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

	logs.Info("Found policies for user", logs.Int("count", len(documents)), logs.String("user", userName))
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

// authorizeSession resolves an STS temporary credential (ASIA-prefixed access
// key) via the session resolver and evaluates the assumed role's policies
// intersected with any session policy the caller supplied. Session tags
// surface as aws:PrincipalTag/<key> values and SourceIdentity as
// sts:SourceIdentity in the policy evaluation context.
//
// The Root principal type bypasses policy evaluation, matching the
// permanent-key root-user behaviour. Federated users (GetFederationToken)
// resolve to the underlying IAM user via the session's stored principal
// name when the caller was a permanent IAM user.
func (a *Authorizer) authorizeSession(
	ctx context.Context,
	reqCtx *request.RequestContext,
	parsedReq *request.ParsedRequest,
	serviceName string,
	r *http.Request,
	accessKeyID string,
) (bool, error) {
	sessionCreds, err := a.sessionResolver.ResolveSession(accessKeyID)
	if err != nil || sessionCreds == nil {
		logs.Warn("Session resolution failed during authorization",
			logs.String("accessKeyId", accessKeyID), logs.Err(err))
		return false, nil
	}

	// Root sessions bypass policy evaluation.
	if sessionCreds.PrincipalType == "Root" || strings.HasSuffix(sessionCreds.PrincipalArn, ":root") {
		reqCtx.Principal = iam.RootUserName
		reqCtx.PrincipalID = iam.RootUserName
		reqCtx.PrincipalType = request.PrincipalTypeUser
		return true, nil
	}

	// For User / FederatedUser sessions (GetSessionToken,
	// GetFederationToken), evaluate the underlying IAM user's effective
	// policies — same policies as the permanent access key. For
	// AssumedRole / SAML / WebIdentity sessions, evaluate the assumed
	// role's identity-based policies.
	var effectivePolicies []*policy.Document
	if sessionCreds.PrincipalType == "User" || sessionCreds.PrincipalType == "FederatedUser" {
		userName := extractUserNameFromArn(sessionCreds.PrincipalArn)
		if userName == "" {
			logs.Warn("Could not extract user name from session principal ARN",
				logs.String("arn", sessionCreds.PrincipalArn))
			return false, nil
		}
		effectivePolicies, err = a.getEffectivePolicies(ctx, userName)
		if err != nil {
			return false, nil
		}
		reqCtx.Principal = userName
		reqCtx.PrincipalID = userName
		reqCtx.PrincipalType = request.PrincipalTypeUser
	} else {
		effectivePolicies, _ = a.fetchEffectiveRolePolicies(ctx, sessionCreds.PrincipalArn)
		reqCtx.Principal = sessionCreds.PrincipalArn
		reqCtx.PrincipalID = sessionCreds.PrincipalArn
		reqCtx.PrincipalType = request.PrincipalTypeRole
	}

	// Build the session evaluation context.
	evalCtx := a.buildSessionEvaluationContext(reqCtx, parsedReq, serviceName, sessionCreds, r)

	logs.Info("Evaluating session policies",
		logs.String("principal", sessionCreds.PrincipalArn),
		logs.String("principalType", sessionCreds.PrincipalType),
		logs.String("action", evalCtx.Action),
		logs.String("resource", evalCtx.Resource),
	)

	decision := a.policyEvaluator.Evaluate(evalCtx, effectivePolicies)
	if decision.Effect != policy.DecisionEffectAllow {
		return false, nil
	}

	// Session-scoping policies (inline Policy + managed PolicyArns).
	// AWS intersection semantics: the role's identity-based policies
	// AND every session policy must independently Allow. If any session
	// policy denies (or fails to resolve), the request is denied.
	sessionDocs, ok := a.collectSessionPolicyDocuments(sessionCreds)
	if !ok {
		return false, nil
	}
	for _, doc := range sessionDocs {
		sessionDecision := a.policyEvaluator.Evaluate(evalCtx, []*policy.Document{doc})
		if sessionDecision.Effect != policy.DecisionEffectAllow {
			logs.Info("Session policy denied action",
				logs.String("principal", sessionCreds.PrincipalArn),
				logs.String("reason", sessionDecision.Reason))
			return false, nil
		}
	}

	return true, nil
}

// collectSessionPolicyDocuments gathers all session-scoping policy
// documents attached to a temporary credential: the inline Policy string
// and each managed policy referenced by PolicyArns.
//
// Returns (docs, true) on success — docs may be nil/empty when the
// session has no scoping policies, which is the normal case. Returns
// (nil, false) when a document cannot be parsed or an ARN cannot be
// resolved; the caller MUST treat ok=false as deny-closed.
func (a *Authorizer) collectSessionPolicyDocuments(sessionCreds *auth.SessionCredentials) ([]*policy.Document, bool) {
	var docs []*policy.Document

	// Inline session policy (--policy on AssumeRole et al.).
	if sessionCreds.Policy != "" {
		doc, err := policy.ParseDocument(sessionCreds.Policy)
		if err != nil {
			logs.Warn("Failed to parse inline session policy",
				logs.String("principal", sessionCreds.PrincipalArn), logs.Err(err))
			return nil, false
		}
		docs = append(docs, doc)
	}

	// Managed-policy ARNs (--policy-arns on AssumeRole et al.). Each
	// ARN references a managed policy whose default-version document
	// acts as an additional session-scoping policy.
	for _, arn := range sessionCreds.PolicyArns {
		version, err := a.iamStore.Policies().GetDefaultVersion(arn)
		if err != nil || version == nil {
			logs.Warn("Failed to resolve session policy ARN",
				logs.String("arn", arn),
				logs.String("principal", sessionCreds.PrincipalArn), logs.Err(err))
			return nil, false
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil || doc == nil {
			logs.Warn("Failed to parse session policy ARN document",
				logs.String("arn", arn),
				logs.String("principal", sessionCreds.PrincipalArn), logs.Err(err))
			return nil, false
		}
		docs = append(docs, doc)
	}

	return docs, true
}

// fetchEffectiveRolePolicies resolves the role referenced by the session
// principal ARN and returns the role's identity-based policy documents
// (inline + attached). Returns an empty document slice when the role
// cannot be found so the caller fails closed (no implicit allow).
func (a *Authorizer) fetchEffectiveRolePolicies(ctx context.Context, sessionPrincipalArn string) ([]*policy.Document, string) {
	roleName := extractRoleNameFromAssumedRoleArn(sessionPrincipalArn)
	if roleName == "" {
		// Not an assumed-role ARN; fall through with an empty policy set.
		return nil, ""
	}
	cacheKey := "role:" + roleName
	if cached, ok := a.policyCache.Load(cacheKey); ok {
		cp := cached.(*cachedPolicies)
		if time.Since(cp.cachedAt) < a.cacheTTL {
			return cp.policies, roleName
		}
	}

	role, err := a.iamStore.Roles().Get(roleName)
	if err != nil || role == nil {
		return nil, roleName
	}

	var documents []*policy.Document

	inlineNames, _ := a.iamStore.InlinePolicies().List("role", roleName)
	for _, name := range inlineNames {
		inline, err := a.iamStore.InlinePolicies().Get("role", roleName, name)
		if err != nil || inline == nil {
			continue
		}
		doc, err := policy.ParseDocument(inline.PolicyDocument)
		if err != nil || doc == nil {
			continue
		}
		documents = append(documents, doc)
	}

	attachedARNs, _ := a.iamStore.AttachedPolicies().ListAttachedPolicies("role", roleName)
	for _, arn := range attachedARNs {
		version, err := a.iamStore.Policies().GetDefaultVersion(arn)
		if err != nil || version == nil {
			continue
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil || doc == nil {
			continue
		}
		documents = append(documents, doc)
	}

	a.policyCache.Store(cacheKey, &cachedPolicies{
		policies: documents,
		cachedAt: time.Now(),
	})
	a.enforceCacheSize()

	logs.Info("Found policies for role", logs.Int("count", len(documents)), logs.String("role", roleName))
	return documents, roleName
}

// buildSessionEvaluationContext constructs an EvaluationContext for an STS
// temporary credential. Session tags populate the SessionContext map so
// aws:PrincipalTag/<key> condition keys resolve; SourceIdentity surfaces
// as sts:SourceIdentity.
func (a *Authorizer) buildSessionEvaluationContext(
	reqCtx *request.RequestContext,
	parsedReq *request.ParsedRequest,
	serviceName string,
	sessionCreds *auth.SessionCredentials,
	r *http.Request,
) *policy.EvaluationContext {
	accountID := reqCtx.GetAccountID()
	if accountID == "" {
		// Fallback: extract the account from the assumed-role ARN
		// (arn:aws:iam::<account>:role/<name>).
		accountID = extractAccountIDFromArn(sessionCreds.PrincipalArn)
	}
	action := a.actionMapper.Map(serviceName, parsedReq.Operation)
	resource := a.resourceExtractor.Extract(
		serviceName,
		parsedReq.Operation,
		parsedReq.Parameters,
		accountID,
		parsedReq.GetRegion(),
	)

	sessionContext := make(map[string]string, len(sessionCreds.Tags)+1)
	for k, v := range sessionCreds.Tags {
		sessionContext["aws:PrincipalTag/"+k] = v
	}
	if sessionCreds.SourceIdentity != "" {
		sessionContext["sts:SourceIdentity"] = sessionCreds.SourceIdentity
	}

	return &policy.EvaluationContext{
		Principal:        sessionCreds.PrincipalArn,
		PrincipalAccount: accountID,
		Action:           action,
		Resource:         resource,
		RequestTime:      time.Now(),
		SourceIP:         extractSourceIP(r),
		UserAgent:        r.UserAgent(),
		SessionContext:   sessionContext,
		ServiceContext: map[string]string{
			"region":  parsedReq.GetRegion(),
			"service": serviceName,
		},
	}
}

// extractRoleNameFromAssumedRoleARN returns the role name from an assumed-
// role ARN (arn:aws:iam::<account>:role/<name>/<session>) or empty when
// the ARN does not match the role shape.
func extractRoleNameFromAssumedRoleArn(arn string) string {
	// arn:aws:iam::<account>:role/<name> or
	// arn:aws:sts::<account>:assumed-role/<name>/<session>
	const stsAssumedPrefix = ":assumed-role/"
	if idx := strings.Index(arn, stsAssumedPrefix); idx >= 0 {
		rest := arn[idx+len(stsAssumedPrefix):]
		if slash := strings.Index(rest, "/"); slash >= 0 {
			return rest[:slash]
		}
		return rest
	}
	const roleSuffix = ":role/"
	if idx := strings.Index(arn, roleSuffix); idx >= 0 {
		return arn[idx+len(roleSuffix):]
	}
	return ""
}

// extractAccountIDFromArn returns the AWS account ID embedded in an IAM/STS
// ARN (the 5th colon-separated segment) or empty when the ARN is malformed.
func extractAccountIDFromArn(arn string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}

// extractUserNameFromArn extracts the IAM user name from an ARN of the form
// arn:aws:iam::<account>:user/<userName>. Returns empty for other ARN types.
func extractUserNameFromArn(arn string) string {
	const userSuffix = ":user/"
	if idx := strings.Index(arn, userSuffix); idx >= 0 {
		return arn[idx+len(userSuffix):]
	}
	return ""
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
	a.stopOnce.Do(func() { close(a.stopCleanup) })
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
				func() {
					defer func() {
						if re := recover(); re != nil {
							logs.Error("authorizer cache cleanup panic recovered", logs.Any("panic", re))
						}
					}()
					a.cleanupExpiredEntries()
				}()
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
		logs.Info("Cleaned up expired cache entries", logs.Int("expired", len(toDelete)), logs.Int("remaining", count-len(toDelete)))
	}

	if count > a.maxCacheSize {
		a.evictOldestEntries(count-a.maxCacheSize, count)
	}
}

func (a *Authorizer) enforceCacheSize() {
	count := 0
	a.policyCache.Range(func(_, _ interface{}) bool {
		count++
		return count <= a.maxCacheSize*2
	})

	if count <= a.maxCacheSize {
		return
	}

	a.evictOldestEntries(count-a.maxCacheSize, count)
}

func (a *Authorizer) evictOldestEntries(toEvict int, maxScan int) {
	type entry struct {
		key      string
		cachedAt time.Time
	}

	entries := make([]entry, 0, toEvict)
	a.policyCache.Range(func(key, value interface{}) bool {
		if cp, ok := value.(*cachedPolicies); ok {
			entries = append(entries, entry{key: key.(string), cachedAt: cp.cachedAt})
		}
		if len(entries) >= maxScan {
			return false
		}
		return true
	})

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].cachedAt.Before(entries[j].cachedAt)
	})

	evicted := 0
	for i := 0; i < toEvict && i < len(entries); i++ {
		a.policyCache.Delete(entries[i].key)
		evicted++
	}

	if evicted > 0 {
		logs.Info("Evicted oldest cache entries", logs.Int("evicted", evicted))
	}
}

// trustForwardedHeaders controls whether X-Forwarded-For / X-Real-IP
// headers are trusted. Default is false (secure) — use RemoteAddr.
// Set TRUST_FORWARDED_HEADERS=true when deploying behind a reverse
// proxy (nginx, ALB) that overwrites client-supplied forwarding headers.
var trustForwardedHeaders = os.Getenv("TRUST_FORWARDED_HEADERS") == "true"

func extractSourceIP(r *http.Request) string {
	if trustForwardedHeaders {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ips := strings.Split(xff, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
