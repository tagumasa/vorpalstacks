package policy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"vorpalstacks/internal/common/iotutil"
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
	// iot:Connect resources are client ARNs, not MQTT topics. They use
	// AWS policy wildcards (* and ?), not MQTT wildcards (+ and #).
	// iot:Publish/Subscribe/Receive resources are MQTT topic filters.
	resourceMatchFn := topicMatch
	if req.Action == "iot:Connect" {
		resourceMatchFn = wildcardMatch
	}
	for _, policy := range req.Policies {
		for _, stmt := range policy.Statement {
			actionMatch := len(stmt.Action) == 0 || matchAny(stmt.Action, req.Action, wildcardMatch)
			// AWS IoT policies allow variables in Resource patterns (e.g.
			// client/${iot:ClientId}); substitute before matching.
			resources := substituteResources(stmt.Resource, req)
			resourceMatch := len(resources) == 0 || matchAny(resources, req.Resource, resourceMatchFn)

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

// substituteResources applies AWS IoT policy variables (${iot:ClientId},
// ${iot:topic}, ${iot:SourceIp}) to each resource pattern so that wildcard
// matching operates on the concrete values.
func substituteResources(resources []string, req *EvaluateRequest) []string {
	if len(resources) == 0 {
		return resources
	}
	out := make([]string, len(resources))
	for i, r := range resources {
		r = strings.ReplaceAll(r, "${iot:ClientId}", req.ClientID)
		r = strings.ReplaceAll(r, "${iot:topic}", req.Topic)
		r = strings.ReplaceAll(r, "${iot:SourceIp}", req.SourceIP)
		out[i] = r
	}
	return out
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

// regexpCache memoises compiled wildcard patterns to avoid repeated
// regexp.MustCompile on every policy evaluation pass.
var regexpCache sync.Map

func wildcardMatch(s, pattern string) bool {
	return wildcardToRegexp(pattern).MatchString(s)
}

// wildcardToRegexp compiles an AWS wildcard pattern into an anchored regexp.
// * → .* (match any characters), all other characters are escaped.
// Results are cached in regexpCache for reuse across evaluations.
func wildcardToRegexp(pattern string) *regexp.Regexp {
	if cached, ok := regexpCache.Load(pattern); ok {
		return cached.(*regexp.Regexp)
	}

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
	re := regexp.MustCompile(buf.String())
	regexpCache.Store(pattern, re)
	return re
}

func topicMatch(resource, pattern string) bool {
	return iotutil.MQTTTopicMatch(resource, pattern)
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
