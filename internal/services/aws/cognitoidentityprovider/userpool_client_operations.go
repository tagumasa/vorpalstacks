package cognitoidentityprovider

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateUserPoolClient creates a user pool client for a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolClient.html
func (s *CognitoService) CreateUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientName := req.GetParam("ClientName")
	if userPoolID == "" || clientName == "" {
		return nil, ErrInvalidParameter
	}

	client := cognitostore.NewUserPoolClient(userPoolID, clientName)
	if err := applyUserPoolClientParams(req, client); err != nil {
		return nil, err
	}

	// Suppress the client secret when GenerateSecret is explicitly false.
	// The value may arrive as a JSON bool or a query-string string.
	if !client.GenerateSecret {
		client.ClientSecret = ""
	}

	if _, err := s.createUserPoolClientCore(reqCtx.GetRegion(), client); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client, true),
	}, nil
}

// DescribeUserPoolClient returns information about a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeUserPoolClient.html
func (s *CognitoService) DescribeUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	client, err := s.describeUserPoolClientCore(reqCtx.GetRegion(), getUserPoolID(req), getClientId(req))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client, false),
	}, nil
}

// UpdateUserPoolClient updates a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPoolClient.html
func (s *CognitoService) UpdateUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientID := getClientId(req)
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	client, err := s.describeUserPoolClientCore(reqCtx.GetRegion(), userPoolID, clientID)
	if err != nil {
		return nil, err
	}

	if clientName := req.GetParam("ClientName"); clientName != "" {
		client.ClientName = clientName
	}
	if err := applyUserPoolClientParams(req, client); err != nil {
		return nil, err
	}

	if err := s.updateUserPoolClientCore(reqCtx.GetRegion(), client); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client, false),
	}, nil
}

func applyUserPoolClientParams(req *request.ParsedRequest, client *cognitostore.UserPoolClient) error {
	if val := getIntParam(req, "RefreshTokenValidity"); val > 0 {
		if !validateRefreshTokenValidity(val) {
			return ErrInvalidParameter
		}
		client.RefreshTokenValidity = val
	}
	if val := getIntParam(req, "AccessTokenValidity"); val > 0 {
		if !validateAccessTokenValidity(val) {
			return ErrInvalidParameter
		}
		client.AccessTokenValidity = val
	}
	if val := getIntParam(req, "IdTokenValidity"); val > 0 {
		if !validateIdTokenValidity(val) {
			return ErrInvalidParameter
		}
		client.IDTokenValidity = val
	}
	if flows := getStringSliceParam(req, "ExplicitAuthFlows"); len(flows) > 0 {
		for _, f := range flows {
			if !validateExplicitAuthFlow(f) {
				return ErrInvalidParameter
			}
		}
		client.ExplicitAuthFlows = flows
	}
	if flows := getStringSliceParam(req, "AllowedOAuthFlows"); len(flows) > 0 {
		for _, f := range flows {
			if !validateOAuthFlow(f) {
				return ErrInvalidParameter
			}
		}
		client.AllowedOAuthFlows = flows
	}
	if urls := getStringSliceParam(req, "CallbackURLs"); len(urls) > 0 {
		client.CallbackURLs = urls
	}
	if urls := getStringSliceParam(req, "LogoutURLs"); len(urls) > 0 {
		client.LogoutURLs = urls
	}
	if uri := req.GetParam("DefaultRedirectURI"); uri != "" {
		client.DefaultRedirectURI = uri
	}
	if providers := getStringSliceParam(req, "SupportedIdentityProviders"); len(providers) > 0 {
		client.SupportedIdentityProviders = providers
	}
	if scopes := getStringSliceParam(req, "AllowedOAuthScopes"); len(scopes) > 0 {
		client.AllowedOAuthScopes = scopes
	}
	// Parse AllowedOAuthFlowsUserPoolClient.
	client.AllowedOAuthFlowsUserPoolClient = getBoolParam(req, "AllowedOAuthFlowsUserPoolClient")
	// Parse PreventUserExistenceErrors.
	if v := req.GetParam("PreventUserExistenceErrors"); v != "" {
		if !validatePreventUserExistenceErrors(v) {
			return ErrInvalidParameter
		}
		client.PreventUserExistenceErrors = v
	}
	// Parse missing Smithy fields.
	if val := getIntParam(req, "AuthSessionValidity"); val > 0 {
		client.AuthSessionValidity = val
	}
	if attrs := getStringSliceParam(req, "ReadAttributes"); len(attrs) > 0 {
		client.ReadAttributes = attrs
	}
	if attrs := getStringSliceParam(req, "WriteAttributes"); len(attrs) > 0 {
		client.WriteAttributes = attrs
	}
	client.EnablePropagateAdditionalUserContextData = getBoolParam(req, "EnablePropagateAdditionalUserContextData")
	client.EnableTokenRevocation = getBoolParam(req, "EnableTokenRevocation")
	client.GenerateSecret = getBoolParam(req, "GenerateSecret")
	if m, ok := req.Parameters["AnalyticsConfiguration"].(map[string]interface{}); ok {
		ac := &cognitostore.AnalyticsConfiguration{}
		if v, ok := m["ApplicationArn"].(string); ok {
			ac.ApplicationArn = v
		}
		if v, ok := m["ApplicationId"].(string); ok {
			ac.ApplicationId = v
		}
		if v, ok := m["ExternalId"].(string); ok {
			ac.ExternalId = v
		}
		if v, ok := m["RoleArn"].(string); ok {
			ac.RoleArn = v
		}
		if v, ok := m["UserDataShared"].(bool); ok {
			ac.UserDataShared = v
		}
		client.AnalyticsConfiguration = ac
	}
	if m, ok := req.Parameters["TokenValidityUnits"].(map[string]interface{}); ok {
		tvu := &cognitostore.TokenValidityUnits{}
		if v, ok := m["AccessToken"].(string); ok {
			if !validateTimeUnit(v) {
				return ErrInvalidParameter
			}
			tvu.AccessToken = v
		}
		if v, ok := m["IdToken"].(string); ok {
			if !validateTimeUnit(v) {
				return ErrInvalidParameter
			}
			tvu.IdToken = v
		}
		if v, ok := m["RefreshToken"].(string); ok {
			if !validateTimeUnit(v) {
				return ErrInvalidParameter
			}
			tvu.RefreshToken = v
		}
		client.TokenValidityUnits = tvu
	}
	if m, ok := req.Parameters["RefreshTokenRotation"].(map[string]interface{}); ok {
		rtr := &cognitostore.RefreshTokenRotation{}
		if v, ok := m["Feature"].(string); ok {
			if !validateFeatureType(v) {
				return ErrInvalidParameter
			}
			rtr.Feature = v
		}
		if v, ok := m["RetryGracePeriodSeconds"]; ok {
			switch n := v.(type) {
			case int:
				rtr.RetryGracePeriodSeconds = n
			case float64:
				rtr.RetryGracePeriodSeconds = int(n)
			}
		}
		client.RefreshTokenRotation = rtr
	}
	return nil
}

// DeleteUserPoolClient deletes a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolClient.html
func (s *CognitoService) DeleteUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteUserPoolClientCore(reqCtx.GetRegion(), getUserPoolID(req), getClientId(req)); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListUserPoolClients lists the user pool clients for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolClients.html
func (s *CognitoService) ListUserPoolClients(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy QueryLimit: range {min: 1, max: 60}
	maxResults, err := parseStrictListLimit(req.Parameters, "MaxResults", 60)
	if err != nil {
		return nil, err
	}
	result, err := s.listUserPoolClientsCore(reqCtx.GetRegion(), ListUserPoolClientsInput{
		UserPoolID: getUserPoolID(req),
		MaxResults: maxResults,
		NextToken:  request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	clientList := make([]map[string]interface{}, 0, len(result.Clients))
	for _, c := range result.Clients {
		clientList = append(clientList, map[string]interface{}{
			"ClientId":   c.ClientID,
			"UserPoolId": c.UserPoolID,
			"ClientName": c.ClientName,
		})
	}

	resp := map[string]interface{}{
		"UserPoolClients": clientList,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

func getStringSliceParam(req *request.ParsedRequest, key string) []string {
	var result []string
	for i := 1; ; i++ {
		idx := strconv.Itoa(i)
		itemKey := key + "." + idx
		item := req.GetParam(itemKey)
		if item == "" {
			break
		}
		result = append(result, item)
	}

	if len(result) == 0 {
		if arr, ok := req.Parameters[key].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					result = append(result, s)
				}
			}
		}
	}

	return result
}

func formatUserPoolClient(client *cognitostore.UserPoolClient, includeSecret bool) map[string]interface{} {
	result := map[string]interface{}{
		"ClientId":             client.ClientID,
		"UserPoolId":           client.UserPoolID,
		"ClientName":           client.ClientName,
		"RefreshTokenValidity": client.RefreshTokenValidity,
		"AccessTokenValidity":  client.AccessTokenValidity,
		"IdTokenValidity":      client.IDTokenValidity,
		"CreationDate":         client.CreationDate.Unix(),
		"LastModifiedDate":     client.LastModifiedDate.Unix(),
	}

	if includeSecret && client.ClientSecret != "" {
		result["ClientSecret"] = client.ClientSecret
	}

	if len(client.ExplicitAuthFlows) > 0 {
		result["ExplicitAuthFlows"] = client.ExplicitAuthFlows
	}
	if len(client.AllowedOAuthFlows) > 0 {
		result["AllowedOAuthFlows"] = client.AllowedOAuthFlows
	}
	if len(client.CallbackURLs) > 0 {
		result["CallbackURLs"] = client.CallbackURLs
	}
	if len(client.LogoutURLs) > 0 {
		result["LogoutURLs"] = client.LogoutURLs
	}
	if client.DefaultRedirectURI != "" {
		result["DefaultRedirectURI"] = client.DefaultRedirectURI
	}
	if len(client.SupportedIdentityProviders) > 0 {
		result["SupportedIdentityProviders"] = client.SupportedIdentityProviders
	}
	if len(client.AllowedOAuthScopes) > 0 {
		result["AllowedOAuthScopes"] = client.AllowedOAuthScopes
	}
	if client.PreventUserExistenceErrors != "" {
		result["PreventUserExistenceErrors"] = client.PreventUserExistenceErrors
	}
	// Always include AllowedOAuthFlowsUserPoolClient.
	result["AllowedOAuthFlowsUserPoolClient"] = client.AllowedOAuthFlowsUserPoolClient
	result["EnablePropagateAdditionalUserContextData"] = client.EnablePropagateAdditionalUserContextData
	result["EnableTokenRevocation"] = client.EnableTokenRevocation

	if client.AuthSessionValidity > 0 {
		result["AuthSessionValidity"] = client.AuthSessionValidity
	}
	if len(client.ReadAttributes) > 0 {
		result["ReadAttributes"] = client.ReadAttributes
	}
	if len(client.WriteAttributes) > 0 {
		result["WriteAttributes"] = client.WriteAttributes
	}
	if client.AnalyticsConfiguration != nil {
		ac := map[string]interface{}{}
		if client.AnalyticsConfiguration.ApplicationArn != "" {
			ac["ApplicationArn"] = client.AnalyticsConfiguration.ApplicationArn
		}
		if client.AnalyticsConfiguration.ApplicationId != "" {
			ac["ApplicationId"] = client.AnalyticsConfiguration.ApplicationId
		}
		if client.AnalyticsConfiguration.ExternalId != "" {
			ac["ExternalId"] = client.AnalyticsConfiguration.ExternalId
		}
		if client.AnalyticsConfiguration.RoleArn != "" {
			ac["RoleArn"] = client.AnalyticsConfiguration.RoleArn
		}
		ac["UserDataShared"] = client.AnalyticsConfiguration.UserDataShared
		result["AnalyticsConfiguration"] = ac
	}
	if client.TokenValidityUnits != nil {
		tvu := map[string]interface{}{}
		if client.TokenValidityUnits.AccessToken != "" {
			tvu["AccessToken"] = client.TokenValidityUnits.AccessToken
		}
		if client.TokenValidityUnits.IdToken != "" {
			tvu["IdToken"] = client.TokenValidityUnits.IdToken
		}
		if client.TokenValidityUnits.RefreshToken != "" {
			tvu["RefreshToken"] = client.TokenValidityUnits.RefreshToken
		}
		result["TokenValidityUnits"] = tvu
	}
	if client.RefreshTokenRotation != nil {
		rtr := map[string]interface{}{}
		if client.RefreshTokenRotation.Feature != "" {
			rtr["Feature"] = client.RefreshTokenRotation.Feature
		}
		if client.RefreshTokenRotation.RetryGracePeriodSeconds > 0 {
			rtr["RetryGracePeriodSeconds"] = client.RefreshTokenRotation.RetryGracePeriodSeconds
		}
		result["RefreshTokenRotation"] = rtr
	}

	return result
}
