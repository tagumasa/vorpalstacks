package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

func (s *IoTService) CreatePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createPolicyVersionCore(store, CreatePolicyVersionInput{
		PolicyName:     request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		PolicyDocument: request.GetParamCaseInsensitive(req.Parameters, "policyDocument"),
		SetAsDefault:   request.GetBoolParam(req.Parameters, "setAsDefault"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"policyArn":        result.PolicyARN,
		"policyDocument":   result.PolicyDocument,
		"policyVersionId":  result.PolicyVersionID,
		"isDefaultVersion": result.IsDefault,
	}, nil
}

func (s *IoTService) DeletePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deletePolicyVersionCore(store, PolicyVersionInput{
		PolicyName:      request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		PolicyVersionID: request.GetParamCaseInsensitive(req.Parameters, "policyVersionId"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) SetDefaultPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setDefaultPolicyVersionCore(store, PolicyVersionInput{
		PolicyName:      request.GetParamCaseInsensitive(req.Parameters, "policyName"),
		PolicyVersionID: request.GetParamCaseInsensitive(req.Parameters, "policyVersionId"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListTargetsForPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principals, err := s.listPolicyPrincipalsCore(store, request.GetParamCaseInsensitive(req.Parameters, "policyName"))
	if err != nil {
		return nil, err
	}

	// AWS SDK v2 expects Targets as a flat []string (PolicyTarget is a
	// string shape in the Smithy model), not an array of objects.
	return paginatedStrings("targets", principals, req.Parameters)
}

// AttachPrincipalPolicy is the legacy alias of AttachPolicy. AWS accepts a
// "principal" parameter instead of "target" (Smithy AttachPrincipalPolicyRequest
// vs AttachPolicyRequest). Map the parameter so the alias dispatches through
// AttachPolicy without losing the principal identifier.
func (s *IoTService) AttachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.AttachPolicy(ctx, reqCtx, req)
}

// DetachPrincipalPolicy is the legacy alias of DetachPolicy. Same parameter
// remap as AttachPrincipalPolicy.
func (s *IoTService) DetachPrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if _, ok := req.Parameters["target"]; !ok {
		if principal, ok := req.Parameters["principal"]; ok {
			req.Parameters["target"] = principal
		}
	}
	return s.DetachPolicy(ctx, reqCtx, req)
}

func (s *IoTService) TestInvokeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.testInvokeAuthorizerCore(ctx, store, TestInvokeAuthorizerInput{
		AuthorizerName: request.GetParamCaseInsensitive(req.Parameters, "authorizerName"),
		Token:          request.GetParamCaseInsensitive(req.Parameters, "token"),
		TokenSignature: request.GetParamCaseInsensitive(req.Parameters, "tokenSignature"),
		MqttContext:    request.GetMapParamCaseInsensitive(req.Parameters, "mqttContext"),
		HttpContext:    request.GetMapParamCaseInsensitive(req.Parameters, "httpContext"),
		TlsContext:     request.GetMapParamCaseInsensitive(req.Parameters, "tlsContext"),
	})
}

func (s *IoTService) TestAuthorization(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// authInfos entries carry actionType plus a resources list.
	var authInfos []TestAuthorizationAuthInfo
	if raw, ok := req.Parameters["authInfos"].([]interface{}); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				actionType, _ := m["actionType"].(string)
				authInfos = append(authInfos, TestAuthorizationAuthInfo{
					ActionType: actionType,
					Resources:  request.GetStringList(m, "resources"),
				})
			}
		}
	}

	results, err := s.testAuthorizationCore(store, TestAuthorizationInput{
		Principal:             request.GetParamCaseInsensitive(req.Parameters, "principal"),
		ClientID:              request.GetParamCaseInsensitive(req.Parameters, "clientId"),
		CognitoIdentityPoolID: request.GetParamCaseInsensitive(req.Parameters, "cognitoIdentityPoolId"),
		PolicyNamesToAdd:      request.GetStringList(req.Parameters, "policyNamesToAdd"),
		PolicyNamesToSkip:     request.GetStringList(req.Parameters, "policyNamesToSkip"),
		AuthInfos:             authInfos,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"authResults": results}, nil
}
