package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
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
// Handles both flat map format ({"mykey": "myval"}) and CLI single-tag format ({"Key": "k", "Value": "v"}).
func parseTags(params map[string]interface{}) map[string]string {
	raw := request.GetMapParam(params, "tags")
	if raw == nil {
		raw = request.GetMapParam(params, "Tags")
	}
	if raw == nil {
		return nil
	}
	if k, ok := raw["key"]; ok {
		if v, ok2 := raw["value"]; ok2 {
			if ks, ok3 := k.(string); ok3 {
				if vs, ok4 := v.(string); ok4 {
					return map[string]string{ks: vs}
				}
			}
		}
	}
	if k, ok := raw["Key"]; ok {
		if v, ok2 := raw["Value"]; ok2 {
			if ks, ok3 := k.(string); ok3 {
				if vs, ok4 := v.(string); ok4 {
					return map[string]string{ks: vs}
				}
			}
		}
	}
	result := make(map[string]string)
	for k, v := range raw {
		if vs, ok := v.(string); ok {
			result[k] = vs
		}
	}
	return result
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

// --- GraphQL API parse helpers ---

// parseAdditionalAuthProviders parses additional authentication providers from request parameters.
func parseAdditionalAuthProviders(params map[string]interface{}) []appsyncstore.AdditionalAuthenticationProvider {
	raw := request.GetArrayParam(params, "additionalAuthenticationProviders")
	var providers []appsyncstore.AdditionalAuthenticationProvider
	for _, p := range raw {
		if pMap, ok := p.(map[string]interface{}); ok {
			ap := appsyncstore.AdditionalAuthenticationProvider{
				AuthenticationType: request.GetStringParam(pMap, "authenticationType"),
			}
			if lambdaRaw := request.GetMapParam(pMap, "lambdaAuthorizerConfig"); lambdaRaw != nil {
				ap.LambdaAuthorizerConfig = parseLambdaAuthorizerConfigFromMap(lambdaRaw)
			}
			if oidcRaw := request.GetMapParam(pMap, "openIDConnectConfig"); oidcRaw != nil {
				ap.OpenIDConnectConfig = parseOpenIDConnectConfigFromMap(oidcRaw)
			}
			if cognitoRaw := request.GetMapParam(pMap, "userPoolConfig"); cognitoRaw != nil {
				ap.UserPoolConfig = &appsyncstore.CognitoUserPoolConfig{
					AwsRegion:        request.GetStringParam(cognitoRaw, "awsRegion"),
					UserPoolId:       request.GetStringParam(cognitoRaw, "userPoolId"),
					AppIdClientRegex: request.GetStringParam(cognitoRaw, "appIdClientRegex"),
				}
			}
			providers = append(providers, ap)
		}
	}
	return providers
}

// parseEnhancedMetricsConfig parses an EnhancedMetricsConfig from request parameters.
func parseEnhancedMetricsConfig(params map[string]interface{}) *appsyncstore.EnhancedMetricsConfig {
	raw := request.GetMapParam(params, "enhancedMetricsConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.EnhancedMetricsConfig{
		DataSourceLevelMetricsBehavior: request.GetStringParam(raw, "dataSourceLevelMetricsBehavior"),
		OperationLevelMetricsConfig:    request.GetStringParam(raw, "operationLevelMetricsConfig"),
		ResolverLevelMetricsBehavior:   request.GetStringParam(raw, "resolverLevelMetricsBehavior"),
	}
}

// parseLambdaAuthorizerConfig parses a LambdaAuthorizerConfig from request parameters.
func parseLambdaAuthorizerConfig(params map[string]interface{}) *appsyncstore.LambdaAuthorizerConfig {
	raw := request.GetMapParam(params, "lambdaAuthorizerConfig")
	if raw == nil {
		return nil
	}
	return parseLambdaAuthorizerConfigFromMap(raw)
}

// parseLambdaAuthorizerConfigFromMap parses a LambdaAuthorizerConfig from a raw map.
func parseLambdaAuthorizerConfigFromMap(raw map[string]interface{}) *appsyncstore.LambdaAuthorizerConfig {
	return &appsyncstore.LambdaAuthorizerConfig{
		AuthorizerUri:                request.GetStringParam(raw, "authorizerUri"),
		AuthorizerResultTtlInSeconds: int32(request.GetIntParam(raw, "authorizerResultTtlInSeconds")),
		IdentityValidationExpression: request.GetStringParam(raw, "identityValidationExpression"),
	}
}

// parseLogConfig parses a LogConfig from request parameters.
func parseLogConfig(params map[string]interface{}) *appsyncstore.LogConfig {
	raw := request.GetMapParam(params, "logConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.LogConfig{
		CloudWatchLogsRoleArn: request.GetStringParam(raw, "cloudWatchLogsRoleArn"),
		FieldLogLevel:         request.GetStringParam(raw, "fieldLogLevel"),
		ExcludeVerboseContent: request.GetBoolParam(raw, "excludeVerboseContent"),
	}
}

// parseOpenIDConnectConfig parses an OpenIDConnectConfig from request parameters.
func parseOpenIDConnectConfig(params map[string]interface{}) *appsyncstore.OpenIDConnectConfig {
	raw := request.GetMapParam(params, "openIDConnectConfig")
	if raw == nil {
		return nil
	}
	return parseOpenIDConnectConfigFromMap(raw)
}

// parseOpenIDConnectConfigFromMap parses an OpenIDConnectConfig from a raw map.
func parseOpenIDConnectConfigFromMap(raw map[string]interface{}) *appsyncstore.OpenIDConnectConfig {
	return &appsyncstore.OpenIDConnectConfig{
		Issuer:   request.GetStringParam(raw, "issuer"),
		AuthTTL:  request.GetInt64Param(raw, "authTTL"),
		ClientId: request.GetStringParam(raw, "clientId"),
		IatTTL:   request.GetInt64Param(raw, "iatTTL"),
	}
}

// parseUserPoolConfig parses a UserPoolConfig from request parameters.
func parseUserPoolConfig(params map[string]interface{}) *appsyncstore.UserPoolConfig {
	raw := request.GetMapParam(params, "userPoolConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.UserPoolConfig{
		AwsRegion:        request.GetStringParam(raw, "awsRegion"),
		DefaultAction:    request.GetStringParam(raw, "defaultAction"),
		UserPoolId:       request.GetStringParam(raw, "userPoolId"),
		AppIdClientRegex: request.GetStringParam(raw, "appIdClientRegex"),
	}
}

// --- Data Source parse helpers ---

// parseDynamoDBConfig parses a DynamoDB data source config from request parameters.
func parseDynamoDBConfig(params map[string]interface{}) *appsyncstore.DynamodbDataSourceConfig {
	raw := request.GetMapParam(params, "dynamodbConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.DynamodbDataSourceConfig{
		AwsRegion:            request.GetStringParam(raw, "awsRegion"),
		TableName:            request.GetStringParam(raw, "tableName"),
		UseCallerCredentials: request.GetBoolParam(raw, "useCallerCredentials"),
		Versioned:            request.GetBoolParam(raw, "versioned"),
	}
	if deltaRaw := request.GetMapParam(raw, "deltaSyncConfig"); deltaRaw != nil {
		cfg.DeltaSyncConfig = &appsyncstore.DeltaSyncConfig{
			BaseTableTTL:       request.GetInt64Param(deltaRaw, "baseTableTTL"),
			DeltaSyncTableName: request.GetStringParam(deltaRaw, "deltaSyncTableName"),
			DeltaSyncTableTTL:  request.GetInt64Param(deltaRaw, "deltaSyncTableTTL"),
		}
	}
	return cfg
}

// parseElasticsearchConfig parses an Elasticsearch data source config from request parameters.
func parseElasticsearchConfig(params map[string]interface{}) *appsyncstore.ElasticsearchDataSourceConfig {
	raw := request.GetMapParam(params, "elasticsearchConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.ElasticsearchDataSourceConfig{
		AwsRegion: request.GetStringParam(raw, "awsRegion"),
		Endpoint:  request.GetStringParam(raw, "endpoint"),
	}
}

// parseEventBridgeConfig parses an EventBridge data source config from request parameters.
func parseEventBridgeConfig(params map[string]interface{}) *appsyncstore.EventBridgeDataSourceConfig {
	raw := request.GetMapParam(params, "eventBridgeConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.EventBridgeDataSourceConfig{
		EventBusArn: request.GetStringParam(raw, "eventBusArn"),
	}
}

// parseHttpConfig parses an HTTP data source config from request parameters.
func parseHttpConfig(params map[string]interface{}) *appsyncstore.HttpDataSourceConfig {
	raw := request.GetMapParam(params, "httpConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.HttpDataSourceConfig{
		Endpoint: request.GetStringParam(raw, "endpoint"),
	}
	if authRaw := request.GetMapParam(raw, "authorizationConfig"); authRaw != nil {
		cfg.AuthorizationConfig = &appsyncstore.AuthorizationConfig{
			AuthorizationType: request.GetStringParam(authRaw, "authorizationType"),
		}
		if iamRaw := request.GetMapParam(authRaw, "awsIamConfig"); iamRaw != nil {
			cfg.AuthorizationConfig.AwsIamConfig = &appsyncstore.AwsIamConfig{
				SigningRegion:      request.GetStringParam(iamRaw, "signingRegion"),
				SigningServiceName: request.GetStringParam(iamRaw, "signingServiceName"),
			}
		}
	}
	return cfg
}

// parseLambdaDataSourceConfig parses a Lambda data source config from request parameters.
func parseLambdaDataSourceConfig(params map[string]interface{}) *appsyncstore.LambdaDataSourceConfig {
	raw := request.GetMapParam(params, "lambdaConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.LambdaDataSourceConfig{
		LambdaFunctionArn: request.GetStringParam(raw, "lambdaFunctionArn"),
	}
}

// parseOpenSearchServiceConfig parses an OpenSearch Service data source config from request parameters.
func parseOpenSearchServiceConfig(params map[string]interface{}) *appsyncstore.OpenSearchServiceDataSourceConfig {
	raw := request.GetMapParam(params, "openSearchServiceConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.OpenSearchServiceDataSourceConfig{
		AwsRegion: request.GetStringParam(raw, "awsRegion"),
		Endpoint:  request.GetStringParam(raw, "endpoint"),
	}
}

// parseRelationalDatabaseConfig parses a relational database data source config from request parameters.
func parseRelationalDatabaseConfig(params map[string]interface{}) *appsyncstore.RelationalDatabaseDataSourceConfig {
	raw := request.GetMapParam(params, "relationalDatabaseConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.RelationalDatabaseDataSourceConfig{
		RelationalDatabaseSourceType: request.GetStringParam(raw, "relationalDatabaseSourceType"),
	}
	if rdsRaw := request.GetMapParam(raw, "rdsHttpEndpointConfig"); rdsRaw != nil {
		cfg.RdsHttpEndpointConfig = &appsyncstore.RdsHttpEndpointConfig{
			AwsRegion:           request.GetStringParam(rdsRaw, "awsRegion"),
			AwsSecretStoreArn:   request.GetStringParam(rdsRaw, "awsSecretStoreArn"),
			DatabaseName:        request.GetStringParam(rdsRaw, "databaseName"),
			DbClusterIdentifier: request.GetStringParam(rdsRaw, "dbClusterIdentifier"),
			Schema:              request.GetStringParam(rdsRaw, "schema"),
		}
	}
	return cfg
}

// parseNeptuneConfig parses a Neptune data source config from request parameters.
func parseNeptuneConfig(params map[string]interface{}) *appsyncstore.NeptuneDataSourceConfig {
	raw := request.GetMapParam(params, "neptuneConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.NeptuneDataSourceConfig{
		GraphID: request.GetStringParam(raw, "graphId"),
	}
}

// --- Resolver parse helpers ---

// parseAppSyncRuntime parses an AppSyncRuntime from request parameters.
func parseAppSyncRuntime(params map[string]interface{}) *appsyncstore.AppSyncRuntime {
	raw := request.GetMapParam(params, "runtime")
	if raw == nil {
		return nil
	}
	return &appsyncstore.AppSyncRuntime{
		Name:           request.GetStringParam(raw, "name"),
		RuntimeVersion: request.GetStringParam(raw, "runtimeVersion"),
	}
}

// parsePipelineConfig parses a PipelineConfig from request parameters.
func parsePipelineConfig(params map[string]interface{}) *appsyncstore.PipelineConfig {
	raw := request.GetMapParam(params, "pipelineConfig")
	if raw == nil {
		return nil
	}
	functions := request.GetStringList(raw, "functions")
	if len(functions) == 0 {
		return nil
	}
	return &appsyncstore.PipelineConfig{
		Functions: functions,
	}
}

// parseCachingConfig parses a CachingConfig from request parameters.
func parseCachingConfig(params map[string]interface{}) *appsyncstore.CachingConfig {
	raw := request.GetMapParam(params, "cachingConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.CachingConfig{
		CachingKeys: request.GetStringList(raw, "cachingKeys"),
		Ttl:         request.GetInt64Param(raw, "ttl"),
	}
}

// parseSyncConfig parses a SyncConfig from request parameters.
func parseSyncConfig(params map[string]interface{}) *appsyncstore.SyncConfig {
	raw := request.GetMapParam(params, "syncConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.SyncConfig{
		ConflictDetection: request.GetStringParam(raw, "conflictDetection"),
		ConflictHandler:   request.GetStringParam(raw, "conflictHandler"),
	}
	if lambdaRaw := request.GetMapParam(raw, "lambdaConflictHandlerConfig"); lambdaRaw != nil {
		cfg.LambdaConflictHandlerConfig = &appsyncstore.LambdaConflictHandlerConfig{
			LambdaConflictHandlerArn: request.GetStringParam(lambdaRaw, "lambdaConflictHandlerArn"),
		}
	}
	return cfg
}
