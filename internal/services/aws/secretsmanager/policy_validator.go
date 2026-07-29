package secretsmanager

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/iam/policy"
)

// policyCheck mirrors the Smithy ValidationErrorsEntry shape (CheckName +
// ErrorMessage).  Multiple checks are returned by validatePolicyDocument so
// that ValidateResourcePolicy can report all issues at once, matching the
// AWS behaviour where PolicyValidationPassed=false and ValidationErrors
// contains one entry per failed check.
type policyCheck struct {
	CheckName    string
	ErrorMessage string
}

const (
	checkNameSyntax         = "RESOURCE_POLICY_SYNTAX"
	checkNameMissingVersion = "RESOURCE_POLICY_MISSING_VERSION"
)

// parsePolicyDocument parses a JSON resource-policy string into a
// policy.Document.  Returns a MalformedPolicyDocument-style sentinel so
// callers can map it to the correct AWS error type without inspecting the
// inner error text.
//
// Returns:
//   - *policy.Document on success
//   - error describing the parse failure (never nil when doc is nil)
func parsePolicyDocument(policyJSON string) (*policy.Document, error) {
	if policyJSON == "" {
		return nil, fmt.Errorf("resource policy is empty")
	}
	doc, err := policy.ParseDocument(policyJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid resource policy JSON: %w", err)
	}
	return doc, nil
}

// isPolicyPublic reports whether any Allow statement in the policy grants
// access to everyone (Principal "*" or {"AWS": "*"}) without a Condition
// block that would restrict the scope.  This mirrors the AWS Secrets Manager
// definition of a "public" resource policy used by BlockPublicPolicy and
// ValidateResourcePolicy.
func isPolicyPublic(doc *policy.Document) bool {
	if doc == nil {
		return false
	}
	for _, stmt := range doc.Statement {
		if stmt.Effect != policy.EffectAllow {
			continue
		}
		if !isPrincipalPublic(stmt.Principal) {
			continue
		}
		// A non-empty Condition block may scope down the grant (e.g.
		// StringEquals on aws:SourceAccount).  If present we do not
		// flag the statement as public.
		if hasCondition(stmt.Condition) {
			continue
		}
		return true
	}
	return false
}

// isPrincipalPublic checks whether the principal element grants access to
// everyone.  The principal is public when:
//   - Principal is the bare string "*" (parsed as Everyone=true)
//   - Principal.AWS contains "*"
func isPrincipalPublic(p *policy.Principal) bool {
	if p == nil {
		return false
	}
	if p.Everyone {
		return true
	}
	for _, arn := range p.AWS {
		if arn == "*" {
			return true
		}
	}
	return false
}

// hasCondition reports whether the condition map contains any entries.
func hasCondition(c policy.ConditionMap) bool {
	if c == nil {
		return false
	}
	for range c {
		return true
	}
	return false
}

// validatePolicyDocument runs structural validation checks on a
// resource-policy JSON string and returns a slice of policyCheck entries.
// The slice is empty when the policy passes all checks.
//
// ValidateResourcePolicy checks structural validity (syntax, missing
// version).  Public access detection is handled separately by
// PutResourcePolicy + BlockPublicPolicy, not by ValidateResourcePolicy.
//
// Checks performed:
//  1. RESOURCE_POLICY_SYNTAX — JSON must parse into a valid policy.
//  2. RESOURCE_POLICY_MISSING_VERSION — the policy should declare a
//     Version field (AWS recommends "2012-10-17").
func validatePolicyDocument(policyJSON string) []policyCheck {
	var checks []policyCheck

	doc, err := parsePolicyDocument(policyJSON)
	if err != nil {
		checks = append(checks, policyCheck{
			CheckName:    checkNameSyntax,
			ErrorMessage: err.Error(),
		})
		return checks
	}

	if doc.Version == "" {
		checks = append(checks, policyCheck{
			CheckName:    checkNameMissingVersion,
			ErrorMessage: "The policy is missing the Version field.",
		})
	}

	return checks
}

// policyChecksToResponse converts a slice of policyCheck into the
// ValidationErrors list format used in API responses.
func policyChecksToResponse(checks []policyCheck) []interface{} {
	result := make([]interface{}, 0, len(checks))
	for _, c := range checks {
		result = append(result, map[string]interface{}{
			"CheckName":    c.CheckName,
			"ErrorMessage": c.ErrorMessage,
		})
	}
	return result
}

// ensurePolicyJSONValid is a lightweight guard used by PutResourcePolicy
// to reject malformed JSON early with the correct error type.  It returns
// the parsed document so the caller can perform further checks (e.g.
// BlockPublicPolicy) without re-parsing.
func ensurePolicyJSONValid(policyJSON string) (*policy.Document, error) {
	doc, err := parsePolicyDocument(policyJSON)
	if err != nil {
		return nil, err
	}
	// Additional structural validation beyond ParseDocument.
	raw := make(map[string]json.RawMessage)
	if jErr := json.Unmarshal([]byte(policyJSON), &raw); jErr != nil {
		return nil, fmt.Errorf("resource policy is not a valid JSON object: %w", jErr)
	}
	return doc, nil
}
