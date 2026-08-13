package lambda

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
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

	// Pre-warm the function container for the resolved qualifier.  This
	// eliminates cold-start latency on the first invocation.
	if s.dockerClient != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logs.Error("Panic in provisioned concurrency pre-warm goroutine",
						logs.String("function", functionName),
						logs.String("qualifier", qualifier),
						logs.Any("panic", r))
				}
			}()
			fn, err := store.Functions.Get(functionName)
			if err != nil {
				return
			}
			_, version, alias, err := store.Functions.ResolveQualifier(functionName, qualifier)
			if err != nil {
				return
			}
			ver := version
			if alias != nil {
				ver = resolveAliasTargetVersion(fn, alias)
			}
			_, _ = s.ensureFunctionContainer(fn, ver, store.Functions, reqCtx.GetRegion())
		}()
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

	maxItems := validateMaxItems(request.GetIntParam(req.Parameters, "MaxItems"))
	marker := request.GetStringParam(req.Parameters, "Marker")

	pageResult := pagination.PaginateSlice(configs, marker, maxItems, func(c lambdastore.ProvisionedConcurrencyConfig) string {
		return c.Qualifier
	})

	items := make([]map[string]interface{}, 0, len(pageResult.Items))
	for _, c := range pageResult.Items {
		items = append(items, s.toProvisionedConcurrencyConfig(&c))
	}

	response := map[string]interface{}{
		"ProvisionedConcurrencyConfigs": items,
	}
	if pageResult.IsTruncated {
		response["NextMarker"] = pageResult.NextMarker
	}

	return response, nil
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
