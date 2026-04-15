package apigateway

import (
	"net/http"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	svcapigatewayruntime "vorpalstacks/internal/services/aws/apigateway/runtime"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// apiGatewayStores holds the various API Gateway stores.
type apiGatewayStores struct {
	restApis *apigatewaystore.RestApiStore
	usage    *apigatewaystore.UsageStore
	domains  *apigatewaystore.DomainStore
}

// APIGatewayService provides AWS API Gateway operations.
type APIGatewayService struct {
	accountID     string
	region        string
	stores        sync.Map // region → *apiGatewayStores
	runtimeServer *svcapigatewayruntime.RuntimeServer
}

// NewAPIGatewayService creates a new API Gateway service instance.
func NewAPIGatewayService(accountID, region string) *APIGatewayService {
	return &APIGatewayService{
		accountID: accountID,
		region:    region,
	}
}

// InitRuntimeServer creates the runtime server using the same stores as the
// management service. The runtime server is owned by this service and must be
// initialised exactly once before use.
func (s *APIGatewayService) InitRuntimeServer(
	storageManager *storage.RegionStorageManager,
	lambdaInvoker svcapigatewayruntime.LambdaInvoker,
	sqsStore sqsstore.SQSStoreInterface,
	snsStore snsstore.SNSStoreInterface,
	bus eventbus.Bus,
) {
	storage, err := storageManager.GetStorage(s.region)
	if err != nil {
		return
	}

	stores := &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(storage, s.accountID, s.region),
		usage:    apigatewaystore.NewUsageStore(storage, s.accountID, s.region),
		domains:  apigatewaystore.NewDomainStore(storage, s.accountID, s.region),
	}
	if actual, loaded := s.stores.LoadOrStore(s.region, stores); loaded {
		stores = actual.(*apiGatewayStores)
	}

	s.runtimeServer = svcapigatewayruntime.NewRuntimeServer(stores.restApis, stores.usage, lambdaInvoker)
	s.runtimeServer.SetAccountID(s.accountID)
	if sqsStore != nil {
		s.runtimeServer.SetSQSStore(sqsStore, s.accountID, s.region)
	}
	if snsStore != nil {
		s.runtimeServer.SetSNSStore(snsStore, s.accountID, s.region)
	}
	s.runtimeServer.SetEventBus(bus)
}

// RuntimeHandler returns an http.Handler for the API Gateway runtime, or nil
// if the runtime server has not been initialised.
func (s *APIGatewayService) RuntimeHandler() http.Handler {
	if s.runtimeServer == nil {
		return nil
	}
	return http.HandlerFunc(s.runtimeServer.HandleRequest)
}

// CloseRuntimeServer stops background goroutines in the runtime server.
func (s *APIGatewayService) CloseRuntimeServer() {
	if s.runtimeServer != nil {
		s.runtimeServer.Close()
	}
}

func (s *APIGatewayService) store(reqCtx *request.RequestContext) (*apiGatewayStores, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*apiGatewayStores, error) {
		st, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return &apiGatewayStores{
			restApis: apigatewaystore.NewRestApiStore(st, s.accountID, reqCtx.GetRegion()),
			usage:    apigatewaystore.NewUsageStore(st, s.accountID, reqCtx.GetRegion()),
			domains:  apigatewaystore.NewDomainStore(st, s.accountID, reqCtx.GetRegion()),
		}, nil
	})
}

// RegisterHandlers registers the API Gateway service handlers with the dispatcher.
func (s *APIGatewayService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("apigateway", "CreateRestApi", s.CreateRestApi)
	d.RegisterHandlerForService("apigateway", "GetRestApi", s.GetRestApi)
	d.RegisterHandlerForService("apigateway", "DeleteRestApi", s.DeleteRestApi)
	d.RegisterHandlerForService("apigateway", "UpdateRestApi", s.UpdateRestApi)
	d.RegisterHandlerForService("apigateway", "GetRestApis", s.GetRestApis)

	d.RegisterHandlerForService("apigateway", "CreateResource", s.CreateResource)
	d.RegisterHandlerForService("apigateway", "GetResource", s.GetResource)
	d.RegisterHandlerForService("apigateway", "DeleteResource", s.DeleteResource)
	d.RegisterHandlerForService("apigateway", "UpdateResource", s.UpdateResource)
	d.RegisterHandlerForService("apigateway", "GetResources", s.GetResources)

	d.RegisterHandlerForService("apigateway", "PutMethod", s.PutMethod)
	d.RegisterHandlerForService("apigateway", "GetMethod", s.GetMethod)
	d.RegisterHandlerForService("apigateway", "DeleteMethod", s.DeleteMethod)
	d.RegisterHandlerForService("apigateway", "UpdateMethod", s.UpdateMethod)

	d.RegisterHandlerForService("apigateway", "PutMethodResponse", s.PutMethodResponse)
	d.RegisterHandlerForService("apigateway", "GetMethodResponse", s.GetMethodResponse)
	d.RegisterHandlerForService("apigateway", "DeleteMethodResponse", s.DeleteMethodResponse)

	d.RegisterHandlerForService("apigateway", "PutIntegration", s.PutIntegration)
	d.RegisterHandlerForService("apigateway", "GetIntegration", s.GetIntegration)
	d.RegisterHandlerForService("apigateway", "DeleteIntegration", s.DeleteIntegration)
	d.RegisterHandlerForService("apigateway", "UpdateIntegration", s.UpdateIntegration)
	d.RegisterHandlerForService("apigateway", "PutIntegrationResponse", s.PutIntegrationResponse)
	d.RegisterHandlerForService("apigateway", "GetIntegrationResponse", s.GetIntegrationResponse)
	d.RegisterHandlerForService("apigateway", "DeleteIntegrationResponse", s.DeleteIntegrationResponse)
	d.RegisterHandlerForService("apigateway", "UpdateIntegrationResponse", s.UpdateIntegrationResponse)

	d.RegisterHandlerForService("apigateway", "CreateDeployment", s.CreateDeployment)
	d.RegisterHandlerForService("apigateway", "GetDeployment", s.GetDeployment)
	d.RegisterHandlerForService("apigateway", "DeleteDeployment", s.DeleteDeployment)
	d.RegisterHandlerForService("apigateway", "UpdateDeployment", s.UpdateDeployment)
	d.RegisterHandlerForService("apigateway", "GetDeployments", s.GetDeployments)

	d.RegisterHandlerForService("apigateway", "CreateStage", s.CreateStage)
	d.RegisterHandlerForService("apigateway", "GetStage", s.GetStage)
	d.RegisterHandlerForService("apigateway", "DeleteStage", s.DeleteStage)
	d.RegisterHandlerForService("apigateway", "UpdateStage", s.UpdateStage)
	d.RegisterHandlerForService("apigateway", "GetStages", s.GetStages)

	d.RegisterHandlerForService("apigateway", "CreateRequestValidator", s.CreateRequestValidator)
	d.RegisterHandlerForService("apigateway", "GetRequestValidator", s.GetRequestValidator)
	d.RegisterHandlerForService("apigateway", "DeleteRequestValidator", s.DeleteRequestValidator)
	d.RegisterHandlerForService("apigateway", "UpdateRequestValidator", s.UpdateRequestValidator)
	d.RegisterHandlerForService("apigateway", "GetRequestValidators", s.GetRequestValidators)

	d.RegisterHandlerForService("apigateway", "CreateModel", s.CreateModel)
	d.RegisterHandlerForService("apigateway", "GetModel", s.GetModel)
	d.RegisterHandlerForService("apigateway", "DeleteModel", s.DeleteModel)
	d.RegisterHandlerForService("apigateway", "UpdateModel", s.UpdateModel)
	d.RegisterHandlerForService("apigateway", "GetModels", s.GetModels)

	d.RegisterHandlerForService("apigateway", "CreateAuthorizer", s.CreateAuthorizer)
	d.RegisterHandlerForService("apigateway", "GetAuthorizer", s.GetAuthorizer)
	d.RegisterHandlerForService("apigateway", "DeleteAuthorizer", s.DeleteAuthorizer)
	d.RegisterHandlerForService("apigateway", "UpdateAuthorizer", s.UpdateAuthorizer)
	d.RegisterHandlerForService("apigateway", "GetAuthorizers", s.GetAuthorizers)
	d.RegisterHandlerForService("apigateway", "TestInvokeAuthorizer", s.TestInvokeAuthorizer)
	d.RegisterHandlerForService("apigateway", "TestInvokeMethod", s.TestInvokeMethod)

	d.RegisterHandlerForService("apigateway", "CreateApiKey", s.CreateApiKey)
	d.RegisterHandlerForService("apigateway", "GetApiKey", s.GetApiKey)
	d.RegisterHandlerForService("apigateway", "DeleteApiKey", s.DeleteApiKey)
	d.RegisterHandlerForService("apigateway", "UpdateApiKey", s.UpdateApiKey)
	d.RegisterHandlerForService("apigateway", "GetApiKeys", s.GetApiKeys)

	d.RegisterHandlerForService("apigateway", "CreateUsagePlan", s.CreateUsagePlan)
	d.RegisterHandlerForService("apigateway", "GetUsagePlan", s.GetUsagePlan)
	d.RegisterHandlerForService("apigateway", "DeleteUsagePlan", s.DeleteUsagePlan)
	d.RegisterHandlerForService("apigateway", "UpdateUsagePlan", s.UpdateUsagePlan)
	d.RegisterHandlerForService("apigateway", "GetUsagePlans", s.GetUsagePlans)

	d.RegisterHandlerForService("apigateway", "CreateUsagePlanKey", s.CreateUsagePlanKey)
	d.RegisterHandlerForService("apigateway", "GetUsagePlanKey", s.GetUsagePlanKey)
	d.RegisterHandlerForService("apigateway", "DeleteUsagePlanKey", s.DeleteUsagePlanKey)
	d.RegisterHandlerForService("apigateway", "GetUsagePlanKeys", s.GetUsagePlanKeys)
	d.RegisterHandlerForService("apigateway", "GetUsage", s.GetUsage)

	d.RegisterHandlerForService("apigateway", "CreateDomainName", s.CreateDomainName)
	d.RegisterHandlerForService("apigateway", "GetDomainName", s.GetDomainName)
	d.RegisterHandlerForService("apigateway", "DeleteDomainName", s.DeleteDomainName)
	d.RegisterHandlerForService("apigateway", "UpdateDomainName", s.UpdateDomainName)
	d.RegisterHandlerForService("apigateway", "GetDomainNames", s.GetDomainNames)

	d.RegisterHandlerForService("apigateway", "CreateBasePathMapping", s.CreateBasePathMapping)
	d.RegisterHandlerForService("apigateway", "GetBasePathMapping", s.GetBasePathMapping)
	d.RegisterHandlerForService("apigateway", "DeleteBasePathMapping", s.DeleteBasePathMapping)
	d.RegisterHandlerForService("apigateway", "UpdateBasePathMapping", s.UpdateBasePathMapping)
	d.RegisterHandlerForService("apigateway", "GetBasePathMappings", s.GetBasePathMappings)

	d.RegisterHandlerForService("apigateway", "TagResource", s.TagResource)
	d.RegisterHandlerForService("apigateway", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("apigateway", "ListTagsForResource", s.ListTagsForResource)
}
