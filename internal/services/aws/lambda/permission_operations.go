package lambda

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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
	principal := request.GetStringParam(req.Parameters, "Principal")
	action := request.GetStringParam(req.Parameters, "Action")
	condition := parsePermissionCondition(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resource, revisionId, err := s.addPermissionCore(store, function, &AddPermissionInput{
		StatementId:         statementId,
		Qualifier:           request.GetStringParam(req.Parameters, "Qualifier"),
		FunctionUrlAuthType: request.GetStringParam(req.Parameters, "FunctionUrlAuthType"),
		Principal:           principal,
		Action:              action,
		Statement:           request.GetStringParam(req.Parameters, "Statement"),
		Condition:           condition,
	})
	if err != nil {
		return nil, err
	}

	statement := buildPermissionStatement(statementId, principal, action, resource, condition)

	statementJSON, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"Statement": string(statementJSON),
	}
	if revisionId != "" {
		result["RevisionId"] = revisionId
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
	if err := s.removePermissionCore(store, function, statementId); err != nil {
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

	policies, revisionId, err := s.getPolicyCore(store, functionName)
	if err != nil {
		return nil, err
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
	if revisionId != "" {
		result["RevisionId"] = revisionId
	}

	return result, nil
}
