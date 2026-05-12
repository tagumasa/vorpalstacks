package ec2

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// ec2Filter represents a parsed EC2 Filter.N.Name + Filter.N.Value.M pair.
type ec2Filter struct {
	Name   string
	Values []string
}

// parseFilters extracts EC2 Filter.N.Name/Filter.N.Value.M from request params.
func parseFilters(params map[string]interface{}) []ec2Filter {
	var filters []ec2Filter
	for i := 1; ; i++ {
		name := request.GetStringParam(params, fmt.Sprintf("Filter.%d.Name", i))
		if name == "" {
			break
		}
		var values []string
		for j := 1; ; j++ {
			v := request.GetStringParam(params, fmt.Sprintf("Filter.%d.Value.%d", i, j))
			if v == "" {
				break
			}
			values = append(values, v)
		}
		if len(values) > 0 {
			filters = append(filters, ec2Filter{Name: name, Values: values})
		}
	}
	return filters
}

// anyMatch returns true if target matches any of the values (case-insensitive).
func anyMatch(values []string, target string) bool {
	for _, v := range values {
		if strings.EqualFold(v, target) {
			return true
		}
	}
	return false
}

// hasTagKey returns true if any tag has a key matching any of the given values.
func hasTagKey(tags []types.Tag, values []string) bool {
	for _, t := range tags {
		if anyMatch(values, t.Key) {
			return true
		}
	}
	return false
}

// hasTagValue returns true if any tag has a value matching any of the given values.
func hasTagValue(tags []types.Tag, values []string) bool {
	for _, t := range tags {
		if anyMatch(values, t.Value) {
			return true
		}
	}
	return false
}

// hasTagKeyValue returns true if any tag matches any "key:value" entry in values.
func hasTagKeyValue(tags []types.Tag, values []string) bool {
	for _, v := range values {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 {
			continue
		}
		for _, t := range tags {
			if t.Key == parts[0] && t.Value == parts[1] {
				return true
			}
		}
	}
	return false
}
