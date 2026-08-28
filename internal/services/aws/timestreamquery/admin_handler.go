package timestreamquery

import (
	"context"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/timestreamquery"
	timestreamqueryconnect "vorpalstacks/internal/pb/aws/timestreamquery/timestreamqueryconnect"
)

// AdminHandler provides Timestream Query service administration functionality.
// It implements the TimestreamQueryServiceHandler interface for gRPC-Web
// communication. This file contains zero store package imports; all store
// access and proto conversion is delegated to
// admin_handler_convert.go.
type AdminHandler struct {
	timestreamqueryconnect.UnimplementedTimestreamQueryServiceHandler
	service *TimestreamQueryService
}

var _ timestreamqueryconnect.TimestreamQueryServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Timestream Query AdminHandler backed by the
// given service instance.
func NewAdminHandler(svc *TimestreamQueryService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

// ListScheduledQueries lists scheduled queries in Timestream Query.
func (h *AdminHandler) ListScheduledQueries(ctx context.Context, req *connect.Request[pb.ListScheduledQueriesRequest]) (*connect.Response[pb.ListScheduledQueriesResponse], error) {
	region := resolveRegion(req.Header())
	stores, err := h.getStores(region)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.ListScheduledQueriesCore(stores, ListScheduledQueriesInput{
		Region: region,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var summaries []*pb.ScheduledQuery
	for _, sq := range result.Summaries {
		summaries = append(summaries, toPbScheduledQuery(sq))
	}

	resp := &pb.ListScheduledQueriesResponse{
		Scheduledqueries: summaries,
	}
	return connect.NewResponse(resp), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Timestream
// Query admin console.
func NewConnectHandler(svc *TimestreamQueryService) (string, http.Handler) {
	return timestreamqueryconnect.NewTimestreamQueryServiceHandler(NewAdminHandler(svc))
}
