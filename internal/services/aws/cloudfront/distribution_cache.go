package cloudfront

import (
	"container/list"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	cfstore "vorpalstacks/internal/store/aws/cloudfront"
)

// Platform bounds for the in-process edge cache. Unlike the AWS quotas in
// the store package these are vorpalstacks resource limits: cached
// responses live in server memory, so both the entry count and the
// per-object size are capped.
const (
	maxCacheEntries      = 4096
	maxCachedObjectBytes = 10 << 20
)

// cacheEntry is one stored origin response. expiresAt is the deadline the
// resolved TTL produced; an expired entry stays in the map solely to serve
// the documented stale-on-origin-failure path.
type cacheEntry struct {
	key            string
	distribution   string
	path           string
	status         int
	header         http.Header
	body           []byte
	fetchedAt      time.Time
	expiresAt      time.Time
	stripSetCookie bool
}

func (e *cacheEntry) expired(now time.Time) bool {
	return !now.Before(e.expiresAt)
}

// responseCache is the per-server edge cache. Entries are keyed by the
// full CloudFront cache key (distribution, cache behaviour, path, and the
// forwarded query/cookie/header components) and evicted least-recently-
// used beyond maxCacheEntries.
type responseCache struct {
	mu      sync.Mutex
	entries map[string]*list.Element
	order   *list.List
}

func newResponseCache() *responseCache {
	return &responseCache{
		entries: map[string]*list.Element{},
		order:   list.New(),
	}
}

// lookup returns the entry for a key regardless of expiry, or nil. A hit
// refreshes the entry's LRU position.
func (c *responseCache) lookup(key string) *cacheEntry {
	if key == "" {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.entries[key]
	if !ok {
		return nil
	}
	c.order.MoveToFront(el)
	return el.Value.(*cacheEntry)
}

// put inserts or replaces an entry and evicts past the entry bound.
func (c *responseCache) put(entry *cacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[entry.key]; ok {
		c.order.Remove(el)
		delete(c.entries, entry.key)
	}
	el := c.order.PushFront(entry)
	c.entries[entry.key] = el
	for c.order.Len() > maxCacheEntries {
		last := c.order.Back()
		if last == nil {
			break
		}
		c.order.Remove(last)
		delete(c.entries, last.Value.(*cacheEntry).key)
	}
}

// rearm extends an expired entry's usability by the error-caching
// duration. The documented stale-serving window resets each time the
// origin keeps failing, so the deadline moves rather than extends from
// the original expiry.
func (c *responseCache) rearm(key string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.entries[key]
	if !ok {
		return
	}
	entry := el.Value.(*cacheEntry)
	entry.expiresAt = time.Now().Add(ttl)
	c.order.MoveToFront(el)
}

// purge removes the entries of one distribution whose request path matches
// any invalidation path. It returns the number of removed entries.
func (c *responseCache) purge(distributionID string, paths []string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	removed := 0
	for key, el := range c.entries {
		entry := el.Value.(*cacheEntry)
		if entry.distribution != distributionID {
			continue
		}
		for _, pattern := range paths {
			if invalidationPathMatches(pattern, entry.path) {
				c.order.Remove(el)
				delete(c.entries, key)
				removed++
				break
			}
		}
	}
	return removed
}

// purgeDistribution removes every entry of one distribution. Distribution
// configuration changes invalidate the whole edge cache for that
// distribution.
func (c *responseCache) purgeDistribution(distributionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, el := range c.entries {
		if el.Value.(*cacheEntry).distribution == distributionID {
			c.order.Remove(el)
			delete(c.entries, key)
		}
	}
}

// InvalidatePaths applies an invalidation batch to the live cache.
func (s *DistributionServer) InvalidatePaths(distributionID string, paths []string) {
	s.cache.purge(distributionID, paths)
}

// PurgeDistribution drops all cached entries of one distribution.
func (s *DistributionServer) PurgeDistribution(distributionID string) {
	s.cache.purgeDistribution(distributionID)
}

// invalidationPathMatches reports whether an invalidation Path matches a
// cached request path. Invalidation paths start with "/" and may contain
// "*" wildcards; "/*" matches every path of the distribution.
func invalidationPathMatches(pattern, path string) bool {
	if pattern == "" {
		return false
	}
	if pattern == "/*" {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return pattern == path
	}
	parts := strings.Split(pattern, "*")
	if !strings.HasPrefix(path, parts[0]) {
		return false
	}
	rest := path[len(parts[0]):]
	for i := 1; i < len(parts)-1; i++ {
		idx := strings.Index(rest, parts[i])
		if idx < 0 {
			return false
		}
		rest = rest[idx+len(parts[i]):]
	}
	return strings.HasSuffix(rest, parts[len(parts)-1])
}

// ttlSet holds the TTL bounds in effect for a request: from the attached
// cache policy when present, from the deprecated behaviour fields
// otherwise. A zero max means no upper bound.
type ttlSet struct {
	min time.Duration
	max time.Duration
	def time.Duration
}

// disabled reports the all-zero configuration (the CachingDisabled
// shape) under which no response is cached at all, errors included.
func (t ttlSet) disabled() bool {
	return t.min == 0 && t.max == 0 && t.def == 0
}

func secondsToDuration(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

func behaviourTTLs(behavior *cfstore.CacheBehavior, policy *cfstore.CachePolicy) ttlSet {
	var minSeconds, maxSeconds, defaultSeconds int64
	if behavior.CachePolicyId != "" && policy != nil && policy.CachePolicyConfig != nil {
		minSeconds = policy.CachePolicyConfig.MinTTL
		maxSeconds = policy.CachePolicyConfig.MaxTTL
		defaultSeconds = policy.CachePolicyConfig.DefaultTTL
	} else {
		minSeconds = int64(behavior.MinTTL)
		maxSeconds = int64(behavior.MaxTTL)
		defaultSeconds = int64(behavior.DefaultTTL)
	}
	return ttlSet{
		min: secondsToDuration(minSeconds),
		max: secondsToDuration(maxSeconds),
		def: secondsToDuration(defaultSeconds),
	}
}

// originDirectives carries the cache-relevant directives of one origin
// response.
type originDirectives struct {
	noStore       bool
	hasMaxAge     bool
	maxAge        time.Duration
	hasSMaxAge    bool
	sMaxAge       time.Duration
	hasExpires    bool
	expiresOffset time.Duration // Expires relative to the fetch time
}

// parseOriginDirectives extracts the s-maxage, max-age, no-cache/no-store/
// private, and Expires values from an origin response per the documented
// precedence: s-maxage over max-age, and Cache-Control over Expires.
func parseOriginDirectives(header http.Header) originDirectives {
	var d originDirectives
	for _, value := range header.Values("Cache-Control") {
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			name := part
			val := ""
			if i := strings.Index(part, "="); i >= 0 {
				name = strings.TrimSpace(part[:i])
				val = strings.Trim(strings.TrimSpace(part[i+1:]), "\"")
			}
			switch strings.ToLower(name) {
			case "no-store", "no-cache", "private":
				d.noStore = true
			case "s-maxage":
				if seconds, err := strconv.ParseInt(val, 10, 64); err == nil && seconds >= 0 {
					d.hasSMaxAge = true
					d.sMaxAge = secondsToDuration(seconds)
				}
			case "max-age":
				if seconds, err := strconv.ParseInt(val, 10, 64); err == nil && seconds >= 0 {
					d.hasMaxAge = true
					d.maxAge = secondsToDuration(seconds)
				}
			}
		}
	}
	if expires := header.Get("Expires"); expires != "" {
		if t, err := http.ParseTime(expires); err == nil {
			d.hasExpires = true
			d.expiresOffset = time.Until(t)
		}
	}
	return d
}

// resolveTTL returns the cache duration for an object response. With a
// no-cache/no-store/private directive the response is cached only when
// the minimum TTL forces it (the documented warning); otherwise the
// origin directive (or the default TTL when absent) is clamped to the
// minimum and maximum bounds.
func resolveTTL(d originDirectives, ttls ttlSet) time.Duration {
	if d.noStore {
		if ttls.min > 0 {
			return ttls.min
		}
		return 0
	}
	var origin time.Duration
	switch {
	case d.hasSMaxAge:
		origin = d.sMaxAge
	case d.hasMaxAge:
		origin = d.maxAge
	case d.hasExpires:
		origin = d.expiresOffset
	default:
		origin = ttls.def
	}
	if origin < ttls.min {
		origin = ttls.min
	}
	if ttls.max > 0 && origin > ttls.max {
		origin = ttls.max
	}
	if origin < 0 {
		return 0
	}
	return origin
}

// Statuses whose responses CloudFront caches as objects (HTTP 200 plus
// the origin redirects the custom-origin behaviour documents as cached).
func isObjectStatus(status int) bool {
	switch status {
	case http.StatusOK, http.StatusMovedPermanently, http.StatusFound,
		http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return true
	}
	return false
}

// errorStatusCacheable reports whether a 4xx/5xx response is cacheable:
// 404/414/500/501/502/503/504 always, and 400/403/405/412/415 only when
// the origin sent a Cache-Control max-age or s-maxage directive. An
// all-zero TTL configuration (the CachingDisabled shape) caches nothing.
func errorStatusCacheable(status int, d originDirectives, ttls ttlSet) bool {
	if ttls.disabled() {
		return false
	}
	switch status {
	case http.StatusNotFound, 414, 500, 501, 502, 503, 504:
		return true
	case http.StatusBadRequest, http.StatusForbidden, http.StatusMethodNotAllowed, 412, 415:
		return d.hasSMaxAge || d.hasMaxAge
	}
	return false
}

// errorCacheTTL is the duration a cached error response keeps: the
// maximum of the configured ErrorCachingMinTTL (10 seconds by default)
// and the origin Cache-Control max-age or s-maxage directive.
func errorCacheTTL(custom *cfstore.CustomErrorResponses, status int, d originDirectives) time.Duration {
	ttl := secondsToDuration(cfstore.DefaultErrorCachingTTLSeconds)
	if custom != nil {
		for _, item := range custom.Items {
			if item.ErrorCode == status && item.ErrorCachingMinTTL > 0 {
				ttl = secondsToDuration(int64(item.ErrorCachingMinTTL))
			}
		}
	}
	if d.hasSMaxAge && d.sMaxAge > ttl {
		ttl = d.sMaxAge
	}
	if d.hasMaxAge && d.maxAge > ttl {
		ttl = d.maxAge
	}
	return ttl
}

// keyPolicy captures which request components form the cache key and
// reach the origin for one cache behaviour. Query keys under the legacy
// ForwardedValues settings are cache-key-only: every parameter is
// forwarded while caching keys on the listed subset.
type keyPolicy struct {
	forwardAllQuery  bool
	keyAllQuery      bool
	keyQuery         []string
	excludeQuery     []string
	cookiesForwarded bool
	keyAllCookies    bool
	keyCookies       []string
	excludeCookies   []string
	keyHeaders       []string
}

func buildKeyPolicy(behavior *cfstore.CacheBehavior, policy *cfstore.CachePolicy) keyPolicy {
	var kp keyPolicy
	// An attached cache policy fully determines keying and forwarding;
	// the deprecated ForwardedValues are mutually exclusive with it and
	// are ignored when a policy is present.
	if behavior.CachePolicyId != "" && policy != nil && policy.CachePolicyConfig != nil {
		if params := policy.CachePolicyConfig.ParametersInCacheKeyParametersInCacheKey; params != nil {
			if qs := params.QueryStringsConfig; qs != nil {
				switch qs.QueryStringBehavior {
				case "all":
					kp.forwardAllQuery = true
					kp.keyAllQuery = true
				case "whitelist":
					if qs.QueryStrings != nil {
						kp.keyQuery = qs.QueryStrings.Items
					}
				case "allExcept":
					if qs.QueryStrings != nil {
						kp.excludeQuery = qs.QueryStrings.Items
						kp.forwardAllQuery = true
					}
				}
			}
			if cc := params.CookiesConfig; cc != nil {
				switch cc.CookieBehavior {
				case "all":
					kp.cookiesForwarded = true
					kp.keyAllCookies = true
				case "whitelist":
					kp.cookiesForwarded = true
					if cc.Cookies != nil {
						kp.keyCookies = cc.Cookies.Items
					}
				case "allExcept":
					kp.cookiesForwarded = true
					if cc.Cookies != nil {
						kp.excludeCookies = cc.Cookies.Items
					}
				}
			}
			if hc := params.HeadersConfig; hc != nil && hc.HeaderBehavior != "none" {
				if hc.Headers != nil && hc.Headers.Quantity > 0 {
					kp.keyHeaders = hc.Headers.Items
				}
			}
		}
		return kp
	}
	if fv := behavior.ForwardedValues; fv != nil {
		if fv.QueryString {
			kp.forwardAllQuery = true
			kp.keyAllQuery = true
			if qk := fv.QueryStringCacheKeys; qk != nil && qk.Quantity > 0 {
				kp.keyAllQuery = false
				kp.keyQuery = qk.Items
			}
		}
		if ck := fv.Cookies; ck != nil && ck.Forward != "none" {
			kp.cookiesForwarded = true
			if ck.Forward == "all" {
				kp.keyAllCookies = true
			} else if wl := ck.WhitelistedNames; wl != nil {
				kp.keyCookies = wl.Items
			}
		}
		if hs := fv.Headers; hs != nil && hs.Quantity > 0 {
			kp.keyHeaders = hs.Items
		}
	}
	return kp
}

// selectedQuery returns the query parameters that participate in the
// cache key, in canonical (sorted) form.
func (kp keyPolicy) selectedQuery(query url.Values) url.Values {
	switch {
	case kp.keyAllQuery:
		return query
	case len(kp.keyQuery) > 0:
		selected := url.Values{}
		for _, name := range kp.keyQuery {
			if vals, ok := query[name]; ok {
				selected[name] = append([]string(nil), vals...)
			}
		}
		return selected
	case len(kp.excludeQuery) > 0:
		excluded := map[string]bool{}
		for _, name := range kp.excludeQuery {
			excluded[strings.ToLower(name)] = true
		}
		selected := url.Values{}
		for name, vals := range query {
			if !excluded[strings.ToLower(name)] {
				selected[name] = append([]string(nil), vals...)
			}
		}
		return selected
	}
	return nil
}

// originRawQuery builds the query string forwarded to the origin.
func (kp keyPolicy) originRawQuery(rawQuery string) string {
	if kp.forwardAllQuery {
		return rawQuery
	}
	if len(kp.keyQuery) == 0 && len(kp.excludeQuery) == 0 {
		return ""
	}
	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		return rawQuery
	}
	return kp.selectedQuery(query).Encode()
}

// buildCacheKey assembles the cache key from the distribution, cache
// behaviour, request path, and the key-policy components present in the
// request. The viewer protocol is deliberately absent: CloudFront caches
// an object once regardless of the protocol the fetch used.
func buildCacheKey(distributionID, pathPattern, requestPath string, r *http.Request, kp keyPolicy) string {
	var b strings.Builder
	b.WriteString(distributionID)
	b.WriteByte('\n')
	b.WriteString(pathPattern)
	b.WriteByte('\n')
	b.WriteString(requestPath)
	if q := kp.selectedQuery(r.URL.Query()); len(q) > 0 {
		b.WriteString("\nq:")
		b.WriteString(canonicalQuery(q))
	}
	if kp.keyAllCookies || len(kp.keyCookies) > 0 || len(kp.excludeCookies) > 0 {
		if cookies := kp.selectedCookies(r.Header.Get("Cookie")); len(cookies) > 0 {
			b.WriteString("\nc:")
			b.WriteString(cookies)
		}
	}
	for _, name := range sortedHeaderNames(kp.keyHeaders) {
		if values := r.Header.Values(name); len(values) > 0 {
			b.WriteString("\nh:")
			b.WriteString(strings.ToLower(name))
			b.WriteByte('=')
			b.WriteString(strings.Join(values, ","))
		}
	}
	return b.String()
}

func canonicalQuery(query url.Values) string {
	names := make([]string, 0, len(query))
	for name := range query {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		values := append([]string(nil), query[name]...)
		sort.Strings(values)
		parts = append(parts, name+"="+strings.Join(values, ","))
	}
	return strings.Join(parts, "&")
}

// selectedCookies returns the canonical "name=value; ..." list of the
// cookies that participate in the cache key.
func (kp keyPolicy) selectedCookies(cookieHeader string) string {
	if cookieHeader == "" {
		return ""
	}
	cookies := parseCookieHeader(cookieHeader)
	names := make([]string, 0, len(cookies))
	for name := range cookies {
		switch {
		case kp.keyAllCookies:
			names = append(names, name)
		case len(kp.keyCookies) > 0:
			if containsFold(kp.keyCookies, name) {
				names = append(names, name)
			}
		case len(kp.excludeCookies) > 0:
			if !containsFold(kp.excludeCookies, name) {
				names = append(names, name)
			}
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, name+"="+cookies[name])
	}
	return strings.Join(parts, "; ")
}

// cookieHeaderForOrigin filters the viewer Cookie header down to the
// cookies the behaviour forwards to the origin.
func (kp keyPolicy) cookieHeaderForOrigin(cookieHeader string) string {
	if !kp.cookiesForwarded || cookieHeader == "" {
		return ""
	}
	if kp.keyAllCookies {
		return cookieHeader
	}
	cookies := parseCookieHeader(cookieHeader)
	names := make([]string, 0, len(cookies))
	for name := range cookies {
		if len(kp.keyCookies) > 0 && containsFold(kp.keyCookies, name) {
			names = append(names, name)
		} else if len(kp.excludeCookies) > 0 && !containsFold(kp.excludeCookies, name) {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, name+"="+cookies[name])
	}
	return strings.Join(parts, "; ")
}

func parseCookieHeader(header string) map[string]string {
	cookies := map[string]string{}
	for _, part := range strings.Split(header, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if i := strings.Index(part, "="); i > 0 {
			cookies[strings.TrimSpace(part[:i])] = strings.TrimSpace(part[i+1:])
		} else {
			cookies[part] = ""
		}
	}
	return cookies
}

func containsFold(names []string, name string) bool {
	for _, candidate := range names {
		if strings.EqualFold(candidate, name) {
			return true
		}
	}
	return false
}

func sortedHeaderNames(names []string) []string {
	sorted := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, name := range names {
		lower := strings.ToLower(name)
		if !seen[lower] {
			seen[lower] = true
			sorted = append(sorted, name)
		}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i]) < strings.ToLower(sorted[j])
	})
	return sorted
}

// methodCacheable reports whether CloudFront caches responses to the
// method: GET and HEAD always, OPTIONS when the behaviour's cached
// methods include it.
func methodCacheable(method string, behavior *cfstore.CacheBehavior) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return true
	case http.MethodOptions:
		if behavior.AllowedMethods == nil {
			return false
		}
		for _, cached := range behavior.AllowedMethods.CachedMethods {
			if strings.EqualFold(cached, "OPTIONS") {
				return true
			}
		}
	}
	return false
}

// singleByteRange parses a "bytes=a-b" / "bytes=a-" / "bytes=-n" header
// against a body of the given length and reports the slice offsets. Only
// a single range is resolved; multi-range requests bypass the cache.
func singleByteRange(header string, size int) (start, length int, ok bool) {
	if !strings.HasPrefix(header, "bytes=") || size <= 0 {
		return 0, 0, false
	}
	spec := strings.TrimPrefix(header, "bytes=")
	if strings.Contains(spec, ",") {
		return 0, 0, false
	}
	dash := strings.Index(spec, "-")
	if dash < 0 {
		return 0, 0, false
	}
	first, last := strings.TrimSpace(spec[:dash]), strings.TrimSpace(spec[dash+1:])
	switch {
	case first == "" && last != "":
		suffix, err := strconv.Atoi(last)
		if err != nil || suffix <= 0 {
			return 0, 0, false
		}
		if suffix > size {
			suffix = size
		}
		return size - suffix, suffix, true
	case first != "" && last == "":
		begin, err := strconv.Atoi(first)
		if err != nil || begin < 0 || begin >= size {
			return 0, 0, false
		}
		return begin, size - begin, true
	case first != "" && last != "":
		begin, err1 := strconv.Atoi(first)
		end, err2 := strconv.Atoi(last)
		if err1 != nil || err2 != nil || begin < 0 || end < begin || begin >= size {
			return 0, 0, false
		}
		if end >= size {
			end = size - 1
		}
		return begin, end - begin + 1, true
	}
	return 0, 0, false
}

// hopByHopHeader names the response headers CloudFront strips before
// returning an origin response to the viewer.
func hopByHopHeader(name string) bool {
	switch strings.ToLower(name) {
	case "transfer-encoding", "trailer", "upgrade", "connection":
		return true
	}
	return false
}
