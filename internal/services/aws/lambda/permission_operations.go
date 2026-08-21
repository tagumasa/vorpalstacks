package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// buildPermissionStatement constructs the IAM statement map for a
// resource-based policy entry, including the Condition block derived from
// the AddPermission scoping members.
func buildPermissionStatement(sid, principal, action, resource string, condition map[string]interface{}) map[string]interface{} {
	// IAM ARN → {"AWS": ...}, service principal → {"Service": ...},
	// wildcard "*" → "*" (direct value, no wrapper key).
	statement := map[string]interface{}{
		"Sid":       sid,
		"Effect":    "Allow",
		"Principal": buildPrincipalField(principal),
		"Action":    action,
		"Resource":  resource,
	}
	if len(condition) > 0 {
		statement["Condition"] = condition
	}
	return statement
}

// parsePermissionCondition derives the IAM Condition block from the
// AddPermission scoping members. SourceArn uses the StringLike operator;
// the remaining members compare with StringEquals ("Lambda configures the
// [SourceArn] comparison using the StringLike operator" — API model).
func parsePermissionCondition(params map[string]interface{}) map[string]interface{} {
	condition := map[string]interface{}{}
	if sourceArn := request.GetStringParam(params, "SourceArn"); sourceArn != "" {
		condition["StringLike"] = map[string]interface{}{"aws:SourceArn": sourceArn}
	}
	stringEquals := map[string]interface{}{}
	if sourceAccount := request.GetStringParam(params, "SourceAccount"); sourceAccount != "" {
		stringEquals["aws:SourceAccount"] = sourceAccount
	}
	if principalOrgID := request.GetStringParam(params, "PrincipalOrgID"); principalOrgID != "" {
		stringEquals["aws:PrincipalOrgID"] = principalOrgID
	}
	if authType := request.GetStringParam(params, "FunctionUrlAuthType"); authType != "" {
		stringEquals["lambda:FunctionUrlAuthType"] = authType
	}
	if token := request.GetStringParam(params, "EventSourceToken"); token != "" {
		stringEquals["lambda:EventSourceToken"] = token
	}
	if len(stringEquals) > 0 {
		condition["StringEquals"] = stringEquals
	}
	return condition
}

// AddPermission adds a permission to a Lambda function's resource-based policy.
// Allows another AWS service or principal to invoke the function.
func (s *LambdaService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	statementId := request.GetStringParam(req.Parameters, "StatementId")
	if statementId == "" {
		return nil, NewInvalidParameter("StatementId", "Statement ID is required")
	}
	if err := validateStatementId(statementId); err != nil {
		return nil, err
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if functionUrlAuthType := request.GetStringParam(req.Parameters, "FunctionUrlAuthType"); functionUrlAuthType != "" {
		if err := validateAuthType(functionUrlAuthType); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// A qualifier scopes the permission to a published version or alias
	// and must resolve before the statement is recorded.
	resource := function.FunctionArn
	if qualifier != "" && qualifier != "$LATEST" {
		if _, _, _, err := store.Functions.ResolveQualifier(function.FunctionName, qualifier); err != nil {
			return nil, NewResourceNotFound("Qualifier", qualifier)
		}
		resource = function.FunctionArn + ":" + qualifier
	}

	principal := request.GetStringParam(req.Parameters, "Principal")
	action := request.GetStringParam(req.Parameters, "Action")

	condition := parsePermissionCondition(req.Parameters)

	policy := &lambdastore.FunctionPolicy{
		Id:        statementId,
		Principal: principal,
		Action:    action,
		Statement: request.GetStringParam(req.Parameters, "Statement"),
		Resource:  resource,
		Condition: condition,
	}

	if err := validatePermission(policy); err != nil {
		return nil, err
	}

	if err := store.Functions.AddPolicyAtomically(function.FunctionName, policy); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("StatementId %s already exists", statementId))
		}
		return nil, mapStoreError(err)
	}

	statement := buildPermissionStatement(statementId, principal, action, resource, condition)

	statementJSON, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"Statement": string(statementJSON),
	}

	// Include RevisionId so clients can use conditional updates.
	updated, getErr := store.Functions.Get(function.FunctionName)
	if getErr != nil {
		logs.Warn("Failed to fetch function for RevisionId after AddPermission",
			logs.String("function", function.FunctionName),
			logs.Err(getErr))
	}
	if getErr == nil && updated.RevisionId != "" {
		result["RevisionId"] = updated.RevisionId
	}

	return result, nil
}

// buildPrincipalField returns the IAM policy Principal value for the
// given principal string. Returns a wrapped map for ARN and service
// principals, or the bare string for the wildcard.
func buildPrincipalField(principal string) interface{} {
	pt := principalType(principal)
	if pt == "" {
		return principal
	}
	return map[string]interface{}{pt: principal}
}

// RemovePermission removes a permission from a Lambda function's resource-based policy.
func (s *LambdaService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	statementId := request.GetStringParam(req.Parameters, "StatementId")
	if statementId == "" {
		return nil, NewInvalidParameter("StatementId", "Statement ID is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Functions.RemovePolicy(function.FunctionName, statementId); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyNotFound) {
			return nil, NewResourceNotFound("Statement", statementId)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetPolicy returns the resource-based policy for a Lambda function.
func (s *LambdaService) GetPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	functionName = extractFunctionName(functionName)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policies, err := store.Functions.GetPolicy(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if len(policies) == 0 {
		return nil, ErrResourceNotFound
	}

	statements := make([]map[string]interface{}, 0)
	for _, p := range policies {
		statements = append(statements, buildPermissionStatement(p.Id, p.Principal, p.Action, p.Resource, p.Condition))
	}

	policyDoc := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statements,
	}

	policyJSON, err := json.Marshal(policyDoc)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"Policy": string(policyJSON),
	}
	if fn, getErr := store.Functions.Get(functionName); getErr == nil && fn.RevisionId != "" {
		result["RevisionId"] = fn.RevisionId
	}

	return result, nil
}
