package lambda

import (
	"encoding/json"
	"errors"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// PutResourcePolicyInput carries the wire members of a PutResourcePolicy
// request: the complete ARN of the resource, the policy document as a
// JSON string, and the optional policy-revision precondition.
type PutResourcePolicyInput struct {
	ResourceArn string
	Policy      string
	RevisionId  string
}

// DeleteResourcePolicyInput carries the wire members of a
// DeleteResourcePolicy request.
type DeleteResourcePolicyInput struct {
	ResourceArn string
	RevisionId  string
}

// ResourcePolicyResult carries a rendered policy document together with
// the policy revision that revision-preconditioned callers must present
// on the next mutation.
type ResourcePolicyResult struct {
	Policy     string
	RevisionId string
}

// parseResourcePolicyDocument parses a submitted policy document into
// policy statements. Each statement keeps its verbatim JSON in Raw so
// renderings reproduce it as submitted; the Principal, Action and
// Resource fields are populated for every Allow-effect statement with a
// string-shaped member because the eventbus dispatch authorisation
// matches on those fields. The IAM policy grammar allows Statement as an
// array or a single statement object; both are accepted here.
func parseResourcePolicyDocument(doc string) ([]lambdastore.FunctionPolicy, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(doc), &top); err != nil {
		return nil, NewInvalidParameter("Policy", "The policy document must be a valid JSON object")
	}
	rawStatements, ok := top["Statement"]
	if !ok {
		return nil, NewInvalidParameter("Policy", "The policy document must contain a Statement member")
	}

	var raws []json.RawMessage
	if err := json.Unmarshal(rawStatements, &raws); err != nil {
		var single map[string]json.RawMessage
		if err := json.Unmarshal(rawStatements, &single); err != nil {
			return nil, NewInvalidParameter("Policy", "Statement must be an array of statement objects or a single statement object")
		}
		raws = []json.RawMessage{rawStatements}
	}
	if len(raws) == 0 {
		return nil, NewInvalidParameter("Policy", "The policy document must contain at least one statement")
	}

	policies := make([]lambdastore.FunctionPolicy, 0, len(raws))
	for _, raw := range raws {
		var members map[string]json.RawMessage
		if err := json.Unmarshal(raw, &members); err != nil {
			return nil, NewInvalidParameter("Policy", "Every policy statement must be a JSON object")
		}

		entry := lambdastore.FunctionPolicy{Raw: string(raw)}
		if sidRaw, ok := members["Sid"]; ok {
			var sid string
			if err := json.Unmarshal(sidRaw, &sid); err != nil {
				return nil, NewInvalidParameter("Policy", "Statement Sid must be a string")
			}
			entry.Id = sid
		}

		effectRaw, ok := members["Effect"]
		if !ok {
			return nil, NewInvalidParameter("Policy", "Statement Effect is required")
		}
		var effect string
		if err := json.Unmarshal(effectRaw, &effect); err != nil {
			return nil, NewInvalidParameter("Policy", "Statement Effect must be a string")
		}
		entry.Effect = effect

		// The IAM policy grammar requires Principal and either Action or
		// NotAction in every resource-policy statement, whatever the
		// effect. Members that do not reduce to a single string (arrays,
		// NotAction) stay render-only: they cannot be evaluated, so they
		// grant nothing through the dispatch authorisation.
		if _, ok := members["Principal"]; !ok {
			return nil, NewInvalidParameter("Policy", "Statement Principal is required")
		}
		_, hasAction := members["Action"]
		_, hasNotAction := members["NotAction"]
		if !hasAction && !hasNotAction {
			return nil, NewInvalidParameter("Policy", "Statement Action or NotAction is required")
		}

		switch effect {
		case "Allow":
			entry.Principal = stringPrincipal(members["Principal"])
			entry.Action = jsonStringMember(members["Action"])
			entry.Resource = jsonStringMember(members["Resource"])
		case "Deny":
			// A Deny statement's matching fields stay empty: the
			// dispatch authorisation turns the entry into a
			// universal-deny statement instead of evaluating its
			// members.
		default:
			return nil, NewInvalidParameter("Policy", "Statement Effect must be Allow or Deny")
		}
		policies = append(policies, entry)
	}
	return policies, nil
}

// jsonStringMember decodes a member that must be a JSON string. Any other
// shape (absent, array, object, number) returns the empty string.
func jsonStringMember(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return ""
	}
	return s
}

// stringPrincipal reduces a statement Principal member to the single
// string the dispatch authorisation matches on: a bare string principal
// as-is, or the value of a single-entry {"AWS": "..."} or
// {"Service": "..."} object. Any other shape (multi-key object, arrays)
// returns the empty string — the statement stays renderable through its
// verbatim JSON and grants nothing through the matcher.
func stringPrincipal(raw json.RawMessage) string {
	if s := jsonStringMember(raw); s != "" {
		return s
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil || len(m) != 1 {
		return ""
	}
	for _, key := range []string{"AWS", "Service"} {
		if v, ok := m[key]; ok {
			return jsonStringMember(v)
		}
	}
	return ""
}

// renderResourcePolicyDocument renders a function's policy statements as
// a policy document. Statements that arrived in a PutResourcePolicy
// document are reproduced from their verbatim JSON; statements added via
// AddPermission are rendered from the structured fields.
func renderResourcePolicyDocument(policies []lambdastore.FunctionPolicy) (string, error) {
	statements := make([]interface{}, 0, len(policies))
	for _, p := range policies {
		if p.Raw != "" {
			statements = append(statements, json.RawMessage(p.Raw))
		} else {
			statements = append(statements, buildPermissionStatement(p.Id, p.Principal, p.Action, p.Resource, p.Condition))
		}
	}
	policyJSON, err := json.Marshal(map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statements,
	})
	if err != nil {
		return "", err
	}
	return string(policyJSON), nil
}

// resourcePolicyTarget validates the ResourceArn of a resource-policy
// operation and resolves it to the stored function. An embedded qualifier
// resolves through the same helper AddPermission uses.
func (s *LambdaService) resourcePolicyTarget(stores *lambdaStore, resourceArn string) (*lambdastore.Function, error) {
	functionName, qualifier, err := validateResourcePolicyArn(resourceArn)
	if err != nil {
		return nil, err
	}
	function, err := stores.Functions.Get(functionName)
	if err != nil {
		return nil, mapStoreError(err)
	}
	// A qualified resource ARN addresses a published version or alias;
	// the qualifier resolves through the shared helper.
	if _, err := qualifiedPolicyQualifier(stores, functionName, qualifier); err != nil {
		return nil, err
	}
	return function, nil
}

// putResourcePolicyCore replaces a resource's entire resource-based
// policy with the submitted document under the revision precondition.
func (s *LambdaService) putResourcePolicyCore(stores *lambdaStore, in *PutResourcePolicyInput) (*ResourcePolicyResult, error) {
	if in.Policy == "" {
		return nil, NewInvalidParameter("Policy", "Policy is required")
	}
	if len(in.Policy) > lambdastore.MaxResourcePolicyLength {
		return nil, NewPolicyLengthExceeded("The permissions policy is too large.")
	}
	entries, err := parseResourcePolicyDocument(in.Policy)
	if err != nil {
		return nil, err
	}
	function, err := s.resourcePolicyTarget(stores, in.ResourceArn)
	if err != nil {
		return nil, err
	}

	revision, err := stores.Functions.SetResourcePolicy(function.FunctionName, in.RevisionId, entries)
	if err != nil {
		if errors.Is(err, lambdastore.ErrPolicyRevisionMismatch) {
			return nil, NewPreconditionFailed(revisionMismatchMessage)
		}
		return nil, mapStoreError(err)
	}

	policyJSON, err := renderResourcePolicyDocument(entries)
	if err != nil {
		return nil, err
	}
	return &ResourcePolicyResult{Policy: policyJSON, RevisionId: revision}, nil
}

// getResourcePolicyCore renders the resource-based policy attached to a
// resource together with the policy revision. A resource without a policy
// answers ResourceNotFound, the same contract GetPolicy holds.
func (s *LambdaService) getResourcePolicyCore(stores *lambdaStore, resourceArn string) (*ResourcePolicyResult, error) {
	function, err := s.resourcePolicyTarget(stores, resourceArn)
	if err != nil {
		return nil, err
	}
	if len(function.Policies) == 0 {
		return nil, ErrResourceNotFound
	}
	policyJSON, err := renderResourcePolicyDocument(function.Policies)
	if err != nil {
		return nil, err
	}
	return &ResourcePolicyResult{Policy: policyJSON, RevisionId: function.PolicyRevisionId}, nil
}

// deleteResourcePolicyCore removes a resource's entire resource-based
// policy under the revision precondition. Deleting a policy-less resource
// succeeds so repeated deletes stay idempotent.
func (s *LambdaService) deleteResourcePolicyCore(stores *lambdaStore, in *DeleteResourcePolicyInput) error {
	function, err := s.resourcePolicyTarget(stores, in.ResourceArn)
	if err != nil {
		return err
	}
	if err := stores.Functions.DeleteResourcePolicy(function.FunctionName, in.RevisionId); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyRevisionMismatch) {
			return NewPreconditionFailed(revisionMismatchMessage)
		}
		return mapStoreError(err)
	}
	return nil
}
