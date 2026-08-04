package apigateway

import (
	"net/http"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
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

// storeErr maps a service-layer error to a connect gRPC error.
// Core functions may return raw store sentinel errors (e.g.
// ErrRestApiNotFound) that are not *AWSError values.  toApiGatewayError
// converts those sentinels into *ApiGatewayError via the
// storeErrorMappings table; then AWSErrorToGRPC performs the canonical
// HTTP-status-to-connect-code mapping, including 409 sub-code
// distinctions (DeleteConflict, LimitExceeded, *AlreadyExists).
func storeErr(err error) error {
	return svcerrors.AWSErrorToGRPC(toApiGatewayError(err))
}

// NewConnectHandler returns the connect RPC path and handler for API Gateway admin.
func NewConnectHandler(svc *APIGatewayService) (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler(svc))
}
