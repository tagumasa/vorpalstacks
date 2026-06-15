package policy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Effect string

const (
	EffectAllow Effect = "Allow"
	EffectDeny  Effect = "Deny"
)

type PolicyVersion struct {
	Version   string            `json:"Version"`
	Statement []PolicyStatement `json:"Statement"`
}

type PolicyStatement struct {
	Effect    Effect                       `json:"Effect"`
	Action    StringOrSlice                `json:"Action,omitempty"`
	Resource  StringOrSlice                `json:"Resource,omitempty"`
	Condition map[string]map[string]string `json:"Condition,omitempty"`
	Principal map[string]string            `json:"Principal,omitempty"`
}

// StringOrSlice represents a JSON value that may be either a single string
// or an array of strings, matching the AWS IAM/IoT policy syntax where
// both "Action": "iot:*" and "Action": ["iot:Connect"] are valid.
type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}
	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return fmt.Errorf("expected string or array of strings: %w", err)
	}
	*s = multi
	return nil
}

type EvaluateRequest struct {
	Policies []*PolicyVersion
	Action   string
	Resource string
	ClientID string
	SourceIP string
	Topic    string
}

func Evaluate(req *EvaluateRequest) (bool, error) {
	allowed := false
	for _, policy := range req.Policies {
		for _, stmt := range policy.Statement {
			actionMatch := len(stmt.Action) == 0 || matchAny(stmt.Action, req.Action, wildcardMatch)
			resourceMatch := len(stmt.Resource) == 0 || matchAny(stmt.Resource, req.Resource, topicMatch)

			condMatch := true
			if stmt.Condition != nil {
				var err error
				condMatch, err = evaluateCondition(stmt.Condition, req)
				if err != nil {
					return false, fmt.Errorf("condition evaluation failed: %w", err)
				}
			}

			if actionMatch && resourceMatch && condMatch {
				if stmt.Effect == EffectDeny {
					return false, nil
				}
				allowed = true
			}
		}
	}
	return allowed, nil
}

func evaluateCondition(conditions map[string]map[string]string, req *EvaluateRequest) (bool, error) {
	for key, operators := range conditions {
		for op, value := range operators {
			var actual string
			switch key {
			case "iot:ClientId":
				actual = req.ClientID
			case "iot:SourceIp":
				actual = req.SourceIP
			case "iot:Topic":
				actual = req.Topic
			default:
				continue
			}

			fn, ok := conditionOperators[op]
			if !ok {
				return false, fmt.Errorf("unsupported condition operator: %s", op)
			}
			if !fn(actual, value) {
				return false, nil
			}
		}
	}
	return true, nil
}

func wildcardMatch(s, pattern string) bool {
	return wildcardToRegexp(pattern).MatchString(s)
}

// wildcardToRegexp compiles an AWS wildcard pattern into an anchored regexp.
// * → .* (match any characters), all other characters are escaped.
func wildcardToRegexp(pattern string) *regexp.Regexp {
	var buf strings.Builder
	buf.Grow(len(pattern) + 8)
	buf.WriteString("^")
	for _, ch := range pattern {
		if ch == '*' {
			buf.WriteString(".*")
		} else if strings.ContainsRune(`\.+?()[]{}^$|`, ch) {
			buf.WriteByte('\\')
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(ch)
		}
	}
	buf.WriteString("$")
	return regexp.MustCompile(buf.String())
}

func topicMatch(resource, pattern string) bool {
	return mqttTopicMatch(resource, pattern)
}

func mqttTopicMatch(topic, pattern string) bool {
	topicParts := strings.Split(topic, "/")
	patternParts := strings.Split(pattern, "/")

	for i := 0; i < len(patternParts); i++ {
		if patternParts[i] == "#" {
			return true
		}
		if i >= len(topicParts) {
			return false
		}
		if patternParts[i] != "+" && patternParts[i] != topicParts[i] {
			return false
		}
	}
	return len(topicParts) == len(patternParts)
}

func matchAny(patterns []string, value string, matchFn func(string, string) bool) bool {
	for _, p := range patterns {
		if matchFn(value, p) {
			return true
		}
	}
	return false
}

func ParsePolicyVersion(doc []byte) (*PolicyVersion, error) {
	var pv PolicyVersion
	if err := json.Unmarshal(doc, &pv); err != nil {
		return nil, fmt.Errorf("failed to parse policy document: %w", err)
	}
	if pv.Version == "" {
		return nil, fmt.Errorf("policy document must specify a Version")
	}
	for _, stmt := range pv.Statement {
		if stmt.Effect == "" {
			return nil, fmt.Errorf("policy statement must specify an Effect")
		}
	}
	return &pv, nil
}
