// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"vorpalstacks/internal/common/iam/policy"
)

// BuildEvaluationContext creates an evaluation context for trust policy evaluation.
func BuildEvaluationContext(sourceAccount string, principalArn string) *policy.EvaluationContext {
	return &policy.EvaluationContext{
		Principal:        principalArn,
		PrincipalAccount: sourceAccount,
		Variables: map[string]string{
			"aws:SourceAccount": sourceAccount,
		},
	}
}

// EvaluateTrustPolicy evaluates if a trust policy allows assuming a role.
func EvaluateTrustPolicy(doc *policy.Document, callerPrincipal string, evalCtx *policy.EvaluationContext) error {
	return EvaluateTrustPolicyForAction(doc, callerPrincipal, "sts:AssumeRole", evalCtx)
}

// EvaluateTrustPolicyForAction evaluates if a trust policy allows a specific
// STS action. An explicit Deny that matches the caller principal, action, and
// conditions always wins over an Allow, mirroring standard IAM evaluation.
func EvaluateTrustPolicyForAction(doc *policy.Document, callerPrincipal string, action string, evalCtx *policy.EvaluationContext) error {
	hasAllow := false

	for _, stmt := range doc.Statement {
		if !matchesPrincipal(stmt.Principal, callerPrincipal) {
			continue
		}

		if !matchesAction(stmt.Action, action) {
			continue
		}

		if stmt.Condition != nil {
			if !evaluateCondition(stmt.Condition, evalCtx) {
				continue
			}
		}

		if stmt.Effect == policy.EffectDeny {
			return ErrRoleCannotBeAssumed
		}

		if stmt.Effect == policy.EffectAllow {
			hasAllow = true
		}
	}

	if hasAllow {
		return nil
	}
	return ErrRoleCannotBeAssumed
}

func evaluateCondition(conditions policy.ConditionMap, evalCtx *policy.EvaluationContext) bool {
	evaluator := policy.NewConditionEvaluator()
	return evaluator.Evaluate(conditions, evalCtx)
}

func matchesPrincipal(principal *policy.Principal, callerPrincipal string) bool {
	if principal == nil {
		return false
	}
	if principal.Everyone {
		return true
	}
	if principal.Matches(callerPrincipal) {
		return true
	}
	for _, s := range principal.Service {
		if s == callerPrincipal || s == "*" {
			return true
		}
	}
	for _, f := range principal.Federated {
		if f == callerPrincipal || f == "*" {
			return true
		}
	}
	return false
}

func matchesAction(actions policy.ActionList, action string) bool {
	for _, a := range actions {
		if a == action || a == "*" {
			return true
		}
	}
	return false
}
