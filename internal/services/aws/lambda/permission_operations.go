package lambda

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

	policy := &lambdastore.FunctionPolicy{
		Id:        statementId,
		Principal: request.GetStringParam(req.Parameters, "Principal"),
		Action:    request.GetStringParam(req.Parameters, "Action"),
		Statement: request.GetStringParam(req.Parameters, "Statement"),
		Resource:  function.FunctionArn,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Functions.AddPolicyAtomically(function.FunctionName, policy); err != nil {
		return nil, err
	}

	statement := map[string]interface{}{
		"Sid":       statementId,
		"Effect":    "Allow",
		"Principal": map[string]interface{}{"Service": request.GetStringParam(req.Parameters, "Principal")},
		"Action":    request.GetStringParam(req.Parameters, "Action"),
		"Resource":  function.FunctionArn,
	}

	statementJSON, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Statement": string(statementJSON),
	}, nil
}

// RemovePermission removes a permission from a Lambda function's resource-based policy.
func (s *LambdaService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	statementId := request.GetStringParam(req.Parameters, "StatementId")
	if statementId == "" {
		return nil, NewInvalidParameter("StatementId", "Statement ID is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Functions.RemovePolicy(functionName, statementId); err != nil {
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
			"Principal": p.Principal,
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
