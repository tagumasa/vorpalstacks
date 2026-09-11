// Package policy provides IAM policy evaluation functionality for vorpalstacks.
package policy

import (
	"strconv"
	"strings"
	"time"
)

// DecisionEffect represents the effect of an IAM policy decision.
type DecisionEffect string

// Decision effect constants.
const (
	DecisionEffectAllow       DecisionEffect = "Allow"
	DecisionEffectDeny        DecisionEffect = "Deny"
	DecisionEffectDefaultDeny DecisionEffect = "DefaultDeny"
)

// Decision represents the result of policy evaluation.
type Decision struct {
	Effect DecisionEffect
	Reason string
	// Matched holds the statements that determined the effect, each
	// identified by the index of its policy in the Evaluate argument
	// slice: when any matched statement denies, only the denying
	// statements appear (an explicit deny overrides every allow);
	// otherwise every matched allow statement appears. Callers that know
	// the identity of each evaluated document map the index to the
	// statement's source policy.
	Matched []MatchedStatement
}

// MatchedStatement identifies one statement that determined a Decision.
type MatchedStatement struct {
	// PolicyIndex is the position of the matched statement's document in
	// the policy slice passed to Evaluate.
	PolicyIndex int
	// Effect is the matched statement's own Effect (Allow or Deny).
	Effect DecisionEffect
}

// EvaluationContext holds the context information for policy evaluation.
type EvaluationContext struct {
	Principal              string
	PrincipalAccount       string
	Action                 string
	Resource               string
	RequestTime            time.Time
	SourceIP               string
	UserAgent              string
	UserID                 string
	UserName               string
	Referer                string
	SecureTransport        bool
	TokenIssueTime         time.Time
	MultiFactorAuthPresent bool
	// SessionContext maps a (lowercase) condition key to its full value
	// list: a single-valued key carries a one-element list, a multivalued
	// key carries every value as its own element. Values are verbatim —
	// a value that itself contains a comma stays one value, so set
	// operators must consume ContextValues, never split a joined string.
	SessionContext    map[string][]string
	EncryptionContext map[string]string
	ServiceContext    map[string]string
	Variables         map[string]string
}

// ResolveVariable resolves a variable in an IAM policy condition key.
func (ctx *EvaluationContext) ResolveVariable(key string) string {
	key = strings.TrimPrefix(key, "${")
	key = strings.TrimSuffix(key, "}")

	if value, ok := ctx.Variables[key]; ok {
		return value
	}

	if strings.HasPrefix(key, "aws:") {
		return ctx.resolveAWSVariable(key)
	}

	if strings.HasPrefix(key, "kms:") {
		return ctx.resolveKMSVariable(key)
	}

	return ctx.GetContextValue(key)
}

func (ctx *EvaluationContext) resolveAWSVariable(key string) string {
	switch key {
	case "aws:username":
		return ctx.UserName
	case "aws:userid":
		return ctx.UserID
	case "aws:principalarn", "aws:PrincipalArn":
		return ctx.Principal
	case "aws:principalaccount", "aws:PrincipalAccount":
		return ctx.PrincipalAccount
	case "aws:currenttime", "aws:CurrentTime":
		return ctx.RequestTime.Format(time.RFC3339)
	case "aws:epochtime", "aws:EpochTime":
		return strconv.FormatInt(ctx.RequestTime.Unix(), 10)
	case "aws:sourceip", "aws:SourceIp":
		return ctx.SourceIP
	case "aws:useragent", "aws:UserAgent":
		return ctx.UserAgent
	case "aws:RequestedRegion":
		if ctx.ServiceContext != nil {
			return ctx.ServiceContext["region"]
		}
	case "aws:RequestedAccount":
		return ctx.PrincipalAccount
	case "aws:referer", "aws:Referer":
		return ctx.Referer
	case "aws:securetransport", "aws:SecureTransport":
		if ctx.SecureTransport {
			return "true"
		}
		return "false"
	case "aws:tokenissuetime", "aws:TokenIssueTime":
		if !ctx.TokenIssueTime.IsZero() {
			return ctx.TokenIssueTime.Format(time.RFC3339)
		}
	case "aws:multifactorauthpresent", "aws:MultiFactorAuthPresent":
		if ctx.MultiFactorAuthPresent {
			return "true"
		}
		return "false"
	}
	return ""
}

func (ctx *EvaluationContext) resolveKMSVariable(key string) string {
	if !strings.HasPrefix(key, "kms:EncryptionContext:") {
		return ""
	}

	contextKey := strings.TrimPrefix(key, "kms:EncryptionContext:")
	if ctx.EncryptionContext != nil {
		return ctx.EncryptionContext[contextKey]
	}
	return ""
}

// GetContextValue retrieves a context value by key from the evaluation
// context as a single string: a multivalued key resolves to its values
// joined with commas, which serves scalar comparisons and the Null
// operator's presence check. Set operators (ForAnyValue:/ForAllValues:)
// must use ContextValues instead, so a value containing a comma is never
// mistaken for two values.
func (ctx *EvaluationContext) GetContextValue(key string) string {
	key = strings.ToLower(key)

	switch key {
	case "principal", "aws:principalarn":
		return ctx.Principal
	case "action", "aws:action":
		return ctx.Action
	case "resource", "aws:resource":
		return ctx.Resource
	case "sourceip", "aws:sourceip":
		return ctx.SourceIP
	case "currenttime", "aws:currenttime":
		return ctx.RequestTime.Format(time.RFC3339)
	case "useragent", "aws:useragent":
		return ctx.UserAgent
	}

	if ctx.SessionContext != nil {
		if val, ok := ctx.SessionContext[key]; ok {
			return strings.Join(val, ",")
		}
	}

	if ctx.EncryptionContext != nil {
		if val, ok := ctx.EncryptionContext[key]; ok {
			return val
		}
	}

	return ""
}

// ContextValues returns the full value list a condition key carries in the
// session context, verbatim and in supplied order; a key absent from the
// session context yields nil. The key is lowercased to match the
// session-context convention.
func (ctx *EvaluationContext) ContextValues(key string) []string {
	key = strings.ToLower(key)
	if ctx.SessionContext == nil {
		return nil
	}
	return ctx.SessionContext[key]
}

// PolicyEvaluator evaluates IAM policies against an evaluation context.
type PolicyEvaluator struct {
	conditionEvaluator *ConditionEvaluator
}

// NewPolicyEvaluator creates a new PolicyEvaluator instance.
func NewPolicyEvaluator() *PolicyEvaluator {
	return &PolicyEvaluator{
		conditionEvaluator: NewConditionEvaluator(),
	}
}

// Evaluate evaluates IAM policies against an evaluation context and returns a decision.
func (e *PolicyEvaluator) Evaluate(ctx *EvaluationContext, policies []*Document) *Decision {
	var matchedAllows, matchedDenies []MatchedStatement

	for i, policy := range policies {
		allows, denies := e.evaluatePolicy(ctx, policy, i)
		matchedAllows = append(matchedAllows, allows...)
		matchedDenies = append(matchedDenies, denies...)
	}

	if len(matchedDenies) > 0 {
		return &Decision{
			Effect:  DecisionEffectDeny,
			Reason:  "Explicit deny in policy",
			Matched: matchedDenies,
		}
	}

	if len(matchedAllows) > 0 {
		return &Decision{
			Effect:  DecisionEffectAllow,
			Reason:  "Allowed by policy",
			Matched: matchedAllows,
		}
	}

	return &Decision{
		Effect: DecisionEffectDefaultDeny,
		Reason: "No matching allow statement",
	}
}

// evaluatePolicy returns the statements of one document that matched the
// context, split by their own effect. It scans ALL statements so that an
// explicit Deny always wins over an Allow within the same policy,
// regardless of statement order; the caller keeps only the matched denies
// when any exist, which is where that override is applied across
// documents. policyIndex stamps each returned statement with the
// document's position in the Evaluate argument slice.
func (e *PolicyEvaluator) evaluatePolicy(ctx *EvaluationContext, policy *Document, policyIndex int) (allows, denies []MatchedStatement) {
	for _, stmt := range policy.Statement {
		if !e.statementMatches(ctx, &stmt) {
			continue
		}
		if stmt.Effect == EffectDeny {
			denies = append(denies, MatchedStatement{PolicyIndex: policyIndex, Effect: DecisionEffectDeny})
		} else {
			allows = append(allows, MatchedStatement{PolicyIndex: policyIndex, Effect: DecisionEffectAllow})
		}
	}
	return allows, denies
}

func (e *PolicyEvaluator) statementMatches(ctx *EvaluationContext, stmt *Statement) bool {
	if stmt.Principal != nil && !stmt.Principal.Matches(ctx.Principal) {
		return false
	}
	if stmt.NotPrincipal != nil && stmt.NotPrincipal.Matches(ctx.Principal) {
		return false
	}
	if len(stmt.Action) > 0 {
		if !stmt.Action.Matches(ctx.Action) {
			return false
		}
	} else if stmt.NotAction != nil {
		if stmt.NotAction.Matches(ctx.Action) {
			return false
		}
	}
	if len(stmt.Resource) > 0 {
		if !stmt.Resource.Matches(ctx.Resource) {
			return false
		}
	} else if stmt.NotResource != nil {
		if stmt.NotResource.Matches(ctx.Resource) {
			return false
		}
	}
	if !e.conditionEvaluator.Evaluate(stmt.Condition, ctx) {
		return false
	}

	return true
}

// EvaluateKeyPolicy evaluates a key policy document against an evaluation context.
func (e *PolicyEvaluator) EvaluateKeyPolicy(ctx *EvaluationContext, policyJSON string) *Decision {
	policy, err := ParseDocument(policyJSON)
	if err != nil {
		return &Decision{
			Effect: DecisionEffectDefaultDeny,
			Reason: "Invalid policy document: " + err.Error(),
		}
	}
	return e.Evaluate(ctx, []*Document{policy})
}
