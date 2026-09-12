package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Handlers for the provisioned-limit family; validation and store access
// live in provisioned_limit_core.go.

// parseProvisionedLimitInput extracts the LimitDefinition wire members
// shared by the provisioned-limit operations; update requests also carry
// the required RequestedLimitValue, whose presence the Core validates.
func parseProvisionedLimitInput(region string, req *request.ParsedRequest, withRequestedValue bool) (ProvisionedLimitInput, error) {
	in := ProvisionedLimitInput{Region: region}

	limitDef, ok := req.Parameters["LimitDefinition"].(map[string]interface{})
	if !ok {
		return in, ErrInvalidParameter
	}
	in.LimitClass = getStringParam(limitDef, "LimitClass")
	attrs, _ := limitDef["Attributes"].(map[string]interface{})
	in.Category = getStringParam(attrs, "Category")

	if withRequestedValue {
		value, present := parseIntParam(req, "RequestedLimitValue")
		in.RequestedValue = value
		in.RequestedValueSet = present
	}
	return in, nil
}

// formatProvisionedLimit serialises the LimitType response member.
func formatProvisionedLimit(limit *ProvisionedLimitResult) map[string]interface{} {
	return map[string]interface{}{
		"LimitDefinition": map[string]interface{}{
			"LimitClass": limit.LimitClass,
			"Attributes": map[string]interface{}{"Category": limit.Category},
		},
		"ProvisionedLimitValue": limit.ProvisionedLimitValue,
		"FreeLimitValue":        limit.FreeLimitValue,
	}
}

// GetProvisionedLimit retrieves a provisioned limit.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_GetProvisionedLimit.html
func (s *CognitoService) GetProvisionedLimit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in, err := parseProvisionedLimitInput(reqCtx.GetRegion(), req, false)
	if err != nil {
		return nil, err
	}

	limit, err := s.getProvisionedLimitCore(in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Limit": formatProvisionedLimit(limit),
	}, nil
}

// UpdateProvisionedLimit updates a provisioned limit.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateProvisionedLimit.html
func (s *CognitoService) UpdateProvisionedLimit(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in, err := parseProvisionedLimitInput(reqCtx.GetRegion(), req, true)
	if err != nil {
		return nil, err
	}

	limit, err := s.updateProvisionedLimitCore(in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Limit": formatProvisionedLimit(limit),
	}, nil
}
