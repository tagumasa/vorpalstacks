package cognitoidentityprovider

import (
	"context"
	"errors"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/internal/store/aws/common"
)

// CreateUserPoolClient creates a user pool client for a Cognito user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolClient.html
func (s *CognitoService) CreateUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientName := req.GetParam("ClientName")
	if userPoolID == "" || clientName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	client := cognitostore.NewUserPoolClient(userPoolID, clientName)
	applyUserPoolClientParams(req, client)

	// Suppress client secret when GenerateSecret is explicitly false (CIPD-12)
	if v, ok := req.Parameters["GenerateSecret"].(bool); ok && !v {
		client.ClientSecret = ""
	} else if gs := req.GetParam("GenerateSecret"); gs == "false" || gs == "False" {
		client.ClientSecret = ""
	}

	if err := store.CreateUserPoolClient(client); err != nil {
		if errors.Is(err, cognitostore.ErrClientAlreadyExists) {
			return nil, ErrClientAlreadyExists
		}
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client, true),
	}, nil
}

// DescribeUserPoolClient returns information about a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DescribeUserPoolClient.html
func (s *CognitoService) DescribeUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientID := getClientId(req)
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrClientNotFound
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	client, err := store.GetUserPoolClient(userPoolID, clientID)
	if err != nil {
		return nil, ErrClientNotFound
	}

	if clientName := req.GetParam("ClientName"); clientName != "" {
		client.ClientName = clientName
	}
	applyUserPoolClientParams(req, client)

	if err := store.UpdateUserPoolClient(client); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"UserPoolClient": formatUserPoolClient(client, false),
	}, nil
}

func applyUserPoolClientParams(req *request.ParsedRequest, client *cognitostore.UserPoolClient) {
	if val := getIntParam(req, "RefreshTokenValidity"); val > 0 {
		client.RefreshTokenValidity = val
	}
	if val := getIntParam(req, "AccessTokenValidity"); val > 0 {
		client.AccessTokenValidity = val
	}
	if val := getIntParam(req, "IdTokenValidity"); val > 0 {
		client.IDTokenValidity = val
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
	// Parse AllowedOAuthFlowsUserPoolClient (CIPD-3)
	client.AllowedOAuthFlowsUserPoolClient = getBoolParam(req, "AllowedOAuthFlowsUserPoolClient")
	// Parse PreventUserExistenceErrors (CIPD-4)
	if v := req.GetParam("PreventUserExistenceErrors"); v != "" {
		client.PreventUserExistenceErrors = v
	}
}

// DeleteUserPoolClient deletes a user pool client.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolClient.html
func (s *CognitoService) DeleteUserPoolClient(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	clientID := getClientId(req)
	if userPoolID == "" || clientID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteUserPoolClient(userPoolID, clientID); err != nil {
		return nil, ErrClientNotFound
	}

	return response.EmptyResponse(), nil
}

// ListUserPoolClients lists the user pool clients for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolClients.html
func (s *CognitoService) ListUserPoolClients(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userPoolID := getUserPoolID(req)
	if userPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 || maxResults > 60 {
		maxResults = 60
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")

	opts := common.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	}

	result, err := store.ListUserPoolClientsPaginated(userPoolID, opts)
	if err != nil {
		return nil, ErrInternalError
	}

	clientList := make([]map[string]interface{}, 0, len(result.Items))
	for _, client := range result.Items {
		clientList = append(clientList, formatUserPoolClientSummary(client))
	}

	resp := map[string]interface{}{
		"UserPoolClients": clientList,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
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
	// Always include AllowedOAuthFlowsUserPoolClient (CIPD-5)
	result["AllowedOAuthFlowsUserPoolClient"] = client.AllowedOAuthFlowsUserPoolClient

	return result
}

func formatUserPoolClientSummary(client *cognitostore.UserPoolClient) map[string]interface{} {
	return map[string]interface{}{
		"ClientId":   client.ClientID,
		"UserPoolId": client.UserPoolID,
		"ClientName": client.ClientName,
	}
}
