package appsync

import (
	"time"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	"vorpalstacks/internal/utils/timeutils"
)

// --- Event API (v2) ---

// apiToMap converts an Api struct to a response map with correct wire format.
// Timestamps are serialised as epoch seconds (float64) per REST-JSON 1.0 protocol.
func apiToMap(api *appsyncstore.Api) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":       api.ApiId,
		"name":        api.Name,
		"apiArn":      api.Arn,
		"created":     timeutils.FormatEpochSeconds(api.Created),
		"xrayEnabled": api.XrayEnabled,
	}

	if api.OwnerContact != "" {
		m["ownerContact"] = api.OwnerContact
	}
	if api.WafWebAclArn != "" {
		m["wafWebAclArn"] = api.WafWebAclArn
	}
	if len(api.Dns) > 0 {
		m["dns"] = api.Dns
	}
	if api.EventConfig != nil {
		m["eventConfig"] = eventConfigToMap(api.EventConfig)
	}

	return m
}

// eventConfigToMap converts an EventConfig to a wire-format map.
func eventConfigToMap(ec *appsyncstore.EventConfig) map[string]interface{} {
	m := map[string]interface{}{}

	if len(ec.AuthProviders) > 0 {
		providers := make([]interface{}, 0, len(ec.AuthProviders))
		for _, ap := range ec.AuthProviders {
			providers = append(providers, authProviderToMap(&ap))
		}
		m["authProviders"] = providers
	}
	if len(ec.ConnectionAuthModes) > 0 {
		m["connectionAuthModes"] = authModesToMap(ec.ConnectionAuthModes)
	}
	if len(ec.DefaultPublishAuthModes) > 0 {
		m["defaultPublishAuthModes"] = authModesToMap(ec.DefaultPublishAuthModes)
	}
	if len(ec.DefaultSubscribeAuthModes) > 0 {
		m["defaultSubscribeAuthModes"] = authModesToMap(ec.DefaultSubscribeAuthModes)
	}
	if ec.LogConfig != nil {
		m["logConfig"] = map[string]interface{}{
			"cloudWatchLogsRoleArn": ec.LogConfig.CloudWatchLogsRoleArn,
			"logLevel":              ec.LogConfig.LogLevel,
		}
	}

	return m
}

// authProviderToMap converts an AuthProvider to a wire-format map.
func authProviderToMap(ap *appsyncstore.AuthProvider) map[string]interface{} {
	m := map[string]interface{}{
		"authType": ap.AuthType,
	}
	if ap.CognitoConfig != nil {
		m["cognitoConfig"] = map[string]interface{}{
			"awsRegion":        ap.CognitoConfig.AwsRegion,
			"userPoolId":       ap.CognitoConfig.UserPoolId,
			"appIdClientRegex": ap.CognitoConfig.AppIdClientRegex,
		}
	}
	if ap.LambdaAuthorizerConfig != nil {
		lc := map[string]interface{}{
			"authorizerUri": ap.LambdaAuthorizerConfig.AuthorizerUri,
		}
		if ap.LambdaAuthorizerConfig.AuthorizerResultTtlInSeconds > 0 {
			lc["authorizerResultTtlInSeconds"] = ap.LambdaAuthorizerConfig.AuthorizerResultTtlInSeconds
		}
		if ap.LambdaAuthorizerConfig.IdentityValidationExpression != "" {
			lc["identityValidationExpression"] = ap.LambdaAuthorizerConfig.IdentityValidationExpression
		}
		m["lambdaAuthorizerConfig"] = lc
	}
	if ap.OpenIDConnectConfig != nil {
		oidc := map[string]interface{}{
			"issuer": ap.OpenIDConnectConfig.Issuer,
		}
		if ap.OpenIDConnectConfig.AuthTTL > 0 {
			oidc["authTTL"] = ap.OpenIDConnectConfig.AuthTTL
		}
		if ap.OpenIDConnectConfig.ClientId != "" {
			oidc["clientId"] = ap.OpenIDConnectConfig.ClientId
		}
		if ap.OpenIDConnectConfig.IatTTL > 0 {
			oidc["iatTTL"] = ap.OpenIDConnectConfig.IatTTL
		}
		m["openIDConnectConfig"] = oidc
	}
	return m
}

// authModesToMap converts a slice of AuthMode to a wire-format slice.
func authModesToMap(modes []appsyncstore.AuthMode) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(modes))
	for _, mode := range modes {
		result = append(result, map[string]interface{}{
			"authType": mode.AuthType,
		})
	}
	return result
}

// --- Channel Namespace ---

// channelNamespaceToMap converts a ChannelNamespace struct to a response map
// with correct wire format. Timestamps are serialised as epoch seconds.
func channelNamespaceToMap(ns *appsyncstore.ChannelNamespace) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":               ns.ApiId,
		"name":                ns.Name,
		"channelNamespaceArn": ns.ChannelNamespaceArn,
		"created":             timeutils.FormatEpochSeconds(ns.Created),
		"lastModified":        timeutils.FormatEpochSeconds(ns.LastModified),
	}

	if ns.CodeHandlers != "" {
		m["codeHandlers"] = ns.CodeHandlers
	}
	if ns.HandlerConfigs != nil {
		m["handlerConfigs"] = handlerConfigsToMap(ns.HandlerConfigs)
	}
	if len(ns.PublishAuthModes) > 0 {
		m["publishAuthModes"] = authModesToMap(ns.PublishAuthModes)
	}
	if len(ns.SubscribeAuthModes) > 0 {
		m["subscribeAuthModes"] = authModesToMap(ns.SubscribeAuthModes)
	}

	return m
}

// handlerConfigsToMap converts HandlerConfigs to a wire-format map.
func handlerConfigsToMap(hc *appsyncstore.HandlerConfigs) map[string]interface{} {
	m := map[string]interface{}{}
	if hc.OnPublish != nil {
		m["onPublish"] = handlerConfigToMap(hc.OnPublish)
	}
	if hc.OnSubscribe != nil {
		m["onSubscribe"] = handlerConfigToMap(hc.OnSubscribe)
	}
	return m
}

// handlerConfigToMap converts a HandlerConfig to a wire-format map.
func handlerConfigToMap(hc *appsyncstore.HandlerConfig) map[string]interface{} {
	m := map[string]interface{}{
		"behavior": hc.Behavior,
	}
	if hc.Integration != nil {
		integration := map[string]interface{}{
			"dataSourceName": hc.Integration.DataSourceName,
		}
		if hc.Integration.LambdaConfig != nil {
			integration["lambdaConfig"] = map[string]interface{}{
				"invokeType": hc.Integration.LambdaConfig.InvokeType,
			}
		}
		m["integration"] = integration
	}
	return m
}

// --- GraphQL API ---

// graphqlApiToMap converts a GraphqlApi struct to a response map with correct wire format.
func graphqlApiToMap(api *appsyncstore.GraphqlApi) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":              api.ApiId,
		"name":               api.Name,
		"arn":                api.Arn,
		"authenticationType": api.AuthenticationType,
		"xrayEnabled":        api.XrayEnabled,
	}

	if api.ApiType != "" {
		m["apiType"] = api.ApiType
	}
	if len(api.AdditionalAuthenticationProviders) > 0 {
		providers := make([]interface{}, 0, len(api.AdditionalAuthenticationProviders))
		for _, p := range api.AdditionalAuthenticationProviders {
			providers = append(providers, additionalAuthProviderToMap(&p))
		}
		m["additionalAuthenticationProviders"] = providers
	}
	if len(api.Dns) > 0 {
		m["dns"] = api.Dns
	}
	if api.EnhancedMetricsConfig != nil {
		m["enhancedMetricsConfig"] = enhancedMetricsConfigToMap(api.EnhancedMetricsConfig)
	}
	if api.IntrospectionConfig != "" {
		m["introspectionConfig"] = api.IntrospectionConfig
	}
	if api.LambdaAuthorizerConfig != nil {
		m["lambdaAuthorizerConfig"] = lambdaAuthorizerConfigToMap(api.LambdaAuthorizerConfig)
	}
	if api.LogConfig != nil {
		m["logConfig"] = logConfigToMap(api.LogConfig)
	}
	if api.MergedApiExecutionRoleArn != "" {
		m["mergedApiExecutionRoleArn"] = api.MergedApiExecutionRoleArn
	}
	if api.OpenIDConnectConfig != nil {
		m["openIDConnectConfig"] = openIDConnectConfigToMap(api.OpenIDConnectConfig)
	}
	if api.Owner != "" {
		m["owner"] = api.Owner
	}
	if api.OwnerContact != "" {
		m["ownerContact"] = api.OwnerContact
	}
	if api.QueryDepthLimit > 0 {
		m["queryDepthLimit"] = api.QueryDepthLimit
	}
	if api.ResolverCountLimit > 0 {
		m["resolverCountLimit"] = api.ResolverCountLimit
	}
	if len(api.Uris) > 0 {
		m["uris"] = api.Uris
	}
	if api.UserPoolConfig != nil {
		m["userPoolConfig"] = userPoolConfigToMap(api.UserPoolConfig)
	}
	if api.Visibility != "" {
		m["visibility"] = api.Visibility
	}
	if api.WafWebAclArn != "" {
		m["wafWebAclArn"] = api.WafWebAclArn
	}

	return m
}

func additionalAuthProviderToMap(ap *appsyncstore.AdditionalAuthenticationProvider) map[string]interface{} {
	m := map[string]interface{}{
		"authenticationType": ap.AuthenticationType,
	}
	if ap.LambdaAuthorizerConfig != nil {
		m["lambdaAuthorizerConfig"] = lambdaAuthorizerConfigToMap(ap.LambdaAuthorizerConfig)
	}
	if ap.OpenIDConnectConfig != nil {
		m["openIDConnectConfig"] = openIDConnectConfigToMap(ap.OpenIDConnectConfig)
	}
	if ap.UserPoolConfig != nil {
		m["userPoolConfig"] = map[string]interface{}{
			"awsRegion":        ap.UserPoolConfig.AwsRegion,
			"userPoolId":       ap.UserPoolConfig.UserPoolId,
			"appIdClientRegex": ap.UserPoolConfig.AppIdClientRegex,
		}
	}
	return m
}

func enhancedMetricsConfigToMap(c *appsyncstore.EnhancedMetricsConfig) map[string]interface{} {
	return map[string]interface{}{
		"dataSourceLevelMetricsBehavior": c.DataSourceLevelMetricsBehavior,
		"operationLevelMetricsConfig":    c.OperationLevelMetricsConfig,
		"resolverLevelMetricsBehavior":   c.ResolverLevelMetricsBehavior,
	}
}

func lambdaAuthorizerConfigToMap(c *appsyncstore.LambdaAuthorizerConfig) map[string]interface{} {
	m := map[string]interface{}{
		"authorizerUri": c.AuthorizerUri,
	}
	if c.AuthorizerResultTtlInSeconds > 0 {
		m["authorizerResultTtlInSeconds"] = c.AuthorizerResultTtlInSeconds
	}
	if c.IdentityValidationExpression != "" {
		m["identityValidationExpression"] = c.IdentityValidationExpression
	}
	return m
}

func logConfigToMap(c *appsyncstore.LogConfig) map[string]interface{} {
	m := map[string]interface{}{
		"cloudWatchLogsRoleArn": c.CloudWatchLogsRoleArn,
		"fieldLogLevel":         c.FieldLogLevel,
	}
	if c.ExcludeVerboseContent {
		m["excludeVerboseContent"] = c.ExcludeVerboseContent
	}
	return m
}

func openIDConnectConfigToMap(c *appsyncstore.OpenIDConnectConfig) map[string]interface{} {
	m := map[string]interface{}{
		"issuer": c.Issuer,
	}
	if c.AuthTTL > 0 {
		m["authTTL"] = c.AuthTTL
	}
	if c.ClientId != "" {
		m["clientId"] = c.ClientId
	}
	if c.IatTTL > 0 {
		m["iatTTL"] = c.IatTTL
	}
	return m
}

func userPoolConfigToMap(c *appsyncstore.UserPoolConfig) map[string]interface{} {
	m := map[string]interface{}{
		"awsRegion":     c.AwsRegion,
		"defaultAction": c.DefaultAction,
		"userPoolId":    c.UserPoolId,
	}
	if c.AppIdClientRegex != "" {
		m["appIdClientRegex"] = c.AppIdClientRegex
	}
	return m
}

// --- Data Source ---

// dataSourceToMap converts a DataSource struct to a response map with correct wire format.
func dataSourceToMap(ds *appsyncstore.DataSource) map[string]interface{} {
	m := map[string]interface{}{
		"name":          ds.Name,
		"type":          ds.Type,
		"dataSourceArn": ds.DataSourceArn,
	}

	if ds.Description != "" {
		m["description"] = ds.Description
	}
	if ds.ServiceRoleArn != "" {
		m["serviceRoleArn"] = ds.ServiceRoleArn
	}
	if ds.DynamodbConfig != nil {
		m["dynamodbConfig"] = dynamodbConfigToMap(ds.DynamodbConfig)
	}
	if ds.ElasticsearchConfig != nil {
		m["elasticsearchConfig"] = elasticsearchConfigToMap(ds.ElasticsearchConfig)
	}
	if ds.EventBridgeConfig != nil {
		m["eventBridgeConfig"] = eventBridgeConfigToMap(ds.EventBridgeConfig)
	}
	if ds.HttpConfig != nil {
		m["httpConfig"] = httpConfigToMap(ds.HttpConfig)
	}
	if ds.LambdaConfig != nil {
		m["lambdaConfig"] = lambdaDataSourceConfigToMap(ds.LambdaConfig)
	}
	if ds.MetricsConfig != "" {
		m["metricsConfig"] = ds.MetricsConfig
	}
	if ds.NeptuneConfig != nil {
		m["neptuneConfig"] = neptuneConfigToMap(ds.NeptuneConfig)
	}
	if ds.OpenSearchServiceConfig != nil {
		m["openSearchServiceConfig"] = openSearchServiceConfigToMap(ds.OpenSearchServiceConfig)
	}
	if ds.RelationalDatabaseConfig != nil {
		m["relationalDatabaseConfig"] = relationalDatabaseConfigToMap(ds.RelationalDatabaseConfig)
	}
	if len(ds.Tags) > 0 {
		m["tags"] = ds.Tags
	}

	return m
}

// dynamodbConfigToMap converts a DynamodbDataSourceConfig to a wire-format map.
func dynamodbConfigToMap(c *appsyncstore.DynamodbDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"tableName": c.TableName,
	}
	if c.UseCallerCredentials {
		m["useCallerCredentials"] = c.UseCallerCredentials
	}
	if c.Versioned {
		m["versioned"] = c.Versioned
	}
	if c.DeltaSyncConfig != nil {
		dsc := map[string]interface{}{}
		if c.DeltaSyncConfig.BaseTableTTL > 0 {
			dsc["baseTableTTL"] = c.DeltaSyncConfig.BaseTableTTL
		}
		if c.DeltaSyncConfig.DeltaSyncTableName != "" {
			dsc["deltaSyncTableName"] = c.DeltaSyncConfig.DeltaSyncTableName
		}
		if c.DeltaSyncConfig.DeltaSyncTableTTL > 0 {
			dsc["deltaSyncTableTTL"] = c.DeltaSyncConfig.DeltaSyncTableTTL
		}
		m["deltaSyncConfig"] = dsc
	}
	return m
}

// elasticsearchConfigToMap converts an ElasticsearchDataSourceConfig to a wire-format map.
func elasticsearchConfigToMap(c *appsyncstore.ElasticsearchDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"endpoint":  c.Endpoint,
	}
}

// eventBridgeConfigToMap converts an EventBridgeDataSourceConfig to a wire-format map.
func eventBridgeConfigToMap(c *appsyncstore.EventBridgeDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"eventBusArn": c.EventBusArn,
	}
}

// httpConfigToMap converts an HttpDataSourceConfig to a wire-format map.
func httpConfigToMap(c *appsyncstore.HttpDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.Endpoint != "" {
		m["endpoint"] = c.Endpoint
	}
	if c.AuthorizationConfig != nil {
		ac := map[string]interface{}{
			"authorizationType": c.AuthorizationConfig.AuthorizationType,
		}
		if c.AuthorizationConfig.AwsIamConfig != nil {
			iam := map[string]interface{}{}
			if c.AuthorizationConfig.AwsIamConfig.SigningRegion != "" {
				iam["signingRegion"] = c.AuthorizationConfig.AwsIamConfig.SigningRegion
			}
			if c.AuthorizationConfig.AwsIamConfig.SigningServiceName != "" {
				iam["signingServiceName"] = c.AuthorizationConfig.AwsIamConfig.SigningServiceName
			}
			ac["awsIamConfig"] = iam
		}
		m["authorizationConfig"] = ac
	}
	return m
}

// lambdaDataSourceConfigToMap converts a LambdaDataSourceConfig to a wire-format map.
func lambdaDataSourceConfigToMap(c *appsyncstore.LambdaDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"lambdaFunctionArn": c.LambdaFunctionArn,
	}
}

// openSearchServiceConfigToMap converts an OpenSearchServiceDataSourceConfig to a wire-format map.
func openSearchServiceConfigToMap(c *appsyncstore.OpenSearchServiceDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"endpoint":  c.Endpoint,
	}
}

// relationalDatabaseConfigToMap converts a RelationalDatabaseDataSourceConfig to a wire-format map.
func relationalDatabaseConfigToMap(c *appsyncstore.RelationalDatabaseDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.RelationalDatabaseSourceType != "" {
		m["relationalDatabaseSourceType"] = c.RelationalDatabaseSourceType
	}
	if c.RdsHttpEndpointConfig != nil {
		rds := map[string]interface{}{}
		if c.RdsHttpEndpointConfig.AwsRegion != "" {
			rds["awsRegion"] = c.RdsHttpEndpointConfig.AwsRegion
		}
		if c.RdsHttpEndpointConfig.AwsSecretStoreArn != "" {
			rds["awsSecretStoreArn"] = c.RdsHttpEndpointConfig.AwsSecretStoreArn
		}
		if c.RdsHttpEndpointConfig.DatabaseName != "" {
			rds["databaseName"] = c.RdsHttpEndpointConfig.DatabaseName
		}
		if c.RdsHttpEndpointConfig.DbClusterIdentifier != "" {
			rds["dbClusterIdentifier"] = c.RdsHttpEndpointConfig.DbClusterIdentifier
		}
		if c.RdsHttpEndpointConfig.Schema != "" {
			rds["schema"] = c.RdsHttpEndpointConfig.Schema
		}
		m["rdsHttpEndpointConfig"] = rds
	}
	return m
}

// neptuneConfigToMap converts a NeptuneDataSourceConfig to a wire-format map.
func neptuneConfigToMap(c *appsyncstore.NeptuneDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.GraphID != "" {
		m["graphId"] = c.GraphID
	}
	return m
}

// --- Resolver ---

// resolverToMap converts a Resolver struct to a response map with correct wire format.
func resolverToMap(r *appsyncstore.Resolver) map[string]interface{} {
	m := map[string]interface{}{
		"typeName":  r.TypeName,
		"fieldName": r.FieldName,
	}

	if r.ResolverArn != "" {
		m["resolverArn"] = r.ResolverArn
	}
	if r.Kind != "" {
		m["kind"] = r.Kind
	}
	if r.DataSourceName != "" {
		m["dataSourceName"] = r.DataSourceName
	}
	if r.RequestMappingTemplate != "" {
		m["requestMappingTemplate"] = r.RequestMappingTemplate
	}
	if r.ResponseMappingTemplate != "" {
		m["responseMappingTemplate"] = r.ResponseMappingTemplate
	}
	if r.PipelineConfig != nil {
		m["pipelineConfig"] = pipelineConfigToMap(r.PipelineConfig)
	}
	if r.Runtime != nil {
		m["runtime"] = appSyncRuntimeToMap(r.Runtime)
	}
	if r.Code != "" {
		m["code"] = r.Code
	}
	if r.CachingConfig != nil {
		m["cachingConfig"] = cachingConfigToMap(r.CachingConfig)
	}
	if r.MaxBatchSize > 0 {
		m["maxBatchSize"] = r.MaxBatchSize
	}
	if r.MetricsConfig != "" {
		m["metricsConfig"] = r.MetricsConfig
	}
	if r.SyncConfig != nil {
		m["syncConfig"] = syncConfigToMap(r.SyncConfig)
	}

	return m
}

func appSyncRuntimeToMap(r *appsyncstore.AppSyncRuntime) map[string]interface{} {
	return map[string]interface{}{
		"name":           r.Name,
		"runtimeVersion": r.RuntimeVersion,
	}
}

func pipelineConfigToMap(c *appsyncstore.PipelineConfig) map[string]interface{} {
	functions := make([]interface{}, 0, len(c.Functions))
	for _, f := range c.Functions {
		functions = append(functions, f)
	}
	return map[string]interface{}{
		"functions": functions,
	}
}

func cachingConfigToMap(c *appsyncstore.CachingConfig) map[string]interface{} {
	m := map[string]interface{}{
		"ttl": c.Ttl,
	}
	if len(c.CachingKeys) > 0 {
		keys := make([]interface{}, 0, len(c.CachingKeys))
		for _, k := range c.CachingKeys {
			keys = append(keys, k)
		}
		m["cachingKeys"] = keys
	}
	return m
}

func syncConfigToMap(c *appsyncstore.SyncConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.ConflictDetection != "" {
		m["conflictDetection"] = c.ConflictDetection
	}
	if c.ConflictHandler != "" {
		m["conflictHandler"] = c.ConflictHandler
	}
	if c.LambdaConflictHandlerConfig != nil {
		m["lambdaConflictHandlerConfig"] = map[string]interface{}{
			"lambdaConflictHandlerArn": c.LambdaConflictHandlerConfig.LambdaConflictHandlerArn,
		}
	}
	return m
}

// --- Function ---

// functionToMap converts a FunctionConfiguration struct to a response map with correct wire format.
func functionToMap(f *appsyncstore.FunctionConfiguration) map[string]interface{} {
	m := map[string]interface{}{
		"functionId":     f.FunctionId,
		"name":           f.Name,
		"dataSourceName": f.DataSourceName,
	}

	if f.FunctionArn != "" {
		m["functionArn"] = f.FunctionArn
	}
	if f.FunctionVersion != "" {
		m["functionVersion"] = f.FunctionVersion
	}
	if f.Description != "" {
		m["description"] = f.Description
	}
	if f.RequestMappingTemplate != "" {
		m["requestMappingTemplate"] = f.RequestMappingTemplate
	}
	if f.ResponseMappingTemplate != "" {
		m["responseMappingTemplate"] = f.ResponseMappingTemplate
	}
	if f.Runtime != nil {
		m["runtime"] = appSyncRuntimeToMap(f.Runtime)
	}
	if f.Code != "" {
		m["code"] = f.Code
	}
	if f.MaxBatchSize > 0 {
		m["maxBatchSize"] = f.MaxBatchSize
	}
	if f.SyncConfig != nil {
		m["syncConfig"] = syncConfigToMap(f.SyncConfig)
	}

	return m
}

// --- Type ---

// typeToMap converts a Type struct to a response map with correct wire format.
func typeToMap(t *appsyncstore.Type) map[string]interface{} {
	m := map[string]interface{}{
		"name":   t.Name,
		"format": t.Format,
	}

	if t.Arn != "" {
		m["arn"] = t.Arn
	}
	if t.Definition != "" {
		m["definition"] = t.Definition
	}
	if t.Description != "" {
		m["description"] = t.Description
	}

	return m
}

// --- API Key ---

// apiKeyToMap converts an API key to a serialisable map.
// The deletes field is only emitted when non-zero, matching AWS behaviour.
func apiKeyToMap(k *appsyncstore.ApiKey) map[string]interface{} {
	result := map[string]interface{}{
		"id": k.Id,
	}
	if k.Description != "" {
		result["description"] = k.Description
	}
	if k.Expires != 0 {
		result["expires"] = k.Expires
	}
	if k.Deletes != 0 {
		result["deletes"] = k.Deletes
	}
	return result
}

// --- API Cache ---

func apiCacheToMap(c *appsyncstore.ApiCache) map[string]interface{} {
	result := map[string]interface{}{
		"type":               c.Type,
		"ttl":                c.Ttl,
		"apiCachingBehavior": c.ApiCachingBehavior,
		"status":             c.Status,
	}
	if c.AtRestEncryptionEnabled {
		result["atRestEncryptionEnabled"] = c.AtRestEncryptionEnabled
	}
	if c.TransitEncryptionEnabled {
		result["transitEncryptionEnabled"] = c.TransitEncryptionEnabled
	}
	if c.HealthMetricsConfig != "" {
		result["healthMetricsConfig"] = c.HealthMetricsConfig
	}
	return result
}

// --- Domain Name ---

func domainNameConfigToMap(c *appsyncstore.DomainNameConfig) map[string]interface{} {
	result := map[string]interface{}{
		"domainName":        c.DomainName,
		"appsyncDomainName": c.AppsyncDomainName,
	}
	if c.CertificateArn != "" {
		result["certificateArn"] = c.CertificateArn
	}
	if c.Description != "" {
		result["description"] = c.Description
	}
	if c.DomainNameArn != "" {
		result["domainNameArn"] = c.DomainNameArn
	}
	if c.HostedZoneId != "" {
		result["hostedZoneId"] = c.HostedZoneId
	}
	if len(c.Tags) > 0 {
		result["tags"] = c.Tags
	}
	return result
}

func apiAssociationToMap(a *appsyncstore.ApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"domainName":        a.DomainName,
		"associationStatus": a.AssociationStatus,
	}
	if a.ApiId != "" {
		result["apiId"] = a.ApiId
	}
	if a.DeploymentDetail != "" {
		result["deploymentDetail"] = a.DeploymentDetail
	}
	return result
}

// --- Merged API Association ---

func mergedApiAssociationToMap(a *appsyncstore.SourceApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"associationId":              a.AssociationId,
		"mergedApiId":                a.MergedApiId,
		"sourceApiId":                a.SourceApiId,
		"sourceApiAssociationStatus": a.SourceApiAssociationStatus,
	}
	if a.AssociationArn != "" {
		result["associationArn"] = a.AssociationArn
	}
	if a.MergedApiArn != "" {
		result["mergedApiArn"] = a.MergedApiArn
	}
	if a.SourceApiArn != "" {
		result["sourceApiArn"] = a.SourceApiArn
	}
	if a.Description != "" {
		result["description"] = a.Description
	}
	if a.SourceApiAssociationConfig != nil {
		result["sourceApiAssociationConfig"] = map[string]interface{}{
			"mergeType": a.SourceApiAssociationConfig.MergeType,
		}
	}
	if a.SourceApiAssociationStatusDetail != "" {
		result["sourceApiAssociationStatusDetail"] = a.SourceApiAssociationStatusDetail
	}
	if a.LastSuccessfulMergeDate != nil {
		result["lastSuccessfulMergeDate"] = a.LastSuccessfulMergeDate.Format(time.RFC3339)
	}
	return result
}

func mergedApiAssociationSummaryToMap(a *appsyncstore.SourceApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"associationId": a.AssociationId,
		"mergedApiId":   a.MergedApiId,
		"sourceApiId":   a.SourceApiId,
	}
	if a.AssociationArn != "" {
		result["associationArn"] = a.AssociationArn
	}
	if a.MergedApiArn != "" {
		result["mergedApiArn"] = a.MergedApiArn
	}
	if a.SourceApiArn != "" {
		result["sourceApiArn"] = a.SourceApiArn
	}
	if a.Description != "" {
		result["description"] = a.Description
	}
	return result
}
