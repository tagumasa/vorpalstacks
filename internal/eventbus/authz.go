// Package eventbus provides an internal service bus for cross-service
// communication within vorpalstacks, including IAM authorisation
// infrastructure (role resolution, resource policy evaluation).
package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// RoleResolver validates IAM role ARNs by checking role existence
// and trust policy presence.
type RoleResolver interface {
	ValidateRole(ctx context.Context, roleARN string) error
}

// BusPolicyEvaluator evaluates a BusPolicyDocument against a given
// principal, action, and resource to determine whether access is allowed.
type BusPolicyEvaluator interface {
	Evaluate(ctx context.Context, policy *BusPolicyDocument, principal string, action string, resource string) (bool, error)
}

// ResourcePolicyFunc retrieves the resource-based policy document for a
// target identified by its ARN. Returns nil (not an error) when no policy
// exists, which means the target is open to all principals.
type ResourcePolicyFunc func(ctx context.Context, targetARN string) (*BusPolicyDocument, error)

// BusPolicyDocument represents a normalised IAM resource-based policy
// document suitable for evaluation by BusPolicyEvaluator.
type BusPolicyDocument struct {
	Version   string               `json:"Version"`
	Statement []BusPolicyStatement `json:"Statement"`
}

// BusPolicyStatement represents a single statement within a
// BusPolicyDocument. Principal is stored as a normalised string
// extracted from either the string form ("*") or the object form
// ({"Service": "sns.amazonaws.com"}).
type BusPolicyStatement struct {
	Effect    string                       `json:"Effect"`
	Principal string                       `json:"Principal,omitempty"`
	Action    []string                     `json:"Action,omitempty"`
	Resource  []string                     `json:"Resource,omitempty"`
	Condition map[string]map[string]string `json:"Condition,omitempty"`
}

// EmptyBusPolicyDocument returns a policy document with the default
// version and no statements. A policy with no statements allows all
// access (open policy).
func EmptyBusPolicyDocument() *BusPolicyDocument {
	return &BusPolicyDocument{
		Version:   "2012-10-17",
		Statement: nil,
	}
}

// iamPolicyEvaluatorAdapter wraps an arbitrary evaluation function to
// satisfy the BusPolicyEvaluator interface. This avoids importing
// internal/common/iam/policy from the eventbus package.
type iamPolicyEvaluatorAdapter struct {
	evaluateFn func(ctx context.Context, policy *BusPolicyDocument, principal string, action string, resource string) (bool, error)
}

// NewIAMPolicyEvaluatorAdapter creates a BusPolicyEvaluator from the
// given evaluation function.
func NewIAMPolicyEvaluatorAdapter(evaluateFn func(ctx context.Context, policy *BusPolicyDocument, principal string, action string, resource string) (bool, error)) BusPolicyEvaluator {
	return &iamPolicyEvaluatorAdapter{evaluateFn: evaluateFn}
}

// Evaluate delegates to the wrapped evaluation function.
func (a *iamPolicyEvaluatorAdapter) Evaluate(ctx context.Context, policy *BusPolicyDocument, principal string, action string, resource string) (bool, error) {
	return a.evaluateFn(ctx, policy, principal, action, resource)
}

// ---------------------------------------------------------------------------
// Phase 4A: RoleResolver
// ---------------------------------------------------------------------------

// RoleLookupFunc retrieves a role's assume-role trust policy document by
// ARN. Returns the policy JSON string, or an error if the role does not
// exist or cannot be accessed.
type RoleLookupFunc func(ctx context.Context, roleARN string) (assumeRolePolicyDocument string, err error)

// IAMRoleResolver validates IAM role ARNs by checking role existence
// and trust policy presence via a RoleLookupFunc backed by the IAM store.
type IAMRoleResolver struct {
	lookup RoleLookupFunc
}

// NewIAMRoleResolver creates a new role resolver. The lookup function
// may be nil initially and set later via SetLookup; until set, all
// validations pass (backward compatible).
func NewIAMRoleResolver(lookup RoleLookupFunc) *IAMRoleResolver {
	return &IAMRoleResolver{lookup: lookup}
}

// SetLookup updates the role lookup function. This allows deferred wiring
// when the IAM store is not available at resolver creation time.
func (r *IAMRoleResolver) SetLookup(lookup RoleLookupFunc) {
	r.lookup = lookup
}

// ValidateRole checks that the specified role ARN refers to a valid role
// with a non-empty trust policy document. A nil lookup function causes
// all validations to pass.
func (r *IAMRoleResolver) ValidateRole(ctx context.Context, roleARN string) error {
	if r.lookup == nil {
		return nil
	}
	policy, err := r.lookup(ctx, roleARN)
	if err != nil {
		return fmt.Errorf("eventbus: role validation failed for %q: %w", roleARN, err)
	}
	if policy == "" {
		return fmt.Errorf("eventbus: role %q has no trust policy", roleARN)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Phase 4B: Policy document parsing and evaluation
// ---------------------------------------------------------------------------

// ParsePolicyDocument parses a JSON IAM policy document string into a
// BusPolicyDocument. Handles Principal as either a string (e.g. "*") or
// an object (e.g. {"Service": "sns.amazonaws.com"}). Action and Resource
// are normalised to string slices.
func ParsePolicyDocument(policyJSON string) (*BusPolicyDocument, error) {
	if policyJSON == "" {
		return EmptyBusPolicyDocument(), nil
	}

	var raw struct {
		Version   string          `json:"Version"`
		Statement json.RawMessage `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policyJSON), &raw); err != nil {
		return nil, fmt.Errorf("eventbus: failed to parse policy document: %w", err)
	}

	var rawStmts []json.RawMessage
	if err := json.Unmarshal(raw.Statement, &rawStmts); err != nil {
		rawStmts = []json.RawMessage{raw.Statement}
	}

	doc := &BusPolicyDocument{Version: raw.Version}
	for _, rs := range rawStmts {
		stmt, err := parsePolicyStatement(rs)
		if err != nil {
			return nil, err
		}
		doc.Statement = append(doc.Statement, *stmt)
	}
	return doc, nil
}

// parsePolicyStatement unmarshals a raw JSON object into a
// BusPolicyStatement, extracting Effect, Principal, Action, Resource,
// and Condition fields.
func parsePolicyStatement(raw json.RawMessage) (*BusPolicyStatement, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("eventbus: failed to parse policy statement: %w", err)
	}

	stmt := &BusPolicyStatement{}

	if v, ok := m["Effect"]; ok {
		_ = json.Unmarshal(v, &stmt.Effect)
	}
	stmt.Principal = extractPrincipalField(m["Principal"])
	stmt.Action = extractStringOrStringSlice(m["Action"])
	stmt.Resource = extractStringOrStringSlice(m["Resource"])

	if rawCond, ok := m["Condition"]; ok {
		stmt.Condition = parseConditionBlock(rawCond)
	}

	return stmt, nil
}

// extractPrincipalField normalises the IAM Principal field from raw JSON
// into a string. Handles string form ("*"), object form
// ({"Service": "..."}), and returns an empty string for nil input.
func extractPrincipalField(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}

	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}

	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) == nil {
		for _, key := range []string{"Service", "AWS", "Federated", "CanonicalUser"} {
			if v, ok := obj[key]; ok {
				var val string
				if json.Unmarshal(v, &val) == nil {
					return val
				}
			}
		}
	}
	return ""
}

// extractStringOrStringSlice unmarshals a raw JSON value that may be a
// single string or an array of strings, returning a consistent slice.
// Returns nil for nil input.
func extractStringOrStringSlice(raw json.RawMessage) []string {
	if raw == nil {
		return nil
	}

	var s string
	if json.Unmarshal(raw, &s) == nil {
		return []string{s}
	}

	var arr []string
	if json.Unmarshal(raw, &arr) == nil {
		return arr
	}
	return nil
}

// parseConditionBlock parses a raw JSON Condition block into a nested
// map[string]map[string]string. Only supports single-string values;
// complex condition values are silently skipped.
func parseConditionBlock(raw json.RawMessage) map[string]map[string]string {
	var operators map[string]json.RawMessage
	if json.Unmarshal(raw, &operators) != nil {
		return nil
	}
	result := make(map[string]map[string]string)
	for op, rawKV := range operators {
		var kv map[string]json.RawMessage
		if json.Unmarshal(rawKV, &kv) != nil {
			continue
		}
		inner := make(map[string]string)
		for k, v := range kv {
			var s string
			if json.Unmarshal(v, &s) == nil {
				inner[k] = s
			}
		}
		if len(inner) > 0 {
			result[op] = inner
		}
	}
	return result
}

// SimplePolicyEvaluator provides basic IAM policy statement matching
// with support for Effect, Principal, Action, and Resource fields.
// Condition evaluation is intentionally not supported; statements
// containing non-empty Condition blocks are treated as non-matching
// (fail-closed for security) since their intent cannot be verified.
type SimplePolicyEvaluator struct{}

// NewSimplePolicyEvaluator creates a new simple policy evaluator.
func NewSimplePolicyEvaluator() BusPolicyEvaluator {
	return &SimplePolicyEvaluator{}
}

// Evaluate checks whether the given principal is authorised to perform
// the specified action on the specified resource according to the policy
// document. Explicit Deny statements take precedence over Allow. When
// the policy exists but contains no matching statement, access is
// denied (default-deny). When the policy is nil or empty, access is
// allowed (no policy = open).
func (e *SimplePolicyEvaluator) Evaluate(ctx context.Context, policy *BusPolicyDocument, principal string, action string, resource string) (bool, error) {
	if policy == nil || len(policy.Statement) == 0 {
		return true, nil
	}

	hasAllow := false
	for _, stmt := range policy.Statement {
		if !policyStatementMatches(stmt, principal, action, resource) {
			continue
		}
		switch stmt.Effect {
		case "Deny":
			return false, nil
		case "Allow":
			hasAllow = true
		}
	}
	return hasAllow, nil
}

// policyStatementMatches checks whether a policy statement matches the
// given principal, action, and resource. Statements with non-empty
// Condition blocks are treated as non-matching (fail-closed).
func policyStatementMatches(stmt BusPolicyStatement, principal, action, resource string) bool {
	if len(stmt.Condition) > 0 {
		return false
	}
	if !principalMatches(stmt.Principal, principal) {
		return false
	}
	if !stringSliceMatches(stmt.Action, action) {
		return false
	}
	if !stringSliceMatches(stmt.Resource, resource) {
		return false
	}
	return true
}

// principalMatches checks whether the request principal matches the
// policy principal. An empty policy principal matches all (wildcard).
func principalMatches(policyPrincipal, requestPrincipal string) bool {
	if policyPrincipal == "" {
		return true
	}
	if policyPrincipal == "*" {
		return true
	}
	return policyPrincipal == requestPrincipal
}

// stringSliceMatches checks whether the given value matches any of the
// pattern strings. Supports exact match, wildcard "*", and suffix
// wildcard (e.g. "lambda:Invoke*"). An empty pattern slice matches all.
func stringSliceMatches(patterns []string, value string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if p == "*" || p == value {
			return true
		}
		if strings.HasSuffix(p, ":*") && strings.HasPrefix(value, p[:len(p)-1]) {
			return true
		}
		if strings.HasSuffix(p, "*") && strings.HasPrefix(value, p[:len(p)-1]) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Phase 4B: Resource policy factory functions
// ---------------------------------------------------------------------------

// LambdaPolicyLookupFunc retrieves a Lambda function's resource-based
// policy entries by function ARN. Returns the individual policy entries,
// or an error if the function does not exist.
type LambdaPolicyLookupFunc func(ctx context.Context, functionARN string) (entries []LambdaPolicyEntry, err error)

// LambdaPolicyEntry represents a single statement from a Lambda
// function's resource-based policy as stored in the Lambda function
// metadata.
type LambdaPolicyEntry struct {
	Statement string
	Principal string
	Action    string
	Resource  string
	Condition map[string]interface{}
}

// LambdaResourcePolicyFn creates a ResourcePolicyFunc that retrieves and
// converts a Lambda function's resource-based policy into a
// BusPolicyDocument. Each LambdaPolicyEntry becomes a statement with
// Effect "Allow".
func LambdaResourcePolicyFn(lookup LambdaPolicyLookupFunc) ResourcePolicyFunc {
	return func(ctx context.Context, functionARN string) (*BusPolicyDocument, error) {
		if lookup == nil {
			return EmptyBusPolicyDocument(), nil
		}
		entries, err := lookup(ctx, functionARN)
		if err != nil {
			return nil, fmt.Errorf("eventbus: failed to lookup lambda policy for %q: %w", functionARN, err)
		}
		if len(entries) == 0 {
			return EmptyBusPolicyDocument(), nil
		}
		doc := &BusPolicyDocument{Version: "2012-10-17"}
		for _, entry := range entries {
			stmt := BusPolicyStatement{
				Effect:    "Allow",
				Principal: entry.Principal,
				Action:    splitCommaSeparated(entry.Action),
				Resource:  splitCommaSeparated(entry.Resource),
			}
			doc.Statement = append(doc.Statement, stmt)
		}
		return doc, nil
	}
}

// SQSPolicyLookupFunc retrieves an SQS queue's resource-based policy
// document as raw JSON by queue ARN.
type SQSPolicyLookupFunc func(ctx context.Context, queueARN string) (policyJSON string, err error)

// SQSResourcePolicyFn creates a ResourcePolicyFunc that retrieves an SQS
// queue's policy by looking up the queue by name extracted from the ARN
// and reading its Policy attribute.
func SQSResourcePolicyFn(lookup SQSPolicyLookupFunc) ResourcePolicyFunc {
	return func(ctx context.Context, queueARN string) (*BusPolicyDocument, error) {
		if lookup == nil {
			return EmptyBusPolicyDocument(), nil
		}
		policyJSON, err := lookup(ctx, queueARN)
		if err != nil {
			return nil, fmt.Errorf("eventbus: failed to lookup SQS policy for %q: %w", queueARN, err)
		}
		return ParsePolicyDocument(policyJSON)
	}
}

// SNSTopicPolicyLookupFunc retrieves an SNS topic's resource-based policy
// document as raw JSON by topic ARN.
type SNSTopicPolicyLookupFunc func(ctx context.Context, topicARN string) (policyJSON string, err error)

// SNSTopicResourcePolicyFn creates a ResourcePolicyFunc that retrieves an
// SNS topic's policy by ARN and parses it into a BusPolicyDocument.
func SNSTopicResourcePolicyFn(lookup SNSTopicPolicyLookupFunc) ResourcePolicyFunc {
	return func(ctx context.Context, topicARN string) (*BusPolicyDocument, error) {
		if lookup == nil {
			return EmptyBusPolicyDocument(), nil
		}
		policyJSON, err := lookup(ctx, topicARN)
		if err != nil {
			return nil, fmt.Errorf("eventbus: failed to lookup SNS policy for %q: %w", topicARN, err)
		}
		return ParsePolicyDocument(policyJSON)
	}
}

// EventBridgeBusPolicyLookupFunc retrieves an EventBridge event bus's
// resource-based policy document as raw JSON by bus name.
type EventBridgeBusPolicyLookupFunc func(ctx context.Context, busName string) (policyJSON string, err error)

// EventBridgeBusResourcePolicyFn creates a ResourcePolicyFunc that
// retrieves an EventBridge event bus's policy by extracting the bus name
// from the ARN and parsing the policy JSON into a BusPolicyDocument.
func EventBridgeBusResourcePolicyFn(lookup EventBridgeBusPolicyLookupFunc) ResourcePolicyFunc {
	return func(ctx context.Context, busARN string) (*BusPolicyDocument, error) {
		if lookup == nil {
			return EmptyBusPolicyDocument(), nil
		}
		busName := svcarn.ExtractEventBusNameFromARN(busARN)
		if busName == "" {
			return EmptyBusPolicyDocument(), nil
		}
		policyJSON, err := lookup(ctx, busName)
		if err != nil {
			return nil, fmt.Errorf("eventbus: failed to lookup EventBridge policy for %q: %w", busARN, err)
		}
		return ParsePolicyDocument(policyJSON)
	}
}

// splitCommaSeparated splits a comma-separated action or resource string
// into individual entries, trimming whitespace from each.
func splitCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
