// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"strings"

	"vorpalstacks/internal/services/aws/iam/policy"
)

// BuildEvaluationContext creates an evaluation context for trust policy evaluation.
func BuildEvaluationContext(sourceAccount string) *policy.EvaluationContext {
	return &policy.EvaluationContext{
		PrincipalAccount: sourceAccount,
		Variables: map[string]string{
			"aws:SourceAccount": sourceAccount,
		},
	}
}

// EvaluateTrustPolicy evaluates if a trust policy allows assuming a role.
func EvaluateTrustPolicy(doc *policy.Document, servicePrincipal ServicePrincipal, evalCtx *policy.EvaluationContext) error {
	for _, stmt := range doc.Statement {
		if stmt.Effect != policy.EffectAllow {
			continue
		}

		if !matchesPrincipal(stmt.Principal, string(servicePrincipal)) {
			continue
		}

		if !matchesAction(stmt.Action, "sts:AssumeRole") {
			continue
		}

		if stmt.Condition != nil {
			if !evaluateCondition(stmt.Condition, evalCtx) {
				continue
			}
		}

		return nil
	}
	return ErrRoleCannotBeAssumed
}

func evaluateCondition(conditions policy.ConditionMap, evalCtx *policy.EvaluationContext) bool {
	evaluator := policy.NewConditionEvaluator()
	return evaluator.Evaluate(conditions, evalCtx)
}

func matchConditionValue(operator string, value string, expected []string) bool {
	for _, e := range expected {
		if operator == "StringEquals" && value == e {
			return true
		}
		if operator == "StringLike" && matchPattern(e, value) {
			return true
		}
	}
	return false
}

func matchesPrincipal(principal *policy.Principal, servicePrincipal string) bool {
	if principal == nil {
		return false
	}
	if principal.Everyone {
		return true
	}
	for _, s := range principal.Service {
		if s == servicePrincipal || s == "*" {
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

func matchPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 0 {
			return true
		}
		if !strings.HasPrefix(value, parts[0]) {
			return false
		}
		if len(parts) > 1 && parts[len(parts)-1] != "" {
			return strings.HasSuffix(value, parts[len(parts)-1])
		}
		return true
	}
	return value == pattern
}
