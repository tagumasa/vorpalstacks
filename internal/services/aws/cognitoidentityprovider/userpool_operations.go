package cognitoidentityprovider

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateUserPool creates a new Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPool.html
func (s *CognitoService) CreateUserPool(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	poolName := getParam(req, "PoolName")
	if poolName == "" {
		return nil, ErrInvalidParameter
	}

	userPool := cognitostore.NewUserPool(poolName, reqCtx.GetRegion())
	userPool.PasswordPolicy = parsePasswordPolicy(req)
	userPool.LambdaConfig = parseLambdaConfig(req)

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
		if err := store.TagResource(created.Arn, tags); err != nil {
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

	tags, _ := store.ListTagsAsSlice(userPool.Arn)
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

	mfaConfig := "OFF"
	if userPool.MfaConfiguration != "" {
		mfaConfig = userPool.MfaConfiguration
	}

	result := map[string]interface{}{
		"MfaConfiguration":    mfaConfig,
		"SmsMfaConfiguration": map[string]interface{}{"SmsConfiguration": map[string]interface{}{"SmsConfigType": "EXTERNAL"}},
	}

	return result, nil
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

	if mfaConfig := getParam(req, "MfaConfiguration"); mfaConfig != "" {
		userPool.MfaConfiguration = mfaConfig
	}

	if err := store.UpdateUserPool(userPool); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"MfaConfiguration":    userPool.MfaConfiguration,
		"SmsMfaConfiguration": map[string]interface{}{"SmsConfiguration": map[string]interface{}{"SmsConfigType": "EXTERNAL"}},
	}, nil
}
