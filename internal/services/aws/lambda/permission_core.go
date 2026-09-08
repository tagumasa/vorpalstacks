package lambda

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// AddPermissionInput carries the wire members of an AddPermission request.
// The condition map arrives already parsed from the scoping members.
type AddPermissionInput struct {
	StatementId         string
	Qualifier           string
	FunctionUrlAuthType string
	Principal           string
	Action              string
	Condition           map[string]interface{}
}

// qualifiedPolicyQualifier resolves a policy-scoping qualifier and
// reports whether it qualifies the policy resource: a published version
// or alias does, while $LATEST and the empty qualifier do not.
// AddPermission and the resource-policy operations resolve qualifiers
// identically, so they share this one path.
func qualifiedPolicyQualifier(stores *lambdaStore, functionName, qualifier string) (bool, error) {
	if qualifier == "" || qualifier == "$LATEST" {
		return false, nil
	}
	if _, _, _, err := stores.Functions.ResolveQualifier(functionName, qualifier); err != nil {
		return false, NewResourceNotFound("Qualifier", qualifier)
	}
	return true, nil
}

// addPermissionCore records a permission statement in a function's
// resource-based policy. It returns the qualified resource the statement
// was recorded for and the function's current RevisionId.
func (s *LambdaService) addPermissionCore(stores *lambdaStore, function *lambdastore.Function, in *AddPermissionInput) (resource, revisionId string, err error) {
	if in.StatementId == "" {
		return "", "", NewInvalidParameter("StatementId", "Statement ID is required")
	}
	if err := validateStatementId(in.StatementId); err != nil {
		return "", "", err
	}

	if in.FunctionUrlAuthType != "" {
		if err := validateAuthType(in.FunctionUrlAuthType); err != nil {
			return "", "", err
		}
	}

	// A qualifier scopes the permission to a published version or alias
	// and must resolve before the statement is recorded.
	resource = function.FunctionArn
	qualified, err := qualifiedPolicyQualifier(stores, function.FunctionName, in.Qualifier)
	if err != nil {
		return "", "", err
	}
	if qualified {
		resource = function.FunctionArn + ":" + in.Qualifier
	}

	policy := &lambdastore.FunctionPolicy{
		Id:        in.StatementId,
		Effect:    "Allow",
		Principal: in.Principal,
		Action:    in.Action,
		Resource:  resource,
		Condition: in.Condition,
	}

	if err := validatePermission(policy); err != nil {
		return "", "", err
	}

	if err := stores.Functions.AddPolicyAtomically(function.FunctionName, policy); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyAlreadyExists) {
			return "", "", NewResourceConflict(fmt.Sprintf("StatementId %s already exists", in.StatementId))
		}
		return "", "", mapStoreError(err)
	}

	// Include RevisionId so clients can use conditional updates.
	updated, getErr := stores.Functions.Get(function.FunctionName)
	if getErr != nil {
		logs.Warn("Failed to fetch function for RevisionId after AddPermission",
			logs.String("function", function.FunctionName),
			logs.Err(getErr))
	}
	if getErr == nil && updated.RevisionId != "" {
		revisionId = updated.RevisionId
	}

	return resource, revisionId, nil
}

// removePermissionCore removes a permission statement from a function's
// resource-based policy.
func (s *LambdaService) removePermissionCore(stores *lambdaStore, function *lambdastore.Function, statementId string) error {
	if statementId == "" {
		return NewInvalidParameter("StatementId", "Statement ID is required")
	}
	if err := stores.Functions.RemovePolicy(function.FunctionName, statementId); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyNotFound) {
			return NewResourceNotFound("Statement", statementId)
		}
		return err
	}
	return nil
}

// getPolicyCore retrieves a function's resource-based policy statements
// together with the function's current RevisionId. The function reference
// arrives in any namespaced form; the qualifier is irrelevant because the
// policy attaches to the whole function.
func (s *LambdaService) getPolicyCore(stores *lambdaStore, functionNameRaw string) (policies []lambdastore.FunctionPolicy, revisionId string, err error) {
	if functionNameRaw == "" {
		return nil, "", NewInvalidParameter("FunctionName", "Function name is required")
	}
	if err := validateNamespacedFunctionName(functionNameRaw); err != nil {
		return nil, "", err
	}
	functionName, _ := resolveNamespacedFunctionRef(functionNameRaw)

	policies, err = stores.Functions.GetPolicy(functionName)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	if len(policies) == 0 {
		return nil, "", ErrResourceNotFound
	}

	if fn, getErr := stores.Functions.Get(functionName); getErr == nil && fn.RevisionId != "" {
		revisionId = fn.RevisionId
	}

	return policies, revisionId, nil
}
