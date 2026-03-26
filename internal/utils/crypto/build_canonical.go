package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// AWS signature version 4 canonical request building.

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// BuildCanonicalRequest builds a canonical request string for AWS Signature Version 4.
func BuildCanonicalRequest(r *http.Request, signedHeaders string, body []byte) string {
	method := r.Method

	uri := r.URL.EscapedPath()
	if uri == "" {
		uri = "/"
	}

	query := buildCanonicalQueryString(r.URL.Query())

	canonicalHeaders := buildCanonicalHeaders(r, signedHeaders)

	signedHeadersList := signedHeaders

	payloadHash := SHA256HashString(string(body))

	return method + "\n" +
		uri + "\n" +
		query + "\n" +
		canonicalHeaders + "\n" +
		signedHeadersList + "\n" +
		payloadHash
}

func buildCanonicalQueryString(query url.Values) string {
	if len(query) == 0 {
		return ""
	}

	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var result strings.Builder
	for i, key := range keys {
		if i > 0 {
			result.WriteString("&")
		}
		result.WriteString(url.QueryEscape(key))
		result.WriteString("=")

		values := query[key]
		for j, value := range values {
			if j > 0 {
				result.WriteString("&")
				result.WriteString(url.QueryEscape(key))
				result.WriteString("=")
			}
			result.WriteString(url.QueryEscape(value))
		}
	}

	return result.String()
}

func buildCanonicalHeaders(r *http.Request, signedHeaders string) string {
	headers := make(map[string]string)
	for _, header := range strings.Split(signedHeaders, ";") {
		header = strings.TrimSpace(header)
		if header != "" {
			value := r.Header.Get(header)
			headers[strings.ToLower(header)] = strings.TrimSpace(value)
		}
	}

	sortedHeaders := make([]string, 0, len(headers))
	for header := range headers {
		sortedHeaders = append(sortedHeaders, header)
	}
	sort.Strings(sortedHeaders)

	var result strings.Builder
	for _, header := range sortedHeaders {
		result.WriteString(header)
		result.WriteString(":")
		result.WriteString(headers[header])
		result.WriteString("\n")
	}

	return result.String()
}
