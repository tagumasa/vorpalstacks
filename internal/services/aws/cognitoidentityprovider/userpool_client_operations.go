package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CreateUserPoolClient creates a user pool client for a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolClient.html
func (s *CognitoService) CreateUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	client, err := s.newUserPoolClientCore(getUserPoolID(req), req.GetParam("ClientName"))
	if err != nil {
		return nil, err
	}
	applyUserPoolClientParams(req, client)
	// GenerateSecret and ClientSecret are members of the create request
	// only — the update request model defines neither, so the update
	// handler never reads them.
	if v, ok := getBoolParamOK(req, "GenerateSecret"); ok {
		client.GenerateSecret = v
	}

	if _, err := s.createUserPoolClientCore(reqCtx.GetRegion(), client, req.GetParam("ClientSecret")); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client),
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
		"UserPoolClient": formatUserPoolClient(client),
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
	applyUserPoolClientParams(req, client)

	if err := s.updateUserPoolClientCore(reqCtx.GetRegion(), client); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client),
	}, nil
}

// applyUserPoolClientParams applies the client-configuration members that
// the create and update request models share onto the store-level client.
// It is pure parameter application: the assembled configuration is
// validated by validateUserPoolClientConfig inside the create/update Cores,
// never here. Boolean members are applied on presence — an omitted member
// keeps the stored value — and create-only members (GenerateSecret) are
// handled by the create handler, not here.
func applyUserPoolClientParams(req *request.ParsedRequest, client *cognitostore.UserPoolClient) {
	// Token-validity members apply on presence so an explicit out-of-range
	// value reaches the Core's validation (the model's minimum is 1); an
	// omitted member keeps the stored value.
	if _, ok := req.Parameters["RefreshTokenValidity"]; ok {
		client.RefreshTokenValidity = getIntParam(req, "RefreshTokenValidity")
	}
	if _, ok := req.Parameters["AccessTokenValidity"]; ok {
		client.AccessTokenValidity = getIntParam(req, "AccessTokenValidity")
	}
	if _, ok := req.Parameters["IdTokenValidity"]; ok {
		client.IDTokenValidity = getIntParam(req, "IdTokenValidity")
	}
	if flows := getStringSliceParam(req, "ExplicitAuthFlows"); len(flows) > 0 {
		client.ExplicitAuthFlows = flows
	}
	if flows := getStringSliceParam(req, "AllowedOAuthFlows"); len(flows) > 0 {
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
	if v, ok := getBoolParamOK(req, "AllowedOAuthFlowsUserPoolClient"); ok {
		client.AllowedOAuthFlowsUserPoolClient = v
	}
	// Parse PreventUserExistenceErrors.
	if v := req.GetParam("PreventUserExistenceErrors"); v != "" {
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
	if v, ok := getBoolParamOK(req, "EnablePropagateAdditionalUserContextData"); ok {
		client.EnablePropagateAdditionalUserContextData = v
	}
	if v, ok := getBoolParamOK(req, "EnableTokenRevocation"); ok {
		client.EnableTokenRevocation = v
	}
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
			tvu.AccessToken = v
		}
		if v, ok := m["IdToken"].(string); ok {
			tvu.IdToken = v
		}
		if v, ok := m["RefreshToken"].(string); ok {
			tvu.RefreshToken = v
		}
		client.TokenValidityUnits = tvu
	}
	if m, ok := req.Parameters["RefreshTokenRotation"].(map[string]interface{}); ok {
		rtr := &cognitostore.RefreshTokenRotation{}
		if v, ok := m["Feature"].(string); ok {
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
	maxResults, err := parseStrictListLimit(req.Parameters, "MaxResults", listLimitMax)
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

// getStringSliceParam reads a JSON array of strings under its awsJson1_1
// member name; anything else (including a query-form Key.N family) reads
// as absent.
func getStringSliceParam(req *request.ParsedRequest, key string) []string {
	arr, ok := req.Parameters[key].([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func formatUserPoolClient(client *cognitostore.UserPoolClient) map[string]interface{} {
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

	// ClientSecret is a member of the shared response type: it is
	// returned whenever the stored client carries one, on create, update
	// and describe alike, and omitted for clients without a secret.
	if client.ClientSecret != "" {
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
