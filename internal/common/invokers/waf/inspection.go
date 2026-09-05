package waf

import (
	"net/http"
	"strings"
)

// RequestHeadersWithHost returns the header map WAF inspection must
// see: a copy of the request's headers with the Host header restored.
// net/http promotes Host out of the header map into Request.Host, but
// the Host header is part of the inspected request on the wire —
// header-based statements and the HeaderOrder component both include
// it. The map is copied, so the caller's map is never modified.
func RequestHeadersWithHost(headers map[string][]string, host string) map[string][]string {
	out := make(map[string][]string, len(headers)+1)
	for name, values := range headers {
		out[name] = values
	}
	if host != "" {
		out["Host"] = []string{host}
	}
	return out
}

// BuildWebACLInspectionRequest assembles the inspection request from
// the parts of an incoming HTTP request. Cookie headers are parsed
// into the Cookies field (for the Cookies request component) and also
// kept in Headers, because header-based statements inspect the raw
// Cookie header. When the caller captured the wire order of the
// request's header names (headerorder.FromContext), the headers are
// emitted in that order so order-sensitive components such as WAF's
// HeaderOrder match the request as the client sent it; a nil order
// falls back to the header map's unspecified order.
func BuildWebACLInspectionRequest(method, uriPath, rawQuery, sourceIP, httpVersion string, headers map[string][]string, orderedHeaderNames []string, body []byte, bodyTruncated bool) WebACLInspectionRequest {
	req := WebACLInspectionRequest{
		Method:        method,
		URIPath:       uriPath,
		RawQuery:      rawQuery,
		SourceIP:      sourceIP,
		HTTPVersion:   httpVersion,
		Body:          body,
		BodyTruncated: bodyTruncated,
	}
	for _, header := range resolveHeaderNames(headers, orderedHeaderNames) {
		for _, value := range headers[header.key] {
			req.Headers = append(req.Headers, WebACLHTTPHeader{Name: header.wire, Value: value})
			if strings.EqualFold(header.wire, "Cookie") {
				for _, pair := range strings.Split(value, ";") {
					if cookieName, cookieValue, found := strings.Cut(strings.TrimSpace(pair), "="); found && cookieName != "" {
						req.Cookies = append(req.Cookies, WebACLHTTPHeader{Name: cookieName, Value: cookieValue})
					}
				}
			}
		}
	}
	return req
}

// resolvedHeader pairs a header's wire name (as the client sent it,
// which inspection must see) with its canonical map key (which the
// lookup needs).
type resolvedHeader struct {
	wire string
	key  string
}

// resolveHeaderNames resolves the header iteration order: the captured
// wire order when present, with any header the capture missed
// appended, and otherwise the map's unspecified order. Wire names are
// canonicalised for the map lookup but emitted as sent, so inspection
// sees the request's own casing.
func resolveHeaderNames(headers map[string][]string, captured []string) []resolvedHeader {
	if captured == nil {
		out := make([]resolvedHeader, 0, len(headers))
		for name := range headers {
			out = append(out, resolvedHeader{wire: name, key: name})
		}
		return out
	}
	seen := make(map[string]bool, len(headers))
	out := make([]resolvedHeader, 0, len(headers))
	for _, wireName := range captured {
		canonical := http.CanonicalHeaderKey(wireName)
		if _, ok := headers[canonical]; !ok || seen[canonical] {
			continue
		}
		seen[canonical] = true
		out = append(out, resolvedHeader{wire: wireName, key: canonical})
	}
	if len(out) == len(headers) {
		return out
	}
	for name := range headers {
		if !seen[name] {
			out = append(out, resolvedHeader{wire: name, key: name})
		}
	}
	return out
}
