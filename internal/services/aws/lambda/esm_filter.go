package lambda

import (
	"encoding/json"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/invokers"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// filterSQSRecords applies the event source mapping's FilterCriteria to a
// list of SQS records, returning only those that match at least one filter
// pattern.  If no criteria are configured, all records pass through
// unchanged.
func filterSQSRecords(records []esmSQSRecord, criteria *lambdastore.FilterCriteria) []esmSQSRecord {
	if criteria == nil || len(criteria.Filters) == 0 {
		return records
	}

	var result []esmSQSRecord
	for _, rec := range records {
		if matchesAnyFilter(sqsRecordToFilterable(rec), criteria) {
			result = append(result, rec)
		}
	}
	return result
}

// filterKinesisRecords applies FilterCriteria to Kinesis event records.
func filterKinesisRecords(records []map[string]interface{}, criteria *lambdastore.FilterCriteria) []map[string]interface{} {
	if criteria == nil || len(criteria.Filters) == 0 {
		return records
	}

	var result []map[string]interface{}
	for _, rec := range records {
		if matchesAnyFilter(rec, criteria) {
			result = append(result, rec)
		}
	}
	return result
}

// filterDynamoDBRecords applies FilterCriteria to DynamoDB Streams records.
func filterDynamoDBRecords(records []invokers.DynamoDBStreamRecord, criteria *lambdastore.FilterCriteria) []invokers.DynamoDBStreamRecord {
	if criteria == nil || len(criteria.Filters) == 0 {
		return records
	}

	var result []invokers.DynamoDBStreamRecord
	for _, rec := range records {
		// Convert struct to map for filter evaluation.
		m, err := structToMap(rec)
		if err != nil {
			result = append(result, rec)
			continue
		}
		if matchesAnyFilter(m, criteria) {
			result = append(result, rec)
		}
	}
	return result
}

// sqsRecordToFilterable converts an esmSQSRecord into a map suitable for
// filter pattern evaluation.  The message body is parsed as JSON when
// possible so that patterns can address individual fields within it.
func sqsRecordToFilterable(r esmSQSRecord) map[string]interface{} {
	m := map[string]interface{}{
		"messageId":      r.MessageID,
		"eventSource":    r.EventSource,
		"eventSourceArn": r.EventSourceARN,
		"awsRegion":      r.AWSRegion,
	}

	// Try to parse body as JSON; fall back to raw string.
	var parsed interface{}
	if json.Unmarshal([]byte(r.Body), &parsed) == nil {
		m["body"] = parsed
	} else {
		m["body"] = r.Body
	}

	if r.MessageAttributes != nil {
		m["messageAttributes"] = r.MessageAttributes
	}

	return m
}

// matchesAnyFilter returns true if the record matches at least one filter
// pattern (OR semantics across Filters).
func matchesAnyFilter(record map[string]interface{}, criteria *lambdastore.FilterCriteria) bool {
	for _, filter := range criteria.Filters {
		if filter.Pattern == "" {
			continue
		}
		var pattern map[string]interface{}
		if err := json.Unmarshal([]byte(filter.Pattern), &pattern); err != nil {
			continue // invalid pattern JSON — skip this filter
		}
		if matchesPattern(record, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern evaluates a single filter pattern against a record.
// All keys in the pattern must match (AND semantics).
func matchesPattern(record map[string]interface{}, pattern map[string]interface{}) bool {
	for key, patternVal := range pattern {
		recordVal, exists := record[key]
		if !exists {
			// Special handling for "exists": false criteria.
			if pm, ok := patternVal.([]interface{}); ok {
				for _, c := range pm {
					if cm, ok := c.(map[string]interface{}); ok {
						if e, ok := cm["exists"].(bool); ok && !e {
							return true // exists:false matches absent key
						}
					}
				}
			}
			return false
		}
		if !matchNode(recordVal, patternVal) {
			return false
		}
	}
	return true
}

// matchNode recursively evaluates a pattern node against a record value.
func matchNode(recordVal interface{}, patternVal interface{}) bool {
	switch pv := patternVal.(type) {
	case map[string]interface{}:
		// Nested object pattern — recurse into record value.
		rv, ok := recordVal.(map[string]interface{})
		if !ok {
			return false
		}
		return matchesPattern(rv, pv)
	case []interface{}:
		// Leaf criteria array — evaluate OR semantics.
		return matchesCriteria(recordVal, pv)
	default:
		return false
	}
}

// matchesCriteria evaluates a leaf criteria array.  Any single criterion
// matching is sufficient (OR within the criteria array).
func matchesCriteria(recordVal interface{}, criteria []interface{}) bool {
	for _, c := range criteria {
		if matchSingleCriterion(recordVal, c) {
			return true
		}
	}
	return false
}

// matchSingleCriterion evaluates one criterion value against the record value.
func matchSingleCriterion(recordVal interface{}, criterion interface{}) bool {
	switch c := criterion.(type) {
	case string:
		return fmt.Sprintf("%v", recordVal) == c
	case float64:
		return fmt.Sprintf("%v", recordVal) == fmt.Sprintf("%v", c)
	case bool:
		return fmt.Sprintf("%v", recordVal) == fmt.Sprintf("%v", c)
	case map[string]interface{}:
		for op, val := range c {
			result := matchOperator(recordVal, op, val)
			if result {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// matchOperator evaluates a single filter operator (prefix, numeric, exists,
// suffix, equals-ignore-case).
func matchOperator(recordVal interface{}, op string, val interface{}) bool {
	switch op {
	case "prefix":
		s, ok := recordVal.(string)
		if !ok {
			return false
		}
		prefix, ok := val.(string)
		if !ok {
			return false
		}
		return strings.HasPrefix(s, prefix)

	case "suffix":
		s, ok := recordVal.(string)
		if !ok {
			return false
		}
		suffix, ok := val.(string)
		if !ok {
			return false
		}
		return strings.HasSuffix(s, suffix)

	case "equals-ignore-case":
		s, ok := recordVal.(string)
		if !ok {
			return false
		}
		target, ok := val.(string)
		if !ok {
			return false
		}
		return strings.EqualFold(s, target)

	case "exists":
		// At this point, recordVal is the value found in the record.
		// exists:true → value is present (always true here).
		// exists:false → should never reach here (handled in matchesPattern).
		b, ok := val.(bool)
		if !ok {
			return false
		}
		return b

	case "numeric":
		return matchNumericOperator(recordVal, val)

	case "anything-but":
		// Match anything except the specified values.
		switch v := val.(type) {
		case string:
			return fmt.Sprintf("%v", recordVal) != v
		case []interface{}:
			recordStr := fmt.Sprintf("%v", recordVal)
			for _, item := range v {
				if fmt.Sprintf("%v", item) == recordStr {
					return false
				}
			}
			return true
		}
		return false

	default:
		return false
	}
}

// matchNumericOperator evaluates numeric range criteria.
// Format: [">", 0, "<=", 100] or [">", 50] etc.
func matchNumericOperator(recordVal interface{}, val interface{}) bool {
	parts, ok := val.([]interface{})
	if !ok || len(parts) < 2 {
		return false
	}

	// Convert recordVal to float64 for comparison.
	recordFloat, ok := toFloat64(recordVal)
	if !ok {
		return false
	}

	for i := 0; i+1 < len(parts); i += 2 {
		op, ok := parts[i].(string)
		if !ok {
			return false
		}
		threshold, ok := toFloat64(parts[i+1])
		if !ok {
			return false
		}
		switch op {
		case ">":
			if !(recordFloat > threshold) {
				return false
			}
		case ">=":
			if !(recordFloat >= threshold) {
				return false
			}
		case "<":
			if !(recordFloat < threshold) {
				return false
			}
		case "<=":
			if !(recordFloat <= threshold) {
				return false
			}
		case "=":
			if recordFloat != threshold {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// toFloat64 attempts to convert various numeric types to float64.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := parseFloatSafe(n)
		return f, err == nil
	default:
		return 0, false
	}
}

// structToMap converts a struct to a map[string]interface{} via JSON
// round-trip.  This is used for filter evaluation on typed records.
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func parseFloatSafe(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
