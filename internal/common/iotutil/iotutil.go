// Package iotutil provides shared utility functions used across IoT service
// packages (iot, iotevents, policy, rules) to eliminate duplication.
package iotutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var templatePattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ToInt converts a numeric value of any common Go type (float64, int, int64,
// float32, int32) to int. Returns 0 for unrecognised types or nil. This is
// used to handle iteration counters that may arrive as different numeric types
// depending on whether the message traversed JSON serialisation (float64) or
// stayed in-process (int).
func ToInt(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case float32:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case int32:
		return int(n)
	default:
		return 0
	}
}

// ResolveTemplate replaces AWS IoT substitution tokens (${topic()},
// ${payload.x}, ${timestamp()}, ${clientid()}) in the input string with
// runtime values from the MQTT message context. Unresolvable tokens are
// left as-is.
func ResolveTemplate(input, topic, clientID string, payload map[string]interface{}) string {
	return templatePattern.ReplaceAllStringFunc(input, func(match string) string {
		inner := match[2 : len(match)-1]
		return resolveToken(inner, topic, clientID, payload)
	})
}

func resolveToken(token, topic, clientID string, payload map[string]interface{}) string {
	if strings.HasPrefix(token, "topic(") {
		return resolveTopicFunction(token, topic)
	}
	switch token {
	case "topic()":
		return topic
	case "timestamp()":
		return fmt.Sprintf("%d", time.Now().Unix())
	case "timestamp":
		return fmt.Sprintf("%d", time.Now().Unix())
	case "clientid()":
		return clientID
	case "clientid":
		return clientID
	case "payload()":
		return fmt.Sprintf("%v", payload)
	default:
		if strings.HasPrefix(token, "payload.") {
			key := strings.TrimPrefix(token, "payload.")
			if val, ok := payload[key]; ok {
				return fmt.Sprintf("%v", val)
			}
		}
		if val, ok := payload[token]; ok {
			return fmt.Sprintf("%v", val)
		}
		return "${" + token + "}"
	}
}

func resolveTopicFunction(token, topic string) string {
	// Guard against malformed templates like ${topic(} where the closing
	// parenthesis is missing, which would cause token[6:5] to panic.
	if len(token) < 8 {
		return topic
	}
	argStr := token[6 : len(token)-1]
	if argStr == "" {
		return topic
	}
	parts := strings.Split(topic, "/")
	if strings.HasPrefix(argStr, "'") && strings.HasSuffix(argStr, "'") {
		return argStr[1 : len(argStr)-1]
	}
	idx, err := strconv.Atoi(argStr)
	if err != nil {
		return topic
	}
	if idx < 0 {
		idx = len(parts) + idx
	}
	if idx >= 0 && idx < len(parts) {
		return parts[idx]
	}
	return ""
}

// StrFromMap returns the first non-empty string value found for any of the
// given keys in the map. Returns empty string when none of the keys are
// present or all values are empty.
func StrFromMap(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// ListResponse builds a standard AWS list operation response with optional
// nextToken pagination marker.
// ListResponse builds a paginated list response map. IoT operations use two
// different marker conventions: marker-based ops (ListCertificates,
// ListPolicies, ListTopicRules) expect "nextMarker", while token-based ops
// (ListThings, ListThingGroups) expect "nextToken". Including both ensures
// the SDK deserialiser finds the continuation token regardless of convention.
func ListResponse(key string, items []map[string]interface{}, nextMarker string) map[string]interface{} {
	resp := map[string]interface{}{key: items}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
		resp["nextMarker"] = nextMarker
	}
	return resp
}

// MQTTTopicMatch reports whether topic matches the given MQTT topic filter
// pattern. Supports MQTT wildcards:
//   - "+" matches exactly one topic level
//   - "#" matches all remaining levels (must be the last level in the pattern)
func MQTTTopicMatch(topic, pattern string) bool {
	topicParts := splitSlash(topic)
	patternParts := splitSlash(pattern)
	return matchTopicParts(patternParts, topicParts, 0, 0)
}

func splitSlash(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func matchTopicParts(pattern, topic []string, pi, ti int) bool {
	for pi < len(pattern) && ti < len(topic) {
		// AWS IoT policies allow * as a wildcard equivalent to # when
		// used as a standalone pattern segment in topic resources.
		if pattern[pi] == "#" || pattern[pi] == "*" {
			return pi == len(pattern)-1
		}
		if pattern[pi] == "+" || pattern[pi] == topic[ti] {
			pi++
			ti++
			continue
		}
		return false
	}
	if pi < len(pattern) && pattern[pi] == "#" && pi == len(pattern)-1 {
		return true
	}
	return pi == len(pattern) && ti == len(topic)
}
