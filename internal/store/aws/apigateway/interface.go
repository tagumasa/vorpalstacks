package apigateway

import (
	"vorpalstacks/internal/store/aws/common"
)

// RestApiStoreInterface defines operations for managing REST APIs.
type RestApiStoreInterface interface {
	Create(api *RestApi) (*RestApi, error)
	Get(apiId string) (*RestApi, error)
	Update(api *RestApi) error
	Delete(apiId string) error
	List(opts common.ListOptions) (*common.ListResult[RestApi], error)
	GetTags(apiId string) ([]common.Tag, error)
	TagResource(apiId string, tags map[string]string) error
	UntagResource(apiId string, tagKeys []string) error
	GenerateRevisionId() string
	CreateRequestValidator(apiId string, validator *RequestValidator) (*RequestValidator, error)
	GetRequestValidator(apiId, validatorId string) (*RequestValidator, error)
	UpdateRequestValidator(apiId string, validator *RequestValidator) error
	DeleteRequestValidator(apiId, validatorId string) error
	ListRequestValidators(apiId string) ([]*RequestValidator, error)
	CreateModel(apiId string, model *Model) (*Model, error)
	GetModel(apiId, modelName string) (*Model, error)
	DeleteModel(apiId, modelName string) error
	ListModels(apiId string) ([]*Model, error)
	UpdateModel(apiId string, model *Model) error
	CreateAuthorizer(apiId string, authorizer *Authorizer) (*Authorizer, error)
	GetAuthorizer(apiId, authorizerId string) (*Authorizer, error)
	UpdateAuthorizer(apiId string, authorizer *Authorizer) error
	DeleteAuthorizer(apiId, authorizerId string) error
	ListAuthorizers(apiId string) ([]*Authorizer, error)
}

// UsageStoreInterface defines operations for managing API keys and usage plans.
type UsageStoreInterface interface {
	CreateApiKey(apiKey *ApiKey) (*ApiKey, error)
	GetApiKey(apiKeyId string) (*ApiKey, error)
	UpdateApiKey(apiKey *ApiKey) error
	DeleteApiKey(apiKeyId string) error
	ListApiKeys(opts common.ListOptions) (*common.ListResult[ApiKey], error)
	GetApiKeyByValue(value string) (*ApiKey, error)
	CreateUsagePlan(usagePlan *UsagePlan) (*UsagePlan, error)
	GetUsagePlan(usagePlanId string) (*UsagePlan, error)
	UpdateUsagePlan(usagePlan *UsagePlan) error
	DeleteUsagePlan(usagePlanId string) error
	ListUsagePlans(opts common.ListOptions) (*common.ListResult[UsagePlan], error)
	CreateUsagePlanKey(usagePlanId string, key *UsagePlanKey) (*UsagePlanKey, error)
	GetUsagePlanKey(usagePlanId, keyId string) (*UsagePlanKey, error)
	DeleteUsagePlanKey(usagePlanId, keyId string) error
	ListUsagePlanKeys(usagePlanId string, opts common.ListOptions) (*common.ListResult[UsagePlanKey], error)
	ListUsagePlansForAPIKey(apiKeyId string) ([]*UsagePlan, error)
	ListUsageRecordsForAPIKey(usagePlanId, apiKeyId, startDate, endDate string) ([]*UsageRecord, error)
	GetUsage(usagePlanId, apiKeyId, date string) (*UsageRecord, error)
	RecordUsage(record *UsageRecord) error
}

// DomainStoreInterface defines operations for managing custom domain names.
type DomainStoreInterface interface {
	CreateDomainName(domain *DomainName) (*DomainName, error)
	GetDomainName(domainName string) (*DomainName, error)
	UpdateDomainName(domain *DomainName) error
	DeleteDomainName(domainName string) error
	ListDomainNames(opts common.ListOptions) (*common.ListResult[DomainName], error)
	CreateBasePathMapping(domainName string, mapping *BasePathMapping) (*BasePathMapping, error)
	GetBasePathMapping(domainName, basePath string) (*BasePathMapping, error)
	UpdateBasePathMapping(domainName, basePath string, mapping *BasePathMapping) error
	DeleteBasePathMapping(domainName, basePath string) error
	ListBasePathMappings(domainName string, opts common.ListOptions) (*common.ListResult[BasePathMapping], error)
}

// APIGatewayStoresInterface defines access to all API Gateway stores.
type APIGatewayStoresInterface interface {
	RestApis() RestApiStoreInterface
	Usage() UsageStoreInterface
	Domains() DomainStoreInterface
	RestApisRaw() *RestApiStore
	UsageRaw() *UsageStore
	DomainsRaw() *DomainStore
}

// APIGatewayStore provides access to all API Gateway stores.
type APIGatewayStore struct {
	restApis *RestApiStore
	usage    *UsageStore
	domains  *DomainStore
}

// NewAPIGatewayStore creates a new APIGatewayStore with the given stores.
func NewAPIGatewayStore(restApis *RestApiStore, usage *UsageStore, domains *DomainStore) *APIGatewayStore {
	return &APIGatewayStore{
		restApis: restApis,
		usage:    usage,
		domains:  domains,
	}
}

// RestApis returns the REST API store.
func (s *APIGatewayStore) RestApis() RestApiStoreInterface {
	return s.restApis
}

// Usage returns the Usage store.
func (s *APIGatewayStore) Usage() UsageStoreInterface {
	return s.usage
}

// Domains returns the Domain store.
func (s *APIGatewayStore) Domains() DomainStoreInterface {
	return s.domains
}

// RestApisRaw returns the raw REST API store.
func (s *APIGatewayStore) RestApisRaw() *RestApiStore {
	return s.restApis
}

// UsageRaw returns the raw Usage store.
func (s *APIGatewayStore) UsageRaw() *UsageStore {
	return s.usage
}

// DomainsRaw returns the raw Domain store.
func (s *APIGatewayStore) DomainsRaw() *DomainStore {
	return s.domains
}

var (
	_ RestApiStoreInterface     = (*RestApiStore)(nil)
	_ UsageStoreInterface       = (*UsageStore)(nil)
	_ DomainStoreInterface      = (*DomainStore)(nil)
	_ APIGatewayStoresInterface = (*APIGatewayStore)(nil)
)

// WrapRestApiStore wraps a RestApiStore in its interface.
func WrapRestApiStore(s *RestApiStore) RestApiStoreInterface {
	return s
}

// WrapUsageStore wraps a UsageStore in its interface.
func WrapUsageStore(s *UsageStore) UsageStoreInterface {
	return s
}

// WrapDomainStore wraps a DomainStore in its interface.
func WrapDomainStore(s *DomainStore) DomainStoreInterface {
	return s
}
