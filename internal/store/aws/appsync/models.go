package appsync

import "time"

// Event API (v2) configuration types.

// EventConfig defines the authentication and connection configuration for an Event API.
type EventConfig struct {
	AuthProviders             []AuthProvider  `json:"authProviders"`
	ConnectionAuthModes       []AuthMode      `json:"connectionAuthModes"`
	DefaultPublishAuthModes   []AuthMode      `json:"defaultPublishAuthModes"`
	DefaultSubscribeAuthModes []AuthMode      `json:"defaultSubscribeAuthModes"`
	LogConfig                 *EventLogConfig `json:"logConfig,omitempty"`
}

// AuthProvider defines a single authentication provider within an EventConfig.
type AuthProvider struct {
	AuthType               string                  `json:"authType"`
	CognitoConfig          *CognitoEventConfig     `json:"cognitoConfig,omitempty"`
	LambdaAuthorizerConfig *LambdaAuthorizerConfig `json:"lambdaAuthorizerConfig,omitempty"`
	OpenIDConnectConfig    *OpenIDConnectConfig    `json:"openIDConnectConfig,omitempty"`
}

// AuthMode specifies an authentication type for publish or subscribe operations.
type AuthMode struct {
	AuthType string `json:"authType"`
}

// EventLogConfig configures CloudWatch Logs delivery for an Event API.
type EventLogConfig struct {
	CloudWatchLogsRoleArn string `json:"cloudWatchLogsRoleArn"`
	LogLevel              string `json:"logLevel"`
}

// CognitoEventConfig holds Amazon Cognito settings for Event API authentication.
// This is distinct from CognitoUserPoolConfig and UserPoolConfig used by GraphQL APIs.
type CognitoEventConfig struct {
	AwsRegion        string `json:"awsRegion"`
	UserPoolId       string `json:"userPoolId"`
	AppIdClientRegex string `json:"appIdClientRegex,omitempty"`
}

// LambdaAuthorizerConfig configures a Lambda function as an authoriser.
// Shared between Event API (v2) and GraphQL API (v1).
type LambdaAuthorizerConfig struct {
	AuthorizerUri                string `json:"authorizerUri"`
	AuthorizerResultTtlInSeconds int32  `json:"authorizerResultTtlInSeconds,omitempty"`
	IdentityValidationExpression string `json:"identityValidationExpression,omitempty"`
}

// OpenIDConnectConfig specifies an OpenID Connect (OIDC) identity provider.
// Shared between Event API (v2) and GraphQL API (v1).
type OpenIDConnectConfig struct {
	Issuer   string `json:"issuer"`
	AuthTTL  int64  `json:"authTTL,omitempty"`
	ClientId string `json:"clientId,omitempty"`
	IatTTL   int64  `json:"iatTTL,omitempty"`
}

// Api represents an Event API (v2) resource.
type Api struct {
	ApiId        string            `json:"apiId"`
	Name         string            `json:"name"`
	Arn          string            `json:"apiArn"`
	Created      time.Time         `json:"created"`
	Dns          map[string]string `json:"dns,omitempty"`
	EventConfig  *EventConfig      `json:"eventConfig"`
	OwnerContact string            `json:"ownerContact,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	WafWebAclArn string            `json:"wafWebAclArn,omitempty"`
	XrayEnabled  bool              `json:"xrayEnabled"`
}

// ChannelNamespace represents a channel namespace within an Event API (v2).
// Channel paths are scoped to a namespace (e.g. /myNamespace/ch1/ch2).
type ChannelNamespace struct {
	ApiId               string            `json:"apiId"`
	Name                string            `json:"name"`
	ChannelNamespaceArn string            `json:"channelNamespaceArn"`
	CodeHandlers        string            `json:"codeHandlers,omitempty"`
	Created             time.Time         `json:"created"`
	HandlerConfigs      *HandlerConfigs   `json:"handlerConfigs,omitempty"`
	LastModified        time.Time         `json:"lastModified"`
	PublishAuthModes    []AuthMode        `json:"publishAuthModes,omitempty"`
	SubscribeAuthModes  []AuthMode        `json:"subscribeAuthModes,omitempty"`
	Tags                map[string]string `json:"tags,omitempty"`
}

// HandlerConfigs defines optional handlers for publish and subscribe events
// within a channel namespace.
type HandlerConfigs struct {
	OnPublish   *HandlerConfig `json:"onPublish,omitempty"`
	OnSubscribe *HandlerConfig `json:"onSubscribe,omitempty"`
}

// HandlerConfig defines how an event (publish or subscribe) is processed.
type HandlerConfig struct {
	Behavior    string       `json:"behavior"`
	Integration *Integration `json:"integration"`
}

// Integration specifies the data source that handles events for a handler.
type Integration struct {
	DataSourceName string           `json:"dataSourceName"`
	LambdaConfig   *LambdaIntConfig `json:"lambdaConfig,omitempty"`
}

// LambdaIntConfig specifies invocation settings for a Lambda-based handler integration.
type LambdaIntConfig struct {
	InvokeType string `json:"invokeType,omitempty"`
}

// GraphQL API (v1) types.

// GraphqlApi represents a GraphQL API (v1) resource.
type GraphqlApi struct {
	ApiId                             string                             `json:"apiId"`
	Name                              string                             `json:"name"`
	Arn                               string                             `json:"arn"`
	ApiType                           string                             `json:"apiType,omitempty"`
	AuthenticationType                string                             `json:"authenticationType"`
	AdditionalAuthenticationProviders []AdditionalAuthenticationProvider `json:"additionalAuthenticationProviders,omitempty"`
	Dns                               map[string]string                  `json:"dns,omitempty"`
	EnhancedMetricsConfig             *EnhancedMetricsConfig             `json:"enhancedMetricsConfig,omitempty"`
	IntrospectionConfig               string                             `json:"introspectionConfig,omitempty"`
	LambdaAuthorizerConfig            *LambdaAuthorizerConfig            `json:"lambdaAuthorizerConfig,omitempty"`
	LogConfig                         *LogConfig                         `json:"logConfig,omitempty"`
	MergedApiExecutionRoleArn         string                             `json:"mergedApiExecutionRoleArn,omitempty"`
	OpenIDConnectConfig               *OpenIDConnectConfig               `json:"openIDConnectConfig,omitempty"`
	Owner                             string                             `json:"owner,omitempty"`
	OwnerContact                      string                             `json:"ownerContact,omitempty"`
	QueryDepthLimit                   int32                              `json:"queryDepthLimit,omitempty"`
	ResolverCountLimit                int32                              `json:"resolverCountLimit,omitempty"`
	Tags                              map[string]string                  `json:"tags,omitempty"`
	Uris                              map[string]string                  `json:"uris,omitempty"`
	UserPoolConfig                    *UserPoolConfig                    `json:"userPoolConfig,omitempty"`
	Visibility                        string                             `json:"visibility,omitempty"`
	WafWebAclArn                      string                             `json:"wafWebAclArn,omitempty"`
	XrayEnabled                       bool                               `json:"xrayEnabled"`
}

// AdditionalAuthenticationProvider defines an extra authentication mechanism
// beyond the primary authenticationType on a GraphQL API.
type AdditionalAuthenticationProvider struct {
	AuthenticationType     string                  `json:"authenticationType"`
	LambdaAuthorizerConfig *LambdaAuthorizerConfig `json:"lambdaAuthorizerConfig,omitempty"`
	OpenIDConnectConfig    *OpenIDConnectConfig    `json:"openIDConnectConfig,omitempty"`
	UserPoolConfig         *CognitoUserPoolConfig  `json:"userPoolConfig,omitempty"`
}

// CognitoUserPoolConfig configures Cognito user pool authentication for an
// additional authentication provider. Differs from CognitoEventConfig (Event API)
// and UserPoolConfig (primary GraphQL auth with DefaultAction).
type CognitoUserPoolConfig struct {
	AwsRegion        string `json:"awsRegion"`
	UserPoolId       string `json:"userPoolId"`
	AppIdClientRegex string `json:"appIdClientRegex,omitempty"`
}

// UserPoolConfig configures Amazon Cognito user pool as the primary authentication
// for a GraphQL API. Includes DefaultAction (ALLOW/DENY) which the other two
// Cognito config types do not have.
type UserPoolConfig struct {
	AwsRegion        string `json:"awsRegion"`
	DefaultAction    string `json:"defaultAction"`
	UserPoolId       string `json:"userPoolId"`
	AppIdClientRegex string `json:"appIdClientRegex,omitempty"`
}

// EnhancedMetricsConfig controls enhanced metric collection levels.
type EnhancedMetricsConfig struct {
	DataSourceLevelMetricsBehavior string `json:"dataSourceLevelMetricsBehavior"`
	OperationLevelMetricsConfig    string `json:"operationLevelMetricsConfig"`
	ResolverLevelMetricsBehavior   string `json:"resolverLevelMetricsBehavior"`
}

// LogConfig configures CloudWatch Logs for a GraphQL API.
type LogConfig struct {
	CloudWatchLogsRoleArn string `json:"cloudWatchLogsRoleArn"`
	FieldLogLevel         string `json:"fieldLogLevel"`
	ExcludeVerboseContent bool   `json:"excludeVerboseContent,omitempty"`
}

// DataSource types.

// DataSource represents a data source connected to a GraphQL API.
// The ApiId field is excluded from JSON serialisation (populated from URI).
type DataSource struct {
	ApiId                    string                              `json:"-"`
	Name                     string                              `json:"name"`
	Type                     string                              `json:"type"`
	DataSourceArn            string                              `json:"dataSourceArn"`
	Description              string                              `json:"description,omitempty"`
	ServiceRoleArn           string                              `json:"serviceRoleArn,omitempty"`
	DynamodbConfig           *DynamodbDataSourceConfig           `json:"dynamodbConfig,omitempty"`
	ElasticsearchConfig      *ElasticsearchDataSourceConfig      `json:"elasticsearchConfig,omitempty"`
	EventBridgeConfig        *EventBridgeDataSourceConfig        `json:"eventBridgeConfig,omitempty"`
	HttpConfig               *HttpDataSourceConfig               `json:"httpConfig,omitempty"`
	LambdaConfig             *LambdaDataSourceConfig             `json:"lambdaConfig,omitempty"`
	MetricsConfig            string                              `json:"metricsConfig,omitempty"`
	OpenSearchServiceConfig  *OpenSearchServiceDataSourceConfig  `json:"openSearchServiceConfig,omitempty"`
	RelationalDatabaseConfig *RelationalDatabaseDataSourceConfig `json:"relationalDatabaseConfig,omitempty"`
	Tags                     map[string]string                   `json:"tags,omitempty"`
}

// DynamodbDataSourceConfig specifies an Amazon DynamoDB table as a data source.
type DynamodbDataSourceConfig struct {
	AwsRegion            string           `json:"awsRegion"`
	TableName            string           `json:"tableName"`
	UseCallerCredentials bool             `json:"useCallerCredentials,omitempty"`
	Versioned            bool             `json:"versioned,omitempty"`
	DeltaSyncConfig      *DeltaSyncConfig `json:"deltaSyncConfig,omitempty"`
}

// DeltaSyncConfig configures Delta Sync for a DynamoDB data source.
type DeltaSyncConfig struct {
	BaseTableTTL       int64  `json:"baseTableTTL,omitempty"`
	DeltaSyncTableName string `json:"deltaSyncTableName,omitempty"`
	DeltaSyncTableTTL  int64  `json:"deltaSyncTableTTL,omitempty"`
}

// LambdaDataSourceConfig specifies an AWS Lambda function as a data source.
type LambdaDataSourceConfig struct {
	LambdaFunctionArn string `json:"lambdaFunctionArn"`
}

// HttpDataSourceConfig specifies an HTTP endpoint as a data source.
type HttpDataSourceConfig struct {
	Endpoint            string               `json:"endpoint,omitempty"`
	AuthorizationConfig *AuthorizationConfig `json:"authorizationConfig,omitempty"`
}

// AuthorizationConfig specifies how an HTTP data source authenticates requests.
type AuthorizationConfig struct {
	AuthorizationType string        `json:"authorizationType"`
	AwsIamConfig      *AwsIamConfig `json:"awsIamConfig,omitempty"`
}

// AwsIamConfig configures AWS IAM signing for HTTP data source requests.
type AwsIamConfig struct {
	SigningRegion      string `json:"signingRegion,omitempty"`
	SigningServiceName string `json:"signingServiceName,omitempty"`
}

// ElasticsearchDataSourceConfig specifies an Amazon Elasticsearch domain as a data source.
type ElasticsearchDataSourceConfig struct {
	AwsRegion string `json:"awsRegion"`
	Endpoint  string `json:"endpoint"`
}

// OpenSearchServiceDataSourceConfig specifies an Amazon OpenSearch Service domain as a data source.
type OpenSearchServiceDataSourceConfig struct {
	AwsRegion string `json:"awsRegion"`
	Endpoint  string `json:"endpoint"`
}

// EventBridgeDataSourceConfig specifies an Amazon EventBridge event bus as a data source.
type EventBridgeDataSourceConfig struct {
	EventBusArn string `json:"eventBusArn"`
}

// RelationalDatabaseDataSourceConfig specifies an RDS or Aurora relational database as a data source.
type RelationalDatabaseDataSourceConfig struct {
	RelationalDatabaseSourceType string                 `json:"relationalDatabaseSourceType,omitempty"`
	RdsHttpEndpointConfig        *RdsHttpEndpointConfig `json:"rdsHttpEndpointConfig,omitempty"`
}

// RdsHttpEndpointConfig configures the RDS Data API HTTP endpoint for a relational database.
type RdsHttpEndpointConfig struct {
	AwsRegion           string `json:"awsRegion,omitempty"`
	AwsSecretStoreArn   string `json:"awsSecretStoreArn,omitempty"`
	DatabaseName        string `json:"databaseName,omitempty"`
	DbClusterIdentifier string `json:"dbClusterIdentifier,omitempty"`
	Schema              string `json:"schema,omitempty"`
}

// Resolver types.

// Resolver maps a GraphQL field to a data source via VTL or APPSYNC_JS templates.
// ApiId is excluded from JSON serialisation (populated from URI).
type Resolver struct {
	ApiId                   string          `json:"-"`
	TypeName                string          `json:"typeName"`
	FieldName               string          `json:"fieldName"`
	ResolverArn             string          `json:"resolverArn,omitempty"`
	Kind                    string          `json:"kind,omitempty"`
	DataSourceName          string          `json:"dataSourceName,omitempty"`
	RequestMappingTemplate  string          `json:"requestMappingTemplate,omitempty"`
	RequestMappingTemplateS string          `json:"requestMappingTemplateS,omitempty"`
	ResponseMappingTemplate string          `json:"responseMappingTemplate,omitempty"`
	PipelineConfig          *PipelineConfig `json:"pipelineConfig,omitempty"`
	Runtime                 *AppSyncRuntime `json:"runtime,omitempty"`
	Code                    string          `json:"code,omitempty"`
	CachingConfig           *CachingConfig  `json:"cachingConfig,omitempty"`
	MaxBatchSize            int32           `json:"maxBatchSize,omitempty"`
	MetricsConfig           string          `json:"metricsConfig,omitempty"`
	SyncConfig              *SyncConfig     `json:"syncConfig,omitempty"`
}

// PipelineConfig lists the function IDs that form a pipeline resolver's execution chain.
type PipelineConfig struct {
	Functions []string `json:"functions"`
}

// AppSyncRuntime specifies the resolver runtime (e.g. APPSYNC_JS).
type AppSyncRuntime struct {
	Name           string `json:"name"`
	RuntimeVersion string `json:"runtimeVersion"`
}

// CachingConfig controls per-resolver caching behaviour.
type CachingConfig struct {
	CachingKeys []string `json:"cachingKeys,omitempty"`
	Ttl         int64    `json:"ttl"`
}

// SyncConfig configures conflict detection and resolution for a resolver.
type SyncConfig struct {
	ConflictDetection           string                       `json:"conflictDetection,omitempty"`
	ConflictHandler             string                       `json:"conflictHandler,omitempty"`
	LambdaConflictHandlerConfig *LambdaConflictHandlerConfig `json:"lambdaConflictHandlerConfig,omitempty"`
}

// LambdaConflictHandlerConfig specifies the Lambda function used for conflict resolution.
type LambdaConflictHandlerConfig struct {
	LambdaConflictHandlerArn string `json:"lambdaConflictHandlerArn"`
}

// FunctionConfiguration represents an AppSync Function (a reusable resolver unit
// within a pipeline resolver). Not to be confused with an AWS Lambda function.
// ApiId is excluded from JSON serialisation (populated from URI).
type FunctionConfiguration struct {
	ApiId                   string          `json:"-"`
	FunctionId              string          `json:"functionId"`
	Name                    string          `json:"name"`
	FunctionArn             string          `json:"functionArn,omitempty"`
	FunctionVersion         string          `json:"functionVersion,omitempty"`
	Description             string          `json:"description,omitempty"`
	DataSourceName          string          `json:"dataSourceName"`
	RequestMappingTemplate  string          `json:"requestMappingTemplate,omitempty"`
	ResponseMappingTemplate string          `json:"responseMappingTemplate,omitempty"`
	Runtime                 *AppSyncRuntime `json:"runtime,omitempty"`
	Code                    string          `json:"code,omitempty"`
	MaxBatchSize            int32           `json:"maxBatchSize,omitempty"`
	SyncConfig              *SyncConfig     `json:"syncConfig,omitempty"`
}

// Type represents a GraphQL type definition attached to an API.
// ApiId is excluded from JSON serialisation (populated from URI).
type Type struct {
	ApiId       string `json:"-"`
	Name        string `json:"name"`
	Arn         string `json:"arn,omitempty"`
	Definition  string `json:"definition,omitempty"`
	Description string `json:"description,omitempty"`
	Format      string `json:"format"`
}

// ApiKey represents an API key for GraphQL API authentication.
type ApiKey struct {
	Id          string `json:"id"`
	Description string `json:"description,omitempty"`
	Expires     int64  `json:"expires,omitempty"`
	Deletes     int64  `json:"deletes,omitempty"`
}

// ApiCache represents caching configuration for a GraphQL API.
type ApiCache struct {
	Type                     string `json:"type"`
	Ttl                      int64  `json:"ttl"`
	ApiCachingBehavior       string `json:"apiCachingBehavior"`
	AtRestEncryptionEnabled  bool   `json:"atRestEncryptionEnabled,omitempty"`
	TransitEncryptionEnabled bool   `json:"transitEncryptionEnabled,omitempty"`
	Status                   string `json:"status,omitempty"`
	HealthMetricsConfig      string `json:"healthMetricsConfig,omitempty"`
}

// Domain name and API association types.

// DomainNameConfig represents a custom domain name for an AppSync API.
type DomainNameConfig struct {
	DomainName        string            `json:"domainName"`
	AppsyncDomainName string            `json:"appsyncDomainName,omitempty"`
	CertificateArn    string            `json:"certificateArn,omitempty"`
	Description       string            `json:"description,omitempty"`
	DomainNameArn     string            `json:"domainNameArn,omitempty"`
	HostedZoneId      string            `json:"hostedZoneId,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// ApiAssociation links a domain name to a GraphQL API.
type ApiAssociation struct {
	ApiId             string `json:"apiId,omitempty"`
	DomainName        string `json:"domainName,omitempty"`
	AssociationStatus string `json:"associationStatus,omitempty"`
	DeploymentDetail  string `json:"deploymentDetail,omitempty"`
}

// Schema and environment variable types.

// SchemaCreationStatus tracks the status of an async schema creation operation.
// ApiId is excluded from JSON serialisation (populated from URI).
type SchemaCreationStatus struct {
	ApiId      string `json:"-"`
	Status     string `json:"status"`
	Details    string `json:"details,omitempty"`
	Definition string `json:"definition,omitempty"`
}

// EnvironmentVariables stores key-value pairs accessible in resolver templates via $stageVariables.
// ApiId is excluded from JSON serialisation (populated from URI).
type EnvironmentVariables struct {
	ApiId                string            `json:"-"`
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

// Merged API types.

// SourceApiAssociation represents a source API linked to a merged API.
type SourceApiAssociation struct {
	AssociationId                    string                      `json:"associationId"`
	MergedApiId                      string                      `json:"mergedApiId"`
	SourceApiId                      string                      `json:"sourceApiId"`
	AssociationArn                   string                      `json:"associationArn,omitempty"`
	MergedApiArn                     string                      `json:"mergedApiArn,omitempty"`
	SourceApiArn                     string                      `json:"sourceApiArn,omitempty"`
	Description                      string                      `json:"description,omitempty"`
	SourceApiAssociationStatus       string                      `json:"sourceApiAssociationStatus,omitempty"`
	SourceApiAssociationStatusDetail string                      `json:"sourceApiAssociationStatusDetail,omitempty"`
	SourceApiAssociationConfig       *SourceApiAssociationConfig `json:"sourceApiAssociationConfig,omitempty"`
	LastSuccessfulMergeDate          *time.Time                  `json:"lastSuccessfulMergeDate,omitempty"`
}

// SourceApiAssociationSummary is a lightweight representation used in list operations.
type SourceApiAssociationSummary struct {
	AssociationId  string `json:"associationId"`
	MergedApiId    string `json:"mergedApiId"`
	SourceApiId    string `json:"sourceApiId"`
	AssociationArn string `json:"associationArn,omitempty"`
	MergedApiArn   string `json:"mergedApiArn,omitempty"`
	SourceApiArn   string `json:"sourceApiArn,omitempty"`
	Description    string `json:"description,omitempty"`
}

// SourceApiAssociationConfig controls merge behaviour (MANUAL_MERGE or AUTO_MERGE).
type SourceApiAssociationConfig struct {
	MergeType string `json:"mergeType"`
}

// DataSource introspection types.

// DataSourceIntrospectionResult holds the status and output of an introspection run.
type DataSourceIntrospectionResult struct {
	IntrospectionId           string                              `json:"introspectionId,omitempty"`
	IntrospectionStatus       string                              `json:"introspectionStatus,omitempty"`
	IntrospectionStatusDetail string                              `json:"introspectionStatusDetail,omitempty"`
	IntrospectionResult       *DataSourceIntrospectionModelResult `json:"introspectionResult,omitempty"`
}

// DataSourceIntrospectionModelResult contains the discovered models from introspection.
type DataSourceIntrospectionModelResult struct {
	Models    []DataSourceIntrospectionModel `json:"models,omitempty"`
	NextToken string                         `json:"nextToken,omitempty"`
}

// DataSourceIntrospectionModel represents a single database table or view discovered by introspection.
type DataSourceIntrospectionModel struct {
	Name       string                              `json:"name,omitempty"`
	Fields     []DataSourceIntrospectionModelField `json:"fields,omitempty"`
	Indexes    []DataSourceIntrospectionModelIndex `json:"indexes,omitempty"`
	PrimaryKey *DataSourceIntrospectionModelIndex  `json:"primaryKey,omitempty"`
	Sdl        string                              `json:"sdl,omitempty"`
}

// DataSourceIntrospectionModelField describes a column within an introspected model.
type DataSourceIntrospectionModelField struct {
	Name   string                                 `json:"name,omitempty"`
	Length int64                                  `json:"length,omitempty"`
	Type   *DataSourceIntrospectionModelFieldType `json:"type,omitempty"`
}

// DataSourceIntrospectionModelFieldType describes the type of an introspected field.
// Self-referential: Type field is present when Kind is NonNull or List.
type DataSourceIntrospectionModelFieldType struct {
	Kind   string                                 `json:"kind,omitempty"`
	Name   string                                 `json:"name,omitempty"`
	Type   *DataSourceIntrospectionModelFieldType `json:"type,omitempty"`
	Values []string                               `json:"values,omitempty"`
}

// DataSourceIntrospectionModelIndex represents a database index discovered during introspection.
type DataSourceIntrospectionModelIndex struct {
	Name   string   `json:"name,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

// RdsDataApiConfig specifies the RDS Data API connection for a data source introspection request.
type RdsDataApiConfig struct {
	DatabaseName string `json:"databaseName"`
	ResourceArn  string `json:"resourceArn"`
	SecretArn    string `json:"secretArn"`
}
