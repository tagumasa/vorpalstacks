package iot

import (
	"encoding/json"
	"strings"

	pagination "vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/iot/iotutil"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// parseAttributePayload extracts the attributes map from an attributePayload parameter.
func parseAttributePayload(params map[string]interface{}) (map[string]string, bool, bool) {
	attrs := make(map[string]string)

	// Try map form first (AWS SDK JSON protocol passes nested structs as maps).
	if m := request.GetMapParamCaseInsensitive(params, "attributePayload"); m != nil {
		if a, ok := m["attributes"].(map[string]interface{}); ok {
			for k, v := range a {
				if s, ok := v.(string); ok {
					attrs[k] = s
				}
			}
			merge := true
			if mb, ok := m["merge"].(bool); ok {
				merge = mb
			} else if ms, ok := m["merge"].(string); ok {
				merge = strings.ToLower(ms) == "true"
			}
			return attrs, merge, true
		}
	}

	attrStr := request.GetParamCaseInsensitive(params, "attributePayload")
	if attrStr == "" {
		return attrs, true, false
	}

	var payload struct {
		Attributes map[string]string `json:"attributes"`
		Merge      *bool             `json:"merge"`
	}
	if json.Unmarshal([]byte(attrStr), &payload) == nil && payload.Attributes != nil {
		merge := true
		if payload.Merge != nil {
			merge = *payload.Merge
		}
		return payload.Attributes, merge, true
	}

	var direct map[string]string
	if json.Unmarshal([]byte(attrStr), &direct) == nil {
		return direct, true, true
	}

	return attrs, true, false
}

// mapKeys returns the keys of a string map as an unordered slice.
// Map iteration order is non-deterministic in Go; callers that require
// stable ordering must sort the result themselves.
func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// parseListOptions extracts pagination parameters (nextToken, maxResults) from
// the request parameters into ListOptions.
// parseListOptions extracts pagination parameters from the request.
// IoT operations use two different parameter naming conventions:
//   - awsJson1_1 ops (ListCertificates, ListPolicies, ListTopicRules): marker / pageSize
//   - REST-JSON ops (ListThings, ListThingGroups): nextToken / maxResults
//
// Both conventions are supported here.
func parseListOptions(params map[string]interface{}) storecommon.ListOptions {
	opts := storecommon.ListOptions{}
	if token := request.GetParamCaseInsensitive(params, "nextToken"); token != "" {
		opts.Marker = token
	} else if token := request.GetParamCaseInsensitive(params, "marker"); token != "" {
		opts.Marker = token
	}
	if max := request.GetIntParam(params, "maxResults"); max > 0 {
		opts.MaxItems = max
	} else if max := request.GetIntParam(params, "pageSize"); max > 0 {
		opts.MaxItems = max
	}
	return opts
}

func ensureMap(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

// unwrapProps extracts a nested properties wrapper from request parameters.
// The AWS IoT awsJson1_1 protocol wraps operation parameters in structures
// (e.g., topicRulePayload, thingGroupProperties). The framework's JSON parser
// stores these as map[string]interface{} values at the top level of
// req.Parameters, so handlers must explicitly unwrap them.
// If wrapperKey is not found in params, the original params are returned
// unchanged.
func unwrapProps(params map[string]interface{}, wrapperKey string) map[string]interface{} {
	if m := request.GetMapParamCaseInsensitive(params, wrapperKey); m != nil {
		return m
	}
	return params
}

func listResponse(key string, items []map[string]interface{}, nextMarker string) map[string]interface{} {
	return iotutil.ListResponse(key, items, nextMarker)
}

func paginatedStrings(key string, items []string, params map[string]interface{}) map[string]interface{} {
	return pagination.PaginateOffsetStrings(key, items, params)
}

func paginatedMaps(key string, items []map[string]interface{}, params map[string]interface{}) map[string]interface{} {
	return pagination.PaginateOffsetMaps(key, items, params)
}

func boolToActiveStatus(active bool) string {
	if active {
		return "ACTIVE"
	}
	return "INACTIVE"
}

func principalFromParams(params map[string]interface{}) (string, bool) {
	if p := request.GetParamCaseInsensitive(params, "principal"); p != "" {
		return p, true
	}
	if t := request.GetParamCaseInsensitive(params, "target"); t != "" {
		return t, true
	}
	return "", false
}

// strVal extracts a string value from an interface{}, returning empty string for nil or non-string.
func strVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// float64Val extracts a float64 value from an interface{}, returning zero for nil or non-numeric.
func float64Val(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case json.Number:
		f, _ := n.Float64()
		return f
	}
	return 0
}

// int64Val extracts an int64 value from an interface{}, returning zero for nil or non-numeric.
func int64Val(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	}
	return 0
}

// boolVal extracts a boolean value from an interface{}, returning false for nil or non-bool.
func boolVal(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
