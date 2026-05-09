package apigateway

import (
	"errors"
	"net/http"
	"sync"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	stores         sync.Map
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

type adminStores struct {
	restApis *apigatewaystore.RestApiStore
	usage    *apigatewaystore.UsageStore
	domains  *apigatewaystore.DomainStore
}

func (h *AdminHandler) getStores(headers http.Header) (*adminStores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return storecommon.GetOrCreateStoreE(&h.stores, region, func() (*adminStores, error) {
		regionStorage, err := h.storageManager.GetStorage(region)
		if err != nil {
			return nil, err
		}
		return &adminStores{
			restApis: apigatewaystore.NewRestApiStore(regionStorage, h.accountId, region),
			usage:    apigatewaystore.NewUsageStore(regionStorage, h.accountId, region),
			domains:  apigatewaystore.NewDomainStore(regionStorage, h.accountId, region),
		}, nil
	})
}

func storeErr(err error) error {
	notFoundErrors := []error{
		apigatewaystore.ErrRestApiNotFound,
		apigatewaystore.ErrResourceNotFound,
		apigatewaystore.ErrMethodNotFound,
		apigatewaystore.ErrIntegrationNotFound,
		apigatewaystore.ErrDeploymentNotFound,
		apigatewaystore.ErrStageNotFound,
		apigatewaystore.ErrRequestValidatorNotFound,
		apigatewaystore.ErrModelNotFound,
		apigatewaystore.ErrApiKeyNotFound,
		apigatewaystore.ErrUsagePlanNotFound,
		apigatewaystore.ErrUsagePlanKeyNotFound,
		apigatewaystore.ErrDomainNameNotFound,
		apigatewaystore.ErrBasePathMappingNotFound,
		apigatewaystore.ErrAuthorizerNotFound,
		apigatewaystore.ErrMethodResponseNotFound,
		apigatewaystore.ErrIntegrationResponseNotFound,
	}
	for _, nf := range notFoundErrors {
		if errors.Is(err, nf) {
			return connect.NewError(connect.CodeNotFound, err)
		}
	}

	alreadyExistsErrors := []error{
		apigatewaystore.ErrRestApiAlreadyExists,
		apigatewaystore.ErrResourceAlreadyExists,
		apigatewaystore.ErrDeploymentAlreadyExists,
		apigatewaystore.ErrStageAlreadyExists,
		apigatewaystore.ErrRequestValidatorAlreadyExists,
		apigatewaystore.ErrModelAlreadyExists,
		apigatewaystore.ErrApiKeyAlreadyExists,
		apigatewaystore.ErrUsagePlanAlreadyExists,
		apigatewaystore.ErrUsagePlanKeyAlreadyExists,
		apigatewaystore.ErrDomainNameAlreadyExists,
		apigatewaystore.ErrBasePathMappingAlreadyExists,
		apigatewaystore.ErrAuthorizerAlreadyExists,
	}
	for _, ae := range alreadyExistsErrors {
		if errors.Is(err, ae) {
			return connect.NewError(connect.CodeAlreadyExists, err)
		}
	}

	if errors.Is(err, apigatewaystore.ErrDeploymentInUse) {
		return connect.NewError(connect.CodeFailedPrecondition, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}

func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler(sm, accountID))
}
