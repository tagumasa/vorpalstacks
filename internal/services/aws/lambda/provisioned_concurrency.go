package lambda

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// resolveProvisionedConcurrencyTarget resolves the FunctionName reference
// forms and merges the Qualifier parameter with an embedded qualifier.
// Provisioned concurrency targets a published version or alias, so a
// qualifier is mandatory.
func resolveProvisionedConcurrencyTarget(params map[string]interface{}) (string, string, error) {
	functionNameRaw := request.GetStringParam(params, "FunctionName")
	functionName, embeddedQualifier := resolveFunctionRef(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return "", "", err
	}
	qualifier := mergeQualifier(request.GetStringParam(params, "Qualifier"), embeddedQualifier)
	return functionName, qualifier, nil
}

// PutProvisionedConcurrencyConfig configures provisioned concurrency for a Lambda function alias or version.
func (s *LambdaService) PutProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveProvisionedConcurrencyTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	config, err := s.putProvisionedConcurrencyCore(reqCtx, &ProvisionedConcurrencyInput{
		FunctionName:                    functionName,
		Qualifier:                       qualifier,
		ProvisionedConcurrentExecutions: int32(request.GetIntParam(req.Parameters, "ProvisionedConcurrentExecutions")),
		Region:                          reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return s.toProvisionedConcurrencyConfig(config), nil
}

// GetProvisionedConcurrencyConfig retrieves the provisioned concurrency configuration for a Lambda function alias or version.
func (s *LambdaService) GetProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveProvisionedConcurrencyTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := s.getProvisionedConcurrencyCore(store, functionName, qualifier)
	if err != nil {
		return nil, err
	}

	return s.toProvisionedConcurrencyConfig(config), nil
}

// DeleteProvisionedConcurrencyConfig removes the provisioned concurrency configuration for a Lambda function alias or version.
func (s *LambdaService) DeleteProvisionedConcurrencyConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveProvisionedConcurrencyTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteProvisionedConcurrencyCore(store, functionName, qualifier); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListProvisionedConcurrencyConfigs lists the provisioned concurrency configurations for a Lambda function.
func (s *LambdaService) ListProvisionedConcurrencyConfigs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameRaw := request.GetStringParam(req.Parameters, "FunctionName")
	functionName := extractFunctionName(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configs, err := s.listProvisionedConcurrencyConfigsCore(store, functionName)
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

	resp := map[string]interface{}{
		"ProvisionedConcurrencyConfigs": items,
	}
	if pageResult.IsTruncated {
		resp["NextMarker"] = pageResult.NextMarker
	}

	return resp, nil
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
