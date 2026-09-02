package inspection

import (
	"encoding/json"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// The WAF store keeps several rule fields as untyped interface values
// because they round-trip through JSON storage: after a load they hold
// map[string]interface{} shapes, after in-memory construction they hold
// typed pointers. normaliseThroughJSON converts such a value to its
// typed form with the same JSON round-trip the service layer uses for
// statement conversion.
func normaliseThroughJSON[T any](value interface{}) (T, bool) {
	var zero T
	if value == nil {
		return zero, false
	}
	if typed, ok := value.(T); ok {
		return typed, true
	}
	data, err := json.Marshal(value)
	if err != nil {
		return zero, false
	}
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return zero, false
	}
	return out, true
}

// ruleAction returns the typed action for a rule. A rule carries
// either Action (its own disposition) or, for rule group reference
// statements at the web ACL level, OverrideAction.
func ruleAction(rule *wafstore.Rule) *wafstore.Action {
	if action, ok := normaliseThroughJSON[*wafstore.Action](rule.Action); ok && action != nil {
		return action
	}
	return nil
}

func ruleOverrideAction(rule *wafstore.Rule) *wafstore.Action {
	if rule.OverrideAction != nil {
		return rule.OverrideAction
	}
	if action, ok := normaliseThroughJSON[*wafstore.Action](rule.OverrideAction); ok {
		return action
	}
	return nil
}

// webACLDefaultAction returns the web ACL's default action, which AWS
// restricts to Allow or Block.
func webACLDefaultAction(acl *wafstore.WebACL) *wafstore.Action {
	if action, ok := normaliseThroughJSON[*wafstore.Action](acl.DefaultAction); ok {
		return action
	}
	return nil
}

// ruleLabels returns the label declarations attached to a rule.
func ruleLabels(rule *wafstore.Rule) []RuleLabel {
	labels, ok := normaliseThroughJSON[[]RuleLabel](rule.RuleLabels)
	if !ok {
		return nil
	}
	return labels
}

// customResponseBodies returns the web ACL's custom response body map.
func customResponseBodies(acl *wafstore.WebACL) map[string]CustomResponseBody {
	bodies, ok := normaliseThroughJSON[map[string]CustomResponseBody](acl.CustomResponseBodies)
	if !ok {
		return nil
	}
	return bodies
}

// actionKind classifies an action value into one of the action name
// constants, treating the AWS-only Captcha, Challenge and Monetize
// actions as their literal names.
func actionKind(action *wafstore.Action) string {
	switch {
	case action == nil:
		return ""
	case action.Allow != nil:
		return ActionAllow
	case action.Block != nil:
		return ActionBlock
	case action.Count != nil:
		return ActionCount
	case action.Captcha != nil:
		return "Captcha"
	case action.Challenge != nil:
		return "Challenge"
	case action.Monetize != nil:
		return "Monetize"
	}
	return ""
}
