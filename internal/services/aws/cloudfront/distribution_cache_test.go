package cloudfront

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	cfstore "vorpalstacks/internal/store/aws/cloudfront"
)

func TestResolveTTLOriginDirectives(t *testing.T) {
	hour := time.Hour
	cases := []struct {
		name   string
		header http.Header
		ttls   ttlSet
		want   time.Duration
	}{
		{
			name:   "max-age under zero bounds is used as-is",
			header: http.Header{"Cache-Control": {"max-age=3600"}},
			ttls:   ttlSet{},
			want:   hour,
		},
		{
			name:   "max-age above maximum is capped",
			header: http.Header{"Cache-Control": {"max-age=7200"}},
			ttls:   ttlSet{max: hour},
			want:   hour,
		},
		{
			name:   "max-age below minimum is raised",
			header: http.Header{"Cache-Control": {"max-age=600"}},
			ttls:   ttlSet{min: hour},
			want:   hour,
		},
		{
			name:   "max-age within bounds is kept",
			header: http.Header{"Cache-Control": {"max-age=1800"}},
			ttls:   ttlSet{min: hour / 4, max: hour},
			want:   30 * time.Minute,
		},
		{
			name:   "s-maxage wins over max-age",
			header: http.Header{"Cache-Control": {"max-age=600, s-maxage=1200"}},
			ttls:   ttlSet{},
			want:   20 * time.Minute,
		},
		{
			name:   "no directive uses default TTL",
			header: http.Header{},
			ttls:   ttlSet{def: hour},
			want:   hour,
		},
		{
			name:   "no directive with minimum raises to minimum when larger",
			header: http.Header{},
			ttls:   ttlSet{min: 2 * time.Hour, def: hour},
			want:   2 * time.Hour,
		},
		{
			name:   "no-store with zero minimum is not cached",
			header: http.Header{"Cache-Control": {"no-store"}},
			ttls:   ttlSet{def: hour},
			want:   0,
		},
		{
			name:   "no-store with positive minimum caches for the minimum",
			header: http.Header{"Cache-Control": {"private, no-cache"}},
			ttls:   ttlSet{min: time.Minute, def: hour},
			want:   time.Minute,
		},
		{
			name:   "max-age zero is not cached",
			header: http.Header{"Cache-Control": {"max-age=0"}},
			ttls:   ttlSet{},
			want:   0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveTTL(parseOriginDirectives(tc.header), tc.ttls)
			if got != tc.want {
				t.Fatalf("resolveTTL = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolveTTLExpires(t *testing.T) {
	future := time.Now().Add(30 * time.Minute).UTC().Format(http.TimeFormat)
	header := http.Header{"Expires": {future}}
	got := resolveTTL(parseOriginDirectives(header), ttlSet{})
	if got <= 29*time.Minute || got > 30*time.Minute {
		t.Fatalf("Expires resolution = %v, want just under 30 minutes", got)
	}

	past := time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)
	got = resolveTTL(parseOriginDirectives(http.Header{"Expires": {past}}), ttlSet{})
	if got != 0 {
		t.Fatalf("expired Expires resolution = %v, want 0", got)
	}

	// Cache-Control max-age takes precedence over Expires.
	both := http.Header{
		"Cache-Control": {"max-age=120"},
		"Expires":       {future},
	}
	if got := resolveTTL(parseOriginDirectives(both), ttlSet{}); got != 2*time.Minute {
		t.Fatalf("max-age precedence = %v, want 2 minutes", got)
	}
}

func TestErrorStatusCacheable(t *testing.T) {
	enabled := ttlSet{def: time.Hour}
	if !errorStatusCacheable(404, originDirectives{}, enabled) {
		t.Error("404 without directives must be cacheable")
	}
	if errorStatusCacheable(400, originDirectives{}, enabled) {
		t.Error("400 without directives must not be cacheable")
	}
	if !errorStatusCacheable(400, originDirectives{hasMaxAge: true, maxAge: time.Minute}, enabled) {
		t.Error("400 with max-age must be cacheable")
	}
	disabled := ttlSet{}
	if errorStatusCacheable(404, originDirectives{}, disabled) {
		t.Error("all-zero TTL configuration must cache nothing")
	}
	if errorStatusCacheable(418, originDirectives{}, enabled) {
		t.Error("418 is not a cacheable error status")
	}
}

func TestErrorCacheTTL(t *testing.T) {
	got := errorCacheTTL(nil, 404, originDirectives{})
	if got != cfstore.DefaultErrorCachingTTLSeconds*time.Second {
		t.Fatalf("default error TTL = %v", got)
	}
	custom := &cfstore.CustomErrorResponses{Items: []cfstore.CustomErrorResponse{
		{ErrorCode: 404, ErrorCachingMinTTL: 30},
	}}
	if got := errorCacheTTL(custom, 404, originDirectives{}); got != 30*time.Second {
		t.Fatalf("ErrorCachingMinTTL override = %v", got)
	}
	if got := errorCacheTTL(custom, 404, originDirectives{hasMaxAge: true, maxAge: time.Minute}); got != time.Minute {
		t.Fatalf("max-age floor = %v", got)
	}
}

func TestInvalidationPathMatches(t *testing.T) {
	cases := []struct {
		pattern, path string
		want          bool
	}{
		{"/*", "/anything/at/all", true},
		{"/foo", "/foo", true},
		{"/foo", "/foobar", false},
		{"/foo*", "/foobar", true},
		{"/foo/*", "/foo/bar", true},
		{"/foo/*", "/foobar", false},
		{"/a*b", "/axxb", true},
		{"/a*b", "/axx", false},
		{"", "/foo", false},
	}
	for _, tc := range cases {
		if got := invalidationPathMatches(tc.pattern, tc.path); got != tc.want {
			t.Errorf("invalidationPathMatches(%q, %q) = %v, want %v", tc.pattern, tc.path, got, tc.want)
		}
	}
}

func TestSingleByteRange(t *testing.T) {
	cases := []struct {
		header     string
		size       int
		wantOK     bool
		wantStart  int
		wantLength int
	}{
		{"bytes=0-4", 10, true, 0, 5},
		{"bytes=2-", 10, true, 2, 8},
		{"bytes=-3", 10, true, 7, 3},
		{"bytes=8-99", 10, true, 8, 2},
		{"bytes=0-1,5-6", 10, false, 0, 0},
		{"bytes=15-", 10, false, 0, 0},
		{"bytes=a-b", 10, false, 0, 0},
	}
	for _, tc := range cases {
		start, length, ok := singleByteRange(tc.header, tc.size)
		if ok != tc.wantOK || (ok && (start != tc.wantStart || length != tc.wantLength)) {
			t.Errorf("singleByteRange(%q) = (%d, %d, %v), want (%d, %d, %v)",
				tc.header, start, length, ok, tc.wantStart, tc.wantLength, tc.wantOK)
		}
	}
}

func TestResponseCacheLifecycle(t *testing.T) {
	cache := newResponseCache()
	entry := &cacheEntry{
		key:          "dist\n\n/x",
		distribution: "dist",
		path:         "/x",
		status:       200,
		body:         []byte("body"),
		expiresAt:    time.Now().Add(time.Minute),
	}
	cache.put(entry)
	if got := cache.lookup(entry.key); got == nil || got.status != 200 {
		t.Fatal("entry not found after put")
	}
	if got := cache.lookup("missing"); got != nil {
		t.Fatal("lookup of a missing key must return nil")
	}
	if got := cache.lookup(""); got != nil {
		t.Fatal("lookup of an empty key must return nil")
	}

	// Expiry does not remove the entry (stale serving needs it).
	entry.expiresAt = time.Now().Add(-time.Second)
	if got := cache.lookup(entry.key); got == nil {
		t.Fatal("expired entry must remain retrievable for stale serving")
	}

	cache.rearm(entry.key, time.Minute)
	if cache.lookup(entry.key).expired(time.Now()) {
		t.Fatal("rearm must extend the entry deadline")
	}

	if removed := cache.purge("dist", []string{"/*"}); removed != 1 {
		t.Fatalf("purge removed %d entries, want 1", removed)
	}
	if cache.lookup(entry.key) != nil {
		t.Fatal("entry must be gone after purge")
	}

	cache.put(entry)
	cache.purgeDistribution("dist")
	if cache.lookup(entry.key) != nil {
		t.Fatal("entry must be gone after purgeDistribution")
	}
	cache.put(entry)
	cache.purgeDistribution("other")
	if cache.lookup(entry.key) == nil {
		t.Fatal("purgeDistribution of another distribution must not remove entries")
	}
}

func TestResponseCacheLRUEviction(t *testing.T) {
	cache := newResponseCache()
	for i := 0; i <= maxCacheEntries; i++ {
		cache.put(&cacheEntry{
			key:          fmt.Sprintf("dist\n\n/%d", i),
			distribution: "dist",
			path:         fmt.Sprintf("/%d", i),
			expiresAt:    time.Now().Add(time.Hour),
		})
	}
	if cache.order.Len() != maxCacheEntries {
		t.Fatalf("cache holds %d entries, want the bound %d", cache.order.Len(), maxCacheEntries)
	}
	if cache.lookup("dist\n\n/0") != nil {
		t.Fatal("the oldest entry must have been evicted")
	}
	if cache.lookup("dist\n\n/1") == nil {
		t.Fatal("the second-oldest entry must survive")
	}
}

func TestBuildCacheKeyComponents(t *testing.T) {
	r := httptest.NewRequest("GET", "http://d.example.test/p?a=1&b=2&c=3", nil)
	r.Header.Set("Cookie", "session=abc; theme=dark")
	r.Header.Set("X-Variant", "one")

	legacyAll := keyPolicy{keyAllQuery: true, forwardAllQuery: true}
	legacyKeys := keyPolicy{forwardAllQuery: true, keyQuery: []string{"b"}}
	noQuery := keyPolicy{}
	cookies := keyPolicy{cookiesForwarded: true, keyAllCookies: true}
	cookieWhitelist := keyPolicy{cookiesForwarded: true, keyCookies: []string{"theme"}}
	headers := keyPolicy{keyHeaders: []string{"X-Variant"}}

	base := buildCacheKey("dist", "", "/p", r, noQuery)
	if key := buildCacheKey("dist", "", "/p", r, legacyAll); key == base {
		t.Error("full query keying must differ from the query-less key")
	}
	if key := buildCacheKey("dist", "", "/p", r, legacyKeys); !strings.Contains(key, "b=2") || strings.Contains(key, "a=1") {
		t.Errorf("query whitelist key = %q, want only b=2 keyed", key)
	}
	if key := buildCacheKey("dist", "", "/p", r, cookies); !strings.Contains(key, "session=abc") {
		t.Errorf("cookie key = %q, want session cookie included", key)
	}
	if key := buildCacheKey("dist", "", "/p", r, cookieWhitelist); strings.Contains(key, "session=") || !strings.Contains(key, "theme=dark") {
		t.Errorf("cookie whitelist key = %q, want only theme keyed", key)
	}
	if key := buildCacheKey("dist", "", "/p", r, headers); !strings.Contains(key, "x-variant=one") {
		t.Errorf("header key = %q, want the header value included", key)
	}
}

func TestParseCacheBehaviorDefaultTTL(t *testing.T) {
	cb := parseCacheBehavior(map[string]interface{}{
		"TargetOriginId": "o", "ViewerProtocolPolicy": "allow-all",
	})
	if cb.DefaultTTL != cfstore.LegacyDefaultTTLSeconds {
		t.Fatalf("omitted DefaultTTL = %d, want the 24-hour default", cb.DefaultTTL)
	}

	cb = parseCacheBehavior(map[string]interface{}{
		"TargetOriginId": "o", "CachePolicyId": "policy-1",
	})
	if cb.DefaultTTL != 0 {
		t.Fatalf("DefaultTTL with a cache policy = %d, want 0", cb.DefaultTTL)
	}

	cb = parseCacheBehavior(map[string]interface{}{
		"TargetOriginId": "o", "DefaultTTL": float64(0),
	})
	if cb.DefaultTTL != 0 {
		t.Fatalf("explicit zero DefaultTTL = %d, want 0", cb.DefaultTTL)
	}
}

// newCacheTestServer builds a DistributionServer wired to an httptest
// origin and a temporary store holding one enabled distribution.
func newCacheTestServer(t *testing.T, originHandler http.HandlerFunc, mutateConfig func(*cfstore.DistributionConfig)) (*DistributionServer, *cfstore.Distribution, *httptest.Server) {
	t.Helper()
	origin := httptest.NewServer(originHandler)
	t.Cleanup(origin.Close)

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })

	config := &cfstore.DistributionConfig{
		Enabled: true,
		Origins: cfstore.Origins{Quantity: 1, Items: []*cfstore.Origin{{
			ID:                 "origin-1",
			DomainName:         strings.TrimPrefix(origin.URL, "http://"),
			CustomOriginConfig: &cfstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
		}}},
		DefaultCacheBehavior: &cfstore.CacheBehavior{
			TargetOriginId:       "origin-1",
			ViewerProtocolPolicy: "allow-all",
			ForwardedValues:      &cfstore.ForwardedValues{QueryString: false, Cookies: &cfstore.CookiePreferences{Forward: "none"}},
		},
	}
	if mutateConfig != nil {
		mutateConfig(config)
	}

	distStore := cfstore.NewDistributionStore(st, "123456789012")
	dist, err := distStore.Create("cache-test", config)
	if err != nil {
		t.Fatal(err)
	}

	server := &DistributionServer{
		accountID: "123456789012",
		client:    &http.Client{Timeout: 5 * time.Second},
		cache:     newResponseCache(),
	}
	server.SetDistributionStore(distStore)
	return server, dist, origin
}

func cacheTestRequest(server *DistributionServer, dist *cfstore.Distribution, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://edge.example.test"+path, nil)
	r.Host = dist.ID + ".cloudfront.net"
	for name, value := range headers {
		r.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	server.HandleRequest(recorder, r)
	return recorder
}

func TestDistributionCacheHitAndMiss(t *testing.T) {
	originHits := 0
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		originHits++
		w.Header().Set("X-Origin-Hits", fmt.Sprint(originHits))
		_, _ = w.Write([]byte("payload"))
	}, func(config *cfstore.DistributionConfig) {
		// Explicit TTLs; the parse-time default only applies on the wire.
		config.DefaultCacheBehavior.MinTTL = 0
		config.DefaultCacheBehavior.DefaultTTL = 3600
		config.DefaultCacheBehavior.MaxTTL = 86400
	})

	first := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if got := first.Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("first response X-Cache = %q", got)
	}
	if first.Body.String() != "payload" || first.Code != http.StatusOK {
		t.Fatalf("first response = %d %q", first.Code, first.Body.String())
	}

	second := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if got := second.Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("second response X-Cache = %q, want Hit", got)
	}
	if second.Header().Get("X-Origin-Hits") != "1" {
		t.Fatalf("cached entry must not refetch the origin (origin hits header = %q)", second.Header().Get("X-Origin-Hits"))
	}
	if second.Body.String() != "payload" {
		t.Fatalf("cached body = %q", second.Body.String())
	}
	if second.Header().Get("Age") == "" {
		t.Fatal("hit responses must carry an Age header")
	}
	if originHits != 1 {
		t.Fatalf("origin hits = %d, want 1", originHits)
	}

	third := cacheTestRequest(server, dist, http.MethodGet, "/other", nil)
	if got := third.Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("a different path must miss, got %q", got)
	}

	// HEAD is served from the GET entry without a body.
	head := cacheTestRequest(server, dist, http.MethodHead, "/obj", nil)
	if got := head.Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("HEAD from a GET entry X-Cache = %q, want Hit", got)
	}
	if head.Body.Len() != 0 {
		t.Fatalf("HEAD response must carry no body, got %q", head.Body.String())
	}
}

func TestDistributionCacheControlHonoured(t *testing.T) {
	cases := []struct {
		name         string
		originHeader string
		minTTL       int
		wantSecond   string
	}{
		{"no-store is not cached", "Cache-Control: no-store", 0, "Miss from cloudfront"},
		{"max-age zero is not cached", "Cache-Control: max-age=0", 0, "Miss from cloudfront"},
		{"positive minimum overrides no-store", "Cache-Control: no-store", 60, "Hit from cloudfront"},
		{"origin max-age caches", "Cache-Control: max-age=600", 0, "Hit from cloudfront"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", strings.TrimPrefix(tc.originHeader, "Cache-Control: "))
				_, _ = w.Write([]byte("body"))
			}, func(config *cfstore.DistributionConfig) {
				config.DefaultCacheBehavior.MinTTL = tc.minTTL
				config.DefaultCacheBehavior.DefaultTTL = 3600
			})
			cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
			second := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
			if got := second.Header().Get("X-Cache"); got != tc.wantSecond {
				t.Fatalf("second X-Cache = %q, want %q", got, tc.wantSecond)
			}
		})
	}
}

func TestDistributionInvalidationPurge(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("body"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	cacheTestRequest(server, dist, http.MethodGet, "/keep", nil)
	cacheTestRequest(server, dist, http.MethodGet, "/drop", nil)

	server.InvalidatePaths(dist.ID, []string{"/drop"})

	if got := cacheTestRequest(server, dist, http.MethodGet, "/drop", nil).Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("invalidated path X-Cache = %q, want Miss", got)
	}
	if got := cacheTestRequest(server, dist, http.MethodGet, "/keep", nil).Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("untouched path X-Cache = %q, want Hit", got)
	}

	server.InvalidatePaths(dist.ID, []string{"/*"})
	if got := cacheTestRequest(server, dist, http.MethodGet, "/keep", nil).Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("after /* invalidation X-Cache = %q, want Miss", got)
	}
}

func TestDistributionSetCookiePolicy(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "origin=cookie")
		_, _ = w.Write([]byte("body"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	first := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if first.Header().Get("Set-Cookie") != "" {
		t.Fatal("Set-Cookie must be stripped when cookies are not forwarded")
	}
	second := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if second.Header().Get("Set-Cookie") != "" {
		t.Fatal("Set-Cookie must stay stripped on cache hits")
	}

	forwarding, forwardingDist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Set-Cookie", "origin=cookie")
		_, _ = w.Write([]byte("body"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
		config.DefaultCacheBehavior.ForwardedValues.Cookies.Forward = "all"
		config.DefaultCacheBehavior.ForwardedValues.Cookies.WhitelistedNames = nil
	})
	response := cacheTestRequest(forwarding, forwardingDist, http.MethodGet, "/obj", nil)
	if response.Header().Get("Set-Cookie") == "" {
		t.Fatal("Set-Cookie must pass through when cookies are forwarded")
	}
}

func TestDistributionQueryForwarding(t *testing.T) {
	var originQuery string
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		originQuery = r.URL.RawQuery
		_, _ = w.Write([]byte("body"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
		config.DefaultCacheBehavior.ForwardedValues.QueryString = false
	})

	cacheTestRequest(server, dist, http.MethodGet, "/obj?a=1", nil)
	if originQuery != "" {
		t.Fatalf("origin received query %q although forwarding is disabled", originQuery)
	}
	if got := cacheTestRequest(server, dist, http.MethodGet, "/obj?a=2", nil).Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatal("query values must not split the cache key when the query is not forwarded")
	}

	forwarding, forwardingDist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("body"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
		config.DefaultCacheBehavior.ForwardedValues.QueryString = true
	})
	cacheTestRequest(forwarding, forwardingDist, http.MethodGet, "/obj?a=1", nil)
	if got := cacheTestRequest(forwarding, forwardingDist, http.MethodGet, "/obj?a=1", nil).Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatal("identical queries must share the cache entry")
	}
	if got := cacheTestRequest(forwarding, forwardingDist, http.MethodGet, "/obj?a=2", nil).Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatal("different query values must split the cache key when the query is keyed")
	}
}

func TestDistributionStaleOnOriginFailure(t *testing.T) {
	originDown := false
	server, dist, origin := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if originDown {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("origin broken"))
			return
		}
		_, _ = w.Write([]byte("good object"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)

	// Force the entry into the expired state an elapsed TTL would leave
	// it in, then make the origin fail.
	server.cache.mu.Lock()
	for _, el := range server.cache.entries {
		el.Value.(*cacheEntry).expiresAt = time.Now().Add(-time.Second)
	}
	server.cache.mu.Unlock()
	originDown = true

	response := cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if response.Code != http.StatusOK || response.Body.String() != "good object" {
		t.Fatalf("stale serve = %d %q", response.Code, response.Body.String())
	}
	if got := response.Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("stale serve X-Cache = %q, want Hit", got)
	}

	// A fully unreachable origin serves stale content as well.
	server.cache.mu.Lock()
	for _, el := range server.cache.entries {
		el.Value.(*cacheEntry).expiresAt = time.Now().Add(-time.Second)
	}
	server.cache.mu.Unlock()
	origin.Close()
	response = cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	if response.Code != http.StatusOK || response.Body.String() != "good object" {
		t.Fatalf("unreachable-origin stale serve = %d %q", response.Code, response.Body.String())
	}
}

func TestDistributionRangeFromCache(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("0123456789"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	cacheTestRequest(server, dist, http.MethodGet, "/obj", nil)
	response := cacheTestRequest(server, dist, http.MethodGet, "/obj", map[string]string{"Range": "bytes=2-5"})
	if response.Code != http.StatusPartialContent {
		t.Fatalf("range hit status = %d, want 206", response.Code)
	}
	if response.Body.String() != "2345" {
		t.Fatalf("range hit body = %q", response.Body.String())
	}
	if got := response.Header().Get("Content-Range"); got != "bytes 2-5/10" {
		t.Fatalf("Content-Range = %q", got)
	}
	if got := response.Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("range hit X-Cache = %q", got)
	}
}

func TestDistributionErrorCaching(t *testing.T) {
	originHits := 0
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		originHits++
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	first := cacheTestRequest(server, dist, http.MethodGet, "/missing", nil)
	if first.Code != http.StatusNotFound {
		t.Fatalf("first response = %d", first.Code)
	}
	second := cacheTestRequest(server, dist, http.MethodGet, "/missing", nil)
	if second.Code != http.StatusNotFound {
		t.Fatalf("cached error = %d", second.Code)
	}
	if got := second.Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("cached 404 X-Cache = %q, want Hit", got)
	}
	if originHits != 1 {
		t.Fatalf("origin hits = %d, want 1", originHits)
	}
}

func TestDistributionNonCacheableMethodBypassesCache(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("posted"))
	}, func(config *cfstore.DistributionConfig) {
		config.DefaultCacheBehavior.AllowedMethods = &cfstore.AllowedMethods{
			Quantity:      7,
			Items:         []string{"GET", "HEAD", "OPTIONS", "PUT", "PATCH", "POST", "DELETE"},
			CachedMethods: []string{"GET", "HEAD"},
		}
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	post := func() *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "http://edge.example.test/obj", strings.NewReader("data"))
		r.Host = dist.ID + ".cloudfront.net"
		recorder := httptest.NewRecorder()
		server.HandleRequest(recorder, r)
		return recorder
	}
	if got := post().Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("first POST X-Cache = %q", got)
	}
	if got := post().Header().Get("X-Cache"); got != "Miss from cloudfront" {
		t.Fatalf("POST responses must never be served from cache, got %q", got)
	}
}

func TestDistributionAliasResolution(t *testing.T) {
	server, dist, _ := newCacheTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("alias body"))
	}, func(config *cfstore.DistributionConfig) {
		config.Aliases = &cfstore.Aliases{Quantity: 1, Items: []string{"www.example.com"}}
		config.DefaultCacheBehavior.DefaultTTL = 3600
	})

	aliasRequest := func(host string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "http://edge.example.test/alias", nil)
		r.Host = host
		recorder := httptest.NewRecorder()
		server.HandleRequest(recorder, r)
		return recorder
	}

	if response := aliasRequest("www.example.com"); response.Code != http.StatusOK || response.Body.String() != "alias body" {
		t.Fatalf("CNAME alias request = %d %q", response.Code, response.Body.String())
	}
	if response := aliasRequest("WWW.EXAMPLE.COM:50104"); response.Code != http.StatusOK {
		t.Fatalf("alias matching must be case-insensitive and port-tolerant, got %d", response.Code)
	}
	if response := aliasRequest("unknown.example.com"); response.Code != http.StatusNotFound {
		t.Fatalf("an unregistered host must 404, got %d", response.Code)
	}
	if response := aliasRequest("localhost"); response.Code != http.StatusNotFound {
		t.Fatalf("localhost must not resolve to a distribution, got %d", response.Code)
	}

	// The alias and the distribution domain share the same cache.
	if got := cacheTestRequest(server, dist, http.MethodGet, "/alias", nil).Header().Get("X-Cache"); got != "Hit from cloudfront" {
		t.Fatalf("the alias-fetched object must be shared with the distribution domain key, got %q", got)
	}
}
