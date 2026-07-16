package timestreamquery

import (
	"context"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/timestreamquery"
	timestreamqueryconnect "vorpalstacks/internal/pb/aws/timestreamquery/timestreamqueryconnect"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
)

// AdminHandler provides Timestream Query service administration functionality.
// It implements the TimestreamQueryServiceHandler interface for gRPC-Web communication.
// It delegates to the shared TimestreamQueryService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
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

// getScheduledQueryStore retrieves the Timestream Query scheduled query store
// for the request region from the shared service cache.
func (h *AdminHandler) getScheduledQueryStore(req *connect.Request[pb.ListScheduledQueriesRequest]) (*timestreamstore.ScheduledQueryStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	return h.service.GetScheduledQueryStoreForRegion(region)
}

// ListScheduledQueries lists scheduled queries in Timestream Query.
func (h *AdminHandler) ListScheduledQueries(ctx context.Context, req *connect.Request[pb.ListScheduledQueriesRequest]) (*connect.Response[pb.ListScheduledQueriesResponse], error) {
	store, err := h.getScheduledQueryStore(req)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	queries, err := store.ListScheduledQueries()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var summaries []*pb.ScheduledQuery
	for _, sq := range queries {
		summary := &pb.ScheduledQuery{
			Arn:          sq.ARN,
			Name:         sq.Name,
			Creationtime: sq.CreationTime.Format(timeutils.ISO8601UTCFormat),
		}
		if !sq.PreviousRunTime.IsZero() {
			summary.Previousinvocationtime = sq.PreviousRunTime.Format(timeutils.ISO8601UTCFormat)
		}
		if !sq.NextRunTime.IsZero() {
			summary.Nextinvocationtime = sq.NextRunTime.Format(timeutils.ISO8601UTCFormat)
		}
		summaries = append(summaries, summary)
	}

	return connect.NewResponse(&pb.ListScheduledQueriesResponse{
		Scheduledqueries: summaries,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Timestream Query admin console.
func NewConnectHandler(svc *TimestreamQueryService) (string, http.Handler) {
	return timestreamqueryconnect.NewTimestreamQueryServiceHandler(NewAdminHandler(svc))
}
