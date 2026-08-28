package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateUserPool creates a new Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPool.html
func (s *CognitoService) CreateUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := req.GetParam("PoolName")

	userPool, err := s.newUserPoolCore(poolName, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	// CreateUserPool is the only operation that carries the Schema member;
	// apply it before the shared update path so the whole-pool validation
	// still sees the schema definitions.
	userPool.SchemaAttributes = parseSchemaAttributes(req)
	if err := applyUserPoolUpdates(userPool, req); err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "UserPoolTags"))

	created, err := s.createUserPoolCore(CreateUserPoolInput{
		Pool:   userPool,
		Region: reqCtx.GetRegion(),
		Tags:   tags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPool": formatUserPool(created),
	}, nil
}

// DescribeUserPool returns the description of a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeUserPool.html
func (s *CognitoService) DescribeUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	userPool, err := s.describeUserPoolCore(reqCtx.GetRegion(), userPoolID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPool": formatUserPool(userPool),
	}, nil
}

// DeleteUserPool deletes a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPool.html
func (s *CognitoService) DeleteUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteUserPoolCore(reqCtx.GetRegion(), getUserPoolID(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UpdateUserPool updates a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPool.html
func (s *CognitoService) UpdateUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)

	userPool, err := s.getUserPoolCore(reqCtx.GetRegion(), userPoolID)
	if err != nil {
		return nil, err
	}

	if err := applyUserPoolUpdates(userPool, req); err != nil {
		return nil, err
	}

	if err := s.updateUserPoolPersistCore(reqCtx.GetRegion(), userPool); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListUserPools lists the Cognito user pools.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPools.html
func (s *CognitoService) ListUserPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy PoolQueryLimitType: range {min: 1, max: 60}
	maxResults, err := parseStrictListLimit(req.Parameters, "MaxResults", 60)
	if err != nil {
		return nil, err
	}
	result, err := s.listUserPoolsCore(reqCtx.GetRegion(), ListUserPoolsInput{
		MaxResults: maxResults,
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	pools := make([]map[string]interface{}, 0, len(result.UserPools))
	for _, pool := range result.UserPools {
		pools = append(pools, map[string]interface{}{
			"Id":               pool.ID,
			"Name":             pool.Name,
			"Arn":              pool.Arn,
			"Status":           pool.Status,
			"CreationDate":     pool.CreationDate.Unix(),
			"LastModifiedDate": pool.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"UserPools": pools,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

// GetUserPoolMfaConfig retrieves the multi-factor authentication configuration for a Cognito user pool.
func (s *CognitoService) GetUserPoolMfaConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)

	userPool, err := s.getUserPoolCore(reqCtx.GetRegion(), userPoolID)
	if err != nil {
		return nil, err
	}

	return formatMfaConfigResponse(userPool), nil
}

// SetUserPoolMfaConfig updates the multi-factor authentication configuration for a Cognito user pool.
func (s *CognitoService) SetUserPoolMfaConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := SetUserPoolMfaConfigInput{
		Region:           reqCtx.GetRegion(),
		UserPoolID:       getUserPoolID(req),
		MfaConfiguration: req.GetParam("MfaConfiguration"),
	}
	if m, ok := req.Parameters["SmsMfaConfiguration"].(map[string]interface{}); ok {
		in.SmsMfaConfiguration = m
	}
	if m, ok := req.Parameters["SoftwareTokenMfaConfiguration"].(map[string]interface{}); ok {
		in.SoftwareTokenMfaConfiguration = m
	}
	if m, ok := req.Parameters["EmailMfaConfiguration"].(map[string]interface{}); ok {
		in.EmailMfaConfiguration = m
	}
	if m, ok := req.Parameters["WebAuthnConfiguration"].(map[string]interface{}); ok {
		in.WebAuthnConfiguration = m
	}

	userPool, err := s.setUserPoolMfaConfigCore(in)
	if err != nil {
		return nil, err
	}

	return formatMfaConfigResponse(userPool), nil
}

func formatMfaConfigResponse(pool *cognitostore.UserPool) map[string]interface{} {
	mfaConfig := "OFF"
	if pool.MfaConfiguration != "" {
		mfaConfig = pool.MfaConfiguration
	}
	result := map[string]interface{}{
		"MfaConfiguration": mfaConfig,
	}
	if pool.MfaConfigurationSms != nil {
		smsEntry := map[string]interface{}{}
		if pool.MfaConfigurationSms.SmsAuthenticationMessage != "" {
			smsEntry["SmsAuthenticationMessage"] = pool.MfaConfigurationSms.SmsAuthenticationMessage
		}
		if pool.MfaConfigurationSms.SmsConfiguration != nil {
			smsConfig := map[string]interface{}{}
			if pool.MfaConfigurationSms.SmsConfiguration.SnsCallerArn != "" {
				smsConfig["SnsCallerArn"] = pool.MfaConfigurationSms.SmsConfiguration.SnsCallerArn
			}
			if pool.MfaConfigurationSms.SmsConfiguration.ExternalId != "" {
				smsConfig["ExternalId"] = pool.MfaConfigurationSms.SmsConfiguration.ExternalId
			}
			if pool.MfaConfigurationSms.SmsConfiguration.SnsRegion != "" {
				smsConfig["SnsRegion"] = pool.MfaConfigurationSms.SmsConfiguration.SnsRegion
			}
			if len(smsConfig) > 0 {
				smsEntry["SmsConfiguration"] = smsConfig
			}
		}
		result["SmsMfaConfiguration"] = smsEntry
	}
	if pool.MfaConfigurationSoftwareToken != nil {
		result["SoftwareTokenMfaConfiguration"] = map[string]interface{}{
			"Enabled": pool.MfaConfigurationSoftwareToken.Enabled,
		}
	}
	if pool.EmailMfaConfig != nil {
		emailEntry := map[string]interface{}{}
		if pool.EmailMfaConfig.Message != "" {
			emailEntry["Message"] = pool.EmailMfaConfig.Message
		}
		if pool.EmailMfaConfig.Subject != "" {
			emailEntry["Subject"] = pool.EmailMfaConfig.Subject
		}
		result["EmailMfaConfiguration"] = emailEntry
	}
	if pool.WebAuthnConfiguration != nil {
		waEntry := map[string]interface{}{}
		if pool.WebAuthnConfiguration.RelyingPartyId != "" {
			waEntry["RelyingPartyId"] = pool.WebAuthnConfiguration.RelyingPartyId
		}
		if pool.WebAuthnConfiguration.UserVerification != "" {
			waEntry["UserVerification"] = pool.WebAuthnConfiguration.UserVerification
		}
		result["WebAuthnConfiguration"] = waEntry
	}
	return result
}
