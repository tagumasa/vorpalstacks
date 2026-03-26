// Package policy provides IAM policy evaluation functionality for vorpalstacks.
package policy

import (
	"encoding/json"
)

// Effect represents whether a policy statement allows or denies access.
type Effect string

// Effect constants.
const (
	EffectAllow Effect = "Allow"
	EffectDeny  Effect = "Deny"
)

// Document represents an IAM policy document.
type Document struct {
	Version   string      `json:"Version"`
	Id        string      `json:"Id,omitempty"`
	Statement []Statement `json:"Statement"`
}

// Statement represents a single statement in an IAM policy document.
type Statement struct {
	Sid          string       `json:"Sid,omitempty"`
	Effect       Effect       `json:"Effect"`
	Principal    *Principal   `json:"Principal,omitempty"`
	NotPrincipal *Principal   `json:"NotPrincipal,omitempty"`
	Action       ActionList   `json:"Action,omitempty"`
	NotAction    ActionList   `json:"NotAction,omitempty"`
	Resource     ResourceList `json:"Resource,omitempty"`
	NotResource  ResourceList `json:"NotResource,omitempty"`
	Condition    ConditionMap `json:"Condition,omitempty"`
}

// Principal represents the principal element in a policy statement.
type Principal struct {
	AWS       StringList `json:"AWS,omitempty"`
	Service   StringList `json:"Service,omitempty"`
	Federated StringList `json:"Federated,omitempty"`
	Account   StringList `json:"Account,omitempty"`
	Everyone  bool       `json:"-"`
}

// ActionList represents a list of actions in a policy statement.
type ActionList []string

// ResourceList represents a list of resources in a policy statement.
type ResourceList []string

// StringList represents a list of strings in a policy statement.
type StringList []string

// ConditionMap represents the condition block in a policy statement.
type ConditionMap map[string]ConditionKeyValue

// ConditionKeyValue represents condition keys and their values.
type ConditionKeyValue map[string]ConditionValue

// ConditionValue represents condition values.
type ConditionValue []string

// ParseDocument parses a JSON policy document and returns a Document struct.
func ParseDocument(jsonStr string) (*Document, error) {
	var doc Document
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// UnmarshalJSON implements custom JSON unmarshalling for Principal.
func (p *Principal) UnmarshalJSON(data []byte) error {
	if string(data) == `"*"` {
		p.Everyone = true
		return nil
	}

	type alias Principal
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*p = Principal(tmp)
	return nil
}

// UnmarshalJSON implements custom JSON unmarshalling for ActionList.
// It handles both single string and array of strings formats.
func (al *ActionList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*al = ActionList{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}
	*al = ActionList(multi)
	return nil
}

// UnmarshalJSON implements custom JSON unmarshalling for ResourceList.
// It handles both single string and array of strings formats.
func (rl *ResourceList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*rl = ResourceList{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}
	*rl = ResourceList(multi)
	return nil
}

// Matches checks if the principal matches the given ARN.
func (p *Principal) Matches(principalArn string) bool {
	if p == nil {
		return true
	}
	if p.Everyone {
		return true
	}
	for _, arn := range p.AWS {
		if arn == "*" || arn == principalArn {
			return true
		}
		if matchArnPattern(arn, principalArn) {
			return true
		}
	}
	return false
}

// Matches checks if the action list contains the given action.
func (al ActionList) Matches(action string) bool {
	if len(al) == 0 {
		return true
	}
	for _, a := range al {
		if a == "*" || a == action {
			return true
		}
		if matchActionPattern(a, action) {
			return true
		}
	}
	return false
}

// Matches checks if the resource list contains the given resource.
func (rl ResourceList) Matches(resource string) bool {
	if len(rl) == 0 {
		return true
	}
	for _, r := range rl {
		if r == "*" || r == resource {
			return true
		}
		if matchArnPattern(r, resource) {
			return true
		}
	}
	return false
}

func matchActionPattern(pattern, action string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == "" {
		return false
	}
	if pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(action) >= len(prefix) && action[:len(prefix)] == prefix
	}
	return pattern == action
}

func matchArnPattern(pattern, arn string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == "" {
		return false
	}
	if pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(arn) >= len(prefix) && arn[:len(prefix)] == prefix
	}
	return pattern == arn
}
