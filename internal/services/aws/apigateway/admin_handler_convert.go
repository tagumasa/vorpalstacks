package apigateway

import (
	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	aws_types "vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

func toPbRestApi(api *apigatewaystore.RestApi) *pb.RestApi {
	pbApi := &pb.RestApi{
		Id:                     api.Id,
		Name:                   api.Name,
		Description:            api.Description,
		Version:                api.Version,
		Warnings:               api.Warnings,
		Createddate:            api.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		Binarymediatypes:       api.BinaryMediaTypes,
		Minimumcompressionsize: derefInt32(api.MinimumCompressionSize),
		Apikeysource:           toPbApiKeySourceType(api.ApiKeySource),
		Policy:                 api.Policy,
		Tags:                   tagsToPbMap(api.Tags),
		Securitypolicy:         toPbSecurityPolicy(api.SecurityPolicy),
		Endpointaccessmode:     toPbEndpointAccessMode(api.EndpointAccessMode),
	}
	if api.DisableExecuteApiEndpoint {
		pbApi.Disableexecuteapiendpoint = proto.Bool(api.DisableExecuteApiEndpoint)
	}
	if api.EndpointConfiguration != nil {
		types := make([]pb.EndpointType, len(api.EndpointConfiguration.Types))
		for i, t := range api.EndpointConfiguration.Types {
			switch t {
			case "PRIVATE":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_PRIVATE
			case "REGIONAL":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_REGIONAL
			case "EDGE":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_EDGE
			}
		}
		pbApi.Endpointconfiguration = &pb.EndpointConfiguration{Types: types}
	}
	return pbApi
}

func toPbResource(r *apigatewaystore.Resource) *pb.Resource {
	pbR := &pb.Resource{
		Id:       r.Id,
		Parentid: r.ParentId,
		Path:     r.Path,
		Pathpart: r.PathPart,
	}
	if len(r.ResourceMethods) > 0 {
		pbR.Resourcemethods = make(map[string]*pb.Method)
		for k, m := range r.ResourceMethods {
			pbR.Resourcemethods[k] = toPbMethod(m)
		}
	}
	return pbR
}

func toPbMethod(m *apigatewaystore.Method) *pb.Method {
	pbM := &pb.Method{
		Httpmethod:         m.HttpMethod,
		Authorizationtype:  m.AuthorizationType,
		Apikeyrequired:     proto.Bool(m.ApiKeyRequired),
		Authorizerid:       m.AuthorizerId,
		Requestvalidatorid: m.RequestValidatorId,
		Operationname:      m.OperationName,
		Requestparameters:  m.RequestParameters,
		Requestmodels:      m.RequestModels,
	}
	if m.MethodIntegration != nil {
		pbM.Methodintegration = toPbIntegration(m.MethodIntegration)
	}
	if len(m.MethodResponses) > 0 {
		pbM.Methodresponses = make(map[string]*pb.MethodResponse)
		for k, r := range m.MethodResponses {
			pbM.Methodresponses[k] = toPbMethodResponse(r)
		}
	}
	return pbM
}

func toPbIntegration(i *apigatewaystore.Integration) *pb.Integration {
	pbI := &pb.Integration{
		Type:                toPbIntegrationType(i.Type),
		Httpmethod:          i.IntegrationHttpMethod,
		Uri:                 i.Uri,
		Credentials:         i.Credentials,
		Passthroughbehavior: i.PassthroughBehavior,
		Cachenamespace:      i.CacheNamespace,
		Connectiontype:      toPbConnectionType(i.ConnectionType),
		Connectionid:        i.ConnectionId,
		Requestparameters:   i.RequestParameters,
		Requesttemplates:    i.RequestTemplates,
		Cachekeyparameters:  i.CacheKeyParameters,
	}
	if i.ContentHandling != "" {
		pbI.Contenthandling = toPbContentHandling(i.ContentHandling)
	}
	if i.TimeoutInMillis > 0 {
		pbI.Timeoutinmillis = i.TimeoutInMillis
	}
	if len(i.IntegrationResponses) > 0 {
		pbI.Integrationresponses = make(map[string]*pb.IntegrationResponse)
		for k, r := range i.IntegrationResponses {
			pbI.Integrationresponses[k] = toPbIntegrationResponse(r)
		}
	}
	return pbI
}

func toPbIntegrationResponse(r *apigatewaystore.IntegrationResponse) *pb.IntegrationResponse {
	pbR := &pb.IntegrationResponse{
		Statuscode:       r.StatusCode,
		Selectionpattern: r.SelectionPattern,
	}
	if r.ContentHandling != "" {
		pbR.Contenthandling = toPbContentHandling(r.ContentHandling)
	}
	if len(r.ResponseParameters) > 0 {
		pbR.Responseparameters = r.ResponseParameters
	}
	if len(r.ResponseTemplates) > 0 {
		pbR.Responsetemplates = r.ResponseTemplates
	}
	return pbR
}

func toPbMethodResponse(r *apigatewaystore.MethodResponse) *pb.MethodResponse {
	return &pb.MethodResponse{
		Statuscode:         r.StatusCode,
		Responseparameters: r.ResponseParameters,
		Responsemodels:     r.ResponseModels,
	}
}

func fromPbAuthorizerType(t pb.AuthorizerType) string {
	switch t {
	case pb.AuthorizerType_AUTHORIZER_TYPE_TOKEN:
		return "TOKEN"
	case pb.AuthorizerType_AUTHORIZER_TYPE_REQUEST:
		return "REQUEST"
	case pb.AuthorizerType_AUTHORIZER_TYPE_COGNITO_USER_POOLS:
		return "COGNITO_USER_POOLS"
	default:
		return ""
	}
}

func fromPbIntegrationType(t pb.IntegrationType) string {
	switch t {
	case pb.IntegrationType_INTEGRATION_TYPE_HTTP:
		return "HTTP"
	case pb.IntegrationType_INTEGRATION_TYPE_HTTP_PROXY:
		return "HTTP_PROXY"
	case pb.IntegrationType_INTEGRATION_TYPE_AWS:
		return "AWS"
	case pb.IntegrationType_INTEGRATION_TYPE_AWS_PROXY:
		return "AWS_PROXY"
	case pb.IntegrationType_INTEGRATION_TYPE_MOCK:
		return "MOCK"
	default:
		return ""
	}
}

func fromPbConnectionType(t pb.ConnectionType) string {
	switch t {
	case pb.ConnectionType_CONNECTION_TYPE_INTERNET:
		return "INTERNET"
	case pb.ConnectionType_CONNECTION_TYPE_VPC_LINK:
		return "VPC_LINK"
	default:
		return ""
	}
}

func fromPbContentHandling(ch pb.ContentHandlingStrategy) string {
	switch ch {
	case pb.ContentHandlingStrategy_CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY:
		return "CONVERT_TO_BINARY"
	case pb.ContentHandlingStrategy_CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT:
		return "CONVERT_TO_TEXT"
	default:
		return ""
	}
}

func toPbConnectionType(t string) pb.ConnectionType {
	switch t {
	case "VPC_LINK":
		return pb.ConnectionType_CONNECTION_TYPE_VPC_LINK
	default:
		return pb.ConnectionType_CONNECTION_TYPE_INTERNET
	}
}

func toPbIntegrationType(t string) pb.IntegrationType {
	switch t {
	case "HTTP":
		return pb.IntegrationType_INTEGRATION_TYPE_HTTP
	case "HTTP_PROXY":
		return pb.IntegrationType_INTEGRATION_TYPE_HTTP_PROXY
	case "AWS":
		return pb.IntegrationType_INTEGRATION_TYPE_AWS
	case "AWS_PROXY":
		return pb.IntegrationType_INTEGRATION_TYPE_AWS_PROXY
	case "MOCK":
		return pb.IntegrationType_INTEGRATION_TYPE_MOCK
	default:
		return pb.IntegrationType(0)
	}
}

func toPbContentHandling(ch string) pb.ContentHandlingStrategy {
	switch ch {
	case "CONVERT_TO_BINARY":
		return pb.ContentHandlingStrategy_CONTENT_HANDLING_STRATEGY_CONVERT_TO_BINARY
	case "CONVERT_TO_TEXT":
		return pb.ContentHandlingStrategy_CONTENT_HANDLING_STRATEGY_CONVERT_TO_TEXT
	default:
		return pb.ContentHandlingStrategy(0)
	}
}

func toPbDeployment(d *apigatewaystore.Deployment) *pb.Deployment {
	return &pb.Deployment{
		Id:          d.Id,
		Description: d.Description,
		Createddate: d.CreatedDate.Format(timeutils.ISO8601UTCFormat),
	}
}

func toPbStage(s *apigatewaystore.Stage) *pb.Stage {
	pbS := &pb.Stage{
		Stagename:            s.StageName,
		Deploymentid:         s.DeploymentId,
		Description:          s.Description,
		Cacheclusterenabled:  proto.Bool(s.CacheClusterEnabled),
		Tracingenabled:       proto.Bool(s.TracingEnabled),
		Createddate:          s.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		Lastupdateddate:      s.LastUpdatedDate.Format(timeutils.ISO8601UTCFormat),
		Documentationversion: s.DocumentationVersion,
		Webaclarn:            s.WebAclArn,
		Tags:                 tagsToPbMap(s.Tags),
	}
	if s.CacheClusterSize != "" {
		pbS.Cacheclustersize = toPbCacheClusterSize(s.CacheClusterSize)
	}
	if s.CacheClusterStatus != "" {
		pbS.Cacheclusterstatus = toPbCacheClusterStatus(s.CacheClusterStatus)
	}
	if len(s.Variables) > 0 {
		pbS.Variables = s.Variables
	}
	if s.AccessLogSettings != nil {
		pbS.Accesslogsettings = &pb.AccessLogSettings{
			Destinationarn: s.AccessLogSettings.DestinationArn,
			Format:         s.AccessLogSettings.Format,
		}
	}
	if s.CanarySettings != nil {
		pbS.Canarysettings = &pb.CanarySettings{
			Percenttraffic:         s.CanarySettings.PercentTraffic,
			Deploymentid:           s.CanarySettings.DeploymentId,
			Stagevariableoverrides: s.CanarySettings.StageVariableOverrides,
			Usestagecache:          proto.Bool(s.CanarySettings.UseStageCache),
		}
	}
	if len(s.MethodSettings) > 0 {
		pbS.Methodsettings = make(map[string]*pb.MethodSetting)
		for k, v := range s.MethodSettings {
			pbS.Methodsettings[k] = toPbMethodSetting(v)
		}
	}
	return pbS
}

func toPbApiKey(k *apigatewaystore.ApiKey, includeValue bool) *pb.ApiKey {
	pbK := &pb.ApiKey{
		Id:              k.Id,
		Name:            k.Name,
		Enabled:         proto.Bool(k.Enabled),
		Createddate:     k.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		Lastupdateddate: k.LastUpdatedDate.Format(timeutils.ISO8601UTCFormat),
		Description:     k.Description,
		Customerid:      k.CustomerId,
		Stagekeys:       k.StageKeys,
		Tags:            tagsToPbMap(k.Tags),
	}
	if includeValue && k.Value != "" {
		pbK.Value = k.Value
	}
	return pbK
}

func toPbUsagePlan(p *apigatewaystore.UsagePlan) *pb.UsagePlan {
	pbP := &pb.UsagePlan{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Productcode: p.ProductCode,
		Tags:        tagsToPbMap(p.Tags),
	}
	for _, as := range p.ApiStages {
		pbStage := &pb.ApiStage{
			Apiid: as.ApiId,
			Stage: as.Stage,
		}
		if len(as.Throttle) > 0 {
			pbStage.Throttle = make(map[string]*pb.ThrottleSettings)
			for k, v := range as.Throttle {
				pbStage.Throttle[k] = &pb.ThrottleSettings{
					Burstlimit: int32(v.BurstLimit),
					Ratelimit:  v.RateLimit,
				}
			}
		}
		pbP.Apistages = append(pbP.Apistages, pbStage)
	}
	if p.Quota != nil {
		pbP.Quota = &pb.QuotaSettings{
			Limit:  int32(p.Quota.Limit),
			Offset: int32(p.Quota.Offset),
			Period: toPbQuotaPeriodType(p.Quota.Period),
		}
	}
	if p.Throttle != nil {
		pbP.Throttle = &pb.ThrottleSettings{
			Burstlimit: int32(p.Throttle.BurstLimit),
			Ratelimit:  p.Throttle.RateLimit,
		}
	}
	return pbP
}

func toPbUsagePlanKey(k *apigatewaystore.UsagePlanKey) *pb.UsagePlanKey {
	return &pb.UsagePlanKey{
		Id:    k.Id,
		Type:  k.Type,
		Value: k.Value,
		Name:  k.Name,
	}
}

func toPbAuthorizer(a *apigatewaystore.Authorizer) *pb.Authorizer {
	return &pb.Authorizer{
		Id:                           a.Id,
		Name:                         a.Name,
		Type:                         toPbAuthorizerType(a.Type),
		Authtype:                     a.AuthType,
		Authorizeruri:                a.AuthorizerUri,
		Authorizercredentials:        a.AuthorizerCredentials,
		Identitysource:               a.IdentitySource,
		Identityvalidationexpression: a.IdentityValidationExpression,
		Authorizerresultttlinseconds: a.AuthorizerResultTtlInSeconds,
		Providerarns:                 a.ProviderArns,
	}
}

func toPbAuthorizerType(t string) pb.AuthorizerType {
	switch t {
	case "TOKEN":
		return pb.AuthorizerType_AUTHORIZER_TYPE_TOKEN
	case "REQUEST":
		return pb.AuthorizerType_AUTHORIZER_TYPE_REQUEST
	case "COGNITO_USER_POOLS":
		return pb.AuthorizerType_AUTHORIZER_TYPE_COGNITO_USER_POOLS
	default:
		return pb.AuthorizerType(0)
	}
}

func tagsToPbMap(tags []aws_types.Tag) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	m := make(map[string]string, len(tags))
	for _, t := range tags {
		m[t.Key] = t.Value
	}
	return m
}

func toPbApiKeySourceType(s string) pb.ApiKeySourceType {
	switch s {
	case "HEADER":
		return pb.ApiKeySourceType_API_KEY_SOURCE_TYPE_HEADER
	case "AUTHORIZER":
		return pb.ApiKeySourceType_API_KEY_SOURCE_TYPE_AUTHORIZER
	default:
		return pb.ApiKeySourceType(0)
	}
}

func toPbQuotaPeriodType(s string) pb.QuotaPeriodType {
	switch s {
	case "DAY":
		return pb.QuotaPeriodType_QUOTA_PERIOD_TYPE_DAY
	case "WEEK":
		return pb.QuotaPeriodType_QUOTA_PERIOD_TYPE_WEEK
	case "MONTH":
		return pb.QuotaPeriodType_QUOTA_PERIOD_TYPE_MONTH
	default:
		return pb.QuotaPeriodType(0)
	}
}

func toPbCacheClusterSize(s string) pb.CacheClusterSize {
	switch s {
	case "0.5":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB
	case "1.6":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB
	case "6.1":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB
	case "13.5":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB
	case "28.4":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB
	case "58.2":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB
	case "118":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_118_GB
	case "237":
		return pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_237_GB
	default:
		return pb.CacheClusterSize(0)
	}
}

func toPbCacheClusterStatus(s string) pb.CacheClusterStatus {
	switch s {
	case "CREATE_IN_PROGRESS":
		return pb.CacheClusterStatus_CACHE_CLUSTER_STATUS_CREATE_IN_PROGRESS
	case "AVAILABLE":
		return pb.CacheClusterStatus_CACHE_CLUSTER_STATUS_AVAILABLE
	case "DELETE_IN_PROGRESS":
		return pb.CacheClusterStatus_CACHE_CLUSTER_STATUS_DELETE_IN_PROGRESS
	case "FLUSH_IN_PROGRESS":
		return pb.CacheClusterStatus_CACHE_CLUSTER_STATUS_FLUSH_IN_PROGRESS
	case "NOT_AVAILABLE":
		return pb.CacheClusterStatus_CACHE_CLUSTER_STATUS_NOT_AVAILABLE
	default:
		return pb.CacheClusterStatus(0)
	}
}

func toPbMethodSetting(ms *apigatewaystore.MethodSetting) *pb.MethodSetting {
	if ms == nil {
		return nil
	}
	return &pb.MethodSetting{
		Metricsenabled:                      proto.Bool(ms.MetricsEnabled),
		Logginglevel:                        ms.LoggingLevel,
		Datatraceenabled:                    proto.Bool(ms.DataTraceEnabled),
		Throttlingburstlimit:                ms.ThrottlingBurstLimit,
		Throttlingratelimit:                 ms.ThrottlingRateLimit,
		Cachingenabled:                      proto.Bool(ms.CachingEnabled),
		Cachettlinseconds:                   ms.CacheTtlInSeconds,
		Cachedataencrypted:                  proto.Bool(ms.CacheDataEncrypted),
		Requireauthorizationforcachecontrol: proto.Bool(ms.RequireAuthorizationForCacheControl),
	}
}

// Reverse enum conversion functions for admin handler request parsing.

func apiKeySourceFromPb(v pb.ApiKeySourceType) string {
	switch v {
	case pb.ApiKeySourceType_API_KEY_SOURCE_TYPE_HEADER:
		return "HEADER"
	case pb.ApiKeySourceType_API_KEY_SOURCE_TYPE_AUTHORIZER:
		return "AUTHORIZER"
	default:
		return ""
	}
}

func cacheClusterSizeFromPb(v pb.CacheClusterSize) string {
	switch v {
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_0_POINT_5_GB:
		return "0.5"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_1_POINT_6_GB:
		return "1.6"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_6_POINT_1_GB:
		return "6.1"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_13_POINT_5_GB:
		return "13.5"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_28_POINT_4_GB:
		return "28.4"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_58_POINT_2_GB:
		return "58.2"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_118_GB:
		return "118"
	case pb.CacheClusterSize_CACHE_CLUSTER_SIZE_SIZE_237_GB:
		return "237"
	default:
		return ""
	}
}

func securityPolicyFromPb(v pb.SecurityPolicy) string {
	switch v {
	case pb.SecurityPolicy_SECURITY_POLICY_TLS_1_0:
		return "TLS_1_0"
	case pb.SecurityPolicy_SECURITY_POLICY_TLS_1_2:
		return "TLS_1_2"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09:
		return "SecurityPolicy_TLS13_1_2_PQ_2025_09"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE:
		return "SecurityPolicy_TLS13_2025_EDGE"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09:
		return "SecurityPolicy_TLS13_1_2_FIPS_PQ_2025_09"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE:
		return "SecurityPolicy_TLS12_2018_EDGE"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE:
		return "SecurityPolicy_TLS12_PFS_2025_EDGE"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06:
		return "SecurityPolicy_TLS13_1_2_2021_06"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09:
		return "SecurityPolicy_TLS13_1_3_2025_09"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09:
		return "SecurityPolicy_TLS13_1_2_PFS_PQ_2025_09"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09:
		return "SecurityPolicy_TLS13_1_3_FIPS_2025_09"
	case pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PFS_PQ_2025_09:
		return "SecurityPolicy_TLS13_1_2_FIPS_PFS_PQ_2025_09"
	default:
		return ""
	}
}

func endpointAccessModeFromPb(v pb.EndpointAccessMode) string {
	switch v {
	case pb.EndpointAccessMode_ENDPOINT_ACCESS_MODE_BASIC:
		return "BASIC"
	case pb.EndpointAccessMode_ENDPOINT_ACCESS_MODE_STRICT:
		return "STRICT"
	default:
		return ""
	}
}

func toPbSecurityPolicy(s string) pb.SecurityPolicy {
	switch s {
	case "TLS_1_0":
		return pb.SecurityPolicy_SECURITY_POLICY_TLS_1_0
	case "TLS_1_2":
		return pb.SecurityPolicy_SECURITY_POLICY_TLS_1_2
	case "SecurityPolicy_TLS13_1_2_PQ_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PQ_2025_09
	case "SecurityPolicy_TLS13_2025_EDGE":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_2025_EDGE
	case "SecurityPolicy_TLS13_1_2_FIPS_PQ_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PQ_2025_09
	case "SecurityPolicy_TLS12_2018_EDGE":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS12_2018_EDGE
	case "SecurityPolicy_TLS12_PFS_2025_EDGE":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS12_PFS_2025_EDGE
	case "SecurityPolicy_TLS13_1_2_2021_06":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_2021_06
	case "SecurityPolicy_TLS13_1_3_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_2025_09
	case "SecurityPolicy_TLS13_1_2_PFS_PQ_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_PFS_PQ_2025_09
	case "SecurityPolicy_TLS13_1_3_FIPS_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_3_FIPS_2025_09
	case "SecurityPolicy_TLS13_1_2_FIPS_PFS_PQ_2025_09":
		return pb.SecurityPolicy_SECURITY_POLICY_SECURITYPOLICY_TLS13_1_2_FIPS_PFS_PQ_2025_09
	default:
		return pb.SecurityPolicy(0)
	}
}

func toPbEndpointAccessMode(s string) pb.EndpointAccessMode {
	switch s {
	case "BASIC":
		return pb.EndpointAccessMode_ENDPOINT_ACCESS_MODE_BASIC
	case "STRICT":
		return pb.EndpointAccessMode_ENDPOINT_ACCESS_MODE_STRICT
	default:
		return pb.EndpointAccessMode(0)
	}
}

// fromPbEndpointType converts a proto EndpointType enum to the short
// string form used in storage ("PRIVATE", "REGIONAL", "EDGE").
// Returns "" for unspecified or unknown values.
func fromPbEndpointType(t pb.EndpointType) string {
	switch t {
	case pb.EndpointType_ENDPOINT_TYPE_PRIVATE:
		return "PRIVATE"
	case pb.EndpointType_ENDPOINT_TYPE_REGIONAL:
		return "REGIONAL"
	case pb.EndpointType_ENDPOINT_TYPE_EDGE:
		return "EDGE"
	default:
		return ""
	}
}
