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
	if poolName == "" {
		return nil, ErrInvalidParameter
	}

	userPool := cognitostore.NewUserPool(poolName, reqCtx.GetRegion())
	userPool.PasswordPolicy = parsePasswordPolicy(req)
	userPool.LambdaConfig = parseLambdaConfig(req)
	if v := request.GetStringList(req.Parameters, "AliasAttributes"); v != nil {
		userPool.AliasAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "UsernameAttributes"); v != nil {
		userPool.UsernameAttributes = v
	}
	if v := request.GetStringList(req.Parameters, "AutoVerifiedAttributes"); v != nil {
		userPool.AutoVerifiedAttributes = v
	}
	if v := req.GetParam("MfaConfiguration"); v != "" {
		userPool.MfaConfiguration = v
	}
	userPool.DeletionProtection = req.GetParam("DeletionProtection")
	userPool.EmailVerificationMessage = req.GetParam("EmailVerificationMessage")
	userPool.EmailVerificationSubject = req.GetParam("EmailVerificationSubject")
	userPool.SmsVerificationMessage = req.GetParam("SmsVerificationMessage")
	userPool.SmsAuthenticationMessage = req.GetParam("SmsAuthenticationMessage")
	userPool.EmailConfiguration = parseEmailConfiguration(req)
	userPool.SmsConfiguration = parseSmsConfiguration(req)
	userPool.AdminCreateUserConfig = parseAdminCreateUserConfig(req)
	userPool.VerificationMessageTemplate = parseVerificationMessageTemplate(req)
	userPool.UserAttributeUpdateSettings = parseUserAttributeUpdateSettings(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.CreateUserPool(userPool)
	if err != nil {
		return nil, ErrInternalError
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "UserPoolTags"))
	if len(tags) > 0 {
		if err := store.Tag(created.Arn, tags); err != nil {
			return nil, ErrInternalError
		}
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tags, _ := store.ListAsSlice(userPool.Arn)
	userPool.Tags = tags

	return map[string]interface{}{
		"UserPool": formatUserPool(userPool),
	}, nil
}

// DeleteUserPool deletes a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPool.html
func (s *CognitoService) DeleteUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteUserPool(userPoolID); err != nil {
		return nil, ErrResourceNotFound
	}

	return response.EmptyResponse(), nil
}

// UpdateUserPool updates a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPool.html
func (s *CognitoService) UpdateUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if poolName := req.GetParam("PoolName"); poolName != "" {
		userPool.Name = poolName
	}

	if policy := parsePasswordPolicy(req); policy != nil {
		userPool.PasswordPolicy = policy
	}

	if config := parseLambdaConfig(req); config != nil {
		userPool.LambdaConfig = config
	}

	if err := store.UpdateUserPool(userPool); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}

// ListUserPools lists the Cognito user pools.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPools.html
func (s *CognitoService) ListUserPools(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPools, err := store.ListUserPools()
	if err != nil {
		return nil, ErrInternalError
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 50
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	pools := make([]map[string]interface{}, 0)
	started := nextToken == ""
	for _, pool := range userPools {
		if !started {
			if pool.ID == nextToken {
				started = true
			}
			continue
		}
		pools = append(pools, map[string]interface{}{
			"Id":               pool.ID,
			"Name":             pool.Name,
			"Arn":              pool.Arn,
			"Status":           pool.Status,
			"CreationDate":     pool.CreationDate.Unix(),
			"LastModifiedDate": pool.LastModifiedDate.Unix(),
		})
		if len(pools) >= maxResults {
			break
		}
	}

	result := map[string]interface{}{
		"UserPools": pools,
	}
	if len(pools) >= maxResults && len(pools) > 0 {
		result["NextToken"] = pools[len(pools)-1]["Id"]
	}

	return result, nil
}

// GetUserPoolMfaConfig retrieves the multi-factor authentication configuration for a Cognito user pool.
func (s *CognitoService) GetUserPoolMfaConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return formatMfaConfigResponse(userPool), nil
}

// SetUserPoolMfaConfig updates the multi-factor authentication configuration for a Cognito user pool.
func (s *CognitoService) SetUserPoolMfaConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if mfaConfig := req.GetParam("MfaConfiguration"); mfaConfig != "" {
		userPool.MfaConfiguration = mfaConfig
	}

	if m, ok := req.Parameters["SmsMfaConfiguration"].(map[string]interface{}); ok {
		smsMfa := &cognitostore.MfaConfigurationType{}
		if enabled, ok := m["Enabled"].(bool); ok {
			smsMfa.Enabled = enabled
		}
		if smsConfig, ok := m["SmsConfiguration"].(map[string]interface{}); ok {
			poolSmsConfig := &cognitostore.SmsConfiguration{}
			if v, ok := smsConfig["SnsCallerArn"].(string); ok {
				poolSmsConfig.SnsCallerArn = v
			}
			if v, ok := smsConfig["ExternalId"].(string); ok {
				poolSmsConfig.ExternalId = v
			}
			userPool.SmsConfiguration = poolSmsConfig
			smsMfa.Enabled = true
		}
		userPool.MfaConfigurationSms = smsMfa
	}
	if m, ok := req.Parameters["SoftwareTokenMfaConfiguration"].(map[string]interface{}); ok {
		swMfa := &cognitostore.MfaConfigurationType{}
		if enabled, ok := m["Enabled"].(bool); ok {
			swMfa.Enabled = enabled
		}
		userPool.MfaConfigurationSoftwareToken = swMfa
	}

	if err := store.UpdateUserPool(userPool); err != nil {
		return nil, ErrInternalError
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
		smsEntry := map[string]interface{}{
			"Enabled": pool.MfaConfigurationSms.Enabled,
		}
		if pool.SmsConfiguration != nil {
			smsConfig := map[string]interface{}{}
			if pool.SmsConfiguration.SnsCallerArn != "" {
				smsConfig["SnsCallerArn"] = pool.SmsConfiguration.SnsCallerArn
			}
			if pool.SmsConfiguration.ExternalId != "" {
				smsConfig["ExternalId"] = pool.SmsConfiguration.ExternalId
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
	return result
}
