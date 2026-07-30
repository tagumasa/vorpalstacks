package ssm

import (
	"strconv"

	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// parseParameterFilters extracts ParameterFilters from a request in either
// shape: a JSON-decoded "ParameterFilters" key holding a list of maps (the
// AWS SDK JSON protocol), or flat query-string keys with the
// "ParameterFilters.N.Key" / "ParameterFilters.N.Option" /
// "ParameterFilters.N.Values.M" pattern (the AWS SDK query protocol).
//
// The deprecated "Filters" field uses the same query-string shape but with
// no Option member; this parser accepts both field names so SDK clients that
// still use the deprecated ParametersFilter keep working.
func parseParameterFilters(params map[string]interface{}) ([]ssmstore.ParameterFilter, error) {
	if filters, err := filtersFromList(params["ParameterFilters"]); filters != nil || err != nil {
		return filters, err
	}
	if filters, err := filtersFromList(params["Filters"]); filters != nil || err != nil {
		return filters, err
	}
	return filtersFromQueryParams(params, "ParameterFilters")
}

func filtersFromQueryParams(params map[string]interface{}, field string) ([]ssmstore.ParameterFilter, error) {
	var filters []ssmstore.ParameterFilter
	for i := 1; ; i++ {
		// AWS Query protocol uses ".member." between the field and the index.
		// SSM itself is awsJson1_1 and reaches this branch only when the
		// caller sends flat query-string keys; the .member. prefix restores
		// the canonical query protocol shape.
		key := paramString(params, field+".member."+strconv.Itoa(i)+".Key")
		if key == "" {
			return filters, nil
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, ErrInvalidFilterKey
		}
		option := paramString(params, field+".member."+strconv.Itoa(i)+".Option")
		if option != "" && !ssmstore.ValidateParameterFilterOption(option) {
			return nil, ErrInvalidFilterOption
		}
		var values []string
		for j := 1; ; j++ {
			v := paramString(params, field+".member."+strconv.Itoa(i)+".Values.member."+strconv.Itoa(j))
			if v == "" {
				break
			}
			values = append(values, v)
		}
		if len(values) == 0 {
			break
		}
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: option,
			Values: values,
		})
	}
	return filters, nil
}

// filtersFromList extracts filters from a JSON-decoded list under the given
// key. Returns (nil, nil) when the key is absent so callers can fall through
// to other formats; returns (nil, err) on validation failures.
func filtersFromList(raw interface{}) ([]ssmstore.ParameterFilter, error) {
	if raw == nil {
		return nil, nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, nil
	}
	var filters []ssmstore.ParameterFilter
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := m["Key"].(string)
		if key == "" {
			continue
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, ErrInvalidFilterKey
		}
		option, _ := m["Option"].(string)
		if option != "" && !ssmstore.ValidateParameterFilterOption(option) {
			return nil, ErrInvalidFilterOption
		}
		var values []string
		if vs, ok := m["Values"].([]interface{}); ok {
			for _, v := range vs {
				if s, ok := v.(string); ok && s != "" {
					values = append(values, s)
				}
			}
		}
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: option,
			Values: values,
		})
	}
	return filters, nil
}

// paramString returns the params map value as a string, handling both
// single-value and slice encodings that come out of url.ParseQuery.
func paramString(params map[string]interface{}, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		if len(val) == 0 {
			return ""
		}
		if s, ok := val[0].(string); ok {
			return s
		}
	}
	return ""
}
