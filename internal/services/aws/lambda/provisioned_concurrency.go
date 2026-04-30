package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// PutProvisionedConcurrencyConfig configures provisioned concurrency for a Lambda function alias or version.
func (s *LambdaService) PutProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		return nil, NewInvalidParameter("Qualifier", "Qualifier is required")
	}

	concurrentExecutions := int32(request.GetIntParam(req.Parameters, "ProvisionedConcurrentExecutions"))
	if concurrentExecutions < 1 {
		return nil, NewInvalidParameter("ProvisionedConcurrentExecutions", "Provisioned concurrent executions must be at least 1")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.Functions.SetProvisionedConcurrency(functionName, qualifier, concurrentExecutions); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	config, err := store.Functions.GetProvisionedConcurrency(functionName, qualifier)
	if err != nil {
		return nil, err
	}

	return s.toProvisionedConcurrencyConfig(config), nil
}

// GetProvisionedConcurrencyConfig retrieves the provisioned concurrency configuration for a Lambda function alias or version.
func (s *LambdaService) GetProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		return nil, NewInvalidParameter("Qualifier", "Qualifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := store.Functions.GetProvisionedConcurrency(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrProvisionedConcurrencyNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return s.toProvisionedConcurrencyConfig(config), nil
}

// DeleteProvisionedConcurrencyConfig removes the provisioned concurrency configuration for a Lambda function alias or version.
func (s *LambdaService) DeleteProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		return nil, NewInvalidParameter("Qualifier", "Qualifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.Functions.DeleteProvisionedConcurrency(functionName, qualifier); err != nil {
		if err == lambdastore.ErrProvisionedConcurrencyNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListProvisionedConcurrencyConfigs lists the provisioned concurrency configurations for a Lambda function.
func (s *LambdaService) ListProvisionedConcurrencyConfigs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configs, err := store.Functions.ListProvisionedConcurrency(functionName)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		items = append(items, s.toProvisionedConcurrencyConfig(&c))
	}

	return map[string]interface{}{
		"ProvisionedConcurrencyConfigs": items,
	}, nil
}

func (s *LambdaService) toProvisionedConcurrencyConfig(c *lambdastore.ProvisionedConcurrencyConfig) map[string]interface{} {
	return map[string]interface{}{
		"FunctionName": c.FunctionName,
		"FunctionArn":  c.FunctionArn,
		"Qualifier":    c.Qualifier,
		"AllocatedProvisionedConcurrentExecutions": c.AllocatedProvisionedConcurrentExecutions,
		"AvailableProvisionedConcurrentExecutions": c.AvailableProvisionedConcurrentExecutions,
		"RequestedProvisionedConcurrentExecutions": c.RequestedProvisionedConcurrentExecutions,
		"Status":       c.Status,
		"StatusReason": c.StatusReason,
		"LastModified": c.LastModified.Format(timeutils.ISO8601UTCFormat),
	}
}
