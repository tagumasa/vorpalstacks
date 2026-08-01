package s3

import (
	"net/http"
	"strconv"
	"strings"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// handleCORSPreflight processes OPTIONS requests for CORS preflight checks.
// It evaluates the request against the bucket's CORS configuration and returns
// the appropriate Access-Control-* headers. Returns true if CORS is configured
// and the request matched a rule, false otherwise.
func (h *S3Handler) handleCORSPreflight(w http.ResponseWriter, r *http.Request, bucket string, stores *s3Stores) bool {
	b, err := stores.buckets.Get(bucket)
	if err != nil || b.CORSConfiguration == nil || len(b.CORSConfiguration.CORSRules) == 0 {
		return false
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}

	requestMethod := r.Header.Get("Access-Control-Request-Method")
	requestHeaders := r.Header.Get("Access-Control-Request-Headers")

	rule := matchCORSRule(b.CORSConfiguration.CORSRules, origin, requestMethod, requestHeaders)
	if rule == nil {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.WriteHeader(http.StatusForbidden)
		return true
	}

	writeCORSHeaders(w, rule, origin)
	w.WriteHeader(http.StatusOK)
	return true
}

// applyCORSForActualRequest evaluates CORS for non-preflight requests.
// If the Origin header matches a CORS rule, Access-Control-* headers are
// added to the response. Returns true if CORS headers were applied.
func (h *S3Handler) applyCORSForActualRequest(w http.ResponseWriter, r *http.Request, bucket string, stores *s3Stores) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}

	b, err := stores.buckets.Get(bucket)
	if err != nil || b.CORSConfiguration == nil || len(b.CORSConfiguration.CORSRules) == 0 {
		return false
	}

	rule := matchCORSRule(b.CORSConfiguration.CORSRules, origin, r.Method, "")
	if rule == nil {
		return false
	}

	writeCORSHeaders(w, rule, origin)
	return true
}

// matchCORSRule finds the first CORS rule that matches the given origin, method,
// and headers. The method and headers parameters are matched against the rule's
// AllowedMethods and AllowedHeaders lists. An empty allowedHeadersStr matches
// any request that passes the origin and method checks.
func matchCORSRule(rules []s3store.CORSRule, origin, method, allowedHeadersStr string) *s3store.CORSRule {
	for i := range rules {
		rule := &rules[i]
		if !originMatches(rule.AllowedOrigins, origin) {
			continue
		}
		if method != "" && !methodMatches(rule.AllowedMethods, method) {
			continue
		}
		if allowedHeadersStr != "" && !headersMatch(rule.AllowedHeaders, allowedHeadersStr) {
			continue
		}
		return rule
	}
	return nil
}

// originMatches checks if the origin matches any of the allowed origins.
// Supports exact match and wildcard patterns (e.g., "*.example.com").
func originMatches(allowed []string, origin string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}
		if a == origin {
			return true
		}
		if strings.HasPrefix(a, "*.") {
			suffix := a[1:]
			if strings.HasSuffix(origin, suffix) {
				return true
			}
		}
	}
	return false
}

// methodMatches checks if the HTTP method is in the allowed methods list.
func methodMatches(allowed []string, method string) bool {
	for _, m := range allowed {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}

// headersMatch checks if all requested headers are in the allowed headers list.
// The allowed list may contain wildcards (e.g., "x-amz-*").
func headersMatch(allowed []string, requested string) bool {
	if requested == "" {
		return true
	}
	requestedHeaders := strings.Split(requested, ",")
	for _, rh := range requestedHeaders {
		rh = strings.TrimSpace(strings.ToLower(rh))
		matched := false
		for _, ah := range allowed {
			ah = strings.TrimSpace(strings.ToLower(ah))
			if ah == "*" {
				return true
			}
			if ah == rh {
				matched = true
				break
			}
			if strings.HasSuffix(ah, "*") && strings.HasPrefix(rh, ah[:len(ah)-1]) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// writeCORSHeaders sets the Access-Control-* response headers based on the matched rule.
func writeCORSHeaders(w http.ResponseWriter, rule *s3store.CORSRule, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Add("Vary", "Origin")

	if len(rule.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(rule.AllowedMethods, ", "))
	}
	if len(rule.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(rule.AllowedHeaders, ", "))
	}
	if len(rule.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(rule.ExposeHeaders, ", "))
	}
	if rule.MaxAgeSeconds != nil && *rule.MaxAgeSeconds > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(*rule.MaxAgeSeconds)))
	}
}
