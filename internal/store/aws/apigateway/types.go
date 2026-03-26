// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// RestApi represents an API Gateway REST API.
type RestApi struct {
	Id                     string                       `json:"id"`
	Name                   string                       `json:"name"`
	Description            string                       `json:"description,omitempty"`
	CreatedDate            time.Time                    `json:"created_date"`
	Version                string                       `json:"version,omitempty"`
	Warnings               []string                     `json:"warnings,omitempty"`
	BinaryMediaTypes       []string                     `json:"binary_media_types,omitempty"`
	MinimumCompressionSize int32                        `json:"minimum_compression_size,omitempty"`
	ApiKeySource           string                       `json:"api_key_source,omitempty"`
	EndpointConfiguration  *EndpointConfiguration       `json:"endpoint_configuration,omitempty"`
	Policy                 string                       `json:"policy,omitempty"`
	Tags                   []common.Tag                 `json:"tags,omitempty"`
	Resources              map[string]*Resource         `json:"resources,omitempty"`
	Deployments            map[string]*Deployment       `json:"deployments,omitempty"`
	Stages                 map[string]*Stage            `json:"stages,omitempty"`
	RequestValidators      map[string]*RequestValidator `json:"request_validators,omitempty"`
	Models                 map[string]*Model            `json:"models,omitempty"`
	Authorizers            map[string]*Authorizer       `json:"authorizers,omitempty"`
}

// EndpointConfiguration defines the endpoint configuration for an API.
type EndpointConfiguration struct {
	Types          []string `json:"types,omitempty"`
	VpcEndpointIds []string `json:"vpc_endpoint_ids,omitempty"`
}

// Resource represents an API Gateway resource.
type Resource struct {
	Id              string             `json:"id"`
	RestApiId       string             `json:"rest_api_id"`
	ParentId        string             `json:"parent_id,omitempty"`
	Path            string             `json:"path"`
	PathPart        string             `json:"path_part,omitempty"`
	ResourceMethods map[string]*Method `json:"resource_methods,omitempty"`
}

// Method represents an API Gateway method.
type Method struct {
	RestApiId          string                     `json:"rest_api_id"`
	ResourceId         string                     `json:"resource_id"`
	HttpMethod         string                     `json:"http_method"`
	AuthorizationType  string                     `json:"authorization_type,omitempty"`
	AuthorizerId       string                     `json:"authorizer_id,omitempty"`
	ApiKeyRequired     bool                       `json:"api_key_required"`
	RequestValidatorId string                     `json:"request_validator_id,omitempty"`
	OperationName      string                     `json:"operation_name,omitempty"`
	RequestParameters  map[string]bool            `json:"request_parameters,omitempty"`
	RequestModels      map[string]string          `json:"request_models,omitempty"`
	MethodIntegration  *Integration               `json:"method_integration,omitempty"`
	MethodResponses    map[string]*MethodResponse `json:"method_responses,omitempty"`
}

// Integration represents an API Gateway integration.
type Integration struct {
	RestApiId             string                          `json:"rest_api_id"`
	ResourceId            string                          `json:"resource_id"`
	HttpMethod            string                          `json:"http_method"`
	Type                  string                          `json:"type"`
	IntegrationHttpMethod string                          `json:"integration_http_method,omitempty"`
	Uri                   string                          `json:"uri,omitempty"`
	Credentials           string                          `json:"credentials,omitempty"`
	RequestParameters     map[string]string               `json:"request_parameters,omitempty"`
	RequestTemplates      map[string]string               `json:"request_templates,omitempty"`
	PassthroughBehavior   string                          `json:"passthrough_behavior,omitempty"`
	ContentHandling       string                          `json:"content_handling,omitempty"`
	CacheNamespace        string                          `json:"cache_namespace,omitempty"`
	CacheKeyParameters    []string                        `json:"cache_key_parameters,omitempty"`
	IntegrationResponses  map[string]*IntegrationResponse `json:"integration_responses,omitempty"`
	TimeoutInMillis       int32                           `json:"timeout_in_millis,omitempty"`
	ConnectionType        string                          `json:"connection_type,omitempty"`
	ConnectionId          string                          `json:"connection_id,omitempty"`
	TlsConfig             *TlsConfig                      `json:"tls_config,omitempty"`
}

// IntegrationResponse represents an API Gateway integration response.
type IntegrationResponse struct {
	StatusCode         string            `json:"status_code"`
	SelectionPattern   string            `json:"selection_pattern,omitempty"`
	ResponseParameters map[string]string `json:"response_parameters,omitempty"`
	ResponseTemplates  map[string]string `json:"response_templates,omitempty"`
	ContentHandling    string            `json:"content_handling,omitempty"`
}

// MethodResponse represents an API Gateway method response.
type MethodResponse struct {
	StatusCode         string            `json:"status_code"`
	ResponseParameters map[string]bool   `json:"response_parameters,omitempty"`
	ResponseModels     map[string]string `json:"response_models,omitempty"`
}

// TlsConfig defines the TLS configuration for an integration.
type TlsConfig struct {
	InsecureSkipVerification bool `json:"insecure_skip_verification"`
}

// Deployment represents an API Gateway deployment.
type Deployment struct {
	Id          string                 `json:"id"`
	RestApiId   string                 `json:"rest_api_id"`
	Description string                 `json:"description,omitempty"`
	CreatedDate time.Time              `json:"created_date"`
	ApiSummary  map[string]interface{} `json:"api_summary,omitempty"`
}

// Stage represents an API Gateway stage.
type Stage struct {
	DeploymentId         string                    `json:"deployment_id,omitempty"`
	RestApiId            string                    `json:"rest_api_id"`
	StageName            string                    `json:"stage_name"`
	Description          string                    `json:"description,omitempty"`
	CacheClusterEnabled  bool                      `json:"cache_cluster_enabled"`
	CacheClusterSize     string                    `json:"cache_cluster_size,omitempty"`
	CacheClusterStatus   string                    `json:"cache_cluster_status,omitempty"`
	MethodSettings       map[string]*MethodSetting `json:"method_settings,omitempty"`
	Variables            map[string]string         `json:"variables,omitempty"`
	DocumentationVersion string                    `json:"documentation_version,omitempty"`
	AccessLogSettings    *AccessLogSettings        `json:"access_log_settings,omitempty"`
	CanarySettings       *CanarySettings           `json:"canary_settings,omitempty"`
	TracingEnabled       bool                      `json:"tracing_enabled"`
	CreatedDate          time.Time                 `json:"created_date"`
	LastUpdatedDate      time.Time                 `json:"last_updated_date"`
	WebAclArn            string                    `json:"web_acl_arn,omitempty"`
	Tags                 []common.Tag              `json:"tags,omitempty"`
}

// MethodSetting defines the method-level settings for caching and throttling.
type MethodSetting struct {
	MetricsEnabled                      bool              `json:"metrics_enabled"`
	LoggingLevel                        string            `json:"logging_level,omitempty"`
	DataTraceEnabled                    bool              `json:"data_trace_enabled"`
	ThrottlingBurstLimit                int32             `json:"throttling_burst_limit,omitempty"`
	ThrottlingRateLimit                 float64           `json:"throttling_rate_limit,omitempty"`
	CachingEnabled                      bool              `json:"caching_enabled"`
	CacheTtlInSeconds                   int32             `json:"cache_ttl_in_seconds,omitempty"`
	CacheDataEncrypted                  bool              `json:"cache_data_encrypted"`
	RequireAuthorizationForCacheControl bool              `json:"require_authorization_for_cache_control"`
	UnreservedCacheParameters           map[string]string `json:"unreserved_cache_parameters,omitempty"`
}

// AccessLogSettings defines the access log settings for a stage.
type AccessLogSettings struct {
	DestinationArn string `json:"destination_arn,omitempty"`
	Format         string `json:"format,omitempty"`
}

// CanarySettings defines the canary deployment settings for a stage.
type CanarySettings struct {
	PercentTraffic         float64           `json:"percent_traffic,omitempty"`
	DeploymentId           string            `json:"deployment_id,omitempty"`
	StageVariableOverrides map[string]string `json:"stage_variable_overrides,omitempty"`
	UseStageCache          bool              `json:"use_stage_cache"`
}

// RequestValidator defines the request validation settings for an API.
type RequestValidator struct {
	Id                        string `json:"id"`
	RestApiId                 string `json:"rest_api_id"`
	Name                      string `json:"name"`
	ValidateRequestBody       bool   `json:"validate_request_body"`
	ValidateRequestParameters bool   `json:"validate_request_parameters"`
}

// Model defines a request or response model for an API.
type Model struct {
	Id          string `json:"id"`
	RestApiId   string `json:"rest_api_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
	ContentType string `json:"content_type,omitempty"`
}

// Authorizer defines an authoriser for an API.
type Authorizer struct {
	Id                           string   `json:"id"`
	RestApiId                    string   `json:"rest_api_id"`
	Name                         string   `json:"name"`
	Type                         string   `json:"type"`
	ProviderArns                 []string `json:"provider_arns,omitempty"`
	AuthType                     string   `json:"auth_type,omitempty"`
	AuthorizerUri                string   `json:"authorizer_uri,omitempty"`
	AuthorizerCredentials        string   `json:"authorizer_credentials,omitempty"`
	IdentitySource               string   `json:"identity_source,omitempty"`
	IdentityValidationExpression string   `json:"identity_validation_expression,omitempty"`
	AuthorizerResultTtlInSeconds int32    `json:"authorizer_result_ttl_in_seconds"`
}

// ApiKey represents an API key for API Gateway.
type ApiKey struct {
	Id              string       `json:"id"`
	Value           string       `json:"value"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	Enabled         bool         `json:"enabled"`
	CreatedDate     time.Time    `json:"created_date"`
	LastUpdatedDate time.Time    `json:"last_updated_date"`
	StageKeys       []string     `json:"stage_keys,omitempty"`
	Tags            []common.Tag `json:"tags,omitempty"`
	CustomerId      string       `json:"customer_id,omitempty"`
}

// UsagePlan defines a usage plan for API keys.
type UsagePlan struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	ApiStages   []ApiStage   `json:"api_stages,omitempty"`
	Quota       *Quota       `json:"quota,omitempty"`
	Throttle    *Throttle    `json:"throttle,omitempty"`
	Tags        []common.Tag `json:"tags,omitempty"`
	ProductCode string       `json:"product_code,omitempty"`
}

// ApiStage defines an API stage within a usage plan.
type ApiStage struct {
	ApiId    string               `json:"api_id"`
	Stage    string               `json:"stage"`
	Throttle map[string]*Throttle `json:"throttle,omitempty"`
}

// Quota defines the quota settings for a usage plan.
type Quota struct {
	Limit  int64  `json:"limit,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Period string `json:"period,omitempty"`
}

// Throttle defines the throttling settings for a usage plan.
type Throttle struct {
	BurstLimit int64   `json:"burst_limit,omitempty"`
	RateLimit  float64 `json:"rate_limit,omitempty"`
}

// UsagePlanKey associates an API key with a usage plan.
type UsagePlanKey struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

// DomainName represents a custom domain name for API Gateway.
type DomainName struct {
	DomainName               string                 `json:"domain_name"`
	DomainNameArn            string                 `json:"domain_name_arn,omitempty"`
	DomainNameId             string                 `json:"domain_name_id,omitempty"`
	CertificateArn           string                 `json:"certificate_arn,omitempty"`
	CertificateName          string                 `json:"certificate_name,omitempty"`
	CertificateUploadDate    time.Time              `json:"certificate_upload_date,omitempty"`
	DistributionDomainName   string                 `json:"distribution_domain_name,omitempty"`
	DistributionHostedZoneId string                 `json:"distribution_hosted_zone_id,omitempty"`
	RegionalDomainName       string                 `json:"regional_domain_name,omitempty"`
	RegionalHostedZoneId     string                 `json:"regional_hosted_zone_id,omitempty"`
	RegionalCertificateArn   string                 `json:"regional_certificate_arn,omitempty"`
	RegionalCertificateName  string                 `json:"regional_certificate_name,omitempty"`
	EndpointConfiguration    *EndpointConfiguration `json:"endpoint_configuration,omitempty"`
	DomainNameStatus         string                 `json:"domain_name_status,omitempty"`
	DomainNameStatusMessage  string                 `json:"domain_name_status_message,omitempty"`
	SecurityPolicy           string                 `json:"security_policy,omitempty"`
	Tags                     []common.Tag           `json:"tags,omitempty"`
}

// BasePathMapping maps a base path to an API Gateway API.
type BasePathMapping struct {
	DomainName string `json:"domain_name"`
	BasePath   string `json:"base_path"`
	RestApiId  string `json:"rest_api_id"`
	Stage      string `json:"stage"`
}

// RestApiListResult represents the result of listing REST APIs.
type RestApiListResult struct {
	Items       []*RestApi
	NextMarker  string
	IsTruncated bool
}

// ResourceListResult represents the result of listing API Gateway resources.
type ResourceListResult struct {
	Items       []*Resource
	NextMarker  string
	IsTruncated bool
}

// DeploymentListResult represents the result of listing API Gateway deployments.
type DeploymentListResult struct {
	Items       []*Deployment
	NextMarker  string
	IsTruncated bool
}

// StageListResult represents the result of listing API Gateway stages.
type StageListResult struct {
	Items       []*Stage
	NextMarker  string
	IsTruncated bool
}

// ApiKeyListResult represents the result of listing API keys.
type ApiKeyListResult struct {
	Items       []*ApiKey
	NextMarker  string
	IsTruncated bool
}

// UsagePlanListResult represents the result of listing usage plans.
type UsagePlanListResult struct {
	Items       []*UsagePlan
	NextMarker  string
	IsTruncated bool
}
