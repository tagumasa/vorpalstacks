// Package ssm provides Systems Manager Parameter Store storage functionality
// for vorpalstacks. The ParameterFilter type mirrors the Smithy
// ParameterStringFilter shape (Key, Option, Values) and is shared with the
// service-layer filter parser.
package ssm

import (
	"regexp"
	"strings"
)

// parameterStringFilterKeyPattern is the Smithy ParameterStringFilterKey
// constraint: tag-prefixed user keys plus a fixed set of system keys.
var parameterStringFilterKeyPattern = regexp.MustCompile(`^(tag:.+|Name|Type|KeyId|Path|Label|Tier|DataType)$`)

// validFilterOptions is the AWS-documented Option set for ParameterStringFilter.
// For Name, Contains is also supported. For Path, Recursive/OneLevel apply.
var validFilterOptions = map[string]struct{}{
	"Equals":     {},
	"BeginsWith": {},
	"Contains":   {},
	"Recursive":  {},
	"OneLevel":   {},
}

// ParameterFilter mirrors the Smithy ParameterStringFilter shape.
type ParameterFilter struct {
	Key    string
	Option string
	Values []string
}

// ValidateParameterFilterKey returns true when the given filter Key matches
// the Smithy ParameterStringFilterKey pattern.
func ValidateParameterFilterKey(key string) bool {
	return parameterStringFilterKeyPattern.MatchString(key)
}

// ValidateParameterFilterOption returns true when the given Option is one of
// the AWS-documented values.
func ValidateParameterFilterOption(option string) bool {
	_, ok := validFilterOptions[option]
	return ok
}

// Matches returns true when p matches every filter under the supplied context.
// The forGetByPath flag enforces the Smithy doc rule: GetParametersByPath
// excludes tag:, DataType, Name, Path and Tier keys. DescribeParameters
// excludes Label.
func (p *Parameter) Matches(filters []ParameterFilter, forGetByPath bool) bool {
	for _, f := range filters {
		if forGetByPath && !validKeyForGetByPath(f.Key) {
			continue
		}
		if !forGetByPath && f.Key == "Label" {
			continue
		}
		if !matchSingleFilter(p, f) {
			return false
		}
	}
	return true
}

func validKeyForGetByPath(key string) bool {
	if strings.HasPrefix(key, "tag:") {
		return false
	}
	switch key {
	case "DataType", "Name", "Path", "Tier":
		return false
	}
	return true
}

func matchSingleFilter(p *Parameter, f ParameterFilter) bool {
	// An omitted Values member is a key-existence filter: the member is
	// optional on the wire shape, so a filter carrying only a Key matches
	// every parameter that has the key's attribute (or tag key) at all.
	// Attribute keys every parameter carries (Name, Type, Tier, Path)
	// therefore match everything, while KeyId, DataType, Label and tag:
	// keys match only parameters where the attribute exists.
	if len(f.Values) == 0 {
		return keyExistsOnParameter(p, f.Key)
	}
	option := f.Option
	if option == "" {
		option = "Equals"
	}
	switch f.Key {
	case "Name":
		return matchStringOption(p.Name, f.Values, option)
	case "Type":
		return matchStringOption(string(p.Type), f.Values, option)
	case "KeyId":
		return matchStringOption(p.KeyID, f.Values, option)
	case "Tier":
		return matchStringOption(string(p.Tier), f.Values, option)
	case "DataType":
		return matchStringOption(p.DataType, f.Values, option)
	case "Path":
		return matchPathOption(p.Name, f.Values, option)
	case "Label":
		return matchLabelOption(p.VersionLabels, f.Values)
	}
	if strings.HasPrefix(f.Key, "tag:") {
		tagKey := strings.TrimPrefix(f.Key, "tag:")
		return matchStringOption(p.Tags[tagKey], f.Values, option)
	}
	return false
}

// keyExistsOnParameter reports whether the filter key's attribute exists on
// the parameter at all, for filters whose Values member was omitted.
func keyExistsOnParameter(p *Parameter, key string) bool {
	switch key {
	case "Name", "Type", "Tier", "Path":
		return true
	case "KeyId":
		return p.KeyID != ""
	case "DataType":
		return p.DataType != ""
	case "Label":
		return len(p.VersionLabels) > 0
	}
	if strings.HasPrefix(key, "tag:") {
		_, ok := p.Tags[strings.TrimPrefix(key, "tag:")]
		return ok
	}
	return false
}

func matchStringOption(field string, values []string, option string) bool {
	for _, v := range values {
		switch option {
		case "Equals":
			if field == v {
				return true
			}
		case "BeginsWith":
			if strings.HasPrefix(field, v) {
				return true
			}
		case "Contains":
			if strings.Contains(field, v) {
				return true
			}
		}
	}
	return false
}

func matchPathOption(name string, values []string, option string) bool {
	if option == "" {
		option = "OneLevel"
	}
	for _, v := range values {
		switch option {
		case "Equals", "BeginsWith":
			if strings.HasPrefix(name, v) {
				return true
			}
		case "Recursive":
			if strings.HasPrefix(name, v) && (strings.HasSuffix(v, "/") || (len(name) > len(v) && name[len(v)] == '/')) {
				return true
			}
		case "OneLevel":
			if strings.HasPrefix(name, v) {
				rest := name[len(v):]
				if rest == "" {
					return true
				}
				rest = strings.TrimPrefix(rest, "/")
				return !strings.Contains(rest, "/")
			}
		}
	}
	return false
}

func matchLabelOption(versionLabels map[int64][]string, values []string) bool {
	for _, vs := range versionLabels {
		for _, l := range vs {
			for _, v := range values {
				if l == v {
					return true
				}
			}
		}
	}
	return false
}
