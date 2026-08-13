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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principal := request.GetStringParam(req.Parameters, "Principal")
	action := request.GetStringParam(req.Parameters, "Action")

	policy := &lambdastore.FunctionPolicy{
		Id:        statementId,
		Principal: principal,
		Action:    action,
		Statement: request.GetStringParam(req.Parameters, "Statement"),
		Resource:  function.FunctionArn,
	}

	if err := validatePermission(policy); err != nil {
		return nil, err
	}

	if err := store.Functions.AddPolicyAtomically(function.FunctionName, policy); err != nil {
		if errors.Is(err, lambdastore.ErrPolicyAlreadyExists) {
			return nil, NewResourceConflict(fmt.Sprintf("StatementId %s already exists", statementId))
		}
		return nil, err
	}

	// Build the IAM policy statement with the correct Principal key.
	// IAM ARN → {"AWS": ...}, service principal → {"Service": ...},
	// wildcard "*" → "*" (direct value, no wrapper key).
	principalField := buildPrincipalField(principal)

	statement := map[string]interface{}{
		"Sid":       statementId,
		"Effect":    "Allow",
		"Principal": principalField,
		"Action":    action,
		"Resource":  function.FunctionArn,
	}

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
		stmt := map[string]interface{}{
			"Sid":       p.Id,
			"Effect":    "Allow",
			"Principal": buildPrincipalField(p.Principal),
			"Action":    p.Action,
		}
		if p.Resource != "" {
			stmt["Resource"] = p.Resource
		}
		if p.Condition != nil {
			stmt["Condition"] = p.Condition
		}
		statements = append(statements, stmt)
	}

	policyDoc := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statements,
	}

	policyJSON, err := json.Marshal(policyDoc)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy": string(policyJSON),
	}, nil
}
