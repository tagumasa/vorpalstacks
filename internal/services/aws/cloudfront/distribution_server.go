package cloudfront

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/headerorder"
	"vorpalstacks/internal/common/waflimits"
	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/server/fqdnrouter"
	cfstore "vorpalstacks/internal/store/aws/cloudfront"
)

// DistributionServer serves requests for CloudFront distributions by proxying
// to the configured origin, applying cache behaviour headers, and serving
// repeat requests from the edge cache.
type DistributionServer struct {
	storageManager     *storage.RegionStorageManager
	accountID          string
	client             *http.Client
	distributionMu     sync.RWMutex
	distribution       *cfstore.DistributionStore
	policyMu           sync.Mutex
	cachePolicies      *cfstore.CachePolicyStore
	deploymentPolicies *cfstore.ContinuousDeploymentPolicyStore
	cache              *responseCache
	inspectorMu        sync.RWMutex
	inspector          eventbus.WebACLInspector
	providerMu         sync.RWMutex
	acmCertificates    eventbus.ACMCertificateProvider
	iamCertificates    eventbus.IAMServerCertificateProvider
	certCache          *tlsCertificateCache
}

// NewDistributionServer creates a new DistributionServer. CloudFront is a
// global service, so the server always reads from the global Pebble DB.
func NewDistributionServer(storageManager *storage.RegionStorageManager, accountID string) *DistributionServer {
	return &DistributionServer{
		storageManager: storageManager,
		accountID:      accountID,
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		cache:     newResponseCache(),
		certCache: newTLSCertificateCache(),
	}
}

// SetDistributionStore injects a DistributionStore instance, typically the same
// one used by CloudFrontService, so that both components share a single store.
func (s *DistributionServer) SetDistributionStore(store *cfstore.DistributionStore) {
	s.distributionMu.Lock()
	s.distribution = store
	s.distributionMu.Unlock()
}

// SetWebACLInspector injects the WAF request-inspection entry point. The
// wiring runs after the WAFv2 service initialises and before listeners
// start serving traffic.
func (s *DistributionServer) SetWebACLInspector(inspector eventbus.WebACLInspector) {
	s.inspectorMu.Lock()
	s.inspector = inspector
	s.inspectorMu.Unlock()
}

func (s *DistributionServer) webACLInspector() eventbus.WebACLInspector {
	s.inspectorMu.RLock()
	defer s.inspectorMu.RUnlock()
	return s.inspector
}

// enforceWebACL inspects the request against the WebACL associated with
// the distribution. It answers the blocked response itself and returns
// false; an allowed request returns true with any custom request
// headers applied to r. Inspection failures fail open (the request is
// served): the AWS WAF Developer Guide documents that when WAF
// encounters an internal error, CloudFront typically allows the request
// or serves the content, while Regional services typically deny it.
// The body is buffered up to the inspection limit for the evaluation
// and recombined with the unread tail so origin proxying still streams
// the full request.
func (s *DistributionServer) enforceWebACL(w http.ResponseWriter, r *http.Request, distributionARN string) bool {
	inspector := s.webACLInspector()
	if inspector == nil {
		return true
	}

	var body []byte
	bodyTruncated := false
	if r.Body != nil {
		buffered, err := io.ReadAll(io.LimitReader(r.Body, waflimits.DefaultBodyInspectionLimit+1))
		if err != nil {
			// The body cannot be read for inspection; with no evidence
			// to match on, serve the request.
			logs.Warn("cloudfront waf inspection body read failed", logs.String("distribution", distributionARN), logs.Err(err))
			return true
		}
		body = buffered
		if int64(len(body)) > waflimits.DefaultBodyInspectionLimit {
			body = body[:waflimits.DefaultBodyInspectionLimit]
			bodyTruncated = true
			r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buffered), r.Body))
		} else {
			r.Body = io.NopCloser(bytes.NewReader(buffered))
		}
	}

	inspHeaders := eventbus.RequestHeadersWithHost(r.Header, r.Host)
	headerOrder, _ := headerorder.FromContext(r.Context(), inspHeaders)
	result, err := inspector.InspectWebACLRequest(r.Context(), "", distributionARN, eventbus.BuildWebACLInspectionRequest(
		r.Method, r.URL.Path, r.URL.RawQuery, remoteAddrHost(r.RemoteAddr), r.Proto,
		inspHeaders, headerOrder, body, bodyTruncated,
	))
	if err != nil {
		logs.Warn("cloudfront waf inspection failed, serving request", logs.String("distribution", distributionARN), logs.Err(err))
		return true
	}
	if result.Interrupts() {
		status := result.ResponseCode
		if status == 0 {
			// The default Block action response is 403 Forbidden.
			status = http.StatusForbidden
		}
		for _, h := range result.ResponseHeaders {
			w.Header().Set(h.Name, h.Value)
		}
		w.WriteHeader(status)
		if result.ResponseBody != "" {
			_, _ = io.WriteString(w, result.ResponseBody)
		}
		return false
	}
	for _, h := range result.InsertHeaders {
		// Inserted header names arrive prefixed with x-amzn-waf-, so
		// adding cannot overwrite a header the client sent.
		r.Header.Add(h.Name, h.Value)
	}
	return true
}

func remoteAddrHost(remoteAddr string) string {
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return remoteAddr
}

// HandleRequest proxies an incoming request to the matching CloudFront
// distribution's origin. Cacheable methods are first looked up in the
// edge cache; misses populate it per the behaviour TTL settings and the
// origin Cache-Control/Expires directives.
func (s *DistributionServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if eventbus.ServeWAFTokenExchange(r.Context(), s.webACLInspector(), w, r) {
		return
	}

	distributionID := fqdnrouter.ResourceIDFromContext(r.Context())
	if distributionID == "" {
		distributionID = s.resolveDistributionID(r.Host)
	}
	if distributionID == "" {
		http.Error(w, `{"message":"Distribution not found"}`, http.StatusNotFound)
		return
	}

	store := s.getDistributionStore()
	dist, err := store.Get(distributionID)
	if err != nil || dist == nil {
		http.Error(w, `{"message":"Distribution not found"}`, http.StatusNotFound)
		return
	}

	if !dist.Enabled {
		http.Error(w, `{"message":"Distribution is disabled"}`, http.StatusForbidden)
		return
	}

	config := dist.DistributionConfig
	if config == nil {
		http.Error(w, `{"message":"Distribution has no configuration"}`, http.StatusInternalServerError)
		return
	}

	if dist.Staging {
		// Viewers cannot reach a staging distribution directly; only the
		// primary distribution's continuous deployment policy routes
		// traffic to it.
		http.Error(w, `{"message":"Access denied. Staging distributions are not directly accessible."}`, http.StatusForbidden)
		return
	}

	if !s.enforceWebACL(w, r, dist.ARN) {
		return
	}

	// Continuous deployment: an attached, enabled policy routes part of
	// the traffic to the staging distribution. The staging distribution
	// has its own cache, so the routing decision also switches the
	// distribution identity the rest of the handler uses.
	if config.ContinuousDeploymentPolicyId != "" {
		if staged := s.applyContinuousDeployment(w, r, store, config); staged != nil {
			dist = staged
			distributionID = dist.ID
			config = dist.DistributionConfig
			if config == nil {
				http.Error(w, `{"message":"Distribution has no configuration"}`, http.StatusInternalServerError)
				return
			}
		}
	}

	requestPath := r.URL.Path
	if requestPath == "" || requestPath == "/" {
		if config.DefaultRootObject != "" {
			requestPath = "/" + config.DefaultRootObject
		}
	}

	behavior := s.matchCacheBehavior(config, requestPath)
	if behavior == nil {
		http.Error(w, `{"message":"No matching cache behavior"}`, http.StatusNotFound)
		return
	}

	// The behaviour's viewer protocol policy applies to plain-HTTP
	// requests on both planes: https-only refuses them, redirect-to-https
	// answers with the redirect contract before any cache or origin work.
	if r.TLS == nil && !s.enforceViewerProtocolPolicy(w, r, behavior.ViewerProtocolPolicy) {
		return
	}

	origin := s.resolveOrigin(config.Origins, behavior.TargetOriginId)
	if origin == nil {
		http.Error(w, `{"message":"Origin not found"}`, http.StatusNotFound)
		return
	}

	policy := s.cachePolicyFor(behavior)
	ttls := behaviourTTLs(behavior, policy)
	keys := buildKeyPolicy(behavior, policy)
	cacheable := methodCacheable(r.Method, behavior)

	var cacheKey string
	if cacheable {
		cacheKey = buildCacheKey(distributionID, behavior.PathPattern, requestPath, r, keys)
		// A fresh entry answers without an origin round trip. Range
		// requests are served by slicing the cached body; partial
		// origin responses are never stored, so entries always hold
		// complete bodies.
		if entry := s.cache.lookup(cacheKey); entry != nil && !entry.expired(time.Now()) {
			s.serveFromCache(w, r, entry, distributionID)
			return
		}
	}

	viewerScheme := "http"
	if r.TLS != nil {
		viewerScheme = "https"
	}
	targetURL := s.buildOriginURL(origin, requestPath, keys.originRawQuery(r.URL.RawQuery), viewerScheme)

	originReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, `{"message":"Failed to create origin request"}`, http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			originReq.Header.Add(key, value)
		}
	}
	if cookieHeader := keys.cookieHeaderForOrigin(r.Header.Get("Cookie")); cookieHeader != "" {
		originReq.Header.Set("Cookie", cookieHeader)
	} else {
		originReq.Header.Del("Cookie")
	}
	originReq.Header.Set("X-Forwarded-Host", r.Host)
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		originReq.Header.Set("X-Forwarded-For", host)
	} else {
		originReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	}

	resp, err := s.client.Do(originReq)
	if err != nil {
		logs.Error("CloudFront origin request failed", logs.String("url", targetURL), logs.Err(err))
		// An unreachable origin serves an expired entry for the
		// error-caching duration when one is present.
		if entry := s.cache.lookup(cacheKey); entry != nil && entry.expired(time.Now()) {
			s.cache.rearm(cacheKey, errorCacheTTL(config.CustomErrorResponses, 504, originDirectives{}))
			s.serveFromCache(w, r, entry, distributionID)
			return
		}
		http.Error(w, `{"message":"Origin request failed"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	now := time.Now()

	// An expired entry plus a 5xx origin response serves the stale
	// object; a 4xx response replaces the entry instead.
	if resp.StatusCode >= 500 {
		if entry := s.cache.lookup(cacheKey); entry != nil && entry.expired(now) {
			directives := parseOriginDirectives(resp.Header)
			s.cache.rearm(cacheKey, errorCacheTTL(config.CustomErrorResponses, resp.StatusCode, directives))
			s.serveFromCache(w, r, entry, distributionID)
			return
		}
	}

	for key, values := range resp.Header {
		if hopByHopHeader(key) {
			continue
		}
		// A response whose behaviour forwards no cookies must not carry
		// the origin's Set-Cookie. Filtering at copy time (rather than
		// deleting afterwards) preserves the session-stickiness cookie a
		// continuous deployment policy may have added.
		if !keys.cookiesForwarded && strings.EqualFold(key, "Set-Cookie") {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("X-Cache", "Miss from cloudfront")
	w.Header().Set("Via", "1.1 "+distributionID+".cloudfront.net (CloudFront)")

	// Cacheable full-body responses are buffered so they can be stored;
	// everything else streams.
	var body []byte
	var prebuffered []byte
	if cacheable && r.Method != http.MethodHead && r.Header.Get("Range") == "" {
		buffered, readErr := io.ReadAll(io.LimitReader(resp.Body, maxCachedObjectBytes+1))
		switch {
		case readErr == nil && len(buffered) <= maxCachedObjectBytes:
			body = buffered
		case readErr == nil:
			// Oversized: stream what was read plus the unread tail.
			prebuffered = buffered
		}
	}

	if body != nil {
		directives := parseOriginDirectives(resp.Header)
		if entry, ok := buildCacheEntry(distributionID, cacheKey, behavior.PathPattern, requestPath, resp, body, directives, ttls, keys, config, now); ok {
			s.cache.put(entry)
		}
		w.WriteHeader(resp.StatusCode)
		if r.Method != http.MethodHead {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(resp.StatusCode)
	if len(prebuffered) > 0 {
		_, _ = w.Write(prebuffered)
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		logs.Warn("cloudfront distribution proxy copy error", logs.String("distributionId", distributionID), logs.Err(err))
	}
}

// serveFromCache writes a cached entry to the viewer with the Hit
// marker, the entry age, and the CloudFront Via header. HEAD requests
// receive the headers only; a single-range GET against a 200 entry is
// answered with 206 and the requested slice.
func (s *DistributionServer) serveFromCache(w http.ResponseWriter, r *http.Request, entry *cacheEntry, distributionID string) {
	for key, values := range entry.header {
		// Filter at copy time (rather than deleting afterwards) so the
		// session-stickiness cookie of a continuous deployment policy
		// survives Set-Cookie stripping.
		if entry.stripSetCookie && strings.EqualFold(key, "Set-Cookie") {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Del("Transfer-Encoding")

	age := int(time.Since(entry.fetchedAt).Round(time.Second) / time.Second)
	if age < 0 {
		age = 0
	}
	w.Header().Set("Age", strconv.Itoa(age))
	w.Header().Set("X-Cache", "Hit from cloudfront")
	w.Header().Set("Via", "1.1 "+distributionID+".cloudfront.net (CloudFront)")

	if header := r.Header.Get("Range"); header != "" && entry.status == http.StatusOK {
		if start, length, ok := singleByteRange(header, len(entry.body)); ok {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, start+length-1, len(entry.body)))
			w.Header().Set("Content-Length", strconv.Itoa(length))
			w.WriteHeader(http.StatusPartialContent)
			if r.Method != http.MethodHead {
				_, _ = w.Write(entry.body[start : start+length])
			}
			return
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(entry.body)))
	w.WriteHeader(entry.status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(entry.body)
	}
}

// buildCacheEntry decides whether an origin response is cacheable and
// constructs the entry. Object statuses use the resolved TTL; error
// statuses follow the error-caching rules.
func buildCacheEntry(distributionID, cacheKey, pathPattern, requestPath string, resp *http.Response, body []byte, directives originDirectives, ttls ttlSet, keys keyPolicy, config *cfstore.DistributionConfig, now time.Time) (*cacheEntry, bool) {
	var ttl time.Duration
	switch {
	case isObjectStatus(resp.StatusCode):
		ttl = resolveTTL(directives, ttls)
		if ttl <= 0 {
			return nil, false
		}
	case errorStatusCacheable(resp.StatusCode, directives, ttls):
		ttl = errorCacheTTL(config.CustomErrorResponses, resp.StatusCode, directives)
	default:
		return nil, false
	}

	header := make(http.Header, len(resp.Header))
	for key, values := range resp.Header {
		if hopByHopHeader(key) {
			continue
		}
		header[key] = append([]string(nil), values...)
	}

	return &cacheEntry{
		key:            cacheKey,
		distribution:   distributionID,
		path:           requestPath,
		status:         resp.StatusCode,
		header:         header,
		body:           body,
		fetchedAt:      now,
		expiresAt:      now.Add(ttl),
		stripSetCookie: !keys.cookiesForwarded,
	}, true
}

// cachePolicyFor resolves the cache policy attached to a behaviour, if
// any. A missing or unknown policy id leaves the deprecated behaviour
// fields in effect.
func (s *DistributionServer) cachePolicyFor(behavior *cfstore.CacheBehavior) *cfstore.CachePolicy {
	if behavior == nil || behavior.CachePolicyId == "" {
		return nil
	}
	store := s.getCachePolicyStore()
	if store == nil {
		return nil
	}
	policy, err := store.Get(behavior.CachePolicyId)
	if err != nil || policy == nil {
		return nil
	}
	return policy
}

func (s *DistributionServer) getCachePolicyStore() *cfstore.CachePolicyStore {
	s.policyMu.Lock()
	defer s.policyMu.Unlock()
	if s.cachePolicies != nil {
		return s.cachePolicies
	}
	st, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		logs.Error("CloudFront distribution server: failed to get global storage for cache policies", logs.Err(err))
		return nil
	}
	s.cachePolicies = cfstore.NewCachePolicyStore(st, s.accountID)
	return s.cachePolicies
}

func (s *DistributionServer) extractDistributionID(host string) string {
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if len(parts) >= 1 && parts[0] != "" && parts[0] != "localhost" {
		return parts[0]
	}
	return ""
}

// resolveDistributionID maps the request Host to a distribution: a host
// under .cloudfront.net (or a bare single-label host) carries the
// distribution ID as its first label, and any other fully-qualified host
// resolves through the configured CNAME aliases.
func (s *DistributionServer) resolveDistributionID(host string) string {
	hostname := strings.ToLower(strings.Split(host, ":")[0])
	hostname = strings.TrimSuffix(hostname, ".")
	if strings.HasSuffix(hostname, ".cloudfront.net") || !strings.Contains(hostname, ".") {
		return s.extractDistributionID(hostname)
	}
	store := s.getDistributionStore()
	if store == nil {
		return ""
	}
	if dist, err := store.GetByAlias(hostname); err == nil && dist != nil {
		return dist.ID
	}
	return ""
}

func (s *DistributionServer) matchCacheBehavior(config *cfstore.DistributionConfig, path string) *cfstore.CacheBehavior {
	if config.CacheBehaviors != nil {
		for _, behavior := range config.CacheBehaviors.Items {
			if s.pathMatches(behavior.PathPattern, path) {
				return behavior
			}
		}
	}
	return config.DefaultCacheBehavior
}

func (s *DistributionServer) pathMatches(pattern, path string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return path == pattern
}

func (s *DistributionServer) resolveOrigin(origins cfstore.Origins, targetOriginID string) *cfstore.Origin {
	for _, origin := range origins.Items {
		if origin.ID == targetOriginID {
			return origin
		}
	}
	return nil
}

func (s *DistributionServer) buildOriginURL(origin *cfstore.Origin, path, query string, viewerScheme string) string {
	scheme := "http"
	if origin.CustomOriginConfig != nil {
		switch origin.CustomOriginConfig.OriginProtocolPolicy {
		case "https-only":
			scheme = "https"
		case "match-viewer":
			scheme = viewerScheme
		}
	}

	host := origin.DomainName
	if strings.Contains(host, "s3-website") || strings.Contains(host, ".s3-website.") {
		host = fmt.Sprintf("localhost:%d", config.GetInt("ports.s3_website"))
	} else if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
	} else if strings.Contains(host, "lambda-url") {
		host = fmt.Sprintf("localhost:%d", config.GetInt("ports.lambda_url"))
	} else if strings.Contains(host, "s3.") && !strings.Contains(host, "website") {
		host = fmt.Sprintf("localhost:%d", config.ServerPort())
	}

	originPath := strings.TrimSuffix(origin.OriginPath, "/")
	url := fmt.Sprintf("%s://%s%s%s", scheme, host, originPath, path)
	if query != "" {
		url += "?" + query
	}
	return url
}

func (s *DistributionServer) getDistributionStore() *cfstore.DistributionStore {
	s.distributionMu.RLock()
	if s.distribution != nil {
		cached := s.distribution
		s.distributionMu.RUnlock()
		return cached
	}
	s.distributionMu.RUnlock()

	s.distributionMu.Lock()
	defer s.distributionMu.Unlock()
	if s.distribution != nil {
		return s.distribution
	}
	store, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		logs.Error("CloudFront distribution server: failed to get global storage, returning empty store", logs.Err(err))
		return s.distribution
	}
	s.distribution = cfstore.NewDistributionStore(store, s.accountID)
	return s.distribution
}

// stickyCDPCookieName is the session-stickiness cookie a weight-based
// continuous deployment policy pins its routing assignment with.
const stickyCDPCookieName = "aws-cf-cd-sticky"

// applyContinuousDeployment evaluates the primary distribution's
// continuous deployment policy: an enabled policy routes matching
// requests to the staging distribution. Header-based policies route the
// requests carrying the configured header and value; weight-based
// policies assign viewers deterministically by client IP and, when
// session stickiness is configured, pin the assignment with a cookie
// bounded by the idle and maximum TTLs. The returned distribution is
// the staging distribution to serve from; nil keeps the primary.
func (s *DistributionServer) applyContinuousDeployment(w http.ResponseWriter, r *http.Request, store *cfstore.DistributionStore, config *cfstore.DistributionConfig) *cfstore.Distribution {
	policies := s.getContinuousDeploymentPolicyStore()
	if policies == nil {
		return nil
	}
	policy, err := policies.Get(config.ContinuousDeploymentPolicyId)
	if err != nil || policy.ContinuousDeploymentPolicyConfig == nil {
		return nil
	}
	policyConfig := policy.ContinuousDeploymentPolicyConfig
	if !policyConfig.Enabled || policyConfig.TrafficConfig == nil ||
		policyConfig.StagingDistributionDnsNames == nil || len(policyConfig.StagingDistributionDnsNames.Items) == 0 {
		return nil
	}
	staging, err := store.GetByDomainName(policyConfig.StagingDistributionDnsNames.Items[0])
	if err != nil || staging == nil || !staging.Staging || !staging.Enabled || staging.DistributionConfig == nil {
		return nil
	}

	switch policyConfig.TrafficConfig.Type {
	case "SingleHeader":
		if header := policyConfig.TrafficConfig.SingleHeaderConfig; header != nil && r.Header.Get(header.Header) == header.Value {
			return staging
		}
	case "SingleWeight":
		weight := policyConfig.TrafficConfig.SingleWeightConfig
		if weight == nil || weight.Weight <= 0 {
			return nil
		}
		useStaging := weightBasedAssignment(r, weight.Weight)
		if stickiness := weight.SessionStickinessConfig; stickiness != nil {
			useStaging = stickyAssignment(w, r, useStaging, stickiness, time.Now())
		}
		if useStaging {
			return staging
		}
	}
	return nil
}

// weightBasedAssignment maps a viewer deterministically into the weight
// window: FNV-1a over the client IP yields a position in the traffic
// space, and positions below the weight route to staging. The same
// viewer keeps its assignment across requests.
func weightBasedAssignment(r *http.Request, weight float64) bool {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(remoteAddrHost(r.RemoteAddr)))
	fraction := float64(hash.Sum32()%stickyAssignmentSpace) / float64(stickyAssignmentSpace)
	return fraction < weight
}

const stickyAssignmentSpace = 10000

// stickyAssignment honours and refreshes the session-stickiness cookie.
// The cookie carries the session start and the last refresh, so the
// idle TTL is measured from the last request while the maximum TTL is
// measured from the first assignment: a session the viewer keeps
// refreshing still ends once the maximum session time passes, per the
// SessionStickinessConfig semantics. Anything outside the bounds takes
// the fresh weight assignment and reissues the cookie.
func stickyAssignment(w http.ResponseWriter, r *http.Request, assigned bool, cfg *cfstore.SessionStickinessConfig, now time.Time) bool {
	if cookie, err := r.Cookie(stickyCDPCookieName); err == nil && cookie.Value != "" {
		value, start, last, ok := parseStickyCookie(cookie.Value)
		if ok {
			idleAge, maxAge := now.Sub(last), now.Sub(start)
			if idleAge >= 0 && maxAge >= 0 &&
				idleAge <= time.Duration(cfg.IdleTTL)*time.Second &&
				maxAge <= time.Duration(cfg.MaximumTTL)*time.Second {
				w.Header().Add("Set-Cookie", formatStickyCookie(value, start, now, cfg))
				return value
			}
		}
	}
	w.Header().Add("Set-Cookie", formatStickyCookie(assigned, now, now, cfg))
	return assigned
}

func formatStickyCookie(staging bool, start, last time.Time, cfg *cfstore.SessionStickinessConfig) string {
	value := "primary"
	if staging {
		value = "staging"
	}
	return fmt.Sprintf("%s=%s.%d.%d; Path=/; Max-Age=%d", stickyCDPCookieName, value, start.Unix(), last.Unix(), cfg.MaximumTTL)
}

func parseStickyCookie(value string) (staging bool, start, last time.Time, ok bool) {
	assignment, rest, found := strings.Cut(value, ".")
	if !found {
		return false, time.Time{}, time.Time{}, false
	}
	startSeconds, lastSeconds, found := strings.Cut(rest, ".")
	if !found {
		return false, time.Time{}, time.Time{}, false
	}
	startParsed, err := strconv.ParseInt(startSeconds, 10, 64)
	if err != nil || (assignment != "staging" && assignment != "primary") {
		return false, time.Time{}, time.Time{}, false
	}
	lastParsed, err := strconv.ParseInt(lastSeconds, 10, 64)
	if err != nil {
		return false, time.Time{}, time.Time{}, false
	}
	return assignment == "staging", time.Unix(startParsed, 0), time.Unix(lastParsed, 0), true
}

// getContinuousDeploymentPolicyStore lazily builds the policy store from
// the global storage; policyMu also guards this construction.
func (s *DistributionServer) getContinuousDeploymentPolicyStore() *cfstore.ContinuousDeploymentPolicyStore {
	s.policyMu.Lock()
	defer s.policyMu.Unlock()
	if s.deploymentPolicies != nil {
		return s.deploymentPolicies
	}
	st, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		logs.Error("CloudFront distribution server: failed to get global storage for deployment policies", logs.Err(err))
		return nil
	}
	s.deploymentPolicies = cfstore.NewContinuousDeploymentPolicyStore(st)
	return s.deploymentPolicies
}
