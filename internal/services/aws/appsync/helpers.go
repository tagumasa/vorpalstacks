package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/store/aws/common"
)

// parsePaginationOptions extracts list pagination parameters from the request.
// AppSync uses maxResults (int) and nextToken (string) in query params.
func parsePaginationOptions(req *request.ParsedRequest) common.ListOptions {
	opts := common.ListOptions{
		MaxItems: request.GetIntParam(req.Parameters, "maxResults"),
		Marker:   request.GetStringParam(req.Parameters, "nextToken"),
	}
	if opts.MaxItems <= 0 {
		opts.MaxItems = 50
	}
	return opts
}

// parseEventConfig parses an EventConfig from the request parameters.
// EventConfig is required for CreateApi and UpdateApi (v2).
func parseEventConfig(params map[string]interface{}) (*appsyncstore.EventConfig, error) {
	ecRaw := request.GetMapParam(params, "eventConfig")
	if ecRaw == nil {
		return nil, NewBadRequestException("eventConfig is required")
	}

	ec := &appsyncstore.EventConfig{}

	if authProvidersRaw := request.GetArrayParam(ecRaw, "authProviders"); len(authProvidersRaw) > 0 {
		for _, apRaw := range authProvidersRaw {
			if apMap, ok := apRaw.(map[string]interface{}); ok {
				ap := appsyncstore.AuthProvider{
					AuthType: request.GetStringParam(apMap, "authType"),
				}
				if cognitoRaw := request.GetMapParam(apMap, "cognitoConfig"); cognitoRaw != nil {
					ap.CognitoConfig = &appsyncstore.CognitoEventConfig{
						AwsRegion:        request.GetStringParam(cognitoRaw, "awsRegion"),
						UserPoolId:       request.GetStringParam(cognitoRaw, "userPoolId"),
						AppIdClientRegex: request.GetStringParam(cognitoRaw, "appIdClientRegex"),
					}
				}
				if lambdaRaw := request.GetMapParam(apMap, "lambdaAuthorizerConfig"); lambdaRaw != nil {
					ap.LambdaAuthorizerConfig = &appsyncstore.LambdaAuthorizerConfig{
						AuthorizerUri:                request.GetStringParam(lambdaRaw, "authorizerUri"),
						AuthorizerResultTtlInSeconds: int32(request.GetIntParam(lambdaRaw, "authorizerResultTtlInSeconds")),
						IdentityValidationExpression: request.GetStringParam(lambdaRaw, "identityValidationExpression"),
					}
				}
				if oidcRaw := request.GetMapParam(apMap, "openIDConnectConfig"); oidcRaw != nil {
					ap.OpenIDConnectConfig = &appsyncstore.OpenIDConnectConfig{
						Issuer:   request.GetStringParam(oidcRaw, "issuer"),
						AuthTTL:  request.GetInt64Param(oidcRaw, "authTTL"),
						ClientId: request.GetStringParam(oidcRaw, "clientId"),
						IatTTL:   request.GetInt64Param(oidcRaw, "iatTTL"),
					}
				}
				ec.AuthProviders = append(ec.AuthProviders, ap)
			}
		}
	}

	ec.ConnectionAuthModes = parseAuthModes(request.GetArrayParam(ecRaw, "connectionAuthModes"))
	ec.DefaultPublishAuthModes = parseAuthModes(request.GetArrayParam(ecRaw, "defaultPublishAuthModes"))
	ec.DefaultSubscribeAuthModes = parseAuthModes(request.GetArrayParam(ecRaw, "defaultSubscribeAuthModes"))

	if logRaw := request.GetMapParam(ecRaw, "logConfig"); logRaw != nil {
		ec.LogConfig = &appsyncstore.EventLogConfig{
			CloudWatchLogsRoleArn: request.GetStringParam(logRaw, "cloudWatchLogsRoleArn"),
			LogLevel:              request.GetStringParam(logRaw, "logLevel"),
		}
	}

	return ec, nil
}

// parseAuthModes converts a raw array of auth mode maps into a slice of AuthMode.
func parseAuthModes(raw []interface{}) []appsyncstore.AuthMode {
	var modes []appsyncstore.AuthMode
	for _, m := range raw {
		if mMap, ok := m.(map[string]interface{}); ok {
			modes = append(modes, appsyncstore.AuthMode{
				AuthType: request.GetStringParam(mMap, "authType"),
			})
		}
	}
	return modes
}

// parseTags extracts a map[string]string from the "tags" request parameter.
// AppSync uses flat map tags (not the []Tag format used by most other services).
func parseTags(params map[string]interface{}) map[string]string {
	return request.ParseAttributes(params, "tags")
}

// parseHandlerConfigs parses HandlerConfigs from the request parameters.
func parseHandlerConfigs(params map[string]interface{}) *appsyncstore.HandlerConfigs {
	hcRaw := request.GetMapParam(params, "handlerConfigs")
	if hcRaw == nil {
		return nil
	}
	hc := &appsyncstore.HandlerConfigs{}

	if pubRaw := request.GetMapParam(hcRaw, "onPublish"); pubRaw != nil {
		hc.OnPublish = parseHandlerConfig(pubRaw)
	}
	if subRaw := request.GetMapParam(hcRaw, "onSubscribe"); subRaw != nil {
		hc.OnSubscribe = parseHandlerConfig(subRaw)
	}
	return hc
}

// parseHandlerConfig parses a single HandlerConfig from a map.
func parseHandlerConfig(raw map[string]interface{}) *appsyncstore.HandlerConfig {
	hc := &appsyncstore.HandlerConfig{
		Behavior: request.GetStringParam(raw, "behavior"),
	}
	if intRaw := request.GetMapParam(raw, "integration"); intRaw != nil {
		hc.Integration = &appsyncstore.Integration{
			DataSourceName: request.GetStringParam(intRaw, "dataSourceName"),
		}
		if lambdaRaw := request.GetMapParam(intRaw, "lambdaConfig"); lambdaRaw != nil {
			hc.Integration.LambdaConfig = &appsyncstore.LambdaIntConfig{
				InvokeType: request.GetStringParam(lambdaRaw, "invokeType"),
			}
		}
	}
	return hc
}
