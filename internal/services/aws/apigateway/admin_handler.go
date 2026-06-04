package apigateway

import (
	"errors"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// AdminHandler implements the gRPC admin console handlers for API Gateway.
type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
	service *APIGatewayService
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new API Gateway admin handler.
func NewAdminHandler(svc *APIGatewayService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStores(headers http.Header) (*apiGatewayStores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
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

// NewConnectHandler returns the connect RPC path and handler for API Gateway admin.
func NewConnectHandler(svc *APIGatewayService) (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler(svc))
}
